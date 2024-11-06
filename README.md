# Config Framework

[![Build][build-badge]][build-link]
[![Coverage][coveralls-badge]][coveralls-link]
[![Coverage][coverage-badge]][coverage-link]
[![Quality][quality-badge]][quality-link]
[![Report][report-badge]][report-link]
[![FOSSA][fossa-badge]][fossa-link]
[![License][license-badge]][license-link]
[![Docs][docs-badge]][docs-link]
<!--
[![Libraries][libs-badge]][libs-link]
[![Security][security-badge]][security-link]
-->

[build-badge]: https://github.com/tkrop/go-config/actions/workflows/build.yaml/badge.svg
[build-link]: https://github.com/tkrop/go-config/actions/workflows/build.yaml

[coveralls-badge]: https://coveralls.io/repos/github/tkrop/go-config/badge.svg?branch=main
[coveralls-link]: https://coveralls.io/github/tkrop/go-config?branch=main

[coverage-badge]: https://app.codacy.com/project/badge/Coverage/b2bb898346ae4bb4be6414cd6dfe4932
[coverage-link]: https://app.codacy.com/gh/tkrop/go-config/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage

[quality-badge]: https://app.codacy.com/project/badge/Grade/b2bb898346ae4bb4be6414cd6dfe4932
[quality-link]: https://app.codacy.com/gh/tkrop/go-config/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade

[report-badge]: https://goreportcard.com/badge/github.com/tkrop/go-config
[report-link]: https://goreportcard.com/report/github.com/tkrop/go-config

[fossa-badge]: https://app.fossa.com/api/projects/git%2Bgithub.com%2Ftkrop%2Fgo-config.svg?type=shield&issueType=license
[fossa-link]: https://app.fossa.com/projects/git%2Bgithub.com%2Ftkrop%2Fgo-config?ref=badge_shield&issueType=license

[license-badge]: https://img.shields.io/badge/License-MIT-yellow.svg
[license-link]: https://opensource.org/licenses/MIT

[docs-badge]: https://pkg.go.dev/badge/github.com/tkrop/go-config.svg
[docs-link]: https://pkg.go.dev/github.com/tkrop/go-config

<!--
[libs-badge]: https://img.shields.io/librariesio/release/github/tkrop/go-config
[libs-link]: https://libraries.io/github/tkrop/go-config

[security-badge]: https://snyk.io/test/github/tkrop/go-config/main/badge.svg
[security-link]: https://snyk.io/test/github/tkrop/go-config
-->

## Introduction

Goal of [`go-config`][go-config] is to provide an easy to use and extensible
config framework with fluent interface based on [Viper][viper] for services,
jobs, and commands. It is supporting simple `default`-tags and prototype
config to set up and change the reader defaults quickly.

[viper]: <https://github.com/spf13/viper>
[go-config]: <https://github.com/tkrop/go-config>


## How to start

In [`go-config`][go-config] you simply create your config as an extension of
the config provided in this package as a base line as follows:

```go
// Import for config prototype and config reader.
import "github.com/tkrop/go-config/config"

// Config root element for configuration.
type Config struct {
    config.Config `mapstructure:",squash"`

    Int int       `default:"31"`
    String string `default:"my-value"`
    Dur Duration  `default:"1m"`

    Service *my.ServiceConfig
}
```

**Note:**  [`go-config`][go-config] makes it very simple to reuse the simple
config `struct`s provided by other libraries and components, since you can
easily create any hierarchy of `struct`s, `slice`s, and even `map[string]`s
containing native types, based on `int`, `float`, `byte`, `rune`, `complex`,
and `string`. You can also use `time.Time` and `time.Duration`.

As usual in [Viper][viper], you can create your config using the reader that
allows creating multiple configs while applying the setup mechanisms for
defaults using the following convenience functions:

```go
    reader := config.New[config.Config]("<prefix>", "<app-name>").
        SetDefaults(func(c *config.ConfigReader[config.Config]{
            c.SetDefault("int", 32)
        }).ReadConfig("main")

    config := reader.GetConfig("main")
```

This creates a standard config reader with defaults from the given config
prototype reading in additional defaults from the `<app-name>[-env].yaml`-file
and environment variables.

The defaults provided by the different options are overwriting each other in
the following order:

1. First, the values provided via the `default`-tags are applied.
2. Second the values provided by the config prototype instance are applied.
3. Third the values provided by [Viper][viper] custom setup calls are applied.
   This also includes the convenient methods provided in this package.
4. Forth the values provided in the `<app-name>[-env].yaml`-file are applied.
5. And finally the values provided via environment variables are applied
   taking the highest precedence.

**Note**: While yo declare the reader with a default config structure, it is
still possible to customize the reader arbitrarily, e.g. with flag support, and
setup any other config structure by using the original [Viper][viper] interface
functions.

