package main

import (
	"fmt"
	"jj/tools"
	"strings"

	"github.com/monopolly/structs"
)

func (a *Settings) Parse(in string) {

	// ast parse
	res := structs.Parse(in)

	a.PackageName = res.Package
	a.Origin = res.Name

	a.Debug = res.Options["debug"] != nil
	helps.Struct["debug"] = "Debug mode"

	// options
	helps.Struct["lock"] = "Lock model for generation. Can't change model."
	if res.Options["lock"] != nil {
		panic("Locked for generate! Be carefull!")
	}

	a.Options.Noinit = res.Options["noinit"] != nil
	helps.Struct["noinit"] = "No New() init function for struct"

	a.Go.Gotiny = res.Options["gotiny"] != nil
	helps.Struct["gotiny"] = "Create gotiny marshal/unmarshal"

	a.Go.MessagePack = res.Options["msgp"] != nil
	helps.Struct["msgp"] = "Create message pack marshal/unmarshal"

	a.Go.NoOmit = res.Options["!omit"] != nil
	helps.Struct["!omit"] = "No omit tag for json"

	a.Go.StructName, _ = res.Options["go"].(string)
	helps.Struct["go=News"] = "Set golang struct names. Ex: go=News1 > type News1 struct{}"

	a.Go.StructName, _ = res.Options["go"].(string)
	helps.Struct["go=News"] = "Set golang struct names. Ex: go=News1 > type News1 struct{}"

	a.JS.Name, _ = res.Options["js"].(string)
	helps.Struct["js=NewsJson"] = "Change json struct names. Ex: js=NewsJson"

	a.TS.Name, _ = res.Options["ts"].(string)
	helps.Struct["ts=news"] = "Set typescript struct names. Ex: ts=NewsJson"

	a.Swift.Model = res.Options["swift"] != nil
	helps.Struct["swift"] = "Generate swift model"

	a.Swift.Enum, _ = res.Options["enum"].(string)
	helps.Struct["enum"] = "Generate swift enum"

	a.Options.Demo = res.Options["demo"] != nil
	helps.Struct["demo"] = "Generate demo json file with default values"

	a.Go.NoPrefix = res.Options["noprefix"] != nil
	helps.Struct["noprefix"] = "Generate simple index IndexID instead IndexNewsID"

	a.SQL.Table, _ = res.Options["sql"].(string)
	helps.Struct["sql=news"] = "Set sql table name. Ex: sql=accounts"

	a.Clickhouse.Table, _ = res.Options["ch"].(string)
	helps.Struct["ch=views"] = "Set clickhouse sql table name. Ex: ch=views"
	helps.Struct["chengine=MergeTree"] = "Optional, mergeTree by default"

	if a.Go.StructName == "" {
		a.Go.StructName = tools.Title(a.Origin)
	}

	// sql class
	// a.SQL.ClassVarName = fmt.Sprintf("%sSQL", a.Go.StructName)
	a.SQL.Class = fmt.Sprintf("%sSQL", a.Go.StructName) //fmt.Sprintf("sql%s", a.Go.StructName)
	a.SQL.QueryName = fmt.Sprintf("%sQuery", a.Go.StructName)

	// sql class
	a.Clickhouse.ClassVarName = fmt.Sprintf("%sCHSQL", a.Go.StructName)
	a.Clickhouse.Class = fmt.Sprintf("sqlCH%s", a.Go.StructName)
	a.Clickhouse.QueryName = fmt.Sprintf("%sCHQuery", a.Go.StructName)
	a.Clickhouse.Engine, _ = res.Options["chengine"].(string)

	// index type name
	a.IndexTypeName = fmt.Sprintf("%sIndexType", a.Go.StructName)

	// fields
	for _, x := range res.Fields {
		var p Field
		p.Name = x.Name
		p.MessagePack.Name = x.Name
		p.Comment = x.Comment
		p.Type = x.Type
		p.Go.Name = tools.Title(p.Name)
		p.Parse()

		if len(p.Name) > a.Padding {
			a.Padding = len(p.Name)
		}

		fields = append(fields, &p)
	}

	if len(fieldnameSQLErrors) > 0 {
		panic(fmt.Sprintf("%s — can't use as SQL field name", strings.Join(fieldnameSQLErrors, ", ")))
	}

	if a.Go.StructName != "" {
		return
	}

}
