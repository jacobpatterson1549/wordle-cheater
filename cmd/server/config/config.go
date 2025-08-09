package config

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Config struct {
	fs   *flag.FlagSet
	Host string
	Port string
}

func New() (*Config, error) {
	return newConfig(flag.ExitOnError, os.Stdout, os.Args...)
}

func newConfig(flagSetErrorHandling flag.ErrorHandling, out io.Writer, args ...string) (*Config, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("missing program name in args")
	}
	programName, args := args[0], args[1:]

	fs := flag.NewFlagSet(programName, flagSetErrorHandling)
	fs.SetOutput(out)
	fs.Usage = func() {
		w := fs.Output()
		fmt.Fprintln(w, "runs "+programName)
		fs.PrintDefaults()
	}

	var cfg Config
	fs.StringVar(&cfg.Host, "host", "", "the server to run on (usually leave empty)")
	fs.StringVar(&cfg.Port, "port", "8000", "the port to run on (required)")
	cfg.fs = fs

	if err := cfg.parse(args...); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return &cfg, nil
}

func (cfg *Config) parse(args ...string) error {
	if err := cfg.fs.Parse(args); err != nil {
		return fmt.Errorf("parsing program args: %w", err)
	}
	if err := cfg.parseEnv(); err != nil {
		return fmt.Errorf("setting value from environment variable: %w", err)
	}
	return nil
}

func (cfg *Config) parseEnv() error {
	var lastErr error
	cfg.fs.VisitAll(func(f *flag.Flag) {
		upperName := strings.ToUpper(f.Name)
		name := strings.ReplaceAll(upperName, "-", "_")
		val, ok := os.LookupEnv(name)
		if !ok {
			return
		}
		if err := f.Value.Set(val); err != nil {
			lastErr = err
		}
	})
	return lastErr
}
