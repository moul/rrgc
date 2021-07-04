package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/climan"
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
	debug     bool
	printKeep bool
	logger    *zap.Logger
}

func run(args []string) error {
	// setup CLI
	root := &climan.Command{
		Name: "rrgc",
		FlagSetBuilder: func(fs *flag.FlagSet) {
			fs.BoolVar(&opts.debug, "debug", opts.debug, "debug")
			fs.BoolVar(&opts.printKeep, "keep", opts.printKeep, "print list of files to keep instead of files to delete")
		},
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
	// check if enough args
	var (
		rawWindows = []string{}
		windows    = make([]rrgc.Window, 0)
		globs      = make([]string, 0)

		modeGlob = false
	)
	for _, arg := range args {
		switch {
		case !modeGlob && arg == "--":
			modeGlob = true
		case modeGlob && arg == "--":
			return flag.ErrHelp
		case !modeGlob:
			rawWindows = append(rawWindows, arg)
		default:
			globs = append(globs, arg)
		}
	}
	if len(rawWindows) == 0 || len(globs) == 0 {
		return flag.ErrHelp
	}

	// parse input formats;
	for _, rawWindow := range rawWindows {
		window, err := rrgc.ParseWindow(rawWindow)
		if err != nil {
			return fmt.Errorf("parse window: %q: %w", rawWindow, err)
		}
		windows = append(windows, window)
	}
	// FIXME: check glob input?

	// perform rrgc with sanitized inputs
	opts.logger.Debug(
		"args",
		zap.Any("windows", windows),
		zap.Strings("globs", globs),
	)
	toKeep, toDelete, err := rrgc.GCListByPathGlobs(globs, windows)
	if err != nil {
		return fmt.Errorf("compute GC list: %w", err)
	}
	opts.logger.Debug("to keep", zap.Strings("paths", toKeep))
	opts.logger.Debug("to delete", zap.Strings("paths", toDelete))
	if opts.printKeep {
		for _, path := range toKeep {
			fmt.Println(path)
		}
	} else {
		for _, path := range toDelete {
			fmt.Println(path)
		}
	}
	return nil
}
