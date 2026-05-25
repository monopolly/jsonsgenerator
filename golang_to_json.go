package main

import (
	"fmt"
	"strings"
)

// parser from the original struct to []interface{}
func (a *Golang) ToJson() []byte {

	js := []string{
		"\n",
		"//Struct to json",
		fmt.Sprintf("func (a *%s) ToJson() (r []byte){ ", settings.Go.StructName),
		`js := jsons.Create().`,
	}

	for pos, x := range fields {

		switch x.Json.Raw {
		case true:
			if x.Type == "[]byte" {
				js = append(js, fmt.Sprintf(`AddRaw(%s,a.%s)`, x.FieldName, x.Go.Name))
			} else {
				js = append(js, fmt.Sprintf(`Add(%s,a.%s)`, x.FieldName, x.Go.Name))
			}
		default:
			if x.Type == "date.Int" {
				js = append(js, fmt.Sprintf(`Add(%s,int(a.%s))`, x.FieldName, x.Go.Name))
			} else {
				js = append(js, fmt.Sprintf(`Add(%s,a.%s)`, x.FieldName, x.Go.Name))
			}

		}

		if pos < len(fields)-1 {
			js[len(js)-1] += "."
		}

	}
	js = append(js, "return js.Bytes()")
	js = append(js, "}\n")

	return []byte(strings.Join(js, "\n"))
}
