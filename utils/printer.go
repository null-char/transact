package utils

import (
	"fmt"

	s "github.com/null-char/transact/store"
)

// PrintValue formats and prints a mappable value
func PrintValue(v s.Mappable) {
	switch v := v.(type) {
	case s.Number:
		fmt.Printf("(integer) %d \n", v)
	case s.String:
		fmt.Printf("(string) %s \n", v)
	default:
		fmt.Printf("(unknown) %v \n", v)
	}
}
