package compiler

import (
	"fmt"
	"os/exec"

	"github.com/HicaroD/Telia/config"
	"github.com/HicaroD/Telia/internal/ast"
	"github.com/HicaroD/Telia/internal/codegen/llvm"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/parser"
	"github.com/HicaroD/Telia/internal/sema"
)

var initialized bool

func init() {
	defer func() {
		initialized = true
	}()

	config.SetDevMode(true)
	err := config.SetupConfigDir()
	if err != nil {
		panic(fmt.Sprintf("failed to setup config dir: %v", err))
	}
	err = config.SetupEnvFile()
	if err != nil {
		panic(fmt.Sprintf("failed to setup env file: %v", err))
	}
}

func CompileFile(path string) (string, *diagnostics.Collector) {
	return CompileFileWithBuildType(path, config.BUILD_OPT_DEBUG)
}

func CompileFileWithBuildType(path string, buildType config.BuildOptimizationType) (string, *diagnostics.Collector) {
	collector := diagnostics.New()

	loc, err := ast.LocFromPath(path)
	if err != nil {
		collector.ReportAndSave(diagnostics.Diag{Message: fmt.Sprintf("failed to get location: %v", err)})
		return "", collector
	}

	exePath, err := compilePipeline(loc, buildType, collector)
	if err != nil {
		collector.ReportAndSave(diagnostics.Diag{Message: err.Error()})
		return "", collector
	}

	output, err := RunBinary(exePath)
	if err != nil {
		collector.ReportAndSave(diagnostics.Diag{Message: fmt.Sprintf("failed to run binary: %v", err)})
		return "", collector
	}

	return output, collector
}

func RunBinary(path string) (string, error) {
	cmd := exec.Command(path)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return string(exitErr.Stderr), exitErr
		}
		return "", err
	}
	return string(output), nil
}

func compilePipeline(loc *ast.Loc, buildType config.BuildOptimizationType, collector *diagnostics.Collector) (string, error) {
	p := parser.New(collector)
	program, runtime, err := p.ParseFileAsProgram(loc.Path, loc, collector)
	if err != nil {
		return "", err
	}

	checker := sema.New(collector)
	err = checker.Check(program, runtime)
	if err != nil {
		return "", err
	}

	cg := llvm.NewCG(loc, program, runtime)
	err = cg.Generate(buildType)
	if err != nil {
		return "", fmt.Errorf("codegen failed: %v", err)
	}

	return cg.ExePath(), nil
}
