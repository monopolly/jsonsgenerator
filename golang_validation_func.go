package main

import (
	"fmt"
	"strings"
)

// парсер для изначальной структуры в []intefacr{}
func (a *Golang) Valid() []byte {

	lines := []string{
		"\n",
		"//Valid empty values for selected fields by go{must}",
		fmt.Sprintf("func (a *%s) Must() (emptyfields []string){ ", settings.Go.StructName),
	}
	var fieldlist []string
	for _, x := range fields {
		if !x.Go.Must {
			continue
		}
		var condition string
		var v string
		v = x.Type
		if x.Go.Type != "" {
			v = x.Go.Type
		}
		switch v {
		case "int", "int8", "int32", "int64", "uint", "uint64", "uint32", "time.Duration", "date.Int", "float64", "byte":
			condition = "< 1"
		case "string":
			condition = `== ""`
		case "bool":
			condition = " == false"
		case "any", "inteface{}":
			condition = "== nil"
		default:
			switch {
			case strings.Contains(v, "[]"):
				condition = "== nil"
			case strings.Contains(v, "map["):
				condition = "== nil"
			}
		}
		fieldlist = append(fieldlist, fmt.Sprintf(`if a.%s %s {
			emptyfields = append(emptyfields, "%s")
		}`, x.Go.Name, condition, x.Name))
	}

	// no must fields
	if fieldlist == nil {
		return nil
	}

	lines = append(lines, fieldlist...)

	lines = append(lines, `return }`)

	return []byte(strings.Join(lines, "\n"))
}
