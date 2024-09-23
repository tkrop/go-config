package config_test

import (
	"fmt"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/tkrop/go-config/config"
	"github.com/tkrop/go-config/internal/filepath"
	"github.com/tkrop/go-config/log"
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
			test.Error(viper.ConfigFileNotFoundError{}).Set("name", "test").
				Set("locations", fmt.Sprintf("%s", configPaths)).
				Get("").(error))),
	},

	"panic after unmarschal failure": {
		setup: func(r *config.Reader[config.Config]) {
			r.AddConfigPath("fixtures")
			r.SetDefault("viper.panic.unmarshal", true)
			r.SetDefault("info.dirty", "5s")
		},
		expect: test.Panic(config.NewErrConfig("unmarshal config",
			"test", &mapstructure.Error{
				Errors: []string{"cannot parse 'Info.Dirty' as bool: " +
					"strconv.ParseBool: parsing \"5s\": invalid syntax"},
			})),
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
			reader := config.New[config.Config]("TC", "test").
				SetDefaults(param.setup)

			// When
			reader.LoadConfig("test")

			// Then
			assert.Equal(t, param.expectEnv, reader.GetString("env"))
			assert.Equal(t, param.expectLogLevel, reader.GetString("log.level"))
		})
}

func TestSetupLogger(t *testing.T) {
	t.Parallel()

	// Given
	logger := log.New()
	config := config.New[config.Config]("TC", "test").
		SetDefaults(func(r *config.Reader[config.Config]) {
			r.AddConfigPath("fixtures")
			r.SetDefault("log.level", "trace")
		}).
		GetConfig(t.Name())

	// When
	config.SetupLogger(logger)

	// Then
	assert.Equal(t, log.TraceLevel, logger.GetLevel())
}
