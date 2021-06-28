package rrgc_test

import (
	"os"
	"time"

	"moul.io/rrgc/rrgc"
)

func Example() {
	logGlobs := []string{
		"*/*.log",
		"*/*.log.gz",
	}
	windows := []rrgc.Window{
		{Every: 2 * time.Hour, MaxKeep: 5},
		{Every: time.Hour * 24, MaxKeep: 4},
		{Every: time.Hour * 24 * 7, MaxKeep: 3},
	}
	toDelete, _ := rrgc.GCListByPathGlobs(logGlobs, windows)
	for _, path := range toDelete {
		_ = os.Remove(path)
	}
}
