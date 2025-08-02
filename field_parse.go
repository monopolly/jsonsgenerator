package main

import (
	"fmt"

	"github.com/monopolly/jsonsgenerator/tools"
	"strings"
)

func (a *Field) Parse() {

	a.Go.Tags = make(map[string]bool)
	a.Go.UserLists = make(map[string][]string)

	//потому что go не дает type ставить
	if a.Name == "types" {
		a.Name = "type"
	}

	a.SQL.Name = a.Name

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

	// js
	a.parseJSOptions(tools.Between(a.Comment, "js{", "}"))

	// swift
	a.parseSwiftOptions(tools.Between(a.Comment, "swift{", "}"))

	// users
	a.parseUserOptions()

	// tags
	a.parseTagsOptions()

	if a.Go.UpperCase {
		a.Go.Name = strings.ToUpper(a.Name)
	}

	a.FieldName = fmt.Sprintf("Field%s%s", tools.Title(settings.Go.StructName), a.Go.Name) //ProfileFieldID

	// check sql stop words
	if settings.SQL.Table != "" && tools.SqlWord(a.Name) {
		panic(fmt.Sprintf("%s is a private sql word. Change it!", a.Name))
	}

	// if has ID int field (for sql)
	if a.SQL.Inc {
		settings.Fields.IncField = a.SQL.Name
	}

}
