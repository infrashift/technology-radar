package cli

import (
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Validate the config and blip content against the CUE schema",
		RunE:  runValidate,
	}
}

func runValidate(cmd *cobra.Command, args []string) error {
	cfg, err := loadConfig(cmd)
	if err != nil {
		return err
	}
	_, _, err = loadBlips(cmd, cfg)
	return err
}
