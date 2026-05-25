package main

import (
	"fmt"
	"strings"
)

var jsn JSON

type JSON int

// generates New() for jsons struct
func (a *JSON) generateJsonsInit() []byte {
	res := `
	//NewFuncName create struct
	func NewFuncName() structName {
		return []byte("type")
	}
`
	if settings.Options.Noinit {
		res = fmt.Sprintf(`
	/*
		%s
	*/
	`, res)
	}

	res = strings.ReplaceAll(res, "structName", settings.JS.Name)
	res = strings.ReplaceAll(res, "FuncName", settings.JS.Name)
	res = strings.ReplaceAll(res, "type", "{}")
	return []byte(res)
}

// struct, get, set
func (a *JSON) generateJsonsGeneral() []byte {
	var functions []string

	//type jsons
	functions = append(functions, fmt.Sprintf(`
	//%[1]s is a struct
	type %[1]s []byte`, settings.JS.Name))

	//set value
	functions = append(functions, fmt.Sprintf(`
		//Set value
		func (a *%s) Set(k string, v any) *%s {
			(*a) = jsons.Set((*a), k,v)
			return a
		}
	`, settings.JS.Name, settings.JS.Name))

	//get value
	functions = append(functions, fmt.Sprintf(`
		//Get value
		func (a *%s) Get(k string) jsons.Result {
			return jsons.Get((*a), k)
		}
	`, settings.JS.Name))

	//delete fields value
	functions = append(functions, fmt.Sprintf(`
		//Get value
		func (a *%s) DeleteFields(fields ...string) {
			(*a) = jsons.Delete((*a), fields...)
		}
	`, settings.JS.Name))

	return []byte(strings.Join(functions, "\n"))
}

