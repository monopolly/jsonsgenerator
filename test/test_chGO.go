package news

import (
	"context"
	json "encoding/json"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/jackc/pgx/v5/pgxpool"
	jsoniter "github.com/json-iterator/go"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	"github.com/monopolly/cast"
	"github.com/monopolly/jsons"
	"github.com/niubaoshu/gotiny"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"
)

/*
	Auto-generate. Do not change!
	Help to auto-generated correct models for golang structures.
	Sergey Keplin (c) 2026
	a@senthy.com, lava.mobi@gmail.com, https://t.me/martinprestone
	github.com/monopolly/jsons
	Ex: //! [] go=News js=NewsJson sql=news noinit up...

	STRUCT:
	!omit                         No omit tag for json
	demo                          Generate demo json file with default values
	noinit                        No New() init function for struct
	go=News                       Set golang struct names. Ex: go=News1 > type News1 struct{}
	js=NewsJson                   Change json struct names. Ex: js=NewsJson
	ts=news                       Set typescript struct names. Ex: ts=NewsJson
	optimize                      Try to create golang padding optimized struct
	sql=news                      Set sql table name. Ex: sql=accounts
	lock                          Lock model for generation. Can't change model.
	gotiny                        Create gotiny marshal/unmarshal
	msgp                          Create message pack marshal/unmarshal
	proto=NewsProto               Generate proto file and compile to current package
	test                          Generate go tests for generated model
	noprefix                      Generate simple index IndexID instead IndexNewsID
	chengine=MergeTree            Optional, mergeTree by default
	debug                         Debug mode
	swift                         Generate swift model
	enum                          Generate swift enum
	ch=views                      Set clickhouse sql table name. Ex: ch=views

	Field:
	title{}                       Add custom title for index. title{Nice & Sweet}
	desc{}                        Add custom desc for index. desc{This field for success}

	GO:
	@                             Add custom structs for fields. Ex: id int //@public @team...
	up                            Make uppercase for functions
	must                          Create one validation function for all must fields
	name=""                       Replace golang struct name. Ex: name="NewsList"
	desc=""                       Add custom desc for index
	req{}                         Add user required fields
	#                             Add lists for fields. Ex: id int //#readonly #must...
	nofunc                        Do not create any jsons functions for fiels
	type=""                       Replace golang type. Ex: type="[]*News"
	title=""                      Add custom title for index
	title{}                       Add field user title title{Nice}
	desc{}                        Add field user title desc{Use it nice}

	SQL:
	unique="groupname"            Add unique fields constrains by group (you need set group name, then generator join fields). Ex: unique="group1" unique="group2"
	renames="oldname"             Create an alter table record in SQL file. Rename table column.
	index                         Create simple default index or gin for jsonb
	altertable                    Add field line Alter table to SQL file with current time comment
	idx                           Add index fields by group name. Ex: idx="nameIndex" idx="credsIndex"
	defaults                      Add default value based on field type
	skip                          Skip field for sql queries
	type="jsonb"                  Rewrite sql type for field. Ex: type="jsonb"
	add="primary key"             Append sql for field. Ex: replace="primary key"
	ver="v4"                      Create an alter table record in SQL file. Add new column. With new version in comment. Ex: ver="2"
	noinsert                      Do use field for insert function
	search                        Add tsvector index by group. Ex: search="tsv": tsv tsvector GENERATED ALWAYS AS (to_tsvector('simple', title || ' ' || brand)). For search: SELECT brand, title FROM assets WHERE search @@ to_tsquery('english', 'f8');
	unix                          Add default value: extract(epoch from now())
	replace="bigint primary key"  Rewrite sql for field. Ex: replace="bigint primary key"

	Clickhouse:
	low                           Wrap field type in LowCardinality(...)
	type=ip                       Use ClickHouse IPv4 type for IP address fields

	PROTO:
	name="field"                  Rewrite proto field name
	skip                          Skip field for proto
	type="string"                 Rewrite proto type for field. Ex: type="google.protobuf.Value"

	SWIFT:
	skip                          Ignore field for Swift
	type=""                       Replace type for Swift model. Ex: type="string"
	file=""                       Create another swift file for this field Swift model. file="f1", file="f2"
	must                          Swift field with required values

	JSON:
	time                          Create convert jsons function for unixtime fields
	raw                           Set Raw json function inside field
	name=""                       Replace json field name. Ex: name="sid"
	skip                          Skip json field
	inc                           Add inc jsons function for numbers fields
	bool                          Add set jsons function for bool fields
*/

// field type
type StatIndexType int

// int index
const (
	IndexStatID = StatIndexType(iota)
	IndexStatUID
	IndexStatType
	IndexStatName
	IndexStatSource
	IndexStatIP
	IndexStatImage
	IndexStatStack
	IndexStatUints
	IndexStatUints8
	IndexStatUints16
	IndexStatUints32
	IndexStatUints64
	IndexStatInts8
	IndexStatInts16
	IndexStatInts32
	IndexStatInts64
)

// string index
const (
	FieldStatID      = "id"      // int sql{inc}  #readonly
	FieldStatUID     = "uid"     // int ch{primary}
	FieldStatType    = "type"    // int ch{low}
	FieldStatName    = "name"    // string
	FieldStatSource  = "source"  // string
	FieldStatIP      = "ip"      // string
	FieldStatImage   = "image"   // bool
	FieldStatStack   = "stack"   // []string
	FieldStatUints   = "uints"   // []uint
	FieldStatUints8  = "uints8"  // []uint8
	FieldStatUints16 = "uints16" // []uint16
	FieldStatUints32 = "uints32" // []uint32
	FieldStatUints64 = "uints64" // []uint64
	FieldStatInts8   = "ints8"   // []int8
	FieldStatInts16  = "ints16"  // []int16
	FieldStatInts32  = "ints32"  // []int32
	FieldStatInts64  = "ints64"  // []int64
)

// index func
func StatIndexes() []StatIndexType {
	return []StatIndexType{IndexStatID, IndexStatUID, IndexStatType, IndexStatName, IndexStatSource, IndexStatIP, IndexStatImage, IndexStatStack, IndexStatUints, IndexStatUints8, IndexStatUints16, IndexStatUints32, IndexStatUints64, IndexStatInts8, IndexStatInts16, IndexStatInts32, IndexStatInts64}

}

// 160 bytes (go padding)
//
//easyjson:json
type Stat struct {
	ID      int      `json:"id,omitempty"`      // sql{inc}  #readonly
	UID     int      `json:"uid,omitempty"`     // ch{primary}
	Type    int      `json:"type,omitempty"`    // ch{low}
	Name    string   `json:"name,omitempty"`    //
	Source  string   `json:"source,omitempty"`  //
	IP      string   `json:"ip,omitempty"`      //
	Image   bool     `json:"image,omitempty"`   //
	Stack   []string `json:"stack,omitempty"`   //
	Uints   []uint   `json:"uints,omitempty"`   //
	Uints8  []uint8  `json:"uints8,omitempty"`  //
	Uints16 []uint16 `json:"uints16,omitempty"` //
	Uints32 []uint32 `json:"uints32,omitempty"` //
	Uints64 []uint64 `json:"uints64,omitempty"` //
	Ints8   []int8   `json:"ints8,omitempty"`   //
	Ints16  []int16  `json:"ints16,omitempty"`  //
	Ints32  []int32  `json:"ints32,omitempty"`  //
	Ints64  []int64  `json:"ints64,omitempty"`  //
}

// Parse []any to ID struct
func ParseStatToStruct(r []any) (a *Stat) {
	a = new(Stat)

	for pos, x := range r {
		switch StatIndexType(pos) {
		case IndexStatID:
			cast.Convert(&a.ID, x) //int
		case IndexStatUID:
			cast.Convert(&a.UID, x) //int
		case IndexStatType:
			cast.Convert(&a.Type, x) //int
		case IndexStatName:
			cast.Convert(&a.Name, x) //string
		case IndexStatSource:
			cast.Convert(&a.Source, x) //string
		case IndexStatIP:
			cast.Convert(&a.IP, x) //string
		case IndexStatImage:
			cast.Convert(&a.Image, x) //bool
		case IndexStatStack:
			cast.Convert(&a.Stack, x) //[]string
		case IndexStatUints:
			b, _ := jsoniter.Marshal(x)
			_ = jsoniter.Unmarshal(b, &a.Uints) //[]uint
		case IndexStatUints8:
			b, _ := jsoniter.Marshal(x)
			_ = jsoniter.Unmarshal(b, &a.Uints8) //[]uint8
		case IndexStatUints16:
			b, _ := jsoniter.Marshal(x)
			_ = jsoniter.Unmarshal(b, &a.Uints16) //[]uint16
		case IndexStatUints32:
			b, _ := jsoniter.Marshal(x)
			_ = jsoniter.Unmarshal(b, &a.Uints32) //[]uint32
		case IndexStatUints64:
			b, _ := jsoniter.Marshal(x)
			_ = jsoniter.Unmarshal(b, &a.Uints64) //[]uint64
		case IndexStatInts8:
			b, _ := jsoniter.Marshal(x)
			_ = jsoniter.Unmarshal(b, &a.Ints8) //[]int8
		case IndexStatInts16:
			b, _ := jsoniter.Marshal(x)
			_ = jsoniter.Unmarshal(b, &a.Ints16) //[]int16
		case IndexStatInts32:
			b, _ := jsoniter.Marshal(x)
			_ = jsoniter.Unmarshal(b, &a.Ints32) //[]int32
		case IndexStatInts64:
			b, _ := jsoniter.Marshal(x)
			_ = jsoniter.Unmarshal(b, &a.Ints64) //[]int64
		}
	}
	return
}

// Tuple create an array from struct
func (a *Stat) Tuple() (r []any) {
	return []any{a.ID, a.UID, a.Type, a.Name, a.Source, a.IP, a.Image, a.Stack, a.Uints, a.Uints8, a.Uints16, a.Uints32, a.Uints64, a.Ints8, a.Ints16, a.Ints32, a.Ints64}
}

// Tuple create an array from struct
func (a *Stat) sqlTuple() (r []any) {
	return []any{a.UID, a.Type, a.Name, a.Source, a.IP, a.Image, a.Stack, a.Uints, a.Uints8, a.Uints16, a.Uints32, a.Uints64, a.Ints8, a.Ints16, a.Ints32, a.Ints64}
}

// Tuple create an array from struct
func (a *Stat) sqlAllTuples() (r []any) {
	return []any{a.ID, a.UID, a.Type, a.Name, a.Source, a.IP, a.Image, a.Stack, a.Uints, a.Uints8, a.Uints16, a.Uints32, a.Uints64, a.Ints8, a.Ints16, a.Ints32, a.Ints64}
}

