package subtleties

import "fmt"

func BuildMessage[T ~string](message T) string {
	return fmt.Sprintf("%s", message)
}

// We can use the ~ operator to constrain a generic type signature. For instance, for typed constants
