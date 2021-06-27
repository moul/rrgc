package rrgc

import (
	"fmt"
	"sort"
	"time"
)

type Window struct {
	Every   time.Duration
	MaxKeep int
}

type file struct {
	path string
	time time.Time
}

func GCListByPathGlobs(inputs []string, windows []Window) ([]string, error) {
	files, err := fileListByPathGlobs(inputs)
	if err != nil {
		return nil, fmt.Errorf("file list by path globs: %w", err)
	}

	_, drop, err := filterFilesByWindows(files, windows)
	if err != nil {
		return nil, fmt.Errorf("filter files by windows: %w", err)
	}

	return drop, nil
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
	return nil, nil
}

func (f file) String() string {
	return fmt.Sprintf("%s: %q", f.path, f.time)
}