// jsons functions
func (a *JSON) generateJsonsFunctions() []byte {
	var functions []string

	for _, x := range fields {

		var jsonsValue string

		result := JSONdefaultfunc

		fieldType := x.Type
		if x.Go.Type != "" {
			fieldType = x.Go.Type
		}
		varType := fieldType
		var funcTypeMapKey, funcTypeMapValue string
		switch fieldType {
		case "float64":
			jsonsValue = "Float64"
		case "float32":
			jsonsValue = "Float32"
		case "int8", "int16", "int32":
			jsonsValue = "Int"
			result = JSONconvertfunc
		case "int":
			jsonsValue = "Int"
			// increment and decrement functions
			if x.Json.Inc {
				result = result + JSONinc
			}
		case "int64":
			jsonsValue = "Int64"
			result = `
				//funcName set or get value
				func (a *structName) funcName(v ...funcType) (res funcType) {
					if v == nil {
						return jsons.jsonsFunc((*a), fieldKey)
					}
					a.Set(fieldKey, v[0])
					return
				}
			 `

			// increment and decrement functions
			if x.Json.Inc {
				result = result + JSONinc
			}

			// add time converter when needed
			switch strings.ToLower(x.Name) {
			case
				"created", "create",
				"time", "times", "timer", "date",
				"publish", "published",
				"up", "updated", "update",
				"deleted":
				result = result + `
					//funcNameTime get value as time
					func (a *structName) funcNameTime() (res time.Time) {
						return time.Unix(a.funcName(),0)
					}

				`
				x.Json.Time = true
			default:
				if x.Json.Time {
					result = result + `
					//funcNameTime get value as time
					func (a *structName) funcNameTime() (res time.Time) {
						return time.Unix(a.funcName(),0)
					}

				`
				}
			}

		case "uint64":
			jsonsValue = "Uint64"
			// increment and decrement functions
			if x.Json.Inc {
				result = result + JSONinc
			}
		case "uint32":
			jsonsValue = "Uint32"
			// increment and decrement functions
			if x.Json.Inc {
				result = result + JSONinc
			}
		case "uint16":
			jsonsValue = "Uint"
			result = JSONconvertfunc
		case "uint":
			jsonsValue = "Uint"
			// increment and decrement functions
			if x.Json.Inc {
				result = result + JSONinc
			}
		case "string":
			jsonsValue = "String"
		case "[]byte":
			jsonsValue = "Bytes"
		case "byte":
			jsonsValue = "Byte"
		case "uint8":
			jsonsValue = "Uint8"
		case "bool":
			jsonsValue = "Bool"
			if x.Json.Bool {
				result = result + JSONbool
			}
		case "[]string":
			varType = "string"
			jsonsValue = "ArrayString"
			result = JSONarray
		case "[]int":
			jsonsValue = "ArrayInt"
			varType = "int"
			result = JSONarray
		case "[]int8":
			jsonsValue = "ArrayInt8"
			varType = "int8"
			result = JSONarray
		case "[]int16":
			jsonsValue = "ArrayInt16"
			varType = "int16"
			result = JSONarray
		case "[]int32":
			jsonsValue = "ArrayInt32"
			varType = "int32"
			result = JSONarray
		case "[]int64":
			jsonsValue = "ArrayInt64"
			varType = "int64"
			result = JSONarray

		case "[]uint":
			jsonsValue = "ArrayUint"
			varType = "uint"
			result = JSONarray
		case "[]uint8":
			jsonsValue = "ArrayUint8"
			varType = "uint8"
			result = JSONarray
		case "[]uint16":
			jsonsValue = "ArrayUint16"
			varType = "uint16"
			result = JSONarray
		case "[]uint32":
			jsonsValue = "ArrayUint32"
			varType = "uint32"
			result = JSONarray
		case "[]uint64":
			jsonsValue = "ArrayUint64"
			varType = "uint64"
			result = JSONarray

		case "[]float64":
			varType = "float64"
			result = JSONarrayGeneric
		case "map[string]string":
			jsonsValue = "MapString"
			funcTypeMapKey = "string"
			funcTypeMapValue = "string"
			result = JSONmap
		case "map[string][]byte":
			result = JSONgeneric
		case "map[string]bool":
			jsonsValue = "MapBool"
			funcTypeMapKey = "string"
			funcTypeMapValue = "bool"
			result = JSONmap
		case "map[string]int":
			jsonsValue = "MapInt"
			funcTypeMapKey = "string"
			funcTypeMapValue = "int"
			result = JSONmap
		case "map[string]float64":
			jsonsValue = "MapFloats"
			funcTypeMapKey = "string"
			funcTypeMapValue = "float64"
			result = JSONmap
		case "map[string]any":
			jsonsValue = "MapAny"
			funcTypeMapKey = "string"
			funcTypeMapValue = "any"
			result = JSONmap
		case "map[int]int":
			jsonsValue = "MapIntInt"
			funcTypeMapKey = "int"
			funcTypeMapValue = "int"
			result = JSONmap
		case "map[int]bool":
			result = JSONgeneric
		case "map[int]string":
			jsonsValue = "MapIntString"
			funcTypeMapKey = "int"
			funcTypeMapValue = "string"
			result = JSONmap
		case "date.Int":
			jsonsValue = "Int"
			varType = "int"
			result = JSONdateint
		case "any":
			result = JSONany
		case "time.Duration":
			jsonsValue = "TimeDuration"
			result = result + JSONnamefunc
		}

		result = strings.ReplaceAll(result, "funcTypeMapKey", funcTypeMapKey)
		result = strings.ReplaceAll(result, "funcTypeMapValue", funcTypeMapValue)
		result = strings.ReplaceAll(result, "structName", settings.JS.Name)
		result = strings.ReplaceAll(result, "funcName", x.Go.Name)
		result = strings.ReplaceAll(result, "funcType", fieldType)
		result = strings.ReplaceAll(result, "jsonsFunc", jsonsValue)
		result = strings.ReplaceAll(result, "fieldKey", x.FieldName)
		result = strings.ReplaceAll(result, "varType", varType)

		if x.Go.Nofunc {
			result = fmt.Sprintf(`/*
				%s
			*/`, result)
		}
		functions = append(functions, result)
	}
	return []byte(strings.Join(functions, "\n"))
}
