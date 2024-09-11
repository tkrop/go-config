package log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tkrop/go-config/config"
	"github.com/tkrop/go-config/log"
	"github.com/tkrop/go-testing/test"
)

type setupLoggingParams struct {
	logger           *log.Logger
	config           *log.Config
	expectTimeFormat string
	expectLogLevel   string
	expectLogCaller  bool
}

var testSetupLoggingParams = map[string]setupLoggingParams{
	"read default log config to std logger": {
		config:           &log.Config{},
		expectLogLevel:   log.DefaultLogLevel,
		expectTimeFormat: log.DefaultLogTimeFormat,
		expectLogCaller:  log.DefaultLogCaller,
	},

	"read default log config": {
		logger:           log.New(),
		config:           &log.Config{},
		expectLogLevel:   log.DefaultLogLevel,
		expectTimeFormat: log.DefaultLogTimeFormat,
		expectLogCaller:  log.DefaultLogCaller,
	},

	"change log level debug": {
		logger: log.New(),
		config: &log.Config{
			Level: "debug",
		},
		expectLogLevel:   "debug",
		expectTimeFormat: log.DefaultLogTimeFormat,
	},

	"change time format date": {
		logger: log.New(),
		config: &log.Config{
			TimeFormat: "2024-12-31",
		},
		expectLogLevel:   log.DefaultLogLevel,
		expectTimeFormat: "2024-12-31",
	},

	"change caller to true": {
		logger: log.New(),
		config: &log.Config{
			Caller: true,
		},
		expectLogLevel:   log.DefaultLogLevel,
		expectTimeFormat: log.DefaultLogTimeFormat,
		expectLogCaller:  true,
	},
}

func TestSetupLogging(t *testing.T) {
	test.Map(t, testSetupLoggingParams).
		Run(func(t test.Test, param setupLoggingParams) {
			// Given
			config.New("TEST", "test", &config.Config{}).
				SetDefaults(func(c *config.ConfigReader[config.Config]) {
					c.AddConfigPath("../fixtures")
				}).LoadConfig(t.Name()).Log.Setup(param.logger)

			// When
			param.config.Setup(param.logger)

			// Then
			logger := param.logger
			if logger == nil {
				logger = log.StandardLogger()
			}
			assert.Equal(t, param.expectTimeFormat,
				logger.Formatter.(*log.TextFormatter).TimestampFormat)
			assert.Equal(t, param.expectLogLevel, logger.GetLevel().String())
			assert.Equal(t, param.expectLogCaller, logger.ReportCaller)
		})
}
