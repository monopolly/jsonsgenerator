package main

import (
	"fmt"
	"strings"
)

// create generated json marshal
func (a *Golang) EasyJsonMarshal() []byte {
	var list []string
	list = append(list, "\n")
	list = append(list, "//easyjson marshal")
	list = append(list, fmt.Sprintf("func (a *%s) MarshalJSON() []byte {", settings.Go.StructName))
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
	list = append(list, fmt.Sprintf("func %sUnmarshalJSON(src []byte) (a *%s) {", settings.Go.StructName, settings.Go.StructName))
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
