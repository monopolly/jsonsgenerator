package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// field names for json
func (a *Golang) StructCustom() []byte {

	all := map[string]bool{}
	for _, x := range fields {
		for k := range x.Go.CustomStructs {
			all[k] = true
		}
	}

	var b bytes.Buffer
	for k := range all {
		// fmt.Println("add custom", k)
		b.Write(a.structCustom(k))
		b.WriteString("\n")
		b.Write(a.structCustomParseFunction(k))
		b.WriteString("\n")
		b.Write(a.structCustomParseFunctionReverse(k))
		b.WriteString("\n")
		b.Write(a.structCustomParseFunctionPack(k))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	return b.Bytes()

}

// field names for json
func (a *Golang) structCustom(name string) []byte {

	// list
	var list, structComments []string

	// struct
	list = append(list, "\n")
	list = append(list, fmt.Sprintf("//%s", strings.Join(structComments, " ")))
	list = append(list, fmt.Sprintf(`type %s%s struct {`, settings.Go.StructName, tools.Title(name)))

	for _, x := range fields {

		if !x.Go.CustomStructs[name] {
			continue
		}

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

// field names for json
func (a *Golang) structCustomParseFunction(name string) []byte {

	// list
	var list, structComments []string

	// struct
	list = append(list, "\n")
	list = append(list, fmt.Sprintf("//%s", strings.Join(structComments, " ")))
	list = append(list, fmt.Sprintf(`func (a *%[1]s) %[1]s%[2]s() (res *%[1]s%[2]s) {`, settings.Go.StructName, tools.Title(name)))
	list = append(list, fmt.Sprintf(`return &%[1]s%[2]s{`, settings.Go.StructName, tools.Title(name)))

	for _, x := range fields {

		if !x.Go.CustomStructs[name] {
			continue
		}

		var line string

		line = fmt.Sprintf(`%[1]s: a.%[1]s,`, x.Go.Name)
		list = append(list, line)
	}
	list = append(list, "}")
	list = append(list, "}\n")
	return []byte(strings.Join(list, "\n"))
}

// field names for json
func (a *Golang) structCustomParseFunctionReverse(name string) []byte {

	// list
	var list, structComments []string

	// struct
	list = append(list, "\n")
	list = append(list, fmt.Sprintf("//%s", strings.Join(structComments, " ")))
	list = append(list, fmt.Sprintf(`func (a *%[1]s%[2]s) %[1]s() (res *%[1]s) {`, settings.Go.StructName, tools.Title(name)))
	list = append(list, fmt.Sprintf(`return &%[1]s{`, settings.Go.StructName, tools.Title(name)))

	for _, x := range fields {

		if !x.Go.CustomStructs[name] {
			continue
		}

		var line string

		line = fmt.Sprintf(`%[1]s: a.%[1]s,`, x.Go.Name)
		list = append(list, line)
	}
	list = append(list, "}")
	list = append(list, "}\n")
	return []byte(strings.Join(list, "\n"))
}

// field names for json
func (a *Golang) structCustomParseFunctionPack(name string) []byte {
	return fmt.Appendf(nil, `
		func (a *%s%s) Pack() (res []byte){
			res,_ = jsoniter.Marshal(a)
			return
		}
	`, settings.Go.StructName, tools.Title(name))
}
