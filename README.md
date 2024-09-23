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
    reader := config.New("<prefix>", "<app-name>", &Config{}).
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

A special feature provided by [`go-config`][go-config] is to set up the
defaults using a partial or complete config prototype. While you must provide a
complete prototype in the `New` constructor, you can provide any sub-prototype
in the `SetSubDefaults` method as follows:

```go
    reader := config.New("<prefix>", "<app-name>", &config.Config{
            Env: "prod",
        }).SetSubDefaults("<sub-path>", &log.Config{
            Level: "debug",
        }, false).
```


## Logger setup

The [`go-config`][go-config] framework supports to set up a [logrus][logrus]
`Logger`_out-of-the-box using the following two approaches:

```go
    logger := config.SetupLogger(log.New())
    logger := config.Log.Setup(log.New())
```

If no logger is provided, the standard logger is configured and returned.

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

[go-make]: <https://github.com/tkrop/go-make>
