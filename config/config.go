// Package config provides common configuration handling for services, jobs,
// and commands.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/tkrop/go-config/info"
	"github.com/tkrop/go-config/internal/reflect"
	clog "github.com/tkrop/go-config/log"
)

// ErrConfig is a common error to indicate a configuration error.
var ErrConfig = errors.New("config")

// NewErrConfig is a convenience method to create a new config error with the
// given context wrapping the original error.
func NewErrConfig(message, context string, err error) error {
	return fmt.Errorf("%w - %s [%s]: %w", ErrConfig, message, context, err)
}

// Config common application configuration.
type Config struct {
	// Env contains the execution environment, e.g. local, prod, test.
	Env string `default:"prod"`
	// Info default build information.
	Info *info.Info
	// Log default logger setup.
	Log *clog.Config
}

// Reader common config reader based on viper.
type Reader[C any] struct {
	*viper.Viper
}

// GetEnvName returns the environment specific configuration file name using
// the given environment prefix and base filename. The filename is extended
// with the environment specific suffix for loading the config file in `yaml`
// format.
func GetEnvName(prefix string, name string) string {
	if env := strings.ToLower(os.Getenv(prefix + "_ENV")); env != "" {
		return fmt.Sprintf("%s-%s", name, env)
	}
	return name
}

// New creates a new config reader with the given prefix, name, and config
// struct. The config struct is evaluate for default config tags and available
// config values to initialize the map. The `default` tags are only used, if
// the config values are zero.
func New[C any](
	prefix, name string, setup ...func(*Reader[C]),
) *Reader[C] {
	r := &Reader[C]{
		Viper: viper.New(),
	}

	r.AutomaticEnv()
	r.AllowEmptyEnv(true)
	r.SetEnvPrefix(prefix)
	r.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	r.SetConfigName(GetEnvName(prefix, name))
	r.SetConfigType("yaml")
	r.AddConfigPath(".")
	r.SetDefaultConfig("", new(C), true)
	r.SetDefaults(setup...)

	return r
}

// SetDefaults is a convenience method to configure the reader with defaults
// and standard values. It is also calling the provide function to customize
// values and add more defaults.
func (r *Reader[C]) SetDefaults(
	setup ...func(*Reader[C]),
) *Reader[C] {
	for _, s := range setup {
		if s != nil {
			s(r)
		}
	}
	return r
}

// SetDefaultConfig is a convenience method to update the default values of
// config in the reader by using the given config struct. The config struct is
// scanned for `default`-tags and non-zero values to set the defaults using the
// given key as prefix for constructing the config key-value pairs. This way
// the default config of a whole config struct as well as any sub-config can be
// updated.
//
// Depending on the `zero` flag the default values are either include setting
// zero values or ignoring them.
func (r *Reader[C]) SetDefaultConfig(
	key string, config any, zero bool,
) *Reader[C] {
	info := info.GetDefault()
	r.SetDefault("info.path", info.Path)
	r.SetDefault("info.version", info.Version)
	r.SetDefault("info.revision", info.Revision)
	r.SetDefault("info.build", info.Build)
	r.SetDefault("info.commit", info.Commit)
	r.SetDefault("info.dirty", info.Dirty)
	r.SetDefault("info.go", info.Go)
	r.SetDefault("info.platform", info.Platform)
	r.SetDefault("info.compiler", info.Compiler)

	reflect.NewTagWalker("default", "mapstructure", zero).
		Walk(key, config, r.SetDefault)

	return r
}

// SetDefault is a convenience method to set the default value for the given
// key in the config reader and return the config reader.
//
// *Note:* This method is primarily kept to simplify debugging and testing.
// Currently, it contains no additional logic.
func (r *Reader[C]) SetDefault(key string, value any) {
	r.Viper.SetDefault(key, value)
}

// ReadConfig is a convenience method to read the environment specific config
// file to extend the default config. The context is used to distinguish
// different calls in case of a failure loading the config file.
func (r *Reader[C]) ReadConfig(context string) *Reader[C] {
	if err := r.ReadInConfig(); err != nil {
		err := NewErrConfig("loading file", context, err)
		log.WithFields(log.Fields{
			"context": context,
		}).WithError(err).Warn("no config file found")
		if r.GetBool("viper.panic.load") {
			panic(err)
		}
	}

	return r
}

// GetConfig is a convenience method to return the config without loading the
// environment specific config file. The context is used to distinguish
// different calls in case of a panic created by failures while unmarschalling
// the config.
func (r *Reader[C]) GetConfig(context string) *C {
	config := new(C)
	if err := r.Unmarshal(config); err != nil {
		err := NewErrConfig("unmarshal config", context, err)
		log.WithFields(log.Fields{
			"context": context,
		}).WithError(err).Error("unmarshal config")
		if r.GetBool("viper.panic.unmarshal") {
			panic(err)
		}
	}

	log.WithFields(log.Fields{
		"context": context,
		"config":  config,
	}).Debugf("config loaded")

	return config
}

// LoadConfig is a convenience method to load the environment specific config
// file and returns the config. The context is used to distinguish different
// calls in case of a panic created by failures loading the config file or
// umarshalling the config.
func (r *Reader[C]) LoadConfig(context string) *C {
	return r.ReadConfig(context).GetConfig(context)
}
