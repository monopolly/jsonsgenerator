package main

import (
	"fmt"
	"strings"
)

// create gotiny marshal
func (a *Golang) Gotiny() []byte {
	if !settings.Go.Gotiny {
		return nil
	}

	var list []string
	list = append(list, "\n")
	list = append(list, "//gotiny marshal")
	list = append(list, fmt.Sprintf("func (a *%s) MarshalGotiny() []byte { ", settings.Go.StructName))

	var names []string
	for _, x := range fields {
		names = append(names, fmt.Sprintf("&a.%s", x.Go.Name))
	}
	list = append(list, fmt.Sprintf("return gotiny.Marshal(%s)", strings.Join(names, ", ")))

	list = append(list, "}")
	list = append(list, "\n")

	return []byte(strings.Join(list, "\n"))
}

// create gotiny unmarshal
func (a *Golang) GotinyUnmarshall() []byte {
	if !settings.Go.Gotiny {
		return nil
	}

	var list []string
	list = append(list, "\n")
	list = append(list, "//parse gotiny")
	list = append(list, fmt.Sprintf("func Unmarshal%[1]sGotiny(v []byte) (a %[1]s) { ", settings.Go.StructName))
	// list = append(list, fmt.Sprintf( "var a %s",settings.Go.StructName))

	var names []string
	for _, x := range fields {
		names = append(names, fmt.Sprintf("&a.%s", x.Go.Name))
	}
	list = append(list, fmt.Sprintf("gotiny.Unmarshal(v, %s)", strings.Join(names, ", ")))
	list = append(list, "return")

	list = append(list, "}")
	list = append(list, "\n")

	return []byte(strings.Join(list, "\n"))
}