// Tuple create an array from struct for clickhouse
func (a *Stat) clickhouseTuple() (r []any) {
	return []any{a.ID, a.UID, a.Type, a.Name, a.Source, net.ParseIP(a.IP), a.Image, a.Stack, a.Uints, a.Uints8, a.Uints16, a.Uints32, a.Uints64, a.Ints8, a.Ints16, a.Ints32, a.Ints64}
}

// update struct with function
func (a *Stat) Update(k string, x any) {
	switch k {
	case "id":
		cast.Convert(&a.ID, x) //int
	case "uid":
		cast.Convert(&a.UID, x) //int
	case "type":
		cast.Convert(&a.Type, x) //int
	case "name":
		cast.Convert(&a.Name, x) //string
	case "source":
		cast.Convert(&a.Source, x) //string
	case "ip":
		cast.Convert(&a.IP, x) //string
	case "image":
		cast.Convert(&a.Image, x) //bool
	case "stack":
		cast.Convert(&a.Stack, x) //[]string
	case "uints":
		b, _ := jsoniter.Marshal(x)
		_ = jsoniter.Unmarshal(b, &a.Uints) //[]uint
	case "uints8":
		b, _ := jsoniter.Marshal(x)
		_ = jsoniter.Unmarshal(b, &a.Uints8) //[]uint8
	case "uints16":
		b, _ := jsoniter.Marshal(x)
		_ = jsoniter.Unmarshal(b, &a.Uints16) //[]uint16
	case "uints32":
		b, _ := jsoniter.Marshal(x)
		_ = jsoniter.Unmarshal(b, &a.Uints32) //[]uint32
	case "uints64":
		b, _ := jsoniter.Marshal(x)
		_ = jsoniter.Unmarshal(b, &a.Uints64) //[]uint64
	case "ints8":
		b, _ := jsoniter.Marshal(x)
		_ = jsoniter.Unmarshal(b, &a.Ints8) //[]int8
	case "ints16":
		b, _ := jsoniter.Marshal(x)
		_ = jsoniter.Unmarshal(b, &a.Ints16) //[]int16
	case "ints32":
		b, _ := jsoniter.Marshal(x)
		_ = jsoniter.Unmarshal(b, &a.Ints32) //[]int32
	case "ints64":
		b, _ := jsoniter.Marshal(x)
		_ = jsoniter.Unmarshal(b, &a.Ints64) //[]int64
	}
}

// get struct value with function
func (a *Stat) Get(k string) (v any) {
	switch k {
	case "id":
		return a.ID //int
	case "uid":
		return a.UID //int
	case "type":
		return a.Type //int
	case "name":
		return a.Name //string
	case "source":
		return a.Source //string
	case "ip":
		return a.IP //string
	case "image":
		return a.Image //bool
	case "stack":
		return a.Stack //[]string
	case "uints":
		return a.Uints //[]uint
	case "uints8":
		return a.Uints8 //[]uint8
	case "uints16":
		return a.Uints16 //[]uint16
	case "uints32":
		return a.Uints32 //[]uint32
	case "uints64":
		return a.Uints64 //[]uint64
	case "ints8":
		return a.Ints8 //[]int8
	case "ints16":
		return a.Ints16 //[]int16
	case "ints32":
		return a.Ints32 //[]int32
	case "ints64":
		return a.Ints64 //[]int64
	}
	return
}

// get any struct value as string
func (a *Stat) String(k string) (v string) {
	switch k {
	case "id":
		return fmt.Sprint(a.ID) //int
	case "uid":
		return fmt.Sprint(a.UID) //int
	case "type":
		return fmt.Sprint(a.Type) //int
	case "name":
		return fmt.Sprint(a.Name) //string
	case "source":
		return fmt.Sprint(a.Source) //string
	case "ip":
		return fmt.Sprint(a.IP) //string
	case "image":
		return fmt.Sprint(a.Image) //bool
	case "stack":
		return fmt.Sprint(a.Stack) //[]string
	case "uints":
		return fmt.Sprint(a.Uints) //[]uint
	case "uints8":
		return fmt.Sprint(a.Uints8) //[]uint8
	case "uints16":
		return fmt.Sprint(a.Uints16) //[]uint16
	case "uints32":
		return fmt.Sprint(a.Uints32) //[]uint32
	case "uints64":
		return fmt.Sprint(a.Uints64) //[]uint64
	case "ints8":
		return fmt.Sprint(a.Ints8) //[]int8
	case "ints16":
		return fmt.Sprint(a.Ints16) //[]int16
	case "ints32":
		return fmt.Sprint(a.Ints32) //[]int32
	case "ints64":
		return fmt.Sprint(a.Ints64) //[]int64
	}
	return
}

// Struct to json
func (a *Stat) ToJson() (r []byte) {
	js := jsons.Create().
		Add(FieldStatID, a.ID).
		Add(FieldStatUID, a.UID).
		Add(FieldStatType, a.Type).
		Add(FieldStatName, a.Name).
		Add(FieldStatSource, a.Source).
		Add(FieldStatIP, a.IP).
		Add(FieldStatImage, a.Image).
		Add(FieldStatStack, a.Stack).
		Add(FieldStatUints, a.Uints).
		Add(FieldStatUints8, a.Uints8).
		Add(FieldStatUints16, a.Uints16).
		Add(FieldStatUints32, a.Uints32).
		Add(FieldStatUints64, a.Uints64).
		Add(FieldStatInts8, a.Ints8).
		Add(FieldStatInts16, a.Ints16).
		Add(FieldStatInts32, a.Ints32).
		Add(FieldStatInts64, a.Ints64)
	return js.Bytes()
}

func StatReadonlyList() []StatIndexType {
	return []StatIndexType{IndexStatID}
}

func (a StatIndexType) Readonly() bool {
	switch a {
	case IndexStatID:
		return true
	default:
		return false
	}
}

// key index string
func (a StatIndexType) String() string {
	switch a {
	case IndexStatID:
		return "id"
	case IndexStatUID:
		return "uid"
	case IndexStatType:
		return "type"
	case IndexStatName:
		return "name"
	case IndexStatSource:
		return "source"
	case IndexStatIP:
		return "ip"
	case IndexStatImage:
		return "image"
	case IndexStatStack:
		return "stack"
	case IndexStatUints:
		return "uints"
	case IndexStatUints8:
		return "uints8"
	case IndexStatUints16:
		return "uints16"
	case IndexStatUints32:
		return "uints32"
	case IndexStatUints64:
		return "uints64"
	case IndexStatInts8:
		return "ints8"
	case IndexStatInts16:
		return "ints16"
	case IndexStatInts32:
		return "ints32"
	case IndexStatInts64:
		return "ints64"
	default:
		return ""
	}
}

// key index string
func (a StatIndexType) SQLName() string {
	switch a {
	case IndexStatID:
		return "id"
	case IndexStatUID:
		return "uid"
	case IndexStatType:
		return "type"
	case IndexStatName:
		return "name"
	case IndexStatSource:
		return "source"
	case IndexStatIP:
		return "ip"
	case IndexStatImage:
		return "image"
	case IndexStatStack:
		return "stack"
	case IndexStatUints:
		return "uints"
	case IndexStatUints8:
		return "uints8"
	case IndexStatUints16:
		return "uints16"
	case IndexStatUints32:
		return "uints32"
	case IndexStatUints64:
		return "uints64"
	case IndexStatInts8:
		return "ints8"
	case IndexStatInts16:
		return "ints16"
	case IndexStatInts32:
		return "ints32"
	case IndexStatInts64:
		return "ints64"
	default:
		return ""
	}
}

// key index clickhouse string
func (a StatIndexType) ClickhouseName() string {
	switch a {
	case IndexStatID:
		return "id"
	case IndexStatUID:
		return "uid"
	case IndexStatType:
		return "type"
	case IndexStatName:
		return "name"
	case IndexStatSource:
		return "source"
	case IndexStatIP:
		return "ip"
	case IndexStatImage:
		return "image"
	case IndexStatStack:
		return "stack"
	case IndexStatUints:
		return "uints"
	case IndexStatUints8:
		return "uints8"
	case IndexStatUints16:
		return "uints16"
	case IndexStatUints32:
		return "uints32"
	case IndexStatUints64:
		return "uints64"
	case IndexStatInts8:
		return "ints8"
	case IndexStatInts16:
		return "ints16"
	case IndexStatInts32:
		return "ints32"
	case IndexStatInts64:
		return "ints64"
	default:
		return ""
	}
}

// key index type
func (a StatIndexType) Type() string {
	switch a {
	case IndexStatID:
		return "int"
	case IndexStatUID:
		return "int"
	case IndexStatType:
		return "int"
	case IndexStatName:
		return "string"
	case IndexStatSource:
		return "string"
	case IndexStatIP:
		return "string"
	case IndexStatImage:
		return "bool"
	case IndexStatStack:
		return "[]string"
	case IndexStatUints:
		return "[]uint"
	case IndexStatUints8:
		return "[]uint8"
	case IndexStatUints16:
		return "[]uint16"
	case IndexStatUints32:
		return "[]uint32"
	case IndexStatUints64:
		return "[]uint64"
	case IndexStatInts8:
		return "[]int8"
	case IndexStatInts16:
		return "[]int16"
	case IndexStatInts32:
		return "[]int32"
	case IndexStatInts64:
		return "[]int64"
	default:
		return ""
	}
}

// custom title
func (a StatIndexType) Title() string {
	switch a {
	case IndexStatID:
		return "ID"
	case IndexStatUID:
		return "UID"
	case IndexStatType:
		return "Type"
	case IndexStatName:
		return "Name"
	case IndexStatSource:
		return "Source"
	case IndexStatIP:
		return "IP"
	case IndexStatImage:
		return "Image"
	case IndexStatStack:
		return "Stack"
	case IndexStatUints:
		return "Uints"
	case IndexStatUints8:
		return "Uints8"
	case IndexStatUints16:
		return "Uints16"
	case IndexStatUints32:
		return "Uints32"
	case IndexStatUints64:
		return "Uints64"
	case IndexStatInts8:
		return "Ints8"
	case IndexStatInts16:
		return "Ints16"
	case IndexStatInts32:
		return "Ints32"
	case IndexStatInts64:
		return "Ints64"
	default:
		return ""
	}
}

// custom desc
func (a StatIndexType) Desc() string {
	switch a {
	default:
		return ""
	}
}

