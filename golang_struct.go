package main

import (
	"fmt"
	"strings"
)

// field names for json
func (a *Golang) Struct() []byte {

	// list
	var list, structComments []string

	// struct
	list = append(list, "\n")
	list = append(list, fmt.Sprintf("//%s", strings.Join(append(structComments, fmt.Sprintf("%d bytes (go padding)", settings.StructSize)), " ")))
	list = append(list, "//easyjson:json")
	list = append(list, fmt.Sprintf(`type %s struct {`, settings.Go.StructName))

	for _, x := range fields {

		//Name int `json:"name,omitempty"` //comment
		var structTags string

		var JsonFieldName string
		switch x.Json.Name == "" {
		case true:
			JsonFieldName = x.Name
		case false:
			JsonFieldName = x.Json.Name
		}

		var taglist []string
		omit := ",omitempty"
		if settings.Go.NoOmit {
			omit = ""
		}

		taglist = append(taglist, fmt.Sprintf("json:\"%s%s\"", JsonFieldName, omit))
		if settings.Go.MessagePack {
			taglist = append(taglist, fmt.Sprintf("msg:\"%s%s\"", JsonFieldName, omit))
		}
		structTags = fmt.Sprintf("`%s`", strings.Join(taglist, " "))

		var line string

		// go type
		var Type string
		switch x.Go.Type != "" {
		case true:
			Type = x.Go.Type
		default:
			Type = x.Type
		}

		// go type
		// var Name string
		// switch x.Go.Name != "" {
		// case true:
		// 	Name = x.Go.Name
		// default:
		// 	Name = x.Title
		// }

		line = fmt.Sprintf(`%s %s %s // %s`, x.Go.Name, Type, structTags, x.Comment)
		list = append(list, line)
	}
	list = append(list, "}\n")
	return []byte(strings.Join(list, "\n"))
}
