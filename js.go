package main

import (
	"fmt"
	"strings"
)

// js model
func generateJSModel() []byte {

	structname := settings.Origin

	/*
		type Profile = {
			id: Number,
			email: String
			created: Number
			updated: number
			verified: Boolean
			business: boolean
		}

	*/

	var list []string
	//var switcher []string

	//switcher = append(switcher, "function update(k: string, v: any){\n")
	//switcher = append(switcher, "switch (k) {")

	list = append(list, "export default{\n")
	list = append(list, fmt.Sprintf("\t%s: {", structname))

	for _, x := range fields {
		var Type string
		switch x.Type {
		case "int", "int64", "uint", "uint64", "time.Duration", "int16", "int32", "uint32", "int8", "uint8":
			Type = "undefined"
		case "bool":
			Type = "false"
		case "float64", "float32":
			Type = "undefined"
		case "string":
			Type = "undefined"
		default:
			switch {
			case strings.HasPrefix(x.Type, "[]"):
				Type = "[]"
			case strings.HasPrefix(x.Type, "map["):
				Type = "{}"
			default:
				Type = "undefined"
			}
		}

		// var line string
		line := fmt.Sprintf("\t\t%s: %s, //%s %s", x.Name, Type, x.Type, x.Comment)
		list = append(list, line)

	}

	res := strings.Join(list, "\n") + "\n\t},\n\n}\n"

	/* switcher = append(switcher, "default: \n}")
	switcher = append(switcher, "}")
	res += "\n" + strings.Join(switcher, "\n") */
	return []byte(res)
}