// struct key to index
func StatKeyIndex(key string) StatIndexType {
	switch key {
	case "id":
		return IndexStatID
	case "uid":
		return IndexStatUID
	case "type":
		return IndexStatType
	case "name":
		return IndexStatName
	case "source":
		return IndexStatSource
	case "ip":
		return IndexStatIP
	case "image":
		return IndexStatImage
	case "stack":
		return IndexStatStack
	case "uints":
		return IndexStatUints
	case "uints8":
		return IndexStatUints8
	case "uints16":
		return IndexStatUints16
	case "uints32":
		return IndexStatUints32
	case "uints64":
		return IndexStatUints64
	case "ints8":
		return IndexStatInts8
	case "ints16":
		return IndexStatInts16
	case "ints32":
		return IndexStatInts32
	case "ints64":
		return IndexStatInts64
	default:
		return 0
	}
}

// valid struct key check
func StatValidKey(key string) bool {
	switch key {
	case "id", "uid", "type", "name", "source", "ip", "image", "stack", "uints", "uints8", "uints16", "uints32", "uints64", "ints8", "ints16", "ints32", "ints64":
		return true
	default:
		return false
	}
}

// struct to map
func (a *Stat) Map() map[string]any {
	return map[string]any{
		"id":      a.ID,
		"uid":     a.UID,
		"type":    a.Type,
		"name":    a.Name,
		"source":  a.Source,
		"ip":      a.IP,
		"image":   a.Image,
		"stack":   a.Stack,
		"uints":   a.Uints,
		"uints8":  a.Uints8,
		"uints16": a.Uints16,
		"uints32": a.Uints32,
		"uints64": a.Uints64,
		"ints8":   a.Ints8,
		"ints16":  a.Ints16,
		"ints32":  a.Ints32,
		"ints64":  a.Ints64,
	}
}

// struct to map
func (a *Stat) Iterate(f func(k StatIndexType, v any)) {
	for _, x := range StatIndexes() {
		f(x, a.Get(x.String()))
	}
}

// gotiny marshal
func (a *Stat) MarshalGotiny() []byte {
	return gotiny.Marshal(&a.ID, &a.UID, &a.Type, &a.Name, &a.Source, &a.IP, &a.Image, &a.Stack, &a.Uints, &a.Uints8, &a.Uints16, &a.Uints32, &a.Uints64, &a.Ints8, &a.Ints16, &a.Ints32, &a.Ints64)
}

// parse gotiny
func UnmarshalStatGotiny(v []byte) (a Stat) {
	gotiny.Unmarshal(v, &a.ID, &a.UID, &a.Type, &a.Name, &a.Source, &a.IP, &a.Image, &a.Stack, &a.Uints, &a.Uints8, &a.Uints16, &a.Uints32, &a.Uints64, &a.Ints8, &a.Ints16, &a.Ints32, &a.Ints64)
	return
}

// fast json marshal
func (a *Stat) Pack() []byte {
	var jsoner = jsoniter.ConfigCompatibleWithStandardLibrary
	b, _ := jsoner.Marshal(a)
	return b
}

// fast json unmarshal
func ParseStat(v []byte) (a Stat, err error) {
	var jsoner = jsoniter.ConfigCompatibleWithStandardLibrary
	err = jsoner.Unmarshal(v, &a)
	return
}

// Parse []any to json
func ParseTupleToStatJson(r []any) (a StatJson) {
	for pos, x := range r {
		switch StatIndexType(pos) {
		case IndexStatID:
			a.Set(FieldStatID, x)
		case IndexStatUID:
			a.Set(FieldStatUID, x)
		case IndexStatType:
			a.Set(FieldStatType, x)
		case IndexStatName:
			a.Set(FieldStatName, x)
		case IndexStatSource:
			a.Set(FieldStatSource, x)
		case IndexStatIP:
			a.Set(FieldStatIP, x)
		case IndexStatImage:
			a.Set(FieldStatImage, x)
		case IndexStatStack:
			a.Set(FieldStatStack, x)
		case IndexStatUints:
			a.Set(FieldStatUints, x)
		case IndexStatUints8:
			a.Set(FieldStatUints8, x)
		case IndexStatUints16:
			a.Set(FieldStatUints16, x)
		case IndexStatUints32:
			a.Set(FieldStatUints32, x)
		case IndexStatUints64:
			a.Set(FieldStatUints64, x)
		case IndexStatInts8:
			a.Set(FieldStatInts8, x)
		case IndexStatInts16:
			a.Set(FieldStatInts16, x)
		case IndexStatInts32:
			a.Set(FieldStatInts32, x)
		case IndexStatInts64:
			a.Set(FieldStatInts64, x)
		}
	}
	return
}

//NewStatJson create struct
func NewStatJson() StatJson {
	return []byte("{}")
}

//StatJson is a struct
type StatJson []byte

//Set value
func (a *StatJson) Set(k string, v any) *StatJson {
	(*a) = jsons.Set((*a), k, v)
	return a
}

//Get value
func (a *StatJson) Get(k string) jsons.Result {
	return jsons.Get((*a), k)
}

//Get value
func (a *StatJson) DeleteFields(fields ...string) {
	(*a) = jsons.Delete((*a), fields...)
}

//ID set or get value
func (a *StatJson) ID(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldStatID)
	}
	a.Set(FieldStatID, v[0])
	return
}

//UID set or get value
func (a *StatJson) UID(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldStatUID)
	}
	a.Set(FieldStatUID, v[0])
	return
}

//Type set or get value
func (a *StatJson) Type(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldStatType)
	}
	a.Set(FieldStatType, v[0])
	return
}

//Name set or get value
func (a *StatJson) Name(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldStatName)
	}
	a.Set(FieldStatName, v[0])
	return
}

//Source set or get value
func (a *StatJson) Source(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldStatSource)
	}
	a.Set(FieldStatSource, v[0])
	return
}

//IP set or get value
func (a *StatJson) IP(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldStatIP)
	}
	a.Set(FieldStatIP, v[0])
	return
}

//Image set or get value
func (a *StatJson) Image(v ...bool) (res bool) {
	if v == nil {
		return jsons.Bool((*a), FieldStatImage)
	}
	a.Set(FieldStatImage, v[0])
	return
}

//Stack set or get value
func (a *StatJson) Stack(v ...string) (res []string) {
	if v == nil {
		return jsons.ArrayString((*a), FieldStatStack)
	}
	a.Set(FieldStatStack, v)
	return
}

//StackAdd add values
func (a *StatJson) StackAdd(v ...string) {
	a.Set(FieldStatStack, append(a.Stack(), v...))
}

//StackPrepend add values
func (a *StatJson) StackPrepend(v ...string) {
	a.Set(FieldStatStack, append(v, a.Stack()...))
}

//StackAddUnique add unique values only
func (a *StatJson) StackAddUnique(v ...string) {
	var list []string
	un := map[string]bool{}
	for _, x := range a.Stack() {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}
	for _, x := range v {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}

	a.Stack(list...)
}

//StackDelete add unique values only
func (a *StatJson) StackDelete(v string) {
	var list []string
	for _, x := range a.Stack() {
		if x == v {
			continue
		}
		list = append(list, x)
	}
	a.Stack(list...)
}

//Uints set or get value
func (a *StatJson) Uints(v ...uint) (res []uint) {
	if v == nil {
		return jsons.ArrayUint((*a), FieldStatUints)
	}
	a.Set(FieldStatUints, v)
	return
}

//UintsAdd add values
func (a *StatJson) UintsAdd(v ...uint) {
	a.Set(FieldStatUints, append(a.Uints(), v...))
}

//UintsPrepend add values
func (a *StatJson) UintsPrepend(v ...uint) {
	a.Set(FieldStatUints, append(v, a.Uints()...))
}

//UintsAddUnique add unique values only
func (a *StatJson) UintsAddUnique(v ...uint) {
	var list []uint
	un := map[uint]bool{}
	for _, x := range a.Uints() {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}
	for _, x := range v {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}

	a.Uints(list...)
}

//UintsDelete add unique values only
func (a *StatJson) UintsDelete(v uint) {
	var list []uint
	for _, x := range a.Uints() {
		if x == v {
			continue
		}
		list = append(list, x)
	}
	a.Uints(list...)
}

//Uints8 set or get value
func (a *StatJson) Uints8(v ...uint8) (res []uint8) {
	if v == nil {
		return jsons.ArrayUint8((*a), FieldStatUints8)
	}
	a.Set(FieldStatUints8, v)
	return
}

//Uints8Add add values
func (a *StatJson) Uints8Add(v ...uint8) {
	a.Set(FieldStatUints8, append(a.Uints8(), v...))
}

//Uints8Prepend add values
func (a *StatJson) Uints8Prepend(v ...uint8) {
	a.Set(FieldStatUints8, append(v, a.Uints8()...))
}

//Uints8AddUnique add unique values only
func (a *StatJson) Uints8AddUnique(v ...uint8) {
	var list []uint8
	un := map[uint8]bool{}
	for _, x := range a.Uints8() {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}
	for _, x := range v {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}

	a.Uints8(list...)
}

//Uints8Delete add unique values only
func (a *StatJson) Uints8Delete(v uint8) {
	var list []uint8
	for _, x := range a.Uints8() {
		if x == v {
			continue
		}
		list = append(list, x)
	}
	a.Uints8(list...)
}

//Uints16 set or get value
func (a *StatJson) Uints16(v ...uint16) (res []uint16) {
	if v == nil {
		return jsons.ArrayUint16((*a), FieldStatUints16)
	}
	a.Set(FieldStatUints16, v)
	return
}

//Uints16Add add values
func (a *StatJson) Uints16Add(v ...uint16) {
	a.Set(FieldStatUints16, append(a.Uints16(), v...))
}

//Uints16Prepend add values
func (a *StatJson) Uints16Prepend(v ...uint16) {
	a.Set(FieldStatUints16, append(v, a.Uints16()...))
}

//Uints16AddUnique add unique values only
func (a *StatJson) Uints16AddUnique(v ...uint16) {
	var list []uint16
	un := map[uint16]bool{}
	for _, x := range a.Uints16() {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}
	for _, x := range v {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}

	a.Uints16(list...)
}

//Uints16Delete add unique values only
func (a *StatJson) Uints16Delete(v uint16) {
	var list []uint16
	for _, x := range a.Uints16() {
		if x == v {
			continue
		}
		list = append(list, x)
	}
	a.Uints16(list...)
}

//Uints32 set or get value
func (a *StatJson) Uints32(v ...uint32) (res []uint32) {
	if v == nil {
		return jsons.ArrayUint32((*a), FieldStatUints32)
	}
	a.Set(FieldStatUints32, v)
	return
}

