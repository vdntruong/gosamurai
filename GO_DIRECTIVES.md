# Go Directives Reference

A comprehensive guide to all special Go directives (magic comments) that control compiler behavior, code generation, and build constraints.

## Table of Contents

- [Build Constraints](#build-constraints)
- [Code Generation](#code-generation)
- [Static Assets](#static-assets)
- [Compiler Optimizations](#compiler-optimizations)
- [Memory Management](#memory-management)
- [Low-Level Directives](#low-level-directives)
- [Testing & Debugging](#testing--debugging)

---

## Build Constraints

### `//go:build`

**Since:** Go 1.17  
**Purpose:** Conditional compilation based on build tags, OS, architecture, and other conditions.

**Syntax:**
```go
//go:build <expression>

package mypackage
```

**Examples:**
```go
// Platform-specific
//go:build linux
//go:build windows
//go:build darwin

// Architecture-specific
//go:build amd64
//go:build arm64
//go:build 386

// Go version constraints
//go:build go1.18
//go:build go1.21

// Logical operators
//go:build linux && amd64           // AND
//go:build linux || darwin          // OR
//go:build !windows                 // NOT
//go:build (linux || darwin) && amd64

// Custom tags
//go:build integration             // Only when: go test -tags=integration
//go:build debug                    // Only when: go build -tags=debug
//go:build tools                    // Never compiled (tracks dependencies)
```

**Common Use Cases:**
- Platform-specific implementations
- Feature flags
- Test builds
- Tool dependency tracking

**Notes:**
- Must be at the top of the file before `package`
- Replaces the older `// +build` syntax
- If both `//go:build` and `// +build` exist, `//go:build` takes precedence
- No space between `//` and `go:`

---

### `// +build` (Legacy)

**Purpose:** Old syntax for build constraints (still supported but deprecated).

**Syntax:**
```go
// +build <tags>

package mypackage
```

**Examples:**
```go
// +build linux
// +build linux,amd64    // AND (comma)
// +build linux darwin   // OR (space)
// +build !windows       // NOT
```

**Migration:**
```go
// Old:
// +build linux,amd64

// New:
//go:build linux && amd64
```

**Notes:**
- Deprecated in favor of `//go:build`
- If both are present, `//go:build` must come first
- Use `gofmt` to auto-generate compatible `// +build` lines

---

## Code Generation

### `//go:generate`

**Purpose:** Run commands during `go generate` to auto-generate code.

**Syntax:**
```go
//go:generate <command> [arguments...]
```

**Examples:**
```go
// Generate String() methods for enums
//go:generate stringer -type=Status

// Generate mocks
//go:generate mockgen -source=interface.go -destination=mock_test.go

// Run multiple generators
//go:generate go run generate_models.go
//go:generate gofmt -w generated_models.go

// With environment variables
//go:generate go run -ldflags "-X main.Version=$VERSION" gen.go

// Custom scripts
//go:generate bash scripts/generate_proto.sh
```

**Special Variables:**
- `$GOARCH` - target architecture
- `$GOOS` - target operating system
- `$GOFILE` - basename of the current file
- `$GOLINE` - line number of the directive
- `$GOPACKAGE` - name of the package
- `$DOLLAR` - literal `$` character

**Example with variables:**
```go
//go:generate echo "Generating for $GOFILE in package $GOPACKAGE"
```

**Run generation:**
```bash
go generate ./...              # Generate for all packages
go generate .                  # Current package only
go generate ./path/to/package  # Specific package
```

**Notes:**
- Commands run in the package directory
- Can appear anywhere in `.go` files
- Commonly placed near the types they generate code for
- Not run automatically during `go build` - must explicitly run `go generate`

---

## Static Assets

### `//go:embed`

**Since:** Go 1.16  
**Purpose:** Embed files and directories into the compiled binary.

**Syntax:**
```go
//go:embed <pattern>
var varName <type>
```

**Requires:**
```go
import "embed"
```

**Examples:**

**Single file:**
```go
import _ "embed"

//go:embed version.txt
var version string

//go:embed logo.png
var logo []byte
```

**Multiple files:**
```go
import "embed"

//go:embed templates/*.html
var templates embed.FS

//go:embed static/css/*.css static/js/*.js
var assets embed.FS
```

**Entire directory:**
```go
import "embed"

//go:embed all:static
var staticFiles embed.FS
```

**Usage:**
```go
// Read embedded file
content, _ := templates.ReadFile("templates/index.html")

// Use with http.FileServer
http.Handle("/static/", http.FileServer(http.FS(staticFiles)))

// Use with template parsing
tmpl, _ := template.ParseFS(templates, "templates/*.html")
```

**Patterns:**
- `*.txt` - all .txt files in current directory
- `dir/*.go` - all .go files in dir
- `dir/**/*` - all files recursively
- `all:dir` - includes dot files (`.gitignore`, etc.)

**Notes:**
- Files must be in the same module
- Cannot embed files outside module root
- Embedded files are compressed automatically
- Can only embed into `string`, `[]byte`, or `embed.FS`
- Variable must be package-level (not local)

---

## Compiler Optimizations

### `//go:noinline`

**Purpose:** Prevent the compiler from inlining a function.

**Syntax:**
```go
//go:noinline
func functionName() {
    // ...
}
```

**Example:**
```go
//go:noinline
func expensiveOperation() int {
    // Force this to be a separate function call
    // Useful for benchmarking or debugging
    return compute()
}
```

**Use Cases:**
- Benchmarking specific functions
- Debugging with proper stack traces
- Prevent optimization that might hide bugs
- Control binary size

---

### `//go:inline`

**Purpose:** Hint to the compiler that a function should be inlined (not guaranteed).

**Syntax:**
```go
//go:inline
func functionName() {
    // ...
}
```

**Example:**
```go
//go:inline
func fastPath(x int) int {
    return x * 2
}
```

**Notes:**
- Only a hint, compiler may ignore it
- Inlining happens automatically for small functions
- Less commonly needed than `//go:noinline`

---

## Memory Management

### `//go:noescape`

**Purpose:** Assert that pointer arguments don't escape to the heap.

**Syntax:**
```go
//go:noescape
func externalFunc(ptr *Type)
```

**Example:**
```go
//go:noescape
func libc_malloc(size uintptr) unsafe.Pointer

//go:noescape
func memcpy(dst, src unsafe.Pointer, n uintptr)
```

**Use Cases:**
- Assembly function declarations
- External (C) function declarations
- Performance optimization by keeping data on stack

**Notes:**
- Typically used with `//go:linkname`
- Helps escape analysis optimization
- Incorrect usage can lead to memory corruption
- Used in runtime and syscall packages

---

### `//go:nosplit`

**Purpose:** Disable stack overflow checking for a function.

**Syntax:**
```go
//go:nosplit
func functionName() {
    // ...
}
```

**Example:**
```go
//go:nosplit
func criticalRuntimeFunc() {
    // Runtime code that can't check stack
}
```

**Use Cases:**
- Low-level runtime code
- Code that runs during stack growth
- Signal handlers

**Notes:**
- Very dangerous if misused
- Primarily for Go runtime internals
- Can cause silent stack overflow if function uses too much stack
- Maximum stack usage is ~800 bytes

---

### `//go:notinheap`

**Purpose:** Declare that a type must not be allocated on the heap.

**Syntax:**
```go
//go:notinheap
type TypeName struct {
    // ...
}
```

**Example:**
```go
//go:notinheap
type notInHeap struct {
    data [4096]byte
}
```

**Use Cases:**
- Runtime internals
- Memory-mapped structures
- Performance-critical code

**Notes:**
- Primarily for Go runtime
- Cannot contain pointers to regular heap objects
- Cannot be used as a type argument for generic types

---

### `//go:uintptrescapes`

**Purpose:** Indicate that `uintptr` arguments represent pointers that escape.

**Syntax:**
```go
//go:uintptrescapes
func functionName(ptr uintptr) {
    // ...
}
```

**Example:**
```go
//go:uintptrescapes
func syscallWithPointer(ptr uintptr, size int) error {
    // Tells compiler that ptr is actually a pointer
    return syscall.Syscall(SYS_SOMETHING, ptr, uintptr(size), 0)
}
```

**Notes:**
- Helps garbage collector track pointers disguised as uintptr
- Used in syscall packages
- Advanced usage only

---

## Low-Level Directives

### `//go:linkname`

**Purpose:** Link to private symbols in other packages or runtime.

**Syntax:**
```go
//go:linkname localname importpath.name
```

**Requires:**
```go
import _ "unsafe"
```

**Examples:**
```go
package mypackage

import _ "unsafe"

//go:linkname nanotime runtime.nanotime
func nanotime() int64

//go:linkname fastrand runtime.fastrand
func fastrand() uint32
```

**Real-world example:**
```go
package main

import (
    _ "unsafe"
)

//go:linkname runtimeNano runtime.nanotime
func runtimeNano() int64

func main() {
    nano := runtimeNano()
    println(nano)
}
```

**Use Cases:**
- Access runtime internals
- Testing private functions
- Performance hacks

**⚠️ Warnings:**
- Breaks encapsulation
- No API stability guarantees
- Can break across Go versions
- Use only as last resort
- Requires `import _ "unsafe"`

---

### `//go:nowritebarrier`

**Purpose:** Assert that a function does not contain write barriers.

**Syntax:**
```go
//go:nowritebarrier
func functionName() {
    // ...
}
```

**Use Cases:**
- Garbage collector implementation
- Runtime internals

**Notes:**
- Extremely low-level
- Used only in Go runtime
- Compiler will error if function needs write barriers

---

### `//go:nowritebarrierrec`

**Purpose:** Like `//go:nowritebarrier` but applies recursively to all called functions.

**Syntax:**
```go
//go:nowritebarrierrec
func functionName() {
    // ...
}
```

**Notes:**
- Even more restrictive than `//go:nowritebarrier`
- Runtime internals only

---

### `//go:yeswritebarrierrec`

**Purpose:** Cancel `//go:nowritebarrierrec` from caller.

**Syntax:**
```go
//go:yeswritebarrierrec
func functionName() {
    // ...
}
```

**Notes:**
- Allows write barriers in functions called from `//go:nowritebarrierrec` context
- Runtime internals only

---

## Testing & Debugging

### `//go:norace`

**Purpose:** Disable race detector for a specific function.

**Syntax:**
```go
//go:norace
func functionName() {
    // ...
}
```

**Example:**
```go
//go:norace
func unsafeButCorrect() {
    // Intentional data race that's actually safe
    // Race detector would flag this, but we know it's OK
}
```

**Use Cases:**
- Benign data races
- Performance-critical code
- Lock-free algorithms

**Notes:**
- Use sparingly and document why it's safe
- Prefer fixing race conditions over using this
- Only affects race detector, not actual behavior

---

## Pragma Directives Summary

| Directive              | Stability     | Common Usage          |
| ---------------------- | ------------- | --------------------- |
| `//go:build`           | ✅ Stable      | Production code       |
| `//go:generate`        | ✅ Stable      | Production code       |
| `//go:embed`           | ✅ Stable      | Production code       |
| `//go:noinline`        | ⚠️ Semi-stable | Debugging, benchmarks |
| `//go:inline`          | ⚠️ Semi-stable | Performance tuning    |
| `//go:noescape`        | ⚠️ Semi-stable | Low-level code        |
| `//go:linkname`        | ❌ Unstable    | Hacks, avoid          |
| `//go:nosplit`         | ❌ Unstable    | Runtime only          |
| `//go:notinheap`       | ❌ Unstable    | Runtime only          |
| `//go:norace`          | ⚠️ Semi-stable | Race detector tuning  |
| `//go:nowritebarrier*` | ❌ Unstable    | Runtime only          |

## Best Practices

### ✅ Safe to Use
- `//go:build` - Build constraints
- `//go:generate` - Code generation
- `//go:embed` - Static file embedding

### ⚠️ Use with Caution
- `//go:noinline` - For benchmarking/debugging
- `//go:norace` - Only when you're certain it's safe

### ❌ Avoid Unless You Know What You're Doing
- `//go:linkname` - Breaks API guarantees
- `//go:nosplit` - Can cause crashes
- `//go:notinheap` - Runtime internals
- `//go:nowritebarrier*` - GC internals

## Important Rules

1. **No space after `//`**: `//go:build` not `// go:build`
2. **Placement matters**: Most directives go right before what they affect
3. **Build tags go first**: Before `package` declaration
4. **One per line**: Can't combine directives on same line
5. **Must be comments**: Changing to `/* */` won't work

## Resources

- [Go Build Constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
- [Go Generate](https://pkg.go.dev/cmd/go#hdr-Generate_Go_files_by_processing_source)
- [Embed Package](https://pkg.go.dev/embed)
- [Compiler Directives](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ir/node.go)
- [Go Compiler Source](https://github.com/golang/go/tree/master/src/cmd/compile)

## Examples Repository

For more examples of these directives in action, check:
- [Go Runtime Source](https://github.com/golang/go/tree/master/src/runtime)
- [Go Syscall Package](https://github.com/golang/go/tree/master/src/syscall)
- [Go Testing Package](https://github.com/golang/go/tree/master/src/testing)
