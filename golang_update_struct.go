package main

import (
	"fmt"
	"strings"
)

// update struct with function
func (a *Golang) Update() []byte {

	parse := []string{
		"\n",
		"//update struct with function",
		fmt.Sprintf("func (a *%s) Update(k string, x any) {", settings.Go.StructName),
		"switch k {",
		// fmt.Sprintf("switch %s(pos) {", settings.IndexTypeName),
	}

	for _, x := range fields {
		if x.Go.Type != "" {
			continue
		}
		var parsefunc string
		switch x.Type {
		case "int":
			parsefunc = fmt.Sprintf(`a.%s = cast.Int(x)`, x.Go.Name)
		case "int64":
			parsefunc = fmt.Sprintf(`a.%s = cast.Int64(x)`, x.Go.Name)
		case "time.Duration":
			parsefunc = fmt.Sprintf(`a.%s = cast.Duration(x)`, x.Go.Name)
		case "uint64":
			parsefunc = fmt.Sprintf(`a.%s = cast.Uint64(x)`, x.Go.Name)
		case "uint32":
			parsefunc = fmt.Sprintf(`a.%s = cast.Uint32(x)`, x.Go.Name)
		case "date.Int":
			parsefunc = fmt.Sprintf(`a.%s = cast.Date(x)`, x.Go.Name)
		case "string":
			parsefunc = fmt.Sprintf(`a.%s = cast.String(x)`, x.Go.Name)
		case "[]string":
			parsefunc = fmt.Sprintf(`a.%s = cast.SliceString(x)`, x.Go.Name)
		case "[]int":
			parsefunc = fmt.Sprintf(`a.%s = cast.SliceInt(x)`, x.Go.Name)
		case "[]int64":
			parsefunc = fmt.Sprintf(`a.%s = cast.SliceInt64(x)`, x.Go.Name)
		case "map[string]string":
			parsefunc = fmt.Sprintf(`a.%s = cast.StringMapString(x)`, x.Go.Name)
		case "map[string]any":
			parsefunc = fmt.Sprintf(`a.%s = cast.StringMap(x)`, x.Go.Name)
		case "map[string]bool":
			parsefunc = fmt.Sprintf(`a.%s = cast.StringMapBool(x)`, x.Go.Name)
		case "map[string]int":
			parsefunc = fmt.Sprintf(`a.%s = cast.MapStringInt(x)`, x.Go.Name)
		case "map[string]float64":
			parsefunc = fmt.Sprintf(`a.%s = cast.MapStringFloats(x)`, x.Go.Name)
		case "map[int]string":
			parsefunc = fmt.Sprintf(`a.%s = cast.MapIntString(x)`, x.Go.Name)
		case "map[int]int":
			parsefunc = fmt.Sprintf(`a.%s = cast.MapIntInt(x)`, x.Go.Name)
		case "[]byte":
			parsefunc = fmt.Sprintf(`a.%s = cast.Bytes(x)`, x.Go.Name)
		case "byte":
			parsefunc = fmt.Sprintf(`a.%s = cast.Byte(x)`, x.Go.Name)
		case "bool":
			parsefunc = fmt.Sprintf(`a.%s = cast.Bool(x)`, x.Go.Name)
		case "float64":
			parsefunc = fmt.Sprintf(`a.%s = cast.Float(x)`, x.Go.Name)
		case "any", "interface{}":
			parsefunc = fmt.Sprintf(`a.%s = x`, x.Go.Name)
		default:
		}

		parse = append(parse, fmt.Sprintf(`case "%s": %s //%s`, x.Name, parsefunc, x.Type))
	}
	parse = append(parse, "}")
	// parse = append(parse, "return")
	parse = append(parse, "}")

	return []byte(strings.Join(parse, "\n"))

}
