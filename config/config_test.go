package config_test

import (
	stdjson "encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/ory/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tkrop/go-config/config"
	"github.com/tkrop/go-config/internal/filepath"
	ireflect "github.com/tkrop/go-config/internal/reflect"
	"github.com/tkrop/go-testing/mock"
	ref "github.com/tkrop/go-testing/reflect"
	"github.com/tkrop/go-testing/test"
)

var configPaths = []string{filepath.Normalize(".")}

func newConfig(
	env, level string, setup func(*config.Config),
) *config.Config {
	reader := config.NewReader[config.Config]("TC", "test")
	config := reader.GetConfig("test-helper")

	config.Env = env
	config.Log.Level = level
	if setup != nil {
		setup(config)
	}

	return config
}

func isRepoDirty() bool {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) != ""
}

type ConfigParams struct {
	setup  mock.SetupFunc
	reader func(test.Test) *config.Reader[config.Config]
	expect *config.Config
}

var configTestCases = map[string]ConfigParams{
	// Environment selection and overrides.
	"env-selects-config-and-overrides-leaf": {
		reader: func(t test.Test) *config.Reader[config.Config] {
			t.Setenv("TC_ENV", "test")
			t.Setenv("TC_LOG_LEVEL", "trace")
			return config.NewReader[config.Config]("TC", "test").
				SetDefaults(func(r *config.Reader[config.Config]) {
					r.AddConfigPath("fixtures")
				})
		},
		expect: newConfig("test", "trace", nil),
	},
	"env-structured-over-defaults": {
		reader: func(t test.Test) *config.Reader[config.Config] {
			t.Setenv("TC_LOG", `{"level":"warn",`+
				`"timeFormat":"2006-01-02 15:04:05.999999",`+
				`"file":"/dev/stderr","colorMode":"auto",`+
				`"orderMode":"on","formatter":"pretty"}`)
			return config.NewReader[config.Config]("TC", "test")
		},
		expect: newConfig("prod", "warn", nil),
	},
	"env-structured-over-file": {
		reader: func(t test.Test) *config.Reader[config.Config] {
			t.Setenv("TC_LOG", `{"level":"warn",`+
				`"timeFormat":"2006-01-02 15:04:05.999999",`+
				`"file":"/dev/stderr","colorMode":"auto",`+
				`"orderMode":"on","formatter":"pretty"}`)
			return config.NewReader[config.Config]("TC", "test").
				SetDefaults(func(r *config.Reader[config.Config]) {
					r.AddConfigPath("fixtures")
				})
		},
		expect: newConfig("prod", "warn", func(config *config.Config) {
			config.Info.Path = "github.com/tkrop/go-config"
		}),
	},
	"env-structured-parent-over-leaf": {
		reader: func(t test.Test) *config.Reader[config.Config] {
			t.Setenv("TC_LOG", `{"level":"warn",`+
				`"timeFormat":"2006-01-02 15:04:05.999999",`+
				`"file":"/dev/stderr","colorMode":"auto",`+
				`"orderMode":"on","formatter":"pretty"}`)
			t.Setenv("TC_LOG_LEVEL", "trace")
			return config.NewReader[config.Config]("TC", "test")
		},
		expect: newConfig("prod", "warn", nil),
	},
	"env-structured-omitted-field-no-fallback": {
		reader: func(t test.Test) *config.Reader[config.Config] {
			t.Setenv("TC_LOG", `{"timeFormat":"2006-01-02 15:04:05.999999",`+
				`"file":"/dev/stderr","colorMode":"auto",`+
				`"orderMode":"on","formatter":"pretty"}`)
			t.Setenv("TC_LOG_LEVEL", "trace")
			return config.NewReader[config.Config]("TC", "test")
		},
		expect: newConfig("prod", "", nil),
	},
	"env-structured-null": {
		reader: func(t test.Test) *config.Reader[config.Config] {
			t.Setenv("TC_LOG", `null`)
			return config.NewReader[config.Config]("TC", "test")
		},
		expect: newConfig("prod", "info", func(config *config.Config) {
			config.Log = nil
		}),
	},
	"env-structured-malformed-panics": {
		setup: test.Panic(config.NewErrConfig(
			"merging environment", "test",
			errors.Join(ireflect.NewErrTagWalker(
				"json parsing", "log", `{"caller":`,
				stdjson.Unmarshal(
					[]byte(`{"caller":`), &map[string]any{}),
			)),
		).Error()),
		reader: func(t test.Test) *config.Reader[config.Config] {
			t.Setenv("TC_LOG", `{"caller":`)
			return config.NewReader[config.Config]("TC", "test").
				SetDefaults(func(r *config.Reader[config.Config]) {
					r.SetDefault("viper.panic.unmarshal", true)
				})
		},
	},

	// Default initialization and tags.
	"default-no-file": {
		reader: func(_ test.Test) *config.Reader[config.Config] {
			return config.NewReader[config.Config]("TC", "test")
		},
		expect: newConfig("prod", "info", nil),
	},
	"default-invalid-tag-ignored": {
		reader: func(_ test.Test) *config.Reader[config.Config] {
			return config.NewReader[config.Config]("TC", "test").
				SetDefaults(func(r *config.Reader[config.Config]) {
					r.SetDefaultConfig("", &struct {
						Field []string `default:"invalid"`
					}{}, false)
				})
		},
		expect: newConfig("prod", "info", func(config *config.Config) {
			config.Info.Path = "github.com/tkrop/go-config"
			config.Info.Go = runtime.Version()[2:]
			config.Info.Platform = runtime.GOOS + "/" + runtime.GOARCH
			config.Info.Compiler = runtime.Compiler
			config.Info.Dirty = isRepoDirty()
		}),
	},
	"default-invalid-tag-panics": {
		setup: test.Panic(config.NewErrConfig("creating defaults", "",
			ireflect.NewErrTagWalker("yaml parsing", "field", "invalid",
				errors.New("yaml: unmarshal errors:\n  line 1: "+
					"cannot unmarshal !!str `invalid` into []string"))).
			Error()),
		reader: func(_ test.Test) *config.Reader[config.Config] {
			return config.NewReader[config.Config]("TC", "test").
				SetDefaults(func(r *config.Reader[config.Config]) {
					r.SetDefault("viper.panic.defaults", true)
					r.SetDefaultConfig("", &struct {
						Field []string `default:"invalid"`
					}{}, false)
				})
		},
	},

	// Explicit default setters.
	"set-log-level-default": {
		reader: func(_ test.Test) *config.Reader[config.Config] {
			return config.NewReader[config.Config]("TC", "test").
				SetDefaults(func(r *config.Reader[config.Config]) {
					r.SetDefault("log.level", "trace")
				})
		},
		expect: newConfig("prod", "trace", nil),
	},

	// Config file reads.
	"read-file-over-defaults": {
		reader: func(_ test.Test) *config.Reader[config.Config] {
			return config.NewReader[config.Config]("TC", "test").
				SetDefaults(func(r *config.Reader[config.Config]) {
					r.AddConfigPath("fixtures")
				})
		},
		expect: newConfig("prod", "debug", func(config *config.Config) {
			config.Info.Path = "github.com/tkrop/go-config"
		}),
	},
	"read-file-not-found-panics": {
		setup: test.Panic(config.NewErrConfig("loading file", "test",
			ref.NewBuilder[viper.ConfigFileNotFoundError]().
				Set("locations", fmt.Sprintf("%s", configPaths)).
				Set("name", "test").Build()).Error()),
		reader: func(_ test.Test) *config.Reader[config.Config] {
			return config.NewReader[config.Config]("TC", "test").
				SetDefaults(func(r *config.Reader[config.Config]) {
					r.SetDefault("viper.panic.load", true)
				})
		},
	},

	// Config unmarshalling.
	"get-unmarshal-failure-panics": {
		setup: test.Panic(config.NewErrConfig("unmarshal config",
			"test", fmt.Errorf(
				"decoding failed due to the following error(s):\n\n%v",
				"'Info.Dirty' cannot parse value as 'bool': "+
					"strconv.ParseBool: invalid syntax")).Error()),
		reader: func(_ test.Test) *config.Reader[config.Config] {
			return config.NewReader[config.Config]("TC", "test").
				SetDefaults(func(r *config.Reader[config.Config]) {
					r.AddConfigPath("fixtures")
					r.SetDefault("viper.panic.unmarshal", true)
					r.SetDefault("info.dirty", "5s")
				})
		},
	},
}

