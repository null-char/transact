package utils

import (
	"strconv"

	s "github.com/null-char/transact/store"
)

// ParseValue takes a string and parses it to a mappable value (either a Number or a String)
func ParseValue(v string) s.Mappable {
	x, err := strconv.Atoi(v)
	if err != nil {
		return s.String(v)
	}
	return s.Number(x)
}
