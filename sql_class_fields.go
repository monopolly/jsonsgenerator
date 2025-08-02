package main

import (
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// get list
func (a *SQL) Fields() []byte {

	var list []string
	list = append(list, "\n//parse sql query")
	list = append(list, fmt.Sprintf(`func (a *%s) Fields() (res []%s) {`, settings.SQL.Class, settings.IndexTypeName))

	// if fields == nil
	ifNilFields := `return []indexTypeName { allFields }`
	ifNilFields = tools.Replace(ifNilFields, "allFields", strings.Join(settings.Go.Indexes, ","))
	ifNilFields = tools.Replace(ifNilFields, "indexTypeName", settings.IndexTypeName)
	list = append(list, ifNilFields)
	list = append(list, "}")
	return []byte(strings.Join(list, "\n"))
}
