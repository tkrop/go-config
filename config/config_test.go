package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tkrop/go-config/config"
	"github.com/tkrop/go-config/log"
)

func TestReadConfig(t *testing.T) {
	t.Parallel()

	// Given
	reader := config.New("TC", "test", &config.Config{}).
		SetDefaults(func(c *config.ConfigReader[config.Config]) {
			c.AddConfigPath("../fixtures")
		})

	// When
	config := reader.LoadConfig(t.Name())

	// Than
	assert.Equal(t, "prod", config.Env)
	assert.Equal(t, log.DefaultLogLevel, config.Log.Level)
}

func TestReadConfig_UnmarshalFailure(t *testing.T) {
	t.Parallel()
	type Config struct {
		config.Config `mapstructure:",squash"`
		Content       int
	}

	// Given
	defer func() { _ = recover() }()
	reader := config.New("TC", "test", &Config{}).
		SetDefaults(func(c *config.ConfigReader[Config]) {
			c.AddConfigPath("../fixtures")
		})

	// When
	_ = reader.LoadConfig(t.Name())

	// Then
	require.Fail(t, "no panic after unmarschal failure")
}

func TestReadConfig_FileNotFound(t *testing.T) {
	// Given
	defer func() { _ = recover() }()
	t.Setenv("TC_ENV", "other")
	reader := config.New("TC", "test", &config.Config{}).
		SetDefaults(func(c *config.ConfigReader[config.Config]) {
			c.AddConfigPath("../fixtures")
		})

	// When
	_ = reader.LoadConfig(t.Name())

	// Then
	require.Fail(t, "no panic after missing config file")
}

func TestReadConfig_OverridingEnv(t *testing.T) {
	// Given
	t.Setenv("TC_LOG_LEVEL", "trace")
	reader := config.New("TC", "test", &config.Config{}).
		SetDefaults(func(c *config.ConfigReader[config.Config]) {
			c.AddConfigPath("../fixtures")
		})

	// When
	config := reader.LoadConfig(t.Name())

	// Then
	assert.Equal(t, "prod", config.Env)
	assert.Equal(t, "trace", config.Log.Level)
}

func TestReadConfig_OverridingFunc(t *testing.T) {
	t.Parallel()

	// Given
	reader := config.New("TC", "test", &config.Config{}).
		SetDefaults(func(c *config.ConfigReader[config.Config]) {
			c.SetDefault("log.level", "trace")
		})

	// When
	config := reader.GetConfig(t.Name())

	// Then
	assert.Equal(t, "prod", config.Env)
	assert.Equal(t, "trace", config.Log.Level)
}
