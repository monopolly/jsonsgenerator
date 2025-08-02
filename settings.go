package main

import "fmt"

var settings = NewSettings()

func NewSettings() (a *Settings) {
	a = new(Settings)
	a.Imports = make(map[string]bool)
	return
}

type Settings struct {
	Debug bool

	Origin        string //type news struct{}
	PackageName   string
	IndexTypeName string //NewsIndexType

	Go struct {
		StructName  string //struct name: News
		NoOmit      bool
		NoPrefix    bool
		Gotiny      bool     //fast binary format
		MessagePack bool     //msgpack binary format
		EasyJson    bool     //generate easy json fast code gen
		Indexes     []string //string list all indexes [IndexNewsID, IndexNewsTitle...]
	}

	TS struct {
		Name string //struct name: news
	}
	JS struct {
		Name string //go js struct name: NewsJson
	}

	SQL struct {
		Class string //accountSQL: type accountSQL int{}
		// ClassVarName string //var AccountSQL accountSQL чтобы обращаться без инициализации
		Table     string //sql table from struct
		QueryName string //NewsQuery
	}

	Clickhouse struct {
		Table        string //sql table from struct
		Engine       string //
		Class        string //accountCHSQL: type accountCHSQL int{}
		ClassVarName string //var AccountCHSQL accountCHSQL чтобы обращаться без инициализации
		QueryName    string //NewsCHQuery
	}

	Swift struct {
		Model bool   //generate swift
		Enum  string //CategoryItem in swift
	}

	Options struct {
		Noinit bool //убирает New() вообще
		Demo   bool //demo json file with default values
	}

	Fields struct {
		IncField string //fields has ID int field (for sql insert generates)
	}

	File struct {
		path           string //test/news.go
		golang         string //test/news.go
		golangSQL      string //test/newsSQL.go
		golangSQLQuery string //test/newsQuery.go
		golangJSON     string //test/newsJson.go
		swift          string //test/news.swift
		enum           string //test/newsEnum.swift
		sql            string //test/news.sql
		clickhouse     string //test/news_ch.sql
		ts             string //test/news.ts
		js             string //test/news.js
		demo           string //test/news.json
	}

	Padding int //max len field title

	Imports map[string]bool
}

// github.com/something
func (a *Settings) AddImportPackage(v string) {
	a.Imports[fmt.Sprintf(`"%s"`, v)] = true
}

// jj "github.com/something"
func (a *Settings) AddImportRawPackage(v string) {
	a.Imports[v] = true
}
