package stores

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/logger"
	"github.com/spf13/afero"

	"github.com/discentem/cavorite/internal/metadata"

	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

type S3Downloader interface {
	Download(
		ctx context.Context,
		w io.WriterAt,
		input *s3.GetObjectInput,
		options ...func(*s3manager.Downloader)) (n int64, err error)
}

type S3Uploader interface {
	Upload(ctx context.Context,
		input *s3.PutObjectInput,
		opts ...func(*s3manager.Uploader)) (
		*s3manager.UploadOutput, error,
	)
}

type S3Store struct {
	Options      Options `json:"options" mapstructure:"options"`
	fsys         afero.Fs
	awsRegion    string
	s3Uploader   S3Uploader
	s3Downloader S3Downloader
}

func NewS3StoreClient(ctx context.Context, fsys afero.Fs, opts Options) (*S3Store, error) {
	cfg, err := getConfig(
		ctx,
		opts.Region,
		opts.BackendAddress,
	)
	if err != nil {
		return nil, err
	}

	// TODO (@radsec) - extract this S3 logic to a separate internal client instead of directly from AWS_SDK
	s3Client := s3.NewFromConfig(*cfg)
	s3Uploader := s3manager.NewUploader(
		s3Client,
		func(u *s3manager.Uploader) {
			u.PartSize = 64 * 1024 * 1024 // 64MB per part
		},
	)
	s3Downloader := s3manager.NewDownloader(
		s3Client,
		func(d *s3manager.Downloader) {
			d.Concurrency = 3
		},
	)

	return &S3Store{
		Options:   opts,
		fsys:      fsys,
		awsRegion: opts.Region,
		// s3Uploader meets our interface for S3Uploader
		s3Uploader: s3Uploader,
		// s3Downloader meets our interface for S3Downloader
		s3Downloader: s3Downloader,
	}, nil
}

func getConfig(ctx context.Context, region string, address string) (*aws.Config, error) {
	var cfg aws.Config
	var err error

	switch {
	case strings.HasPrefix(address, "s3://"):
		cfg, err = awsConfig.LoadDefaultConfig(
			ctx,
			config.WithRegion(region),
		)
		if err != nil {
			return nil, err
		}
	case strings.HasPrefix(address, "http://"):
		fallthrough
	case strings.HasPrefix(address, "https://"):
		server, _ := path.Split(address)
		// https://stackoverflow.com/questions/67575681/is-aws-go-sdk-v2-integrated-with-local-minio-server
		resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               server,
				SigningRegion:     region,
				HostnameImmutable: true,
			}, nil
		})

		cfg, err = config.LoadDefaultConfig(
			ctx,
			config.WithRegion(region),
			config.WithEndpointResolverWithOptions(resolver),
		)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("address did not contain s3://, http://, or https:// prefix")
	}

	return &cfg, nil
}

func (s *S3Store) GetOptions() Options {
	return s.Options
}

func (s *S3Store) GetFsys() afero.Fs {
	return s.fsys
}

// TODO(discentem): #34 largely copy-pasted from stores/local/local.go. Can be consolidated
// Upload generates the metadata, writes it to disk and uploads the file to the S3 bucket
func (s *S3Store) Upload(ctx context.Context, objects ...string) error {
	for _, o := range objects {
		fi, err := s.fsys.Open(o)
		if err != nil {
			return err
		}
		cleanup, err := WriteMetadataToFsys(s, o, fi)
		if err != nil {
			return err
		}
		_, err = fi.Seek(0, io.SeekStart)
		if err != nil {
			if err := cleanup(); err != nil {
				return err
			}
			return err
		}

		// Generate S3 struct for object and upload to S3 bucket
		_, buck := path.Split(s.Options.BackendAddress)
		obj := s3.PutObjectInput{
			Bucket: aws.String(buck),
			Key:    &o,
			Body:   fi,
		}
		out, err := s.s3Uploader.Upload(ctx, &obj)
		if err != nil {
			if err := cleanup(); err != nil {
				return fmt.Errorf("cleanup() failed after Upload failure: %w", err)
			}
			logger.Error(out)
			return err
		}
		if err := fi.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Retrieve gets the file from the S3 bucket, validates the hash is correct and writes it to disk
func (s *S3Store) Retrieve(ctx context.Context, objects ...string) error {
	for _, o := range objects {
		// For Retrieve, the object is the cfile itself, which we derive the actual filename from
		objectPath := strings.TrimSuffix(o, filepath.Ext(o))
		// We will either read the file that already exists or download it because it
		// is missing
		f, err := openOrCreateFile(s.fsys, objectPath)
		if err != nil {
			return err
		}
		fileInfo, err := f.Stat()
		if err != nil {
			return err
		}
		if fileInfo.Size() > 0 {
			logger.Infof("%s already exists", objectPath)
		} else { // Create an S3 struct for the file to be retrieved
			_, buck := path.Split(s.Options.BackendAddress)
			obj := &s3.GetObjectInput{
				Bucket: aws.String(buck),
				Key:    aws.String(objectPath),
			}
			// Download the file
			_, err := s.s3Downloader.Download(ctx, f, obj)
			if err != nil {
				return err
			}
		}
		// Get the hash for the downloaded file
		hash, err := metadata.SHA256FromReader(f)
		if err != nil {
			return err
		}
		// Get the metadata from the metadata file
		m, err := metadata.ParseCfile(s.fsys, o)
		if err != nil {
			return err
		}
		// If the hash of the downloaded file does not match the retrieved file, return an error
		if hash != m.Checksum {
			logger.Infof("Hash mismatch, got %s but expected %s", hash, m.Checksum)
			if err := s.fsys.Remove(objectPath); err != nil {
				return err
			}
			return ErrRetrieveFailureHashMismatch
		}
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}
