package main

import (
	"fmt"
	"strings"
)

// название полей для json
func (a *Golang) IndexNamed() []byte {

	prefix := "Json"

	//список
	var list []string
	list = append(list, "")
	list = append(list, "//название полей для json")
	list = append(list, "const (")
	for _, x := range fields {
		//JsonFieldNewsID
		x.Json.Named = fmt.Sprintf("%s%s%s", prefix, settings.JS.Name, x.Go.Name)
		var line string
		//если есть кастомное json поле
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
