package rrgc

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterFilesByWindow(t *testing.T) {
	cases := []struct {
		name         string
		files        []file
		window       Window
		expectedKeep string
		expectedDrop string
		expectError  bool
	}{
		{
			"hourly5",
			testingFiles(),
			Window{Every: time.Hour, MaxKeep: 5},
			"ADFGH",
			"BCEIJKLMNOPQRSTUVWXYZ",
			false,
		},
		{
			"secondly5",
			testingFiles(),
			Window{Every: time.Second, MaxKeep: 5},
			"ABCDE",
			"FGHIJKLMNOPQRSTUVWXYZ",
			false,
		},
		{
			"daily5",
			testingFiles(),
			Window{Every: time.Hour * 24, MaxKeep: 5},
			"ANOPQ",
			"BCDEFGHIJKLMRSTUVWXYZ",
			false,
		},
		{
			"yearly5",
			testingFiles(),
			Window{Every: time.Hour * 24 * 365, MaxKeep: 5},
			"A",
			"BCDEFGHIJKLMNOPQRSTUVWXYZ",
			false,
		},
		{
			"all-in",
			testingFiles(),
			Window{Every: time.Second, MaxKeep: 500},
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			"",
			false,
		},
		{
			"nothing",
			testingFiles(),
			Window{Every: time.Second, MaxKeep: 0},
			"",
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			false,
		},
		// FIXME: test same date
		// FIXME: test invalid value (negative count?)
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			keep, drop, err := filterFilesByWindow(tc.files, tc.window)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				keepStr := strings.Join(keep, "")
				dropStr := strings.Join(drop, "")
				assert.Equal(t, keepStr, tc.expectedKeep)
				assert.Equal(t, dropStr, tc.expectedDrop)
			}
		})
	}
}

func TestFilterFilesByWindows(t *testing.T) {
	cases := []struct {
		name         string
		files        []file
		windows      []Window
		expectedKeep string
		expectedDrop string
		expectError  bool
	}{
		{
			"hourly5",
			testingFiles(),
			[]Window{{Every: time.Hour, MaxKeep: 5}},
			"ADFGH",
			"BCEIJKLMNOPQRSTUVWXYZ",
			false,
		},
		{
			"secondly5",
			testingFiles(),
			[]Window{{Every: time.Second, MaxKeep: 5}},
			"ABCDE",
			"FGHIJKLMNOPQRSTUVWXYZ",
			false,
		},
		{
			"daily5",
			testingFiles(),
			[]Window{{Every: time.Hour * 24, MaxKeep: 5}},
			"ANOPQ",
			"BCDEFGHIJKLMRSTUVWXYZ",
			false,
		},
		{
			"yearly5",
			testingFiles(),
			[]Window{{Every: time.Hour * 24 * 365, MaxKeep: 5}},
			"A",
			"BCDEFGHIJKLMNOPQRSTUVWXYZ",
			false,
		},
		{
			"all-in",
			testingFiles(),
			[]Window{{Every: time.Second, MaxKeep: 500}},
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			"",
			false,
		},
		{
			"nothing",
			testingFiles(),
			[]Window{{Every: time.Second, MaxKeep: 0}},
			"",
			"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			false,
		},
		{
			"combined-hourly5-hourly5",
			testingFiles(),
			[]Window{{Every: time.Hour, MaxKeep: 5}, {Every: time.Hour, MaxKeep: 5}},
			"ADFGH",
			"BCEIJKLMNOPQRSTUVWXYZ",
			false,
		},
		{
			"combined-hourly5-daily5",
			testingFiles(),
			[]Window{{Every: time.Hour, MaxKeep: 5}, {Every: time.Hour * 24, MaxKeep: 5}},
			"ADFGHNOPQ",
			"BCEIJKLMRSTUVWXYZ",
			false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			keep, drop, err := filterFilesByWindows(tc.files, tc.windows)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				keepStr := strings.Join(keep, "")
				dropStr := strings.Join(drop, "")
				assert.Equal(t, keepStr, tc.expectedKeep)
				assert.Equal(t, dropStr, tc.expectedDrop)
			}
		})
	}
}

func mustTimeParse(input string) time.Time {
	parsed, err := time.Parse(time.RFC3339, input)
	if err != nil {
		panic(err)
	}
	return parsed
}

func testingFiles() []file {
	return []file{
		{path: "A", time: mustTimeParse("2006-01-02T15:04:05Z")},
		{path: "B", time: mustTimeParse("2006-01-02T15:34:05Z")},
		{path: "C", time: mustTimeParse("2006-01-02T16:03:05Z")},
		{path: "D", time: mustTimeParse("2006-01-02T16:05:05Z")},
		{path: "E", time: mustTimeParse("2006-01-02T17:04:05Z")},
		{path: "F", time: mustTimeParse("2006-01-02T18:04:05Z")},
		{path: "G", time: mustTimeParse("2006-01-02T19:04:05Z")},
		{path: "H", time: mustTimeParse("2006-01-02T20:04:05Z")},
		{path: "I", time: mustTimeParse("2006-01-02T21:04:05Z")},
		{path: "J", time: mustTimeParse("2006-01-03T18:04:05Z")},
		{path: "K", time: mustTimeParse("2006-01-03T18:05:05Z")},
		{path: "L", time: mustTimeParse("2006-01-03T18:06:05Z")},
		{path: "M", time: mustTimeParse("2006-01-03T17:04:05Z")},
		{path: "N", time: mustTimeParse("2006-01-03T16:04:05Z")},
		{path: "O", time: mustTimeParse("2006-01-04T16:04:05Z")},
		{path: "P", time: mustTimeParse("2006-01-05T16:04:05Z")},
		{path: "Q", time: mustTimeParse("2006-01-06T16:04:05Z")},
		{path: "R", time: mustTimeParse("2006-01-07T16:04:05Z")},
		{path: "S", time: mustTimeParse("2006-01-08T16:04:05Z")},
		{path: "T", time: mustTimeParse("2006-01-09T16:04:05Z")},
		{path: "U", time: mustTimeParse("2006-01-10T16:04:05Z")},
		{path: "V", time: mustTimeParse("2006-01-15T16:04:05Z")},
		{path: "W", time: mustTimeParse("2006-01-18T16:04:05Z")},
		{path: "X", time: mustTimeParse("2006-01-20T16:04:05Z")},
		{path: "Y", time: mustTimeParse("2006-01-25T16:04:05Z")},
		{path: "Z", time: mustTimeParse("2006-01-30T16:04:05Z")},
	}
}
