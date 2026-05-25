package main

import jsoniter "github.com/json-iterator/go"

// generates New() for jsons struct
func jsonDemo() []byte {
	list := map[string]any{}
	for _, x := range fields {
		k, v := x.JsonDemo()
		list[k] = v
	}
	b, _ := jsoniter.Marshal(list)
	return b
}
