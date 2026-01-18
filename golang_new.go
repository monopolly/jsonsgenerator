package main

import (
	"fmt"
	"strings"
)

// get struct value with function
func (a *Golang) New() []byte {

	var need bool
	for _, x := range fields {
		if strings.Contains(x.Type, "map") {
			need = true
		}
	}

	if !need {
		return nil
	}

	var list []string
	list = append(list, "\n")
	list = append(list, "//init struct")
	list = append(list, fmt.Sprintf("func New%[1]s() (a *%[1]s){ ", settings.Go.StructName))
	list = append(list, fmt.Sprintf("a = new(%s)", settings.Go.StructName))

	for _, x := range fields {
		if !strings.Contains(x.Type, "map") {
			continue
		}
		list = append(list, fmt.Sprintf(`a.%s = make(%s)`, x.Go.Name, x.Type))
	}
	list = append(list, "return")
	list = append(list, "}")
	list = append(list, "\n")

	return []byte(strings.Join(list, "\n"))
}
