package log

import (
	"bytes"
	"io"
	"maps"
	"slices"
	"sort"

	"github.com/sirupsen/logrus"
)

// SetupRus is setting up and returning the given logger. It particular sets up
// the log level, the report caller flag, as well as the formatter with color
// and order mode. If no logger is given, the standard logger is set up.
func (c *Config) SetupRus(writer io.Writer, logger *logrus.Logger) *logrus.Logger {
	// Uses the standard logger if no logger is given.
	if logger == nil {
		logger = logrus.StandardLogger()
	}

	logger.SetOutput(writer)
	// #nosec G115 // cannot happen.
	logger.SetLevel(logrus.Level(ParseLevel(c.Level)))
	logger.SetReportCaller(c.Caller)

	// Sets up the log output format.
	switch c.Formatter {
	case FormatterText:
		color := c.ColorMode.Parse(IsTerminal(logger.Out))
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: c.TimeFormat,
			FullTimestamp:   true,
			ForceColors:     color&ColorOn == ColorOn,
			DisableColors:   color&ColorOff == ColorOff,
		})
	case FormatterJSON:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: c.TimeFormat,
		})
	case FormatterPretty:
		fallthrough
	default:
		logger.SetFormatter(NewLogRusPretty(c, writer))
	}

	return logger
}

// LogRusPretty formats logs into a pretty format.
type LogRusPretty struct {
	*Setup
}

// NewLogRusPretty creates a new pretty formatter for logrus.
func NewLogRusPretty(c *Config, writer io.Writer) *LogRusPretty {
	return &LogRusPretty{
		Setup: c.Setup(writer),
	}
}

// Format formats the log entry to a pretty format.
func (p *LogRusPretty) Format(entry *logrus.Entry) ([]byte, error) {
	buffer := NewBuffer(p.Setup, &bytes.Buffer{})
	buffer.WriteString(entry.Time.Format(p.TimeFormat)).
		WriteByte(' ').WriteLevel(Level(entry.Level))
	if entry.HasCaller() {
		buffer.WriteCaller(entry.Caller)
	}
	buffer.WriteByte(' ').WriteString(entry.Message)

	for _, key := range p.getSortedKeys(entry.Data) {
		buffer.WriteByte(' ').WriteData(key, entry.Data[key])
	}
	return buffer.WriteByte('\n').Bytes()
}

// getSortedKeys returns the keys of the given data.
func (p *LogRusPretty) getSortedKeys(data logrus.Fields) []string {
	keys := slices.Collect(maps.Keys(data))
	if p.OrderMode.CheckFlag(OrderOn) {
		sort.Strings(keys)
	}
	return keys
}
