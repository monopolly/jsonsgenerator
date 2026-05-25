package main

import "fmt"

func goConvertCode(fieldName, sourceName, fieldType string) string {
	switch fieldType {
	case "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64",
		"[]int8", "[]int16", "[]int32", "[]int64",
		"[]float32":
		return fmt.Sprintf(`b, _ := jsoniter.Marshal(%s)
			_ = jsoniter.Unmarshal(b, &a.%s)`, sourceName, fieldName)
	default:
		return fmt.Sprintf(`cast.Convert(&a.%s, %s)`, fieldName, sourceName)
	}
}

func goConvertSupported(fieldType string) bool {
	switch fieldType {
	case "[]uint", "[]uint8", "[]uint16", "[]uint32", "[]uint64",
		"[]int8", "[]int16", "[]int32", "[]int64",
		"[]float32":
		return true
	default:
		return false
	}
}
