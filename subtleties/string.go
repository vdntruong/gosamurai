package subtleties

import "fmt"

// We can do indexed based string interpolation in Go

func IndexedBasedString() {
	fmt.Printf("%[1]s %[2]s %[2]s\n", "Hello", "World", "!")
}

// IndexedBasedString Output: Hello World World
