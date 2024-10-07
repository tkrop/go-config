package log_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/tkrop/go-config/config"
	"github.com/tkrop/go-config/log"
	"github.com/tkrop/go-config/log/format"
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
	"read default config no logger": {
		config:           &log.Config{},
		expectLogLevel:   DefaultLogLevel,
		expectTimeFormat: DefaultLogTimeFormat,
		expectLogCaller:  DefaultLogCaller,
	},

	"read default config": {
		config:           &log.Config{},
		expectLogLevel:   DefaultLogLevel,
		expectTimeFormat: DefaultLogTimeFormat,
		expectLogCaller:  DefaultLogCaller,
	},

	"log level custom": {
		config: &log.Config{
			Level: "debug",
		},
		expectLogLevel:   "debug",
		expectTimeFormat: DefaultLogTimeFormat,
		expectLogCaller:  DefaultLogCaller,
	},

	"log level invalid": {
		config: &log.Config{
			Level: "detail",
		},
		expectLogLevel:   "info",
		expectTimeFormat: DefaultLogTimeFormat,
		expectLogCaller:  DefaultLogCaller,
	},

	"time format date": {
		config: &log.Config{
			TimeFormat: "2024-12-31",
		},
		expectLogLevel:   DefaultLogLevel,
		expectTimeFormat: "2024-12-31",
		expectLogCaller:  DefaultLogCaller,
	},

	"caller enabled": {
		config: &log.Config{
			Caller: true,
		},
		expectLogLevel:   DefaultLogLevel,
		expectTimeFormat: DefaultLogTimeFormat,
		expectLogCaller:  true,
	},

	"formater text": {
		config: &log.Config{
			Formatter: format.FormatterText,
		},
		expectLogLevel:   DefaultLogLevel,
		expectTimeFormat: DefaultLogTimeFormat,
		expectLogCaller:  DefaultLogCaller,
	},

	"formater json": {
		config: &log.Config{
			Formatter: format.FormatterJSON,
		},
		expectLogLevel:   DefaultLogLevel,
		expectTimeFormat: DefaultLogTimeFormat,
		expectLogCaller:  DefaultLogCaller,
	},

	"formater pretty": {
		config: &log.Config{
			Formatter: format.FormatterPretty,
		},
		expectLogLevel:   DefaultLogLevel,
		expectTimeFormat: DefaultLogTimeFormat,
		expectLogCaller:  DefaultLogCaller,
	},
}

func TestSetupRus(t *testing.T) {
	test.Map(t, testSetupParams).
		Run(func(t test.Test, param setupParams) {
			// Given
			logger := logrus.New()
			logger.SetOutput(&bytes.Buffer{})
			config := config.New[config.Config]("TEST", "test").
				SetSubDefaults("log", param.config, false).
				GetConfig(t.Name())

			// When
			config.Log.SetupRus(logger)

			// Then
			switch param.config.Formatter {
			case format.FormatterText:
				assert.IsType(t, &logrus.TextFormatter{}, logger.Formatter)
				assert.Equal(t, param.expectTimeFormat,
					logger.Formatter.(*logrus.TextFormatter).TimestampFormat)
			case format.FormatterJSON:
				assert.IsType(t, &logrus.JSONFormatter{}, logger.Formatter)
				assert.Equal(t, param.expectTimeFormat,
					logger.Formatter.(*logrus.JSONFormatter).TimestampFormat)
			case format.FormatterPretty:
				assert.IsType(t, &format.Pretty{}, logger.Formatter)
				assert.Equal(t, param.expectTimeFormat,
					logger.Formatter.(*format.Pretty).TimeFormat)
			default:
				assert.IsType(t, &format.Pretty{}, logger.Formatter)
				assert.Equal(t, param.expectTimeFormat,
					logger.Formatter.(*format.Pretty).TimeFormat)
			}
			assert.Equal(t, param.expectLogLevel, logger.GetLevel().String())
			assert.Equal(t, param.expectLogCaller, logger.ReportCaller)
		})
}

func TestSetupNil(t *testing.T) {
	// Given
	config := config.New[config.Config]("TEST", "test").
		GetConfig(t.Name())

	// When
	config.Log.SetupRus(nil)

	// Then
	assert.True(t, true)
}
