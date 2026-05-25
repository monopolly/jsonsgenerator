package main

import (
	"fmt"
	"strings"

	"github.com/monopolly/cast"
	"github.com/monopolly/jsonsgenerator/tools"

	"github.com/monopolly/structs"
)

func (a *Settings) Parse(in string) {

	// ast parse
	res := structs.Parse(in)

	// package name "model"
	a.PackageName = res.Package

	// origin struct name "news"
	a.Origin = res.Name

	// is debug
	a.Debug = res.Options["debug"] != nil
	helps.Struct["debug"] = "Debug mode"

	// lock
	helps.Struct["lock"] = "Lock model for generation. Can't change model."
	if res.Options["lock"] != nil {
		panic("Locked for generate! Be carefull!")
	}

	// no init
	a.Options.Noinit = res.Options["noinit"] != nil
	helps.Struct["noinit"] = "No New() init function for struct"

	// gotiny
	a.Go.Gotiny = res.Options["gotiny"] != nil
	helps.Struct["gotiny"] = "Create gotiny marshal/unmarshal"

	// message pack
	a.Go.MessagePack = res.Options["msgp"] != nil
	helps.Struct["msgp"] = "Create message pack marshal/unmarshal"

	// no omit for json
	a.Go.NoOmit = res.Options["!omit"] != nil
	helps.Struct["!omit"] = "No omit tag for json"

	// goname
	a.Go.StructName = optionName(res.Options, "go", tools.Title(a.Origin))
	helps.Struct["go=News"] = "Set golang struct names. Ex: go=News1 > type News1 struct{}"

	// js
	a.JS.Name = optionName(res.Options, "js", fmt.Sprintf("%sJson", tools.Title(a.Origin)))
	helps.Struct["js=NewsJson"] = "Change json struct names. Ex: js=NewsJson"

	// proto
	a.Proto.Name = optionName(res.Options, "proto", fmt.Sprintf("%sProto", a.Go.StructName))
	a.Proto.Package = a.PackageName
	helps.Struct["proto=NewsProto"] = "Generate proto file and compile to current package"

	// ts
	a.TS.Name = cast.String(res.Options["ts"])
	helps.Struct["ts=news"] = "Set typescript struct names. Ex: ts=NewsJson"

	// swift
	a.Swift.Model = res.Options["swift"] != nil
	helps.Struct["swift"] = "Generate swift model"

	// swift enum
	a.Swift.Enum = cast.String(res.Options["enum"])
	helps.Struct["enum"] = "Generate swift enum"

	// generate json demo
	a.Options.Demo = res.Options["demo"] != nil
	helps.Struct["demo"] = "Generate demo json file with default values"

	// generate json demo
	a.Options.Optimize = res.Options["optimize"] != nil
	helps.Struct["optimize"] = "Try to create golang padding optimized struct"

	// generate tests
	a.Options.Test = res.Options["test"] != nil
	helps.Struct["test"] = "Generate go tests for generated model"

	// simple index name NewsIndexID > IndexID (if 1 struct in package)
	a.Go.NoPrefix = res.Options["noprefix"] != nil
	helps.Struct["noprefix"] = "Generate simple index IndexID instead IndexNewsID"

	// sql
	a.SQL.Table = optionName(res.Options, "sql", a.Origin)
	helps.Struct["sql=news"] = "Set sql table name. Ex: sql=accounts"

	// clickhouse
	a.Clickhouse.Table = optionName(res.Options, "ch", a.Origin)
	a.Clickhouse.Engine = cast.String(res.Options["chengine"])
	helps.Struct["ch=views"] = "Set clickhouse sql table name. Ex: ch=views"
	helps.Struct["chengine=MergeTree"] = "Optional, mergeTree by default"

	if a.Go.StructName == "" {
		a.Go.StructName = tools.Title(a.Origin)
	}

	// sql class
	a.SQL.Class = fmt.Sprintf("%sSQL", a.Go.StructName) //fmt.Sprintf("sql%s", a.Go.StructName)
	a.SQL.QueryName = fmt.Sprintf("%sQuery", a.Go.StructName)

	// clickhouse class
	a.Clickhouse.Class = fmt.Sprintf("%sCQL", a.Go.StructName)
	a.Clickhouse.QueryName = fmt.Sprintf("%sClickhouseQuery", a.Go.StructName)

	// indextype
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

	// fmt.Println("before", fields)
	a.StructSize = StructSize(fields)
	if a.Options.Optimize {
		Optimize(fields)
	}

	// fmt.Println("after", fields)

	if len(fieldnameSQLErrors) > 0 {
		panic(fmt.Sprintf("%s — can't use as SQL field name", strings.Join(fieldnameSQLErrors, ", ")))
	}

	if a.Go.StructName != "" {
		return
	}

}

func optionName(options map[string]any, key, defaultValue string) string {
	v, ok := options[key]
	if !ok {
		return ""
	}
	switch value := v.(type) {
	case bool:
		if value {
			return defaultValue
		}
		return ""
	case string:
		if value == "" {
			return defaultValue
		}
		return value
	default:
		res := cast.String(v)
		if res == "" {
			return defaultValue
		}
		return res
	}
}
