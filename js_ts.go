package main

import (
	"fmt"
	"strings"
)

// insert sql query
func generateTypeScript() []byte {

	var structname string
	if settings.TS.Name != "" {
		structname = settings.TS.Name
	} else {
		structname = settings.Origin
	}

	/*
				type Profile = {
					id: number
					email: string //ext
					created: number
					updated: number
					name: string
					line: string
					about: string
					link: string
					image: string
					geoid: number
					country: string
					city: string
					address: string
					phone: string
					lat: number
					long: number
					cat: string
					sub: string
					verified: boolean
					business: boolean
				}

				switch (k) {
		          case 'name'     : this.profile.name       = v; return;
		          case 'line'     : this.profile.line       = v; return;
		          case 'link'     : this.profile.link       = v; return;
		          case 'image'    : this.profile.image      = v; return;
		          case 'about'    : this.profile.about      = v; return;
		          case 'business' : this.profile.business  = v; return;
		          default:
		        }
	*/

	var list []string
	//var switcher []string

	//switcher = append(switcher, "function update(k: string, v: any){\n")
	//switcher = append(switcher, "switch (k) {")

	list = append(list, fmt.Sprintf("type %s = {", structname))

	var pad int
	for _, x := range fields {
		lens := len(x.Name)
		if pad < lens {
			pad = lens
		}
	}

	// pads := strings.Repeat(" ", pad)

	for _, x := range fields {

		var tsType string
		switch x.PlainType() {
		case "int", "int64", "uint", "uint64", "time.Duration", "int16", "int32", "uint32", "int8", "uint8", "float64", "float32":
			tsType = "number"
		case "bool":
			tsType = "boolean"
		case "string", uuidType, "time.Time":
			tsType = "string"
		default:
			tsType = "any"
		}

		if strings.HasPrefix(x.PlainType(), "[]") {
			if x.PlainType() == "[]time.Time" {
				tsType = "string[]"
			} else if strings.Contains(x.PlainType(), "int") || strings.Contains(x.PlainType(), "float") || strings.Contains(x.PlainType(), "time") {
				tsType = "number[]"
			} else if strings.Contains(x.PlainType(), "str") || strings.Contains(x.PlainType(), uuidType) {
				tsType = "string[]"
			} else if strings.Contains(x.PlainType(), "byte") {
				tsType = "string"
			} else {
				tsType = "Array<any>"
			}
		}

		if strings.HasPrefix(x.PlainType(), "map[") {
			if strings.Contains(x.PlainType(), "inter") || strings.Contains(x.PlainType(), "any") {
				tsType = "Record<string, any>"
			} else if strings.Contains(x.PlainType(), "int") || strings.Contains(x.PlainType(), "float") || strings.Contains(x.PlainType(), "time") {
				tsType = "Record<string, number>"
			} else if strings.Contains(x.PlainType(), "bool") {
				tsType = "Record<string, boolean>"
			} else if strings.Contains(x.PlainType(), "string") {
				tsType = "Record<string, string>"
			} else {
				tsType = "Array<any>"
			}
		}

		list = append(list, fmt.Sprintf("\t%s:%s%s", x.Name, strings.Repeat(" ", (pad+5)-len(x.Name)), tsType))

		//switcher = append(switcher, formatline(fmt.Sprintf(`case '%s': `, x.Name), fmt.Sprintf(`this.%s.%s = v; return;`, structname, x.Name)))

	}

	res := strings.Join(list, "\n") + "\n}"
	/* switcher = append(switcher, "default: \n}")
	switcher = append(switcher, "}")
	res += "\n" + strings.Join(switcher, "\n") */
	return []byte(res)
}
