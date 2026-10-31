package diagnostics

import "github.com/HicaroD/Telia/internal/source"

type Severity int

const (
	Error Severity = iota
	Warning
	Hint
)

type Related struct {
	Span    source.Span
	Message string
}

type Diag struct {
	Severity Severity
	Title    string
	Message  string
	Span     source.Span
	Help     string
	Related  []Related
}

func (diag Diag) String() string {
	return diag.Message
}

func (diag Diag) Error() string { return diag.Message }
