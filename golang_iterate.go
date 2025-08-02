package main

import (
	"fmt"
	"strings"
)

// get struct value with function
func (a *Golang) Iterate() []byte {
	var list []string
	list = append(list, "\n")
	list = append(list, "//struct to map")
	list = append(list, fmt.Sprintf("func (a *%s) Iterate(f func(k %s, v any)) { ", settings.Go.StructName, settings.IndexTypeName))

	list = append(list, fmt.Sprintf("for _, x := range %sIndexes() {", settings.Go.StructName))
	list = append(list, "f(x, a.Get(x.String()))")
	list = append(list, "}")
	list = append(list, "}")
	list = append(list, "\n")

	return []byte(strings.Join(list, "\n"))
}
