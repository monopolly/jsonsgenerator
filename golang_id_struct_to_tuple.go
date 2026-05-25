package main

import (
	"fmt"
	"strings"
)

// parser from the original struct to []interface{}
func (a *Golang) IDStructToTuple() []byte {

	lines := []string{
		"\n",
		"//Tuple create an array from struct",
		fmt.Sprintf("func (a *%s) Tuple() (r []any){ ", settings.Go.StructName),
	}
	var fieldlist []string
	for _, x := range fields {
		fieldlist = append(fieldlist, fmt.Sprintf(`a.%s`, x.Go.Name))
	}
	lines = append(lines, fmt.Sprintf("return []any{%s}", strings.Join(fieldlist, ",")))
	lines = append(lines, "}")

	return []byte(strings.Join(lines, "\n"))
}
