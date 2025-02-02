package stores

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/discentem/cavorite/internal/testutils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type aferoS3Server struct {
	buckets map[string]afero.Fs
}

func (s aferoS3Server) Upload(ctx context.Context,
	input *s3.PutObjectInput,
	opts ...func(*s3manager.Uploader)) (
	*s3manager.UploadOutput, error,
) {
	bucket := *input.Bucket
	// check if the bucket referenced in input exists
	_, ok := s.buckets[bucket]
	if !ok {
		return nil, fmt.Errorf("%s does not exist in this aferoS3Server", bucket)
	}
	inputBytes, err := io.ReadAll(input.Body)
	if err != nil {
		return nil, err
	}
	// create a filesystem for bucket referenced in input
	bucketfs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		*input.Key: {
			// write input body to bucketfs
			Content: inputBytes,
		},
	})
	if err != nil {
		return nil, err
	}
	// write bucketfs to associated bucket
	s.buckets[bucket] = *bucketfs
	// S3Store doesn't use UploadOutput, so in the test we don't either
	return nil, nil
}

func (s aferoS3Server) Download(
	ctx context.Context,
	w io.WriterAt,
	input *s3.GetObjectInput,
	options ...func(*s3manager.Downloader)) (n int64, err error) {

	bucket := *input.Bucket
	// check if the bucket referenced in input exists
	_, ok := s.buckets[bucket]
	if !ok {
		return 0, fmt.Errorf("%s does not exist in this aferoS3Server", bucket)
	}

	objectHandle, err := s.buckets[bucket].Open(*input.Key)
	if err != nil {
		return 0, fmt.Errorf("could not find %s in bucket %s: %w", *input.Key, bucket, err)
	}
	objInfo, err := objectHandle.Stat()
	if err != nil {
		return 0, err
	}
	b := make([]byte, objInfo.Size())
	_, err = objectHandle.Read(b)
	if err != nil {
		return 0, fmt.Errorf("failed to read bytes from objectHandle: %w", err)
	}
	nbw, err := w.WriteAt(b, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to write objectHandle bytes to w: %w", err)
	}
	return int64(nbw), nil

}

func TestS3StoreUpload(t *testing.T) {
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	memfs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"test": {
			Content: []byte("bla"),
			ModTime: &mTime,
		},
	})
	assert.NoError(t, err)

	fakeS3Server := aferoS3Server{
		buckets: map[string]afero.Fs{
			// create a bucket in our fake s3 server
			"test": afero.NewMemMapFs(),
		},
	}

	store := s3Store{
		Options: Options{
			BackendAddress:        "s3://test",
			MetadataFileExtension: "cfile",
		},
		fsys:         *memfs,
		awsRegion:    "us-east-1",
		s3Uploader:   fakeS3Server,
		s3Downloader: fakeS3Server,
	}

	err = store.Upload(context.Background(), "test")
	assert.NoError(t, err)

	b, _ := afero.ReadFile(*memfs, "test.cfile")
	assert.Equal(t, `{
 "name": "test",
 "checksum": "4df3c3f68fcc83b27e9d42c90431a72499f17875c81a599b566c9889b9696703",
 "date_modified": "2014-11-12T11:45:26.371Z"
}`, string(b))
}

type s3FailServer struct{}

func (s s3FailServer) Download(
	ctx context.Context,
	w io.WriterAt,
	input *s3.GetObjectInput,
	options ...func(*s3manager.Downloader)) (n int64, err error) {
	return 0, errors.New("s3FailServer download failed")
}

func (s s3FailServer) Upload(ctx context.Context,
	input *s3.PutObjectInput,
	opts ...func(*s3manager.Uploader)) (
	*s3manager.UploadOutput, error,
) {
	return nil, errors.New("s3FailServer upload failed")
}

