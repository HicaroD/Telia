package diagnostics

import (
	"fmt"
	"io"
	"strings"

	"github.com/HicaroD/Telia/internal/source"
	"github.com/rivo/uniseg"
)

type Renderer struct{ Color bool }

func (r Renderer) Render(w io.Writer, c *Collector) error {
	var out strings.Builder
	for _, d := range c.Diags {
		color := "31"
		severity := "Error"
		if d.Severity == Warning {
			color, severity = "33", "Warning"
		} else if d.Severity == Hint {
			color, severity = "36", "Hint"
		}
		title := d.Title
		if title == "" {
			title = d.Message
		}
		if d.Span.Valid() {
			fmt.Fprintf(&out, "%s · %s\n\n", location(d.Span), r.paint(color, title))
			r.excerpt(&out, d.Span, color)
		} else {
			fmt.Fprintf(&out, "%s · %s\n\n", severity, r.paint(color, title))
		}
		if d.Message != "" && d.Message != title {
			fmt.Fprintf(&out, "    %s\n", d.Message)
		}
		if d.Help != "" {
			fmt.Fprintf(&out, "    %s\n", r.paint("36", strings.ReplaceAll(d.Help, "\n", "\n    ")))
		}
		if d.Help != "" || (d.Message != "" && d.Message != title) {
			out.WriteByte('\n')
		}
		for _, related := range d.Related {
			if !related.Span.Valid() {
				continue
			}
			fmt.Fprintf(
				&out,
				"    %s\n\n",
				r.paint("36", location(related.Span)+" · "+related.Message),
			)
			r.excerpt(&out, related.Span, "36")
		}
	}
	errors, warnings := c.Counts()
	status := "Compilation succeeded"
	if errors > 0 {
		status = "Compilation failed"
	}
	fmt.Fprintf(&out, "%s · %d %s", status, errors, plural(errors, "error"))
	if warnings > 0 {
		fmt.Fprintf(&out, " · %d %s", warnings, plural(warnings, "warning"))
	}
	hints := 0
	for _, d := range c.Diags {
		if d.Severity == Hint {
			hints++
		}
	}
	if hints > 0 {
		fmt.Fprintf(&out, " · %d %s", hints, plural(hints, "hint"))
	}
	out.WriteByte('\n')
	_, err := io.WriteString(w, out.String())
	return err
}

func plural(n int, word string) string {
	if n == 1 {
		return word
	}
	return word + "s"
}

func (r Renderer) paint(color, text string) string {
	if !r.Color {
		return text
	}
	return "\x1b[" + color + "m" + text + "\x1b[0m"
}

func location(s source.Span) string {
	line, column := s.File.Position(s.Start)
	text, start := s.File.Line(line)
	if column-1 <= len(text) {
		column = displayWidth(text[:s.Start-start]) + 1
	}
	return fmt.Sprintf("%s:%d:%d", s.File.Path, line, column)
}

// Tabs use four-column stops; graphemes use their terminal cell width.
func displayWidth(text string) int {
	n := 0
	for i, part := range strings.Split(text, "\t") {
		if i > 0 {
			n += 4 - n%4
		}
		n += uniseg.StringWidth(part)
	}
	return n
}

func expandTabs(text string) string {
	var out strings.Builder
	n := 0
	for i, part := range strings.Split(text, "\t") {
		if i > 0 {
			width := 4 - n%4
			out.WriteString(strings.Repeat(" ", width))
			n += width
		}
		out.WriteString(part)
		n += uniseg.StringWidth(part)
	}
	return out.String()
}

func (r Renderer) excerpt(out *strings.Builder, s source.Span, color string) {
	first, _ := s.File.Position(s.Start)
	last, _ := s.File.Position(s.End)
	if s.End > s.Start {
		last, _ = s.File.Position(s.End - 1)
	}
	width := len(fmt.Sprint(last))
	for line := first; line <= last; line++ {
		text, offset := s.File.Line(line)
		start, end := max(0, s.Start-offset), min(len(text), s.End-offset)
		start = min(start, len(text))
		end = max(start, end)
		prefix := displayWidth(text[:start])
		length := displayWidth(text[:end]) - prefix
		if length == 0 {
			length = 1
		}
		gutter := fmt.Sprintf("    %*d │ ", width, line)
		fmt.Fprintf(out, "%s%s\n", r.paint("90", gutter), expandTabs(text))
		fmt.Fprintf(
			out,
			"%s%s%s\n",
			r.paint("90", strings.Repeat(" ", width+5)+"│ "),
			strings.Repeat(" ", prefix),
			r.paint(color, strings.Repeat("~", length)),
		)
	}
	out.WriteByte('\n')
}
