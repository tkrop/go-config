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

// Reader common config reader based on viper.
type Reader[C any] struct {
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

// New creates a new config reader with the given prefix, name, and config
// struct. The config struct is evaluate for default config tags and available
// config values to initialize the map. The `default` tags are only used, if
// the config values are zero.
func New[C any](
	prefix, name string, config *C,
) *Reader[C] {
	r := &Reader[C]{
		Viper: viper.New(),
	}

	r.SetConfigName(GetEnvName(prefix, name))
	r.SetConfigType("yaml")
	r.AddConfigPath(".")
	r.SetEnvPrefix(prefix)
	r.AutomaticEnv()
	r.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	r.AllowEmptyEnv(true)

	r.SetSubDefaults("", config, true)

	// Determine parent directory of this code file.
	if _, filename, _, ok := runtime.Caller(1); ok {
		filepath := path.Join(path.Dir(filename), "../..")
		r.AddConfigPath(filepath)
	}

	return r
}

// SetDefaults is a convenience method to configure the reader with defaults
// and standard values. It is also calling the provide function to customize
// values and add more defaults.
func (r *Reader[C]) SetDefaults(
	setup func(*Reader[C]),
) *Reader[C] {
	if setup != nil {
		setup(r)
	}

	return r
}

// SetSubDefaults is a convenience method to update the default values of a
// sub-section configured in the reader by using the given config struct. The
// config struct is scanned for `default`-tags and values to set the defaults.
// Depending on the `zero` flag the default values are either include setting
// zero values or ignoring them.
func (r *Reader[C]) SetSubDefaults(
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
// different calls in case of a panic by failures while loading the config
// file.
func (r *Reader[C]) ReadConfig(context string) *Reader[C] {
	err := r.ReadInConfig()
	if err != nil {
		log.WithError(err).Errorf("failed to load config [%s]", context)
		panic(fmt.Errorf("failed to load config [%s]: %w", context, err))
	}

	return r
}

// GetConfig is a convenience method to return the config without loading the
// environment specific config file. The context is used to distinguish
// different calls in case of a panic created by failures while unmarschalling
// the config.
func (r *Reader[C]) GetConfig(context string) *C {
	config := new(C)
	err := r.Unmarshal(config)
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
func (r *Reader[C]) LoadConfig(context string) *C {
	return r.ReadConfig(context).GetConfig(context)
}
