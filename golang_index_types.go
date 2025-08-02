package main

import (
	"fmt"
	"jj/tools"
	"strings"
)

func (a *Golang) IndexType() []byte {

	tags := map[string][]*Field{}

	for _, field := range fields {
		for tag := range field.Go.Tags {
			tags[tag] = append(tags[tag], field)
		}
	}

	var list []string
	list = append(list, "\n")
	for tag := range tags {
		list = append(list, fmt.Sprintf("func (a %s) %s() bool {", settings.IndexTypeName, tools.Title(tag)))
		list = append(list, "switch a {")
		var trues []string
		for _, x := range fields {
			if x.Go.Tags[tag] {
				trues = append(trues, x.Go.Index)
			}
		}
		list = append(list, fmt.Sprintf(`case %s:`, strings.Join(trues, ", ")))
		list = append(list, "return true")

		list = append(list, "default: return false")
		list = append(list, "}")
		list = append(list, "}")
	}

	//string
	list = append(list, "\n")
	list = append(list, "//key index string")
	list = append(list, fmt.Sprintf("func (a %s) String() string {", settings.IndexTypeName))
	list = append(list, "switch a {")
	for _, x := range fields {
		list = append(list, fmt.Sprintf(`case %s:`, x.Go.Index))
		list = append(list, fmt.Sprintf(`return "%s"`, x.Name))
	}
	list = append(list, `default: return ""`)
	list = append(list, "}")
	list = append(list, "}")

	//sql name
	list = append(list, "\n")
	list = append(list, "//key index string")
	list = append(list, fmt.Sprintf("func (a %s) SQLName() string {", settings.IndexTypeName))
	list = append(list, "switch a {")
	for _, x := range fields {
		list = append(list, fmt.Sprintf(`case %s:`, x.Go.Index))
		switch x.SQL.Name != "" {
		case true:
			list = append(list, fmt.Sprintf(`return "%s"`, x.SQL.Name))
		default:
			list = append(list, fmt.Sprintf(`return "%s"`, x.Name))
		}

	}
	list = append(list, `default: return ""`)
	list = append(list, "}")
	list = append(list, "}")

	//type
	list = append(list, "\n")
	list = append(list, "//key index type")
	list = append(list, fmt.Sprintf("func (a %s) Type() string {", settings.IndexTypeName))
	list = append(list, "switch a {")
	for _, x := range fields {
		list = append(list, fmt.Sprintf(`case %s:`, x.Go.Index))
		list = append(list, fmt.Sprintf(`return "%s"`, x.Type))
	}
	list = append(list, `default: return ""`)
	list = append(list, "}")
	list = append(list, "}")

	//title
	list = append(list, "\n")
	list = append(list, "//custom title")
	list = append(list, fmt.Sprintf("func (a %s) Title() string {", settings.IndexTypeName))
	list = append(list, "switch a {")
	for _, x := range fields {
		list = append(list, fmt.Sprintf(`case %s:`, x.Go.Index))
		switch x.Go.Title != "" {
		case true:
			list = append(list, fmt.Sprintf(`return "%s"`, x.Go.Title))
		case false:
			list = append(list, fmt.Sprintf(`return "%s"`, x.Go.Name))
		}
	}
	list = append(list, `default: return ""`)
	list = append(list, "}")
	list = append(list, "}")

	//desc
	list = append(list, "\n")
	list = append(list, "//custom desc")
	list = append(list, fmt.Sprintf("func (a %s) Desc() string {", settings.IndexTypeName))
	list = append(list, "switch a {")
	for _, x := range fields {
		if x.Go.Desc == "" {
			continue
		}
		list = append(list, fmt.Sprintf(`case %s:`, x.Go.Index))
		list = append(list, fmt.Sprintf(`return "%s"`, x.Go.Desc))
	}
	list = append(list, `default: return ""`)
	list = append(list, "}")
	list = append(list, "}")

	return []byte(strings.Join(list, "\n"))
}