A special feature provided by [`go-config`][go-config] is to set up defaults
using a partial or complete config prototype. While in the `New` constructor
automatically an empty prototype is constructed and parsed for `default`-tags,
you can use th `SetDefaultConfig` method to provide any pre-filled (sub-)config
to updated and extend default values.

```go
    reader := config.New("<prefix>", "<app-name>")).
        SetDefaultConfig("", &config.Config{
            Env: "prod",
        }, false).
        SetDefaultConfig("log", &log.Config{
            Level: "debug",
        }, false)
```


## Logger setup

The [`go-config`][go-config] framework supports to set up a `Logger` in
[zerolog][zerolog] and [logrus][logrus] using the provided standard options
out-of-the-box as follows:

```go
    logger := config.Log.Setup[Rus|Zero](writer[, logger])
```

If no logger is provided, the standard logger is configured and returned.

**Note:** While the config supports [zerolog][zerolog], there is currently no
real benefit of using it aside of its having a modern interface. Performance
wise, the necessary transformations for pretty printing logs are a heavy burden
that likely eats up all performance advantages compared to [logrus][logrus].

[zerolog]: <https://github.com/ra/zerolog>
[logrus]: <https://github.com/sirupsen/logrus>


## Build info

Finally, [`go-config`][go-config] in conjunction with [`go-make`][go-make]
supports a build information to track and access the origin of a command,
service or job. While the build information is also auto-discovered, a full
[`go-make`][go-make] integration provides the following variables in the
`main.go`-file.

```go
// Build information variables set by `go-make`.
var (
    // Path contains the package path (set by `go-make`).
    Path string
    // Version contains the custom version (set by `go-make`).
    Version string
    // Build contains the custom build time (set by `go-make`).
    Build string
    // Revision contains the custom revision (set by `go-make`).
    Revision string
    // Commit contains the custom commit time (set by `go-make`).
    Commit string
    // Dirty contains the custom dirty flag (set by `go-make`).
    Dirty string // Bool not supported by ldflags `-X`.
)
```

You can now use this information to set up the default build information in
the config reader by using `SetInfo` during creation as follows:

```go
func main() {
    reader := config.New("<prefix>", "<app-name>", &Config{}).
        SetInfo(info.New(Path, Version, Build, Revision, Commit, Dirty)).
}
```

If you don't want to use [`go-make`][go-make], you can provide the variable
defaults in the `-ldflags="-X main.Path=... -X main.Version=... ...` manually
during your build.


## Building

This project is using [go-make][go-make], which provides default targets for
most common tasks, to initialize, build, test, and run the software of this
project. Read the [go-make manual][go-make-man] for more information about
targets and configuration options.

[go-make]: <https://github.com/tkrop/go-make>
[go-make-man]: <https://github.com/tkrop/go-make/blob/main/MANUAL.md>

The [`Makefile`](Makefile) depends on a preinstalled [`go`][go] for version
management, and makes heavy use of GNU tools, i.e. [`coretils`][core],
[`findutils`][find], ['(g)make'][make], [`(g)awk`][awk], [`(g)sed`][sed], and
not the least [`bash`][bash]. For certain non-core-features it also requires
[`docker`][docker]/[`podman`][podman] and [`curl`][curl]. On MacOS, it uses
[brew][brew] to ensure that the latest versions with the exception
[`docker`][docker]/[`podman`][podman] are.

[go]: <https://go.dev/>
[brew]: <https://brew.sh/>
[curl]: <https://curl.se/>
[docker]: <https://www.docker.com/>
[podman]: <https://podman.io/>
[make]: <https://www.gnu.org/software/make/>
[bash]: <https://www.gnu.org/software/bash/>
[core]: <https://www.gnu.org/software/coreutils/>
[find]: <https://www.gnu.org/software/findutils/>
[awk]: <https://www.gnu.org/software/awk/>
[sed]: <https://www.gnu.org/software/sed/>

**Not:** [go-make][go-make] automatically installs `pre-commit` and `commit-msg`
[hooks][git-hooks] overwriting and deleting pre-existing hooks (see also
[Customizing Git - Git Hooks][git-hooks]). The `pre-commit` hook calls
`make commit` as an alias for executing  `test-go`, `test-unit`, `lint-<level>`,
and `lint-markdown` to enforce successful testing and linting. The `commit-msg`
hook calls `make git-verify message` for validating whether the commit message
is following the [conventional commit][convent-commit] best practice.

[git-hooks]: <https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks>
[convent-commit]: <https://www.conventionalcommits.org/en/v1.0.0/>


## Terms of Usage

This software is open source under the MIT license. You can use it without
restrictions and liabilities. Please give it a star, so that I know. If the
project has more than 25 Stars, I will introduce semantic versions `v1`.


## Contributing

If you like to contribute, please create an issue and/or pull request with a
proper description of your proposal or contribution. I will review it and
provide feedback on it as fast as possible.
