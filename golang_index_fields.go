package main

import (
	"fmt"
	"strings"
)

// FieldNewsID = "id"
func (a *Golang) IndexStrings() []byte {

	prefix := "Field"

	// list
	var list []string
	list = append(list, "\n")
	list = append(list, "//string index")
	list = append(list, "const (")
	for _, x := range fields {
		//fieldNewsID
		x.FieldName = fmt.Sprintf("%s%s%s", settings.Go.StructName, prefix, x.Go.Name)
		var line string

		switch x.Json.Name != "" {
		case true:
			//fieldName = "options.name" //string comment...
			line = fmt.Sprintf(`%s = "%s" // %s %s`, x.FieldName, x.Json.Name, x.Type, x.Comment)
		case false:
			//fieldName = "name" //string comment...
			line = fmt.Sprintf(`%s = "%s" // %s %s`, x.FieldName, x.Name, x.Type, x.Comment)
		}

		list = append(list, line)
	}
	list = append(list, ")\n\n")
	return []byte(strings.Join(list, "\n"))
}
