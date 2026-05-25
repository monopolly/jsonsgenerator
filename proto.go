package main

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

var proto Proto

type Proto int

type protoImports struct {
	structpb  bool
	timestamp bool
}

func (a *Proto) File() []byte {
	var res bytes.Buffer
	imports := a.imports()
	var fieldList []string

	for _, x := range fields {
		if x.Proto.Skip {
			continue
		}
		fieldType := a.fieldType(x, imports)
		fieldList = append(fieldList, fmt.Sprintf("  %s %s = %d;", fieldType, protoFieldName(x.Proto.Name), len(fieldList)+1))
	}

	res.WriteString("syntax = \"proto3\";\n\n")
	res.WriteString(fmt.Sprintf("package %s;\n\n", settings.Proto.Package))
	res.WriteString(fmt.Sprintf("option go_package = \"./;%s\";\n", settings.Proto.Package))

	if imports.structpb || imports.timestamp {
		res.WriteString("\n")
	}
	if imports.structpb {
		res.WriteString("import \"google/protobuf/struct.proto\";\n")
	}
	if imports.timestamp {
		res.WriteString("import \"google/protobuf/timestamp.proto\";\n")
	}

	res.WriteString(fmt.Sprintf("\nmessage %s {\n", settings.Proto.Name))
	res.WriteString(strings.Join(fieldList, "\n"))
	res.WriteString("\n}\n")

	return res.Bytes()
}

func (a *Proto) Imports() (res []string) {
	imports := a.imports()
	res = append(res,
		`protoreflect "google.golang.org/protobuf/reflect/protoreflect"`,
		`protoimpl "google.golang.org/protobuf/runtime/protoimpl"`,
		"reflect",
		"sync",
		"unsafe",
	)
	if imports.structpb {
		res = append(res, `structpb "google.golang.org/protobuf/types/known/structpb"`)
	}
	if imports.timestamp {
		res = append(res, `timestamppb "google.golang.org/protobuf/types/known/timestamppb"`)
	}
	return
}

func (a *Proto) imports() *protoImports {
	imports := &protoImports{}
	for _, x := range fields {
		if x.Proto.Skip {
			continue
		}
		_ = a.fieldType(x, imports)
	}
	return imports
}

func (a *Proto) fieldType(x *Field, imports *protoImports) string {
	if x.Proto.Type != "" {
		a.detectImports(x.Proto.Type, imports)
		return x.Proto.Type
	}

	t := x.Type
	if x.Go.Type != "" {
		t = x.Go.Type
	}
	return a.goType(t, imports)
}

func (a *Proto) goType(t string, imports *protoImports) string {
	t = stripSpacesInside(strings.TrimSpace(t))
	for strings.HasPrefix(t, "*") {
		t = strings.TrimPrefix(t, "*")
	}

	if strings.HasPrefix(t, "map[") {
		key, value, ok := protoParseMap(t)
		if ok {
			return fmt.Sprintf("map<%s, %s>", a.mapKeyType(key), a.mapValueType(value, imports))
		}
	}

	if strings.HasPrefix(t, "[]") {
		elem := strings.TrimPrefix(t, "[]")
		if elem == "byte" || elem == "uint8" {
			return "bytes"
		}
		return "repeated " + a.scalarType(elem, imports)
	}

	if strings.HasPrefix(t, "[") {
		_, elem, ok := parseArray(t)
		if ok {
			if elem == "byte" || elem == "uint8" {
				return "bytes"
			}
			return "repeated " + a.scalarType(elem, imports)
		}
	}

	return a.scalarType(t, imports)
}

func (a *Proto) mapKeyType(t string) string {
	t = strings.TrimPrefix(t, "*")
	switch t {
	case "string":
		return "string"
	case "bool":
		return "bool"
	case "int", "int64":
		return "int64"
	case "int8", "int16", "int32", "rune":
		return "int32"
	case "uint", "uint64", "uintptr":
		return "uint64"
	case "uint8", "uint16", "uint32", "byte":
		return "uint32"
	default:
		return "string"
	}
}

