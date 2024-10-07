// Package format provides common log formatting based on logrus for services,
// jobs, and commands with integrated configuration loading.
package format

import (
	"bytes"
	"io"
	"maps"
	"os"
	"slices"
	"sort"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// Pretty formats logs into a pretty format.
type Pretty struct {
	// TimeFormat is defining the time format used for printing timestamps.
	TimeFormat string
	// ColorMode is defining the color mode (default = ColorAuto).
	ColorMode ColorMode
	// OrderMode is defining the order mode.
	OrderMode OrderMode

	// LevelNames is defining the names used for marking the different log
	// levels.
	levelNames []string
	// LevelColors is defining the colors used for marking the different log
	// levels.
	levelColors []string

	// initOnce is used to initialize the formatter only once.
	initOnce sync.Once
}

// Init initializes the pretty formatter.
func (p *Pretty) Init(out io.Writer) *Pretty {
	// Set public default fields.
	if p.TimeFormat == "" {
		p.TimeFormat = DefaultTimeFormat
	}
	if p.ColorMode == ColorUnset {
		p.ColorMode = ColorDefault
	}
	if p.OrderMode == OrderOff {
		p.OrderMode = OrderDefault
	}

	// Set default level names and colors.
	if len(p.levelNames) == 0 {
		p.levelNames = DefaultLevelNames
	}
	if len(p.levelColors) == 0 {
		p.levelColors = DefaultLevelColors
	}

	if p.ColorMode&ColorAuto == ColorAuto {
		if IsTerminal(out) {
			p.ColorMode |= ColorOn
		}
	}
	if p.ColorMode&ColorOn == ColorOn {
		p.ColorMode |= ColorLevels | ColorFields
	}
	return p
}

// Format formats the log entry to a pretty format.
func (p *Pretty) Format(entry *log.Entry) ([]byte, error) {
	p.initOnce.Do(func() { p.Init(entry.Logger.Out) })

	buffer := NewBuffer(p, &bytes.Buffer{})
	buffer.WriteString(entry.Time.Format(p.TimeFormat)).
		WriteByte(' ').WriteLevel(entry.Level).WriteCaller(entry).
		WriteByte(' ').WriteString(entry.Message)

	for _, key := range p.getKeys(entry.Data) {
		buffer.WriteByte(' ').WriteData(key, entry.Data[key])
	}

	return buffer.Bytes()
}

// getKeys returns the keys of the given data.
func (p *Pretty) getKeys(data log.Fields) []string {
	keys := slices.Collect(maps.Keys(data))
	if p.OrderMode.CheckFlag(OrderOn) {
		sort.Strings(keys)
	}
	return keys
}

// IsTerminal checks whether the given writer is a terminal.
func IsTerminal(writer io.Writer) bool {
	if file, ok := writer.(*os.File); ok {
		// #nosec G115 // is a safe conversion for files.
		_, err := unix.IoctlGetTermios(int(file.Fd()), unix.TCGETS)
		return err == nil
	}
	return false
}

// Buffer is the interface for writing bytes and strings.
type BufferWriter interface {
	// WriteByte writes the given byte to the writer.
	io.ByteWriter
	// WriteString writes the given string to the writer.
	io.StringWriter

	// Bytes returns the current bytes of the writer.
	Bytes() []byte
}

// Buffer is a buffer for the pretty formatter.
type Buffer struct {
	// pretty is the pretty formatter of the buffer.
	pretty *Pretty
	// buffer is the bytes buffer used for writing.
	buffer BufferWriter

	// err is the error occurred during writing.
	err error
}

// NewBuffer creates a new buffer for the pretty formatter.
func NewBuffer(p *Pretty, b BufferWriter) *Buffer {
	return &Buffer{pretty: p, buffer: b}
}

// WriteByte writes the given byte to the buffer.
//
//nolint:govet // Intentional deviation from the go vet check.
func (b *Buffer) WriteByte(byt byte) *Buffer {
	if b.err != nil {
		return b
	}

	if err := b.buffer.WriteByte(byt); err != nil {
		b.err = err
	}
	return b
}

// WriteString writes the given string to the buffer.
func (b *Buffer) WriteString(str string) *Buffer {
	if b.err != nil {
		return b
	}

	if _, err := b.buffer.WriteString(str); err != nil {
		b.err = err
	}
	return b
}

// WriteColored writes the given text with the given color to the buffer.
func (b *Buffer) WriteColored(color, str string) *Buffer {
	if b.err != nil {
		return b
	}

	// Check if color mode is disabled.
	if b.pretty.ColorMode == ColorOff {
		return b.WriteString(str)
	}

	return b.WriteString("\x1b[").WriteString(color).WriteByte('m').
		WriteString(str).WriteString("\x1b[0m")
}

// WriteLevel writes the given log level to the buffer.
func (b *Buffer) WriteLevel(level log.Level) *Buffer {
	if b.err != nil {
		return b
	}

	if b.pretty.ColorMode.CheckFlag(ColorLevels) {
		return b.WriteColored(b.pretty.levelColors[level],
			b.pretty.levelNames[level])
	}
	return b.WriteString(b.pretty.levelNames[level])
}

// WriteField writes the given key with the given color to the buffer.
func (b *Buffer) WriteField(level log.Level, key string) *Buffer {
	if b.err != nil {
		return b
	}

	if b.pretty.ColorMode.CheckFlag(ColorFields) {
		return b.WriteColored(b.pretty.levelColors[level], key)
	}
	return b.WriteString(key)
}

// WriteCaller writes the caller information to the buffer.
func (b *Buffer) WriteCaller(entry *log.Entry) *Buffer {
	if b.err != nil || !entry.HasCaller() {
		return b
	}

	caller := entry.Caller
	return b.WriteByte(' ').WriteByte('[').
		WriteString(caller.File).WriteByte(':').
		WriteString(strconv.Itoa(caller.Line)).WriteByte('#').
		WriteString(caller.Function).WriteByte(']')
}

// WriteData writes the data to the buffer.
func (b *Buffer) WriteData(key string, value any) *Buffer {
	if b.err != nil {
		return b
	}

	switch key {
	case log.ErrorKey:
		return b.WriteField(log.ErrorLevel, key).
			WriteByte('=').WriteString(value.(error).Error())
	default:
		return b.WriteField(FieldLevel, key).
			WriteByte('=').WriteString(value.(string))
	}
}

// Bytes returns current bytes of the buffer with the current error.
func (b *Buffer) Bytes() ([]byte, error) {
	return b.buffer.Bytes(), b.err
}
