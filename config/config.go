// Package config provides common configuration handling for services, jobs,
// and commands.
package config

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/spf13/viper"
	"github.com/tkrop/go-config/info"
	"github.com/tkrop/go-config/internal/reflect"
	"github.com/tkrop/go-config/log"
)

// Config common application configuration.
type Config struct {
	// Env contains the execution environment, e.g. local, prod, test.
	Env string `default:"prod"`
	// Info default build information.
	Info *info.Info
	// Log default logger setup.
	Log *log.Config
}

// ConfigReader common config reader based on viper.
type ConfigReader[C any] struct {
	*viper.Viper
}

// GetEnvName returns the environment specific configuration file name using
// the given environment prefix and base filename. The filename is extended
// with the environment specific suffix for loading the config file in `yaml`
// format.
func GetEnvName(prefix string, name string) string {
	env := strings.ToLower(os.Getenv(prefix + "_ENV"))
	if env != "" {
		return fmt.Sprintf("%s-%s", name, env)
	}
	return name
}

// New creates a new config reader with the given prefix, name, and pointer to
// a config struct. The config struct is used to evaluate the default config
// values and names from the `default`-tags.
func New[C any](
	prefix, name string, config *C,
) *ConfigReader[C] {
	c := &ConfigReader[C]{
		Viper: viper.New(),
	}

	c.SetConfigName(GetEnvName(prefix, name))
	c.SetConfigType("yaml")
	c.AddConfigPath(".")
	c.SetEnvPrefix(prefix)
	c.AutomaticEnv()
	c.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	c.AllowEmptyEnv(true)

	// Determine parent directory of this code file.
	if _, filename, _, ok := runtime.Caller(1); ok {
		filepath := path.Join(path.Dir(filename), "../..")
		c.AddConfigPath(filepath)
	}

	// Set default values from config type.
	reflect.NewTagWalker("default", nil).
		Walk(config, "", c.setConfigDefault)

	return c
}

// setConfigDefault sets the default value for the given path using the default
// tag value.
func (c *ConfigReader[C]) setConfigDefault(
	_ reflect.Value, path, tag string,
) {
	c.SetDefault(path, tag)
}

// SetDefaults is a convenience method to configure the reader with defaults
// and standard values. It is also calling the provide function to customize
// values and add more defaults.
func (c *ConfigReader[C]) SetDefaults(
	setup func(*ConfigReader[C]),
) *ConfigReader[C] {
	info := info.GetDefault()
	c.SetDefault("info.path", info.Path)
	c.SetDefault("info.version", info.Version)
	c.SetDefault("info.revision", info.Revision)
	c.SetDefault("info.build", info.Build)
	c.SetDefault("info.commit", info.Commit)
	c.SetDefault("info.dirty", info.Dirty)
	c.SetDefault("info.go", info.Go)
	c.SetDefault("info.platform", info.Platform)
	c.SetDefault("info.compiler", info.Compiler)

	if setup != nil {
		setup(c)
	}

	return c
}

// ReadConfig is a convenience method to read the environment specific config
// file to extend the default config. The context is used to distinguish
// different calls in case of a panic by failures while loading the config
// file.
func (c *ConfigReader[C]) ReadConfig(context string) *ConfigReader[C] {
	err := c.ReadInConfig()
	if err != nil {
		log.WithError(err).Errorf("failed to load config [%s]", context)
		panic(fmt.Errorf("failed to load config [%s]: %w", context, err))
	}

	return c
}

// GetConfig is a convenience method to return the config without loading the
// environment specific config file. The context is used to distinguish
// different calls in case of a panic created by failures while unmarschalling
// the config.
func (c *ConfigReader[C]) GetConfig(context string) *C {
	config := new(C)
	err := c.Unmarshal(config)
	if err != nil {
		log.WithError(err).Errorf("failed to unmarshal config [%s]", context)
		panic(fmt.Errorf("failed to unmarshal config [%s]: %w", context, err))
	}

	log.WithFields(log.Fields{
		"config": config,
	}).Debugf("config loaded [%s]", context)

	return config
}

// LoadConfig is a convenience method to load the environment specific config
// file and returns the config. The context is used to distinguish different
// calls in case of a panic created by failures while loading the config file
// and umarshalling the config.
func (c *ConfigReader[C]) LoadConfig(context string) *C {
	return c.ReadConfig(context).GetConfig(context)
}
