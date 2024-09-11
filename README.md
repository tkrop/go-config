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

Goal of `go-config` is to provide an easy and extensible config framework with
fluent interface based on [Viper][viper] supporting simple default tags for
standard services, jobs, and commands.

[viper]: https://github.com/spf13/viper


## How to start

To start with `go-config` you simply create your config as an extension of the
config provided in this package as follows:

```go
// Config root element for configuration.
type Config struct {
    config.Config `mapstructure:",squash"`

    String string `default:"my-value"`
    Int int       `default:"31"`
    Dur Duration  `default:"1m"` 
}
```

As usual in [Viper][viper], you can now create your config using the reader
that allows to create multiple configuration while applying the default setup
mechanisms using the following convenience functions:

```go
    reader := config.New("<PREFIX>", "<app-name>", &Config{}).
        SetDefaults(func(c *config.ConfigReader[config.Config]{
            c.SetDefault("int", 32)
        }).ReadConfig("main")

    config := reader.GetConfig("main")
```

The defaults provided have via different options are overwriting each other in
the following order:

1. First the values provided via the `default`-tags are applied.
2. Second the values provided by [Viper][viper] setup calls are applied.
3. Third the values provided in the `<app>[-env].yaml`-file are applied.
4. And finally the values provided via environment variables are applied
   taking the highest precedence.

**Note**: While yo declare the reader with a default config structure, it is
still possible to customize the reader arbitrarily, e.g. with flag support, and
setup any other config structure by using the original [Viper][viper] interface
functions.


## Logger setup

The `go-config` framework als supports to set up a [logrus][logrus] Logger
out-of-the-box using configured defaults as follows:

```go
    config := config.New("<PREFIX>", "<app-name>", &Config{}).
        LoadConfig("main")

    logger := config.Log.Setup(log.New())
```

If no logger is provided, the standard logger is configured and returned.

[logrus]: <https://github.com/sirupsen/logrus>


## Build info

Finally, the `go-config` framework also supports in conjunction with
`go-make` a build information to track and access the origin of a command,
service or job. However, the build information must be manually enabled
by adding the following variable set to the `main.go` file.

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

In addition, you have to set up the build information in config reader by
utilizing the variables as follows:

```go
func main() {
    reader := config.New("<PREFIX>", "<app-name>", &Config{}).
        SetInfo(info.New(Path, Version, Build, Revision, Commit, Dirty)).
}
```
