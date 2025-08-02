package main

import (
	"fmt"
	"strings"
)

// FieldNewsID = "id"
func (a *Golang) IndexStrings() []byte {

	prefix := "Field"

	//список
	var list []string
	list = append(list, "\n")
	list = append(list, "//string index")
	list = append(list, "const (")
	for _, x := range fields {
		//fieldNewsID

		switch settings.Go.NoPrefix {
		case true:
			x.FieldName = fmt.Sprintf("%s%s", prefix, x.Go.Name)
		case false:
			x.FieldName = fmt.Sprintf("%s%s%s", prefix, settings.Go.StructName, x.Go.Name)
		}

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
