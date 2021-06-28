package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/rrgc/rrgc"
	"moul.io/srand"
	"moul.io/zapconfig"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		}
		os.Exit(1)
	}
}

var opts struct {
	dryRun  bool
	debug   bool
	verbose bool
	logger  *zap.Logger
}

func run(args []string) error {
	// setup CLI
	rootFs := flag.NewFlagSet("rrgc", flag.ExitOnError)
	rootFs.BoolVar(&opts.dryRun, "dry-run", opts.dryRun, "dry-run")
	rootFs.BoolVar(&opts.debug, "debug", opts.debug, "debug")
	rootFs.BoolVar(&opts.verbose, "verbose", opts.verbose, "verbose")
	root := &ffcli.Command{
		Name:       "rrgc",
		FlagSet:    rootFs,
		ShortUsage: "rrgc WINDOWS -- GLOBS",
		Exec:       doRoot,
	}
	if err := root.Parse(args); err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	// init logger
	{
		rand.Seed(srand.Fast())
		config := zapconfig.New().SetPreset("light-console")
		if opts.debug {
			config = config.SetLevel(zapcore.DebugLevel)
		} else {
			config = config.SetLevel(zapcore.InfoLevel)
		}
		var err error
		opts.logger, err = config.Build()
		if err != nil {
			return fmt.Errorf("logger init: %w", err)
		}
	}

	// run
	{
		if err := root.Run(context.Background()); err != nil {
			return fmt.Errorf("run error: %w", err)
		}
	}

	return nil
}

func doRoot(ctx context.Context, args []string) error {
	var (
		windows = make([]rrgc.Window, 0)
		globs   = make([]string, 0)
	)
	modeGlob := false
	for _, arg := range args {
		switch {
		case !modeGlob && arg == "--":
			modeGlob = true
		case modeGlob && arg == "--":
			return flag.ErrHelp
		case !modeGlob:
			window, err := rrgc.ParseWindow(arg)
			if err != nil {
				return fmt.Errorf("parse window: %q: %w", window, err)
			}
			windows = append(windows, window)
		default:
			// FIXME: check glob input?
			globs = append(globs, arg)
		}
	}
	if len(windows) == 0 || len(globs) == 0 {
		return flag.ErrHelp
	}
	opts.logger.Debug(
		"args",
		zap.Any("windows", windows),
		zap.Strings("globs", globs),
	)

	toDelete, err := rrgc.GCListByPathGlobs(globs, windows)
	if err != nil {
		return fmt.Errorf("compute GC list: %w", err)
	}
	opts.logger.Debug(
		"to delete",
		zap.Strings("paths", toDelete),
		zap.Bool("dry-run", opts.dryRun),
		zap.Bool("verbose", opts.verbose),
	)
	var errs error
	for _, path := range toDelete {
		if opts.dryRun || opts.verbose {
			fmt.Printf("rm %q\n", path)
		}
		if !opts.dryRun {
			err := os.Remove(path)
			if err != nil {
				errs = multierr.Append(errs, fmt.Errorf("delete %q: %w", path, err))
			}
		}
	}
	return errs
}
