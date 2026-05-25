package main

import (
	"fmt"
	"regexp"
	"strconv"

	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

func (a *Field) Parse() {

	a.Go.Tags = make(map[string]bool)
	a.Go.CustomStructs = make(map[string]bool)
	a.Go.UserLists = make(map[string][]string)

	// because Go does not allow using type directly here
	switch a.Name {
	case "types":
		a.Name = "type"
	}

	a.SQL.Name = a.Name
	a.Clickhouse.Name = a.Name

	a.Go.Name = tools.Title(a.Name)

	// make upper for id>ID
	if len(a.Name) < 3 {
		a.Go.Name = strings.ToUpper(a.Name)
	}

	// make upper for sid,bid,rid
	if len(a.Name) < 4 {
		if strings.Contains(a.Name, "id") {
			a.Go.Name = strings.ToUpper(a.Name)
		}
	}

	helps.Field[`title{}`] = `Add custom title for index. title{Nice & Sweet}`
	helps.Field[`desc{}`] = `Add custom desc for index. desc{This field for success}`

	a.Go.Title = tools.Between(a.Comment, `title{`, `}`)
	a.Go.Desc = tools.Between(a.Comment, `desc{`, `}`)

	// init comments
	a.Comment = strings.ReplaceAll(a.Comment, "//", "")
	a.Comment = strings.TrimSpace(a.Comment)

	// go
	a.parseGolangOptions(tools.Between(a.Comment, "go{", "}"))

	// sql
	a.parseSQLOptions(tools.Between(a.Comment, "sql{", "}"))

	// clickhouse
	a.parseClickhouseOptions(tools.Between(a.Comment, "ch{", "}"))

	// proto
	a.parseProtoOptions(tools.Between(a.Comment, "proto{", "}"))

	// js
	a.parseJSOptions(tools.Between(a.Comment, "js{", "}"))

	// swift
	a.parseSwiftOptions(tools.Between(a.Comment, "swift{", "}"))

	// users
	a.parseUserOptions()

	// tags
	a.parseTagsOptions()

	// custom structs
	a.parseCustomStructsOptions()

	if a.Go.UpperCase {
		a.Go.Name = strings.ToUpper(a.Name)
	}

	a.FieldName = fmt.Sprintf("Field%s%s", tools.Title(settings.Go.StructName), a.Go.Name) //ProfileFieldID

	a.Go.Bits, a.Go.Align = GoTypeBits(a.Type)
	// fmt.Println(a.Type, a.Go.Bits, "bytes")

	// check sql stop words
	if settings.SQL.Table != "" && tools.SqlWord(a.Name) {
		panic(fmt.Sprintf("%s is a private sql word. Change it!", a.Name))
	}

	// if has ID int field (for sql)
	if a.SQL.Inc {
		settings.Fields.IncField = a.SQL.Name
	}

}

/*
bool	1	Boolean value
int8, uint8, byte	1	8-bit integer/byte
int16, uint16	2	16-bit integer
int32, uint32, float32, rune	4	32-bit integer, float, or rune
int, uint, uintptr, int64, uint64, float64	8	64-bit integer or float
Slices: 24 bytes (contain a pointer to data, length, and capacity).
Maps: 8 bytes (pointer to a hmap struct).
Channels: 8 bytes (pointer to a hchan struct).
Functions: 8 bytes (pointer to function's code).
Interfaces: 16 bytes (two pointers: type info and data pointer).
Strings: 16 bytes (two pointers: pointer to underlying byte array and length).
Empty Struct (struct{}): 0 bytes, used when no memory allocation is needed.

*/

// insert sql query
func GoTypeBits(v string) (bits, align int) {
	return sizeAlign64(v)
}

func sizeAlign64(t string) (size int, align int) {
	t = strings.TrimSpace(t)
	t = stripSpacesInside(t)

	// Pointers: *T
	if strings.HasPrefix(t, "*") {
		return 8, 8
	}

	// Arrays: [N]T
	if strings.HasPrefix(t, "[") {
		n, elemT, ok := parseArray(t)
		if ok {
			es, ea := sizeAlign64(elemT)
			// In Go, element size is a multiple of its alignment, so n*es is safe.
			return n * es, ea
		}
		return 8, 8
	}

	// Slices: []T
	if strings.HasPrefix(t, "[]") {
		// slice header: ptr, len, cap
		return 24, 8
	}

	// Strings
	if t == "string" {
		// string header: ptr, len
		return 16, 8
	}

	// Maps
	if strings.HasPrefix(t, "map[") {
		// map value is a pointer to runtime.hmap
		return 8, 8
	}

	// Channels
	if strings.HasPrefix(t, "chan") || strings.HasPrefix(t, "<-chan") || strings.HasPrefix(t, "chan<-") {
		return 8, 8
	}

	// Funcs
	if strings.HasPrefix(t, "func(") {
		return 8, 8
	}

	// Interfaces
	if t == "interface{}" || strings.HasPrefix(t, "interface{") {
		// empty interface is 16 bytes on 64-bit, non-empty also typically 16
		return 16, 8
	}

	// Common scalar types
	switch t {
	case "bool":
		return 1, 1
	case "byte", "uint8", "int8":
		return 1, 1
	case "uint16", "int16":
		return 2, 2
	case "rune", "uint32", "int32", "float32":
		return 4, 4
	case "uint64", "int64", "float64", "complex64":
		return 8, 8
	case "complex128":
		return 16, 8
	case "uint", "int", "uintptr":
		return 8, 8
	}

	// Common named time-like / atomic-like aliases aren’t handled here.
	// Unknown: assume pointer-sized (good default for many ref types / structs we can’t parse).
	return 8, 8
}

// Removes whitespace everywhere (handy for "map[string] int" etc)
func stripSpacesInside(s string) string {
	return strings.Join(strings.Fields(s), "")
}

// parseArray parses "[N]T" where N is a non-negative int literal.
func parseArray(t string) (n int, elem string, ok bool) {
	// Accept: [10]int, [0]*Foo, [32]byte, etc.
	// No support for [...]T (ellipsis) here.
	re := regexp.MustCompile(`^\[(\d+)\](.+)$`)
	m := re.FindStringSubmatch(t)
	if len(m) != 3 {
		return 0, "", false
	}
	v, err := strconv.Atoi(m[1])
	if err != nil || v < 0 {
		return 0, "", false
	}
	return v, m[2], true
}
