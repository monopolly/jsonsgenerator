package main

import (
	"fmt"
	"strings"

	"github.com/monopolly/cast"
)

// парсер для изначальной структуры в []intefacr{}
func (a *Golang) TupleToStruct() []byte {

	parse := []string{
		"\n",
		"//Parse []any to ID struct",
		fmt.Sprintf(`func Parse%[1]sToStruct(r []any) (a *%[2]s){
			a = new(%[2]s)
			`, settings.Go.StructName, settings.Go.StructName),
		"for pos, x := range r {",
		fmt.Sprintf("switch %s(pos) {", settings.IndexTypeName),
	}

	for _, x := range fields {
		if x.Go.Type != "" {
			continue
		}
		if cast.Supported(x.Type) != nil {
			panic(fmt.Sprintf("cast not support type %s", x.Type))
		}
		var parsefunc string
		parsefunc = fmt.Sprintf(`cast.Convert(&a.%[1]s, x)`, x.Go.Name)
		// switch x.Type {
		// case "int":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.Int(x)`, x.Go.Name)
		// case "int8":
		// 	parsefunc = fmt.Sprintf(`a.%s = int8(cast.Int(x))`, x.Go.Name)
		// case "int16":
		// 	parsefunc = fmt.Sprintf(`a.%s = int16(cast.Int(x))`, x.Go.Name)
		// case "int32":
		// 	parsefunc = fmt.Sprintf(`a.%s = int32(cast.Int(x))`, x.Go.Name)
		// case "int64":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.Int64(x)`, x.Go.Name)
		// case "time.Duration":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.Duration(x)`, x.Go.Name)
		// case "uint":
		// 	parsefunc = fmt.Sprintf(`a.%s = uint(cast.Uint64(x))`, x.Go.Name)
		// case "uint8":
		// 	parsefunc = fmt.Sprintf(`a.%s = uint8(cast.Uint64(x))`, x.Go.Name)
		// case "uint16":
		// 	parsefunc = fmt.Sprintf(`a.%s = uint16(cast.Uint64(x))`, x.Go.Name)
		// case "uint64":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.Uint64(x)`, x.Go.Name)
		// case "uint32":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.Uint32(x)`, x.Go.Name)
		// case "date.Int":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.Date(x)`, x.Go.Name)
		// case "string":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.String(x)`, x.Go.Name)
		// case "[]string":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.SliceString(x)`, x.Go.Name)
		// case "[]int":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.SliceInt(x)`, x.Go.Name)
		// case "[]int64":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.SliceInt64(x)`, x.Go.Name)
		// case "map[string]string":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.StringMapString(x)`, x.Go.Name)
		// case "map[string]any":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.StringMap(x)`, x.Go.Name)
		// case "map[string]bool":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.StringMapBool(x)`, x.Go.Name)
		// case "map[string]int":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.MapStringInt(x)`, x.Go.Name)
		// case "map[string]float64":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.MapStringFloats(x)`, x.Go.Name)
		// case "map[int]string":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.MapIntString(x)`, x.Go.Name)
		// case "map[int]int":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.MapIntInt(x)`, x.Go.Name)
		// case "[]byte":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.Bytes(x)`, x.Go.Name)
		// case "byte":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.Byte(x)`, x.Go.Name)
		// case "bool":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.Bool(x)`, x.Go.Name)
		// case "float64":
		// 	parsefunc = fmt.Sprintf(`a.%s = cast.Float(x)`, x.Go.Name)
		// case "float32":
		// 	parsefunc = fmt.Sprintf(`a.%s = float32(cast.Float(x))`, x.Go.Name)
		// case "any":
		// 	parsefunc = fmt.Sprintf(`a.%s = x`, x.Go.Name)
		// default:
		// 	panic(fmt.Sprintf("tupple to struct: %s is not supported yet", x.Type))
		// }
		parse = append(parse, fmt.Sprintf(`case %s: %s //%s`, x.Go.Index, parsefunc, x.Type))
	}
	parse = append(parse, "}}")
	parse = append(parse, "return")
	parse = append(parse, "}")

	return []byte(strings.Join(parse, "\n"))

}
