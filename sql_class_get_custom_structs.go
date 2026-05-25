package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// get list
func (a *SQL) GetCustomStructs() []byte {

	all := map[string]bool{}
	for _, x := range fields {
		for k := range x.Go.CustomStructs {
			all[k] = true
		}
	}

	var b bytes.Buffer
	for k := range all {
		b.Write(a.getCustomStructs(k))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	return b.Bytes()
}

// get list
func (a *SQL) getCustomStructs(name string) []byte {

	var list []string
	list = append(list, "\n//get custom struct sql query")
	list = append(list, fmt.Sprintf(`func (a *%[1]s) Get%[2]s(id any) (res *%[3]s%[2]s, err error) {`, settings.SQL.Class, tools.Title(name), settings.Go.StructName))
	list = append(list, fmt.Sprintf("p, err := a.Get(id, %s%sList()...)", settings.Go.StructName, tools.Title(name)))
	list = append(list, `if err != nil{return}`)
	list = append(list, fmt.Sprintf("res = p.%s%s()", settings.Go.StructName, tools.Title(name)))
	list = append(list, "return}")
	return []byte(strings.Join(list, "\n"))
}
