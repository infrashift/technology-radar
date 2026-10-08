package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/infrashift/technology-radar/internal/build"
	"github.com/infrashift/technology-radar/internal/config"
	"github.com/infrashift/technology-radar/internal/validate"
)

const defaultConfigFile = ".radar.yaml"

// newRootCmd builds a fresh command tree, so each invocation (and each test)
// starts with clean flag state.
func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:          "radar",
		Short:        "Technology Radar CLI - validate, build, and scaffold radar content",
		SilenceUsage: true,
	}

	pf := root.PersistentFlags()
	pf.String("config", defaultConfigFile, "config file")
	pf.String("content-dir", "", "content directory (overrides config)")
	pf.String("web-dir", "", "web project directory (overrides config)")

	root.AddCommand(newBuildCmd(), newValidateCmd(), newScaffoldCmd())
	return root
}

func Execute() error {
	return newRootCmd().Execute()
}

// loadConfig loads the config file named by --config, then applies any
// directory flags that were set explicitly.
func loadConfig(cmd *cobra.Command) (config.Config, error) {
	flags := cmd.Flags()
	path, _ := flags.GetString("config")
	cfg, err := config.Load(path, flags.Changed("config"))
	if err != nil {
		return config.Config{}, err
	}
	if flags.Changed("content-dir") {
		cfg.ContentDir, _ = flags.GetString("content-dir")
	}
	if flags.Changed("web-dir") {
		cfg.WebDir, _ = flags.GetString("web-dir")
	}
	return cfg, nil
}

// loadBlips reads and validates every blip in the configured content directory.
func loadBlips(cmd *cobra.Command, cfg config.Config) ([]build.RawBlipFile, []build.Blip, error) {
	files, err := build.ReadBlipFiles(cfg.ContentDir)
	if err != nil {
		return nil, nil, fmt.Errorf("reading blips: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Parsed %d blips from %s\n", len(files), cfg.ContentDir)

	blips, err := validate.Blips(cfg, files)
	if err != nil {
		return nil, nil, fmt.Errorf("validation failed: %w", err)
	}
	fmt.Fprintln(cmd.OutOrStdout(), "All blips passed validation")
	return files, blips, nil
}
