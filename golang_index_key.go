package main

import (
	"fmt"
	"strings"
)

// field key to index int
func (a *Golang) IndexKeys() []byte {

	var list []string
	list = append(list, "\n")
	list = append(list, "//struct key to index")
	list = append(list, fmt.Sprintf("func %sKeyIndex(key string) %s {", settings.Go.StructName, settings.IndexTypeName))

	list = append(list, "switch key{")
	for _, x := range fields {
		list = append(list, fmt.Sprintf(`case "%s":`, x.Name))
		list = append(list, fmt.Sprintf(`return %s`, x.Go.Index))
	}
	list = append(list, "default: return 0")
	list = append(list, "}")
	list = append(list, "}\n")

	// valid field
	list = append(list, "\n")
	list = append(list, "//valid struct key check")
	list = append(list, fmt.Sprintf("func %sValidKey(key string) bool {", settings.Go.StructName))

	list = append(list, "switch key{")
	var valids []string
	for _, x := range fields {
		valids = append(valids, fmt.Sprintf(`"%s"`, x.Name))
	}
	list = append(list, fmt.Sprintf(`case %s:`, strings.Join(valids, ",")))
	list = append(list, `return true`)
	list = append(list, "default: return false")
	list = append(list, "}")
	list = append(list, "}\n")

	return []byte(strings.Join(list, "\n"))
}
