package main

import (
	"fmt"
	"strings"
)

// get struct value with function
func (a *Golang) Map() []byte {
	var list []string
	list = append(list, "\n")
	list = append(list, "//struct to map")
	list = append(list, fmt.Sprintf("func (a *%s) Map() map[string]any{ ", settings.Go.StructName))
	list = append(list, "return map[string]any{")

	var name, goname string
	for _, x := range fields {

		switch x.Json.Name != "" {
		case true:
			name = x.Json.Name
		default:
			name = x.Name
		}

		goname = x.Go.Name

		list = append(list, fmt.Sprintf(`"%s": a.%s,`, name, goname))

	}
	list = append(list, "}")
	list = append(list, "}")
	list = append(list, "\n")

	return []byte(strings.Join(list, "\n"))
}
