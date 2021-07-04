package rrgc

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.uber.org/multierr"
)

// Window defines a file preservation rule.
type Window struct {
	Every   time.Duration
	MaxKeep int
}

// ParseWindow parses a human-readable Window definition.
//
// Syntax: "Duration,MaxKeep".
//
// Examples: "1h,5" "1h2m3s,42".
func ParseWindow(input string) (Window, error) {
	parts := strings.Split(input, ",")
	if len(parts) != 2 { // nolint:gomnd
		return Window{}, fmt.Errorf("invalid window format: %q", input)
	}
	duration, err := time.ParseDuration(parts[0])
	if err != nil {
		return Window{}, fmt.Errorf("invalid window duration value: %q: %w", parts[0], err)
	}
	maxKeep, err := strconv.Atoi(parts[1])
	if err != nil {
		return Window{}, fmt.Errorf("invalid window max-keep value: %q: %w", parts[1], err)
	}
	if maxKeep < 0 {
		return Window{}, fmt.Errorf("negative window max keep: %q", maxKeep)
	}
	return Window{
		Every:   duration,
		MaxKeep: maxKeep,
	}, nil
}

// GCListByPathGlobs computes a list of paths that should be kept and deleted, based on a list of window rules.
func GCListByPathGlobs(inputs []string, windows []Window) ([]string, []string, error) {
	files, err := fileListByPathGlobs(inputs)
	if err != nil {
		return nil, nil, fmt.Errorf("file list by path globs: %w", err)
	}

	keep, drop, err := filterFilesByWindows(files, windows)
	if err != nil {
		return nil, nil, fmt.Errorf("filter files by windows: %w", err)
	}

	return keep, drop, nil
}

type file struct {
	path string
	time time.Time
}

func filterFilesByWindows(files []file, windows []Window) ([]string, []string, error) {
	combinedKeep := make(map[string]bool)
	for _, window := range windows {
		keep, _, err := filterFilesByWindow(files, window)
		if err != nil {
			return nil, nil, fmt.Errorf("filter files by window: %w", err)
		}
		for _, file := range keep {
			combinedKeep[file] = true
		}
	}

	keep := make([]string, 0)
	drop := make([]string, 0)
	for _, file := range files {
		if _, found := combinedKeep[file.path]; found {
			keep = append(keep, file.path)
		} else {
			drop = append(drop, file.path)
		}
	}
	sort.Strings(keep)
	sort.Strings(drop)
	return keep, drop, nil
}

func filterFilesByWindow(files []file, window Window) ([]string, []string, error) {
	if window.MaxKeep < 0 {
		return nil, nil, fmt.Errorf("negative window max-keep: %q", window.MaxKeep)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].time.Before(files[j].time)
	})

	keep := []string{}
	drop := []string{}
	var previous time.Time
	for _, file := range files {
		if len(keep) >= window.MaxKeep {
			drop = append(drop, file.path)
			continue
		}
		if len(keep) == 0 {
			keep = append(keep, file.path)
			previous = file.time
			continue
		}
		if file.time.Sub(previous) < window.Every {
			drop = append(drop, file.path)
			continue
		}
		// else
		keep = append(keep, file.path)
		previous = file.time
	}
	sort.Strings(keep)
	sort.Strings(drop)
	return keep, drop, nil
}

func fileListByPathGlobs(inputs []string) ([]file, error) {
	uniquePaths := make(map[string]bool)
	var errs error
	for _, input := range inputs {
		matches, err := filepath.Glob(input)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}
		for _, match := range matches {
			clean := filepath.Clean(match)
			uniquePaths[clean] = true
		}
	}
	if errs != nil {
		return nil, fmt.Errorf("glob error: %w", errs)
	}

	paths := make([]string, len(uniquePaths))
	i := 0
	for path := range uniquePaths {
		paths[i] = path
		i++
	}
	sort.Strings(paths)

	files := make([]file, len(paths))
	for i, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("open %q: %w", path, err))
			continue
		}
		stat, err := f.Stat()
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("stat %q: %w", path, err))
			continue
		}
		if stat.IsDir() {
			errs = multierr.Append(errs, fmt.Errorf("is-dir %q", path))
			continue
		}
		files[i] = file{
			path: path,
			time: stat.ModTime(),
		}
	}
	if errs != nil {
		return nil, fmt.Errorf("analyze files: %w", errs)
	}

	return files, nil
}

func (f file) String() string {
	return fmt.Sprintf("%s: %q", f.path, f.time)
}
