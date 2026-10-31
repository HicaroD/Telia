package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/HicaroD/Telia/config"
	"github.com/HicaroD/Telia/internal/ast"
	ccodegen "github.com/HicaroD/Telia/internal/codegen/c"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/parser"
	"github.com/HicaroD/Telia/internal/sema"
	"github.com/HicaroD/Telia/internal/source"
)

func init() {
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
	collector := diagnostics.New()

	loc, err := ast.LocFromPath(path)
	if err != nil {
		collector.ReportAndSave(
			diagnostics.Diag{Message: fmt.Sprintf("failed to get location: %v", err)},
		)
		return "", collector
	}

	exePath, err := compilePipeline(loc, config.BUILD_OPT_DEBUG, collector)
	if err != nil {
		collector.ReportError(err, source.Span{})
		return "", collector
	}
	defer func() { _ = os.RemoveAll(filepath.Dir(exePath)) }() // Best-effort temporary cleanup.

	output, err := RunBinary(exePath)
	if err != nil {
		collector.ReportAndSave(
			diagnostics.Diag{Message: fmt.Sprintf("failed to run binary: %v", err)},
		)
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

// CompileOnly runs parse→sema→codegen and returns the path to the compiled
// binary without executing it. Use this when the test needs to control
// execution itself (e.g. to inspect exit codes or stderr).
func CompileOnly(path string) (string, *diagnostics.Collector) {
	return CompileOnlyWithCFlags(path)
}

func CompileOnlyWithCFlags(path string, flags ...string) (string, *diagnostics.Collector) {
	collector := diagnostics.New()
	loc, err := ast.LocFromPath(path)
	if err != nil {
		collector.ReportAndSave(
			diagnostics.Diag{Message: fmt.Sprintf("failed to get location: %v", err)},
		)
		return "", collector
	}
	exePath, err := compilePipeline(loc, config.BUILD_OPT_DEBUG, collector, flags...)
	if err != nil {
		collector.ReportError(err, source.Span{})
		return "", collector
	}
	return exePath, collector
}

// CompilePackageOnly builds a package without executing its entry point.
func CompilePackageOnly(dirPath string) (string, *diagnostics.Collector) {
	collector := diagnostics.New()
	loc, err := ast.LocFromPath(dirPath)
	if err != nil {
		collector.ReportAndSave(
			diagnostics.Diag{Message: fmt.Sprintf("failed to get location: %v", err)},
		)
		return "", collector
	}
	exePath, err := compilePackagePipeline(dirPath, loc, config.BUILD_OPT_DEBUG, collector)
	if err != nil {
		collector.ReportError(err, source.Span{})
		return "", collector
	}
	return exePath, collector
}

func CompilePackage(dirPath string) (string, *diagnostics.Collector) {
	exePath, collector := CompilePackageOnly(dirPath)
	if collector.HasErrors() {
		return "", collector
	}
	defer func() { _ = os.RemoveAll(filepath.Dir(exePath)) }() // Best-effort temporary cleanup.

	output, err := RunBinary(exePath)
	if err != nil {
		collector.ReportAndSave(
			diagnostics.Diag{Message: fmt.Sprintf("failed to run binary: %v", err)},
		)
		return "", collector
	}

	return output, collector
}

func compilePipeline(
	loc *ast.Loc,
	buildType config.BuildOptimizationType,
	collector *diagnostics.Collector,
	flags ...string,
) (string, error) {
	p := parser.New(collector)
	program, err := p.ParseFileAsProgram(loc.Path, loc, collector)
	if err != nil {
		return "", err
	}

	checker := sema.New(collector)
	err = checker.Check(program)
	if err != nil {
		return "", err
	}

	cg := ccodegen.NewCG(loc, program)
	cg.SetCFlags(flags...)
	return generateTestBinary(cg, buildType)
}

func compilePackagePipeline(
	dirPath string,
	loc *ast.Loc,
	buildType config.BuildOptimizationType,
	collector *diagnostics.Collector,
) (string, error) {
	p := parser.New(collector)
	program, err := p.ParsePackageAsProgram(dirPath, loc)
	if err != nil {
		return "", err
	}

	checker := sema.New(collector)
	err = checker.Check(program)
	if err != nil {
		return "", err
	}

	cg := ccodegen.NewCG(loc, program)
	return generateTestBinary(cg, buildType)
}

func generateTestBinary(
	cg *ccodegen.CCodegen,
	buildType config.BuildOptimizationType,
) (string, error) {
	dir, err := os.MkdirTemp("", "telia-test-")
	if err != nil {
		return "", err
	}
	cg.SetOutputPath(filepath.Join(dir, "program"))
	err = cg.Generate(buildType)
	if err != nil {
		_ = os.RemoveAll(dir) // Preserve the code generation error if cleanup also fails.
		return "", fmt.Errorf("codegen failed: %w", err)
	}

	return cg.ExePath(), nil
}
