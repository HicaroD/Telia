package diagnostics

import (
	"errors"
	"fmt"
)

var (
	COMPILER_ERROR_FOUND = errors.New("compiler error found")
)

type Collector struct {
	Diags []Diag
}

func New() *Collector {
	return &Collector{
		Diags: nil,
	}
}

func (collector *Collector) ReportAndSave(diag Diag) {
	fmt.Println(diag.Message)
	collector.Diags = append(collector.Diags, diag)
}
