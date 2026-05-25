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
	fields []*Field // list of all fields
	// structed bool
	debug bool

	// indexFile bytes.Buffer
	// sqlFile   bytes.Buffer
	// helpFile bytes.Buffer
	// jsFile    bytes.Buffer // result
	goFile bytes.Buffer // for struct output

	fieldnameSQLErrors []string // fields that cannot be used in SQL

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

	// parse input file
	in := parsePath()

	// parse file and create struct metadata
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

	if settings.Proto.Name != "" {
		file.Save(settings.File.proto, proto.File())
		goFile.Write(protoCompile())
	}

	// gosql
	file.Save(settings.File.golang, goFile.Bytes())
	goFmt()
	goFile.Write(easyjsonCompile())
	file.Save(settings.File.golang, addGoImports(goFile.Bytes(), easyjsonImports()))

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
	if settings.Options.Test {
		file.Save(settings.File.golangTest, goTestFile())
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
	settings.File.golangTest = strings.ReplaceAll(settings.File.golang, ".go", "_test.go")
	settings.File.sql = strings.ReplaceAll(in, ".go", ".sql")
	settings.File.clickhouse = strings.ReplaceAll(in, ".go", ".cql")
	settings.File.proto = strings.ReplaceAll(in, ".go", ".proto")
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
	args := []string{"fmt"}
	if file.Exists(settings.File.golang) {
		args = append(args, settings.File.golang)
	}
	if file.Exists(settings.File.golangTest) {
		args = append(args, settings.File.golangTest)
	}
	if len(args) == 1 {
		return
	}
	p := exec.Command("go", args...)
	p.Run()
}

func easyjsonCompile() []byte {
	outfile := path.Join(settings.File.path, fmt.Sprintf("%s_easy.go", settings.Origin))
	p := exec.Command(
		"easyjson",
		"-output_filename="+outfile,
		settings.File.golang,
	)
	out, err := p.CombinedOutput()
	if err != nil {
		log.Fatalf("easyjson failed: %v\n%s", err, string(out))
	}

	body, err := os.ReadFile(outfile)
	if err != nil {
		log.Fatal(err)
	}
	_ = os.Remove(outfile)
	res := generatedGoBody(body)
	res = append(res, golang.EasyJsonMarshal()...)
	res = append(res, golang.EasyJsonUnmarshal()...)
	return res
}

func easyjsonImports() []string {
	return []string{
		`json "encoding/json"`,
		`easyjson "github.com/mailru/easyjson"`,
		`jlexer "github.com/mailru/easyjson/jlexer"`,
		`jwriter "github.com/mailru/easyjson/jwriter"`,
	}
}

func addGoImports(src []byte, imports []string) []byte {
	s := string(src)
	pos := strings.Index(s, "import (\n")
	if pos == -1 {
		return src
	}
	insertPos := pos + len("import (\n")
	var list []string
	for _, x := range imports {
		if strings.Contains(s, x) {
			continue
		}
		list = append(list, "\t"+x+"\n")
	}
	if len(list) == 0 {
		return src
	}
	s = s[:insertPos] + strings.Join(list, "") + s[insertPos:]
	return []byte(s)
}

func protoCompile() []byte {
	p := exec.Command(
		"protoc",
		"--proto_path="+settings.File.path,
		"--go_out="+settings.File.path,
		"--go_opt=paths=source_relative",
		path.Base(settings.File.proto),
	)
	out, err := p.CombinedOutput()
	if err != nil {
		log.Fatalf("protoc failed: %v\n%s", err, string(out))
	}

	generated := strings.TrimSuffix(settings.File.proto, ".proto") + ".pb.go"
	body, err := os.ReadFile(generated)
	if err != nil {
		log.Fatal(err)
	}
	_ = os.Remove(generated)
	return generatedProtoGoBody(body)
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

var imports = map[string]bool{}

// "github.com/jackc/pgx/v5/pgxpool"
//
// "github.com/monopolly/jsons"

func importers() []byte {

	imports = map[string]bool{
		"fmt":                        true,
		"github.com/monopolly/cast":  true,
		"github.com/monopolly/jsons": true,
	}
	imports["jsoniter \"github.com/json-iterator/go\""] = true

	if settings.SQL.Table != "" {
		imports["context"] = true
		imports["errors"] = true
		imports["strings"] = true
		imports["github.com/jackc/pgx/v5/pgxpool"] = true
	}
	if settings.Clickhouse.Table != "" {
		imports["context"] = true
		imports["errors"] = true
		imports["strings"] = true
		imports["sync"] = true
		imports["time"] = true
		imports["github.com/ClickHouse/clickhouse-go/v2/lib/driver"] = true
		if hasClickhouseIPFields() {
			imports["net"] = true
		}
	}
	if settings.JS.Name != "" {
		for _, x := range fields {
			fieldType := x.Type
			if x.Go.Type != "" {
				fieldType = x.Go.Type
			}
			if strings.HasPrefix(fieldType, "map[") {
				imports["reflect"] = true
				break
			}
		}
	}
	if settings.Go.Gotiny {
		imports["github.com/niubaoshu/gotiny"] = true
	}
	if settings.Go.MessagePack {
		imports["github.com/vmihailenco/msgpack/v5"] = true
	}
	if settings.Proto.Name != "" {
		for _, x := range proto.Imports() {
			imports[x] = true
		}
	}
	for _, x := range fields {
		if strings.Contains(x.Type, "time.") || strings.Contains(x.Go.Type, "time.") {
			imports["time"] = true
		}
	}

	list := []string{}
	for x := range imports {
		switch strings.Contains(x, " ") {
		case true:
			list = append(list, x)
		case false:
			list = append(list, fmt.Sprintf(`"%s"`, x))
		}
	}

	p := fmt.Sprintf(`
	import (
		%s
	)
	`, strings.Join(list, "\n"))

	return []byte(p)
}
