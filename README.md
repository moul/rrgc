# rrgc

 🗑 round-robin garbage-collector

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/moul.io/rrgc)
[![License](https://img.shields.io/badge/license-Apache--2.0%20%2F%20MIT-%2397ca00.svg)](https://github.com/moul/rrgc/blob/main/COPYRIGHT)
[![GitHub release](https://img.shields.io/github/release/moul/rrgc.svg)](https://github.com/moul/rrgc/releases)
[![Docker Metrics](https://images.microbadger.com/badges/image/moul/rrgc.svg)](https://microbadger.com/images/moul/rrgc)
[![Made by Manfred Touron](https://img.shields.io/badge/made%20by-Manfred%20Touron-blue.svg?style=flat)](https://manfred.life/)
n
[![Go](https://github.com/moul/rrgc/workflows/Go/badge.svg)](https://github.com/moul/rrgc/actions?query=workflow%3AGo)
[![Release](https://github.com/moul/rrgc/workflows/Release/badge.svg)](https://github.com/moul/rrgc/actions?query=workflow%3ARelease)
[![PR](https://github.com/moul/rrgc/workflows/PR/badge.svg)](https://github.com/moul/rrgc/actions?query=workflow%3APR)
[![GolangCI](https://golangci.com/badges/github.com/moul/rrgc.svg)](https://golangci.com/r/github.com/moul/rrgc)
[![codecov](https://codecov.io/gh/moul/rrgc/branch/main/graph/badge.svg)](https://codecov.io/gh/moul/rrgc)
[![Go Report Card](https://goreportcard.com/badge/moul.io/rrgc)](https://goreportcard.com/report/moul.io/rrgc)
[![CodeFactor](https://www.codefactor.io/repository/github/moul/rrgc/badge)](https://www.codefactor.io/repository/github/moul/rrgc)

[![Gitpod ready-to-code](https://img.shields.io/badge/Gitpod-ready--to--code-blue?logo=gitpod)](https://gitpod.io/#https://github.com/moul/rrgc)

## Usage

### As a CLI tool

[embedmd]:# (.tmp/usage.txt console)
```console
foo@bar:~$ rrgc -h
USAGE
  rrgc WINDOWS -- GLOBS

FLAGS
  -debug false  debug
foo@bar:~$ ls logs
A.log
B.log
C.log
D.log
E.log
F.log
G.log
H.log
I.log
J.log
K.log
L.log
M.log
N.log
O.log
P.log
Q.log
R.log
foo@bar:~$ rrgc 24h,5 1h,5 -- ./logs/*.log | xargs rm -v
removed 'logs/B.log'
removed 'logs/C.log'
removed 'logs/E.log'
removed 'logs/H.log'
removed 'logs/M.log'
removed 'logs/N.log'
removed 'logs/O.log'
removed 'logs/P.log'
removed 'logs/Q.log'
removed 'logs/R.log'
foo@bar:~$ rrgc 24h,5 1h,5 -- ./logs/*.log
foo@bar:~$ ls logs
A.log
D.log
F.log
G.log
I.log
J.log
K.log
L.log
```

### As a Library

[embedmd]:# (rrgc/example_test.go /import\ / $)
```go
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
```

[embedmd]:# (.tmp/godoc.txt txt /FUNCTIONS/ $)
```txt
FUNCTIONS

func GCListByPathGlobs(inputs []string, windows []Window) ([]string, error)
    GCListByPathGlobs computes a list of paths that should be deleted, based on
    a list of windows.


TYPES

type Window struct {
	Every   time.Duration
	MaxKeep int
}
    Window defines a file preservation rule.

func ParseWindow(input string) (Window, error)
    ParseWindow parses a human-readable Window definition.

    Syntax: "Duration,MaxKeep".

    Examples: "1h,5" "1h2m3s,42".

```

## Install

### Using go

```sh
go get moul.io/rrgc
```

### Releases

See https://github.com/moul/rrgc/releases

## Contribute

![Contribute <3](https://raw.githubusercontent.com/moul/moul/main/contribute.gif)

I really welcome contributions.
Your input is the most precious material.
I'm well aware of that and I thank you in advance.
Everyone is encouraged to look at what they can do on their own scale;
no effort is too small.

Everything on contribution is sum up here: [CONTRIBUTING.md](./.github/CONTRIBUTING.md)

### Dev helpers

Pre-commit script for install: https://pre-commit.com

### Contributors ✨

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-2-orange.svg)](#contributors)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="http://manfred.life"><img src="https://avatars1.githubusercontent.com/u/94029?v=4" width="100px;" alt=""/><br /><sub><b>Manfred Touron</b></sub></a><br /><a href="#maintenance-moul" title="Maintenance">🚧</a> <a href="https://github.com/moul/rrgc/commits?author=moul" title="Documentation">📖</a> <a href="https://github.com/moul/rrgc/commits?author=moul" title="Tests">⚠️</a> <a href="https://github.com/moul/rrgc/commits?author=moul" title="Code">💻</a></td>
    <td align="center"><a href="https://manfred.life/moul-bot"><img src="https://avatars1.githubusercontent.com/u/41326314?v=4" width="100px;" alt=""/><br /><sub><b>moul-bot</b></sub></a><br /><a href="#maintenance-moul-bot" title="Maintenance">🚧</a></td>
  </tr>
</table>

<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors)
specification. Contributions of any kind welcome!

### Stargazers over time

[![Stargazers over time](https://starchart.cc/moul/rrgc.svg)](https://starchart.cc/moul/rrgc)

## License

© 2021   [Manfred Touron](https://manfred.life)

Licensed under the [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)
([`LICENSE-APACHE`](LICENSE-APACHE)) or the [MIT license](https://opensource.org/licenses/MIT)
([`LICENSE-MIT`](LICENSE-MIT)), at your option.
See the [`COPYRIGHT`](COPYRIGHT) file for more details.

`SPDX-License-Identifier: (Apache-2.0 OR MIT)`
