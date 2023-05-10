package config

import (
	"context"
	"errors"
	"testing"

	"github.com/discentem/pantri_but_go/internal/stores"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestValidateFuncNil(t *testing.T) {
	conf := Config{
		Validate: nil,
	}
	err := conf.Write(afero.NewMemMapFs(), "")
	assert.ErrorIs(t, err, ErrValidateNil)
}
func TestValidateFails(t *testing.T) {
	conf := Config{
		Validate: func() error {
			return errors.New("failed")
		},
	}
	err := conf.Write(afero.NewMemMapFs(), "")
	assert.ErrorIs(t, err, ErrValidate)
}

func TestDirExpanderFails(t *testing.T) {
	conf := Config{
		Validate: func() error { return nil },
	}
	expander = func(path string) (string, error) {
		return "", errors.New("borked")
	}
	err := conf.Write(afero.NewMemMapFs(), "")
	assert.ErrorIs(t, err, ErrDirExpander)
}

func TestSuccessfulWrite(t *testing.T) {
	cfg := Config{
		StoreType: stores.StoreTypeS3,
		Options: stores.Options{
			PantriAddress:         "s3://blahaddress/bucket",
			MetaDataFileExtension: "",
			Region:                "us-east-9876",
		},
		Validate: func() error { return nil },
	}
	// override back to a dirExpander that will succeed, as opposed to previous test
	expander = func(path string) (string, error) {
		return path, nil
	}
	fsys := afero.NewMemMapFs()
	err := cfg.Write(fsys, ".")
	assert.NoError(t, err)
	f, err := fsys.Open(".pantri/config")
	assert.NoError(t, err)
	b, err := afero.ReadAll(f)
	assert.NoError(t, err)

	expected := `{
 "store_type": "s3",
 "options": {
  "pantri_address": "s3://blahaddress/bucket",
  "metadata_file_extension": "",
  "region": "us-east-9876"
 }
}`
	assert.Equal(t, expected, string(b))

}

func TestWrite(t *testing.T) {
	t.Run("fail if validate() nil", TestValidateFuncNil)
	t.Run("fail if validate() fails", TestValidateFails)
	t.Run("fail if dirExpander() fails",
		TestDirExpanderFails,
	)
	t.Run("successful write", TestSuccessfulWrite)
}

func TestInitializeStoreTypeS3(t *testing.T) {
	ctx := context.Background()
	fsys := afero.NewMemMapFs()

	opts := stores.Options{
		PantriAddress:         "s3://blahaddress/bucket",
		MetaDataFileExtension: "pfile",
		Region:                "us-east-9876",
	}

	cfg := InitializeStoreTypeS3(
		ctx,
		fsys,
		"~/some_repo_root",
		opts.PantriAddress,
		opts.Region,
		opts,
	)

	// Assert the S3Store Config matches all of the inputs
	assert.Equal(t, cfg.StoreType, stores.StoreTypeS3)
	assert.Equal(t, cfg.Options.PantriAddress, opts.PantriAddress)
	assert.Equal(t, cfg.Options.MetaDataFileExtension, opts.MetaDataFileExtension)
	assert.Equal(t, cfg.Options.Region, opts.Region)

	assert.NoError(t, cfg.Validate())
}

func TestInitializeStoreTypeGCS(t *testing.T) {
	ctx := context.Background()
	fsys := afero.NewMemMapFs()

	opts := stores.Options{
		PantriAddress:         "my-test-bucket",
		MetaDataFileExtension: "pfile",
		Region:                "us-east-9876",
	}

	cfg := InitializeStoreTypeGCS(
		ctx,
		fsys,
		"~/some_repo_root",
		opts.PantriAddress,
		opts.Region,
		opts,
	)

	// Assert the S3Store Config matches all of the inputs
	assert.Equal(t, cfg.StoreType, stores.StoreTypeGCS)
	assert.Equal(t, cfg.Options.PantriAddress, opts.PantriAddress)
	assert.Equal(t, cfg.Options.MetaDataFileExtension, opts.MetaDataFileExtension)
	assert.Equal(t, cfg.Options.Region, opts.Region)

	assert.NoError(t, cfg.Validate())
}
