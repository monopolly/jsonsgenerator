package main

import (
	"fmt"
	"jj/tools"
	"strings"
)

func (a *Golang) Tags() []byte {
	tags := map[string][]*Field{}

	for _, field := range fields {
		for tag := range field.Go.Tags {
			tags[tag] = append(tags[tag], field)
		}
	}

	if len(tags) == 0 {
		return nil
	}

	var list []string
	list = append(list, "\n")
	for tag, filds := range tags {
		list = append(list, fmt.Sprintf("func %s%sList() []%s {", settings.Go.StructName, tools.Title(tag), settings.IndexTypeName))
		var taglist []string
		for _, x := range filds {
			taglist = append(taglist, x.Go.Index)
		}
		list = append(list, fmt.Sprintf(`return []%s{%s}`, settings.IndexTypeName, strings.Join(taglist, ", ")))
		list = append(list, "}\n")
	}

	return []byte(strings.Join(list, "\n"))
}