func (a *Proto) mapValueType(t string, imports *protoImports) string {
	t = strings.TrimPrefix(t, "*")
	if strings.HasPrefix(t, "[]") {
		elem := strings.TrimPrefix(t, "[]")
		if elem == "byte" || elem == "uint8" {
			return "bytes"
		}
		imports.structpb = true
		return "google.protobuf.ListValue"
	}
	return a.scalarType(t, imports)
}

func (a *Proto) scalarType(t string, imports *protoImports) string {
	t = strings.TrimPrefix(t, "*")
	switch t {
	case "bool":
		return "bool"
	case "string":
		return "string"
	case "byte", "uint8", "uint16", "uint32":
		return "uint32"
	case "uint", "uint64", "uintptr":
		return "uint64"
	case "int8", "int16", "int32", "rune":
		return "int32"
	case "int", "int64":
		return "int64"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "any", "interface{}":
		imports.structpb = true
		return "google.protobuf.Value"
	case "time.Duration":
		return "int64"
	case "time.Time":
		imports.timestamp = true
		return "google.protobuf.Timestamp"
	case "json.RawMessage":
		return "bytes"
	default:
		return "bytes"
	}
}

func (a *Proto) detectImports(t string, imports *protoImports) {
	if strings.Contains(t, "google.protobuf.Value") || strings.Contains(t, "google.protobuf.Struct") || strings.Contains(t, "google.protobuf.ListValue") {
		imports.structpb = true
	}
	if strings.Contains(t, "google.protobuf.Timestamp") {
		imports.timestamp = true
	}
}

func protoParseMap(t string) (key, value string, ok bool) {
	if !strings.HasPrefix(t, "map[") {
		return "", "", false
	}
	depth := 0
	for i := len("map["); i < len(t); i++ {
		switch t[i] {
		case '[':
			depth++
		case ']':
			if depth == 0 {
				key = t[len("map["):i]
				value = t[i+1:]
				return key, value, key != "" && value != ""
			}
			depth--
		}
	}
	return "", "", false
}

func generatedGoBody(v []byte) []byte {
	s := string(v)
	importStart := strings.Index(s, "\nimport (")
	if importStart == -1 {
		packageStart := strings.Index(s, "\npackage ")
		if packageStart == -1 {
			return v
		}
		rest := s[packageStart+1:]
		lineEnd := strings.Index(rest, "\n")
		if lineEnd == -1 {
			return nil
		}
		return []byte("\n" + rest[lineEnd+1:])
	}
	rest := s[importStart:]
	importEnd := strings.Index(rest, "\n)\n")
	if importEnd == -1 {
		return v
	}
	return []byte("\n" + rest[importEnd+3:])
}

func generatedProtoGoBody(v []byte) []byte {
	s := string(generatedGoBody(v))
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "func init() { ") || !strings.HasSuffix(line, " }") {
			continue
		}
		call := strings.TrimSuffix(strings.TrimPrefix(line, "func init() { "), " }")
		initName := fmt.Sprintf("init%sProto", settings.Go.StructName)
		lines[i] = fmt.Sprintf(`var _ = %s()

func %s() struct{} {
	%s
	return struct{}{}
}`, initName, initName, call)
		return []byte(strings.Join(lines, "\n"))
	}
	return []byte(s)
}

func protoFieldName(v string) string {
	var out []rune
	runes := []rune(v)
	for i, r := range runes {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if unicode.IsUpper(r) {
				if i > 0 {
					prev := runes[i-1]
					nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
					if prev != '_' && (unicode.IsLower(prev) || unicode.IsDigit(prev) || nextLower) {
						out = append(out, '_')
					}
				}
				r = unicode.ToLower(r)
			}
			out = append(out, r)
			continue
		}
		if len(out) > 0 && out[len(out)-1] != '_' {
			out = append(out, '_')
		}
	}
	res := strings.Trim(string(out), "_")
	if res == "" {
		return "field"
	}
	if unicode.IsDigit([]rune(res)[0]) {
		return "field_" + res
	}
	return res
}
