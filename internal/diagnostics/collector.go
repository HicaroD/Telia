package diagnostics

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/HicaroD/Telia/internal/source"
)

var (
	COMPILER_ERROR_FOUND = errors.New("compiler error found")
)

type Collector struct {
	Diags []Diag
	spans map[*source.File]map[int]source.Span
}

func New() *Collector {
	return &Collector{
		Diags: nil,
	}
}

func (collector *Collector) ReportAndSave(diag Diag) {
	diag = collector.normalize(diag)
	collector.Diags = append(collector.Diags, diag)
}

// RegisterSpan retains original spelling even when token lexemes are decoded.
func (collector *Collector) RegisterSpan(span source.Span) {
	if collector.spans == nil {
		collector.spans = make(map[*source.File]map[int]source.Span)
	}
	if collector.spans[span.File] == nil {
		collector.spans[span.File] = make(map[int]source.Span)
	}
	collector.spans[span.File][span.Start] = span
}

// ReportError adapts stage errors at a recovery boundary. Already-collected
// errors use the sentinel to avoid a second diagnostic.
func (collector *Collector) ReportError(err error, span source.Span) {
	if err == nil || errors.Is(err, COMPILER_ERROR_FOUND) {
		return
	}
	var located *locatedError
	hasLocation := errors.As(err, &located)
	if hasLocation {
		span = located.span
	}
	var diag Diag
	if errors.As(err, &diag) {
		if !diag.Span.Valid() {
			diag.Span = span
		}
		collector.ReportAndSave(diag)
		return
	}
	if hasLocation {
		collector.ReportAndSave(Diag{Message: err.Error(), Span: located.span})
		return
	}
	d := collector.normalize(Diag{Message: err.Error()}, span.File)
	if !d.Span.Valid() {
		d.Span = span
	}
	collector.ReportAndSave(d)
}

type locatedError struct {
	err  error
	span source.Span
}

func (e *locatedError) Error() string { return e.err.Error() }
func (e *locatedError) Unwrap() error { return e.err }

// WithSpan keeps the innermost failing expression's location while errors unwind.
func WithSpan(err error, span source.Span) error {
	if err == nil || !span.Valid() {
		return err
	}
	var located *locatedError
	if errors.As(err, &located) {
		return err
	}
	return &locatedError{err: err, span: span}
}

// Decode legacy position prefixes while stage producers are migrated.
var legacyPosition = regexp.MustCompile(`^(?:\[([^\]]+):(\d+):(\d+)\]|(.+?):(\d+):(\d+):)\s*`)

func (collector *Collector) normalize(d Diag, preferred ...*source.File) Diag {
	d.Message = strings.TrimSpace(d.Message)
	match := legacyPosition.FindStringSubmatch(d.Message)
	if match == nil {
		return d
	}
	path, l, col := match[1], match[2], match[3]
	if path == "" {
		path, l, col = match[4], match[5], match[6]
	}
	path = strings.TrimPrefix(path, ":")
	line, _ := strconv.Atoi(l)
	column, _ := strconv.Atoi(col)
	for file, spans := range collector.spans {
		if d.Span.Valid() {
			break
		}
		if len(preferred) > 0 && preferred[0] != nil && file != preferred[0] {
			continue
		}
		if path != file.Path && !strings.HasSuffix(file.Path, "/"+path) &&
			!strings.HasSuffix(path, ":"+file.Path) &&
			path != ":"+file.Path {
			continue
		}
		_, offset := file.Line(line)
		start := offset + column - 1
		if span, ok := spans[start]; ok {
			d.Span = span
		} else if start >= 0 && start <= len(file.Text) {
			d.Span = source.Span{File: file, Start: start, End: min(start+1, len(file.Text))}
		}
		break
	}
	d.Message = strings.TrimSpace(
		strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(d.Message, match[0])), ":"),
	)
	return d
}

func (collector *Collector) Counts() (errors, warnings int) {
	for _, d := range collector.Diags {
		if d.Severity == Warning {
			warnings++
		} else {
			errors++
		}
	}
	return
}

func (collector *Collector) HasErrors() bool {
	errors, _ := collector.Counts()
	return errors > 0
}
