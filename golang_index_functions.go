package main

import (
	"fmt"
	"strings"
)

func (a *Golang) IndexFunction() []byte {

	var list, index []string
	list = append(list, "//index func")
	list = append(list, fmt.Sprintf("func %sIndexes() []%s{", settings.Go.StructName, settings.IndexTypeName))
	for _, x := range fields {
		index = append(index, x.Go.Index)
	}

	in := fmt.Sprintf("[]%s{%s}", settings.IndexTypeName, strings.Join(index, ", "))
	list = append(list, fmt.Sprintf("return %s\n", in))
	list = append(list, "}\n\n")
	return []byte(strings.Join(list, "\n"))
}
