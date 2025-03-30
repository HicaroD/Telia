# Benchmarking Fibonacci

This directory contains benchmarking tests for the Fibonacci algorithm
implemented in different programming languages: Telia, C, and Go. The benchmark
evaluates both the **compilation time** and **execution time** for Fibonacci
numbers **35, 40, and 45**.

## Setup

To perform the benchmarks, we use
[Hyperfine](https://github.com/sharkdp/hyperfine), a command-line benchmarking
tool. Install Hyperfine using the following commands:

```bash
wget https://github.com/sharkdp/hyperfine/releases/download/v1.19.0/hyperfine_1.19.0_amd64.deb
sudo dpkg -i hyperfine_1.19.0_amd64.deb
```

Ensure that you have the necessary compilers installed:

- **Telia:** Ensure the Telia compiler binary is compiled in the root project directory.
- **C:** Requires `clang`.
- **Go:** Requires the Go compiler (`go`).

## Usage

The following command will run all benchmarks for compilation and execution time.

```bash 
./bench.sh
```

If you need to change the code for testing more inputs, go to `code/fib` and
change according to your needs.

Once the benchmark is executed, the results will provide insights into:

- **Compilation speed**: How fast each language compiles the Fibonacci program.
- **Execution performance**: How efficiently each implementation computes
  Fibonacci numbers.

## Results

[bench.png](./imgs/bench.png)


| Input  | Telia  | C      | Go     |
|--------|--------|--------|--------|
| Fib(35) | ~21 ms | ~24 ms | ~44 ms |
| Fib(40) | ~220 ms | ~230 ms | ~500 ms |
| Fib(45) | ~2500 ms | ~2600 ms | ~5400 ms |

### **Interpretation**
- **Telia is consistently as fast as, or slightly faster than, C.**
- **Go is significantly slower across all Fibonacci inputs**, with the gap
  increasing as the input size grows.
- **The execution time scales exponentially**, which is expected due to the
  recursive nature of the Fibonacci algorithm without optimization techniques
  such as memoization.

These results demonstrate that Telia's performance is comparable to C, while Go
exhibits slower execution times for this specific algorithm.