//Uints32Add add values
func (a *StatJson) Uints32Add(v ...uint32) {
	a.Set(FieldStatUints32, append(a.Uints32(), v...))
}

//Uints32Prepend add values
func (a *StatJson) Uints32Prepend(v ...uint32) {
	a.Set(FieldStatUints32, append(v, a.Uints32()...))
}

//Uints32AddUnique add unique values only
func (a *StatJson) Uints32AddUnique(v ...uint32) {
	var list []uint32
	un := map[uint32]bool{}
	for _, x := range a.Uints32() {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}
	for _, x := range v {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}

	a.Uints32(list...)
}

//Uints32Delete add unique values only
func (a *StatJson) Uints32Delete(v uint32) {
	var list []uint32
	for _, x := range a.Uints32() {
		if x == v {
			continue
		}
		list = append(list, x)
	}
	a.Uints32(list...)
}

//Uints64 set or get value
func (a *StatJson) Uints64(v ...uint64) (res []uint64) {
	if v == nil {
		return jsons.ArrayUint64((*a), FieldStatUints64)
	}
	a.Set(FieldStatUints64, v)
	return
}

//Uints64Add add values
func (a *StatJson) Uints64Add(v ...uint64) {
	a.Set(FieldStatUints64, append(a.Uints64(), v...))
}

//Uints64Prepend add values
func (a *StatJson) Uints64Prepend(v ...uint64) {
	a.Set(FieldStatUints64, append(v, a.Uints64()...))
}

//Uints64AddUnique add unique values only
func (a *StatJson) Uints64AddUnique(v ...uint64) {
	var list []uint64
	un := map[uint64]bool{}
	for _, x := range a.Uints64() {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}
	for _, x := range v {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}

	a.Uints64(list...)
}

//Uints64Delete add unique values only
func (a *StatJson) Uints64Delete(v uint64) {
	var list []uint64
	for _, x := range a.Uints64() {
		if x == v {
			continue
		}
		list = append(list, x)
	}
	a.Uints64(list...)
}

//Ints8 set or get value
func (a *StatJson) Ints8(v ...int8) (res []int8) {
	if v == nil {
		return jsons.ArrayInt8((*a), FieldStatInts8)
	}
	a.Set(FieldStatInts8, v)
	return
}

//Ints8Add add values
func (a *StatJson) Ints8Add(v ...int8) {
	a.Set(FieldStatInts8, append(a.Ints8(), v...))
}

//Ints8Prepend add values
func (a *StatJson) Ints8Prepend(v ...int8) {
	a.Set(FieldStatInts8, append(v, a.Ints8()...))
}

//Ints8AddUnique add unique values only
func (a *StatJson) Ints8AddUnique(v ...int8) {
	var list []int8
	un := map[int8]bool{}
	for _, x := range a.Ints8() {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}
	for _, x := range v {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}

	a.Ints8(list...)
}

//Ints8Delete add unique values only
func (a *StatJson) Ints8Delete(v int8) {
	var list []int8
	for _, x := range a.Ints8() {
		if x == v {
			continue
		}
		list = append(list, x)
	}
	a.Ints8(list...)
}

//Ints16 set or get value
func (a *StatJson) Ints16(v ...int16) (res []int16) {
	if v == nil {
		return jsons.ArrayInt16((*a), FieldStatInts16)
	}
	a.Set(FieldStatInts16, v)
	return
}

//Ints16Add add values
func (a *StatJson) Ints16Add(v ...int16) {
	a.Set(FieldStatInts16, append(a.Ints16(), v...))
}

//Ints16Prepend add values
func (a *StatJson) Ints16Prepend(v ...int16) {
	a.Set(FieldStatInts16, append(v, a.Ints16()...))
}

//Ints16AddUnique add unique values only
func (a *StatJson) Ints16AddUnique(v ...int16) {
	var list []int16
	un := map[int16]bool{}
	for _, x := range a.Ints16() {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}
	for _, x := range v {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}

	a.Ints16(list...)
}

//Ints16Delete add unique values only
func (a *StatJson) Ints16Delete(v int16) {
	var list []int16
	for _, x := range a.Ints16() {
		if x == v {
			continue
		}
		list = append(list, x)
	}
	a.Ints16(list...)
}

//Ints32 set or get value
func (a *StatJson) Ints32(v ...int32) (res []int32) {
	if v == nil {
		return jsons.ArrayInt32((*a), FieldStatInts32)
	}
	a.Set(FieldStatInts32, v)
	return
}

//Ints32Add add values
func (a *StatJson) Ints32Add(v ...int32) {
	a.Set(FieldStatInts32, append(a.Ints32(), v...))
}

//Ints32Prepend add values
func (a *StatJson) Ints32Prepend(v ...int32) {
	a.Set(FieldStatInts32, append(v, a.Ints32()...))
}

//Ints32AddUnique add unique values only
func (a *StatJson) Ints32AddUnique(v ...int32) {
	var list []int32
	un := map[int32]bool{}
	for _, x := range a.Ints32() {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}
	for _, x := range v {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}

	a.Ints32(list...)
}

//Ints32Delete add unique values only
func (a *StatJson) Ints32Delete(v int32) {
	var list []int32
	for _, x := range a.Ints32() {
		if x == v {
			continue
		}
		list = append(list, x)
	}
	a.Ints32(list...)
}

//Ints64 set or get value
func (a *StatJson) Ints64(v ...int64) (res []int64) {
	if v == nil {
		return jsons.ArrayInt64((*a), FieldStatInts64)
	}
	a.Set(FieldStatInts64, v)
	return
}

//Ints64Add add values
func (a *StatJson) Ints64Add(v ...int64) {
	a.Set(FieldStatInts64, append(a.Ints64(), v...))
}

//Ints64Prepend add values
func (a *StatJson) Ints64Prepend(v ...int64) {
	a.Set(FieldStatInts64, append(v, a.Ints64()...))
}

//Ints64AddUnique add unique values only
func (a *StatJson) Ints64AddUnique(v ...int64) {
	var list []int64
	un := map[int64]bool{}
	for _, x := range a.Ints64() {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}
	for _, x := range v {
		if un[x] {
			continue
		}
		list = append(list, x)
		un[x] = true
	}

	a.Ints64(list...)
}

//Ints64Delete add unique values only
func (a *StatJson) Ints64Delete(v int64) {
	var list []int64
	for _, x := range a.Ints64() {
		if x == v {
			continue
		}
		list = append(list, x)
	}
	a.Ints64(list...)
}

// sql StatSQL class
type StatSQL struct{ pool *pgxpool.Pool }

func NewStatSQL(pool *pgxpool.Pool) (a *StatSQL) {
	a = new(StatSQL)
	a.pool = pool
	a.CreateTable()
	return
}

//delete item

func (a *StatSQL) Conn(f func(conn *pgxpool.Conn) (err error)) (err error) {

	conn, err := a.pool.Acquire(context.Background())
	if err != nil {
		return
	}
	defer conn.Release()
	return f(conn)
}

// parse sql query
func (a *StatSQL) TableName() (res string) {
	return "stat"
}

func (a *StatSQL) Count(where ...string) (count int, err error) {

	var q string

	switch len(where) {
	case 0:
		q = "select count(*) from stat"
	default:
		q = "select count(*) from stat where " + strings.Join(where, " ")
	}

	err = a.Conn(func(conn *pgxpool.Conn) (err error) {
		return conn.QueryRow(context.Background(), q).Scan(&count)
	})
	return
}

// Insert struct and return int id
func (a *StatSQL) Insert(v *Stat) (id int, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "insert into stat (uid, type, name, source, ip, image, stack, uints, uints8, uints16, uints32, uints64, ints8, ints16, ints32, ints64) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) returning id"
	err = conn.QueryRow(c, q, v.sqlTuple()...).Scan(&id)
	return
}

// Insert struct and return int id
func (a *StatSQL) InsertNoConflict(v *Stat) (id int, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "insert into stat (uid, type, name, source, ip, image, stack, uints, uints8, uints16, uints32, uints64, ints8, ints16, ints32, ints64) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) on conflict do nothing returning id"
	err = conn.QueryRow(c, q, v.sqlTuple()...).Scan(&id)
	return
}

// Insert full struct ID must be (ignore all skips)
func (a *StatSQL) InsertFull(v *Stat) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "insert into stat (id, uid, type, name, source, ip, image, stack, uints, uints8, uints16, uints32, uints64, ints8, ints16, ints32, ints64) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)"
	_, err = conn.Exec(c, q, v.sqlAllTuples()...)
	return
}

// update serial counter (count all items, add 1 and plus custom int)
// useful if you insert with ID
func (a *StatSQL) ReindexIDSerialCounter(add ...int) (err error) {
	return a.Conn(func(conn *pgxpool.Conn) (err error) {
		plus := 1
		if len(add) > 0 {
			plus = add[0]
		}
		_, err = conn.Exec(context.Background(), "SELECT setval(pg_get_serial_sequence('stat', 'id'), COALESCE((SELECT MAX(id) FROM stat), 0) + $1, false)", plus)
		return
	})
}

// set serial counter
func (a *StatSQL) SetIDSerialCounter(value int) (err error) {
	return a.Conn(func(conn *pgxpool.Conn) (err error) {
		_, err = conn.Exec(context.Background(), "SELECT setval(pg_get_serial_sequence('stat', 'id'), $1, false)", value)
		return
	})
}

// parse sql query
func (a *StatSQL) Get(id any, fields ...StatIndexType) (res *Stat, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if fields == nil {
		fields = a.Fields()
	}

	var list []string
	for _, x := range fields {
		list = append(list, x.SQLName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s ", fieldlist)
	q = q + "from stat where id = $1 limit 1"
	res = new(Stat)
	rows, err := conn.Query(c, q, id)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("not found")
		return
	}
	v, err := rows.Values()
	if err != nil {
		return
	}
	if len(v) != len(fields) {
		err = errors.New("len")
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), v[pos])
	}
	err = rows.Err()
	return
}

// parse sql query
func (a *StatSQL) GetWhere(where string, fields ...StatIndexType) (res *Stat, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if fields == nil {
		fields = a.Fields()
	}

	var list []string
	for _, x := range fields {
		list = append(list, x.SQLName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s", fieldlist)
	q = q + " from stat"
	q = q + fmt.Sprintf(" where %s ", where)
	q = q + " limit 1"
	res = new(Stat)
	rows, err := conn.Query(c, q)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("not found")
		return
	}
	v, err := rows.Values()
	if err != nil {
		return
	}
	if len(v) != len(fields) {
		err = errors.New("len")
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), v[pos])
	}
	err = rows.Err()
	return
}

