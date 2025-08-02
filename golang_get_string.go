package main

import (
	"fmt"
	"strings"
)

// get struct value with function
func (a *Golang) GetString() []byte {

	parse := []string{
		"\n",
		"//get any struct value as string",
		fmt.Sprintf("func (a *%s) String(k string) (v string) {", settings.Go.StructName),
		"switch k {",
		// fmt.Sprintf("switch %s(pos) {", settings.IndexTypeName),
	}

	for _, x := range fields {
		parse = append(parse, fmt.Sprintf(`case "%s": return fmt.Sprint(a.%s) //%s`, x.Name, x.Go.Name, x.Type))
	}
	parse = append(parse, "}")
	parse = append(parse, "return")
	parse = append(parse, "}")
	return []byte(strings.Join(parse, "\n"))
}
