package main

import "fmt"

const (
	uuidType   = "uuid.UUID"
	uuidImport = "github.com/google/uuid"
)

func goConvertCode(fieldName, sourceName, fieldType string) string {
	return goConvertTo("a."+fieldName, sourceName, fieldType)
}

// convert source to target var: cast.Convert(&a.Name, x)
func goConvertTo(target, sourceName, fieldType string) string {
	switch fieldType {
	case uuidType:
		// pgx: [16]byte, clickhouse: uuid.UUID, json: string
		return fmt.Sprintf(`switch v := %[2]s.(type) {
			case uuid.UUID:
				%[1]s = v
			case [16]byte:
				%[1]s = v
			case string:
				%[1]s, _ = uuid.Parse(v)
			case []byte:
				if len(v) == 16 {
					%[1]s, _ = uuid.FromBytes(v)
				} else {
					%[1]s, _ = uuid.ParseBytes(v)
				}
			}`, target, sourceName)
	case "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64",
		"[]int8", "[]int16", "[]int32", "[]int64",
		"[]float32", "[]" + uuidType, "[]time.Time", "[]time.Duration":
		return fmt.Sprintf(`b, _ := jsoniter.Marshal(%s)
			_ = jsoniter.Unmarshal(b, &%s)`, sourceName, target)
	default:
		return fmt.Sprintf(`cast.Convert(&%s, %s)`, target, sourceName)
	}
}

// field converter, custom type via base type:
// var v int; cast.Convert(&v, x); a.Role = role.Role(v)
func goConvertField(x *Field, sourceName string) string {
	if x.Go.Base == "" {
		return goConvertCode(x.Go.Name, sourceName, x.Type)
	}
	return fmt.Sprintf("var v %s\n%s\na.%s = %s(v)", x.Go.Base, goConvertTo("v", sourceName, x.Go.Base), x.Go.Name, x.GoType())
}

func goConvertSupported(fieldType string) bool {
	switch fieldType {
	case uuidType, "[]" + uuidType, "[]time.Time", "[]time.Duration",
		"[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64",
		"[]int8", "[]int16", "[]int32", "[]int64",
		"[]float32":
		return true
	default:
		return false
	}
}
