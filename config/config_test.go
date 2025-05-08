package config_test

import (
	"fmt"
	"testing"

	"github.com/ory/viper"
	"github.com/stretchr/testify/assert"

	"github.com/tkrop/go-config/config"
	"github.com/tkrop/go-config/internal/filepath"
	"github.com/tkrop/go-testing/mock"
	"github.com/tkrop/go-testing/test"
)

var configPaths = []string{filepath.Normalize(".")}

type testConfigParam struct {
	setenv         func(test.Test)
	setup          func(*config.Reader[config.Config])
	expect         mock.SetupFunc
	expectEnv      string
	expectLogLevel string
}

var testConfigParams = map[string]testConfigParam{
	"default config without file": {
		expectEnv:      "prod",
		expectLogLevel: "info",
	},

	"default config with file": {
		setup: func(r *config.Reader[config.Config]) {
			r.AddConfigPath("fixtures")
		},
		expectEnv:      "prod",
		expectLogLevel: "debug",
	},

	"read config with overriding env": {
		setenv: func(t test.Test) {
			t.Setenv("TC_ENV", "test")
			t.Setenv("TC_LOG_LEVEL", "trace")
		},
		setup: func(r *config.Reader[config.Config]) {
			r.AddConfigPath("fixtures")
		},
		expectEnv:      "test",
		expectLogLevel: "trace",
	},

	"read config with overriding func": {
		setup: func(r *config.Reader[config.Config]) {
			r.SetDefault("log.level", "trace")
		},
		expectEnv:      "prod",
		expectLogLevel: "trace",
	},

	"panic after file not found": {
		setup: func(r *config.Reader[config.Config]) {
			r.SetDefault("viper.panic.load", true)
		},
		expect: test.Panic(config.NewErrConfig("loading file", "test",
			test.NewBuilder[viper.ConfigFileNotFoundError]().
				Set("locations", fmt.Sprintf("%s", configPaths)).
				Set("name", "test").Build())),
	},

	"panic after unmarschal failure": {
		setup: func(r *config.Reader[config.Config]) {
			r.AddConfigPath("fixtures")
			r.SetDefault("viper.panic.unmarshal", true)
			r.SetDefault("info.dirty", "5s")
		},
		expect: test.Panic(config.NewErrConfig("unmarshal config",
			"test", fmt.Errorf("decoding failed due to the following error(s):\n\n%v",
				"cannot parse 'Info.Dirty' as bool: "+
					"strconv.ParseBool: parsing \"5s\": invalid syntax",
			))),
	},
}

func TestConfig(t *testing.T) {
	test.Map(t, testConfigParams).
		RunSeq(func(t test.Test, param testConfigParam) {
			// Given
			mock.NewMocks(t).Expect(param.expect)
			if param.setenv != nil {
				param.setenv(t)
			}
			reader := config.NewReader[config.Config]("TC", "test").
				SetDefaults(param.setup)

			// When
			reader.LoadConfig("test")

			// Then
			assert.Equal(t, param.expectEnv, reader.GetString("env"))
			assert.Equal(t, param.expectLogLevel, reader.GetString("log.level"))
		})
}
