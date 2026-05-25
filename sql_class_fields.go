package main

import (
	"fmt"
	"strings"
)

// get list
func (a *SQL) Fields() []byte {

	var list []string
	list = append(list, "\n//parse sql query")
	list = append(list, fmt.Sprintf(`func (a *%s) Fields() (res []%s) {`, settings.SQL.Class, settings.IndexTypeName))

	var fieldlist []string
	for _, x := range fields {
		if x.SQL.Skip {
			continue
		}
		fieldlist = append(fieldlist, x.Go.Index)
	}
	list = append(list, fmt.Sprintf(`return []%s { %s }`, settings.IndexTypeName, strings.Join(fieldlist, ",")))
	list = append(list, "}")
	return []byte(strings.Join(list, "\n"))
}
