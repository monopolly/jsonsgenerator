package main

import (
	"fmt"
	"strings"
)

// create gotiny marshal
func (a *Golang) Msgp() []byte {
	if !settings.Go.MessagePack {
		return nil
	}

	var list []string
	list = append(list, "\n")
	list = append(list, "//msgp marshal")
	list = append(list, fmt.Sprintf("func (a *%s) MessagePack() []byte { ", settings.Go.StructName))
	list = append(list, "b, _ := msgpack.Marshal(a)")
	list = append(list, "return b")
	list = append(list, "}")
	list = append(list, "\n")

	return []byte(strings.Join(list, "\n"))
}

// create gotiny unmarshal
func (a *Golang) MsgpUnmarshall() []byte {
	if !settings.Go.MessagePack {
		return nil
	}

	var list []string
	list = append(list, "\n")
	list = append(list, "//msgp unmarshal")
	list = append(list, fmt.Sprintf("func Parse%[1]sMessagePack(v []byte) (a %[1]s, err error) { ", settings.Go.StructName))
	list = append(list, `err = msgpack.Unmarshal(v, &a)`)
	list = append(list, "return")
	list = append(list, "}")
	list = append(list, "\n")

	return []byte(strings.Join(list, "\n"))
}
