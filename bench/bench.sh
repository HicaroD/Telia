#!/bin/bash

#!/bin/bash

hyperfine --warmup 3 --runs 20 "clang -O3 -o c_out code/fib/fib.c"
hyperfine --warmup 3 --runs 20 "go build code/fib/fib.go -o go_out"
hyperfine --warmup 3 --runs 20 "../telia build code/fib/fib.t"

hyperfine --warmup 3 --runs 20 "./c_out"
hyperfine --warmup 3 --runs 20 "./go_out"
hyperfine --warmup 3 --runs 20 "fib"
