package main

import (
	"math/rand"
)

// json field
func (a *Field) JsonDemo() (k string, v any) {

	switch a.Type {
	case "string":
		v = "string"
	case uuidType:
		v = "8e1f2a4c-6b3d-4f7e-9a0c-2d5b8e1f3a6c"
	case "time.Time":
		v = "2026-01-15T10:30:00Z"
	case "[]time.Time":
		v = []string{"2026-01-15T10:30:00Z", "2026-02-20T18:00:00Z"}
	case "[]time.Duration":
		v = []int{rand.Intn(999999), rand.Intn(999999)}
	case "[]" + uuidType:
		v = []string{"8e1f2a4c-6b3d-4f7e-9a0c-2d5b8e1f3a6c", "3c7d9b2e-1a4f-4e8b-b6d0-5f2a9c3e7b1d"}
	case "bool":
		v = true
	case "int8", "uint8", "int", "int64", "uint", "uint64", "time.Duration", "int16", "int32", "uint32":
		v = rand.Intn(999999)
	case "float64", "float32":
		v = rand.Float64()
	case "map[string]any":
		v = map[string]any{
			"name":  "James",
			"admin": true,
			"year":  2024,
		}
	case "map[string]bool":
		v = map[string]bool{
			"admin": true,
			"top":   true,
		}
	case "map[string]string":
		v = map[string]string{
			"name":    "James",
			"family":  "Prestone",
			"twitter": "@nice",
		}
	case "map[int]string":
		v = map[int]string{
			1:     "James",
			222:   "Prestone",
			55555: "@nice",
		}
	case "map[int]int":
		v = map[int]int{
			1:     331,
			222:   411,
			55555: 42111,
		}
	case "map[string]int":
		v = map[string]int{
			"likes": rand.Intn(999),
			"views": rand.Intn(9999999),
		}
	case "[]byte":
		v = []byte(`{"rawjson": true, "or":"any bytes data"}`)
	case "[]string":
		v = []string{"one", "two", "three"}
	case "[]int":
		v = []int{1, 2, 3, 4, 5}
	case "[]int8":
		v = []int8{1, 2, 3, 4, 5}
	case "[]int16":
		v = []int16{1, 2, 3, 4, 5}
	case "[]int32":
		v = []int32{1, 2, 3, 4, 5}
	case "[]int64":
		v = []int64{1, 2, 3, 4, 5}
	case "[]uint":
		v = []uint{1, 2, 3, 4, 5}
	case "[]uint8":
		v = []uint8{1, 2, 3, 4, 5}
	case "[]uint16":
		v = []uint16{1, 2, 3, 4, 5}
	case "[]uint32":
		v = []uint32{1, 2, 3, 4, 5}
	case "[]uint64":
		v = []uint64{1, 2, 3, 4, 5}
	case "[]float64":
		v = []float64{0.421, 0.2456, 0.24114}
	case "[]float32":
		v = []float32{0.421, 0.2456, 0.24114}
	case "byte":
		v = 'b'
	case "any":
		v = "any type"
	default:
		v = a.Type
	}

	switch a.Json.Name != "" {
	case true:
		k = a.Json.Name
	default:
		k = a.Name
	}

	return
}
