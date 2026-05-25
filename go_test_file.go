package main

import (
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

func goTestFile() []byte {
	var imports []string
	imports = append(imports, `"encoding/json"`)
	imports = append(imports, `"reflect"`)
	imports = append(imports, `"testing"`)
	if settings.Proto.Name != "" {
		imports = append(imports, `gproto "google.golang.org/protobuf/proto"`)
	}

	var list []string
	list = append(list, fmt.Sprintf("package %s", settings.PackageName))
	list = append(list, "")
	list = append(list, "import (")
	for _, x := range imports {
		list = append(list, "\t"+x)
	}
	list = append(list, ")")
	list = append(list, "")
	list = append(list, testFixture())
	if settings.Proto.Name != "" {
		list = append(list, "")
		list = append(list, protoTestFixture())
	}
	list = append(list, "")
	list = append(list, fmt.Sprintf("type test%sStd %s", settings.Go.StructName, settings.Go.StructName))
	list = append(list, "")
	list = append(list, fmt.Sprintf(`func Test%[1]sMarshal(t *testing.T) {
	item := test%[1]s()
	b := item.Marshal()
	if len(b) == 0 {
		t.Fatal("empty marshal result")
	}
}`, settings.Go.StructName))
	list = append(list, "")
	list = append(list, fmt.Sprintf(`func Test%[1]sUnmarshal(t *testing.T) {
	item := test%[1]s()
	b := item.Marshal()
	res := %[1]sUnmarshal(b)
	if res == nil {
		t.Fatal("unmarshal failed")
	}
	if !reflect.DeepEqual(item, *res) {
		t.Fatalf("unmarshal mismatch: %%#v != %%#v", item, *res)
	}
}`, settings.Go.StructName))

	if settings.Proto.Name != "" {
		list = append(list, "")
		list = append(list, fmt.Sprintf(`func Test%[1]sProtoMarshal(t *testing.T) {
	item := test%[2]s()
	b, err := gproto.Marshal(&item)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatal("empty proto marshal result")
	}
}`, settings.Go.StructName, settings.Proto.Name))
		list = append(list, "")
		list = append(list, fmt.Sprintf(`func Test%[1]sProtoUnmarshal(t *testing.T) {
	item := test%[2]s()
	b, err := gproto.Marshal(&item)
	if err != nil {
		t.Fatal(err)
	}
	var res %[2]s
	err = gproto.Unmarshal(b, &res)
	if err != nil {
		t.Fatal(err)
	}
	if !gproto.Equal(&item, &res) {
		t.Fatalf("proto unmarshal mismatch: %%#v != %%#v", item, res)
	}
}`, settings.Go.StructName, settings.Proto.Name))
	}

	if settings.Go.Gotiny {
		list = append(list, "")
		list = append(list, fmt.Sprintf(`func Test%[1]sGotinyMarshal(t *testing.T) {
	item := test%[1]s()
	b := item.MarshalGotiny()
	if len(b) == 0 {
		t.Fatal("empty gotiny marshal result")
	}
}`, settings.Go.StructName))
		list = append(list, "")
		list = append(list, fmt.Sprintf(`func Test%[1]sGotinyUnmarshal(t *testing.T) {
	item := test%[1]s()
	b := item.MarshalGotiny()
	res := Unmarshal%[1]sGotiny(b)
	if !reflect.DeepEqual(item, res) {
		t.Fatalf("gotiny unmarshal mismatch: %%#v != %%#v", item, res)
	}
}`, settings.Go.StructName))
	}

	list = append(list, "")
	list = append(list, benchmarkTests())
	if settings.Proto.Name != "" {
		list = append(list, "")
		list = append(list, protoBenchmarkTests())
	}
	if settings.Go.Gotiny {
		list = append(list, "")
		list = append(list, gotinyBenchmarkTests())
	}
	list = append(list, "")
	return []byte(strings.Join(list, "\n"))
}

func testFixture() string {
	var list []string
	list = append(list, fmt.Sprintf("func test%s() %s {", settings.Go.StructName, settings.Go.StructName))
	list = append(list, fmt.Sprintf("\treturn %s{", settings.Go.StructName))
	for _, x := range fields {
		list = append(list, fmt.Sprintf("\t\t%s: %s,", x.Go.Name, testValue(x)))
	}
	list = append(list, "\t}")
	list = append(list, "}")
	return strings.Join(list, "\n")
}

func protoTestFixture() string {
	var list []string
	list = append(list, fmt.Sprintf("func test%s() %s {", settings.Proto.Name, settings.Proto.Name))
	list = append(list, fmt.Sprintf("\treturn %s{", settings.Proto.Name))
	for _, x := range fields {
		value := protoTestValue(x)
		if value == "" {
			continue
		}
		list = append(list, fmt.Sprintf("\t\t%s: %s,", protoGoFieldName(x.Proto.Name), value))
	}
	list = append(list, "\t}")
	list = append(list, "}")
	return strings.Join(list, "\n")
}

func benchmarkTests() string {
	return fmt.Sprintf(`func Benchmark%[1]sMarshalStd(b *testing.B) {
	item := test%[1]s()
	std := test%[1]sStd(item)
	b.ReportAllocs()
	for b.Loop() {
		_, _ = json.Marshal(&std)
	}
}

func Benchmark%[1]sUnmarshalStd(b *testing.B) {
	item := test%[1]s()
	std := test%[1]sStd(item)
	data, err := json.Marshal(&std)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		var res test%[1]sStd
		err = json.Unmarshal(data, &res)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark%[1]sMarshal(b *testing.B) {
	item := test%[1]s()
	b.ReportAllocs()
	for b.Loop() {
		_ = item.Marshal()
	}
}

func Benchmark%[1]sUnmarshal(b *testing.B) {
	item := test%[1]s()
	data := item.Marshal()
	b.ReportAllocs()
	for b.Loop() {
		res := %[1]sUnmarshal(data)
		if res == nil {
			b.Fatal("unmarshal failed")
		}
	}
}`, settings.Go.StructName)
}

func protoBenchmarkTests() string {
	return fmt.Sprintf(`func Benchmark%[1]sProtoMarshal(b *testing.B) {
	item := test%[2]s()
	b.ReportAllocs()
	for b.Loop() {
		_, err := gproto.Marshal(&item)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark%[1]sProtoUnmarshal(b *testing.B) {
	item := test%[2]s()
	data, err := gproto.Marshal(&item)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		var res %[2]s
		err = gproto.Unmarshal(data, &res)
		if err != nil {
			b.Fatal(err)
		}
	}
}`, settings.Go.StructName, settings.Proto.Name)
}

func gotinyBenchmarkTests() string {
	return fmt.Sprintf(`func Benchmark%[1]sGotinyMarshal(b *testing.B) {
	item := test%[1]s()
	b.ReportAllocs()
	for b.Loop() {
		_ = item.MarshalGotiny()
	}
}

func Benchmark%[1]sGotinyUnmarshal(b *testing.B) {
	item := test%[1]s()
	data := item.MarshalGotiny()
	b.ReportAllocs()
	for b.Loop() {
		_ = Unmarshal%[1]sGotiny(data)
	}
}`, settings.Go.StructName)
}

func testValue(x *Field) string {
	t := testType(x)
	if x.Clickhouse.IP || strings.EqualFold(x.Name, "ip") {
		return `"127.0.0.1"`
	}
	switch t {
	case "string":
		return fmt.Sprintf("%q", x.Name+"_value")
	case "bool":
		return "true"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "byte", "time.Duration":
		return "42"
	case "float32", "float64":
		return "42.5"
	case "[]byte":
		return "[]byte(\"bytes_value\")"
	case "[]string":
		return "[]string{\"one\", \"two\", \"three\"}"
	case "[]int":
		return "[]int{1, 2, 3}"
	case "[]int8":
		return "[]int8{1, 2, 3}"
	case "[]int16":
		return "[]int16{1, 2, 3}"
	case "[]int32":
		return "[]int32{1, 2, 3}"
	case "[]int64":
		return "[]int64{1, 2, 3}"
	case "[]uint":
		return "[]uint{1, 2, 3}"
	case "[]uint8":
		return "[]uint8{1, 2, 3}"
	case "[]uint16":
		return "[]uint16{1, 2, 3}"
	case "[]uint32":
		return "[]uint32{1, 2, 3}"
	case "[]uint64":
		return "[]uint64{1, 2, 3}"
	case "[]float32":
		return "[]float32{1.5, 2.5, 3.5}"
	case "[]float64":
		return "[]float64{1.5, 2.5, 3.5}"
	case "map[string]string":
		return `map[string]string{"one": "first", "two": "second"}`
	case "map[string][]byte":
		return `map[string][]byte{"one": []byte("first"), "two": []byte("second")}`
	case "map[string]bool":
		return `map[string]bool{"one": true, "two": false}`
	case "map[string]int":
		return `map[string]int{"one": 1, "two": 2}`
	case "map[string]float64":
		return `map[string]float64{"one": 1.5, "two": 2.5}`
	case "map[string]any":
		return `map[string]any{"one": "first", "two": true}`
	case "map[int]string":
		return `map[int]string{1: "first", 2: "second"}`
	case "map[int]int":
		return `map[int]int{1: 11, 2: 22}`
	case "map[int]bool":
		return `map[int]bool{1: true, 2: false}`
	case "any", "interface{}":
		return `"any_value"`
	default:
		return fmt.Sprintf("%s{}", t)
	}
}

func protoTestValue(x *Field) string {
	t := testType(x)
	if x.Clickhouse.IP || strings.EqualFold(x.Name, "ip") {
		return `"127.0.0.1"`
	}
	switch t {
	case "string":
		return fmt.Sprintf("%q", x.Name+"_value")
	case "bool":
		return "true"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "byte", "time.Duration":
		return "42"
	case "float32", "float64":
		return "42.5"
	case "[]byte":
		return "[]byte(\"bytes_value\")"
	case "[]string":
		return "[]string{\"one\", \"two\", \"three\"}"
	case "[]int", "[]int64":
		return "[]int64{1, 2, 3}"
	case "[]int8", "[]int16", "[]int32":
		return "[]int32{1, 2, 3}"
	case "[]uint", "[]uint64":
		return "[]uint64{1, 2, 3}"
	case "[]uint8":
		return "[]byte{1, 2, 3}"
	case "[]uint16", "[]uint32":
		return "[]uint32{1, 2, 3}"
	case "[]float32":
		return "[]float32{1.5, 2.5, 3.5}"
	case "[]float64":
		return "[]float64{1.5, 2.5, 3.5}"
	case "map[string]string":
		return `map[string]string{"one": "first", "two": "second"}`
	case "map[string][]byte":
		return `map[string][]byte{"one": []byte("first"), "two": []byte("second")}`
	case "map[string]bool":
		return `map[string]bool{"one": true, "two": false}`
	case "map[string]int", "map[string]int64":
		return `map[string]int64{"one": 1, "two": 2}`
	case "map[string]float64":
		return `map[string]float64{"one": 1.5, "two": 2.5}`
	case "map[int]string":
		return `map[int64]string{1: "first", 2: "second"}`
	case "map[int]int":
		return `map[int64]int64{1: 11, 2: 22}`
	case "map[int]bool":
		return `map[int64]bool{1: true, 2: false}`
	default:
		return ""
	}
}

func testType(x *Field) string {
	if x.Go.Type != "" {
		return x.Go.Type
	}
	return x.Type
}

func protoGoFieldName(v string) string {
	var list []string
	for _, x := range strings.Split(protoFieldName(v), "_") {
		if x == "" {
			continue
		}
		list = append(list, tools.Title(x))
	}
	return strings.Join(list, "")
}
