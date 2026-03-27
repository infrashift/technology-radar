package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/infrashift/technology-radar/internal/build"
	"github.com/infrashift/technology-radar/internal/cuegen"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate blip content against the CUE schema",
	RunE:  runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	contentDir := viper.GetString("content_dir")
	schemaPath := viper.GetString("schema_path")
	version := viper.GetString("version")

	blips, _, err := build.ParseAllBlips(contentDir)
	if err != nil {
		return fmt.Errorf("parsing blips: %w", err)
	}
	fmt.Printf("Parsed %d blips from %s\n", len(blips), contentDir)

	if err := cuegen.Validate(schemaPath, blips, version); err != nil {
		return fmt.Errorf("validation failed:\n%w", err)
	}

	fmt.Println("All blips passed CUE schema validation")
	return nil
}
