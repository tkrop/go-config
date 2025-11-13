package config_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ory/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tkrop/go-config/config"
	"github.com/tkrop/go-config/info"
	"github.com/tkrop/go-config/internal/filepath"
	ireflect "github.com/tkrop/go-config/internal/reflect"
	"github.com/tkrop/go-testing/mock"
	ref "github.com/tkrop/go-testing/reflect"
	"github.com/tkrop/go-testing/test"
)

var configPaths = []string{filepath.Normalize(".")}

func newConfig(
	env, level string, setup func(info *info.Info),
) *config.Config {
	reader := config.NewReader[config.Config]("TC", "test")
	config := reader.GetConfig("test-helper")

	if setup != nil {
		setup(config.Info)
	}

	config.Env = env
	config.Log.Level = level

	return config
}

type ConfigParams struct {
	setup  mock.SetupFunc
	reader func(test.Test) *config.Reader[config.Config]
	expect *config.Config
}

var configTestCases = map[string]ConfigParams{
	"default config without file": {
		reader: func(_ test.Test) *config.Reader[config.Config] {
			return config.NewReader[config.Config]("TC", "test")
		},
		expect: newConfig("prod", "info", nil),
	},

	"default config with file": {
		reader: func(_ test.Test) *config.Reader[config.Config] {
			return config.NewReader[config.Config]("TC", "test").
				SetDefaults(func(r *config.Reader[config.Config]) {
					r.AddConfigPath("fixtures")
				})
		},
		expect: newConfig("prod", "debug", func(info *info.Info) {
			info.Path = "github.com/tkrop/go-config"
		}),
	},

	"read config with overriding env": {
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

	"read config with overriding func": {
		reader: func(_ test.Test) *config.Reader[config.Config] {
			return config.NewReader[config.Config]("TC", "test").
				SetDefaults(func(r *config.Reader[config.Config]) {
					r.SetDefault("log.level", "trace")
				})
		},
		expect: newConfig("prod", "trace", nil),
	},

	"panic after file not found": {
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

	"panic after unmarshal failure next": {
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

	"error on default config with invalid tag": {
		reader: func(_ test.Test) *config.Reader[config.Config] {
			return config.NewReader[config.Config]("TC", "test").
				SetDefaults(func(r *config.Reader[config.Config]) {
					r.SetDefaultConfig("", &struct {
						Field []string `default:"invalid"`
					}{}, false)
				})
		},
		expect: newConfig("prod", "info", nil),
	},

	"panic on default config with invalid tag": {
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
	"read struct tag": {
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
