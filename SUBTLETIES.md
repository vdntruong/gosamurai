# Go Subtleties, Gotchas, and Lesser-Known Features

A collection of subtle behaviors, edge cases, and lesser-known features in Go that developers should be aware of.

## Table of Contents

- [Variables and Types](#variables-and-types)
- [Slices and Arrays](#slices-and-arrays)
- [Maps](#maps)
- [Strings and Runes](#strings-and-runes)
- [Interfaces](#interfaces)
- [Methods and Receivers](#methods-and-receivers)
- [Defer, Panic, Recover](#defer-panic-recover)
- [Channels and Concurrency](#channels-and-concurrency)
- [For Loops](#for-loops)
- [Struct Embedding](#struct-embedding)
- [JSON and Encoding](#json-and-encoding)
- [Error Handling](#error-handling)
- [Testing](#testing)
- [Build and Tooling](#build-and-tooling)
- [Standard Library Surprises](#standard-library-surprises)

---

## Variables and Types

### 1. Short Variable Declaration Redeclaration

You can "redeclare" a variable with `:=` if at least one new variable is also declared:

```go
// This works!
x := 1
x, y := 2, 3  // x is reassigned, y is newly declared

fmt.Println(x, y)  // Output: 2 3
```

But be careful with shadowing:

```go
x := 1
if true {
    x, y := 2, 3  // New x in this scope! Shadows outer x
    fmt.Println(x, y)  // 2 3
}
fmt.Println(x)  // 1 (outer x unchanged)
```

### 2. Nil Slices vs Empty Slices

```go
var s1 []int        // nil slice
s2 := []int{}       // empty slice (not nil)
s3 := make([]int, 0) // empty slice (not nil)

fmt.Println(s1 == nil)  // true
fmt.Println(s2 == nil)  // false
fmt.Println(s3 == nil)  // false

// But they all have len 0!
fmt.Println(len(s1), len(s2), len(s3))  // 0 0 0

// JSON marshaling difference
json.Marshal(s1)  // "null"
json.Marshal(s2)  // "[]"
```

**Best practice:** Use `var s []int` for nil slices, return nil from functions unless you specifically need non-nil empty.

### 3. Untyped Constants Are Super Flexible

```go
const maxSize = 1 << 32  // This fits in int64 but not int32

var x int32 = maxSize / 2  // OK! Constant division at compile time
var y int32 = maxSize      // Compile error: overflow

// Constants can be huge
const huge = 1 << 100  // OK!
// var x int = huge     // Error: doesn't fit
```

### 4. Type Aliases vs Type Definitions

```go
// Type alias (Go 1.9+)
type MyInt = int  // MyInt and int are the same type

// Type definition
type MyInt2 int   // MyInt2 is a different type from int

var a int = 1
var b MyInt = a   // OK! Same type
var c MyInt2 = a  // Error! Different types
```

### 5. Generic Type Constraints with ~ (Go 1.18+)

The `~` operator constrains based on underlying type:

```go
type SomeConstantType string
const someConstant SomeConstantType = "foo"

// Accepts any type whose underlying type is string
func buildMessage[T ~string](value T) string {
    return fmt.Sprintf("The underlying string value is: '%s'", value)
}

msg := buildMessage(someConstant)  // Works!
```

This is useful for typed constants (like enums from other languages).


### 6. Unused Variables Are Errors, But...

```go
// Error: unused variable
func bad() {
    x := 1  // declared but not used
}

// OK: use the blank identifier
func good() {
    x := getValue()
    _ = x  // Explicitly ignored
}

// Package-level variables can be unused
var unused = 42  // OK at package level
```

---

## Slices and Arrays

---

## Slices and Arrays

### 7. Array Values Are Copied

```go
a := [3]int{1, 2, 3}
b := a  // Copies the entire array!
b[0] = 99

fmt.Println(a[0])  // 1 (unchanged)
fmt.Println(b[0])  // 99

// Slices are references
s1 := []int{1, 2, 3}
s2 := s1  // Same underlying array
s2[0] = 99
fmt.Println(s1[0])  // 99 (changed!)
```

### 8. Slice Append Can Reallocate

```go
s1 := []int{1, 2, 3}
s2 := s1

s1 = append(s1, 4)  // May or may not reallocate
s1[0] = 99

// s2 may or may not be affected!
fmt.Println(s2[0])  // Could be 1 or 99, depending on capacity
```

**Always capture append's return value:**
```go
s = append(s, value)  // Good
```

### 9. Slice Bounds Don't Prevent Modification

```go
s := make([]int, 3, 10)  // len=3, cap=10
s[0], s[1], s[2] = 1, 2, 3

s2 := s[:2]    // len=2, cap=10
s3 := s2[:cap(s2)]  // Extend to full capacity!

s3[2] = 99     // Modifies s[2]!
fmt.Println(s)  // [1 2 99]
```

### 10. Nil Slices Are Usable

```go
var s []int  // nil

// All these work on nil slices!
len(s)     // 0
cap(s)     // 0
s = append(s, 1)  // OK! Creates new slice
for range s {}     // OK! Zero iterations

// But dereferencing panics
// s[0] = 1  // panic: index out of range
```

### 11. Three-Index Slicing

Control the capacity of subslices:

```go
s := []int{1, 2, 3, 4, 5}

// Normal slice: s[low:high]
s1 := s[1:3]  // [2 3], cap=4

// Three-index slice: s[low:high:max]
s2 := s[1:3:3]  // [2 3], cap=2 (max-low)

// s2 can't grow beyond index 3
s2 = append(s2, 99)  // Forces new allocation, doesn't affect s
```

---

## Maps

---

## Maps

### 12. Map Iteration Order Is Random

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}

// Order is intentionally randomized each time!
for k, v := range m {
    fmt.Println(k, v)  // Order varies
}
```

**Use sorted keys for deterministic order:**
```go
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Println(k, m[k])
}
```

### 13. Map Values Are Not Addressable

```go
type Point struct{ X, Y int }
m := map[string]Point{"p": {1, 2}}

// Error: cannot assign to struct field
// m["p"].X = 10

// Must reassign whole value
p := m["p"]
p.X = 10
m["p"] = p

// Or use pointer values
m2 := map[string]*Point{"p": {1, 2}}
m2["p"].X = 10  // OK!
```

### 14. Reading Missing Keys Returns Zero Value

```go
m := make(map[string]int)

// Returns 0, not an error!
x := m["missing"]  
fmt.Println(x)  // 0

// Check existence with two-value form
if val, ok := m["key"]; ok {
    fmt.Println("Found:", val)
} else {
    fmt.Println("Not found")
}
```

### 15. Nil Map Reads Work, Writes Panic

```go
var m map[string]int  // nil map

// Reading is OK
x := m["key"]     // 0
_, ok := m["key"] // false, ok
len(m)            // 0

// Writing panics!
// m["key"] = 1  // panic: assignment to entry in nil map

// Always initialize
m = make(map[string]int)
m["key"] = 1  // OK
```

---

## Interfaces

### 16. Map Modification During Range May Not Appear

```go
m := map[int]int{1: 1, 2: 2, 3: 3}

for key, value := range m {
    fmt.Printf("%d = %d\n", key, value)
    if key == 1 {
        for i := 10; i < 20; i++ {
            m[i] = i * 10  // Add entries during iteration
        }
    }
}
// New entries MAY OR MAY NOT appear in this iteration!
```

Go hashes keys into buckets. If the bucket was already visited, new entries won't appear. This is for speed, unlike Python's stable insertion order.

---

## Strings and Runes

### 17. Strings Are Immutable Byte Slices

```go
s := "hello"
// s[0] = 'H'  // Error: cannot assign

// Must create new string
s = "H" + s[1:]  // "Hello"

// Strings are just []byte under the hood
b := []byte(s)  // Copies to byte slice
b[0] = 'h'
s = string(b)   // Copies back
```

### 18. String Length Is Bytes, Not Characters

```go
s := "hello"
fmt.Println(len(s))  // 5

s = "こんにちは"  // Japanese
fmt.Println(len(s))  // 15 (UTF-8 bytes!)

// Use rune count for characters
fmt.Println(utf8.RuneCountInString(s))  // 5
```

### 19. Range Over String Iterates Runes

```go
s := "Go语言"

// By index: gets bytes
for i := 0; i < len(s); i++ {
    fmt.Printf("%x ", s[i])
}
// Output: 47 6f e8 af ad e8 a8 80

// By range: gets runes (characters)
for i, r := range s {
    fmt.Printf("%d:%c ", i, r)
}
// Output: 0:G 1:o 2:语 5:言
// Note: indices jump! (2 -> 5)
```

**Invalid UTF-8 Replacement:**
```go
invalidBytes := []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F, 0xFF} // "Hello" + invalid byte
s := string(invalidBytes)

for _, r := range s {
    fmt.Printf("%c ", r)  // Prints: H e l l o �
}
// Invalid UTF-8 bytes are replaced with replacement character �
```

### 20. String Concatenation with Builder

```go
// Inefficient
var s string
for i := 0; i < 1000; i++ {
    s += "x"  // Creates new string each time
}

// Efficient
var sb strings.Builder
for i := 0; i < 1000; i++ {
    sb.WriteString("x")
}
s := sb.String()
```

---

## Struct Embedding

### 21. Index-Based String Interpolation with fmt

```go
// Reuse arguments in fmt
fmt.Printf("%[1]s %[1]s %[2]s %[2]s %[3]s", "one", "two", "three")
// Output: "one one two two three"

// Useful for reducing repetition
fmt.Printf("Error: %[1]s (code: %[2]d). Please fix %[1]s", "connection failed", 500)
```

---

## Interfaces

### 22. Interface Nil Is Not Always Nil

```go
var p *int = nil
var i interface{} = p

fmt.Println(p == nil)  // true
fmt.Println(i == nil)  // false! (interface holds type info)

// The interface is (type *int, value nil)
```

**Check for real nil:**
```go
func isNil(i interface{}) bool {
    if i == nil {
        return true
    }
    v := reflect.ValueOf(i)
    return v.Kind() == reflect.Ptr && v.IsNil()
}
```

### 23. Empty Interface Accepts Anything

```go
func print(v interface{}) {
    fmt.Println(v)
}

print(42)
print("hello")
print(struct{}{})
print(nil)  // All valid!
```

Go 1.18+ use `any` instead:
```go
func print(v any) {  // any is alias for interface{}
    fmt.Println(v)
}
```

### 24. Type Switches and Type Assertions

```go
var x interface{} = "hello"

// Type assertion
s := x.(string)  // OK
// n := x.(int)  // panic!

// Safe type assertion
if s, ok := x.(string); ok {
    fmt.Println("String:", s)
}

// Type switch
switch v := x.(type) {
case string:
    fmt.Println("String:", v)
case int:
    fmt.Println("Int:", v)
default:
    fmt.Println("Unknown:", v)
}
```

---

## Methods and Receivers

### 25. Method Sets: Values vs Pointers

```go
type Counter struct{ n int }

func (c Counter) Get() int       { return c.n }
func (c *Counter) Increment()    { c.n++ }

// Pointer implements both value and pointer methods
var pc *Counter = &Counter{}
pc.Get()        // OK
pc.Increment()  // OK

// Value only implements value methods
var c Counter = Counter{}
c.Get()         // OK
// c.Increment() // Error if Counter is interface!

// Interface assignment
type Incrementer interface{ Increment() }
var i Incrementer = &Counter{}  // OK
// var i Incrementer = Counter{} // Error! Value doesn't implement
```

---

## Methods and Receivers

### 26. Methods on Nil Receivers

```go
type Tree struct {
    value int
    left  *Tree
    right *Tree
}

func (t *Tree) Sum() int {
    if t == nil {
        return 0  // Nil receiver is OK!
    }
    return t.value + t.left.Sum() + t.right.Sum()
}

var t *Tree  // nil
sum := t.Sum()  // Works! Returns 0
```

### 27. Embedding Promotes Methods

```go
type Reader struct{}
func (r Reader) Read() string { return "reading" }

type Writer struct{}
func (w Writer) Write() string { return "writing" }

type ReadWriter struct {
    Reader  // Embedded
    Writer  // Embedded
}

rw := ReadWriter{}
rw.Read()   // Promoted from Reader
rw.Write()  // Promoted from Writer
```

---

## Defer, Panic, Recover

---

## Defer, Panic, Recover

### 28. Defer Arguments Are Evaluated Immediately

```go
func trace(s string) string {
    fmt.Println("entering:", s)
    return s
}

func a() {
    defer trace("a")()  // "entering: a" printed immediately!
    fmt.Println("in a")
}

// Output:
// entering: a
// in a
```

**Closure defers are evaluated later:**
```go
func b() {
    x := 1
    defer func() {
        fmt.Println(x)  // Reads x when defer runs
    }()
    x = 2
}
// Output: 2
```

### 29. Named Return Values and Defer

```go
func f() (result int) {
    defer func() {
        result++  // Modifies return value!
    }()
    return 0
}

fmt.Println(f())  // 1, not 0!
```

### 30. Defer Executes in LIFO Order

```go
func example() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
}
// Output: 3 2 1
```

### 31. Recover Only Works in Defer

```go
func safeDivide(a, b int) (result int) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
            result = 0
        }
    }()
    
    return a / b  // May panic
}

safeDivide(10, 0)  // Recovered: runtime error: integer divide by zero
```

---

## Channels and Concurrency

---

## Channels and Concurrency

### 32. Closed Channel Behavior

```go
ch := make(chan int, 2)
ch <- 1
ch <- 2
close(ch)

// Reading from closed channel returns zero value + false
v, ok := <-ch  // 1, true
v, ok = <-ch   // 2, true
v, ok = <-ch   // 0, false (closed!)

// Sending to closed channel panics
// ch <- 3  // panic: send on closed channel

// Closing closed channel panics
// close(ch)  // panic: close of closed channel
```

### 33. Nil Channel Blocks Forever

```go
var ch chan int  // nil

// Both block forever!
// <-ch    // blocks
// ch <- 1 // blocks

// Useful in select to disable a case
select {
case <-ch:  // Never selected if ch is nil
    // ...
case <-time.After(1 * time.Second):
    fmt.Println("timeout")
}
```

### 34. Select Chooses Randomly

```go
ch1 := make(chan int, 1)
ch2 := make(chan int, 1)
ch1 <- 1
ch2 <- 2

// If multiple cases are ready, one is chosen randomly!
select {
case v := <-ch1:
    fmt.Println("ch1:", v)
case v := <-ch2:
    fmt.Println("ch2:", v)
}
// Could print either!
```

### 35. For-Range on Channel

```go
ch := make(chan int)

go func() {
    for i := 0; i < 3; i++ {
        ch <- i
    }
    close(ch)  // Must close or loop never exits!
}()

// Reads until channel is closed
for v := range ch {
    fmt.Println(v)
}
```

---

## For Loops

### 36. Context Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()  // Always call cancel to release resources!

// Even if timeout occurs, call cancel to free resources
```

### 37. Context-Aware Functions with Select

Always select on context in channel operations:

```go
func sendSignal(ctx context.Context, ch chan<- string) {
    select {
    case <-time.After(5 * time.Second):
        ch <- "operation complete"
    case <-ctx.Done():
        // Without this, we'd wait 5 seconds even if cancelled!
        ch <- "operation cancelled"
    }
}

ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

go sendSignal(ctx, ch)
// Exits after 1 second, not 5!
```

**Bonus:** Context is canceled after HTTP handler completes, even for successful responses!

### 38. sync.WaitGroup.Go (Go 1.25+)

New convenience method for adding goroutines:

```go
var wg sync.WaitGroup

// Old way
wg.Add(1)
go func() {
    defer wg.Done()
    // work...
}()

// New way - shorter!
wg.Go(func() {
    // work...
})

wg.Wait()
```

Implementation automatically calls `Add(1)` and `defer Done()`.

---

## For Loops

### 39. Loop Variable Gotcha

```go
// Classic mistake
var wg sync.WaitGroup
for i := 0; i < 3; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println(i)  // All print 3!
    }()
}
wg.Wait()

// Fix: capture variable
for i := 0; i < 3; i++ {
    wg.Add(1)
    i := i  // Create new variable
    go func() {
        defer wg.Done()
        fmt.Println(i)  // Correct!
    }()
}

// Or pass as parameter
for i := 0; i < 3; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        fmt.Println(n)
    }(i)
}
```

**Note:** Go 1.22+ fixes this! Loop variables are per-iteration by default.

### 40. Range Directly Over Integers (Go 1.22+)

```go
// Go 1.22 introduced ranging over integers
for i := range 10 {
    fmt.Println(i)  // 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
}

// Useful for simple iterations
for i := range 5 {
    fmt.Println(i + 1)  // 1, 2, 3, 4, 5
}
```

This is cleaner than `for i := 0; i < 10; i++` when you just need the index.

### 41. Range Copies Values

```go
type Point struct{ X, Y int }
points := []Point{{1, 2}, {3, 4}}

// v is a copy!
for _, v := range points {
    v.X = 99  // Doesn't modify original
}
fmt.Println(points)  // [{1 2} {3 4}]

// Use index to modify
for i := range points {
    points[i].X = 99
}
```

### 42. Break and Continue with Labels

```go
outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if i == 1 && j == 1 {
            break outer  // Breaks outer loop!
        }
        fmt.Println(i, j)
    }
}

// Can also use with continue
outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if j == 1 {
            continue outer  // Continues outer loop
        }
        fmt.Println(i, j)
    }
}
```

---

## Strings and Runes

---

## Struct Embedding

### 43. Field Shadowing

```go
type Inner struct {
    X int
}

type Outer struct {
    Inner
    X int  // Shadows Inner.X
}

o := Outer{Inner: Inner{X: 1}, X: 2}
fmt.Println(o.X)        // 2 (Outer.X)
fmt.Println(o.Inner.X)  // 1 (Inner.X)
```

### 44. Embedded Interfaces

```go
type Reader interface {
    Read() string
}

type Writer interface {
    Write(string)
}

type ReadWriter interface {
    Reader  // Embeds Reader methods
    Writer  // Embeds Writer methods
}

// Equivalent to:
// type ReadWriter interface {
//     Read() string
//     Write(string)
// }
```

### 45. Hidden Interface Satisfaction with Embedding

Embedding promotes methods, which can unexpectedly satisfy interfaces:

```go
type Event struct {
    Name      string    `json:"name"`
    time.Time          // Embedded! Promotes all methods
}

event := Event{
    Name: "Launch",
    Time: time.Date(2023, time.November, 10, 23, 0, 0, 0, time.UTC),
}

jsonData, _ := json.Marshal(event)
fmt.Println(string(jsonData))  // "2023-11-10T23:00:00Z"
// Unexpected! time.Time's MarshalJSON method is called
```

The embedded `time.Time` has a `MarshalJSON()` method, so the whole struct uses that method instead of default marshaling.

**Solution:** Use a named field instead:
```go
type Event struct {
    Name      string    `json:"name"`
    Timestamp time.Time `json:"timestamp"`  // Named field
}
```

---

## JSON and Encoding

### 46. Empty Structs Take Zero Bytes

```go
var s struct{}
fmt.Println(unsafe.Sizeof(s))  // 0 bytes!

// All zero-sized allocations return same special address
// Great for signaling on channels
done := make(chan struct{})
done <- struct{}{}  // Sends signal, no data

// Also used for set implementation
set := make(map[string]struct{})
set["key"] = struct{}{}
```

Empty structs are more memory-efficient than booleans which occupy space.

---

## JSON and Encoding

### 47. Unexported Fields Not Marshaled

```go
type User struct {
    Name     string  // Exported: marshaled
    password string  // Unexported: ignored
}

u := User{Name: "Alice", password: "secret"}
json.Marshal(u)  // {"Name":"Alice"}
```

### 48. JSON Tags Control Marshaling

```go
type User struct {
    ID        int    `json:"id"`
    Name      string `json:"name,omitempty"`  // Omit if empty
    Password  string `json:"-"`               // Never marshal
    CreatedAt time.Time `json:"created_at"`
}
```

### 49. Zero Values and omitempty

```go
type Config struct {
    Port   int    `json:"port,omitempty"`
    Debug  bool   `json:"debug,omitempty"`
    Name   string `json:"name,omitempty"`
}

c := Config{}
json.Marshal(c)  // {}

c = Config{Port: 0, Debug: false, Name: ""}
json.Marshal(c)  // {} (all omitted!)

// Be careful with zero values!
c = Config{Port: 8080}
json.Marshal(c)  // {"port":8080}
```

---

## Testing

---

## Error Handling

### 50. Custom Error Types with errors.As

```go
type MyError struct {
    Message string
    Code    int
}

func (e *MyError) Error() string {
    return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}

func someFunction() error {
    return &MyError{Message: "something went wrong", Code: 404}
}

// Check for typed error
err := someFunction()
if err != nil {
    var myErr *MyError
    if errors.As(err, &myErr) {
        fmt.Printf("Custom error code: %d\n", myErr.Code)
    }
}
```

Attach structured data to errors for better debugging.

---

## Testing

### 51. Table-Driven Tests

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 1, 2, 3},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", 
                    tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

### 52. Testing Unexported Functions

```go
// In foo.go
package foo

func helper() int {
    return 42
}

// In foo_test.go - same package!
package foo  // Not foo_test

import "testing"

func TestHelper(t *testing.T) {
    result := helper()  // Can access unexported!
    if result != 42 {
        t.Error("expected 42")
    }
}
```

### 53. Test Main for Setup/Teardown

```go
func TestMain(m *testing.M) {
    // Setup
    fmt.Println("Setting up...")
    
    // Run tests
    code := m.Run()
    
    // Teardown
    fmt.Println("Tearing down...")
    
    os.Exit(code)
}
```

---

## Build and Tooling

---

## Build and Tooling

### 54. Build Tags

```go
// +build linux

package main  // Only built on Linux

// Multiple tags
// +build linux darwin
// +build amd64

// Negation
// +build !windows

// Go 1.17+ syntax
//go:build linux && amd64
```

### 55. go:generate Directive

```go
//go:generate stringer -type=Status

type Status int

const (
    StatusOK Status = iota
    StatusError
)

// Run: go generate ./...
```

### 56. Conditional Compilation Tricks

```go
const debug = false

func log(msg string) {
    if debug {
        fmt.Println(msg)
    }
}

// Compiler eliminates the if block when debug is false!
```

---

## Standard Library Surprises

### 57. The embed Package

Embed files directly into your binary:

```go
import _ "embed"

//go:embed version.txt
var version string

//go:embed templates/*
var templates embed.FS

// No need to read from disk at runtime!
// Great for HTML, JS, CSS, images, etc.
```

**Benefits:**
- Single binary deployment
- No external file dependencies
- Faster startup (no disk I/O)

### 58. Renaming Packages on Import

```go
// Regular import
import "github.com/very/long/package/name"

// Renamed import
import shortname "github.com/very/long/package/name"

// Now use shortname instead
shortname.SomeFunction()

// Useful for:
// - Long package names
// - Package name conflicts
// - Mocking in tests
```

**Bonus:** Go's LSP can rename packages and even the directory!

---


---

---

## Standard Library Surprises

### 59. time.After Leaks

```go
// Bad: creates new timer each iteration
for {
    select {
    case <-time.After(1 * time.Second):
        // Timer is not garbage collected until it fires!
    }
}

// Good: reuse ticker
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()
for {
    select {
    case <-ticker.C:
        // ...
    }
}
```

### 60. Sort Stability

```go
// sort.Sort is NOT stable!
type Person struct {
    Name string
    Age  int
}

people := []Person{
    {"Alice", 25},
    {"Bob", 25},
}

sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})
// Order of Alice and Bob is not guaranteed!

// Use sort.SliceStable for stable sort
sort.SliceStable(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})
```

### 61. HTTP Client Connections

```go
// Default client doesn't set timeouts!
resp, err := http.Get(url)  // Can hang forever

// Always use timeouts
client := &http.Client{
    Timeout: 10 * time.Second,
}

// Or use context
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := client.Do(req)
```


### 62. Comparing Times with time.Equal

```go
t1 := time.Date(2024, 1, 1, 15, 0, 0, 0, time.UTC)
t2 := t1.In(time.FixedZone("EST", -5*3600))

// String comparison includes timezone - FAILS!
fmt.Println(t1.String() == t2.String())  // false

// Equal compares instants - WORKS!
fmt.Println(t1.Equal(t2))  // true

// Always use .Equal() for time comparisons!
```

Time.Equal reports whether times represent the same instant, regardless of timezone.

---

## Quick Reference

### Common Gotchas Checklist

- [ ] Loop variable capture in goroutines
- [ ] Nil interface is not nil
- [ ] Map iteration order is random
- [ ] Slice append may or may not reallocate
- [ ] defer arguments evaluated immediately
- [ ] String length is bytes, not runes
- [ ] Closed channel reads return zero value
- [ ] Method sets differ for values vs pointers
- [ ] Nil slices vs empty slices in JSON
- [ ] time.After leaks in loops

### When to Use What

| Use                | When                           |
| ------------------ | ------------------------------ |
| `var s []int`      | Want nil slice                 |
| `s := []int{}`     | Want non-nil empty slice       |
| `defer func(){}()` | Need to capture variable later |
| `defer fn()`       | Arguments evaluated now        |
| `sync.RWMutex`     | Many reads, few writes         |
| `sync.Mutex`       | Balanced reads/writes          |
| `sync.Map`         | Many concurrent accesses       |
| `map[K]V`          | Single-threaded or protected   |

---

## Resources

- [Go FAQ](https://go.dev/doc/faq)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Spec](https://go.dev/ref/spec)
- [Common Mistakes](https://github.com/golang/go/wiki/CommonMistakes)
- [Go Wiki](https://github.com/golang/go/wiki)

---

**Remember:** When in doubt, read the spec or use the Go Playground to experiment!

