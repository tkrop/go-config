package log

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// ParseZeroLevel parses the log level string and returns the corresponding
// zerolog level.
func (c *Config) ParseZeroLevel() zerolog.Level {
	switch strings.ToLower(c.Level) {
	case LevelPanic:
		return zerolog.PanicLevel
	case LevelFatal:
		return zerolog.FatalLevel
	case LevelError:
		return zerolog.ErrorLevel
	case LevelWarn, LevelWarning:
		return zerolog.WarnLevel
	case LevelInfo:
		return zerolog.InfoLevel
	case LevelDebug:
		return zerolog.DebugLevel
	case LevelTrace:
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}

// SetupZero sets up the zerolog logger. It particular it sets up the log
// level, the report caller flag, as well as the formatter with color and order
// mode.
func (c *Config) SetupZero(writer io.Writer) *Config {
	logger := zerolog.New(writer).Level(c.ParseZeroLevel())

	switch c.Formatter {
	case FormatterText:
		color := c.ColorMode.Parse(IsTerminal(writer))
		logger = logger.Output(zerolog.ConsoleWriter{
			Out:        writer,
			NoColor:    color == ColorOff,
			TimeFormat: c.TimeFormat,
		})
	case FormatterJSON:
		logger = logger.Output(writer)
	case FormatterPretty:
		fallthrough
	default:
		logger = logger.Output(NewZeroLogPretty(c, writer))
	}

	context := logger.With().Timestamp()
	if c.Caller {
		context = context.Caller()
	}

	c.logger = context.Logger()

	return c
}

// Zero returns the zerolog logger.
func (c *Config) Zero() zerolog.Logger {
	return c.logger.(zerolog.Logger)
}

// ZeroLogPretty formats logs into a pretty format.
type ZeroLogPretty struct {
	// Setup provides the setup for formatting logs.
	*Setup
	// ConsoleWriter is the console writer used for writing logs.
	zerolog.ConsoleWriter
}

func NewZeroLogPretty(c *Config, writer io.Writer) *ZeroLogPretty {
	setup := c.Setup(writer)
	return &ZeroLogPretty{
		Setup: setup,
		ConsoleWriter: zerolog.ConsoleWriter{
			Out:                 writer,
			TimeFormat:          setup.TimeFormat,
			FormatTimestamp:     setup.FormatTimestamp,
			FormatLevel:         setup.FormatLevel,
			FormatCaller:        setup.FormatCaller,
			FormatMessage:       setup.FormatMessage,
			FormatErrFieldName:  setup.FormatErrFieldName,
			FormatErrFieldValue: setup.FormatErrFieldValue,
			FormatFieldName:     setup.FormatFieldName,
			FormatFieldValue:    setup.FormatFieldValue,
		},
	}
}

func (s *Setup) FormatTimestamp(i any) string {
	if timestamp, ok := i.(string); ok {
		if ttime, err := time.Parse(time.RFC3339, timestamp); err == nil {
			return ttime.Format(s.TimeFormat)
		}
		return timestamp
	}
	return fmt.Sprintf("%v", i)
}

// Format formats the log entry.
func (s *Setup) FormatLevel(i any) string {
	if level, ok := i.(string); ok {
		level := ParseLevel(level)
		buffer := NewBuffer(s, &bytes.Buffer{})
		if s.ColorMode.CheckFlag(ColorLevels) {
			buffer.WriteColored(s.LevelColors[level], s.LevelNames[level])
		} else {
			buffer.WriteString(s.LevelNames[level])
		}
		return buffer.String()
	}
	return fmt.Sprintf("%v", i)
}

// FormatCaller formats the caller.
func (s *Setup) FormatCaller(i any) string {
	if !s.Caller {
		return ""
	} else if caller, ok := i.(string); ok {
		return `[` + caller + `]`
	}
	return fmt.Sprintf("[%v]", i)
}

// FormatMessage formats the message.
func (*Setup) FormatMessage(i any) string {
	if message, ok := i.(string); ok {
		return message
	}
	return fmt.Sprintf("%v", i)
}

// FormatErrFieldName formats the error field name.
func (s *Setup) FormatErrFieldName(i any) string {
	if name, ok := i.(string); ok {
		buffer := NewBuffer(s, &bytes.Buffer{})
		if s.ColorMode.CheckFlag(ColorFields) {
			buffer.WriteColored(ColorError, name)
		} else {
			buffer.WriteString(name)
		}
		return buffer.WriteByte('=').String()
	}
	return fmt.Sprintf("%v=", i)
}

// FormatErrFieldValue formats the error field value.
func (*Setup) FormatErrFieldValue(i any) string {
	if value, ok := i.(string); ok {
		return value
	}
	return fmt.Sprintf("%v", i)
}

// FormatFieldName formats the field name.
func (s *Setup) FormatFieldName(i any) string {
	if field, ok := i.(string); ok {
		buffer := NewBuffer(s, &bytes.Buffer{})
		if s.ColorMode.CheckFlag(ColorFields) {
			buffer.WriteColored(ColorField, field)
		} else {
			buffer.WriteString(field)
		}
		return buffer.WriteByte('=').String()
	}
	return fmt.Sprintf("%v=", i)
}

func (*Setup) FormatFieldValue(i any) string {
	if value, ok := i.(string); ok {
		return `"` + value + `"`
	}
	return fmt.Sprintf("\"%v\"", i)
}
