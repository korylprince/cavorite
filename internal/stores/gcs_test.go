package stores

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/fsouza/fake-gcs-server/fakestorage"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

type mapFile struct {
	content []byte
	modTime *time.Time
}

func memMapFsWith(files map[string]mapFile) (*afero.Fs, error) {
	memfsys := afero.NewMemMapFs()
	for fname, mfile := range files {
		afile, err := memfsys.Create(fname)
		if err != nil {
			return nil, err
		}
		_, err = afile.Write(mfile.content)
		if err != nil {
			return nil, err
		}
		if mfile.modTime != nil {
			err := memfsys.Chtimes(fname, time.Time{}, *mfile.modTime)
			if err != nil {
				return nil, err
			}
		}
	}
	return &memfsys, nil
}

func fakeGCSStore(t *testing.T) GCSStore {
	// create fake gcs
	server, err := fakestorage.NewServerWithOptions(fakestorage.Options{
		InitialObjects: []fakestorage.Object{},
	})
	server.CreateBucketWithOpts(fakestorage.CreateBucketOpts{
		Name: "test",
	})
	assert.NoError(t, err)

	t.Log(server.URL())
	client, err := storage.NewClient(
		context.Background(),
		option.WithoutAuthentication(),
		option.WithHTTPClient(server.HTTPClient()),
	)
	assert.NoError(t, err)

	mTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2014-11-12T11:45:26.371Z")
	memfs, err := memMapFsWith(map[string]mapFile{
		"/thing": {
			content: []byte(`blah`),
			modTime: &mTime,
		},
	})
	assert.NoError(t, err)

	gcs := &GCSStore{
		fsys:      *memfs,
		gcsClient: client,
		Options: Options{
			BackendAddress: "test",
		},
	}
	return *gcs
}

func TestGCSUpload(t *testing.T) {
	store := fakeGCSStore(t)
	err := store.Upload(context.Background(), "/thing")
	assert.NoError(t, err)
}
