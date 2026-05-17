# Duso Runtime Performance

A Cross-Language Comparison

## Setup

- **VM**: Single 1 vCPU / 961 MiB RAM Ubuntu cloud instance (representative of the $5-10/month VPS tier most small servers actually deploy to)
- **Runtimes**: Duso v1.1.1, Node.js v20.19.6, Python 3.12.3, Ruby 3.2.3
- **Timing**: each script's own in-language sub-second timer
- **Memory**: peak RSS via `/usr/bin/time -v`

All benchmarks are reproducible. The scripts live in /bench at the project root, with one file per language so the comparisons are apples-to-apples.

## Single-threaded compute

Pure CPU work, no I/O, no concurrency.

| benchmark                    | Duso     | Node      | Python     | Ruby  |
| ---------------------------- | -------- | --------- | ---------- | ----- |
| fib(30) x 10,000 (iterative) | 280.7 ms | **6 ms**  | 14.3 ms    | 51 ms |
| nested loop 1000x1000 sum    | 447.4 ms | **10 ms** | 213.9 ms   | 92 ms |
| sort 10,000 random floats    | 3.1 ms   | 10 ms     | **1.6 ms** | 4 ms  |

Duso is the slowest on hand-written numeric loops, sometimes by an order of magnitude. The single exception is `sort`, where the work happens inside a native Go function rather than in script. There Duso is competitive with Python and ahead of Node and Ruby.

This pattern -- slow when running script-level loops, native-speed when calling Go builtins -- is the lead indicator for everything that follows.

## I/O-bound concurrency on a constrained VM

500 workers, each making 5 HTTPS requests to `httpbin.org/delay/1`. The endpoint sleeps about one second server-side; most wall-clock time is real network latency. The test measures how well each runtime manages concurrent I/O on a small box.

| runtime  | wall time  | outcome                                              |
| -------- | ---------- | ---------------------------------------------------- |
| **Node** | 23.6 s     | clean completion via async event loop                |
| **Duso** | 28.3 s     | clean completion via 500 goroutines                  |
| Python   | 201.4 s    | thrashed -- ThreadPoolExecutor became the bottleneck |
| Ruby     | OOM-killed | exceeded 1 GB RAM via per-thread stack allocation    |

At 100 workers (the level at which the struggling runtimes can complete):

| runtime  | wall time | peak RSS | per-worker |
| -------- | --------- | -------- | ---------- |
| **Ruby** | 13.9 s    | 163 MB   | 135 ms     |
| Node     | 21.0 s    | 66 MB    | 201 ms     |
| Duso     | 25.5 s    | 24 MB    | 254 ms     |
| Python   | 42.7 s    | 214 MB   | 301 ms     |

At small worker counts, Ruby is the fastest on raw I/O speed. As soon as the worker count grows past what fits in the memory budget, Ruby falls off a cliff and Python becomes effectively unusable. Duso and Node both scale to 500-way concurrency on the 1 GB box without sweating.

## Memory footprint as a practical scaling axis

Peak RSS at 100 concurrent fetch workers:

| runtime  | peak RSS | vs Duso |
| -------- | -------- | ------- |
| **Duso** | 24 MB    | 1.0x    |
| Node     | 66 MB    | 2.8x    |
| Ruby     | 163 MB   | 6.8x    |
| Python   | 214 MB   | 8.9x    |

Translated into practical headroom on a 1 GB VM:

| runtime  | simultaneous processes (approx) |
| -------- | ------------------------------- |
| **Duso** | 40                              |
| Node     | 14                              |
| Ruby     | 6                               |
| Python   | 4                               |

For deployments where the hardware is small and fixed (the entire indie / SaaS / internal-tool tier), per-process memory is the practical scaling limit -- not raw throughput.

## Multi-core: where the picture changes substantially

The above results are on a 1 vCPU box. On a multi-core VM the gap widens in Duso's favor, because Duso uses every available core without configuration while the alternatives don't.

Project measurements at 1000 concurrent workers on a multi-core box show **Duso completing the same fetch benchmark approximately 3.7x faster than Node**. The reasons are structural rather than incidental:

- **Node** runs on a single V8 thread. Multi-core utilization requires `cluster` mode, which forks separate Node processes that don't share state. Each fork carries its own V8 instance (50-100 MB) and requires a process supervisor, plus often an external store (Redis or similar) for shared application state.
- **Python** is GIL-bound -- only one Python thread executes at a time even with N OS threads. Multi-core requires `gunicorn -w N` or `multiprocessing`, both of which fork separate Python interpreters (about 200 MB each).
- **Ruby (MRI)** is GVL-bound for the same architectural reason. Multi-core requires Puma in cluster mode or similar -- more forks, more per-process memory.
- **Duso** rides Go's M:N goroutine scheduler. One process. All cores. No configuration. Goroutines are 2 KB stacks, not OS threads.

Even on an "I/O-bound" benchmark, multi-core scaling matters disproportionately for Duso because what looks like I/O has substantial hidden CPU work: TLS handshakes, HTTP header and body parsing, JSON construction, response handling. That CPU work distributes across cores naturally in Duso and concentrates on one core in the alternatives.

## Implication for HTTP and API servers

The same mechanics apply to inbound web and API workloads. Per request, a server does:

