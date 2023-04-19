package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/google/logger"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
)

type Config struct {
	PantriAddress         string `json:"pantri_address" mapstructure:"pantri_address"`
	MetaDataFileExtension string `json:"metadata_file_extension" mapstructure:"metadata_file_extension"`
	// TODO(discentem) remove this option. See #15
	RemoveFromSourceRepo *bool        `json:"remove_from_sourcerepo" mapstructure:"remove_from_sourcerepo"`
	Type                 string       `json:"type"`
	Validate             func() error `json:"-"`
}

type dirExpanderer func(string) (string, error)

// dirExpander can be overwritten for tests
var dirExpander dirExpanderer = homedir.Expand

var ErrValidateNil = errors.New("pantri config must have a Validate() function")
var ErrValidate = errors.New("validate() failed")
var ErrDirExpander = errors.New("dirExpander failed")

func (c *Config) Write(fsys afero.Fs, sourceRepo string) error {
	if c.Validate == nil {
		return ErrValidateNil
	}
	if err := c.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrValidate, err)
	}
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}
	esr, err := dirExpander(sourceRepo)
	if err != nil {
		return fmt.Errorf(
			"%w: %v",
			ErrDirExpander,
			err,
		)
	}
	cfile := filepath.Join(esr, ".pantri/config")
	if _, err := fsys.Stat(esr); err != nil {
		return fmt.Errorf("%s does not exist, so we can't make it a pantri repo", esr)
	}

	if err := fsys.MkdirAll(filepath.Dir(cfile), os.ModePerm); err != nil {
		return err
	}
	_, err = fsys.Stat(cfile)
	if errors.Is(err, fs.ErrNotExist) {
		logger.Infof("initialized pantri config at %s", cfile)

	} else {
		logger.Infof("updating pantri config at %s", cfile)
	}
	return afero.WriteFile(fsys, cfile, b, os.ModePerm)

}