func TestConfig(t *testing.T) {
	test.Map(t, configTestCases).
		RunSeq(func(t test.Test, param ConfigParams) {
			// Given
			mock.NewMocks(t).Expect(param.setup)

			// When
			reader := param.reader(t)
			config := reader.LoadConfig("test")

			// Then
			assert.Equal(t, param.expect, config)
		})
}

// TODO: improve test to provide meaningful insights about defaults.
type AnyConfigParams struct {
	// setup  mock.SetupFunc
	config any
	expect any
}

type object struct {
	A any
	B any
	C any
	D []int
}

var anyConfigParams = map[string]AnyConfigParams{
	"read-struct-tag": {
		config: &struct {
			S *object `default:"{a: <default>, d: [1,2,3]}"`
		}{S: &object{}},
		expect: &struct {
			S *object `default:"{a: <default>, d: [1,2,3]}"`
		}{S: &object{
			A: "<default>",
			D: []int{1, 2, 3},
		}},
	},
}

func TestAnyConfig(t *testing.T) {
	test.Map(t, anyConfigParams).
		RunSeq(func(t test.Test, param AnyConfigParams) {
			// Given
			// mock.NewMocks(t).Expect(param.setup)
			reader := config.NewReader[any]("TC", "test")
			reader.SetDefaultConfig("", param.config, true)

			// When
			err := reader.Unmarshal(param.config)

			// Then
			require.NoError(t, err)
			assert.Equal(t, param.expect, param.config)
		})
}