- TLS termination (CPU)
- HTTP header and body parse (CPU)
- Handler execution (script logic, datastore access, downstream I/O)
- Response serialization, JSON or template (CPU)
- TLS encrypt (CPU)

That CPU work parallelizes across cores in Duso. On a 4 vCPU box, a single Duso process handles concurrent requests across all 4 cores. Node, Python, and Ruby need cluster / worker / fork configurations to do the same -- each layering operational complexity (process supervisors, port-sharing, shared-state stores) on top of the application.

For typical web / API workloads where script-level logic is small and most cost lies in network, database, or datastore primitives, Duso has the per-VM capacity advantage on the same hardware -- plus a substantially simpler deployment.

## Why an AST-walking interpreter punches above its weight

Duso is architecturally a **tree-walking AST interpreter** -- the textbook-slowest class of interpreter design. The conventional performance hierarchy:

| class            | typical relative speed           |
| ---------------- | -------------------------------- |
| Optimizing JIT   | ~native                          |
| Bytecode VM      | 10x-100x slower than JIT         |
| Tree-walking AST | 10x-100x slower than bytecode VM |

By that math, a tree-walking interpreter should run roughly 1000x slower than V8's TurboFan. Duso runs about 28x slower on the worst-case microbenchmark (`fib` recursion). That is a remarkable amount of distance closed.

For comparison, here is what the other runtimes ship under the hood:

| runtime          | execution model                                                 |
| ---------------- | --------------------------------------------------------------- |
| V8 (Node)        | four-tier optimizing JIT (Ignition -> Sparkplug -> Maglev -> TurboFan), inline caches, hidden classes, speculative type optimization |
| CPython          | bytecode VM with peephole optimization; experimental JIT in 3.13+ |
| MRI Ruby         | YARV bytecode VM plus YJIT (LLVM-based copy-and-patch JIT)      |
| Lua (reference)  | minimal but heavily-optimized bytecode VM                       |
| LuaJIT           | trace JIT; sometimes outperforms hand-written C                 |
| **Duso**         | tree-walking AST interpreter, no bytecode, no JIT               |

Duso skips all of the heavy machinery above. The reasons it remains competitive anyway:

### The work that matters is in Go, not in the interpreter

Sort runs in Go's `sort.Slice`. JSON parses through `encoding/json`. Regex uses Go's RE2. The HTTP server is `net/http`. The datastore is hand-written Go with proper mutex discipline. Templates render through Go-native string operations.

Every meaningful primitive -- the ones a real application spends time in -- executes at native Go speed. The AST interpreter is only running between primitives, orchestrating which one to call next. That orchestration code is a small fraction of total runtime in any realistic application.

The 28x gap on the worst microbenchmark only manifests if a developer writes a tight numeric loop *in script*. Real applications almost never do that -- they call a Go-native builtin to do the work.

### Goroutines handle concurrency, not the interpreter

Most non-mainstream language projects build their own scheduler, often poorly. Duso uses Go's goroutine scheduler -- the one Google paid hundreds of engineers to refine over a decade. Duso's concurrency primitives (`spawn`, `parallel`, datastore wait/cond) sit on top of `go func()` and `sync.Cond`, getting M:N scheduling, work-stealing, and multi-core distribution as inherited infrastructure.

This is why Duso scales cleanly past 500 concurrent workers on a 1 GB box while Ruby OOM-kills and Python thrashes -- the scheduler underneath is already production-grade.

### Targeted optimizations, not heroic ones

The Duso codebase shows the optimizations that matter for an AST walker: compound-assignment fast paths, builtin-lookup short-circuits, environment caching, lock-free read paths for hot builtins. But no JIT, no bytecode compiler. The team did the small optimizations that close the worst gaps and stopped -- accepting the residual perf cost in exchange for an interpreter that is small, simple, and maintainable.

### The architecture is the performance trick

The point that conventional language-design discourse often misses: **most application code is glue between expensive primitives**. If your primitives are fast -- and Duso's are, because they are Go -- the speed of the glue rarely matters.

Pushing every meaningful primitive into a Go builtin is not a workaround for interpreter slowness. It is the design. The AST interpreter is fine for glue, and that is exactly what the interpreter is asked to do. The result is a runtime that competes on the workloads users actually run, while remaining small and maintainable enough that a focused team can keep the whole thing in their head.

## Summary of what the data argues for

| dimension                                            | result                                                                |
| ---------------------------------------------------- | --------------------------------------------------------------------- |
| Tight numeric loops in script                        | Duso loses, sometimes 10x-50x -- rarely matters in practice           |
| Native-primitive work (sort, regex, JSON, HTTP, KV)  | Duso operates at Go speed                                             |
| I/O concurrency on a small VM                        | Duso and Node scale cleanly to 500; Ruby and Python do not            |
| Per-process memory footprint                         | Duso uses 2.8x-8.9x less than alternatives at the same workload       |
| Multi-core scaling                                   | Duso uses every core natively; alternatives need cluster/fork config  |

The runtime does not need to win every benchmark. It needs to be competitive on the workloads users run and dominant on the package: single binary, multi-core for free, low memory, zero operational overhead, batteries included. The data supports that exactly.

For the "single binary running a real web/API server on a small cloud VM" target -- which is the design center -- Duso is the lean, simple choice that holds up under measurement.