// parse sql query
func (a *StatSQL) Row(eq map[string]any, fields ...StatIndexType) (res *Stat, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if fields == nil {
		fields = a.Fields()
	}

	var list []string
	for _, x := range fields {
		list = append(list, x.SQLName())
	}
	fieldlist := strings.Join(list, ", ")

	var andlist []string
	var values []any
	var count int
	for k, v := range eq {
		if !StatValidKey(k) {
			err = errors.New(k)
			return
		}
		in := StatKeyIndex(k)
		count++
		andlist = append(andlist, fmt.Sprintf("%s = $%d", in.SQLName(), count))
		values = append(values, v)
	}
	if len(andlist) == 0 {
		err = errors.New("emptykeys")
		return
	}
	keys := strings.Join(andlist, " and ")

	q := fmt.Sprintf("select %s from %s where %s limit 1", fieldlist, a.TableName(), keys)
	res = new(Stat)
	rows, err := conn.Query(c, q, values...)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("not found")
		return
	}
	v, err := rows.Values()
	if err != nil {
		return
	}
	if len(v) != len(fields) {
		err = errors.New("len")
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), v[pos])
	}
	err = rows.Err()
	return
}

// get simple list from db
func (a *StatSQL) All(fields ...StatIndexType) (res []*Stat, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if fields == nil {
		fields = a.Fields()
	}

	var list []string
	for _, x := range fields {
		list = append(list, x.SQLName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s", fieldlist)
	q += " from stat"
	rows, err := conn.Query(c, q)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item Stat
		v, err := rows.Values()
		if err != nil {
			continue
		}
		if len(v) != len(fields) {
			continue
		}
		for pos, x := range fields {
			item.Update(x.String(), v[pos])
		}
		res = append(res, &item)
	}
	err = rows.Err()

	return
}

// get simple list from db
func (a *StatSQL) List(limit, offset int, fields ...StatIndexType) (res []*Stat, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if fields == nil {
		fields = a.Fields()
	}

	var list []string
	for _, x := range fields {
		list = append(list, x.SQLName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s", fieldlist)
	q += " from stat limit $1 offset $2"
	rows, err := conn.Query(c, q, limit, offset)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item Stat
		v, err := rows.Values()
		if err != nil {
			continue
		}
		if len(v) != len(fields) {
			continue
		}
		for pos, x := range fields {
			item.Update(x.String(), v[pos])
		}
		res = append(res, &item)
	}
	err = rows.Err()

	return
}

// update sql query
func (a *StatSQL) Update(id any, k string, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if !StatValidKey(k) {
		return errors.New("invalid key")
	}
	in := StatKeyIndex(k)
	q := fmt.Sprintf("update stat set %s = $1 where id = $2", in.SQLName())
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// update sql query
func (a *StatSQL) UpdateWhere(k string, v any, where string) (err error) {
	if !StatValidKey(k) {
		return errors.New("invalid key")
	}
	in := StatKeyIndex(k)
	return a.Conn(func(conn *pgxpool.Conn) (err error) {
		q := fmt.Sprintf("update %s set %s = $1 where %s", a.TableName(), in.SQLName(), where)
		_, err = conn.Exec(context.Background(), q, v)
		return
	})
}

// update sql query
func (a *StatSQL) Updates(id any, keys map[string]any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if keys == nil {
		return errors.New("emptykeys")
	}

	var fields []string
	var values []any
	var count int
	for k, v := range keys {
		if !StatValidKey(k) {
			return errors.New(k)
		}
		in := StatKeyIndex(k)
		count++
		fields = append(fields, fmt.Sprintf("%s = $%d", in.SQLName(), count))
		values = append(values, v)
	}

	list := strings.Join(fields, ", ")
	count++
	values = append(values, id)

	q := fmt.Sprintf("update stat set %s where id = $%d", list, count)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, values...)
	return
}

// update sql query
func (a *StatSQL) UpdatesWhere(keys map[string]any, where string, args ...any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if keys == nil {
		return errors.New("emptykeys")
	}

	var fields []string
	var values []any
	var count int
	for k, v := range keys {
		if !StatValidKey(k) {
			return errors.New(k)
		}
		in := StatKeyIndex(k)
		count++
		fields = append(fields, fmt.Sprintf("%s = $%d", in.SQLName(), count))
		values = append(values, v)
	}

	list := strings.Join(fields, ", ")
	count++

	where = fmt.Sprintf(where, args...)
	q := fmt.Sprintf("update stat set %s where %s", list, where)
	_, err = conn.Exec(c, q, values...)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddStack(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set stack = stack || '%s'::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteStack(id any, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update stat set stack = stack - $1 where id = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddStackWhere(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set stack = stack || '%s'::jsonb where %s", res, where)
	_, err = conn.Exec(c, q)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteStackWhere(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update stat set stack = stack - $1 where %s", where)
	_, err = conn.Exec(c, q, v)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddUints(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set uints = uints || '%s'::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteUints(id any, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update stat set uints = uints - $1 where id = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddUintsWhere(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set uints = uints || '%s'::jsonb where %s", res, where)
	_, err = conn.Exec(c, q)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteUintsWhere(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update stat set uints = uints - $1 where %s", where)
	_, err = conn.Exec(c, q, v)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddUints16(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set uints16 = uints16 || '%s'::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteUints16(id any, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update stat set uints16 = uints16 - $1 where id = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddUints16Where(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set uints16 = uints16 || '%s'::jsonb where %s", res, where)
	_, err = conn.Exec(c, q)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteUints16Where(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update stat set uints16 = uints16 - $1 where %s", where)
	_, err = conn.Exec(c, q, v)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddUints32(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set uints32 = uints32 || '%s'::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteUints32(id any, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update stat set uints32 = uints32 - $1 where id = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddUints32Where(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set uints32 = uints32 || '%s'::jsonb where %s", res, where)
	_, err = conn.Exec(c, q)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteUints32Where(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update stat set uints32 = uints32 - $1 where %s", where)
	_, err = conn.Exec(c, q, v)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddUints64(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set uints64 = uints64 || '%s'::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteUints64(id any, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update stat set uints64 = uints64 - $1 where id = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddUints64Where(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set uints64 = uints64 || '%s'::jsonb where %s", res, where)
	_, err = conn.Exec(c, q)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteUints64Where(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update stat set uints64 = uints64 - $1 where %s", where)
	_, err = conn.Exec(c, q, v)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddInts16(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set ints16 = ints16 || '%s'::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteInts16(id any, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update stat set ints16 = ints16 - $1 where id = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddInts16Where(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set ints16 = ints16 || '%s'::jsonb where %s", res, where)
	_, err = conn.Exec(c, q)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteInts16Where(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update stat set ints16 = ints16 - $1 where %s", where)
	_, err = conn.Exec(c, q, v)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddInts32(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set ints32 = ints32 || '%s'::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteInts32(id any, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update stat set ints32 = ints32 - $1 where id = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddInts32Where(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set ints32 = ints32 || '%s'::jsonb where %s", res, where)
	_, err = conn.Exec(c, q)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteInts32Where(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update stat set ints32 = ints32 - $1 where %s", where)
	_, err = conn.Exec(c, q, v)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddInts64(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set ints64 = ints64 || '%s'::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteInts64(id any, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update stat set ints64 = ints64 - $1 where id = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// add array value to jsonb array
func (a *StatSQL) AddInts64Where(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update stat set ints64 = ints64 || '%s'::jsonb where %s", res, where)
	_, err = conn.Exec(c, q)
	return
}

// delete array value from jsonb array
func (a *StatSQL) DeleteInts64Where(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update stat set ints64 = ints64 - $1 where %s", where)
	_, err = conn.Exec(c, q, v)
	return
}

type StatQuery struct {
	Limit  int                 `json:"limit,omitempty"`
	Offset int                 `json:"offset,omitempty"`
	Sort   string              `json:"sort,omitempty"`
	Desc   bool                `json:"desc,omitempty"`
	EQ     map[string]any      `json:"eq,omitempty"`   //equal
	GT     map[string]any      `json:"gt,omitempty"`   //greater then...
	LT     map[string]any      `json:"lt,omitempty"`   //less then
	NOT    map[string][]any    `json:"not,omitempty"`  //less then
	Like   map[string]string   `json:"like,omitempty"` //full text search
	Custom string              `json:"-"`              //append unsafe where condition
	Fields []string            `json:"fields,omitempty"`
	IN     map[string][]int    `json:"in,omitempty"`
	INS    map[string][]string `json:"ins,omitempty"`
	INSQL  map[string]string   `json:"insql,omitempty"` //unsafe in (condition)
}

func (a *StatQuery) Render() (sql string, fields []StatIndexType, values []any) {

	switch len(a.Fields) == 0 {
	case true:
		fields = []StatIndexType{IndexStatID, IndexStatUID, IndexStatType, IndexStatName, IndexStatSource, IndexStatIP, IndexStatImage, IndexStatStack, IndexStatUints, IndexStatUints8, IndexStatUints16, IndexStatUints32, IndexStatUints64, IndexStatInts8, IndexStatInts16, IndexStatInts32, IndexStatInts64}
	default:
		for _, x := range a.Fields {
			if !StatValidKey(x) {
				continue
			}
			fields = append(fields, StatKeyIndex(x))
		}
	}

	var fieldsStrings []string
	for _, x := range fields {
		fieldsStrings = append(fieldsStrings, x.SQLName())
	}

	for k, v := range a.EQ {
		if fmt.Sprint(v) == "" {
			delete(a.EQ, k)
		}
	}
	for k, v := range a.GT {
		if fmt.Sprint(v) == "" {
			delete(a.GT, k)
		}
	}
	for k, v := range a.LT {
		if fmt.Sprint(v) == "" {
			delete(a.LT, k)
		}
	}
	for k, v := range a.NOT {
		if v == nil {
			delete(a.NOT, k)
		}
	}
	for k, v := range a.Like {
		if v == "" {
			delete(a.Like, k)
		}
	}

	// sql
	var list []string

	var count int

	// select
	list = append(list, "select")
	list = append(list, strings.Join(fieldsStrings, ", "))

	// from
	list = append(list, "from stat")

	var andlist []string

	// EQ where
	for k, v := range a.EQ {
		if !StatValidKey(k) {
			continue
		}
		p := StatKeyIndex(k)
		sqlName := p.SQLName()
		switch p.Type() {
		case "bool":
			switch cast.Bool(v) {
			case true:
				andlist = append(andlist, sqlName)
			case false:
				andlist = append(andlist, fmt.Sprintf("(not %[1]s or %[1]s is null)", sqlName))
			}

		default:
			count++
			andlist = append(andlist, fmt.Sprintf("%s = $%d", sqlName, count))
			values = append(values, v)
		}
	}

	// GT where
	for k, v := range a.GT {

		if !StatValidKey(k) {
			continue
		}
		sqlName := StatKeyIndex(k).SQLName()
		count++
		andlist = append(andlist, fmt.Sprintf("%s > $%d", sqlName, count))
		values = append(values, v)
	}

	// LT where
	for k, v := range a.LT {
		if !StatValidKey(k) {
			continue
		}
		sqlName := StatKeyIndex(k).SQLName()
		count++
		andlist = append(andlist, fmt.Sprintf("%s < $%d", sqlName, count))
		values = append(values, v)
	}

	// NOT where
	for k, v := range a.NOT {
		if !StatValidKey(k) {
			continue
		}
		sqlName := StatKeyIndex(k).SQLName()

		for _, x := range v {
			count++
			andlist = append(andlist, fmt.Sprintf("%s != $%d", sqlName, count))
			values = append(values, x)
		}

	}

	// LIKE where
	for k, v := range a.Like {
		if !StatValidKey(k) {
			continue
		}
		sqlName := StatKeyIndex(k).SQLName()
		count++
		andlist = append(andlist, fmt.Sprintf("%s ilike $%d", sqlName, count))
		values = append(values, "%"+v+"%")
	}

	// IN where
	if a.IN != nil {
		for k, v := range a.IN {
			if !StatValidKey(k) {
				continue
			}
			sqlName := StatKeyIndex(k).SQLName()
			var inlist []string
			for _, num := range v {
				count++
				inlist = append(inlist, fmt.Sprintf("$%d", count))
				values = append(values, num)
			}
			if len(inlist) > 0 {
				andlist = append(andlist, fmt.Sprintf("%s in (%s)", sqlName, strings.Join(inlist, ",")))
			}
		}
	}

	// INS where
	if a.INS != nil {
		for k, v := range a.INS {
			if !StatValidKey(k) {
				continue
			}
			sqlName := StatKeyIndex(k).SQLName()
			var inlist []string
			for _, str := range v {
				count++
				inlist = append(inlist, fmt.Sprintf("$%d", count))
				values = append(values, str)
			}
			if len(inlist) > 0 {
				andlist = append(andlist, fmt.Sprintf("%s in (%s)", sqlName, strings.Join(inlist, ",")))
			}
		}
	}

	// INSQ where
	if a.INSQL != nil {
		for k, v := range a.INSQL {
			if !StatValidKey(k) {
				continue
			}
			sqlName := StatKeyIndex(k).SQLName()
			count++
			andlist = append(andlist, fmt.Sprintf("%s in (%s)", sqlName, v))
		}
	}

	if a.Custom != "" {
		andlist = append(andlist, a.Custom)
	}

	// render where
	if andlist != nil {
		list = append(list, "where")
		list = append(list, strings.Join(andlist, " and "))
	}

	// sort by
	if a.Sort != "" && StatValidKey(a.Sort) {
		list = append(list, fmt.Sprintf("order by %s", StatKeyIndex(a.Sort).SQLName()))
		if a.Desc {
			list = append(list, "desc")
		}
	}

	// limit, offset
	if a.Limit > 0 {
		list = append(list, fmt.Sprintf("limit %d", a.Limit))
	}
	if a.Offset > 0 {
		list = append(list, fmt.Sprintf("offset %d", a.Offset))
	}

	// render sql
	sql = strings.Join(list, " ")
	return
}

func NewStatQuery() *StatQuery {
	a := new(StatQuery)
	a.EQ = make(map[string]any)
	a.GT = make(map[string]any)
	a.LT = make(map[string]any)
	a.NOT = make(map[string][]any)
	a.Like = make(map[string]string)
	a.IN = make(map[string][]int)
	a.INS = make(map[string][]string)
	a.INSQL = make(map[string]string)
	return a
}

func (a *StatSQL) Search(q *StatQuery) (res []*Stat, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	sql, fields, values := q.Render()

	rows, err := conn.Query(c, sql, values...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item Stat
		v, err := rows.Values()
		if err != nil {
			continue
		}
		if len(v) != len(fields) {
			continue
		}

		for pos, x := range fields {
			item.Update(x.String(), v[pos])
		}
		res = append(res, &item)
	}
	err = rows.Err()

	return
}

// select custom sql query
func (a *StatSQL) Select(where string, fields ...StatIndexType) (res []*Stat, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if fields == nil {
		fields = a.Fields()
	}

	var list []string
	for _, x := range fields {
		list = append(list, x.SQLName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s ", fieldlist)
	q = q + "from stat"
	if where != "" {
		q += " where " + where
	}
	rows, err := conn.Query(c, q)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item Stat
		v, err := rows.Values()
		if err != nil {
			continue
		}
		if len(v) != len(fields) {
			continue
		}

		for pos, x := range fields {
			item.Update(x.String(), v[pos])
		}
		res = append(res, &item)
	}
	err = rows.Err()
	return
}

// delete item
func (a *StatSQL) Delete(id any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	_, err = conn.Exec(c, "delete from stat where id = $1", id)
	return

}

// delete item where: i = 1 and w = 'nice'
func (a *StatSQL) DeleteWhere(where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	_, err = conn.Exec(c, fmt.Sprintf("delete from stat where %s", where))
	return

}

// has value in db
func (a *StatSQL) Has(field StatIndexType, v any) (has bool, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("select exists (select 1 from stat where %s = $1 limit 1)", field.SQLName())
	err = conn.QueryRow(c, q, v).Scan(&has)
	return
}

// Create table
func (a *StatSQL) CreateTable() (err error) {
	return a.Conn(func(conn *pgxpool.Conn) (err error) {
		q := `create table if not exists stat (
	id          bigserial primary key,
	uid         bigint,
	type        bigint,
	name        text,
	source      text,
	ip          text,
	image       boolean,
	stack       jsonb default '[]'::jsonb,
	uints       jsonb default '[]'::jsonb,
	uints8      bytea,
	uints16     jsonb default '[]'::jsonb,
	uints32     jsonb default '[]'::jsonb,
	uints64     jsonb default '[]'::jsonb,
	ints8       bytea,
	ints16      jsonb default '[]'::jsonb,
	ints32      jsonb default '[]'::jsonb,
	ints64      jsonb default '[]'::jsonb
)
`
		_, err = conn.Exec(context.Background(), q)
		return
	})
}

// parse sql query
func (a *StatSQL) Fields() (res []StatIndexType) {
	return []StatIndexType{IndexStatID, IndexStatUID, IndexStatType, IndexStatName, IndexStatSource, IndexStatIP, IndexStatImage, IndexStatStack, IndexStatUints, IndexStatUints8, IndexStatUints16, IndexStatUints32, IndexStatUints64, IndexStatInts8, IndexStatInts16, IndexStatInts32, IndexStatInts64}
} //clickhouse StatCQL class
type StatCQL struct {
	conn driver.Conn
	list []Stat
	sync.Mutex
}

func NewStatCQL(conn driver.Conn) (a *StatCQL) {
	a = new(StatCQL)
	a.conn = conn
	a.CreateTable()
	go a.deamon()
	return
}

// parse clickhouse query
func (a *StatCQL) TableName() (res string) {
	return "stat"
}

func (a *StatCQL) Count(where ...string) (count int, err error) {

	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	var q string
	switch len(where) {
	case 0:
		q = "select count() from stat"
	default:
		q = "select count() from stat where " + strings.Join(where, " ")
	}
	var n uint64
	err = a.conn.QueryRow(context.Background(), q).Scan(&n)
	count = int(n)
	return
}

// Add items to clickhouse queue
func (a *StatCQL) Add(items ...Stat) {
	a.Lock()
	defer a.Unlock()
	a.list = append(a.list, items...)
}

// Push queued items to clickhouse
func (a *StatCQL) Push() (err error) {
	a.Lock()
	defer a.Unlock()
	if len(a.list) == 0 {
		return
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	batch, err := a.conn.PrepareBatch(context.Background(), "insert into stat (id, uid, type, name, source, ip, image, stack, uints, uints8, uints16, uints32, uints64, ints8, ints16, ints32, ints64)")
	if err != nil {
		return
	}
	for _, item := range a.list {
		err = batch.Append(item.clickhouseTuple()...)
		if err != nil {
			_ = batch.Abort()
			return
		}
	}
	err = batch.Send()
	if err != nil {
		_ = batch.Abort()
		return
	}
	a.list = a.list[:0]
	return
}

// deamon pushes queued items every 10 minutes
func (a *StatCQL) deamon() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		a.Lock()
		hasItems := len(a.list) > 0
		a.Unlock()
		if hasItems {
			_ = a.Push()
		}
	}
}

// parse clickhouse query
func (a *StatCQL) Get(id any, fields ...StatIndexType) (res *Stat, err error) {

	if fields == nil {
		fields = a.Fields()
	}
	var list []string
	for _, x := range fields {
		list = append(list, x.ClickhouseName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s from stat where uid = ? limit 1", fieldlist)
	res = new(Stat)

	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	rows, err := a.conn.Query(context.Background(), q, id)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("not found")
		return
	}
	rowValues := make([]any, len(fields))
	scan := make([]any, len(fields))
	for i := range rowValues {
		scan[i] = &rowValues[i]
	}
	err = rows.Scan(scan...)
	if err != nil {
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), rowValues[pos])
	}
	err = rows.Err()

	return
}

// parse clickhouse query
func (a *StatCQL) GetWhere(where string, fields ...StatIndexType) (res *Stat, err error) {

	if fields == nil {
		fields = a.Fields()
	}
	var list []string
	for _, x := range fields {
		list = append(list, x.ClickhouseName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s from stat where %s limit 1", fieldlist, where)
	res = new(Stat)

	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	rows, err := a.conn.Query(context.Background(), q)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("not found")
		return
	}
	rowValues := make([]any, len(fields))
	scan := make([]any, len(fields))
	for i := range rowValues {
		scan[i] = &rowValues[i]
	}
	err = rows.Scan(scan...)
	if err != nil {
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), rowValues[pos])
	}
	err = rows.Err()

	return
}

// update clickhouse query
func (a *StatCQL) Update(id any, k string, v any, where ...string) (err error) {
	if !StatValidKey(k) {
		return errors.New("invalid key")
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	in := StatKeyIndex(k)
	q := fmt.Sprintf("alter table stat update %s = ? where uid = ?", in.ClickhouseName())
	if len(where) > 0 {
		q += " and " + where[0]
	}
	return a.conn.Exec(context.Background(), q, v, id)
}

// update clickhouse query
func (a *StatCQL) UpdateWhere(k string, v any, where string) (err error) {
	if !StatValidKey(k) {
		return errors.New("invalid key")
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	in := StatKeyIndex(k)
	q := fmt.Sprintf("alter table %s update %s = ? where %s", a.TableName(), in.ClickhouseName(), where)
	return a.conn.Exec(context.Background(), q, v)
}

// update clickhouse query
func (a *StatCQL) Updates(id any, keys map[string]any, where ...string) (err error) {

	if keys == nil {
		return errors.New("emptykeys")
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	var fields []string
	var values []any
	for k, v := range keys {
		if !StatValidKey(k) {
			return errors.New(k)
		}
		in := StatKeyIndex(k)
		fields = append(fields, fmt.Sprintf("%s = ?", in.ClickhouseName()))
		values = append(values, v)
	}
	list := strings.Join(fields, ", ")
	values = append(values, id)

	q := fmt.Sprintf("alter table stat update %s where uid = ?", list)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	return a.conn.Exec(context.Background(), q, values...)
}

// update clickhouse query
func (a *StatCQL) UpdatesWhere(keys map[string]any, where string, args ...any) (err error) {

	if keys == nil {
		return errors.New("emptykeys")
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	var fields []string
	var values []any
	for k, v := range keys {
		if !StatValidKey(k) {
			return errors.New(k)
		}
		in := StatKeyIndex(k)
		fields = append(fields, fmt.Sprintf("%s = ?", in.ClickhouseName()))
		values = append(values, v)
	}
	list := strings.Join(fields, ", ")

	where = fmt.Sprintf(where, args...)
	q := fmt.Sprintf("alter table stat update %s where %s", list, where)
	return a.conn.Exec(context.Background(), q, values...)
}

// select custom clickhouse query
func (a *StatCQL) Select(where string, fields ...StatIndexType) (res []*Stat, err error) {

	if fields == nil {
		fields = a.Fields()
	}
	var list []string
	for _, x := range fields {
		list = append(list, x.ClickhouseName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s from stat", fieldlist)
	if where != "" {
		q += " where " + where
	}

	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	rows, err := a.conn.Query(context.Background(), q)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var item Stat
		values := make([]any, len(fields))
		scan := make([]any, len(fields))
		for i := range values {
			scan[i] = &values[i]
		}
		err = rows.Scan(scan...)
		if err != nil {
			continue
		}
		for pos, x := range fields {
			item.Update(x.String(), values[pos])
		}
		res = append(res, &item)
	}
	err = rows.Err()

	return
}

// delete clickhouse item
func (a *StatCQL) Delete(id any) (err error) {
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	q := "alter table stat delete where uid = ?"
	return a.conn.Exec(context.Background(), q, id)
}

// delete clickhouse items
func (a *StatCQL) DeleteWhere(where string) (err error) {
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	q := "alter table stat delete where " + where
	return a.conn.Exec(context.Background(), q)
}

// has value in clickhouse
func (a *StatCQL) Has(field StatIndexType, v any) (has bool, err error) {
	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	q := fmt.Sprintf("select count() from stat where %s = ? limit 1", field.ClickhouseName())
	var n uint64
	err = a.conn.QueryRow(context.Background(), q, v).Scan(&n)
	has = n > 0
	return
}

// stat by clickhouse field
func (a *StatCQL) Stat(field StatIndexType) (res map[string]int, err error) {
	res = make(map[string]int)
	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	name := field.ClickhouseName()
	q := fmt.Sprintf("select toString(%[1]s), count() from %[2]s group by %[1]s", name, a.TableName())
	rows, err := a.conn.Query(context.Background(), q)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var k string
		var n uint64
		err = rows.Scan(&k, &n)
		if err != nil {
			return
		}
		res[k] = int(n)
	}
	err = rows.Err()
	return
}

// Create clickhouse table
func (a *StatCQL) CreateTable() (err error) {
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	q := `create table if not exists stat (
	id          Int64,
	uid         Int64,
	type        LowCardinality(Int64),
	name        String,
	source      String,
	ip          IPv4,
	image       Bool,
	stack       Array(String),
	uints       Array(UInt64),
	uints8      Array(UInt8),
	uints16     Array(UInt16),
	uints32     Array(UInt32),
	uints64     Array(UInt64),
	ints8       Array(Int8),
	ints16      Array(Int16),
	ints32      Array(Int32),
	ints64      Array(Int64)
)
engine = MergeTree()
order by (uid)
primary key (uid)`
	return a.conn.Exec(context.Background(), q)
}

// parse clickhouse query
func (a *StatCQL) Fields() (res []StatIndexType) {
	return []StatIndexType{IndexStatID, IndexStatUID, IndexStatType, IndexStatName, IndexStatSource, IndexStatIP, IndexStatImage, IndexStatStack, IndexStatUints, IndexStatUints8, IndexStatUints16, IndexStatUints32, IndexStatUints64, IndexStatInts8, IndexStatInts16, IndexStatInts32, IndexStatInts64}
}

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StatProto struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Uid           int64                  `protobuf:"varint,2,opt,name=uid,proto3" json:"uid,omitempty"`
	Type          int64                  `protobuf:"varint,3,opt,name=type,proto3" json:"type,omitempty"`
	Name          string                 `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Source        string                 `protobuf:"bytes,5,opt,name=source,proto3" json:"source,omitempty"`
	Ip            string                 `protobuf:"bytes,6,opt,name=ip,proto3" json:"ip,omitempty"`
	Image         bool                   `protobuf:"varint,7,opt,name=image,proto3" json:"image,omitempty"`
	Stack         []string               `protobuf:"bytes,8,rep,name=stack,proto3" json:"stack,omitempty"`
	Uints         []uint64               `protobuf:"varint,9,rep,packed,name=uints,proto3" json:"uints,omitempty"`
	Uints8        []byte                 `protobuf:"bytes,10,opt,name=uints8,proto3" json:"uints8,omitempty"`
	Uints16       []uint32               `protobuf:"varint,11,rep,packed,name=uints16,proto3" json:"uints16,omitempty"`
	Uints32       []uint32               `protobuf:"varint,12,rep,packed,name=uints32,proto3" json:"uints32,omitempty"`
	Uints64       []uint64               `protobuf:"varint,13,rep,packed,name=uints64,proto3" json:"uints64,omitempty"`
	Ints8         []int32                `protobuf:"varint,14,rep,packed,name=ints8,proto3" json:"ints8,omitempty"`
	Ints16        []int32                `protobuf:"varint,15,rep,packed,name=ints16,proto3" json:"ints16,omitempty"`
	Ints32        []int32                `protobuf:"varint,16,rep,packed,name=ints32,proto3" json:"ints32,omitempty"`
	Ints64        []int64                `protobuf:"varint,17,rep,packed,name=ints64,proto3" json:"ints64,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StatProto) Reset() {
	*x = StatProto{}
	mi := &file_test_ch_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StatProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatProto) ProtoMessage() {}

func (x *StatProto) ProtoReflect() protoreflect.Message {
	mi := &file_test_ch_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatProto.ProtoReflect.Descriptor instead.
func (*StatProto) Descriptor() ([]byte, []int) {
	return file_test_ch_proto_rawDescGZIP(), []int{0}
}

func (x *StatProto) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *StatProto) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *StatProto) GetType() int64 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *StatProto) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *StatProto) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *StatProto) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *StatProto) GetImage() bool {
	if x != nil {
		return x.Image
	}
	return false
}

func (x *StatProto) GetStack() []string {
	if x != nil {
		return x.Stack
	}
	return nil
}

func (x *StatProto) GetUints() []uint64 {
	if x != nil {
		return x.Uints
	}
	return nil
}

func (x *StatProto) GetUints8() []byte {
	if x != nil {
		return x.Uints8
	}
	return nil
}

func (x *StatProto) GetUints16() []uint32 {
	if x != nil {
		return x.Uints16
	}
	return nil
}

func (x *StatProto) GetUints32() []uint32 {
	if x != nil {
		return x.Uints32
	}
	return nil
}

func (x *StatProto) GetUints64() []uint64 {
	if x != nil {
		return x.Uints64
	}
	return nil
}

func (x *StatProto) GetInts8() []int32 {
	if x != nil {
		return x.Ints8
	}
	return nil
}

func (x *StatProto) GetInts16() []int32 {
	if x != nil {
		return x.Ints16
	}
	return nil
}

func (x *StatProto) GetInts32() []int32 {
	if x != nil {
		return x.Ints32
	}
	return nil
}

func (x *StatProto) GetInts64() []int64 {
	if x != nil {
		return x.Ints64
	}
	return nil
}

var File_test_ch_proto protoreflect.FileDescriptor

const file_test_ch_proto_rawDesc = "" +
	"\n" +
	"\rtest_ch.proto\x12\x04news\"\x83\x03\n" +
	"\tStatProto\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x10\n" +
	"\x03uid\x18\x02 \x01(\x03R\x03uid\x12\x12\n" +
	"\x04type\x18\x03 \x01(\x03R\x04type\x12\x12\n" +
	"\x04name\x18\x04 \x01(\tR\x04name\x12\x16\n" +
	"\x06source\x18\x05 \x01(\tR\x06source\x12\x0e\n" +
	"\x02ip\x18\x06 \x01(\tR\x02ip\x12\x14\n" +
	"\x05image\x18\a \x01(\bR\x05image\x12\x14\n" +
	"\x05stack\x18\b \x03(\tR\x05stack\x12\x14\n" +
	"\x05uints\x18\t \x03(\x04R\x05uints\x12\x16\n" +
	"\x06uints8\x18\n" +
	" \x01(\fR\x06uints8\x12\x18\n" +
	"\auints16\x18\v \x03(\rR\auints16\x12\x18\n" +
	"\auints32\x18\f \x03(\rR\auints32\x12\x18\n" +
	"\auints64\x18\r \x03(\x04R\auints64\x12\x14\n" +
	"\x05ints8\x18\x0e \x03(\x05R\x05ints8\x12\x16\n" +
	"\x06ints16\x18\x0f \x03(\x05R\x06ints16\x12\x16\n" +
	"\x06ints32\x18\x10 \x03(\x05R\x06ints32\x12\x16\n" +
	"\x06ints64\x18\x11 \x03(\x03R\x06ints64B\tZ\a./;newsb\x06proto3"

var (
	file_test_ch_proto_rawDescOnce sync.Once
	file_test_ch_proto_rawDescData []byte
)

func file_test_ch_proto_rawDescGZIP() []byte {
	file_test_ch_proto_rawDescOnce.Do(func() {
		file_test_ch_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_test_ch_proto_rawDesc), len(file_test_ch_proto_rawDesc)))
	})
	return file_test_ch_proto_rawDescData
}

var file_test_ch_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_test_ch_proto_goTypes = []any{
	(*StatProto)(nil), // 0: news.StatProto
}
var file_test_ch_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

var _ = initStatProto()

func initStatProto() struct{} {
	file_test_ch_proto_init()
	return struct{}{}
}
func file_test_ch_proto_init() {
	if File_test_ch_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_test_ch_proto_rawDesc), len(file_test_ch_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_test_ch_proto_goTypes,
		DependencyIndexes: file_test_ch_proto_depIdxs,
		MessageInfos:      file_test_ch_proto_msgTypes,
	}.Build()
	File_test_ch_proto = out.File
	file_test_ch_proto_goTypes = nil
	file_test_ch_proto_depIdxs = nil
}

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson73d7d728DecodeGithubComMonopollyJsonsgeneratorTest(in *jlexer.Lexer, out *Stat) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		switch key {
		case "id":
			if in.IsNull() {
				in.Skip()
			} else {
				out.ID = int(in.Int())
			}
		case "uid":
			if in.IsNull() {
				in.Skip()
			} else {
				out.UID = int(in.Int())
			}
		case "type":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Type = int(in.Int())
			}
		case "name":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Name = string(in.String())
			}
		case "source":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Source = string(in.String())
			}
		case "ip":
			if in.IsNull() {
				in.Skip()
			} else {
				out.IP = string(in.String())
			}
		case "image":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Image = bool(in.Bool())
			}
		case "stack":
			if in.IsNull() {
				in.Skip()
				out.Stack = nil
			} else {
				in.Delim('[')
				if out.Stack == nil {
					if !in.IsDelim(']') {
						out.Stack = make([]string, 0, 4)
					} else {
						out.Stack = []string{}
					}
				} else {
					out.Stack = (out.Stack)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					if in.IsNull() {
						in.Skip()
					} else {
						v1 = string(in.String())
					}
					out.Stack = append(out.Stack, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "uints":
			if in.IsNull() {
				in.Skip()
				out.Uints = nil
			} else {
				in.Delim('[')
				if out.Uints == nil {
					if !in.IsDelim(']') {
						out.Uints = make([]uint, 0, 8)
					} else {
						out.Uints = []uint{}
					}
				} else {
					out.Uints = (out.Uints)[:0]
				}
				for !in.IsDelim(']') {
					var v2 uint
					if in.IsNull() {
						in.Skip()
					} else {
						v2 = uint(in.Uint())
					}
					out.Uints = append(out.Uints, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "uints8":
			if in.IsNull() {
				in.Skip()
				out.Uints8 = nil
			} else {
				out.Uints8 = in.Bytes()
			}
		case "uints16":
			if in.IsNull() {
				in.Skip()
				out.Uints16 = nil
			} else {
				in.Delim('[')
				if out.Uints16 == nil {
					if !in.IsDelim(']') {
						out.Uints16 = make([]uint16, 0, 32)
					} else {
						out.Uints16 = []uint16{}
					}
				} else {
					out.Uints16 = (out.Uints16)[:0]
				}
				for !in.IsDelim(']') {
					var v4 uint16
					if in.IsNull() {
						in.Skip()
					} else {
						v4 = uint16(in.Uint16())
					}
					out.Uints16 = append(out.Uints16, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "uints32":
			if in.IsNull() {
				in.Skip()
				out.Uints32 = nil
			} else {
				in.Delim('[')
				if out.Uints32 == nil {
					if !in.IsDelim(']') {
						out.Uints32 = make([]uint32, 0, 16)
					} else {
						out.Uints32 = []uint32{}
					}
				} else {
					out.Uints32 = (out.Uints32)[:0]
				}
				for !in.IsDelim(']') {
					var v5 uint32
					if in.IsNull() {
						in.Skip()
					} else {
						v5 = uint32(in.Uint32())
					}
					out.Uints32 = append(out.Uints32, v5)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "uints64":
			if in.IsNull() {
				in.Skip()
				out.Uints64 = nil
			} else {
				in.Delim('[')
				if out.Uints64 == nil {
					if !in.IsDelim(']') {
						out.Uints64 = make([]uint64, 0, 8)
					} else {
						out.Uints64 = []uint64{}
					}
				} else {
					out.Uints64 = (out.Uints64)[:0]
				}
				for !in.IsDelim(']') {
					var v6 uint64
					if in.IsNull() {
						in.Skip()
					} else {
						v6 = uint64(in.Uint64())
					}
					out.Uints64 = append(out.Uints64, v6)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "ints8":
			if in.IsNull() {
				in.Skip()
				out.Ints8 = nil
			} else {
				in.Delim('[')
				if out.Ints8 == nil {
					if !in.IsDelim(']') {
						out.Ints8 = make([]int8, 0, 64)
					} else {
						out.Ints8 = []int8{}
					}
				} else {
					out.Ints8 = (out.Ints8)[:0]
				}
				for !in.IsDelim(']') {
					var v7 int8
					if in.IsNull() {
						in.Skip()
					} else {
						v7 = int8(in.Int8())
					}
					out.Ints8 = append(out.Ints8, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "ints16":
			if in.IsNull() {
				in.Skip()
				out.Ints16 = nil
			} else {
				in.Delim('[')
				if out.Ints16 == nil {
					if !in.IsDelim(']') {
						out.Ints16 = make([]int16, 0, 32)
					} else {
						out.Ints16 = []int16{}
					}
				} else {
					out.Ints16 = (out.Ints16)[:0]
				}
				for !in.IsDelim(']') {
					var v8 int16
					if in.IsNull() {
						in.Skip()
					} else {
						v8 = int16(in.Int16())
					}
					out.Ints16 = append(out.Ints16, v8)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "ints32":
			if in.IsNull() {
				in.Skip()
				out.Ints32 = nil
			} else {
				in.Delim('[')
				if out.Ints32 == nil {
					if !in.IsDelim(']') {
						out.Ints32 = make([]int32, 0, 16)
					} else {
						out.Ints32 = []int32{}
					}
				} else {
					out.Ints32 = (out.Ints32)[:0]
				}
				for !in.IsDelim(']') {
					var v9 int32
					if in.IsNull() {
						in.Skip()
					} else {
						v9 = int32(in.Int32())
					}
					out.Ints32 = append(out.Ints32, v9)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "ints64":
			if in.IsNull() {
				in.Skip()
				out.Ints64 = nil
			} else {
				in.Delim('[')
				if out.Ints64 == nil {
					if !in.IsDelim(']') {
						out.Ints64 = make([]int64, 0, 8)
					} else {
						out.Ints64 = []int64{}
					}
				} else {
					out.Ints64 = (out.Ints64)[:0]
				}
				for !in.IsDelim(']') {
					var v10 int64
					if in.IsNull() {
						in.Skip()
					} else {
						v10 = int64(in.Int64())
					}
					out.Ints64 = append(out.Ints64, v10)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson73d7d728EncodeGithubComMonopollyJsonsgeneratorTest(out *jwriter.Writer, in Stat) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != 0 {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	if in.UID != 0 {
		const prefix string = ",\"uid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.UID))
	}
	if in.Type != 0 {
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Type))
	}
	if in.Name != "" {
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	if in.Source != "" {
		const prefix string = ",\"source\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Source))
	}
	if in.IP != "" {
		const prefix string = ",\"ip\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.IP))
	}
	if in.Image {
		const prefix string = ",\"image\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Image))
	}
	if len(in.Stack) != 0 {
		const prefix string = ",\"stack\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v11, v12 := range in.Stack {
				if v11 > 0 {
					out.RawByte(',')
				}
				out.String(string(v12))
			}
			out.RawByte(']')
		}
	}
	if len(in.Uints) != 0 {
		const prefix string = ",\"uints\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v13, v14 := range in.Uints {
				if v13 > 0 {
					out.RawByte(',')
				}
				out.Uint(uint(v14))
			}
			out.RawByte(']')
		}
	}
	if len(in.Uints8) != 0 {
		const prefix string = ",\"uints8\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Base64Bytes(in.Uints8)
	}
	if len(in.Uints16) != 0 {
		const prefix string = ",\"uints16\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v17, v18 := range in.Uints16 {
				if v17 > 0 {
					out.RawByte(',')
				}
				out.Uint16(uint16(v18))
			}
			out.RawByte(']')
		}
	}
	if len(in.Uints32) != 0 {
		const prefix string = ",\"uints32\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v19, v20 := range in.Uints32 {
				if v19 > 0 {
					out.RawByte(',')
				}
				out.Uint32(uint32(v20))
			}
			out.RawByte(']')
		}
	}
	if len(in.Uints64) != 0 {
		const prefix string = ",\"uints64\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v21, v22 := range in.Uints64 {
				if v21 > 0 {
					out.RawByte(',')
				}
				out.Uint64(uint64(v22))
			}
			out.RawByte(']')
		}
	}
	if len(in.Ints8) != 0 {
		const prefix string = ",\"ints8\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v23, v24 := range in.Ints8 {
				if v23 > 0 {
					out.RawByte(',')
				}
				out.Int8(int8(v24))
			}
			out.RawByte(']')
		}
	}
	if len(in.Ints16) != 0 {
		const prefix string = ",\"ints16\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v25, v26 := range in.Ints16 {
				if v25 > 0 {
					out.RawByte(',')
				}
				out.Int16(int16(v26))
			}
			out.RawByte(']')
		}
	}
	if len(in.Ints32) != 0 {
		const prefix string = ",\"ints32\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v27, v28 := range in.Ints32 {
				if v27 > 0 {
					out.RawByte(',')
				}
				out.Int32(int32(v28))
			}
			out.RawByte(']')
		}
	}
	if len(in.Ints64) != 0 {
		const prefix string = ",\"ints64\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v29, v30 := range in.Ints64 {
				if v29 > 0 {
					out.RawByte(',')
				}
				out.Int64(int64(v30))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Stat) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson73d7d728EncodeGithubComMonopollyJsonsgeneratorTest(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Stat) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson73d7d728EncodeGithubComMonopollyJsonsgeneratorTest(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Stat) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson73d7d728DecodeGithubComMonopollyJsonsgeneratorTest(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Stat) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson73d7d728DecodeGithubComMonopollyJsonsgeneratorTest(l, v)
}

// easyjson marshal
func (a *Stat) Marshal() []byte {
	b, _ := easyjson.Marshal(a)
	return b
}

// easyjson unmarshal
func StatUnmarshal(src []byte) (a *Stat) {
	a = new(Stat)
	err := easyjson.Unmarshal(src, a)
	if err != nil {
		return nil
	}
	return
}
