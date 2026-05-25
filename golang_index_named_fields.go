package main

import (
	"fmt"
	"strings"
)

// field names for json
func (a *Golang) IndexNamed() []byte {

	prefix := "Json"

	// list
	var list []string
	list = append(list, "")
	list = append(list, "// field names for json")
	list = append(list, "const (")
	for _, x := range fields {
		//JsonFieldNewsID
		x.Json.Named = fmt.Sprintf("%s%s%s", prefix, settings.JS.Name, x.Go.Name)
		var line string
		// custom json field, if present
		switch x.Json.Name != "" {
		case true:
			//fieldName = "options.name" //string comment...
			line = fmt.Sprintf(`%s = "%s" // %s %s`, x.Json.Named, x.Json.Name, x.Type, x.Comment)
		case false:
			//fieldName = "name" //string comment...
			line = fmt.Sprintf(`%s = "%s" // %s %s`, x.Json.Named, x.Name, x.Type, x.Comment)
		}
		list = append(list, line)
	}
	list = append(list, ")")
	return []byte(strings.Join(list, "\n"))
}
