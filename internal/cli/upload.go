package cli

import (
	"errors"
	"fmt"

	"github.com/google/logger"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/discentem/cavorite/internal/program"
)

func uploadCmd() *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: fmt.Sprintf("Upload a file to %s", program.Name),
		Long:  fmt.Sprintf("Upload a file to %s", program.Name),
		Args:  cobra.MinimumNArgs(1),
		// PersistentPreRunE
		// Loads the config with OsFs
		/*
			// The *Run functions are executed in the following order:
			//   * PersistentPreRunE()
			//   * PreRunE() [X]
			//   * RunE()
			//   * PostRunE()
			//   * PersistentPostRunE()
			// All functions get the same args, the arguments after the command name.
		*/
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return loadConfig(afero.NewOsFs())
		},
		RunE: uploadFn,
	}

	return uploadCmd
}

// uploadFn is the execution runtime for the uploadCmd functionality
// in Cobra, this is the RunE phase
/*
	// The *Run functions are executed in the following order:
	//   * PersistentPreRunE()
	//   * PreRunE()
	//   * RunE() [X]
	//   * PostRunE()
	//   * PersistentPostRunE()
	// All functions get the same args, the arguments after the command name.
*/
func uploadFn(cmd *cobra.Command, objects []string) error {
	fsys := afero.NewOsFs()
	s, err := initStoreFromConfig(
		cmd.Context(),
		cfg,
		fsys,
		cfg.Options,
	)
	if err != nil {
		return err
	}
	defer s.Close()
	sourceRepoRoot, err := rootOfSourceRepo()
	if err != nil {
		return err
	}
	if sourceRepoRoot == nil {
		return errors.New("sourceRepoRoot cannot be nil")
	}

	// We need to remove the prefix from the path so it is relative
	objects, err = removePathPrefix(objects, *sourceRepoRoot)
	if err != nil {
		return fmt.Errorf("upload error: %w", err)
	}

	logger.Info("before GetOptions")
	opts, err := s.GetOptions()
	if err != nil {
		return err
	}
	logger.Info(opts)
	logger.Info("after GetOptions")

	logger.Infof("Uploading to: %s", opts.BackendAddress)
	logger.Infof("Uploading file: %s", objects)
	if err := s.Upload(cmd.Context(), objects...); err != nil {
		return err
	}
	return nil
}
