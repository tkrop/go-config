package log_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tkrop/go-testing/mock"
	"github.com/tkrop/go-testing/test"

	"github.com/tkrop/go-config/config"
	"github.com/tkrop/go-config/log"
)

func TestSetupZero(t *testing.T) {
	test.Map(t, setupTestCases).
		Run(func(t test.Test, param setupParams) {
			// Given
			config := config.NewReader[config.Config]("TEST", "test").
				SetDefaultConfig("log", param.config, false).
				GetConfig(t.Name())

			// When
			logger := config.Log.SetupZero(os.Stderr).ZeroLogger()

			// Then
			assert.Equal(t, log.ParseLevel(param.expectLogLevel),
				log.ParseLevel(logger.GetLevel().String()))

			// Check if the writer is set up correctly.
			writer := test.NewAccessor(logger).Get("w")
			require.IsType(t, zerolog.LevelWriterAdapter{}, writer)
			adapter, ok := writer.(zerolog.LevelWriterAdapter)
			require.True(t, ok)

			switch param.config.Formatter {
			case log.FormatterJSON:
				require.IsType(t, &os.File{}, adapter.Writer)

			case log.FormatterText:
				require.IsType(t, zerolog.ConsoleWriter{}, adapter.Writer)
				writer, ok := adapter.Writer.(zerolog.ConsoleWriter)
				require.True(t, ok)

				assert.Equal(t, os.Stderr, writer.Out)
				assert.Equal(t, param.expectTimeFormat, writer.TimeFormat)
				assert.Equal(t, param.expectColorMode.CheckFlag(log.ColorOff),
					writer.NoColor)

			case log.FormatterPretty:
				fallthrough
			default:
				require.IsType(t, &log.ZeroLogPretty{}, adapter.Writer)
				writer, ok := adapter.Writer.(*log.ZeroLogPretty)
				require.True(t, ok)

				assert.Equal(t, os.Stderr, writer.Out)
				assert.Equal(t, param.expectTimeFormat, writer.Setup.TimeFormat)
				assert.Equal(t, param.expectTimeFormat, writer.ConsoleWriter.TimeFormat)
				assert.Equal(t, param.expectColorMode, writer.ColorMode)
				assert.Equal(t, param.expectOrderMode, writer.OrderMode)
			}

			// Check if the hooks are set up with caller hook.
			hooks := test.NewAccessor(logger).Get("hooks")
			require.IsType(t, []zerolog.Hook{}, hooks)
			hookSlice, ok := hooks.([]zerolog.Hook)
			require.True(t, ok)
			if param.expectLogCaller {
				assert.Len(t, hookSlice, 2)
			} else {
				assert.Len(t, hookSlice, 1)
			}
		})
}

type testZeroLogParam struct {
	config       log.Config
	noTerminal   bool
	setup        func(zerolog.Logger)
	expect       mock.SetupFunc
	expectResult string
}

