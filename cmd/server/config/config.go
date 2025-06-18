package config

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

type (
	Config struct {
		fs   *flag.FlagSet
		Host string
		Port string
	}
	LookUpEnv func(string) (string, bool)
)

func New(out io.Writer, lookUpEnv LookUpEnv, args ...string) (*Config, error) {
	fs := flag.NewFlagSet("wordle-cheater", flag.ExitOnError)
	fs.SetOutput(out)
	fs.Usage = func() {
		w := fs.Output()
		fmt.Fprintln(w, "runs site")
		fs.PrintDefaults()
	}

	var cfg Config
	fs.StringVar(&cfg.Host, "host", "", "the server to run on (usually leave empty)")
	fs.StringVar(&cfg.Port, "port", "8000", "the port to run on (required)")
	cfg.fs = fs

	if err := cfg.parse(lookUpEnv, args...); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return &cfg, nil
}

func (cfg *Config) parse(lookUpEnv LookUpEnv, args ...string) error {
	if err := cfg.fs.Parse(args); err != nil {
		return fmt.Errorf("parsing program args: %w", err)
	}
	if err := cfg.parseEnvVars(lookUpEnv); err != nil {
		return fmt.Errorf("setting value from environment variable: %w", err)
	}
	return nil
}

func (cfg *Config) parseEnvVars(lookUpEnv LookUpEnv) error {
	var lastErr error
	cfg.fs.VisitAll(func(f *flag.Flag) {
		upperName := strings.ToUpper(f.Name)
		name := strings.ReplaceAll(upperName, "-", "_")
		val, ok := lookUpEnv(name)
		if !ok {
			return
		}
		if err := f.Value.Set(val); err != nil {
			lastErr = err
		}
	})
	return lastErr
}
