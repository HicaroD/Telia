// Command profile investigates Telia's compiler without changing its CLI.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/HicaroD/Telia/config"
	"github.com/HicaroD/Telia/internal/ast"
	ccodegen "github.com/HicaroD/Telia/internal/codegen/c"
	"github.com/HicaroD/Telia/internal/diagnostics"
	"github.com/HicaroD/Telia/internal/parser"
	"github.com/HicaroD/Telia/internal/sema"
)

type sample struct {
	Functions      int     `json:"functions"`
	Iteration      int     `json:"iteration"`
	ParseMS        float64 `json:"parse_ms"`
	SemaMS         float64 `json:"sema_ms"`
	BackendMS      float64 `json:"backend_ms"`
	ToolchainMS    float64 `json:"toolchain_ms"`
	AllocatedBytes uint64  `json:"allocated_bytes"`
}

type toolCall struct {
	Args []string `json:"args"`
	MS   float64  `json:"ms"`
}

func main() {
	if filepath.Base(os.Args[0]) == "cc" {
		wrapCompiler()
		return
	}
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func wrapCompiler() {
	start := time.Now()
	cmd := exec.Command(os.Getenv("TELIA_BENCH_CC"), os.Args[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()
	f, openErr := os.OpenFile(os.Getenv("TELIA_BENCH_CALLS"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if openErr != nil {
		panic(openErr)
	}
	encodeErr := json.NewEncoder(f).Encode(toolCall{os.Args[1:], float64(time.Since(start).Nanoseconds()) / 1e6})
	closeErr := f.Close()
	if encodeErr != nil {
		panic(encodeErr)
	}
	if closeErr != nil {
		panic(closeErr)
	}
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			os.Exit(exit.ExitCode())
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	out := flag.String("out", "", "unique results directory (normally set by bench.sh)")
	runs := flag.Int("runs", 5, "measured builds per size, after one warmup")
	sizes := flag.String("sizes", "100,500,1000", "comma-separated function counts")
	profileTime := flag.Duration("profile-time", 3*time.Second, "frontend CPU sampling duration at largest size")
	flag.Parse()
	if *out == "" || *runs < 1 || *profileTime <= 0 {
		return fmt.Errorf("out, positive runs, and positive profile-time are required")
	}
	var counts []int
	for _, text := range strings.Split(*sizes, ",") {
		n, err := strconv.Atoi(text)
		if err != nil || n < 1 || n > 10000 {
			return fmt.Errorf("invalid function count %q (expected 1..10000)", text)
		}
		counts = append(counts, n)
	}
	absOut, err := filepath.Abs(*out)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(absOut, 0755); err != nil {
		return err
	}
	cc, err := exec.LookPath("cc")
	if err != nil {
		return err
	}
	cc, err = filepath.Abs(cc)
	if err != nil {
		return err
	}
	self, err := os.Executable()
	if err != nil {
		return err
	}
	shim := filepath.Join(absOut, "tools")
	if err := os.Mkdir(shim, 0755); err != nil {
		return err
	}
	if err := os.Symlink(self, filepath.Join(shim, "cc")); err != nil {
		return err
	}
	calls := filepath.Join(absOut, "toolchain.jsonl")
	for key, value := range map[string]string{"TELIA_BENCH_CC": cc, "TELIA_BENCH_CALLS": calls, "PATH": shim + string(os.PathListSeparator) + os.Getenv("PATH")} {
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	config.ENVS = &config.Envs{STD: filepath.Join(root, "base", "std")}
	var samples []sample
	var largest string
	maxCount := 0
	for _, n := range counts {
		source := filepath.Join(absOut, fmt.Sprintf("functions_%d.t", n))
		if err := os.WriteFile(source, workload(n), 0600); err != nil {
			return err
		}
		if n > maxCount {
			maxCount, largest = n, source
		}
		for i := 0; i <= *runs; i++ {
			// Explicit GC outside timing reduces cross-iteration carryover; this is not a cold-cache build.
			runtime.GC()
			var before, after runtime.MemStats
			runtime.ReadMemStats(&before)
			program, loc, parseMS, semaMS, err := frontend(source)
			if err != nil {
				return err
			}
			if err := os.WriteFile(calls, nil, 0600); err != nil {
				return err
			}
			cg := ccodegen.NewCG(loc, program)
			binary := filepath.Join(absOut, "workload")
			cg.SetOutputPath(binary)
			start := time.Now()
			if err := cg.Generate(config.BUILD_OPT_RELEASE); err != nil {
				return err
			}
			backendMS := ms(start)
			runtime.ReadMemStats(&after)
			// Execution is a correctness check, deliberately outside build timings.
			if output, err := exec.Command(binary).CombinedOutput(); err != nil {
				return fmt.Errorf("workload failed: %w: %s", err, output)
			}
			f, err := os.Open(calls)
			if err != nil {
				return err
			}
			decoder := json.NewDecoder(f)
			toolMS := 0.0
			for decoder.More() {
				var call toolCall
				if err := decoder.Decode(&call); err != nil {
					_ = f.Close()
					return err
				}
				toolMS += call.MS
			}
			if err := f.Close(); err != nil {
				return err
			}
			if i > 0 {
				samples = append(samples, sample{n, i, parseMS, semaMS, backendMS, toolMS, after.TotalAlloc - before.TotalAlloc})
			}
		}
	}
	// CPU sampling is separate from timing runs. Allocation sampling is process-wide.
	runtime.GC()
	cpu, err := os.Create(filepath.Join(absOut, "cpu.pprof"))
	if err != nil {
		return err
	}
	if err := pprof.StartCPUProfile(cpu); err != nil {
		_ = cpu.Close()
		return err
	}
	deadline := time.Now().Add(*profileTime)
	iterations := 0
	for time.Now().Before(deadline) {
		if _, _, _, _, err := frontend(largest); err != nil {
			pprof.StopCPUProfile()
			_ = cpu.Close()
			return err
		}
		iterations++
	}
	pprof.StopCPUProfile()
	if err := cpu.Close(); err != nil {
		return err
	}
	runtime.GC()
	allocs, err := os.Create(filepath.Join(absOut, "allocs.pprof"))
	if err != nil {
		return err
	}
	if err := pprof.Lookup("allocs").WriteTo(allocs, 0); err != nil {
		_ = allocs.Close()
		return err
	}
	if err := allocs.Close(); err != nil {
		return err
	}
	metadata := map[string]any{
		"created_utc": time.Now().UTC().Format(time.RFC3339), "commit": command("git", "rev-parse", "HEAD"),
		"git_status": command("git", "status", "--short"), "go": runtime.Version(), "os": runtime.GOOS, "arch": runtime.GOARCH,
		"cpu": command("sysctl", "-n", "machdep.cpu.brand_string"), "logical_cpus": runtime.NumCPU(),
		"cc": cc, "cc_version": command(cc, "--version"), "flags": "-O3 -lm", "runs": *runs,
		"warmups": 1, "profile_iterations": iterations, "profile_functions": maxCount,
		"profile_mem_rate": runtime.MemProfileRate, "samples": samples,
		"notes": "Warm filesystem caches; in-process compiler API (excludes CLI setup). Backend includes C emission, probes, compilation/linking, and deployment. Toolchain time measures child commands only; wrapper overhead remains in backend. CPU profile covers repeated frontend work, not native compiler CPU. Allocation profile is sampled and process-cumulative, including timing runs.",
	}
	f, err := os.Create(filepath.Join(absOut, "results.json"))
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(metadata); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func frontend(path string) (*ast.Program, *ast.Loc, float64, float64, error) {
	loc, err := ast.LocFromPath(path)
	if err != nil {
		return nil, nil, 0, 0, err
	}
	collector := diagnostics.New()
	var program *ast.Program
	start := time.Now()
	pprof.Do(context.Background(), pprof.Labels("phase", "parse"), func(context.Context) {
		program, err = parser.New(collector).ParseFileAsProgram(path, loc, collector)
	})
	parseMS := ms(start)
	if err == nil && !collector.HasErrors() {
		start = time.Now()
		pprof.Do(context.Background(), pprof.Labels("phase", "sema"), func(context.Context) { err = sema.New(collector).Check(program) })
	}
	if err != nil || collector.HasErrors() {
		_ = (diagnostics.Renderer{}).Render(os.Stderr, collector)
		return nil, nil, 0, 0, fmt.Errorf("frontend failed: %v", err)
	}
	return program, loc, parseMS, ms(start), nil
}

func workload(n int) []byte {
	var b strings.Builder
	b.WriteString("package main\n")
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "fn step_%d(x i32) i32 {\nreturn x + 1\n}\n", i)
	}
	b.WriteString("fn main() i32 {\nx i32 := 0\n")
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "x = step_%d(x)\n", i)
	}
	fmt.Fprintf(&b, "expected i32 := %d\nif x != expected {\nreturn 1\n}\nreturn 0\n}\n", n)
	return []byte(b.String())
}

func ms(start time.Time) float64 { return float64(time.Since(start).Nanoseconds()) / 1e6 }
func command(name string, args ...string) string {
	out, err := exec.Command(name, args...).CombinedOutput()
	if err != nil {
		return fmt.Sprintf("unavailable: %v", err)
	}
	return strings.TrimSpace(string(out))
}
