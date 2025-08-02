package main

import (
	"fmt"
	"strings"
)

// индекс это номера полей для json
func (a *Golang) generateParseTupleToJsons() []byte {

	parse := []string{
		"\n",
		"//Parse []any to json",
		fmt.Sprintf(`func ParseTupleTo%[1]s(r []any) (a %[1]s){`, settings.JS.Name),
		"for pos, x := range r {",
		fmt.Sprintf("switch %s(pos) {", settings.IndexTypeName),
	}

	for _, x := range fields {
		parse = append(parse, fmt.Sprintf(`case %s: a.Set(%s, x)`, x.Go.Index, x.FieldName))
	}

	parse = append(parse, "}}")
	parse = append(parse, "return")
	parse = append(parse, "}\n")

	return []byte(strings.Join(parse, "\n"))

}