var zeroLogTestCases = map[string]testZeroLogParam{
	// Test levels with default.
	"level panic default": {
		config: log.Config{Level: "panic"},
		setup: func(logger zerolog.Logger) {
			logger.Panic().Msg("panic message")
		},
		expect: test.Panic("panic message"),
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " panic message\n",
	},
	// Fatal is not testable this way since it is calling `os.Exit``. It needs
	// to be tested in spawned process instead.
	// "level fatal default": {
	// 	config: log.Config{Level: "fatal"},
	// 	setup: func(logger zerolog.Logger) {
	// 		logger.Fatal().Msg("fatal message")
	// 	},
	// 	expect: test.Panic("fatal message"),
	// 	expectResult: otime[0:26] + " " +
	// 		levelC(log.FatalLevel) + " fatal message\n",
	// },
	"level error default": {
		config: log.Config{Level: "error"},
		setup: func(logger zerolog.Logger) {
			logger.Error().Msg("error message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.ErrorLevel) + " error message\n",
	},
	"level warn default": {
		config: log.Config{Level: "warn"},
		setup: func(logger zerolog.Logger) {
			logger.Warn().Msg("warn message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.WarnLevel) + " warn message\n",
	},
	"level info default": {
		config: log.Config{Level: "info"},
		setup: func(logger zerolog.Logger) {
			logger.Info().Msg("info message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " info message\n",
	},
	"level debug default": {
		config: log.Config{Level: "debug"},
		setup: func(logger zerolog.Logger) {
			logger.Debug().Msg("debug message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.DebugLevel) + " debug message\n",
	},
	"level trace default": {
		config: log.Config{Level: "trace"},
		setup: func(logger zerolog.Logger) {
			logger.Trace().Msg("trace message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.TraceLevel) + " trace message\n",
	},

	// Test levels with color.
	"level panic color-on": {
		config: log.Config{Level: "panic", ColorMode: log.ColorModeOn},
		setup: func(logger zerolog.Logger) {
			logger.Panic().Msg("panic message")
		},
		expect: test.Panic("panic message"),
		expectResult: otime[0:26] + " " +
			levelC(log.PanicLevel) + " panic message\n",
	},
	// Fatal is not testable this way since it is calling `os.Exit``. It needs
	// to be tested in spawned process instead.
	// "level fatal color-on": {
	// 	config: log.Config{Level: "fatal", ColorMode: log.ColorModeOn},
	// 	setup: func(logger zerolog.Logger) {
	// 		logger.Fatal().Msg("fatal message")
	// 	},
	// 	expect: test.Panic("fatal message"),
	// 	expectResult: otime[0:26] + " " +
	// 		levelC(log.FatalLevel) + " fatal message\n",
	// },
	"level error color-on": {
		config: log.Config{Level: "error", ColorMode: log.ColorModeOn},
		setup: func(logger zerolog.Logger) {
			logger.Error().Msg("error message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.ErrorLevel) + " error message\n",
	},
	"level warn color-on": {
		config: log.Config{Level: "warn", ColorMode: log.ColorModeOn},
		setup: func(logger zerolog.Logger) {
			logger.Warn().Msg("warn message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.WarnLevel) + " warn message\n",
	},
	"level info color-on": {
		config: log.Config{Level: "info", ColorMode: log.ColorModeOn},
		setup: func(logger zerolog.Logger) {
			logger.Info().Msg("info message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " info message\n",
	},
	"level debug color-on": {
		config: log.Config{Level: "debug", ColorMode: log.ColorModeOn},
		setup: func(logger zerolog.Logger) {
			logger.Debug().Msg("debug message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.DebugLevel) + " debug message\n",
	},
	"level trace color-on": {
		config: log.Config{Level: "trace", ColorMode: log.ColorModeOn},
		setup: func(logger zerolog.Logger) {
			logger.Trace().Msg("trace message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.TraceLevel) + " trace message\n",
	},

	// Test levels with color.
	"level panic color-off": {
		config: log.Config{Level: "panic", ColorMode: log.ColorModeOff},
		expect: test.Panic("panic message"),
		setup: func(logger zerolog.Logger) {
			logger.Panic().Msg("panic message")
		},
		expectResult: otime[0:26] + " " +
			level(log.PanicLevel) + " panic message\n",
	},
	// "level fatal color-off": {
	// 	config: log.Config{Level: "fatal", ColorMode: log.ColorModeOff},
	// 	expect: test.Panic("fatal message"),
	// 	setup: func(logger zerolog.Logger) {
	// 		logger.Fatal().Msg("fatal message")
	// 	},
	// 	expectResult: otime[0:26] + " " +
	// 		level(log.FatalLevel) + " fatal message\n",
	// },
	"level error color-off": {
		config: log.Config{Level: "error", ColorMode: log.ColorModeOff},
		setup: func(logger zerolog.Logger) {
			logger.Error().Msg("error message")
		},
		expectResult: otime[0:26] + " " +
			level(log.ErrorLevel) + " error message\n",
	},
	"level warn color-off": {
		config: log.Config{Level: "warning", ColorMode: log.ColorModeOff},
		setup: func(logger zerolog.Logger) {
			logger.Warn().Msg("warn message")
		},
		expectResult: otime[0:26] + " " +
			level(log.WarnLevel) + " warn message\n",
	},
	"level info color-off": {
		config: log.Config{Level: "info", ColorMode: log.ColorModeOff},
		setup: func(logger zerolog.Logger) {
			logger.Info().Msg("info message")
		},
		expectResult: otime[0:26] + " " +
			level(log.InfoLevel) + " info message\n",
	},
	"level debug color-off": {
		config: log.Config{Level: "debug", ColorMode: log.ColorModeOff},
		setup: func(logger zerolog.Logger) {
			logger.Debug().Msg("debug message")
		},
		expectResult: otime[0:26] + " " +
			level(log.DebugLevel) + " debug message\n",
	},
	"level trace color-off": {
		config: log.Config{Level: "trace", ColorMode: log.ColorModeOff},
		setup: func(logger zerolog.Logger) {
			logger.Trace().Msg("trace message")
		},
		expectResult: otime[0:26] + " " +
			level(log.TraceLevel) + " trace message\n",
	},

	// Test order key value data.
	"data default": {
		setup: func(logger zerolog.Logger) {
			logger.Info().Str("key1", "value1").
				Str("key2", "value2").Msg("data message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data ordered": {
		config: log.Config{OrderMode: log.OrderModeOn},
		setup: func(logger zerolog.Logger) {
			logger.Info().Str("key2", "value2").
				Str("key1", "value1").Msg("data message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data unordered": {
		config: log.Config{OrderMode: log.OrderModeOff},
		setup: func(logger zerolog.Logger) {
			logger.Info().Str("key1", "value1").
				Str("key2", "value2").Msg("data message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},

	// Test color modes.
	"data color-off": {
		config: log.Config{ColorMode: log.ColorModeOff},
		setup: func(logger zerolog.Logger) {
			logger.Info().Str("key1", "value1").
				Str("key2", "value2").Msg("data message")
		},
		expectResult: otime[0:26] + " " +
			level(log.InfoLevel) + " data message " +
			data("key1", "value1") + " " +
			data("key2", "value2") + "\n",
	},
	"data color-on": {
		config: log.Config{ColorMode: log.ColorModeOn},
		setup: func(logger zerolog.Logger) {
			logger.Info().Str("key1", "value1").
				Str("key2", "value2").Msg("data message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data color-auto colorized": {
		config: log.Config{ColorMode: log.ColorModeAuto},
		setup: func(logger zerolog.Logger) {
			logger.Info().Str("key1", "value1").
				Str("key2", "value2").Msg("data message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data color-auto not-colorized": {
		noTerminal: true,
		config:     log.Config{ColorMode: log.ColorModeAuto},
		setup: func(logger zerolog.Logger) {
			logger.Info().Str("key1", "value1").
				Str("key2", "value2").Msg("data message")
		},
		expectResult: otime[0:26] + " " +
			level(log.InfoLevel) + " data message " +
			data("key1", "value1") + " " +
			data("key2", "value2") + "\n",
	},
	"data color-levels": {
		config: log.Config{ColorMode: log.ColorModeLevels},
		setup: func(logger zerolog.Logger) {
			logger.Info().Str("key1", "value1").
				Str("key2", "value2").Msg("data message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " data message " +
			data("key1", "value1") + " " +
			data("key2", "value2") + "\n",
	},
	"data color-fields": {
		config: log.Config{ColorMode: log.ColorModeFields},
		setup: func(logger zerolog.Logger) {
			logger.Info().Str("key1", "value1").
				Str("key2", "value2").Msg("data message")
		},
		expectResult: otime[0:26] + " " +
			level(log.InfoLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},
	"data color-levels+fields": {
		config: log.Config{
			ColorMode: log.ColorModeLevels + "|" + log.ColorModeFields,
		},
		setup: func(logger zerolog.Logger) {
			logger.Info().Str("key1", "value1").
				Str("key2", "value2").Msg("data message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " data message " +
			dataC("key1", "value1") + " " +
			dataC("key2", "value2") + "\n",
	},

	// Time format.
	"time default": {
		setup: func(logger zerolog.Logger) {
			logger.Info().Msg("default time message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " " +
			"default time message\n",
	},
	"time short": {
		config: log.Config{
			TimeFormat: "2006-01-02 15:04:05",
		},
		setup: func(logger zerolog.Logger) {
			logger.Info().Msg("short time message")
		},
		expectResult: otime[0:19] + " " +
			levelC(log.InfoLevel) + " " +
			"short time message\n",
	},
	"time long": {
		config: log.Config{
			TimeFormat: "2006-01-02 15:04:05.000000000",
		},
		setup: func(logger zerolog.Logger) {
			logger.Info().Msg("long time message")
		},
		expectResult: otime[0:29] + " " +
			levelC(log.InfoLevel) + " " +
			"long time message\n",
	},

	// Report caller.
	"caller only": {
		setup: func(logger zerolog.Logger) {
			logger.Info().Caller(0).Msg("caller message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " " + "caller message\n",
	},
	"caller report": {
		config: log.Config{Caller: true},
		setup: func(logger zerolog.Logger) {
			logger.Info().Msg("caller report message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " " +
			"[" + caller(-4) + "] caller report message\n",
	},

	// Test error.
	"error output": {
		setup: func(logger zerolog.Logger) {
			logger.Info().Err(assert.AnError).Msg("error message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " error message " +
			dataC("error", assert.AnError.Error()) + "\n",
	},
	"error output color-on": {
		config: log.Config{ColorMode: log.ColorModeOn},
		setup: func(logger zerolog.Logger) {
			logger.Info().Err(assert.AnError).Msg("error message")
		},
		expectResult: otime[0:26] + " " +
			levelC(log.InfoLevel) + " error message " +
			dataC("error", assert.AnError.Error()) + "\n",
	},
	"error output color-off": {
		config: log.Config{ColorMode: log.ColorModeOff},
		setup: func(logger zerolog.Logger) {
			logger.Info().Err(assert.AnError).Msg("error message")
		},
		expectResult: otime[0:26] + " " +
			level(log.InfoLevel) + " error message " +
			data("error", assert.AnError.Error()) + "\n",
	},
}

func TestZeroLog(t *testing.T) {
	assert.NoError(t, terr)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFunc = func() time.Time { return ttime }

	test.Map(t, zeroLogTestCases).
		// Filter("level-panic-color-on", true).
		Run(func(t test.Test, param testZeroLogParam) {
			// Given
			buffer := &bytes.Buffer{}
			config := config.NewReader[config.Config]("X", "app").
				SetDefaultConfig("log", param.config, true).
				SetDefaults(func(r *config.Reader[config.Config]) {
					r.SetDefault("log.level", "trace")
				}).GetConfig("zerolog")
			logger := config.Log.SetupZero(buffer).ZeroLogger()
			pretty := test.NewAccessor(logger).Get("w").(zerolog.LevelWriterAdapter).
				Writer.(*log.ZeroLogPretty)
			pretty.ColorMode = param.config.ColorMode.Parse(!param.noTerminal)

			if param.expect != nil {
				// Then
				mock.NewMocks(t).Expect(param.expect)

				defer func() {
					assert.Equal(t, param.expectResult, buffer.String())
					if err := recover(); err != nil {
						panic(err)
					}
				}()
			}

			// When
			param.setup(logger)

			// Then
			assert.Equal(t, param.expectResult, buffer.String())
		})
}

type testSetupFormatParam struct {
	config *log.Config
	call   func(*log.Setup) string
	expect string
}

var setupFormatTestCases = map[string]testSetupFormatParam{
	// Test time format.
	"time default": {
		config: &log.Config{
			TimeFormat: log.DefaultTimeFormat,
		},
		call: func(s *log.Setup) string {
			return s.FormatTimestamp(itime)
		},
		expect: otime[0:26],
	},
	"time short": {
		config: &log.Config{
			TimeFormat: "2006-01-02 15:04:05",
		},
		call: func(s *log.Setup) string {
			return s.FormatTimestamp(itime)
		},
		expect: otime[0:19],
	},
	"time long": {
		config: &log.Config{
			TimeFormat: "2006-01-02 15:04:05.000000000",
		},
		call: func(s *log.Setup) string {
			return s.FormatTimestamp(itime)
		},
		expect: otime[0:29],
	},
	"time invalid value": {
		config: &log.Config{
			TimeFormat: time.RFC3339,
		},
		call: func(s *log.Setup) string {
			return s.FormatTimestamp("2024-12-31 23:59:59Z+07:00")
		},
		expect: "2024-12-31 23:59:59Z+07:00",
	},
	"time invalid type": {
		config: &log.Config{
			TimeFormat: log.DefaultTimeFormat,
		},
		call: func(s *log.Setup) string {
			return s.FormatTimestamp(1)
		},
		expect: "1",
	},

	// Test level format default.
	"level panic default": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatLevel("panic")
		},
		expect: level(log.PanicLevel),
	},
	"level fatal default": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatLevel("fatal")
		},
		expect: level(log.FatalLevel),
	},
	"level error default": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatLevel("error")
		},
		expect: level(log.ErrorLevel),
	},
	"level warn default": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatLevel("warn")
		},
		expect: level(log.WarnLevel),
	},
	"level info default": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatLevel("info")
		},
		expect: level(log.InfoLevel),
	},
	"level debug default": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatLevel("debug")
		},
		expect: level(log.DebugLevel),
	},
	"level trace default": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatLevel("trace")
		},
		expect: level(log.TraceLevel),
	},
	"level unknown default": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatLevel("unknown")
		},
		expect: level(log.InfoLevel),
	},

	// Test level format color-on.
	"level panic color-on": {
		config: &log.Config{ColorMode: log.ColorModeOn},
		call: func(s *log.Setup) string {
			return s.FormatLevel("panic")
		},
		expect: levelC(log.PanicLevel),
	},
	"level fatal color-on": {
		config: &log.Config{ColorMode: log.ColorModeOn},
		call: func(s *log.Setup) string {
			return s.FormatLevel("fatal")
		},
		expect: levelC(log.FatalLevel),
	},
	"level error color-on": {
		config: &log.Config{ColorMode: log.ColorModeOn},
		call: func(s *log.Setup) string {
			return s.FormatLevel("error")
		},
		expect: levelC(log.ErrorLevel),
	},
	"level warn color-on": {
		config: &log.Config{ColorMode: log.ColorModeOn},
		call: func(s *log.Setup) string {
			return s.FormatLevel("warn")
		},
		expect: levelC(log.WarnLevel),
	},
	"level info color-on": {
		config: &log.Config{ColorMode: log.ColorModeOn},
		call: func(s *log.Setup) string {
			return s.FormatLevel("info")
		},
		expect: levelC(log.InfoLevel),
	},
	"level debug color-on": {
		config: &log.Config{ColorMode: log.ColorModeOn},
		call: func(s *log.Setup) string {
			return s.FormatLevel("debug")
		},
		expect: levelC(log.DebugLevel),
	},
	"level trace color-on": {
		config: &log.Config{ColorMode: log.ColorModeOn},
		call: func(s *log.Setup) string {
			return s.FormatLevel("trace")
		},
		expect: levelC(log.TraceLevel),
	},
	"level unknown color-on": {
		config: &log.Config{ColorMode: log.ColorModeOn},
		call: func(s *log.Setup) string {
			return s.FormatLevel("unknown")
		},
		expect: levelC(log.InfoLevel),
	},

	// Test level format color-off.
	"level panic color-off": {
		config: &log.Config{ColorMode: log.ColorModeOff},
		call: func(s *log.Setup) string {
			return s.FormatLevel("panic")
		},
		expect: level(log.PanicLevel),
	},
	"level fatal color-off": {
		config: &log.Config{ColorMode: log.ColorModeOff},
		call: func(s *log.Setup) string {
			return s.FormatLevel("fatal")
		},
		expect: level(log.FatalLevel),
	},
	"level error color-off": {
		config: &log.Config{ColorMode: log.ColorModeOff},
		call: func(s *log.Setup) string {
			return s.FormatLevel("error")
		},
		expect: level(log.ErrorLevel),
	},
	"level warn color-off": {
		config: &log.Config{ColorMode: log.ColorModeOff},
		call: func(s *log.Setup) string {
			return s.FormatLevel("warn")
		},
		expect: level(log.WarnLevel),
	},
	"level info color-off": {
		config: &log.Config{ColorMode: log.ColorModeOff},
		call: func(s *log.Setup) string {
			return s.FormatLevel("info")
		},
		expect: level(log.InfoLevel),
	},
	"level debug color-off": {
		config: &log.Config{ColorMode: log.ColorModeOff},
		call: func(s *log.Setup) string {
			return s.FormatLevel("debug")
		},
		expect: level(log.DebugLevel),
	},
	"level trace color-off": {
		config: &log.Config{ColorMode: log.ColorModeOff},
		call: func(s *log.Setup) string {
			return s.FormatLevel("trace")
		},
		expect: level(log.TraceLevel),
	},
	"level unknown color-off": {
		config: &log.Config{ColorMode: log.ColorModeOff},
		call: func(s *log.Setup) string {
			return s.FormatLevel("unknown")
		},
		expect: level(log.InfoLevel),
	},

	"level invalid type": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatLevel(1)
		},
		expect: "1",
	},

	// Test caller format.
	"caller report off": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatCaller("caller")
		},
		expect: "",
	},
	"caller report on": {
		config: &log.Config{Caller: true},
		call: func(s *log.Setup) string {
			return s.FormatCaller("caller")
		},
		expect: "[caller]",
	},
	"caller invalid type": {
		config: &log.Config{Caller: true},
		call: func(s *log.Setup) string {
			return s.FormatCaller(1)
		},
		expect: "[1]",
	},

	// Test message format.
	"message default": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatMessage("message")
		},
		expect: "message",
	},
	"message color-on": {
		config: &log.Config{ColorMode: log.ColorModeOn},
		call: func(s *log.Setup) string {
			return s.FormatMessage("message")
		},
		expect: "message",
	},
	"message color-off": {
		config: &log.Config{ColorMode: log.ColorModeOff},
		call: func(s *log.Setup) string {
			return s.FormatMessage("message")
		},
		expect: "message",
	},
	"message invalid type": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatMessage(1)
		},
		expect: "1",
	},

	// Test error field name.
	"error field name": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatErrFieldName("error")
		},
		expect: key("error"),
	},
	"error field name color-on": {
		config: &log.Config{ColorMode: log.ColorModeOn},
		call: func(s *log.Setup) string {
			return s.FormatErrFieldName("error")
		},
		expect: keyC("error"),
	},
	"error field name color-off": {
		config: &log.Config{ColorMode: log.ColorModeOff},
		call: func(s *log.Setup) string {
			return s.FormatErrFieldName("error")
		},
		expect: key("error"),
	},
	"error field name invalid type": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatErrFieldName(1)
		},
		expect: key("1"),
	},

	// Test error field value.
	"error field value": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatErrFieldValue("error")
		},
		expect: "error",
	},
	"error field value color-on": {
		config: &log.Config{ColorMode: log.ColorModeOn},
		call: func(s *log.Setup) string {
			return s.FormatErrFieldValue("error")
		},
		expect: "error",
	},
	"error field value color-off": {
		config: &log.Config{ColorMode: log.ColorModeOff},
		call: func(s *log.Setup) string {
			return s.FormatErrFieldValue("error")
		},
		expect: "error",
	},
	"error field value invalid type": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatErrFieldValue(1)
		},
		expect: "1",
	},

	// Test field name.
	"field name default": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatFieldName("field")
		},
		expect: key("field"),
	},
	"field name color-on": {
		config: &log.Config{ColorMode: log.ColorModeOn},
		call: func(s *log.Setup) string {
			return s.FormatFieldName("field")
		},
		expect: keyC("field"),
	},
	"field name color-off": {
		config: &log.Config{ColorMode: log.ColorModeOff},
		call: func(s *log.Setup) string {
			return s.FormatFieldName("field")
		},
		expect: key("field"),
	},
	"field name invalid type": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatFieldName(1)
		},
		expect: key("1"),
	},

	// Test field value.
	"field value default": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatFieldValue("field")
		},
		expect: `"field"`,
	},
	"field value color-on": {
		config: &log.Config{ColorMode: log.ColorModeOn},
		call: func(s *log.Setup) string {
			return s.FormatFieldValue("field")
		},
		expect: `"field"`,
	},
	"field value color-off": {
		config: &log.Config{ColorMode: log.ColorModeOff},
		call: func(s *log.Setup) string {
			return s.FormatFieldValue("field")
		},
		expect: `"field"`,
	},
	"field value invalid type": {
		config: &log.Config{},
		call: func(s *log.Setup) string {
			return s.FormatFieldValue(1)
		},
		expect: `"1"`,
	},
}

func TestSetupFormat(t *testing.T) {
	test.Map(t, setupFormatTestCases).
		// Filter("level-panic", true).
		RunSeq(func(t test.Test, param testSetupFormatParam) {
			// Given
			s := param.config.Setup(os.Stderr)

			// When
			result := param.call(s)

			// Then
			assert.Equal(t, param.expect, result)
		})
}