func TestS3StoreUploadFailCleanup(t *testing.T) {
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	memfs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"test": {
			Content: []byte("bla"),
			ModTime: &mTime,
		},
		// create .cfile to ensure when Upload fails that this gets cleaned up
		"test.cfile": {
			Content: []byte("bla"),
			ModTime: &mTime,
		},
	})
	assert.NoError(t, err)

	// intentionally broken Uploads
	S3ServerFail := s3FailServer{}

	store := s3Store{
		Options: Options{
			BackendAddress:        "s3://test",
			MetadataFileExtension: "cfile",
		},
		fsys:         *memfs,
		awsRegion:    "us-east-1",
		s3Uploader:   S3ServerFail,
		s3Downloader: S3ServerFail,
	}

	// Upload expected to failed
	err = store.Upload(context.Background(), "test")
	assert.Error(t, err)

	// Ensure this file doesn't exist to confirm cleanup() inside store.Upload succeeded
	_, err = afero.ReadFile(*memfs, "test.cfile")
	assert.Equal(t, true, os.IsNotExist(err))
}

func TestS3StoreRetrieve(t *testing.T) {
	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	// create bucket content
	bucketfs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"someObject": {
			Content: []byte("tla"),
			ModTime: &mTime,
		},
	})
	assert.NoError(t, err)

	fakeS3Server := aferoS3Server{
		buckets: map[string]afero.Fs{
			// create a bucket in our fake s3 server with the content
			"aFakeBucket": *bucketfs,
		},
	}

	localFs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"someObject.cfile": {
			Content: []byte(`{
				"name": "someObject",
				"checksum": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				"date_modified": "2014-11-12T11:45:26.371Z"
			   }`),
			ModTime: &mTime,
		},
	})
	assert.NoError(t, err)

	store := s3Store{
		Options: Options{
			BackendAddress:        "s3://aFakeBucket",
			MetadataFileExtension: "cfile",
		},
		fsys:         *localFs,
		awsRegion:    "us-east-1",
		s3Uploader:   fakeS3Server,
		s3Downloader: fakeS3Server,
	}

	err = store.Retrieve(context.Background(), "someObject.cfile")
	assert.NoError(t, err)

	// ensure the content of the file is correct
	b, _ := afero.ReadFile(*localFs, "someObject")
	assert.Equal(t, `tla`, string(b))

}

func TestS3GetBucketNameWithS3Prefix(t *testing.T) {
	expectedBackendAddress := "s3://aFakeBucket"

	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	memfs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"test": {
			Content: []byte("bla"),
			ModTime: &mTime,
		},
	})
	assert.NoError(t, err)

	fakeS3Server := aferoS3Server{
		buckets: map[string]afero.Fs{
			// create a bucket in our fake s3 server
			"test": afero.NewMemMapFs(),
		},
	}

	store := s3Store{
		Options: Options{
			BackendAddress:        expectedBackendAddress,
			MetadataFileExtension: "cfile",
		},
		fsys:         *memfs,
		awsRegion:    "us-east-1",
		s3Uploader:   fakeS3Server,
		s3Downloader: fakeS3Server,
	}

	bucketName, err := store.getBucketName()

	assert.NoError(t, err)
	assert.Equal(t, "aFakeBucket", bucketName)

}

func TestS3GetBucketNameWithHTTPPrefix(t *testing.T) {
	expectedBackendAddress := "http://127.0.0.1:9000/aFakeBucket"

	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	memfs, err := testutils.MemMapFsWith(map[string]testutils.MapFile{
		"test": {
			Content: []byte("bla"),
			ModTime: &mTime,
		},
	})
	assert.NoError(t, err)

	fakeS3Server := aferoS3Server{
		buckets: map[string]afero.Fs{
			// create a bucket in our fake s3 server
			"test": afero.NewMemMapFs(),
		},
	}

	store := s3Store{
		Options: Options{
			BackendAddress:        expectedBackendAddress,
			MetadataFileExtension: "cfile",
		},
		fsys:         *memfs,
		awsRegion:    "us-east-1",
		s3Uploader:   fakeS3Server,
		s3Downloader: fakeS3Server,
	}

	bucketName, err := store.getBucketName()

	assert.NoError(t, err)
	assert.Equal(t, "aFakeBucket", bucketName)

}
