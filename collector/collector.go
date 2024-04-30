package collector

import (
	"errors"
	"fmt"
)

var (
	COMPILER_ERROR_FOUND = errors.New("compiler error found")
)

type DiagCollector struct {
	Diags []Diag
}

func New() *DiagCollector {
	return &DiagCollector{
		Diags: nil,
	}
}

func (collector *DiagCollector) ReportAndSave(diag Diag) {
	fmt.Println(diag.Message)
	collector.Diags = append(collector.Diags, diag)
}
