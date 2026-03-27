package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "radar",
	Short: "Technology Radar CLI - validate, build, and scaffold radar content",
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("config", "", "config file (default is .radar.yaml)")
	rootCmd.PersistentFlags().String("content-dir", "", "content directory")
	rootCmd.PersistentFlags().String("schema-path", "", "CUE schema path")
	rootCmd.PersistentFlags().String("web-dir", "", "web project directory")

	viper.BindPFlag("content_dir", rootCmd.PersistentFlags().Lookup("content-dir"))
	viper.BindPFlag("schema_path", rootCmd.PersistentFlags().Lookup("schema-path"))
	viper.BindPFlag("web_dir", rootCmd.PersistentFlags().Lookup("web-dir"))
}

func initConfig() {
	if cfgFile, _ := rootCmd.Flags().GetString("config"); cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName(".radar")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	viper.SetEnvPrefix("RADAR")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("content_dir", "content/blips")
	viper.SetDefault("schema_path", "schema/radar.cue")
	viper.SetDefault("web_dir", "web")
	viper.SetDefault("output_json", "web/src/data/blips-generated.json")
	viper.SetDefault("output_content", "web/src/content/blips")
	viper.SetDefault("version", "2025.11")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Warning: error reading config: %v\n", err)
		}
	}
}

func Execute() error {
	return rootCmd.Execute()
}
