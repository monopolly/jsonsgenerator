package main

import (
	"fmt"
	"strings"
)

// get struct value with function
func (a *Golang) Get() []byte {

	parse := []string{
		"\n",
		"//get struct value with function",
		fmt.Sprintf("func (a *%s) Get(k string) (v any) {", settings.Go.StructName),
		"switch k {",
		// fmt.Sprintf("switch %s(pos) {", settings.IndexTypeName),
	}

	for _, x := range fields {
		parse = append(parse, fmt.Sprintf(`case "%s": return a.%s //%s`, x.Name, x.Go.Name, x.Type))
	}
	parse = append(parse, "}")
	parse = append(parse, "return")
	parse = append(parse, "}")
	return []byte(strings.Join(parse, "\n"))
}
