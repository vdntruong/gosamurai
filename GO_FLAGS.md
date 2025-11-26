# Go Command Flags Reference

A comprehensive guide to all command-line flags available in Go toolchain commands.

## Table of Contents

- [Build Flags (`go build`)](#build-flags-go-build)
- [Test Flags (`go test`)](#test-flags-go-test)
- [Run Flags (`go run`)](#run-flags-go-run)
- [Module Flags (`go mod`)](#module-flags-go-mod)
- [Get Flags (`go get`)](#get-flags-go-get)
- [List Flags (`go list`)](#list-flags-go-list)
- [Vet Flags (`go vet`)](#vet-flags-go-vet)
- [Format Flags (`go fmt`)](#format-flags-go-fmt)
- [Compiler Flags (`-gcflags`)](#compiler-flags--gcflags)
- [Linker Flags (`-ldflags`)](#linker-flags--ldflags)
- [Environment Variables](#environment-variables)

---

## Build Flags (`go build`)

**Command:** `go build [flags] [packages]`

### Common Flags

#### `-o <file>`
Specify output file name or location.

```bash
go build -o myapp main.go
go build -o bin/myapp ./cmd/myapp
```

#### `-v`
Print the names of packages as they are compiled (verbose).

```bash
go build -v
```

#### `-a`
Force rebuilding of packages that are already up-to-date.

```bash
go build -a
```

#### `-n`
Print the commands but do not run them (dry run).

```bash
go build -n
```

#### `-x`
Print the commands as they are executed.

```bash
go build -x
```

#### `-race`
Enable data race detection (only on supported platforms).

```bash
go build -race
go test -race
```

**Note:** Adds runtime overhead (~5-10x slower)

#### `-msan`
Enable memory sanitizer (detects uninitialized reads).

```bash
go build -msan
```

#### `-asan`
Enable address sanitizer (detects memory errors).

```bash
go build -asan
```

#### `-work`
Print the name of the temporary work directory and do not delete it.

```bash
go build -work
# Output: WORK=/var/folders/tmp/go-build123456789
```

#### `-tags <tag_list>`
Build with specific build tags (comma-separated).

```bash
go build -tags=integration,debug
go build -tags="mysql postgres"
```

#### `-buildmode=<mode>`
Specify the build mode.

**Modes:**
- `default` - Standard executable
- `archive` - Build as archive (`.a` file)
- `c-archive` - Build as C archive for embedding in C programs
- `c-shared` - Build as C shared library (`.so`, `.dll`, `.dylib`)
- `shared` - Build as Go shared library
- `plugin` - Build as Go plugin
- `pie` - Position Independent Executable

```bash
go build -buildmode=plugin
go build -buildmode=c-shared -o libmylib.so
```

#### `-mod=<mode>`
Module download mode.

**Modes:**
- `readonly` - Don't update `go.mod`
- `vendor` - Use vendor directory
- `mod` - Download modules as needed (default)

```bash
go build -mod=readonly
go build -mod=vendor
```

#### `-trimpath`
Remove all file system paths from the resulting executable.

```bash
go build -trimpath
```

**Use case:** Reproducible builds, hide local paths

#### `-ldflags <flags>`
Pass flags to the linker. See [Linker Flags](#linker-flags--ldflags).

```bash
go build -ldflags="-s -w"
go build -ldflags="-X main.Version=1.0.0"
```

#### `-gcflags <flags>`
Pass flags to the compiler. See [Compiler Flags](#compiler-flags--gcflags).

```bash
go build -gcflags="-m"        # Print optimization decisions
go build -gcflags="-N -l"     # Disable optimizations for debugging
```

#### `-asmflags <flags>`
Pass flags to the assembler.

```bash
go build -asmflags="-S"       # Print assembly listing
```

#### `-gccgoflags <flags>`
Pass flags to gccgo compiler.

```bash
go build -compiler=gccgo -gccgoflags="-O3"
```

#### `-compiler <name>`
Name of compiler to use: `gc` (default) or `gccgo`.

```bash
go build -compiler=gccgo
```

#### `-installsuffix <suffix>`
Add suffix to package installation directory.

```bash
go build -installsuffix=cgo
```

#### `-overlay <file>`
Read a JSON file describing file system overlay.

```bash
go build -overlay=overlay.json
```

#### `-pgo <file>`
Enable Profile-Guided Optimization (Go 1.21+).

```bash
go build -pgo=cpu.pprof
go build -pgo=auto          # Use default.pgo if it exists
```

---

## Test Flags (`go test`)

**Command:** `go test [flags] [packages]`

### Test Selection Flags

#### `-run <regexp>`
Run only tests matching the regular expression.

```bash
go test -run TestFoo              # Run tests with "TestFoo" in name
go test -run TestFoo/subtest      # Run specific subtest
go test -run '^TestFoo$'          # Exact match
```

#### `-bench <regexp>`
Run benchmarks matching the regular expression.

```bash
go test -bench=.                  # Run all benchmarks
go test -bench=BenchmarkSort      # Run specific benchmark
go test -bench='^BenchmarkSort$'  # Exact match
```

#### `-fuzz <regexp>`
Run fuzz tests matching the regular expression (Go 1.18+).

```bash
go test -fuzz=FuzzParse
go test -fuzz=FuzzParse -fuzztime=30s
```

#### `-skip <regexp>`
Skip tests matching the regular expression (Go 1.20+).

```bash
go test -skip=TestSlow
go test -skip='Integration|Slow'
```

### Test Execution Flags

#### `-short`
Tell long-running tests to shorten their run time.

```bash
go test -short
```

**In test code:**
```go
func TestSomething(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping in short mode")
    }
    // Long-running test
}
```

#### `-timeout <duration>`
Set timeout for tests (default: 10m).

```bash
go test -timeout=30s
go test -timeout=1h
```

#### `-parallel <n>`
Set maximum number of tests to run simultaneously.

```bash
go test -parallel=4
```

#### `-count <n>`
Run each test and benchmark n times.

```bash
go test -count=5              # Run each test 5 times
go test -count=1              # Disable test caching
```

#### `-failfast`
Stop running tests after first failure.

```bash
go test -failfast
```

#### `-shuffle <mode>`
Randomize test execution order (Go 1.17+).

```bash
go test -shuffle=on           # Random seed
go test -shuffle=off          # No shuffling
go test -shuffle=123          # Specific seed for reproducibility
```

### Coverage Flags

#### `-cover`
Enable coverage analysis.

```bash
go test -cover
```

#### `-covermode <mode>`
Set coverage mode.

**Modes:**
- `set` - Did each statement run?
- `count` - How many times did each statement run?
- `atomic` - Like count, but thread-safe (use with `-race`)

```bash
go test -covermode=count
go test -race -covermode=atomic
```

#### `-coverprofile <file>`
Write coverage profile to file.

```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out      # View in browser
go tool cover -func=coverage.out      # Show coverage per function
```

#### `-coverpkg <packages>`
Apply coverage analysis to specified packages.

```bash
go test -coverpkg=./...
```

### Benchmark Flags

#### `-benchtime <duration>`
Run benchmarks for specified duration.

```bash
go test -bench=. -benchtime=10s
go test -bench=. -benchtime=1000x    # Run exactly 1000 iterations
```

#### `-benchmem`
Print memory allocation statistics for benchmarks.

```bash
go test -bench=. -benchmem
```

#### `-cpu <list>`
Specify list of GOMAXPROCS values for benchmarks.

```bash
go test -bench=. -cpu=1,2,4,8
```

### Fuzzing Flags

#### `-fuzztime <duration>`
How long to run fuzzing (Go 1.18+).

```bash
go test -fuzz=FuzzParse -fuzztime=30s
go test -fuzz=FuzzParse -fuzztime=1000x
```

#### `-fuzzminimizetime <duration>`
Time to spend minimizing failing inputs.

```bash
go test -fuzz=FuzzParse -fuzzminimizetime=10s
```

### Output Flags

#### `-v`
Verbose output: log all tests.

```bash
go test -v
```

#### `-json`
Output test results in JSON format.

```bash
go test -json
go test -json | jq '.Action'
```

#### `-outputdir <dir>`
Place test outputs (profiles, etc.) in directory.

```bash
go test -outputdir=./testdata
```

### Profiling Flags

#### `-cpuprofile <file>`
Write CPU profile to file.

```bash
go test -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

#### `-memprofile <file>`
Write memory profile to file.

```bash
go test -memprofile=mem.prof
go tool pprof mem.prof
```

#### `-blockprofile <file>`
Write goroutine blocking profile to file.

```bash
go test -blockprofile=block.prof
```

#### `-mutexprofile <file>`
Write mutex contention profile to file.

```bash
go test -mutexprofile=mutex.prof
```

#### `-trace <file>`
Write execution trace to file.

```bash
go test -trace=trace.out
go tool trace trace.out
```

---

## Run Flags (`go run`)

**Command:** `go run [flags] <package> [arguments...]`

`go run` accepts most build flags:

```bash
go run -race main.go
go run -ldflags="-X main.Version=1.0.0" main.go arg1 arg2
go run -gcflags="-m" main.go
```

---

## Module Flags (`go mod`)

### `go mod init`

```bash
go mod init [module-path]
```

**Example:**
```bash
go mod init github.com/user/project
```

### `go mod tidy`

Remove unused dependencies and add missing ones.

```bash
go mod tidy
go mod tidy -v          # Verbose
go mod tidy -go=1.21    # Set minimum Go version
```

#### `-v`
Verbose output.

#### `-e`
Continue even if there are errors.

```bash
go mod tidy -e
```

### `go mod download`

Download modules to local cache.

```bash
go mod download
go mod download golang.org/x/tools@latest
```

#### `-json`
Print JSON metadata.

```bash
go mod download -json
```

#### `-x`
Print executed commands.

```bash
go mod download -x
```

### `go mod verify`

Verify dependencies have expected content.

```bash
go mod verify
```

### `go mod vendor`

Copy dependencies to vendor directory.

```bash
go mod vendor
go mod vendor -v        # Verbose
go mod vendor -e        # Continue on error
```

### `go mod graph`

Print module requirement graph.

```bash
go mod graph
go mod graph | grep golang.org/x/tools
```

### `go mod why`

Explain why packages or modules are needed.

```bash
go mod why golang.org/x/tools
go mod why -m golang.org/x/tools      # Module-level
go mod why -vendor                     # Check vendor
```

### `go mod edit`

Edit `go.mod` programmatically.

```bash
go mod edit -require=golang.org/x/tools@latest
go mod edit -droprequire=golang.org/x/tools
go mod edit -replace=old@v1.0.0=new@v2.0.0
go mod edit -go=1.21
```

---

## Get Flags (`go get`)

**Command:** `go get [flags] [packages]`

#### Basic Usage

```bash
go get golang.org/x/tools/cmd/goimports
go get golang.org/x/tools/cmd/goimports@latest
go get golang.org/x/tools/cmd/goimports@v0.1.0
go get golang.org/x/tools/cmd/goimports@commit_hash
```

#### `-u`
Update packages to latest versions.

```bash
go get -u ./...                      # Update all dependencies
go get -u golang.org/x/tools         # Update specific package
```

#### `-u=patch`
Update to latest patch version only.

```bash
go get -u=patch ./...
```

#### `-d`
Download only, don't install.

```bash
go get -d golang.org/x/tools
```

#### `-t`
Consider test dependencies.

```bash
go get -t ./...
```

#### `-insecure`
Allow fetching from repositories with insecure protocols.

```bash
go get -insecure example.com/pkg
```

---

## List Flags (`go list`)

**Command:** `go list [flags] [packages]`

#### `-f <format>`
Use custom format template.

```bash
go list -f '{{.ImportPath}}'
go list -f '{{.Name}}: {{.Doc}}'
go list -f '{{join .Deps "\n"}}'
```

#### `-json`
Print JSON output.

```bash
go list -json
go list -json | jq '.Deps'
```

#### `-m`
List modules instead of packages.

```bash
go list -m all                       # All modules
go list -m -u all                    # Show available updates
go list -m -versions golang.org/x/tools  # All versions
```

#### `-deps`
Iterate over dependencies too.

```bash
go list -deps ./...
```

#### `-test`
Include test packages.

```bash
go list -test
```

#### `-compiled`
Include compiled packages.

```bash
go list -compiled
```

---

## Vet Flags (`go vet`)

**Command:** `go vet [flags] [packages]`

#### `-n`
Print commands without running.

```bash
go vet -n
```

#### `-x`
Print commands as executed.

```bash
go vet -x
```

#### Specific Checks

```bash
go vet -printf=false ./...           # Disable printf checker
go vet -unreachable=false ./...      # Disable unreachable code checker
```

**Common checks:**
- `asmdecl` - Assembly declarations
- `assign` - Useless assignments
- `atomic` - Atomic operations
- `bools` - Boolean mistakes
- `buildtag` - Build tag errors
- `cgocall` - Cgo pointer passing
- `composites` - Unkeyed composite literals
- `copylocks` - Locks passed by value
- `httpresponse` - HTTP response mistakes
- `loopclosure` - Loop variable capture
- `lostcancel` - Lost context cancellation
- `nilfunc` - Nil function comparison
- `printf` - Printf format errors
- `shift` - Shift operations
- `stdmethods` - Standard method signatures
- `structtag` - Struct tag errors
- `tests` - Test mistakes
- `unmarshal` - Unmarshal errors
- `unreachable` - Unreachable code
- `unsafeptr` - Unsafe pointer usage
- `unusedresult` - Unused result

---

## Format Flags (`go fmt`)

**Command:** `go fmt [flags] [packages]`

`go fmt` is a wrapper around `gofmt -l -w`.

#### `-n`
Print commands without running.

```bash
go fmt -n
```

#### `-x`
Print commands as executed.

```bash
go fmt -x
```

**Note:** For more control, use `gofmt` directly:

```bash
gofmt -w .                  # Format in place
gofmt -d file.go            # Show diff
gofmt -s file.go            # Simplify code
gofmt -r 'rule' file.go     # Rewrite using rule
```

---

## Compiler Flags (`-gcflags`)

Pass flags to the Go compiler using `-gcflags`.

```bash
go build -gcflags="<flags>"
go build -gcflags="all=<flags>"        # Apply to all packages
go build -gcflags="main=<flags>"       # Apply to main package only
```

### Common Compiler Flags

#### `-N`
Disable optimizations.

```bash
go build -gcflags="-N"
```

#### `-l`
Disable inlining.

```bash
go build -gcflags="-l"
```

#### `-N -l` (together)
Disable optimizations and inlining (better for debugging).

```bash
go build -gcflags="-N -l"
dlv debug --build-flags="-gcflags='-N -l'"
```

#### `-m`
Print optimization decisions.

```bash
go build -gcflags="-m"              # Level 1
go build -gcflags="-m -m"           # Level 2 (more detail)
go build -gcflags="-m=2"            # Alternative syntax
```

**Shows:**
- Inlining decisions
- Escape analysis (what goes to heap vs stack)
- Bounds check elimination

#### `-S`
Print assembly listing.

```bash
go build -gcflags="-S" 2>&1 | less
```

#### `-d=<flags>`
Enable debug flags.

```bash
go build -gcflags="-d=ssa/check_bce/debug=1"  # Debug bounds check elimination
```

### Escape Analysis Example

```bash
go build -gcflags="-m" main.go
# Output:
# main.go:10:13: inlining call to fmt.Println
# main.go:8:2: moved to heap: x
```

---

## Linker Flags (`-ldflags`)

Pass flags to the Go linker using `-ldflags`.

```bash
go build -ldflags="<flags>"
```

### Common Linker Flags

#### `-s`
Omit symbol table and debug information.

```bash
go build -ldflags="-s"
```

**Result:** Smaller binary (~25-30% reduction)

#### `-w`
Omit DWARF symbol table.

```bash
go build -ldflags="-w"
```

**Result:** Smaller binary, can't use debuggers

#### `-s -w` (together)
Maximum size reduction.

```bash
go build -ldflags="-s -w"
```

**Result:** Smallest binary, no debugging symbols

#### `-X <importpath.name>=<value>`
Set string variable value at link time.

```bash
go build -ldflags="-X main.Version=1.0.0"
go build -ldflags="-X 'main.BuildTime=$(date)'"
go build -ldflags="-X main.GitCommit=$(git rev-parse HEAD)"
```

**In code:**
```go
package main

var (
    Version   string = "dev"
    BuildTime string = "unknown"
    GitCommit string = "unknown"
)

func main() {
    fmt.Printf("Version: %s\n", Version)
    fmt.Printf("Built: %s\n", BuildTime)
    fmt.Printf("Commit: %s\n", GitCommit)
}
```

#### `-extldflags <flags>`
Pass flags to external linker.

```bash
go build -ldflags="-extldflags=-static"  # Static linking
```

#### `-linkmode <mode>`
Set link mode.

**Modes:**
- `internal` - Use Go linker
- `external` - Use external linker (gcc, clang)

```bash
go build -ldflags="-linkmode=external"
```

#### Combined Example

```bash
go build \
  -ldflags="-s -w -X main.Version=1.0.0 -X 'main.BuildTime=$(date)'" \
  -o myapp
```

---

## Environment Variables

### `GOOS` and `GOARCH`

Cross-compilation targets.

```bash
GOOS=linux GOARCH=amd64 go build
GOOS=windows GOARCH=amd64 go build -o app.exe
GOOS=darwin GOARCH=arm64 go build    # Apple Silicon
```

**Common combinations:**
```bash
GOOS=linux GOARCH=amd64       # Linux 64-bit
GOOS=linux GOARCH=arm64       # Linux ARM64
GOOS=darwin GOARCH=amd64      # macOS Intel
GOOS=darwin GOARCH=arm64      # macOS Apple Silicon
GOOS=windows GOARCH=amd64     # Windows 64-bit
GOOS=freebsd GOARCH=amd64     # FreeBSD 64-bit
```

**List all supported platforms:**
```bash
go tool dist list
```

### `CGO_ENABLED`

Enable/disable cgo.

```bash
CGO_ENABLED=0 go build           # Disable cgo (static binary)
CGO_ENABLED=1 go build           # Enable cgo
```

### `GOPATH`

Go workspace location (less important with modules).

```bash
export GOPATH=$HOME/go
```

### `GOBIN`

Where `go install` puts binaries.

```bash
export GOBIN=$HOME/bin
```

### `GOPROXY`

Module proxy (default: https://proxy.golang.org).

```bash
export GOPROXY=https://proxy.golang.org,direct
export GOPROXY=direct              # Bypass proxy
```

### `GOPRIVATE`

Private modules (bypass proxy and checksum).

```bash
export GOPRIVATE=github.com/mycompany/*
export GOPRIVATE=*.mycompany.com
```

### `GONOSUMDB`

Skip checksum database for these modules.

```bash
export GONOSUMDB=github.com/mycompany/*
```

### `GOTOOLCHAIN`

Specify Go toolchain version (Go 1.21+).

```bash
export GOTOOLCHAIN=go1.21.0
export GOTOOLCHAIN=local          # Use installed version
```

### `GODEBUG`

Runtime debugging options.

```bash
GODEBUG=gctrace=1 go run main.go               # GC trace
GODEBUG=schedtrace=1000 go run main.go         # Scheduler trace
GODEBUG=allocfreetrace=1 go run main.go        # Allocation trace
GODEBUG=http2debug=2 go run main.go            # HTTP/2 debug
```

**Multiple options:**
```bash
GODEBUG=gctrace=1,schedtrace=1000 go run main.go
```

---

## Practical Examples

### Build for Production

```bash
go build \
  -trimpath \
  -ldflags="-s -w -X main.Version=$(git describe --tags)" \
  -o dist/myapp
```

### Cross-Compile for Multiple Platforms

```bash
#!/bin/bash
VERSION=$(git describe --tags)

GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$VERSION" -o dist/myapp-linux-amd64
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$VERSION" -o dist/myapp-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w -X main.Version=$VERSION" -o dist/myapp-darwin-arm64
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.Version=$VERSION" -o dist/myapp-windows-amd64.exe
```

### Test with Coverage and Race Detection

```bash
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -html=coverage.out -o coverage.html
```

### Profile-Guided Optimization (PGO)

```bash
# Step 1: Build with profiling
go build -o myapp

# Step 2: Run and collect profile
./myapp -cpuprofile=cpu.pprof

# Step 3: Rebuild with PGO
go build -pgo=cpu.pprof -o myapp
```

### Debug Build

```bash
go build -gcflags="all=-N -l" -o myapp-debug
dlv exec ./myapp-debug
```

### Static Binary (No Dependencies)

```bash
CGO_ENABLED=0 GOOS=linux go build \
  -a -installsuffix cgo \
  -ldflags="-s -w -extldflags '-static'" \
  -o myapp
```

### Analyze Escape Analysis

```bash
go build -gcflags="-m -m" 2>&1 | grep "escapes to heap"
```

### Check What Would Be Built

```bash
go list -f '{{.ImportPath}}' ./...
go list -f '{{.Deps}}' ./...
go list -json ./... | jq '.Deps'
```

---

## Quick Reference Table

| Command    | Common Flags                                | Example                            |
| ---------- | ------------------------------------------- | ---------------------------------- |
| `go build` | `-o`, `-v`, `-race`, `-ldflags`, `-gcflags` | `go build -o app -ldflags="-s -w"` |
| `go test`  | `-v`, `-run`, `-bench`, `-cover`, `-race`   | `go test -v -cover -race ./...`    |
| `go run`   | Same as build                               | `go run -race main.go`             |
| `go get`   | `-u`, `-d`, `-t`                            | `go get -u ./...`                  |
| `go mod`   | `tidy`, `download`, `verify`                | `go mod tidy -v`                   |
| `go list`  | `-f`, `-json`, `-m`                         | `go list -json -m all`             |
| `go vet`   | `-x`                                        | `go vet ./...`                     |
| `go fmt`   | `-x`                                        | `go fmt ./...`                     |

---

## Resources

- [Go Command Documentation](https://pkg.go.dev/cmd/go)
- [Go Build Constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
- [Go Test Flags](https://pkg.go.dev/cmd/go#hdr-Testing_flags)
- [Compiler Directives](https://pkg.go.dev/cmd/compile)
- [Linker Flags](https://pkg.go.dev/cmd/link)
- [Go Environment Variables](https://pkg.go.dev/cmd/go#hdr-Environment_variables)
