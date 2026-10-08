package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/yaml"

	"github.com/infrashift/technology-radar/schema"
)

// Entry is a configurable quadrant or ring with a URL-safe slug and display label.
type Entry struct {
	Slug  string `json:"slug"`
	Label string `json:"label"`
}

// Config holds the application configuration. Field names and defaults are
// defined by #Config in schema/radar.cue.
type Config struct {
	ContentDir    string  `json:"content_dir"`
	WebDir        string  `json:"web_dir"`
	OutputJSON    string  `json:"output_json"`
	OutputContent string  `json:"output_content"`
	Version       string  `json:"version"`
	Quadrants     []Entry `json:"quadrants"`
	Rings         []Entry `json:"rings"`
}

// QuadrantSlugs returns the slug strings from the config's quadrant definitions.
func (c Config) QuadrantSlugs() []string {
	return slugs(c.Quadrants)
}

// RingSlugs returns the slug strings from the config's ring definitions.
func (c Config) RingSlugs() []string {
	return slugs(c.Rings)
}

func slugs(entries []Entry) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.Slug
	}
	return out
}

// Load reads the YAML config at path, validates it against #Config, and fills
// in defaults. A missing file yields the defaults unless required is set.
func Load(path string, required bool) (Config, error) {
	ctx := cuecontext.New()
	s, err := schema.Compile(ctx)
	if err != nil {
		return Config{}, err
	}
	v := s.LookupPath(cue.ParsePath("config"))

	data, err := os.ReadFile(path)
	switch {
	case err == nil:
		f, err := yaml.Extract(path, data)
		if err != nil {
			return Config{}, fmt.Errorf("parsing %s: %w", path, err)
		}
		v = v.Unify(ctx.BuildFile(f))
	case errors.Is(err, fs.ErrNotExist) && !required:
		// No config file: use schema defaults.
	default:
		return Config{}, fmt.Errorf("reading config: %w", err)
	}

	if err := v.Validate(cue.Concrete(true)); err != nil {
		return Config{}, fmt.Errorf("invalid config %s:\n  %s", path, strings.Join(schema.FormatErrors(err), "\n  "))
	}

	var cfg Config
	if err := v.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("decoding config: %w", err)
	}
	return cfg, nil
}
