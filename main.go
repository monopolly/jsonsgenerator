package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/monopolly/file"
)

var (
	fields []*Field //список всех полей
	// structed bool
	debug bool

	// indexFile bytes.Buffer
	// sqlFile   bytes.Buffer
	// helpFile bytes.Buffer
	// jsFile    bytes.Buffer //результат
	goFile bytes.Buffer //для структуры

	fieldnameSQLErrors []string //список полей которые нельзя использовать в sql

	A = func(v ...any) {
		if debug {
			fmt.Println(v...)
		}
	}
)

//go build -o $GOPATH/bin/jsons

func main() {

	// help init
	helps.init()

	//парсим входящие файл
	in := parsePath()

	//парсит файл и создает структуру
	settings.Parse(in)

	// go
	render()

}

// indexfile
func render() {

	// package
	goFile.Write(header())

	// import
	goFile.Write(importers())

	// help
	goFile.Write(helps.Render())

	//go
	goFile.Write(golang.Generate())

	// js
	if settings.JS.Name != "" {
		A("generate js")
		goFile.Write(golang.generateParseTupleToJsons())
		goFile.Write(jsn.generateJsonsInit())
		goFile.Write(jsn.generateJsonsGeneral())
		goFile.Write(jsn.generateJsonsFunctions())
	}

	// sql
	if settings.SQL.Table != "" {
		A("sql table", settings.SQL.Table)
		goFile.Write(sql.Generate())
		file.Save(settings.File.sql, sql.File())
	}

	// clickhouse
	if settings.Clickhouse.Table != "" {
		A("sql table", settings.Clickhouse.Table)
		goFile.Write(clickhouse.Generate())
		file.Save(settings.File.clickhouse, clickhouse.File())
	}

	// gosql
	file.Save(settings.File.golang, goFile.Bytes())

	// ts
	if settings.TS.Name != "" {
		file.Save(settings.File.ts, generateTypeScript())
		file.Save(settings.File.js, generateJSModel())
	}

	// swift
	if settings.Swift.Model {
		for name, v := range generateSwiftFiles() {
			file.Save(name, v)
		}
	}
	// swift
	if settings.Swift.Enum != "" {
		file.Save(settings.File.enum, generateSwiftEnum())
	}

	if settings.Options.Demo {
		file.Save(settings.File.demo, jsonDemo())
	}

	// msgp
	// msgpackGenerator()

	// easyjson
	// easyjsonGenerator()

	// install
	goModTidy()
	goFmt()
}

func parsePath() (in string) {
	files := os.Args
	if len(files) < 2 {
		log.Fatal("Jsons generate json golang struct\n", "Example:", files[0]+" file_struct.go")
		return
	}

	in = files[1]
	if !file.Exists(in) {
		in = in + ".go"
		if !file.Exists(in) {
			log.Fatal("Error: file not found", in)
		}
		return
	}

	if !strings.HasSuffix(in, ".go") {
		log.Fatal("Error: file must be .go")
		return
	}

	//in = "test.go"
	settings.File.path = path.Dir(in)
	settings.File.golang = strings.ReplaceAll(in, ".go", "GO.go")
	settings.File.sql = strings.ReplaceAll(in, ".go", ".sql")
	settings.File.clickhouse = strings.ReplaceAll(in, ".go", "_ch.sql")
	settings.File.ts = strings.ReplaceAll(in, ".go", ".ts")
	settings.File.js = strings.ReplaceAll(in, ".go", ".js")
	settings.File.demo = strings.ReplaceAll(in, ".go", ".json")
	settings.File.swift = strings.ReplaceAll(in, ".go", ".swift")
	settings.File.enum = strings.ReplaceAll(in, ".go", "Enum.swift")
	return
}

func header() []byte {
	return fmt.Appendf(nil, "package %s\n", settings.PackageName)
}

func goModTidy() {
	p := exec.Command("go", "mod", "tidy")
	p.Run()
}
func goFmt() {
	// go fmt
	p := exec.Command("go", "fmt", settings.File.golang)
	p.Run()
}

// func msgpackGenerator() {
// 	if !settings.Go.MessagePack {
// 		return
// 	}

// 	fmt.Println("Install msgp...")
// 	p := exec.Command("go", "install", "github.com/tinylib/msgp@latest")
// 	p.Run()

// 	// fmt.Println("generate msgp...")

// }

// func easyjsonGenerator() {
// 	if !settings.Go.EasyJson {
// 		return
// 	}

// 	// go get github.com/mailru/easyjson && go install github.com/mailru/easyjson/...@latest

// 	fmt.Println("Install easyjson...")
// 	p := exec.Command("go", "get", "github.com/mailru/easyjson", "&&", "go", "install", "github.com/mailru/easyjson/...@latest")
// 	p.Run()

// 	// easyjson <file>.go
// 	fmt.Println("Generate easyjson...")
// 	// g := exec.Command("easyjson", settings.file.golang, "-output_filename", settings.file.golang+"easyjson.go")
// 	g := exec.Command("easyjson", settings.file.golang)
// 	er := g.Run()
// 	if er != nil {
// 		fmt.Println(er)
// 	}

// }

var imports = map[string]bool{
	"context":                           true,
	"fmt":                               true,
	"reflect":                           true,
	"strings":                           true,
	"time":                              true,
	"errors":                            true,
	"github.com/niubaoshu/gotiny":       true,
	"github.com/vmihailenco/msgpack/v5": true,
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver": true,
	"github.com/json-iterator/go":                       true,
	"github.com/monopolly/cast":                         true,
	"github.com/jackc/pgx/v5/pgxpool":                   true,
	"github.com/monopolly/jsons":                        true,
}

// "github.com/jackc/pgx/v5/pgxpool"
//
// "github.com/monopolly/jsons"

func importers() []byte {

	list := []string{}
	for x := range imports {
		list = append(list, fmt.Sprintf(`"%s"`, x))
	}

	list = append(list, `jsoniter "github.com/json-iterator/go"`)

	p := fmt.Sprintf(`
	import (
		%s
	)
	`, strings.Join(list, "\n"))

	return []byte(p)
}
