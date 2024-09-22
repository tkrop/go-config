package log_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tkrop/go-config/config"
	"github.com/tkrop/go-config/log"
	"github.com/tkrop/go-testing/test"
)

const (
	// Default log level in configuration.
	DefaultLogLevel = "info"
	// Default report caller flag in configuration.
	DefaultLogCaller = false
)

// DefaultLogTimeFormat contains the default timestamp format.
var DefaultLogTimeFormat = time.RFC3339Nano[0:26]

type setupParams struct {
	config           *log.Config
	expectTimeFormat string
	expectLogLevel   string
	expectLogCaller  bool
}

var testSetupParams = map[string]setupParams{
	"read default log config no logger": {
		config:           &log.Config{},
		expectLogLevel:   DefaultLogLevel,
		expectTimeFormat: DefaultLogTimeFormat,
		expectLogCaller:  DefaultLogCaller,
	},

	"read default log config": {
		config:           &log.Config{},
		expectLogLevel:   DefaultLogLevel,
		expectTimeFormat: DefaultLogTimeFormat,
		expectLogCaller:  DefaultLogCaller,
	},

	"change log level debug": {
		config: &log.Config{
			Level: "debug",
		},
		expectLogLevel:   "debug",
		expectTimeFormat: DefaultLogTimeFormat,
		expectLogCaller:  DefaultLogCaller,
	},

	"invalid log level debug": {
		config: &log.Config{
			Level: "detail",
		},
		expectLogLevel:   "info",
		expectTimeFormat: DefaultLogTimeFormat,
		expectLogCaller:  DefaultLogCaller,
	},

	"change time format date": {
		config: &log.Config{
			TimeFormat: "2024-12-31",
		},
		expectLogLevel:   DefaultLogLevel,
		expectTimeFormat: "2024-12-31",
		expectLogCaller:  DefaultLogCaller,
	},

	"change caller to true": {
		config: &log.Config{
			Caller: true,
		},
		expectLogLevel:   DefaultLogLevel,
		expectTimeFormat: DefaultLogTimeFormat,
		expectLogCaller:  true,
	},
}

func TestSetup(t *testing.T) {
	test.Map(t, testSetupParams).
		Run(func(t test.Test, param setupParams) {
			// Given
			logger := log.New()
			config := config.New("TEST", "test", &config.Config{}).
				SetSubDefaults("log", param.config, false).
				GetConfig(t.Name())

			// When
			config.SetupLogger(logger)

			// Then
			assert.Equal(t, param.expectTimeFormat,
				logger.Formatter.(*log.TextFormatter).TimestampFormat)
			assert.Equal(t, param.expectLogLevel, logger.GetLevel().String())
			assert.Equal(t, param.expectLogCaller, logger.ReportCaller)
		})
}

func TestSetupNil(t *testing.T) {
	// Given
	config := config.New("TEST", "test", &config.Config{}).
		GetConfig(t.Name())

	// When
	config.SetupLogger(nil)

	// Then
	assert.True(t, true)
}
