package main

import (
	"fmt"
	"strings"
)

// create fast json marshal
func (a *Golang) FastJson() []byte {

	var list []string
	list = append(list, "\n")
	list = append(list, "//fast json marshal")
	list = append(list, fmt.Sprintf("func (a *%s) Pack() []byte {", settings.Go.StructName))
	list = append(list, "var jsoner = jsoniter.ConfigCompatibleWithStandardLibrary")
	list = append(list, "b, _ := jsoner.Marshal(a)")
	list = append(list, "return b")
	list = append(list, "}")
	list = append(list, "\n")

	return []byte(strings.Join(list, "\n"))
}

// create fast json unmarshal
func (a *Golang) FastJsonUnmarshall() []byte {
	var list []string
	list = append(list, "\n")
	list = append(list, "//fast json unmarshal")
	list = append(list, fmt.Sprintf("func Parse%[1]s(v []byte) (a %[1]s, err error) { ", settings.Go.StructName))
	list = append(list, "var jsoner = jsoniter.ConfigCompatibleWithStandardLibrary")
	list = append(list, `err = jsoner.Unmarshal(v, &a)`)
	list = append(list, "return")
	list = append(list, "}")
	list = append(list, "\n")
	return []byte(strings.Join(list, "\n"))
}

// create generated json marshal
func (a *Golang) EasyJsonMarshal() []byte {
	var list []string
	list = append(list, "\n")
	list = append(list, "//easyjson marshal")
	list = append(list, fmt.Sprintf("func (a *%s) Marshal() []byte {", settings.Go.StructName))
	list = append(list, "b, _ := easyjson.Marshal(a)")
	list = append(list, "return b")
	list = append(list, "}")
	list = append(list, "\n")
	return []byte(strings.Join(list, "\n"))
}

// create generated json unmarshal
func (a *Golang) EasyJsonUnmarshal() []byte {
	var list []string
	list = append(list, "\n")
	list = append(list, "//easyjson unmarshal")
	list = append(list, fmt.Sprintf("func %sUnmarshal(src []byte) (a *%s) {", settings.Go.StructName, settings.Go.StructName))
	list = append(list, "a = new("+settings.Go.StructName+")")
	list = append(list, "err := easyjson.Unmarshal(src, a)")
	list = append(list, "if err != nil {")
	list = append(list, "\treturn nil")
	list = append(list, "}")
	list = append(list, "return")
	list = append(list, "}")
	list = append(list, "\n")
	return []byte(strings.Join(list, "\n"))
}
