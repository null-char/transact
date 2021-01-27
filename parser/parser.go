package parser

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// Parser parses input from the given io.Reader by separating it into an operation and its
// respective arguments if any. It is also capable of parsing text into a Mappable value (String / Number).
type Parser struct {
	reader *bufio.Reader
}

// MakeParser constructs a new parser from the given io.Reader
func MakeParser(rd io.Reader) *Parser {
	reader := bufio.NewReader(rd)
	return &Parser{reader}
}

// Run reads input from stdin and then parses it into an operation and its arguments (if any). Returns an error if parsed input is empty.
func (p Parser) Run() (string, []string, error) {
	input, _ := p.reader.ReadString(byte('\n'))
	xs := strings.Fields(input)
	if len(xs) == 0 {
		return "", nil, errors.New("ERROR: Empty input")
	}

	operation := xs[0]
	var args []string
	if len(xs) > 1 {
		args = append(args, xs[1:]...)
	}

	return operation, args, nil
}
