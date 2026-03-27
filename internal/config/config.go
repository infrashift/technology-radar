package config

import "github.com/spf13/viper"

// Config holds the application configuration.
type Config struct {
	ContentDir    string `mapstructure:"content_dir"`
	SchemaPath    string `mapstructure:"schema_path"`
	WebDir        string `mapstructure:"web_dir"`
	OutputJSON    string `mapstructure:"output_json"`
	OutputContent string `mapstructure:"output_content"`
	Version       string `mapstructure:"version"`
}

// Load returns a Config populated from Viper.
func Load() Config {
	return Config{
		ContentDir:    viper.GetString("content_dir"),
		SchemaPath:    viper.GetString("schema_path"),
		WebDir:        viper.GetString("web_dir"),
		OutputJSON:    viper.GetString("output_json"),
		OutputContent: viper.GetString("output_content"),
		Version:       viper.GetString("version"),
	}
}
