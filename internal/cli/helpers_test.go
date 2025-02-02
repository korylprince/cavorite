package cli

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/discentem/cavorite/internal/config"
	"github.com/discentem/cavorite/internal/metadata"
	"github.com/discentem/cavorite/internal/stores"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestRemovePathPrefix(t *testing.T) {
	pathPrefix, err := os.Getwd()
	assert.NoError(t, err)

	expectedRemovePathPrefixes := []string{"foo/foo.dmg", "bar/bar.pkg"}
	testRemovePathPrefixes, err := removePathPrefix(
		[]string{
			fmt.Sprintf("%s/%s", pathPrefix, "foo/foo.dmg"),
			fmt.Sprintf("%s/%s", pathPrefix, "bar/bar.pkg"),
		},
		pathPrefix,
	)
	assert.NoError(t, err)
	assert.Equal(t, expectedRemovePathPrefixes, testRemovePathPrefixes)
}

func TestInitStoreFromConfig(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		StoreType: stores.StoreTypeS3,
		Options: stores.Options{
			BackendAddress:        "s3://test-bucket",
			MetadataFileExtension: metadata.MetadataFileExtension,
			Region:                "us-east-9876",
		},
	}

	fsys := afero.NewMemMapFs()

	s, err := initStoreFromConfig(
		ctx,
		cfg,
		fsys,
		cfg.Options,
	)
	assert.NoError(t, err)

	opts, err := s.GetOptions()
	assert.NoError(t, err)

	// Test GetOptions from newly created S3Store
	assert.Equal(t, opts.BackendAddress, "s3://test-bucket")

	// Test if type is equal to stores.S3Store
	// to ensure sure the buildStores returns
	// the proper stores based on type
	// assert.Equal(
	// 	t,
	// 	reflect.TypeOf(stores.Store(&stores.S3Store{})).Elem(),
	// 	reflect.TypeOf(s).Elem(),
	// )
}

func TestInitStoreFromConfig_InvalidOptions(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		StoreType: stores.StoreTypeS3,
		Options: stores.Options{
			BackendAddress:        "s4://test-bucket",
			MetadataFileExtension: metadata.MetadataFileExtension,
			Region:                "us-east-9876",
		},
	}

	fsys := afero.NewMemMapFs()

	_, err := initStoreFromConfig(
		ctx,
		cfg,
		fsys,
		cfg.Options,
	)

	assert.Errorf(t, err, `improper stores.S3Client init: cavoriteAddress did not contain s3://, http://, or https:// prefix`)
}

func TestInitStoreFromConfig_InvalidateStoreType(t *testing.T) {
	ctx := context.Background()
	cfg := config.Config{
		StoreType: "s4",
		Options: stores.Options{
			BackendAddress:        "s4://test-bucket",
			MetadataFileExtension: metadata.MetadataFileExtension,
			Region:                "us-east-9876",
		},
	}

	fsys := afero.NewMemMapFs()

	_, err := initStoreFromConfig(
		ctx,
		cfg,
		fsys,
		cfg.Options,
	)

	assert.Errorf(t, err, "type %s is not supported", "s4")
}
