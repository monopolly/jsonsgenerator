package main

import (
	"fmt"
	"strings"
)

// парсер для изначальной структуры в []intefacr{}
func (a *Golang) SQLStructToTuple() []byte {

	lines := []string{
		"\n",
		"//Tuple create an array from struct",
		fmt.Sprintf("func (a *%s) sqlTuple() (r []any){ ", settings.Go.StructName),
	}
	var fieldlist []string
	for _, x := range fields {
		if x.SQL.NoInsert || x.SQL.Skip {
			continue
		}
		fieldlist = append(fieldlist, fmt.Sprintf(`a.%s`, x.Go.Name))
	}
	lines = append(lines, fmt.Sprintf("return []any{%s}", strings.Join(fieldlist, ",")))
	lines = append(lines, "}")

	return []byte(strings.Join(lines, "\n"))
}
