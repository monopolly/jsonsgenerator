package main

import (
	"fmt"
	"strings"
	"time"
)

var helps = NewHelp()

func NewHelp() (a *Help) {
	a = new(Help)
	a.init()
	return
}

type Help struct {
	Struct     map[string]string //value
	Golang     map[string]string
	Field      map[string]string
	JS         map[string]string
	Swift      map[string]string
	SQL        map[string]string
	Clickhouse map[string]string
	Proto      map[string]string
	Custom     map[string]map[string]string
}

func (a *Help) init() (b *Help) {
	a.Struct = make(map[string]string)
	a.Golang = make(map[string]string)
	a.Field = make(map[string]string)
	a.JS = make(map[string]string)
	a.Swift = make(map[string]string)
	a.SQL = make(map[string]string)
	a.Clickhouse = make(map[string]string)
	a.Proto = make(map[string]string)
	a.Custom = make(map[string]map[string]string)
	return a
}

func (a *Help) AddCustom(category, k, v string) {
	if a.Custom[category] == nil {
		a.Custom[category] = make(map[string]string)
	}
	a.Custom[category][k] = v
}

func (a *Help) Render() (res []byte) {

	var list []string
	list = append(list, "/*")
	list = append(list, "Auto-generate. Do not change!")
	list = append(list, "Help to auto-generated correct models for golang structures.")
	list = append(list, fmt.Sprintf("Sergey Keplin (c) %d", time.Now().Year()))
	list = append(list, "a@senthy.com, lava.mobi@gmail.com, https://t.me/martinprestone")
	list = append(list, "github.com/monopolly/jsons")
	list = append(list, "Ex: //! [] go=News js=NewsJson sql=news noinit up...")
	list = append(list, "")

	// structs
	list = append(list, "STRUCT:")
	for k, v := range a.Struct {
		list = append(list, Status(k, v))
	}

	// golang
	list = append(list, "")
	list = append(list, "Field:")
	for k, v := range a.Field {
		list = append(list, Status(k, v))
	}

	// golang
	list = append(list, "")
	list = append(list, "GO:")
	for k, v := range a.Golang {
		list = append(list, Status(k, v))
	}

	// sql
	list = append(list, "")
	list = append(list, "SQL:")
	for k, v := range a.SQL {
		list = append(list, Status(k, v))
	}

	// clickhouse
	list = append(list, "")
	list = append(list, "Clickhouse:")
	for k, v := range a.Clickhouse {
		list = append(list, Status(k, v))
	}

	// proto
	list = append(list, "")
	list = append(list, "PROTO:")
	for k, v := range a.Proto {
		list = append(list, Status(k, v))
	}

	// swift
	list = append(list, "")
	list = append(list, "SWIFT:")
	for k, v := range a.Swift {
		list = append(list, Status(k, v))
	}

	// js
	list = append(list, "")
	list = append(list, "JSON:")
	for k, v := range a.JS {
		list = append(list, Status(k, v))
	}

	// add \t
	for p := range list {
		list[p] = "\t" + list[p]
	}
	list = append(list, "*/\n")

	// render
	res = []byte(strings.Join(list, "\n"))
	return
}

func Status(name string, v interface{}) string {
	return fmt.Sprintf("%s%s%v", name, strings.Repeat(" ", 30-len(name)), v)
}
