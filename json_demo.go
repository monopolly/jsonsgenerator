package main

import jsoniter "github.com/json-iterator/go"

// генерит New() для jsons структуры
func jsonDemo() []byte {
	list := map[string]any{}
	for _, x := range fields {
		k, v := x.JsonDemo()
		list[k] = v
	}
	b, _ := jsoniter.Marshal(list)
	return b
}
