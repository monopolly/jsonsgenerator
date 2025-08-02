package main

import "encoding/json"

// генерит New() для jsons структуры
func jsonDemo() []byte {
	list := map[string]any{}
	for _, x := range fields {
		k, v := x.JsonDemo()
		list[k] = v
	}
	b, _ := json.Marshal(list)
	return b
}
