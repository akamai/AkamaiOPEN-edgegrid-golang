package log

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

var timeFormat = "2006/01/02 15:04:05"

var levels = map[slog.Level]string{
	LevelFatal:      "FATAL",
	slog.LevelError: "ERROR",
	slog.LevelWarn:  "WARN",
	slog.LevelInfo:  "INFO",
	slog.LevelDebug: "DEBUG",
	LevelTrace:      "TRACE",
}

// SlogHandler is custom slog handler
type SlogHandler struct {
	h      slog.Handler
	writer io.Writer
	goas   []groupOrAttrs
}
type groupOrAttrs struct {
	group string
	attrs []slog.Attr
}

// Enabled is a wrapper over slog.Handler.Enabled method
func (h *SlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

// Handle handles the Record
// It will only be called when Enabled returns true
func (h *SlogHandler) Handle(_ context.Context, r slog.Record) error {

	buf := make([]byte, 0, 1024)

	// If the record has no Attrs, remove groups at the end of the list; they are empty.
	goas := normalizeGroupsAndAttributes(h.goas, r.NumAttrs())

	for _, goa := range goas {
		if goa.group != "" {
			buf = fmt.Appendf(buf, "%s:\n", goa.group)
		} else {
			for _, a := range goa.attrs {
				buf = h.appendAttr(buf, a)
			}
		}
	}

	r.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttr(buf, a)
		return true
	})

	levelString := levels[r.Level]

	if len(buf) > 0 {
		_, err := fmt.Fprintf(h.writer, "%s [%s] %s: %s\n", r.Time.Format(timeFormat), levelString, r.Message, buf)
		return err
	}
	_, err := fmt.Fprintf(h.writer, "%s [%s] %s\n", r.Time.Format(timeFormat), levelString, r.Message)
	return err
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
// The Handler owns the slice: it may retain, modify or discard it.
func (h *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
func (h *SlogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{group: name})
}

func (h *SlogHandler) withGroupOrAttrs(goa groupOrAttrs) *SlogHandler {
	h2 := *h
	h2.goas = make([]groupOrAttrs, len(h.goas)+1)
	copy(h2.goas, h.goas)
	h2.goas[len(h2.goas)-1] = goa
	return &h2
}

func normalizeGroupsAndAttributes(groupOfAttrs []groupOrAttrs, numAttrs int) []groupOrAttrs {
	goas := groupOfAttrs
	if numAttrs == 0 {
		for len(goas) > 0 && goas[len(goas)-1].group != "" {
			goas = goas[:len(goas)-1]
		}
	}
	return goas
}

func (h *SlogHandler) appendAttr(buf []byte, a slog.Attr) []byte {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return buf
	}

	if len(buf) != 0 {
		buf = fmt.Append(buf, " ")
	}
	switch a.Value.Kind() {
	case slog.KindTime:
		buf = fmt.Appendf(buf, "%s=%s", a.Key, a.Value.Time().Format(timeFormat))

	default:
		buf = parseDefault(buf, a.Key, a.Value.Any())
	}

	return buf

}

func parseDefault(buf []byte, key string, value any) []byte {

	if f, ok := value.(Fields); ok {
		// Inline Fields data if key is empty
		if key == "" {
			inlined := []string{}
			args := f.Get()
			for i := 0; i < len(args); i += 2 {
				inlined = append(inlined, fmt.Sprintf("%s=%v", args[i], args[i+1]))
			}
			buf = fmt.Appendf(buf, "%v", strings.Join(inlined, " "))
			return buf
		}
		t, _ := json.Marshal(f)
		value = string(t)
	}

	buf = fmt.Appendf(buf, "%v=%v", key, value)
	return buf
}

func suppressDefaults(
	next func([]string, slog.Attr) slog.Attr,
) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

// NewSlogHandler returns configured SlogHandler
func NewSlogHandler(w io.Writer, opts *slog.HandlerOptions) *SlogHandler {

	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	b := &bytes.Buffer{}

	return &SlogHandler{
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
		writer: w,
	}
}
