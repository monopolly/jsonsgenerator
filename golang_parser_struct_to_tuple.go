package main

import (
	"fmt"
	"strings"
)

// парсер для изначальной структуры в []intefacr{}
func (a *Golang) generateParserStructToTuple() []byte {

	lens := len(fields)

	lines := []string{
		"\n",
		"//Full tuple create an array from struct",
		fmt.Sprintf("func (a *%s) Tuple() (r []any){ ", settings.Go.StructName),
		fmt.Sprintf(`r = make([]any, %d)`, lens),
	} //
	for _, x := range fields {
		lines = append(lines, fmt.Sprintf(`r[%s] = a.%s`, x.Go.Index, x.Go.Name))
	}
	lines = append(lines, "return")
	lines = append(lines, "}\n")

	return []byte(strings.Join(lines, "\n"))

}
