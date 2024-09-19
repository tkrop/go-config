package log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tkrop/go-config/config"
	"github.com/tkrop/go-config/log"
	"github.com/tkrop/go-testing/test"
)

type setupParams struct {
	config           *log.Config
	expectTimeFormat string
	expectLogLevel   string
	expectLogCaller  bool
}

var testSetupParams = map[string]setupParams{
	"read default log config no logger": {
		config:           &log.Config{},
		expectLogLevel:   log.DefaultLogLevel,
		expectTimeFormat: log.DefaultLogTimeFormat,
		expectLogCaller:  log.DefaultLogCaller,
	},

	"read default log config": {
		config:           &log.Config{},
		expectLogLevel:   log.DefaultLogLevel,
		expectTimeFormat: log.DefaultLogTimeFormat,
		expectLogCaller:  log.DefaultLogCaller,
	},

	"change log level debug": {
		config: &log.Config{
			Level: "debug",
		},
		expectLogLevel:   "debug",
		expectTimeFormat: log.DefaultLogTimeFormat,
		expectLogCaller:  log.DefaultLogCaller,
	},

	"invalid log level debug": {
		config: &log.Config{
			Level: "detail",
		},
		expectLogLevel:   "info",
		expectTimeFormat: log.DefaultLogTimeFormat,
		expectLogCaller:  log.DefaultLogCaller,
	},

	"change time format date": {
		config: &log.Config{
			TimeFormat: "2024-12-31",
		},
		expectLogLevel:   log.DefaultLogLevel,
		expectTimeFormat: "2024-12-31",
		expectLogCaller:  log.DefaultLogCaller,
	},

	"change caller to true": {
		config: &log.Config{
			Caller: true,
		},
		expectLogLevel:   log.DefaultLogLevel,
		expectTimeFormat: log.DefaultLogTimeFormat,
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
			config.Log.Setup(logger)

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
	config.Log.Setup(nil)

	// Then
	assert.True(t, true)
}
