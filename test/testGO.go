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
	"github.com/vmihailenco/msgpack/v5"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
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
	debug                         Debug mode
	lock                          Lock model for generation. Can't change model.
	gotiny                        Create gotiny marshal/unmarshal
	optimize                      Try to create golang padding optimized struct
	chengine=MergeTree            Optional, mergeTree by default
	msgp                          Create message pack marshal/unmarshal
	!omit                         No omit tag for json
	proto=NewsProto               Generate proto file and compile to current package
	swift                         Generate swift model
	enum                          Generate swift enum
	demo                          Generate demo json file with default values
	ch=views                      Set clickhouse sql table name. Ex: ch=views
	test                          Generate go tests for generated model
	noprefix                      Generate simple index IndexID instead IndexNewsID
	sql=news                      Set sql table name. Ex: sql=accounts
	noinit                        No New() init function for struct
	go=News                       Set golang struct names. Ex: go=News1 > type News1 struct{}
	js=NewsJson                   Change json struct names. Ex: js=NewsJson
	ts=news                       Set typescript struct names. Ex: ts=NewsJson

	Field:
	title{}                       Add custom title for index. title{Nice & Sweet}
	desc{}                        Add custom desc for index. desc{This field for success}

	GO:
	up                            Make uppercase for functions
	type=""                       Replace golang type. Ex: type="[]*News"
	desc{}                        Add field user title desc{Use it nice}
	@                             Add custom structs for fields. Ex: id int //@public @team...
	must                          Create one validation function for all must fields
	name=""                       Replace golang struct name. Ex: name="NewsList"
	title=""                      Add custom title for index
	desc=""                       Add custom desc for index
	req{}                         Add user required fields
	title{}                       Add field user title title{Nice}
	#                             Add lists for fields. Ex: id int //#readonly #must...
	nofunc                        Do not create any jsons functions for fiels

	SQL:
	noinsert                      Do use field for insert function
	altertable                    Add field line Alter table to SQL file with current time comment
	idx                           Add index fields by group name. Ex: idx="nameIndex" idx="credsIndex"
	unique="groupname"            Add unique fields constrains by group (you need set group name, then generator join fields). Ex: unique="group1" unique="group2"
	ver="v4"                      Create an alter table record in SQL file. Add new column. With new version in comment. Ex: ver="2"
	renames="oldname"             Create an alter table record in SQL file. Rename table column.
	index                         Create simple default index or gin for jsonb
	search                        Add tsvector index by group. Ex: search="tsv": tsv tsvector GENERATED ALWAYS AS (to_tsvector('simple', title || ' ' || brand)). For search: SELECT brand, title FROM assets WHERE search @@ to_tsquery('english', 'f8');
	defaults                      Add default value based on field type
	unix                          Add default value: extract(epoch from now())
	skip                          Skip field for sql queries
	replace="bigint primary key"  Rewrite sql for field. Ex: replace="bigint primary key"
	type="jsonb"                  Rewrite sql type for field. Ex: type="jsonb"
	add="primary key"             Append sql for field. Ex: replace="primary key"

	Clickhouse:
	low                           Wrap field type in LowCardinality(...)
	type=ip                       Use ClickHouse IPv4 type for IP address fields

	PROTO:
	type="string"                 Rewrite proto type for field. Ex: type="google.protobuf.Value"
	name="field"                  Rewrite proto field name
	skip                          Skip field for proto

	SWIFT:
	must                          Swift field with required values
	skip                          Ignore field for Swift
	type=""                       Replace type for Swift model. Ex: type="string"
	file=""                       Create another swift file for this field Swift model. file="f1", file="f2"

	JSON:
	skip                          Skip json field
	inc                           Add inc jsons function for numbers fields
	bool                          Add set jsons function for bool fields
	time                          Create convert jsons function for unixtime fields
	raw                           Set Raw json function inside field
	name=""                       Replace json field name. Ex: name="sid"
*/

// field type
type NewsIndexType int

// int index
const (
	IndexInc = NewsIndexType(iota)
	IndexInts
	IndexInts8
	IndexInts16
	IndexInts32
	IndexInts64
	IndexUints
	IndexUints8
	IndexUints16
	IndexUints32
	IndexUints64
	IndexFloats32
	IndexFloats64
	IndexBools
	IndexByte1
	IndexBytes
	IndexList_ints
	IndexList_string
	IndexList_float
	IndexMap_string_string
	IndexMap_string_bytes
	IndexMap_string_bool
	IndexMap_string_int
	IndexMap_string_float64
	IndexMap_string_any
	IndexMap_int_string
	IndexMap_int_int
	IndexMap_int_bool
	IndexRenameSQL
	IndexRenameGO_OK
	IndexRenameJS
	IndexMAST_UPPER_GO
	IndexIntToSmallInt
	IndexSkip
	IndexSql_unique_u1_1
	IndexSql_unique_u1_2
	IndexSql_index1_1
	IndexSql_index1_2
	IndexSql_index1_3
	IndexSql_keys_1
	IndexSql_keys_2
	IndexSql_keys_3
	IndexSql_search
	IndexSql_get
	IndexSql_unique_x1
	IndexSql_unique_x2
	IndexSql_unique_x1_x2
	IndexSql_primary
	IndexSql_jsonb_index
	IndexTime_duration
	IndexGo_type_int_to_strings
	IndexPublic_field1
	IndexPublic_field2
	IndexPublic_field3
	IndexPublic_field_me1
	IndexPublic_field_me2
	IndexPublic_field_me3
)

// string index
const (
	FieldInc                    = "inc"                    // int sql{inc}  #readonly
	FieldInts                   = "ints"                   // int
	FieldInts8                  = "ints8"                  // int8
	FieldInts16                 = "ints16"                 // int16
	FieldInts32                 = "ints32"                 // int32
	FieldInts64                 = "ints64"                 // int64
	FieldUints                  = "uints"                  // uint
	FieldUints8                 = "uints8"                 // uint8
	FieldUints16                = "uints16"                // uint16
	FieldUints32                = "uints32"                // uint32
	FieldUints64                = "uints64"                // uint64
	FieldFloats32               = "floats32"               // float32
	FieldFloats64               = "floats64"               // float64
	FieldBools                  = "bools"                  // bool
	FieldByte1                  = "byte1"                  // byte
	FieldBytes                  = "bytes"                  // []byte
	FieldList_ints              = "list_ints"              // []int
	FieldList_string            = "list_string"            // []string
	FieldList_float             = "list_float"             // []float64
	FieldMap_string_string      = "map_string_string"      // map[string]string
	FieldMap_string_bytes       = "map_string_bytes"       // map[string][]byte
	FieldMap_string_bool        = "map_string_bool"        // map[string]bool
	FieldMap_string_int         = "map_string_int"         // map[string]int
	FieldMap_string_float64     = "map_string_float64"     // map[string]float64
	FieldMap_string_any         = "map_string_any"         // map[string]any
	FieldMap_int_string         = "map_int_string"         // map[int]string
	FieldMap_int_int            = "map_int_int"            // map[int]int
	FieldMap_int_bool           = "map_int_bool"           // map[int]bool
	FieldRenameSQL              = "renameSQL"              // string sql{name="renameSQL_OK"}
	FieldRenameGO_OK            = "renameGO"               // string go{name="RenameGO_OK"}
	FieldRenameJS               = "renameJS_OK"            // string js{name="renameJS_OK"}
	FieldMAST_UPPER_GO          = "mast_upper_go"          // string go{up}
	FieldIntToSmallInt          = "intToSmallInt"          // int sql{type="smallint"} rename
	FieldSkip                   = "skip"                   // string sql{skip} swift{skip} js{skip}
	FieldSql_unique_u1_1        = "sql_unique_u1_1"        // int sql{unique="u1"}
	FieldSql_unique_u1_2        = "sql_unique_u1_2"        // int sql{unique="u1"}
	FieldSql_index1_1           = "sql_index1_1"           // int sql{idx="index1"}
	FieldSql_index1_2           = "sql_index1_2"           // int sql{idx="index1"}
	FieldSql_index1_3           = "sql_index1_3"           // int sql{idx="index1"}
	FieldSql_keys_1             = "sql_keys_1"             // int sql{keys="keys1"}
	FieldSql_keys_2             = "sql_keys_2"             // int sql{keys="keys1"}
	FieldSql_keys_3             = "sql_keys_3"             // int sql{keys="keys1"}
	FieldSql_search             = "sql_search"             // string sql{search="search" get}
	FieldSql_get                = "sql_get"                // string sql{get}
	FieldSql_unique_x1          = "sql_unique_x1"          // int sql{unique="x1"}
	FieldSql_unique_x2          = "sql_unique_x2"          // int sql{unique="x2"}
	FieldSql_unique_x1_x2       = "sql_unique_x1_x2"       // int sql{unique="x1", unique="x2"}
	FieldSql_primary            = "sql_primary"            // float64 sql{primarykey}
	FieldSql_jsonb_index        = "sql_jsonb_index"        // map[string]any sql{index}
	FieldTime_duration          = "time_duration"          // time.Duration
	FieldGo_type_int_to_strings = "go_type_int_to_strings" // int go{type="[]string"}
	FieldPublic_field1          = "public_field1"          // int @public
	FieldPublic_field2          = "public_field2"          // int @public
	FieldPublic_field3          = "public_field3"          // int @public
	FieldPublic_field_me1       = "public_field_me1"       // int @me
	FieldPublic_field_me2       = "public_field_me2"       // int @me
	FieldPublic_field_me3       = "public_field_me3"       // int @me
)

// index func
func NewsIndexes() []NewsIndexType {
	return []NewsIndexType{IndexInc, IndexInts, IndexInts8, IndexInts16, IndexInts32, IndexInts64, IndexUints, IndexUints8, IndexUints16, IndexUints32, IndexUints64, IndexFloats32, IndexFloats64, IndexBools, IndexByte1, IndexBytes, IndexList_ints, IndexList_string, IndexList_float, IndexMap_string_string, IndexMap_string_bytes, IndexMap_string_bool, IndexMap_string_int, IndexMap_string_float64, IndexMap_string_any, IndexMap_int_string, IndexMap_int_int, IndexMap_int_bool, IndexRenameSQL, IndexRenameGO_OK, IndexRenameJS, IndexMAST_UPPER_GO, IndexIntToSmallInt, IndexSkip, IndexSql_unique_u1_1, IndexSql_unique_u1_2, IndexSql_index1_1, IndexSql_index1_2, IndexSql_index1_3, IndexSql_keys_1, IndexSql_keys_2, IndexSql_keys_3, IndexSql_search, IndexSql_get, IndexSql_unique_x1, IndexSql_unique_x2, IndexSql_unique_x1_x2, IndexSql_primary, IndexSql_jsonb_index, IndexTime_duration, IndexGo_type_int_to_strings, IndexPublic_field1, IndexPublic_field2, IndexPublic_field3, IndexPublic_field_me1, IndexPublic_field_me2, IndexPublic_field_me3}

}

// 472 bytes (go padding)
//
//easyjson:json
type News struct {
	Inc                    int                `json:"inc,omitempty" msg:"inc,omitempty"`                                       // sql{inc}  #readonly
	Ints                   int                `json:"ints,omitempty" msg:"ints,omitempty"`                                     //
	Ints8                  int8               `json:"ints8,omitempty" msg:"ints8,omitempty"`                                   //
	Ints16                 int16              `json:"ints16,omitempty" msg:"ints16,omitempty"`                                 //
	Ints32                 int32              `json:"ints32,omitempty" msg:"ints32,omitempty"`                                 //
	Ints64                 int64              `json:"ints64,omitempty" msg:"ints64,omitempty"`                                 //
	Uints                  uint               `json:"uints,omitempty" msg:"uints,omitempty"`                                   //
	Uints8                 uint8              `json:"uints8,omitempty" msg:"uints8,omitempty"`                                 //
	Uints16                uint16             `json:"uints16,omitempty" msg:"uints16,omitempty"`                               //
	Uints32                uint32             `json:"uints32,omitempty" msg:"uints32,omitempty"`                               //
	Uints64                uint64             `json:"uints64,omitempty" msg:"uints64,omitempty"`                               //
	Floats32               float32            `json:"floats32,omitempty" msg:"floats32,omitempty"`                             //
	Floats64               float64            `json:"floats64,omitempty" msg:"floats64,omitempty"`                             //
	Bools                  bool               `json:"bools,omitempty" msg:"bools,omitempty"`                                   //
	Byte1                  byte               `json:"byte1,omitempty" msg:"byte1,omitempty"`                                   //
	Bytes                  []byte             `json:"bytes,omitempty" msg:"bytes,omitempty"`                                   //
	List_ints              []int              `json:"list_ints,omitempty" msg:"list_ints,omitempty"`                           //
	List_string            []string           `json:"list_string,omitempty" msg:"list_string,omitempty"`                       //
	List_float             []float64          `json:"list_float,omitempty" msg:"list_float,omitempty"`                         //
	Map_string_string      map[string]string  `json:"map_string_string,omitempty" msg:"map_string_string,omitempty"`           //
	Map_string_bytes       map[string][]byte  `json:"map_string_bytes,omitempty" msg:"map_string_bytes,omitempty"`             //
	Map_string_bool        map[string]bool    `json:"map_string_bool,omitempty" msg:"map_string_bool,omitempty"`               //
	Map_string_int         map[string]int     `json:"map_string_int,omitempty" msg:"map_string_int,omitempty"`                 //
	Map_string_float64     map[string]float64 `json:"map_string_float64,omitempty" msg:"map_string_float64,omitempty"`         //
	Map_string_any         map[string]any     `json:"map_string_any,omitempty" msg:"map_string_any,omitempty"`                 //
	Map_int_string         map[int]string     `json:"map_int_string,omitempty" msg:"map_int_string,omitempty"`                 //
	Map_int_int            map[int]int        `json:"map_int_int,omitempty" msg:"map_int_int,omitempty"`                       //
	Map_int_bool           map[int]bool       `json:"map_int_bool,omitempty" msg:"map_int_bool,omitempty"`                     //
	RenameSQL              string             `json:"renameSQL,omitempty" msg:"renameSQL,omitempty"`                           // sql{name="renameSQL_OK"}
	RenameGO_OK            string             `json:"renameGO,omitempty" msg:"renameGO,omitempty"`                             // go{name="RenameGO_OK"}
	RenameJS               string             `json:"renameJS_OK,omitempty" msg:"renameJS_OK,omitempty"`                       // js{name="renameJS_OK"}
	MAST_UPPER_GO          string             `json:"mast_upper_go,omitempty" msg:"mast_upper_go,omitempty"`                   // go{up}
	IntToSmallInt          int                `json:"intToSmallInt,omitempty" msg:"intToSmallInt,omitempty"`                   // sql{type="smallint"} rename
	Skip                   string             `json:"skip,omitempty" msg:"skip,omitempty"`                                     // sql{skip} swift{skip} js{skip}
	Sql_unique_u1_1        int                `json:"sql_unique_u1_1,omitempty" msg:"sql_unique_u1_1,omitempty"`               // sql{unique="u1"}
	Sql_unique_u1_2        int                `json:"sql_unique_u1_2,omitempty" msg:"sql_unique_u1_2,omitempty"`               // sql{unique="u1"}
	Sql_index1_1           int                `json:"sql_index1_1,omitempty" msg:"sql_index1_1,omitempty"`                     // sql{idx="index1"}
	Sql_index1_2           int                `json:"sql_index1_2,omitempty" msg:"sql_index1_2,omitempty"`                     // sql{idx="index1"}
	Sql_index1_3           int                `json:"sql_index1_3,omitempty" msg:"sql_index1_3,omitempty"`                     // sql{idx="index1"}
	Sql_keys_1             int                `json:"sql_keys_1,omitempty" msg:"sql_keys_1,omitempty"`                         // sql{keys="keys1"}
	Sql_keys_2             int                `json:"sql_keys_2,omitempty" msg:"sql_keys_2,omitempty"`                         // sql{keys="keys1"}
	Sql_keys_3             int                `json:"sql_keys_3,omitempty" msg:"sql_keys_3,omitempty"`                         // sql{keys="keys1"}
	Sql_search             string             `json:"sql_search,omitempty" msg:"sql_search,omitempty"`                         // sql{search="search" get}
	Sql_get                string             `json:"sql_get,omitempty" msg:"sql_get,omitempty"`                               // sql{get}
	Sql_unique_x1          int                `json:"sql_unique_x1,omitempty" msg:"sql_unique_x1,omitempty"`                   // sql{unique="x1"}
	Sql_unique_x2          int                `json:"sql_unique_x2,omitempty" msg:"sql_unique_x2,omitempty"`                   // sql{unique="x2"}
	Sql_unique_x1_x2       int                `json:"sql_unique_x1_x2,omitempty" msg:"sql_unique_x1_x2,omitempty"`             // sql{unique="x1", unique="x2"}
	Sql_primary            float64            `json:"sql_primary,omitempty" msg:"sql_primary,omitempty"`                       // sql{primarykey}
	Sql_jsonb_index        map[string]any     `json:"sql_jsonb_index,omitempty" msg:"sql_jsonb_index,omitempty"`               // sql{index}
	Time_duration          time.Duration      `json:"time_duration,omitempty" msg:"time_duration,omitempty"`                   //
	Go_type_int_to_strings []string           `json:"go_type_int_to_strings,omitempty" msg:"go_type_int_to_strings,omitempty"` // go{type="[]string"}
	Public_field1          int                `json:"public_field1,omitempty" msg:"public_field1,omitempty"`                   // @public
	Public_field2          int                `json:"public_field2,omitempty" msg:"public_field2,omitempty"`                   // @public
	Public_field3          int                `json:"public_field3,omitempty" msg:"public_field3,omitempty"`                   // @public
	Public_field_me1       int                `json:"public_field_me1,omitempty" msg:"public_field_me1,omitempty"`             // @me
	Public_field_me2       int                `json:"public_field_me2,omitempty" msg:"public_field_me2,omitempty"`             // @me
	Public_field_me3       int                `json:"public_field_me3,omitempty" msg:"public_field_me3,omitempty"`             // @me
}

type NewsPublic struct {
	Public_field1 int `json:"public_field1,omitempty" msg:"public_field1,omitempty"` // @public
	Public_field2 int `json:"public_field2,omitempty" msg:"public_field2,omitempty"` // @public
	Public_field3 int `json:"public_field3,omitempty" msg:"public_field3,omitempty"` // @public
}

func (a *News) NewsPublic() (res *NewsPublic) {
	return &NewsPublic{
		Public_field1: a.Public_field1,
		Public_field2: a.Public_field2,
		Public_field3: a.Public_field3,
	}
}

func (a *NewsPublic) News() (res *News) {
	return &News{
		Public_field1: a.Public_field1,
		Public_field2: a.Public_field2,
		Public_field3: a.Public_field3,
	}
}

func (a *NewsPublic) Pack() (res []byte) {
	res, _ = jsoniter.Marshal(a)
	return
}

type NewsMe struct {
	Public_field_me1 int `json:"public_field_me1,omitempty" msg:"public_field_me1,omitempty"` // @me
	Public_field_me2 int `json:"public_field_me2,omitempty" msg:"public_field_me2,omitempty"` // @me
	Public_field_me3 int `json:"public_field_me3,omitempty" msg:"public_field_me3,omitempty"` // @me
}

func (a *News) NewsMe() (res *NewsMe) {
	return &NewsMe{
		Public_field_me1: a.Public_field_me1,
		Public_field_me2: a.Public_field_me2,
		Public_field_me3: a.Public_field_me3,
	}
}

func (a *NewsMe) News() (res *News) {
	return &News{
		Public_field_me1: a.Public_field_me1,
		Public_field_me2: a.Public_field_me2,
		Public_field_me3: a.Public_field_me3,
	}
}

func (a *NewsMe) Pack() (res []byte) {
	res, _ = jsoniter.Marshal(a)
	return
}

// init struct
func NewNews() (a *News) {
	a = new(News)
	a.Map_string_string = make(map[string]string)
	a.Map_string_bytes = make(map[string][]byte)
	a.Map_string_bool = make(map[string]bool)
	a.Map_string_int = make(map[string]int)
	a.Map_string_float64 = make(map[string]float64)
	a.Map_string_any = make(map[string]any)
	a.Map_int_string = make(map[int]string)
	a.Map_int_int = make(map[int]int)
	a.Map_int_bool = make(map[int]bool)
	a.Sql_jsonb_index = make(map[string]any)
	return
}

// Parse []any to ID struct
func ParseNewsToStruct(r []any) (a *News) {
	a = new(News)

	for pos, x := range r {
		switch NewsIndexType(pos) {
		case IndexInc:
			cast.Convert(&a.Inc, x) //int
		case IndexInts:
			cast.Convert(&a.Ints, x) //int
		case IndexInts8:
			cast.Convert(&a.Ints8, x) //int8
		case IndexInts16:
			cast.Convert(&a.Ints16, x) //int16
		case IndexInts32:
			cast.Convert(&a.Ints32, x) //int32
		case IndexInts64:
			cast.Convert(&a.Ints64, x) //int64
		case IndexUints:
			cast.Convert(&a.Uints, x) //uint
		case IndexUints8:
			cast.Convert(&a.Uints8, x) //uint8
		case IndexUints16:
			cast.Convert(&a.Uints16, x) //uint16
		case IndexUints32:
			cast.Convert(&a.Uints32, x) //uint32
		case IndexUints64:
			cast.Convert(&a.Uints64, x) //uint64
		case IndexFloats32:
			cast.Convert(&a.Floats32, x) //float32
		case IndexFloats64:
			cast.Convert(&a.Floats64, x) //float64
		case IndexBools:
			cast.Convert(&a.Bools, x) //bool
		case IndexByte1:
			cast.Convert(&a.Byte1, x) //byte
		case IndexBytes:
			cast.Convert(&a.Bytes, x) //[]byte
		case IndexList_ints:
			cast.Convert(&a.List_ints, x) //[]int
		case IndexList_string:
			cast.Convert(&a.List_string, x) //[]string
		case IndexList_float:
			cast.Convert(&a.List_float, x) //[]float64
		case IndexMap_string_string:
			cast.Convert(&a.Map_string_string, x) //map[string]string
		case IndexMap_string_bytes:
			cast.Convert(&a.Map_string_bytes, x) //map[string][]byte
		case IndexMap_string_bool:
			cast.Convert(&a.Map_string_bool, x) //map[string]bool
		case IndexMap_string_int:
			cast.Convert(&a.Map_string_int, x) //map[string]int
		case IndexMap_string_float64:
			cast.Convert(&a.Map_string_float64, x) //map[string]float64
		case IndexMap_string_any:
			cast.Convert(&a.Map_string_any, x) //map[string]any
		case IndexMap_int_string:
			cast.Convert(&a.Map_int_string, x) //map[int]string
		case IndexMap_int_int:
			cast.Convert(&a.Map_int_int, x) //map[int]int
		case IndexMap_int_bool:
			cast.Convert(&a.Map_int_bool, x) //map[int]bool
		case IndexRenameSQL:
			cast.Convert(&a.RenameSQL, x) //string
		case IndexRenameGO_OK:
			cast.Convert(&a.RenameGO_OK, x) //string
		case IndexRenameJS:
			cast.Convert(&a.RenameJS, x) //string
		case IndexMAST_UPPER_GO:
			cast.Convert(&a.MAST_UPPER_GO, x) //string
		case IndexIntToSmallInt:
			cast.Convert(&a.IntToSmallInt, x) //int
		case IndexSkip:
			cast.Convert(&a.Skip, x) //string
		case IndexSql_unique_u1_1:
			cast.Convert(&a.Sql_unique_u1_1, x) //int
		case IndexSql_unique_u1_2:
			cast.Convert(&a.Sql_unique_u1_2, x) //int
		case IndexSql_index1_1:
			cast.Convert(&a.Sql_index1_1, x) //int
		case IndexSql_index1_2:
			cast.Convert(&a.Sql_index1_2, x) //int
		case IndexSql_index1_3:
			cast.Convert(&a.Sql_index1_3, x) //int
		case IndexSql_keys_1:
			cast.Convert(&a.Sql_keys_1, x) //int
		case IndexSql_keys_2:
			cast.Convert(&a.Sql_keys_2, x) //int
		case IndexSql_keys_3:
			cast.Convert(&a.Sql_keys_3, x) //int
		case IndexSql_search:
			cast.Convert(&a.Sql_search, x) //string
		case IndexSql_get:
			cast.Convert(&a.Sql_get, x) //string
		case IndexSql_unique_x1:
			cast.Convert(&a.Sql_unique_x1, x) //int
		case IndexSql_unique_x2:
			cast.Convert(&a.Sql_unique_x2, x) //int
		case IndexSql_unique_x1_x2:
			cast.Convert(&a.Sql_unique_x1_x2, x) //int
		case IndexSql_primary:
			cast.Convert(&a.Sql_primary, x) //float64
		case IndexSql_jsonb_index:
			cast.Convert(&a.Sql_jsonb_index, x) //map[string]any
		case IndexTime_duration:
			cast.Convert(&a.Time_duration, x) //time.Duration
		case IndexPublic_field1:
			cast.Convert(&a.Public_field1, x) //int
		case IndexPublic_field2:
			cast.Convert(&a.Public_field2, x) //int
		case IndexPublic_field3:
			cast.Convert(&a.Public_field3, x) //int
		case IndexPublic_field_me1:
			cast.Convert(&a.Public_field_me1, x) //int
		case IndexPublic_field_me2:
			cast.Convert(&a.Public_field_me2, x) //int
		case IndexPublic_field_me3:
			cast.Convert(&a.Public_field_me3, x) //int
		}
	}
	return
}

// Tuple create an array from struct
func (a *News) Tuple() (r []any) {
	return []any{a.Inc, a.Ints, a.Ints8, a.Ints16, a.Ints32, a.Ints64, a.Uints, a.Uints8, a.Uints16, a.Uints32, a.Uints64, a.Floats32, a.Floats64, a.Bools, a.Byte1, a.Bytes, a.List_ints, a.List_string, a.List_float, a.Map_string_string, a.Map_string_bytes, a.Map_string_bool, a.Map_string_int, a.Map_string_float64, a.Map_string_any, a.Map_int_string, a.Map_int_int, a.Map_int_bool, a.RenameSQL, a.RenameGO_OK, a.RenameJS, a.MAST_UPPER_GO, a.IntToSmallInt, a.Skip, a.Sql_unique_u1_1, a.Sql_unique_u1_2, a.Sql_index1_1, a.Sql_index1_2, a.Sql_index1_3, a.Sql_keys_1, a.Sql_keys_2, a.Sql_keys_3, a.Sql_search, a.Sql_get, a.Sql_unique_x1, a.Sql_unique_x2, a.Sql_unique_x1_x2, a.Sql_primary, a.Sql_jsonb_index, a.Time_duration, a.Go_type_int_to_strings, a.Public_field1, a.Public_field2, a.Public_field3, a.Public_field_me1, a.Public_field_me2, a.Public_field_me3}
}

// Tuple create an array from struct
func (a *News) sqlTuple() (r []any) {
	return []any{a.Ints, a.Ints8, a.Ints16, a.Ints32, a.Ints64, a.Uints, a.Uints8, a.Uints16, a.Uints32, a.Uints64, a.Floats32, a.Floats64, a.Bools, a.Byte1, a.Bytes, a.List_ints, a.List_string, a.List_float, a.Map_string_string, a.Map_string_bytes, a.Map_string_bool, a.Map_string_int, a.Map_string_float64, a.Map_string_any, a.Map_int_string, a.Map_int_int, a.Map_int_bool, a.RenameSQL, a.RenameGO_OK, a.RenameJS, a.MAST_UPPER_GO, a.IntToSmallInt, a.Sql_unique_u1_1, a.Sql_unique_u1_2, a.Sql_index1_1, a.Sql_index1_2, a.Sql_index1_3, a.Sql_keys_1, a.Sql_keys_2, a.Sql_keys_3, a.Sql_search, a.Sql_get, a.Sql_unique_x1, a.Sql_unique_x2, a.Sql_unique_x1_x2, a.Sql_primary, a.Sql_jsonb_index, a.Time_duration, a.Go_type_int_to_strings, a.Public_field1, a.Public_field2, a.Public_field3, a.Public_field_me1, a.Public_field_me2, a.Public_field_me3}
}

// Tuple create an array from struct
func (a *News) sqlAllTuples() (r []any) {
	return []any{a.Inc, a.Ints, a.Ints8, a.Ints16, a.Ints32, a.Ints64, a.Uints, a.Uints8, a.Uints16, a.Uints32, a.Uints64, a.Floats32, a.Floats64, a.Bools, a.Byte1, a.Bytes, a.List_ints, a.List_string, a.List_float, a.Map_string_string, a.Map_string_bytes, a.Map_string_bool, a.Map_string_int, a.Map_string_float64, a.Map_string_any, a.Map_int_string, a.Map_int_int, a.Map_int_bool, a.RenameSQL, a.RenameGO_OK, a.RenameJS, a.MAST_UPPER_GO, a.IntToSmallInt, a.Sql_unique_u1_1, a.Sql_unique_u1_2, a.Sql_index1_1, a.Sql_index1_2, a.Sql_index1_3, a.Sql_keys_1, a.Sql_keys_2, a.Sql_keys_3, a.Sql_search, a.Sql_get, a.Sql_unique_x1, a.Sql_unique_x2, a.Sql_unique_x1_x2, a.Sql_primary, a.Sql_jsonb_index, a.Time_duration, a.Go_type_int_to_strings, a.Public_field1, a.Public_field2, a.Public_field3, a.Public_field_me1, a.Public_field_me2, a.Public_field_me3}
}

// Tuple create an array from struct for clickhouse
func (a *News) clickhouseTuple() (r []any) {
	return []any{a.Inc, a.Ints, a.Ints8, a.Ints16, a.Ints32, a.Ints64, a.Uints, a.Uints8, a.Uints16, a.Uints32, a.Uints64, a.Floats32, a.Floats64, a.Bools, a.Byte1, a.Bytes, a.List_ints, a.List_string, a.List_float, a.Map_string_string, a.Map_string_bytes, a.Map_string_bool, a.Map_string_int, a.Map_string_float64, a.Map_string_any, a.Map_int_string, a.Map_int_int, a.Map_int_bool, a.RenameSQL, a.RenameGO_OK, a.RenameJS, a.MAST_UPPER_GO, a.IntToSmallInt, a.Skip, a.Sql_unique_u1_1, a.Sql_unique_u1_2, a.Sql_index1_1, a.Sql_index1_2, a.Sql_index1_3, a.Sql_keys_1, a.Sql_keys_2, a.Sql_keys_3, a.Sql_search, a.Sql_get, a.Sql_unique_x1, a.Sql_unique_x2, a.Sql_unique_x1_x2, a.Sql_primary, a.Sql_jsonb_index, a.Time_duration, a.Go_type_int_to_strings, a.Public_field1, a.Public_field2, a.Public_field3, a.Public_field_me1, a.Public_field_me2, a.Public_field_me3}
}

// update struct with function
func (a *News) Update(k string, x any) {
	switch k {
	case "inc":
		cast.Convert(&a.Inc, x) //int
	case "ints":
		cast.Convert(&a.Ints, x) //int
	case "ints8":
		cast.Convert(&a.Ints8, x) //int8
	case "ints16":
		cast.Convert(&a.Ints16, x) //int16
	case "ints32":
		cast.Convert(&a.Ints32, x) //int32
	case "ints64":
		cast.Convert(&a.Ints64, x) //int64
	case "uints":
		cast.Convert(&a.Uints, x) //uint
	case "uints8":
		cast.Convert(&a.Uints8, x) //uint8
	case "uints16":
		cast.Convert(&a.Uints16, x) //uint16
	case "uints32":
		cast.Convert(&a.Uints32, x) //uint32
	case "uints64":
		cast.Convert(&a.Uints64, x) //uint64
	case "floats32":
		cast.Convert(&a.Floats32, x) //float32
	case "floats64":
		cast.Convert(&a.Floats64, x) //float64
	case "bools":
		cast.Convert(&a.Bools, x) //bool
	case "byte1":
		cast.Convert(&a.Byte1, x) //byte
	case "bytes":
		cast.Convert(&a.Bytes, x) //[]byte
	case "list_ints":
		cast.Convert(&a.List_ints, x) //[]int
	case "list_string":
		cast.Convert(&a.List_string, x) //[]string
	case "list_float":
		cast.Convert(&a.List_float, x) //[]float64
	case "map_string_string":
		cast.Convert(&a.Map_string_string, x) //map[string]string
	case "map_string_bytes":
		cast.Convert(&a.Map_string_bytes, x) //map[string][]byte
	case "map_string_bool":
		cast.Convert(&a.Map_string_bool, x) //map[string]bool
	case "map_string_int":
		cast.Convert(&a.Map_string_int, x) //map[string]int
	case "map_string_float64":
		cast.Convert(&a.Map_string_float64, x) //map[string]float64
	case "map_string_any":
		cast.Convert(&a.Map_string_any, x) //map[string]any
	case "map_int_string":
		cast.Convert(&a.Map_int_string, x) //map[int]string
	case "map_int_int":
		cast.Convert(&a.Map_int_int, x) //map[int]int
	case "map_int_bool":
		cast.Convert(&a.Map_int_bool, x) //map[int]bool
	case "renameSQL":
		cast.Convert(&a.RenameSQL, x) //string
	case "renameGO":
		cast.Convert(&a.RenameGO_OK, x) //string
	case "renameJS":
		cast.Convert(&a.RenameJS, x) //string
	case "mast_upper_go":
		cast.Convert(&a.MAST_UPPER_GO, x) //string
	case "intToSmallInt":
		cast.Convert(&a.IntToSmallInt, x) //int
	case "skip":
		cast.Convert(&a.Skip, x) //string
	case "sql_unique_u1_1":
		cast.Convert(&a.Sql_unique_u1_1, x) //int
	case "sql_unique_u1_2":
		cast.Convert(&a.Sql_unique_u1_2, x) //int
	case "sql_index1_1":
		cast.Convert(&a.Sql_index1_1, x) //int
	case "sql_index1_2":
		cast.Convert(&a.Sql_index1_2, x) //int
	case "sql_index1_3":
		cast.Convert(&a.Sql_index1_3, x) //int
	case "sql_keys_1":
		cast.Convert(&a.Sql_keys_1, x) //int
	case "sql_keys_2":
		cast.Convert(&a.Sql_keys_2, x) //int
	case "sql_keys_3":
		cast.Convert(&a.Sql_keys_3, x) //int
	case "sql_search":
		cast.Convert(&a.Sql_search, x) //string
	case "sql_get":
		cast.Convert(&a.Sql_get, x) //string
	case "sql_unique_x1":
		cast.Convert(&a.Sql_unique_x1, x) //int
	case "sql_unique_x2":
		cast.Convert(&a.Sql_unique_x2, x) //int
	case "sql_unique_x1_x2":
		cast.Convert(&a.Sql_unique_x1_x2, x) //int
	case "sql_primary":
		cast.Convert(&a.Sql_primary, x) //float64
	case "sql_jsonb_index":
		cast.Convert(&a.Sql_jsonb_index, x) //map[string]any
	case "time_duration":
		cast.Convert(&a.Time_duration, x) //time.Duration
	case "public_field1":
		cast.Convert(&a.Public_field1, x) //int
	case "public_field2":
		cast.Convert(&a.Public_field2, x) //int
	case "public_field3":
		cast.Convert(&a.Public_field3, x) //int
	case "public_field_me1":
		cast.Convert(&a.Public_field_me1, x) //int
	case "public_field_me2":
		cast.Convert(&a.Public_field_me2, x) //int
	case "public_field_me3":
		cast.Convert(&a.Public_field_me3, x) //int
	}
}

// get struct value with function
func (a *News) Get(k string) (v any) {
	switch k {
	case "inc":
		return a.Inc //int
	case "ints":
		return a.Ints //int
	case "ints8":
		return a.Ints8 //int8
	case "ints16":
		return a.Ints16 //int16
	case "ints32":
		return a.Ints32 //int32
	case "ints64":
		return a.Ints64 //int64
	case "uints":
		return a.Uints //uint
	case "uints8":
		return a.Uints8 //uint8
	case "uints16":
		return a.Uints16 //uint16
	case "uints32":
		return a.Uints32 //uint32
	case "uints64":
		return a.Uints64 //uint64
	case "floats32":
		return a.Floats32 //float32
	case "floats64":
		return a.Floats64 //float64
	case "bools":
		return a.Bools //bool
	case "byte1":
		return a.Byte1 //byte
	case "bytes":
		return a.Bytes //[]byte
	case "list_ints":
		return a.List_ints //[]int
	case "list_string":
		return a.List_string //[]string
	case "list_float":
		return a.List_float //[]float64
	case "map_string_string":
		return a.Map_string_string //map[string]string
	case "map_string_bytes":
		return a.Map_string_bytes //map[string][]byte
	case "map_string_bool":
		return a.Map_string_bool //map[string]bool
	case "map_string_int":
		return a.Map_string_int //map[string]int
	case "map_string_float64":
		return a.Map_string_float64 //map[string]float64
	case "map_string_any":
		return a.Map_string_any //map[string]any
	case "map_int_string":
		return a.Map_int_string //map[int]string
	case "map_int_int":
		return a.Map_int_int //map[int]int
	case "map_int_bool":
		return a.Map_int_bool //map[int]bool
	case "renameSQL":
		return a.RenameSQL //string
	case "renameGO":
		return a.RenameGO_OK //string
	case "renameJS":
		return a.RenameJS //string
	case "mast_upper_go":
		return a.MAST_UPPER_GO //string
	case "intToSmallInt":
		return a.IntToSmallInt //int
	case "skip":
		return a.Skip //string
	case "sql_unique_u1_1":
		return a.Sql_unique_u1_1 //int
	case "sql_unique_u1_2":
		return a.Sql_unique_u1_2 //int
	case "sql_index1_1":
		return a.Sql_index1_1 //int
	case "sql_index1_2":
		return a.Sql_index1_2 //int
	case "sql_index1_3":
		return a.Sql_index1_3 //int
	case "sql_keys_1":
		return a.Sql_keys_1 //int
	case "sql_keys_2":
		return a.Sql_keys_2 //int
	case "sql_keys_3":
		return a.Sql_keys_3 //int
	case "sql_search":
		return a.Sql_search //string
	case "sql_get":
		return a.Sql_get //string
	case "sql_unique_x1":
		return a.Sql_unique_x1 //int
	case "sql_unique_x2":
		return a.Sql_unique_x2 //int
	case "sql_unique_x1_x2":
		return a.Sql_unique_x1_x2 //int
	case "sql_primary":
		return a.Sql_primary //float64
	case "sql_jsonb_index":
		return a.Sql_jsonb_index //map[string]any
	case "time_duration":
		return a.Time_duration //time.Duration
	case "go_type_int_to_strings":
		return a.Go_type_int_to_strings //int
	case "public_field1":
		return a.Public_field1 //int
	case "public_field2":
		return a.Public_field2 //int
	case "public_field3":
		return a.Public_field3 //int
	case "public_field_me1":
		return a.Public_field_me1 //int
	case "public_field_me2":
		return a.Public_field_me2 //int
	case "public_field_me3":
		return a.Public_field_me3 //int
	}
	return
}

// get any struct value as string
func (a *News) String(k string) (v string) {
	switch k {
	case "inc":
		return fmt.Sprint(a.Inc) //int
	case "ints":
		return fmt.Sprint(a.Ints) //int
	case "ints8":
		return fmt.Sprint(a.Ints8) //int8
	case "ints16":
		return fmt.Sprint(a.Ints16) //int16
	case "ints32":
		return fmt.Sprint(a.Ints32) //int32
	case "ints64":
		return fmt.Sprint(a.Ints64) //int64
	case "uints":
		return fmt.Sprint(a.Uints) //uint
	case "uints8":
		return fmt.Sprint(a.Uints8) //uint8
	case "uints16":
		return fmt.Sprint(a.Uints16) //uint16
	case "uints32":
		return fmt.Sprint(a.Uints32) //uint32
	case "uints64":
		return fmt.Sprint(a.Uints64) //uint64
	case "floats32":
		return fmt.Sprint(a.Floats32) //float32
	case "floats64":
		return fmt.Sprint(a.Floats64) //float64
	case "bools":
		return fmt.Sprint(a.Bools) //bool
	case "byte1":
		return fmt.Sprint(a.Byte1) //byte
	case "bytes":
		return fmt.Sprint(a.Bytes) //[]byte
	case "list_ints":
		return fmt.Sprint(a.List_ints) //[]int
	case "list_string":
		return fmt.Sprint(a.List_string) //[]string
	case "list_float":
		return fmt.Sprint(a.List_float) //[]float64
	case "map_string_string":
		return fmt.Sprint(a.Map_string_string) //map[string]string
	case "map_string_bytes":
		return fmt.Sprint(a.Map_string_bytes) //map[string][]byte
	case "map_string_bool":
		return fmt.Sprint(a.Map_string_bool) //map[string]bool
	case "map_string_int":
		return fmt.Sprint(a.Map_string_int) //map[string]int
	case "map_string_float64":
		return fmt.Sprint(a.Map_string_float64) //map[string]float64
	case "map_string_any":
		return fmt.Sprint(a.Map_string_any) //map[string]any
	case "map_int_string":
		return fmt.Sprint(a.Map_int_string) //map[int]string
	case "map_int_int":
		return fmt.Sprint(a.Map_int_int) //map[int]int
	case "map_int_bool":
		return fmt.Sprint(a.Map_int_bool) //map[int]bool
	case "renameSQL":
		return fmt.Sprint(a.RenameSQL) //string
	case "renameGO":
		return fmt.Sprint(a.RenameGO_OK) //string
	case "renameJS":
		return fmt.Sprint(a.RenameJS) //string
	case "mast_upper_go":
		return fmt.Sprint(a.MAST_UPPER_GO) //string
	case "intToSmallInt":
		return fmt.Sprint(a.IntToSmallInt) //int
	case "skip":
		return fmt.Sprint(a.Skip) //string
	case "sql_unique_u1_1":
		return fmt.Sprint(a.Sql_unique_u1_1) //int
	case "sql_unique_u1_2":
		return fmt.Sprint(a.Sql_unique_u1_2) //int
	case "sql_index1_1":
		return fmt.Sprint(a.Sql_index1_1) //int
	case "sql_index1_2":
		return fmt.Sprint(a.Sql_index1_2) //int
	case "sql_index1_3":
		return fmt.Sprint(a.Sql_index1_3) //int
	case "sql_keys_1":
		return fmt.Sprint(a.Sql_keys_1) //int
	case "sql_keys_2":
		return fmt.Sprint(a.Sql_keys_2) //int
	case "sql_keys_3":
		return fmt.Sprint(a.Sql_keys_3) //int
	case "sql_search":
		return fmt.Sprint(a.Sql_search) //string
	case "sql_get":
		return fmt.Sprint(a.Sql_get) //string
	case "sql_unique_x1":
		return fmt.Sprint(a.Sql_unique_x1) //int
	case "sql_unique_x2":
		return fmt.Sprint(a.Sql_unique_x2) //int
	case "sql_unique_x1_x2":
		return fmt.Sprint(a.Sql_unique_x1_x2) //int
	case "sql_primary":
		return fmt.Sprint(a.Sql_primary) //float64
	case "sql_jsonb_index":
		return fmt.Sprint(a.Sql_jsonb_index) //map[string]any
	case "time_duration":
		return fmt.Sprint(a.Time_duration) //time.Duration
	case "go_type_int_to_strings":
		return fmt.Sprint(a.Go_type_int_to_strings) //int
	case "public_field1":
		return fmt.Sprint(a.Public_field1) //int
	case "public_field2":
		return fmt.Sprint(a.Public_field2) //int
	case "public_field3":
		return fmt.Sprint(a.Public_field3) //int
	case "public_field_me1":
		return fmt.Sprint(a.Public_field_me1) //int
	case "public_field_me2":
		return fmt.Sprint(a.Public_field_me2) //int
	case "public_field_me3":
		return fmt.Sprint(a.Public_field_me3) //int
	}
	return
}

// Struct to json
func (a *News) ToJson() (r []byte) {
	js := jsons.Create().
		Add(FieldInc, a.Inc).
		Add(FieldInts, a.Ints).
		Add(FieldInts8, a.Ints8).
		Add(FieldInts16, a.Ints16).
		Add(FieldInts32, a.Ints32).
		Add(FieldInts64, a.Ints64).
		Add(FieldUints, a.Uints).
		Add(FieldUints8, a.Uints8).
		Add(FieldUints16, a.Uints16).
		Add(FieldUints32, a.Uints32).
		Add(FieldUints64, a.Uints64).
		Add(FieldFloats32, a.Floats32).
		Add(FieldFloats64, a.Floats64).
		Add(FieldBools, a.Bools).
		Add(FieldByte1, a.Byte1).
		Add(FieldBytes, a.Bytes).
		Add(FieldList_ints, a.List_ints).
		Add(FieldList_string, a.List_string).
		Add(FieldList_float, a.List_float).
		Add(FieldMap_string_string, a.Map_string_string).
		Add(FieldMap_string_bytes, a.Map_string_bytes).
		Add(FieldMap_string_bool, a.Map_string_bool).
		Add(FieldMap_string_int, a.Map_string_int).
		Add(FieldMap_string_float64, a.Map_string_float64).
		Add(FieldMap_string_any, a.Map_string_any).
		Add(FieldMap_int_string, a.Map_int_string).
		Add(FieldMap_int_int, a.Map_int_int).
		Add(FieldMap_int_bool, a.Map_int_bool).
		Add(FieldRenameSQL, a.RenameSQL).
		Add(FieldRenameGO_OK, a.RenameGO_OK).
		Add(FieldRenameJS, a.RenameJS).
		Add(FieldMAST_UPPER_GO, a.MAST_UPPER_GO).
		Add(FieldIntToSmallInt, a.IntToSmallInt).
		Add(FieldSkip, a.Skip).
		Add(FieldSql_unique_u1_1, a.Sql_unique_u1_1).
		Add(FieldSql_unique_u1_2, a.Sql_unique_u1_2).
		Add(FieldSql_index1_1, a.Sql_index1_1).
		Add(FieldSql_index1_2, a.Sql_index1_2).
		Add(FieldSql_index1_3, a.Sql_index1_3).
		Add(FieldSql_keys_1, a.Sql_keys_1).
		Add(FieldSql_keys_2, a.Sql_keys_2).
		Add(FieldSql_keys_3, a.Sql_keys_3).
		Add(FieldSql_search, a.Sql_search).
		Add(FieldSql_get, a.Sql_get).
		Add(FieldSql_unique_x1, a.Sql_unique_x1).
		Add(FieldSql_unique_x2, a.Sql_unique_x2).
		Add(FieldSql_unique_x1_x2, a.Sql_unique_x1_x2).
		Add(FieldSql_primary, a.Sql_primary).
		Add(FieldSql_jsonb_index, a.Sql_jsonb_index).
		Add(FieldTime_duration, a.Time_duration).
		Add(FieldGo_type_int_to_strings, a.Go_type_int_to_strings).
		Add(FieldPublic_field1, a.Public_field1).
		Add(FieldPublic_field2, a.Public_field2).
		Add(FieldPublic_field3, a.Public_field3).
		Add(FieldPublic_field_me1, a.Public_field_me1).
		Add(FieldPublic_field_me2, a.Public_field_me2).
		Add(FieldPublic_field_me3, a.Public_field_me3)
	return js.Bytes()
}

func NewsReadonlyList() []NewsIndexType {
	return []NewsIndexType{IndexInc}
}

func NewsPublicList() []NewsIndexType {
	return []NewsIndexType{IndexPublic_field1, IndexPublic_field2, IndexPublic_field3}
}

func NewsMeList() []NewsIndexType {
	return []NewsIndexType{IndexPublic_field_me1, IndexPublic_field_me2, IndexPublic_field_me3}
}

func (a NewsIndexType) Readonly() bool {
	switch a {
	case IndexInc:
		return true
	default:
		return false
	}
}
func (a NewsIndexType) Public() bool {
	switch a {
	case IndexPublic_field1, IndexPublic_field2, IndexPublic_field3:
		return true
	default:
		return false
	}
}
func (a NewsIndexType) Me() bool {
	switch a {
	case IndexPublic_field_me1, IndexPublic_field_me2, IndexPublic_field_me3:
		return true
	default:
		return false
	}
}

// key index string
func (a NewsIndexType) String() string {
	switch a {
	case IndexInc:
		return "inc"
	case IndexInts:
		return "ints"
	case IndexInts8:
		return "ints8"
	case IndexInts16:
		return "ints16"
	case IndexInts32:
		return "ints32"
	case IndexInts64:
		return "ints64"
	case IndexUints:
		return "uints"
	case IndexUints8:
		return "uints8"
	case IndexUints16:
		return "uints16"
	case IndexUints32:
		return "uints32"
	case IndexUints64:
		return "uints64"
	case IndexFloats32:
		return "floats32"
	case IndexFloats64:
		return "floats64"
	case IndexBools:
		return "bools"
	case IndexByte1:
		return "byte1"
	case IndexBytes:
		return "bytes"
	case IndexList_ints:
		return "list_ints"
	case IndexList_string:
		return "list_string"
	case IndexList_float:
		return "list_float"
	case IndexMap_string_string:
		return "map_string_string"
	case IndexMap_string_bytes:
		return "map_string_bytes"
	case IndexMap_string_bool:
		return "map_string_bool"
	case IndexMap_string_int:
		return "map_string_int"
	case IndexMap_string_float64:
		return "map_string_float64"
	case IndexMap_string_any:
		return "map_string_any"
	case IndexMap_int_string:
		return "map_int_string"
	case IndexMap_int_int:
		return "map_int_int"
	case IndexMap_int_bool:
		return "map_int_bool"
	case IndexRenameSQL:
		return "renameSQL"
	case IndexRenameGO_OK:
		return "renameGO"
	case IndexRenameJS:
		return "renameJS"
	case IndexMAST_UPPER_GO:
		return "mast_upper_go"
	case IndexIntToSmallInt:
		return "intToSmallInt"
	case IndexSkip:
		return "skip"
	case IndexSql_unique_u1_1:
		return "sql_unique_u1_1"
	case IndexSql_unique_u1_2:
		return "sql_unique_u1_2"
	case IndexSql_index1_1:
		return "sql_index1_1"
	case IndexSql_index1_2:
		return "sql_index1_2"
	case IndexSql_index1_3:
		return "sql_index1_3"
	case IndexSql_keys_1:
		return "sql_keys_1"
	case IndexSql_keys_2:
		return "sql_keys_2"
	case IndexSql_keys_3:
		return "sql_keys_3"
	case IndexSql_search:
		return "sql_search"
	case IndexSql_get:
		return "sql_get"
	case IndexSql_unique_x1:
		return "sql_unique_x1"
	case IndexSql_unique_x2:
		return "sql_unique_x2"
	case IndexSql_unique_x1_x2:
		return "sql_unique_x1_x2"
	case IndexSql_primary:
		return "sql_primary"
	case IndexSql_jsonb_index:
		return "sql_jsonb_index"
	case IndexTime_duration:
		return "time_duration"
	case IndexGo_type_int_to_strings:
		return "go_type_int_to_strings"
	case IndexPublic_field1:
		return "public_field1"
	case IndexPublic_field2:
		return "public_field2"
	case IndexPublic_field3:
		return "public_field3"
	case IndexPublic_field_me1:
		return "public_field_me1"
	case IndexPublic_field_me2:
		return "public_field_me2"
	case IndexPublic_field_me3:
		return "public_field_me3"
	default:
		return ""
	}
}

// key index string
func (a NewsIndexType) SQLName() string {
	switch a {
	case IndexInc:
		return "inc"
	case IndexInts:
		return "ints"
	case IndexInts8:
		return "ints8"
	case IndexInts16:
		return "ints16"
	case IndexInts32:
		return "ints32"
	case IndexInts64:
		return "ints64"
	case IndexUints:
		return "uints"
	case IndexUints8:
		return "uints8"
	case IndexUints16:
		return "uints16"
	case IndexUints32:
		return "uints32"
	case IndexUints64:
		return "uints64"
	case IndexFloats32:
		return "floats32"
	case IndexFloats64:
		return "floats64"
	case IndexBools:
		return "bools"
	case IndexByte1:
		return "byte1"
	case IndexBytes:
		return "bytes"
	case IndexList_ints:
		return "list_ints"
	case IndexList_string:
		return "list_string"
	case IndexList_float:
		return "list_float"
	case IndexMap_string_string:
		return "map_string_string"
	case IndexMap_string_bytes:
		return "map_string_bytes"
	case IndexMap_string_bool:
		return "map_string_bool"
	case IndexMap_string_int:
		return "map_string_int"
	case IndexMap_string_float64:
		return "map_string_float64"
	case IndexMap_string_any:
		return "map_string_any"
	case IndexMap_int_string:
		return "map_int_string"
	case IndexMap_int_int:
		return "map_int_int"
	case IndexMap_int_bool:
		return "map_int_bool"
	case IndexRenameSQL:
		return "renameSQL_OK"
	case IndexRenameGO_OK:
		return "renameGO"
	case IndexRenameJS:
		return "renameJS"
	case IndexMAST_UPPER_GO:
		return "mast_upper_go"
	case IndexIntToSmallInt:
		return "intToSmallInt"
	case IndexSkip:
		return "skip"
	case IndexSql_unique_u1_1:
		return "sql_unique_u1_1"
	case IndexSql_unique_u1_2:
		return "sql_unique_u1_2"
	case IndexSql_index1_1:
		return "sql_index1_1"
	case IndexSql_index1_2:
		return "sql_index1_2"
	case IndexSql_index1_3:
		return "sql_index1_3"
	case IndexSql_keys_1:
		return "sql_keys_1"
	case IndexSql_keys_2:
		return "sql_keys_2"
	case IndexSql_keys_3:
		return "sql_keys_3"
	case IndexSql_search:
		return "sql_search"
	case IndexSql_get:
		return "sql_get"
	case IndexSql_unique_x1:
		return "sql_unique_x1"
	case IndexSql_unique_x2:
		return "sql_unique_x2"
	case IndexSql_unique_x1_x2:
		return "sql_unique_x1_x2"
	case IndexSql_primary:
		return "sql_primary"
	case IndexSql_jsonb_index:
		return "sql_jsonb_index"
	case IndexTime_duration:
		return "time_duration"
	case IndexGo_type_int_to_strings:
		return "go_type_int_to_strings"
	case IndexPublic_field1:
		return "public_field1"
	case IndexPublic_field2:
		return "public_field2"
	case IndexPublic_field3:
		return "public_field3"
	case IndexPublic_field_me1:
		return "public_field_me1"
	case IndexPublic_field_me2:
		return "public_field_me2"
	case IndexPublic_field_me3:
		return "public_field_me3"
	default:
		return ""
	}
}

// key index clickhouse string
func (a NewsIndexType) ClickhouseName() string {
	switch a {
	case IndexInc:
		return "inc"
	case IndexInts:
		return "ints"
	case IndexInts8:
		return "ints8"
	case IndexInts16:
		return "ints16"
	case IndexInts32:
		return "ints32"
	case IndexInts64:
		return "ints64"
	case IndexUints:
		return "uints"
	case IndexUints8:
		return "uints8"
	case IndexUints16:
		return "uints16"
	case IndexUints32:
		return "uints32"
	case IndexUints64:
		return "uints64"
	case IndexFloats32:
		return "floats32"
	case IndexFloats64:
		return "floats64"
	case IndexBools:
		return "bools"
	case IndexByte1:
		return "byte1"
	case IndexBytes:
		return "bytes"
	case IndexList_ints:
		return "list_ints"
	case IndexList_string:
		return "list_string"
	case IndexList_float:
		return "list_float"
	case IndexMap_string_string:
		return "map_string_string"
	case IndexMap_string_bytes:
		return "map_string_bytes"
	case IndexMap_string_bool:
		return "map_string_bool"
	case IndexMap_string_int:
		return "map_string_int"
	case IndexMap_string_float64:
		return "map_string_float64"
	case IndexMap_string_any:
		return "map_string_any"
	case IndexMap_int_string:
		return "map_int_string"
	case IndexMap_int_int:
		return "map_int_int"
	case IndexMap_int_bool:
		return "map_int_bool"
	case IndexRenameSQL:
		return "renameSQL"
	case IndexRenameGO_OK:
		return "renameGO"
	case IndexRenameJS:
		return "renameJS"
	case IndexMAST_UPPER_GO:
		return "mast_upper_go"
	case IndexIntToSmallInt:
		return "intToSmallInt"
	case IndexSkip:
		return "skip"
	case IndexSql_unique_u1_1:
		return "sql_unique_u1_1"
	case IndexSql_unique_u1_2:
		return "sql_unique_u1_2"
	case IndexSql_index1_1:
		return "sql_index1_1"
	case IndexSql_index1_2:
		return "sql_index1_2"
	case IndexSql_index1_3:
		return "sql_index1_3"
	case IndexSql_keys_1:
		return "sql_keys_1"
	case IndexSql_keys_2:
		return "sql_keys_2"
	case IndexSql_keys_3:
		return "sql_keys_3"
	case IndexSql_search:
		return "sql_search"
	case IndexSql_get:
		return "sql_get"
	case IndexSql_unique_x1:
		return "sql_unique_x1"
	case IndexSql_unique_x2:
		return "sql_unique_x2"
	case IndexSql_unique_x1_x2:
		return "sql_unique_x1_x2"
	case IndexSql_primary:
		return "sql_primary"
	case IndexSql_jsonb_index:
		return "sql_jsonb_index"
	case IndexTime_duration:
		return "time_duration"
	case IndexGo_type_int_to_strings:
		return "go_type_int_to_strings"
	case IndexPublic_field1:
		return "public_field1"
	case IndexPublic_field2:
		return "public_field2"
	case IndexPublic_field3:
		return "public_field3"
	case IndexPublic_field_me1:
		return "public_field_me1"
	case IndexPublic_field_me2:
		return "public_field_me2"
	case IndexPublic_field_me3:
		return "public_field_me3"
	default:
		return ""
	}
}

// key index type
func (a NewsIndexType) Type() string {
	switch a {
	case IndexInc:
		return "int"
	case IndexInts:
		return "int"
	case IndexInts8:
		return "int8"
	case IndexInts16:
		return "int16"
	case IndexInts32:
		return "int32"
	case IndexInts64:
		return "int64"
	case IndexUints:
		return "uint"
	case IndexUints8:
		return "uint8"
	case IndexUints16:
		return "uint16"
	case IndexUints32:
		return "uint32"
	case IndexUints64:
		return "uint64"
	case IndexFloats32:
		return "float32"
	case IndexFloats64:
		return "float64"
	case IndexBools:
		return "bool"
	case IndexByte1:
		return "byte"
	case IndexBytes:
		return "[]byte"
	case IndexList_ints:
		return "[]int"
	case IndexList_string:
		return "[]string"
	case IndexList_float:
		return "[]float64"
	case IndexMap_string_string:
		return "map[string]string"
	case IndexMap_string_bytes:
		return "map[string][]byte"
	case IndexMap_string_bool:
		return "map[string]bool"
	case IndexMap_string_int:
		return "map[string]int"
	case IndexMap_string_float64:
		return "map[string]float64"
	case IndexMap_string_any:
		return "map[string]any"
	case IndexMap_int_string:
		return "map[int]string"
	case IndexMap_int_int:
		return "map[int]int"
	case IndexMap_int_bool:
		return "map[int]bool"
	case IndexRenameSQL:
		return "string"
	case IndexRenameGO_OK:
		return "string"
	case IndexRenameJS:
		return "string"
	case IndexMAST_UPPER_GO:
		return "string"
	case IndexIntToSmallInt:
		return "int"
	case IndexSkip:
		return "string"
	case IndexSql_unique_u1_1:
		return "int"
	case IndexSql_unique_u1_2:
		return "int"
	case IndexSql_index1_1:
		return "int"
	case IndexSql_index1_2:
		return "int"
	case IndexSql_index1_3:
		return "int"
	case IndexSql_keys_1:
		return "int"
	case IndexSql_keys_2:
		return "int"
	case IndexSql_keys_3:
		return "int"
	case IndexSql_search:
		return "string"
	case IndexSql_get:
		return "string"
	case IndexSql_unique_x1:
		return "int"
	case IndexSql_unique_x2:
		return "int"
	case IndexSql_unique_x1_x2:
		return "int"
	case IndexSql_primary:
		return "float64"
	case IndexSql_jsonb_index:
		return "map[string]any"
	case IndexTime_duration:
		return "time.Duration"
	case IndexGo_type_int_to_strings:
		return "int"
	case IndexPublic_field1:
		return "int"
	case IndexPublic_field2:
		return "int"
	case IndexPublic_field3:
		return "int"
	case IndexPublic_field_me1:
		return "int"
	case IndexPublic_field_me2:
		return "int"
	case IndexPublic_field_me3:
		return "int"
	default:
		return ""
	}
}

// custom title
func (a NewsIndexType) Title() string {
	switch a {
	case IndexInc:
		return "Inc"
	case IndexInts:
		return "Ints"
	case IndexInts8:
		return "Ints8"
	case IndexInts16:
		return "Ints16"
	case IndexInts32:
		return "Ints32"
	case IndexInts64:
		return "Ints64"
	case IndexUints:
		return "Uints"
	case IndexUints8:
		return "Uints8"
	case IndexUints16:
		return "Uints16"
	case IndexUints32:
		return "Uints32"
	case IndexUints64:
		return "Uints64"
	case IndexFloats32:
		return "Floats32"
	case IndexFloats64:
		return "Floats64"
	case IndexBools:
		return "Bools"
	case IndexByte1:
		return "Byte1"
	case IndexBytes:
		return "Bytes"
	case IndexList_ints:
		return "List_ints"
	case IndexList_string:
		return "List_string"
	case IndexList_float:
		return "List_float"
	case IndexMap_string_string:
		return "Map_string_string"
	case IndexMap_string_bytes:
		return "Map_string_bytes"
	case IndexMap_string_bool:
		return "Map_string_bool"
	case IndexMap_string_int:
		return "Map_string_int"
	case IndexMap_string_float64:
		return "Map_string_float64"
	case IndexMap_string_any:
		return "Map_string_any"
	case IndexMap_int_string:
		return "Map_int_string"
	case IndexMap_int_int:
		return "Map_int_int"
	case IndexMap_int_bool:
		return "Map_int_bool"
	case IndexRenameSQL:
		return "RenameSQL"
	case IndexRenameGO_OK:
		return "RenameGO_OK"
	case IndexRenameJS:
		return "RenameJS"
	case IndexMAST_UPPER_GO:
		return "MAST_UPPER_GO"
	case IndexIntToSmallInt:
		return "IntToSmallInt"
	case IndexSkip:
		return "Skip"
	case IndexSql_unique_u1_1:
		return "Sql_unique_u1_1"
	case IndexSql_unique_u1_2:
		return "Sql_unique_u1_2"
	case IndexSql_index1_1:
		return "Sql_index1_1"
	case IndexSql_index1_2:
		return "Sql_index1_2"
	case IndexSql_index1_3:
		return "Sql_index1_3"
	case IndexSql_keys_1:
		return "Sql_keys_1"
	case IndexSql_keys_2:
		return "Sql_keys_2"
	case IndexSql_keys_3:
		return "Sql_keys_3"
	case IndexSql_search:
		return "Sql_search"
	case IndexSql_get:
		return "Sql_get"
	case IndexSql_unique_x1:
		return "Sql_unique_x1"
	case IndexSql_unique_x2:
		return "Sql_unique_x2"
	case IndexSql_unique_x1_x2:
		return "Sql_unique_x1_x2"
	case IndexSql_primary:
		return "Sql_primary"
	case IndexSql_jsonb_index:
		return "Sql_jsonb_index"
	case IndexTime_duration:
		return "Time_duration"
	case IndexGo_type_int_to_strings:
		return "Go_type_int_to_strings"
	case IndexPublic_field1:
		return "Public_field1"
	case IndexPublic_field2:
		return "Public_field2"
	case IndexPublic_field3:
		return "Public_field3"
	case IndexPublic_field_me1:
		return "Public_field_me1"
	case IndexPublic_field_me2:
		return "Public_field_me2"
	case IndexPublic_field_me3:
		return "Public_field_me3"
	default:
		return ""
	}
}

// custom desc
func (a NewsIndexType) Desc() string {
	switch a {
	default:
		return ""
	}
}

// struct key to index
func NewsKeyIndex(key string) NewsIndexType {
	switch key {
	case "inc":
		return IndexInc
	case "ints":
		return IndexInts
	case "ints8":
		return IndexInts8
	case "ints16":
		return IndexInts16
	case "ints32":
		return IndexInts32
	case "ints64":
		return IndexInts64
	case "uints":
		return IndexUints
	case "uints8":
		return IndexUints8
	case "uints16":
		return IndexUints16
	case "uints32":
		return IndexUints32
	case "uints64":
		return IndexUints64
	case "floats32":
		return IndexFloats32
	case "floats64":
		return IndexFloats64
	case "bools":
		return IndexBools
	case "byte1":
		return IndexByte1
	case "bytes":
		return IndexBytes
	case "list_ints":
		return IndexList_ints
	case "list_string":
		return IndexList_string
	case "list_float":
		return IndexList_float
	case "map_string_string":
		return IndexMap_string_string
	case "map_string_bytes":
		return IndexMap_string_bytes
	case "map_string_bool":
		return IndexMap_string_bool
	case "map_string_int":
		return IndexMap_string_int
	case "map_string_float64":
		return IndexMap_string_float64
	case "map_string_any":
		return IndexMap_string_any
	case "map_int_string":
		return IndexMap_int_string
	case "map_int_int":
		return IndexMap_int_int
	case "map_int_bool":
		return IndexMap_int_bool
	case "renameSQL":
		return IndexRenameSQL
	case "renameGO":
		return IndexRenameGO_OK
	case "renameJS":
		return IndexRenameJS
	case "mast_upper_go":
		return IndexMAST_UPPER_GO
	case "intToSmallInt":
		return IndexIntToSmallInt
	case "skip":
		return IndexSkip
	case "sql_unique_u1_1":
		return IndexSql_unique_u1_1
	case "sql_unique_u1_2":
		return IndexSql_unique_u1_2
	case "sql_index1_1":
		return IndexSql_index1_1
	case "sql_index1_2":
		return IndexSql_index1_2
	case "sql_index1_3":
		return IndexSql_index1_3
	case "sql_keys_1":
		return IndexSql_keys_1
	case "sql_keys_2":
		return IndexSql_keys_2
	case "sql_keys_3":
		return IndexSql_keys_3
	case "sql_search":
		return IndexSql_search
	case "sql_get":
		return IndexSql_get
	case "sql_unique_x1":
		return IndexSql_unique_x1
	case "sql_unique_x2":
		return IndexSql_unique_x2
	case "sql_unique_x1_x2":
		return IndexSql_unique_x1_x2
	case "sql_primary":
		return IndexSql_primary
	case "sql_jsonb_index":
		return IndexSql_jsonb_index
	case "time_duration":
		return IndexTime_duration
	case "go_type_int_to_strings":
		return IndexGo_type_int_to_strings
	case "public_field1":
		return IndexPublic_field1
	case "public_field2":
		return IndexPublic_field2
	case "public_field3":
		return IndexPublic_field3
	case "public_field_me1":
		return IndexPublic_field_me1
	case "public_field_me2":
		return IndexPublic_field_me2
	case "public_field_me3":
		return IndexPublic_field_me3
	default:
		return 0
	}
}

// valid struct key check
func NewsValidKey(key string) bool {
	switch key {
	case "inc", "ints", "ints8", "ints16", "ints32", "ints64", "uints", "uints8", "uints16", "uints32", "uints64", "floats32", "floats64", "bools", "byte1", "bytes", "list_ints", "list_string", "list_float", "map_string_string", "map_string_bytes", "map_string_bool", "map_string_int", "map_string_float64", "map_string_any", "map_int_string", "map_int_int", "map_int_bool", "renameSQL", "renameGO", "renameJS", "mast_upper_go", "intToSmallInt", "skip", "sql_unique_u1_1", "sql_unique_u1_2", "sql_index1_1", "sql_index1_2", "sql_index1_3", "sql_keys_1", "sql_keys_2", "sql_keys_3", "sql_search", "sql_get", "sql_unique_x1", "sql_unique_x2", "sql_unique_x1_x2", "sql_primary", "sql_jsonb_index", "time_duration", "go_type_int_to_strings", "public_field1", "public_field2", "public_field3", "public_field_me1", "public_field_me2", "public_field_me3":
		return true
	default:
		return false
	}
}

// struct to map
func (a *News) Map() map[string]any {
	return map[string]any{
		"inc":                    a.Inc,
		"ints":                   a.Ints,
		"ints8":                  a.Ints8,
		"ints16":                 a.Ints16,
		"ints32":                 a.Ints32,
		"ints64":                 a.Ints64,
		"uints":                  a.Uints,
		"uints8":                 a.Uints8,
		"uints16":                a.Uints16,
		"uints32":                a.Uints32,
		"uints64":                a.Uints64,
		"floats32":               a.Floats32,
		"floats64":               a.Floats64,
		"bools":                  a.Bools,
		"byte1":                  a.Byte1,
		"bytes":                  a.Bytes,
		"list_ints":              a.List_ints,
		"list_string":            a.List_string,
		"list_float":             a.List_float,
		"map_string_string":      a.Map_string_string,
		"map_string_bytes":       a.Map_string_bytes,
		"map_string_bool":        a.Map_string_bool,
		"map_string_int":         a.Map_string_int,
		"map_string_float64":     a.Map_string_float64,
		"map_string_any":         a.Map_string_any,
		"map_int_string":         a.Map_int_string,
		"map_int_int":            a.Map_int_int,
		"map_int_bool":           a.Map_int_bool,
		"renameSQL":              a.RenameSQL,
		"renameGO":               a.RenameGO_OK,
		"renameJS_OK":            a.RenameJS,
		"mast_upper_go":          a.MAST_UPPER_GO,
		"intToSmallInt":          a.IntToSmallInt,
		"skip":                   a.Skip,
		"sql_unique_u1_1":        a.Sql_unique_u1_1,
		"sql_unique_u1_2":        a.Sql_unique_u1_2,
		"sql_index1_1":           a.Sql_index1_1,
		"sql_index1_2":           a.Sql_index1_2,
		"sql_index1_3":           a.Sql_index1_3,
		"sql_keys_1":             a.Sql_keys_1,
		"sql_keys_2":             a.Sql_keys_2,
		"sql_keys_3":             a.Sql_keys_3,
		"sql_search":             a.Sql_search,
		"sql_get":                a.Sql_get,
		"sql_unique_x1":          a.Sql_unique_x1,
		"sql_unique_x2":          a.Sql_unique_x2,
		"sql_unique_x1_x2":       a.Sql_unique_x1_x2,
		"sql_primary":            a.Sql_primary,
		"sql_jsonb_index":        a.Sql_jsonb_index,
		"time_duration":          a.Time_duration,
		"go_type_int_to_strings": a.Go_type_int_to_strings,
		"public_field1":          a.Public_field1,
		"public_field2":          a.Public_field2,
		"public_field3":          a.Public_field3,
		"public_field_me1":       a.Public_field_me1,
		"public_field_me2":       a.Public_field_me2,
		"public_field_me3":       a.Public_field_me3,
	}
}

// struct to map
func (a *News) Iterate(f func(k NewsIndexType, v any)) {
	for _, x := range NewsIndexes() {
		f(x, a.Get(x.String()))
	}
}

// gotiny marshal
func (a *News) MarshalGotiny() []byte {
	return gotiny.Marshal(&a.Inc, &a.Ints, &a.Ints8, &a.Ints16, &a.Ints32, &a.Ints64, &a.Uints, &a.Uints8, &a.Uints16, &a.Uints32, &a.Uints64, &a.Floats32, &a.Floats64, &a.Bools, &a.Byte1, &a.Bytes, &a.List_ints, &a.List_string, &a.List_float, &a.Map_string_string, &a.Map_string_bytes, &a.Map_string_bool, &a.Map_string_int, &a.Map_string_float64, &a.Map_string_any, &a.Map_int_string, &a.Map_int_int, &a.Map_int_bool, &a.RenameSQL, &a.RenameGO_OK, &a.RenameJS, &a.MAST_UPPER_GO, &a.IntToSmallInt, &a.Skip, &a.Sql_unique_u1_1, &a.Sql_unique_u1_2, &a.Sql_index1_1, &a.Sql_index1_2, &a.Sql_index1_3, &a.Sql_keys_1, &a.Sql_keys_2, &a.Sql_keys_3, &a.Sql_search, &a.Sql_get, &a.Sql_unique_x1, &a.Sql_unique_x2, &a.Sql_unique_x1_x2, &a.Sql_primary, &a.Sql_jsonb_index, &a.Time_duration, &a.Go_type_int_to_strings, &a.Public_field1, &a.Public_field2, &a.Public_field3, &a.Public_field_me1, &a.Public_field_me2, &a.Public_field_me3)
}

// parse gotiny
func UnmarshalNewsGotiny(v []byte) (a News) {
	gotiny.Unmarshal(v, &a.Inc, &a.Ints, &a.Ints8, &a.Ints16, &a.Ints32, &a.Ints64, &a.Uints, &a.Uints8, &a.Uints16, &a.Uints32, &a.Uints64, &a.Floats32, &a.Floats64, &a.Bools, &a.Byte1, &a.Bytes, &a.List_ints, &a.List_string, &a.List_float, &a.Map_string_string, &a.Map_string_bytes, &a.Map_string_bool, &a.Map_string_int, &a.Map_string_float64, &a.Map_string_any, &a.Map_int_string, &a.Map_int_int, &a.Map_int_bool, &a.RenameSQL, &a.RenameGO_OK, &a.RenameJS, &a.MAST_UPPER_GO, &a.IntToSmallInt, &a.Skip, &a.Sql_unique_u1_1, &a.Sql_unique_u1_2, &a.Sql_index1_1, &a.Sql_index1_2, &a.Sql_index1_3, &a.Sql_keys_1, &a.Sql_keys_2, &a.Sql_keys_3, &a.Sql_search, &a.Sql_get, &a.Sql_unique_x1, &a.Sql_unique_x2, &a.Sql_unique_x1_x2, &a.Sql_primary, &a.Sql_jsonb_index, &a.Time_duration, &a.Go_type_int_to_strings, &a.Public_field1, &a.Public_field2, &a.Public_field3, &a.Public_field_me1, &a.Public_field_me2, &a.Public_field_me3)
	return
}

// msgp marshal
func (a *News) MessagePack() []byte {
	b, _ := msgpack.Marshal(a)
	return b
}

// msgp unmarshal
func ParseNewsMessagePack(v []byte) (a News, err error) {
	err = msgpack.Unmarshal(v, &a)
	return
}

// fast json marshal
func (a *News) Pack() []byte {
	var jsoner = jsoniter.ConfigCompatibleWithStandardLibrary
	b, _ := jsoner.Marshal(a)
	return b
}

// fast json unmarshal
func ParseNews(v []byte) (a News, err error) {
	var jsoner = jsoniter.ConfigCompatibleWithStandardLibrary
	err = jsoner.Unmarshal(v, &a)
	return
}

// Parse []any to json
func ParseTupleToNewsJson(r []any) (a NewsJson) {
	for pos, x := range r {
		switch NewsIndexType(pos) {
		case IndexInc:
			a.Set(FieldInc, x)
		case IndexInts:
			a.Set(FieldInts, x)
		case IndexInts8:
			a.Set(FieldInts8, x)
		case IndexInts16:
			a.Set(FieldInts16, x)
		case IndexInts32:
			a.Set(FieldInts32, x)
		case IndexInts64:
			a.Set(FieldInts64, x)
		case IndexUints:
			a.Set(FieldUints, x)
		case IndexUints8:
			a.Set(FieldUints8, x)
		case IndexUints16:
			a.Set(FieldUints16, x)
		case IndexUints32:
			a.Set(FieldUints32, x)
		case IndexUints64:
			a.Set(FieldUints64, x)
		case IndexFloats32:
			a.Set(FieldFloats32, x)
		case IndexFloats64:
			a.Set(FieldFloats64, x)
		case IndexBools:
			a.Set(FieldBools, x)
		case IndexByte1:
			a.Set(FieldByte1, x)
		case IndexBytes:
			a.Set(FieldBytes, x)
		case IndexList_ints:
			a.Set(FieldList_ints, x)
		case IndexList_string:
			a.Set(FieldList_string, x)
		case IndexList_float:
			a.Set(FieldList_float, x)
		case IndexMap_string_string:
			a.Set(FieldMap_string_string, x)
		case IndexMap_string_bytes:
			a.Set(FieldMap_string_bytes, x)
		case IndexMap_string_bool:
			a.Set(FieldMap_string_bool, x)
		case IndexMap_string_int:
			a.Set(FieldMap_string_int, x)
		case IndexMap_string_float64:
			a.Set(FieldMap_string_float64, x)
		case IndexMap_string_any:
			a.Set(FieldMap_string_any, x)
		case IndexMap_int_string:
			a.Set(FieldMap_int_string, x)
		case IndexMap_int_int:
			a.Set(FieldMap_int_int, x)
		case IndexMap_int_bool:
			a.Set(FieldMap_int_bool, x)
		case IndexRenameSQL:
			a.Set(FieldRenameSQL, x)
		case IndexRenameGO_OK:
			a.Set(FieldRenameGO_OK, x)
		case IndexRenameJS:
			a.Set(FieldRenameJS, x)
		case IndexMAST_UPPER_GO:
			a.Set(FieldMAST_UPPER_GO, x)
		case IndexIntToSmallInt:
			a.Set(FieldIntToSmallInt, x)
		case IndexSkip:
			a.Set(FieldSkip, x)
		case IndexSql_unique_u1_1:
			a.Set(FieldSql_unique_u1_1, x)
		case IndexSql_unique_u1_2:
			a.Set(FieldSql_unique_u1_2, x)
		case IndexSql_index1_1:
			a.Set(FieldSql_index1_1, x)
		case IndexSql_index1_2:
			a.Set(FieldSql_index1_2, x)
		case IndexSql_index1_3:
			a.Set(FieldSql_index1_3, x)
		case IndexSql_keys_1:
			a.Set(FieldSql_keys_1, x)
		case IndexSql_keys_2:
			a.Set(FieldSql_keys_2, x)
		case IndexSql_keys_3:
			a.Set(FieldSql_keys_3, x)
		case IndexSql_search:
			a.Set(FieldSql_search, x)
		case IndexSql_get:
			a.Set(FieldSql_get, x)
		case IndexSql_unique_x1:
			a.Set(FieldSql_unique_x1, x)
		case IndexSql_unique_x2:
			a.Set(FieldSql_unique_x2, x)
		case IndexSql_unique_x1_x2:
			a.Set(FieldSql_unique_x1_x2, x)
		case IndexSql_primary:
			a.Set(FieldSql_primary, x)
		case IndexSql_jsonb_index:
			a.Set(FieldSql_jsonb_index, x)
		case IndexTime_duration:
			a.Set(FieldTime_duration, x)
		case IndexGo_type_int_to_strings:
			a.Set(FieldGo_type_int_to_strings, x)
		case IndexPublic_field1:
			a.Set(FieldPublic_field1, x)
		case IndexPublic_field2:
			a.Set(FieldPublic_field2, x)
		case IndexPublic_field3:
			a.Set(FieldPublic_field3, x)
		case IndexPublic_field_me1:
			a.Set(FieldPublic_field_me1, x)
		case IndexPublic_field_me2:
			a.Set(FieldPublic_field_me2, x)
		case IndexPublic_field_me3:
			a.Set(FieldPublic_field_me3, x)
		}
	}
	return
}

//NewNewsJson create struct
func NewNewsJson() NewsJson {
	return []byte("{}")
}

//NewsJson is a struct
type NewsJson []byte

//Set value
func (a *NewsJson) Set(k string, v any) *NewsJson {
	(*a) = jsons.Set((*a), k, v)
	return a
}

//Get value
func (a *NewsJson) Get(k string) jsons.Result {
	return jsons.Get((*a), k)
}

//Get value
func (a *NewsJson) DeleteFields(fields ...string) {
	(*a) = jsons.Delete((*a), fields...)
}

//Inc set or get value
func (a *NewsJson) Inc(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldInc)
	}
	a.Set(FieldInc, v[0])
	return
}

//Ints set or get value
func (a *NewsJson) Ints(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldInts)
	}
	a.Set(FieldInts, v[0])
	return
}

//Ints8 set or get value
func (a *NewsJson) Ints8(v ...int8) (res int8) {
	if v == nil {
		return int8(jsons.Int((*a), FieldInts8))
	}
	a.Set(FieldInts8, v[0])
	return
}

//Ints16 set or get value
func (a *NewsJson) Ints16(v ...int16) (res int16) {
	if v == nil {
		return int16(jsons.Int((*a), FieldInts16))
	}
	a.Set(FieldInts16, v[0])
	return
}

//Ints32 set or get value
func (a *NewsJson) Ints32(v ...int32) (res int32) {
	if v == nil {
		return int32(jsons.Int((*a), FieldInts32))
	}
	a.Set(FieldInts32, v[0])
	return
}

//Ints64 set or get value
func (a *NewsJson) Ints64(v ...int64) (res int64) {
	if v == nil {
		return jsons.Int64((*a), FieldInts64)
	}
	a.Set(FieldInts64, v[0])
	return
}

//Uints set or get value
func (a *NewsJson) Uints(v ...uint) (res uint) {
	if v == nil {
		return jsons.Uint((*a), FieldUints)
	}
	a.Set(FieldUints, v[0])
	return
}

//Uints8 set or get value
func (a *NewsJson) Uints8(v ...uint8) (res uint8) {
	if v == nil {
		return jsons.Uint8((*a), FieldUints8)
	}
	a.Set(FieldUints8, v[0])
	return
}

//Uints16 set or get value
func (a *NewsJson) Uints16(v ...uint16) (res uint16) {
	if v == nil {
		return uint16(jsons.Uint((*a), FieldUints16))
	}
	a.Set(FieldUints16, v[0])
	return
}

//Uints32 set or get value
func (a *NewsJson) Uints32(v ...uint32) (res uint32) {
	if v == nil {
		return jsons.Uint32((*a), FieldUints32)
	}
	a.Set(FieldUints32, v[0])
	return
}

//Uints64 set or get value
func (a *NewsJson) Uints64(v ...uint64) (res uint64) {
	if v == nil {
		return jsons.Uint64((*a), FieldUints64)
	}
	a.Set(FieldUints64, v[0])
	return
}

//Floats32 set or get value
func (a *NewsJson) Floats32(v ...float32) (res float32) {
	if v == nil {
		return jsons.Float32((*a), FieldFloats32)
	}
	a.Set(FieldFloats32, v[0])
	return
}

//Floats64 set or get value
func (a *NewsJson) Floats64(v ...float64) (res float64) {
	if v == nil {
		return jsons.Float64((*a), FieldFloats64)
	}
	a.Set(FieldFloats64, v[0])
	return
}

//Bools set or get value
func (a *NewsJson) Bools(v ...bool) (res bool) {
	if v == nil {
		return jsons.Bool((*a), FieldBools)
	}
	a.Set(FieldBools, v[0])
	return
}

//Byte1 set or get value
func (a *NewsJson) Byte1(v ...byte) (res byte) {
	if v == nil {
		return jsons.Byte((*a), FieldByte1)
	}
	a.Set(FieldByte1, v[0])
	return
}

//Bytes set or get value
func (a *NewsJson) Bytes(v ...[]byte) (res []byte) {
	if v == nil {
		return jsons.Bytes((*a), FieldBytes)
	}
	a.Set(FieldBytes, v[0])
	return
}

//List_ints set or get value
func (a *NewsJson) List_ints(v ...int) (res []int) {
	if v == nil {
		return jsons.ArrayInt((*a), FieldList_ints)
	}
	a.Set(FieldList_ints, v)
	return
}

//List_intsAdd add values
func (a *NewsJson) List_intsAdd(v ...int) {
	a.Set(FieldList_ints, append(a.List_ints(), v...))
}

//List_intsPrepend add values
func (a *NewsJson) List_intsPrepend(v ...int) {
	a.Set(FieldList_ints, append(v, a.List_ints()...))
}

//List_intsAddUnique add unique values only
func (a *NewsJson) List_intsAddUnique(v ...int) {
	var list []int
	un := map[int]bool{}
	for _, x := range a.List_ints() {
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

	a.List_ints(list...)
}

//List_intsDelete add unique values only
func (a *NewsJson) List_intsDelete(v int) {
	var list []int
	for _, x := range a.List_ints() {
		if x == v {
			continue
		}
		list = append(list, x)
	}
	a.List_ints(list...)
}

//List_string set or get value
func (a *NewsJson) List_string(v ...string) (res []string) {
	if v == nil {
		return jsons.ArrayString((*a), FieldList_string)
	}
	a.Set(FieldList_string, v)
	return
}

//List_stringAdd add values
func (a *NewsJson) List_stringAdd(v ...string) {
	a.Set(FieldList_string, append(a.List_string(), v...))
}

//List_stringPrepend add values
func (a *NewsJson) List_stringPrepend(v ...string) {
	a.Set(FieldList_string, append(v, a.List_string()...))
}

//List_stringAddUnique add unique values only
func (a *NewsJson) List_stringAddUnique(v ...string) {
	var list []string
	un := map[string]bool{}
	for _, x := range a.List_string() {
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

	a.List_string(list...)
}

//List_stringDelete add unique values only
func (a *NewsJson) List_stringDelete(v string) {
	var list []string
	for _, x := range a.List_string() {
		if x == v {
			continue
		}
		list = append(list, x)
	}
	a.List_string(list...)
}

//List_float set or get value
func (a *NewsJson) List_float(v ...float64) (res []float64) {
	if v == nil {
		_ = jsoniter.Unmarshal([]byte(jsons.Get((*a), FieldList_float).Raw), &res)
		return
	}
	a.Set(FieldList_float, v)
	return
}

//List_floatAdd add values
func (a *NewsJson) List_floatAdd(v ...float64) {
	a.Set(FieldList_float, append(a.List_float(), v...))
}

//List_floatPrepend add values
func (a *NewsJson) List_floatPrepend(v ...float64) {
	a.Set(FieldList_float, append(v, a.List_float()...))
}

//List_floatAddUnique add unique values only
func (a *NewsJson) List_floatAddUnique(v ...float64) {
	var list []float64
	un := map[float64]bool{}
	for _, x := range a.List_float() {
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

	a.List_float(list...)
}

//List_floatDelete add unique values only
func (a *NewsJson) List_floatDelete(v float64) {
	var list []float64
	for _, x := range a.List_float() {
		if x == v {
			continue
		}
		list = append(list, x)
	}
	a.List_float(list...)
}

// Map_string_string set or get value
func (a *NewsJson) Map_string_string(v ...map[string]string) (res map[string]string) {
	if v == nil {
		return jsons.MapString((*a), FieldMap_string_string)
	}
	a.Set(FieldMap_string_string, v[0])
	return
}

// Map_string_stringAdd add values
func (a *NewsJson) Map_string_stringAdd(k string, v string) {
	maps := a.Map_string_string()
	maps[k] = v
	a.Map_string_string(maps)
}

// Map_string_stringDelete add unique values only
func (a *NewsJson) Map_string_stringDelete(k string) {
	maps := a.Map_string_string()
	delete(maps, k)
	a.Map_string_string(maps)
}

// Map_string_stringHas check value
func (a *NewsJson) Map_string_stringHas(k string) bool {
	maps := a.Map_string_string()
	return !reflect.ValueOf(maps[k]).IsZero()
}

//Map_string_bytes set or get value
func (a *NewsJson) Map_string_bytes(v ...map[string][]byte) (res map[string][]byte) {
	if v == nil {
		_ = jsoniter.Unmarshal([]byte(jsons.Get((*a), FieldMap_string_bytes).Raw), &res)
		return
	}
	a.Set(FieldMap_string_bytes, v[0])
	return
}

// Map_string_bool set or get value
func (a *NewsJson) Map_string_bool(v ...map[string]bool) (res map[string]bool) {
	if v == nil {
		return jsons.MapBool((*a), FieldMap_string_bool)
	}
	a.Set(FieldMap_string_bool, v[0])
	return
}

// Map_string_boolAdd add values
func (a *NewsJson) Map_string_boolAdd(k string, v bool) {
	maps := a.Map_string_bool()
	maps[k] = v
	a.Map_string_bool(maps)
}

// Map_string_boolDelete add unique values only
func (a *NewsJson) Map_string_boolDelete(k string) {
	maps := a.Map_string_bool()
	delete(maps, k)
	a.Map_string_bool(maps)
}

// Map_string_boolHas check value
func (a *NewsJson) Map_string_boolHas(k string) bool {
	maps := a.Map_string_bool()
	return !reflect.ValueOf(maps[k]).IsZero()
}

// Map_string_int set or get value
func (a *NewsJson) Map_string_int(v ...map[string]int) (res map[string]int) {
	if v == nil {
		return jsons.MapInt((*a), FieldMap_string_int)
	}
	a.Set(FieldMap_string_int, v[0])
	return
}

// Map_string_intAdd add values
func (a *NewsJson) Map_string_intAdd(k string, v int) {
	maps := a.Map_string_int()
	maps[k] = v
	a.Map_string_int(maps)
}

// Map_string_intDelete add unique values only
func (a *NewsJson) Map_string_intDelete(k string) {
	maps := a.Map_string_int()
	delete(maps, k)
	a.Map_string_int(maps)
}

// Map_string_intHas check value
func (a *NewsJson) Map_string_intHas(k string) bool {
	maps := a.Map_string_int()
	return !reflect.ValueOf(maps[k]).IsZero()
}

// Map_string_float64 set or get value
func (a *NewsJson) Map_string_float64(v ...map[string]float64) (res map[string]float64) {
	if v == nil {
		return jsons.MapFloats((*a), FieldMap_string_float64)
	}
	a.Set(FieldMap_string_float64, v[0])
	return
}

// Map_string_float64Add add values
func (a *NewsJson) Map_string_float64Add(k string, v float64) {
	maps := a.Map_string_float64()
	maps[k] = v
	a.Map_string_float64(maps)
}

// Map_string_float64Delete add unique values only
func (a *NewsJson) Map_string_float64Delete(k string) {
	maps := a.Map_string_float64()
	delete(maps, k)
	a.Map_string_float64(maps)
}

// Map_string_float64Has check value
func (a *NewsJson) Map_string_float64Has(k string) bool {
	maps := a.Map_string_float64()
	return !reflect.ValueOf(maps[k]).IsZero()
}

// Map_string_any set or get value
func (a *NewsJson) Map_string_any(v ...map[string]any) (res map[string]any) {
	if v == nil {
		return jsons.MapAny((*a), FieldMap_string_any)
	}
	a.Set(FieldMap_string_any, v[0])
	return
}

// Map_string_anyAdd add values
func (a *NewsJson) Map_string_anyAdd(k string, v any) {
	maps := a.Map_string_any()
	maps[k] = v
	a.Map_string_any(maps)
}

// Map_string_anyDelete add unique values only
func (a *NewsJson) Map_string_anyDelete(k string) {
	maps := a.Map_string_any()
	delete(maps, k)
	a.Map_string_any(maps)
}

// Map_string_anyHas check value
func (a *NewsJson) Map_string_anyHas(k string) bool {
	maps := a.Map_string_any()
	return !reflect.ValueOf(maps[k]).IsZero()
}

// Map_int_string set or get value
func (a *NewsJson) Map_int_string(v ...map[int]string) (res map[int]string) {
	if v == nil {
		return jsons.MapIntString((*a), FieldMap_int_string)
	}
	a.Set(FieldMap_int_string, v[0])
	return
}

// Map_int_stringAdd add values
func (a *NewsJson) Map_int_stringAdd(k int, v string) {
	maps := a.Map_int_string()
	maps[k] = v
	a.Map_int_string(maps)
}

// Map_int_stringDelete add unique values only
func (a *NewsJson) Map_int_stringDelete(k int) {
	maps := a.Map_int_string()
	delete(maps, k)
	a.Map_int_string(maps)
}

// Map_int_stringHas check value
func (a *NewsJson) Map_int_stringHas(k int) bool {
	maps := a.Map_int_string()
	return !reflect.ValueOf(maps[k]).IsZero()
}

// Map_int_int set or get value
func (a *NewsJson) Map_int_int(v ...map[int]int) (res map[int]int) {
	if v == nil {
		return jsons.MapIntInt((*a), FieldMap_int_int)
	}
	a.Set(FieldMap_int_int, v[0])
	return
}

// Map_int_intAdd add values
func (a *NewsJson) Map_int_intAdd(k int, v int) {
	maps := a.Map_int_int()
	maps[k] = v
	a.Map_int_int(maps)
}

// Map_int_intDelete add unique values only
func (a *NewsJson) Map_int_intDelete(k int) {
	maps := a.Map_int_int()
	delete(maps, k)
	a.Map_int_int(maps)
}

// Map_int_intHas check value
func (a *NewsJson) Map_int_intHas(k int) bool {
	maps := a.Map_int_int()
	return !reflect.ValueOf(maps[k]).IsZero()
}

//Map_int_bool set or get value
func (a *NewsJson) Map_int_bool(v ...map[int]bool) (res map[int]bool) {
	if v == nil {
		_ = jsoniter.Unmarshal([]byte(jsons.Get((*a), FieldMap_int_bool).Raw), &res)
		return
	}
	a.Set(FieldMap_int_bool, v[0])
	return
}

//RenameSQL set or get value
func (a *NewsJson) RenameSQL(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldRenameSQL)
	}
	a.Set(FieldRenameSQL, v[0])
	return
}

//RenameGO_OK set or get value
func (a *NewsJson) RenameGO_OK(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldRenameGO_OK)
	}
	a.Set(FieldRenameGO_OK, v[0])
	return
}

//RenameJS set or get value
func (a *NewsJson) RenameJS(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldRenameJS)
	}
	a.Set(FieldRenameJS, v[0])
	return
}

//MAST_UPPER_GO set or get value
func (a *NewsJson) MAST_UPPER_GO(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldMAST_UPPER_GO)
	}
	a.Set(FieldMAST_UPPER_GO, v[0])
	return
}

//IntToSmallInt set or get value
func (a *NewsJson) IntToSmallInt(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldIntToSmallInt)
	}
	a.Set(FieldIntToSmallInt, v[0])
	return
}

//Skip set or get value
func (a *NewsJson) Skip(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldSkip)
	}
	a.Set(FieldSkip, v[0])
	return
}

//Sql_unique_u1_1 set or get value
func (a *NewsJson) Sql_unique_u1_1(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldSql_unique_u1_1)
	}
	a.Set(FieldSql_unique_u1_1, v[0])
	return
}

//Sql_unique_u1_2 set or get value
func (a *NewsJson) Sql_unique_u1_2(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldSql_unique_u1_2)
	}
	a.Set(FieldSql_unique_u1_2, v[0])
	return
}

//Sql_index1_1 set or get value
func (a *NewsJson) Sql_index1_1(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldSql_index1_1)
	}
	a.Set(FieldSql_index1_1, v[0])
	return
}

//Sql_index1_2 set or get value
func (a *NewsJson) Sql_index1_2(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldSql_index1_2)
	}
	a.Set(FieldSql_index1_2, v[0])
	return
}

//Sql_index1_3 set or get value
func (a *NewsJson) Sql_index1_3(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldSql_index1_3)
	}
	a.Set(FieldSql_index1_3, v[0])
	return
}

//Sql_keys_1 set or get value
func (a *NewsJson) Sql_keys_1(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldSql_keys_1)
	}
	a.Set(FieldSql_keys_1, v[0])
	return
}

//Sql_keys_2 set or get value
func (a *NewsJson) Sql_keys_2(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldSql_keys_2)
	}
	a.Set(FieldSql_keys_2, v[0])
	return
}

//Sql_keys_3 set or get value
func (a *NewsJson) Sql_keys_3(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldSql_keys_3)
	}
	a.Set(FieldSql_keys_3, v[0])
	return
}

//Sql_search set or get value
func (a *NewsJson) Sql_search(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldSql_search)
	}
	a.Set(FieldSql_search, v[0])
	return
}

//Sql_get set or get value
func (a *NewsJson) Sql_get(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldSql_get)
	}
	a.Set(FieldSql_get, v[0])
	return
}

//Sql_unique_x1 set or get value
func (a *NewsJson) Sql_unique_x1(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldSql_unique_x1)
	}
	a.Set(FieldSql_unique_x1, v[0])
	return
}

//Sql_unique_x2 set or get value
func (a *NewsJson) Sql_unique_x2(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldSql_unique_x2)
	}
	a.Set(FieldSql_unique_x2, v[0])
	return
}

//Sql_unique_x1_x2 set or get value
func (a *NewsJson) Sql_unique_x1_x2(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldSql_unique_x1_x2)
	}
	a.Set(FieldSql_unique_x1_x2, v[0])
	return
}

//Sql_primary set or get value
func (a *NewsJson) Sql_primary(v ...float64) (res float64) {
	if v == nil {
		return jsons.Float64((*a), FieldSql_primary)
	}
	a.Set(FieldSql_primary, v[0])
	return
}

// Sql_jsonb_index set or get value
func (a *NewsJson) Sql_jsonb_index(v ...map[string]any) (res map[string]any) {
	if v == nil {
		return jsons.MapAny((*a), FieldSql_jsonb_index)
	}
	a.Set(FieldSql_jsonb_index, v[0])
	return
}

// Sql_jsonb_indexAdd add values
func (a *NewsJson) Sql_jsonb_indexAdd(k string, v any) {
	maps := a.Sql_jsonb_index()
	maps[k] = v
	a.Sql_jsonb_index(maps)
}

// Sql_jsonb_indexDelete add unique values only
func (a *NewsJson) Sql_jsonb_indexDelete(k string) {
	maps := a.Sql_jsonb_index()
	delete(maps, k)
	a.Sql_jsonb_index(maps)
}

// Sql_jsonb_indexHas check value
func (a *NewsJson) Sql_jsonb_indexHas(k string) bool {
	maps := a.Sql_jsonb_index()
	return !reflect.ValueOf(maps[k]).IsZero()
}

//Time_duration set or get value
func (a *NewsJson) Time_duration(v ...time.Duration) (res time.Duration) {
	if v == nil {
		return jsons.TimeDuration((*a), FieldTime_duration)
	}
	a.Set(FieldTime_duration, v[0])
	return
}

//Time_durationString get value as time
func (a *NewsJson) Time_durationString() (res string) {
	return a.Time_duration().String()
}

//Go_type_int_to_strings set or get value
func (a *NewsJson) Go_type_int_to_strings(v ...string) (res []string) {
	if v == nil {
		return jsons.ArrayString((*a), FieldGo_type_int_to_strings)
	}
	a.Set(FieldGo_type_int_to_strings, v)
	return
}

//Go_type_int_to_stringsAdd add values
func (a *NewsJson) Go_type_int_to_stringsAdd(v ...string) {
	a.Set(FieldGo_type_int_to_strings, append(a.Go_type_int_to_strings(), v...))
}

//Go_type_int_to_stringsPrepend add values
func (a *NewsJson) Go_type_int_to_stringsPrepend(v ...string) {
	a.Set(FieldGo_type_int_to_strings, append(v, a.Go_type_int_to_strings()...))
}

//Go_type_int_to_stringsAddUnique add unique values only
func (a *NewsJson) Go_type_int_to_stringsAddUnique(v ...string) {
	var list []string
	un := map[string]bool{}
	for _, x := range a.Go_type_int_to_strings() {
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

	a.Go_type_int_to_strings(list...)
}

//Go_type_int_to_stringsDelete add unique values only
func (a *NewsJson) Go_type_int_to_stringsDelete(v string) {
	var list []string
	for _, x := range a.Go_type_int_to_strings() {
		if x == v {
			continue
		}
		list = append(list, x)
	}
	a.Go_type_int_to_strings(list...)
}

//Public_field1 set or get value
func (a *NewsJson) Public_field1(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldPublic_field1)
	}
	a.Set(FieldPublic_field1, v[0])
	return
}

//Public_field2 set or get value
func (a *NewsJson) Public_field2(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldPublic_field2)
	}
	a.Set(FieldPublic_field2, v[0])
	return
}

//Public_field3 set or get value
func (a *NewsJson) Public_field3(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldPublic_field3)
	}
	a.Set(FieldPublic_field3, v[0])
	return
}

//Public_field_me1 set or get value
func (a *NewsJson) Public_field_me1(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldPublic_field_me1)
	}
	a.Set(FieldPublic_field_me1, v[0])
	return
}

//Public_field_me2 set or get value
func (a *NewsJson) Public_field_me2(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldPublic_field_me2)
	}
	a.Set(FieldPublic_field_me2, v[0])
	return
}

//Public_field_me3 set or get value
func (a *NewsJson) Public_field_me3(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldPublic_field_me3)
	}
	a.Set(FieldPublic_field_me3, v[0])
	return
}

//sql NewsSQL class
type NewsSQL struct{ pool *pgxpool.Pool }

func NewNewsSQL(pool *pgxpool.Pool) (a *NewsSQL) {
	a = new(NewsSQL)
	a.pool = pool
	a.CreateTable()
	return
}

//delete item

func (a *NewsSQL) Conn(f func(conn *pgxpool.Conn) (err error)) (err error) {

	conn, err := a.pool.Acquire(context.Background())
	if err != nil {
		return
	}
	defer conn.Release()
	return f(conn)
}

// parse sql query
func (a *NewsSQL) TableName() (res string) {
	return "news"
}

func (a *NewsSQL) Count(where ...string) (count int, err error) {

	var q string

	switch len(where) {
	case 0:
		q = "select count(*) from news"
	default:
		q = "select count(*) from news where " + strings.Join(where, " ")
	}

	err = a.Conn(func(conn *pgxpool.Conn) (err error) {
		return conn.QueryRow(context.Background(), q).Scan(&count)
	})
	return
}

// Insert struct and return int id
func (a *NewsSQL) Insert(v *News) (id int, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "insert into news (ints, ints8, ints16, ints32, ints64, uints, uints8, uints16, uints32, uints64, floats32, floats64, bools, byte1, bytes, list_ints, list_string, list_float, map_string_string, map_string_bytes, map_string_bool, map_string_int, map_string_float64, map_string_any, map_int_string, map_int_int, map_int_bool, renameSQL_OK, renameGO, renameJS, mast_upper_go, intToSmallInt, sql_unique_u1_1, sql_unique_u1_2, sql_index1_1, sql_index1_2, sql_index1_3, sql_keys_1, sql_keys_2, sql_keys_3, sql_search, sql_get, sql_unique_x1, sql_unique_x2, sql_unique_x1_x2, sql_primary, sql_jsonb_index, time_duration, go_type_int_to_strings, public_field1, public_field2, public_field3, public_field_me1, public_field_me2, public_field_me3) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55) returning inc"
	err = conn.QueryRow(c, q, v.sqlTuple()...).Scan(&id)
	return
}

// Insert struct and return int id
func (a *NewsSQL) InsertNoConflict(v *News) (id int, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "insert into news (ints, ints8, ints16, ints32, ints64, uints, uints8, uints16, uints32, uints64, floats32, floats64, bools, byte1, bytes, list_ints, list_string, list_float, map_string_string, map_string_bytes, map_string_bool, map_string_int, map_string_float64, map_string_any, map_int_string, map_int_int, map_int_bool, renameSQL_OK, renameGO, renameJS, mast_upper_go, intToSmallInt, sql_unique_u1_1, sql_unique_u1_2, sql_index1_1, sql_index1_2, sql_index1_3, sql_keys_1, sql_keys_2, sql_keys_3, sql_search, sql_get, sql_unique_x1, sql_unique_x2, sql_unique_x1_x2, sql_primary, sql_jsonb_index, time_duration, go_type_int_to_strings, public_field1, public_field2, public_field3, public_field_me1, public_field_me2, public_field_me3) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55) on conflict do nothing returning inc"
	err = conn.QueryRow(c, q, v.sqlTuple()...).Scan(&id)
	return
}

// Insert full struct ID must be (ignore all skips)
func (a *NewsSQL) InsertFull(v *News) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "insert into news (inc, ints, ints8, ints16, ints32, ints64, uints, uints8, uints16, uints32, uints64, floats32, floats64, bools, byte1, bytes, list_ints, list_string, list_float, map_string_string, map_string_bytes, map_string_bool, map_string_int, map_string_float64, map_string_any, map_int_string, map_int_int, map_int_bool, renameSQL_OK, renameGO, renameJS, mast_upper_go, intToSmallInt, sql_unique_u1_1, sql_unique_u1_2, sql_index1_1, sql_index1_2, sql_index1_3, sql_keys_1, sql_keys_2, sql_keys_3, sql_search, sql_get, sql_unique_x1, sql_unique_x2, sql_unique_x1_x2, sql_primary, sql_jsonb_index, time_duration, go_type_int_to_strings, public_field1, public_field2, public_field3, public_field_me1, public_field_me2, public_field_me3) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, $54, $55, $56)"
	_, err = conn.Exec(c, q, v.sqlAllTuples()...)
	return
}

// update serial counter (count all items, add 1 and plus custom int)
// useful if you insert with ID
func (a *NewsSQL) ReindexIncSerialCounter(add ...int) (err error) {
	return a.Conn(func(conn *pgxpool.Conn) (err error) {
		plus := 1
		if len(add) > 0 {
			plus = add[0]
		}
		_, err = conn.Exec(context.Background(), "SELECT setval(pg_get_serial_sequence('news', 'inc'), COALESCE((SELECT MAX(inc) FROM news), 0) + $1, false)", plus)
		return
	})
}

// set serial counter
func (a *NewsSQL) SetIncSerialCounter(value int) (err error) {
	return a.Conn(func(conn *pgxpool.Conn) (err error) {
		_, err = conn.Exec(context.Background(), "SELECT setval(pg_get_serial_sequence('news', 'inc'), $1, false)", value)
		return
	})
}

// parse sql query
func (a *NewsSQL) Get(id any, fields ...NewsIndexType) (res *News, err error) {

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
	q = q + "from news where inc = $1 limit 1"
	res = new(News)
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

// get custom struct sql query
func (a *NewsSQL) GetPublic(id any) (res *NewsPublic, err error) {
	p, err := a.Get(id, NewsPublicList()...)
	if err != nil {
		return
	}
	res = p.NewsPublic()
	return
}

// get custom struct sql query
func (a *NewsSQL) GetMe(id any) (res *NewsMe, err error) {
	p, err := a.Get(id, NewsMeList()...)
	if err != nil {
		return
	}
	res = p.NewsMe()
	return
}

// parse sql query
func (a *NewsSQL) GetWhere(where string, fields ...NewsIndexType) (res *News, err error) {

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
	q = q + " from news"
	q = q + fmt.Sprintf(" where %s ", where)
	q = q + " limit 1"
	res = new(News)
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
func (a *NewsSQL) Row(eq map[string]any, fields ...NewsIndexType) (res *News, err error) {

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
		if !NewsValidKey(k) {
			err = errors.New(k)
			return
		}
		in := NewsKeyIndex(k)
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
	res = new(News)
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
func (a *NewsSQL) All(fields ...NewsIndexType) (res []*News, err error) {

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
	q += " from news"
	rows, err := conn.Query(c, q)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item News
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
func (a *NewsSQL) List(limit, offset int, fields ...NewsIndexType) (res []*News, err error) {

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
	q += " from news limit $1 offset $2"
	rows, err := conn.Query(c, q, limit, offset)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item News
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
func (a *NewsSQL) Update(id any, k string, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if !NewsValidKey(k) {
		return errors.New("invalid key")
	}
	in := NewsKeyIndex(k)
	q := fmt.Sprintf("update news set %s = $1 where inc = $2", in.SQLName())
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateWhere(k string, v any, where string) (err error) {
	if !NewsValidKey(k) {
		return errors.New("invalid key")
	}
	in := NewsKeyIndex(k)
	return a.Conn(func(conn *pgxpool.Conn) (err error) {
		q := fmt.Sprintf("update %s set %s = $1 where %s", a.TableName(), in.SQLName(), where)
		_, err = conn.Exec(context.Background(), q, v)
		return
	})
}

// update sql query
func (a *NewsSQL) Updates(id any, keys map[string]any, where ...string) (err error) {

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
		if !NewsValidKey(k) {
			return errors.New(k)
		}
		in := NewsKeyIndex(k)
		count++
		fields = append(fields, fmt.Sprintf("%s = $%d", in.SQLName(), count))
		values = append(values, v)
	}

	list := strings.Join(fields, ", ")
	count++
	values = append(values, id)

	q := fmt.Sprintf("update news set %s where inc = $%d", list, count)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, values...)
	return
}

// update sql query
func (a *NewsSQL) UpdatesWhere(keys map[string]any, where string, args ...any) (err error) {

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
		if !NewsValidKey(k) {
			return errors.New(k)
		}
		in := NewsKeyIndex(k)
		count++
		fields = append(fields, fmt.Sprintf("%s = $%d", in.SQLName(), count))
		values = append(values, v)
	}

	list := strings.Join(fields, ", ")
	count++

	where = fmt.Sprintf(where, args...)
	q := fmt.Sprintf("update news set %s where %s", list, where)
	_, err = conn.Exec(c, q, values...)
	return
}

// update sql query
func (a *NewsSQL) UpdateUints16(id any, k string, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set uints16 = uints16 || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesUints16(id any, keys map[string]any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := jsoniter.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set uints16 = uints16 || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyUints16(id any, k string, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set uints16 = uints16 - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyUints16(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set uints16 = uints16 - $1 || jsonb_build_object($2, uints16->$1) where inc = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// add array value to jsonb array
func (a *NewsSQL) AddList_ints(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set list_ints = list_ints || '%s'::jsonb where inc = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *NewsSQL) DeleteList_ints(id any, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set list_ints = list_ints - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// add array value to jsonb array
func (a *NewsSQL) AddList_intsWhere(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set list_ints = list_ints || '%s'::jsonb where %s", res, where)
	_, err = conn.Exec(c, q)
	return
}

// delete array value from jsonb array
func (a *NewsSQL) DeleteList_intsWhere(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update news set list_ints = list_ints - $1 where %s", where)
	_, err = conn.Exec(c, q, v)
	return
}

// add array value to jsonb array
func (a *NewsSQL) AddList_string(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set list_string = list_string || '%s'::jsonb where inc = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *NewsSQL) DeleteList_string(id any, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set list_string = list_string - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// add array value to jsonb array
func (a *NewsSQL) AddList_stringWhere(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set list_string = list_string || '%s'::jsonb where %s", res, where)
	_, err = conn.Exec(c, q)
	return
}

// delete array value from jsonb array
func (a *NewsSQL) DeleteList_stringWhere(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update news set list_string = list_string - $1 where %s", where)
	_, err = conn.Exec(c, q, v)
	return
}

// add array value to jsonb array
func (a *NewsSQL) AddList_float(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set list_float = list_float || '%s'::jsonb where inc = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *NewsSQL) DeleteList_float(id any, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set list_float = list_float - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// add array value to jsonb array
func (a *NewsSQL) AddList_floatWhere(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set list_float = list_float || '%s'::jsonb where %s", res, where)
	_, err = conn.Exec(c, q)
	return
}

// delete array value from jsonb array
func (a *NewsSQL) DeleteList_floatWhere(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update news set list_float = list_float - $1 where %s", where)
	_, err = conn.Exec(c, q, v)
	return
}

// update sql query
func (a *NewsSQL) UpdateMap_string_string(id any, k string, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set map_string_string = map_string_string || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesMap_string_string(id any, keys map[string]any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := jsoniter.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set map_string_string = map_string_string || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyMap_string_string(id any, k string, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_string_string = map_string_string - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyMap_string_string(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_string_string = map_string_string - $1 || jsonb_build_object($2, map_string_string->$1) where inc = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateMap_string_bytes(id any, k string, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set map_string_bytes = map_string_bytes || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesMap_string_bytes(id any, keys map[string]any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := jsoniter.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set map_string_bytes = map_string_bytes || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyMap_string_bytes(id any, k string, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_string_bytes = map_string_bytes - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyMap_string_bytes(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_string_bytes = map_string_bytes - $1 || jsonb_build_object($2, map_string_bytes->$1) where inc = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateMap_string_bool(id any, k string, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set map_string_bool = map_string_bool || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesMap_string_bool(id any, keys map[string]any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := jsoniter.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set map_string_bool = map_string_bool || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyMap_string_bool(id any, k string, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_string_bool = map_string_bool - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyMap_string_bool(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_string_bool = map_string_bool - $1 || jsonb_build_object($2, map_string_bool->$1) where inc = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateMap_string_int(id any, k string, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set map_string_int = map_string_int || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesMap_string_int(id any, keys map[string]any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := jsoniter.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set map_string_int = map_string_int || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyMap_string_int(id any, k string, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_string_int = map_string_int - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyMap_string_int(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_string_int = map_string_int - $1 || jsonb_build_object($2, map_string_int->$1) where inc = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateMap_string_float64(id any, k string, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set map_string_float64 = map_string_float64 || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesMap_string_float64(id any, keys map[string]any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := jsoniter.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set map_string_float64 = map_string_float64 || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyMap_string_float64(id any, k string, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_string_float64 = map_string_float64 - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyMap_string_float64(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_string_float64 = map_string_float64 - $1 || jsonb_build_object($2, map_string_float64->$1) where inc = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateMap_string_any(id any, k string, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set map_string_any = map_string_any || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesMap_string_any(id any, keys map[string]any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := jsoniter.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set map_string_any = map_string_any || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyMap_string_any(id any, k string, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_string_any = map_string_any - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyMap_string_any(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_string_any = map_string_any - $1 || jsonb_build_object($2, map_string_any->$1) where inc = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateMap_int_string(id any, k string, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set map_int_string = map_int_string || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesMap_int_string(id any, keys map[string]any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := jsoniter.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set map_int_string = map_int_string || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyMap_int_string(id any, k string, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_int_string = map_int_string - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyMap_int_string(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_int_string = map_int_string - $1 || jsonb_build_object($2, map_int_string->$1) where inc = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateMap_int_int(id any, k string, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set map_int_int = map_int_int || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesMap_int_int(id any, keys map[string]any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := jsoniter.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set map_int_int = map_int_int || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyMap_int_int(id any, k string, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_int_int = map_int_int - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyMap_int_int(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_int_int = map_int_int - $1 || jsonb_build_object($2, map_int_int->$1) where inc = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateMap_int_bool(id any, k string, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set map_int_bool = map_int_bool || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesMap_int_bool(id any, keys map[string]any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := jsoniter.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set map_int_bool = map_int_bool || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyMap_int_bool(id any, k string, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_int_bool = map_int_bool - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyMap_int_bool(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set map_int_bool = map_int_bool - $1 || jsonb_build_object($2, map_int_bool->$1) where inc = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateSql_jsonb_index(id any, k string, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set sql_jsonb_index = sql_jsonb_index || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesSql_jsonb_index(id any, keys map[string]any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := jsoniter.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set sql_jsonb_index = sql_jsonb_index || $$%s$$::jsonb where inc = $1", res)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeySql_jsonb_index(id any, k string, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set sql_jsonb_index = sql_jsonb_index - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeySql_jsonb_index(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set sql_jsonb_index = sql_jsonb_index - $1 || jsonb_build_object($2, sql_jsonb_index->$1) where inc = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// add array value to jsonb array
func (a *NewsSQL) AddGo_type_int_to_strings(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set go_type_int_to_strings = go_type_int_to_strings || '%s'::jsonb where inc = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *NewsSQL) DeleteGo_type_int_to_strings(id any, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set go_type_int_to_strings = go_type_int_to_strings - $1 where inc = $2"
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// add array value to jsonb array
func (a *NewsSQL) AddGo_type_int_to_stringsWhere(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set go_type_int_to_strings = go_type_int_to_strings || '%s'::jsonb where %s", res, where)
	_, err = conn.Exec(c, q)
	return
}

// delete array value from jsonb array
func (a *NewsSQL) DeleteGo_type_int_to_stringsWhere(v any, where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update news set go_type_int_to_strings = go_type_int_to_strings - $1 where %s", where)
	_, err = conn.Exec(c, q, v)
	return
}

type NewsQuery struct {
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

func (a *NewsQuery) Render() (sql string, fields []NewsIndexType, values []any) {

	switch len(a.Fields) == 0 {
	case true:
		fields = []NewsIndexType{IndexInc, IndexInts, IndexInts8, IndexInts16, IndexInts32, IndexInts64, IndexUints, IndexUints8, IndexUints16, IndexUints32, IndexUints64, IndexFloats32, IndexFloats64, IndexBools, IndexByte1, IndexBytes, IndexList_ints, IndexList_string, IndexList_float, IndexMap_string_string, IndexMap_string_bytes, IndexMap_string_bool, IndexMap_string_int, IndexMap_string_float64, IndexMap_string_any, IndexMap_int_string, IndexMap_int_int, IndexMap_int_bool, IndexRenameSQL, IndexRenameGO_OK, IndexRenameJS, IndexMAST_UPPER_GO, IndexIntToSmallInt, IndexSql_unique_u1_1, IndexSql_unique_u1_2, IndexSql_index1_1, IndexSql_index1_2, IndexSql_index1_3, IndexSql_keys_1, IndexSql_keys_2, IndexSql_keys_3, IndexSql_search, IndexSql_get, IndexSql_unique_x1, IndexSql_unique_x2, IndexSql_unique_x1_x2, IndexSql_primary, IndexSql_jsonb_index, IndexTime_duration, IndexGo_type_int_to_strings, IndexPublic_field1, IndexPublic_field2, IndexPublic_field3, IndexPublic_field_me1, IndexPublic_field_me2, IndexPublic_field_me3}
	default:
		for _, x := range a.Fields {
			if !NewsValidKey(x) {
				continue
			}
			fields = append(fields, NewsKeyIndex(x))
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
	list = append(list, "from news")

	var andlist []string

	// EQ where
	for k, v := range a.EQ {
		if !NewsValidKey(k) {
			continue
		}
		p := NewsKeyIndex(k)
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

		if !NewsValidKey(k) {
			continue
		}
		sqlName := NewsKeyIndex(k).SQLName()
		count++
		andlist = append(andlist, fmt.Sprintf("%s > $%d", sqlName, count))
		values = append(values, v)
	}

	// LT where
	for k, v := range a.LT {
		if !NewsValidKey(k) {
			continue
		}
		sqlName := NewsKeyIndex(k).SQLName()
		count++
		andlist = append(andlist, fmt.Sprintf("%s < $%d", sqlName, count))
		values = append(values, v)
	}

	// NOT where
	for k, v := range a.NOT {
		if !NewsValidKey(k) {
			continue
		}
		sqlName := NewsKeyIndex(k).SQLName()

		for _, x := range v {
			count++
			andlist = append(andlist, fmt.Sprintf("%s != $%d", sqlName, count))
			values = append(values, x)
		}

	}

	// LIKE where
	for k, v := range a.Like {
		if !NewsValidKey(k) {
			continue
		}
		sqlName := NewsKeyIndex(k).SQLName()
		count++
		andlist = append(andlist, fmt.Sprintf("%s ilike $%d", sqlName, count))
		values = append(values, "%"+v+"%")
	}

	// IN where
	if a.IN != nil {
		for k, v := range a.IN {
			if !NewsValidKey(k) {
				continue
			}
			sqlName := NewsKeyIndex(k).SQLName()
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
			if !NewsValidKey(k) {
				continue
			}
			sqlName := NewsKeyIndex(k).SQLName()
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
			if !NewsValidKey(k) {
				continue
			}
			sqlName := NewsKeyIndex(k).SQLName()
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
	if a.Sort != "" && NewsValidKey(a.Sort) {
		list = append(list, fmt.Sprintf("order by %s", NewsKeyIndex(a.Sort).SQLName()))
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

func NewNewsQuery() *NewsQuery {
	a := new(NewsQuery)
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

func (a *NewsSQL) Search(q *NewsQuery) (res []*News, err error) {

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
		var item News
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
func (a *NewsSQL) Select(where string, fields ...NewsIndexType) (res []*News, err error) {

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
	q = q + "from news"
	if where != "" {
		q += " where " + where
	}
	rows, err := conn.Query(c, q)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item News
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
func (a *NewsSQL) Delete(id any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	_, err = conn.Exec(c, "delete from news where inc = $1", id)
	return

}

// delete item where: i = 1 and w = 'nice'
func (a *NewsSQL) DeleteWhere(where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	_, err = conn.Exec(c, fmt.Sprintf("delete from news where %s", where))
	return

}

// has value in db
func (a *NewsSQL) Has(field NewsIndexType, v any) (has bool, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("select exists (select 1 from news where %s = $1 limit 1)", field.SQLName())
	err = conn.QueryRow(c, q, v).Scan(&has)
	return
}

// Create table
func (a *NewsSQL) CreateTable() (err error) {
	return a.Conn(func(conn *pgxpool.Conn) (err error) {
		q := `create table if not exists news (
	inc                        bigserial,
	ints                       bigint,
	ints8                      smallint,
	ints16                     int,
	ints32                     int,
	ints64                     bigint,
	uints                      bigint,
	uints8                     smallint,
	uints16                    jsonb default '{}'::jsonb,
	uints32                    int,
	uints64                    bigint,
	floats32                   double precision,
	floats64                   double precision,
	bools                      boolean,
	byte1                      smallint,
	bytes                      bytea,
	list_ints                  jsonb default '[]'::jsonb,
	list_string                jsonb default '[]'::jsonb,
	list_float                 jsonb default '[]'::jsonb,
	map_string_string          jsonb default '{}'::jsonb,
	map_string_bytes           jsonb default '{}'::jsonb,
	map_string_bool            jsonb default '{}'::jsonb,
	map_string_int             jsonb default '{}'::jsonb,
	map_string_float64         jsonb default '{}'::jsonb,
	map_string_any             jsonb default '{}'::jsonb,
	map_int_string             jsonb default '{}'::jsonb,
	map_int_int                jsonb default '{}'::jsonb,
	map_int_bool               jsonb default '{}'::jsonb,
	renameSQL_OK               text,
	renameGO                   text,
	renameJS                   text,
	mast_upper_go              text,
	intToSmallInt              smallint,
	sql_unique_u1_1            bigint,
	sql_unique_u1_2            bigint,
	sql_index1_1               bigint,
	sql_index1_2               bigint,
	sql_index1_3               bigint,
	sql_keys_1                 bigint,
	sql_keys_2                 bigint,
	sql_keys_3                 bigint,
	sql_search                 text,
	sql_get                    text,
	sql_unique_x1              bigint,
	sql_unique_x2              bigint,
	sql_unique_x1_x2           bigint,
	sql_primary                double precision,
	sql_jsonb_index            jsonb default '{}'::jsonb,
	time_duration              bigint,
	go_type_int_to_strings     jsonb default '[]'::jsonb,
	public_field1              bigint,
	public_field2              bigint,
	public_field3              bigint,
	public_field_me1           bigint,
	public_field_me2           bigint,
	public_field_me3           bigint,
	unique(sql_unique_u1_1,sql_unique_u1_2),
	unique(sql_unique_x1,sql_unique_x1_x2),
	unique(sql_unique_x2,sql_unique_x1_x2),
	primary key (sql_primary)
)
`
		_, err = conn.Exec(context.Background(), q)
		return
	})
}

// parse sql query
func (a *NewsSQL) Fields() (res []NewsIndexType) {
	return []NewsIndexType{IndexInc, IndexInts, IndexInts8, IndexInts16, IndexInts32, IndexInts64, IndexUints, IndexUints8, IndexUints16, IndexUints32, IndexUints64, IndexFloats32, IndexFloats64, IndexBools, IndexByte1, IndexBytes, IndexList_ints, IndexList_string, IndexList_float, IndexMap_string_string, IndexMap_string_bytes, IndexMap_string_bool, IndexMap_string_int, IndexMap_string_float64, IndexMap_string_any, IndexMap_int_string, IndexMap_int_int, IndexMap_int_bool, IndexRenameSQL, IndexRenameGO_OK, IndexRenameJS, IndexMAST_UPPER_GO, IndexIntToSmallInt, IndexSql_unique_u1_1, IndexSql_unique_u1_2, IndexSql_index1_1, IndexSql_index1_2, IndexSql_index1_3, IndexSql_keys_1, IndexSql_keys_2, IndexSql_keys_3, IndexSql_search, IndexSql_get, IndexSql_unique_x1, IndexSql_unique_x2, IndexSql_unique_x1_x2, IndexSql_primary, IndexSql_jsonb_index, IndexTime_duration, IndexGo_type_int_to_strings, IndexPublic_field1, IndexPublic_field2, IndexPublic_field3, IndexPublic_field_me1, IndexPublic_field_me2, IndexPublic_field_me3}
} //clickhouse NewsCQL class
type NewsCQL struct {
	conn driver.Conn
	list []News
	sync.Mutex
}

func NewNewsCQL(conn driver.Conn) (a *NewsCQL) {
	a = new(NewsCQL)
	a.conn = conn
	a.CreateTable()
	go a.deamon()
	return
}

// parse clickhouse query
func (a *NewsCQL) TableName() (res string) {
	return "news"
}

func (a *NewsCQL) Count(where ...string) (count int, err error) {

	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	var q string
	switch len(where) {
	case 0:
		q = "select count() from news"
	default:
		q = "select count() from news where " + strings.Join(where, " ")
	}
	var n uint64
	err = a.conn.QueryRow(context.Background(), q).Scan(&n)
	count = int(n)
	return
}

// Add items to clickhouse queue
func (a *NewsCQL) Add(items ...News) {
	a.Lock()
	defer a.Unlock()
	a.list = append(a.list, items...)
}

// Push queued items to clickhouse
func (a *NewsCQL) Push() (err error) {
	a.Lock()
	defer a.Unlock()
	if len(a.list) == 0 {
		return
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	batch, err := a.conn.PrepareBatch(context.Background(), "insert into news (inc, ints, ints8, ints16, ints32, ints64, uints, uints8, uints16, uints32, uints64, floats32, floats64, bools, byte1, bytes, list_ints, list_string, list_float, map_string_string, map_string_bytes, map_string_bool, map_string_int, map_string_float64, map_string_any, map_int_string, map_int_int, map_int_bool, renameSQL, renameGO, renameJS, mast_upper_go, intToSmallInt, skip, sql_unique_u1_1, sql_unique_u1_2, sql_index1_1, sql_index1_2, sql_index1_3, sql_keys_1, sql_keys_2, sql_keys_3, sql_search, sql_get, sql_unique_x1, sql_unique_x2, sql_unique_x1_x2, sql_primary, sql_jsonb_index, time_duration, go_type_int_to_strings, public_field1, public_field2, public_field3, public_field_me1, public_field_me2, public_field_me3)")
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
func (a *NewsCQL) deamon() {
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
func (a *NewsCQL) Get(id any, fields ...NewsIndexType) (res *News, err error) {

	if fields == nil {
		fields = a.Fields()
	}
	var list []string
	for _, x := range fields {
		list = append(list, x.ClickhouseName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s from news where inc = ? limit 1", fieldlist)
	res = new(News)

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
func (a *NewsCQL) GetWhere(where string, fields ...NewsIndexType) (res *News, err error) {

	if fields == nil {
		fields = a.Fields()
	}
	var list []string
	for _, x := range fields {
		list = append(list, x.ClickhouseName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s from news where %s limit 1", fieldlist, where)
	res = new(News)

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
func (a *NewsCQL) Update(id any, k string, v any, where ...string) (err error) {
	if !NewsValidKey(k) {
		return errors.New("invalid key")
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	in := NewsKeyIndex(k)
	q := fmt.Sprintf("alter table news update %s = ? where inc = ?", in.ClickhouseName())
	if len(where) > 0 {
		q += " and " + where[0]
	}
	return a.conn.Exec(context.Background(), q, v, id)
}

// update clickhouse query
func (a *NewsCQL) UpdateWhere(k string, v any, where string) (err error) {
	if !NewsValidKey(k) {
		return errors.New("invalid key")
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	in := NewsKeyIndex(k)
	q := fmt.Sprintf("alter table %s update %s = ? where %s", a.TableName(), in.ClickhouseName(), where)
	return a.conn.Exec(context.Background(), q, v)
}

// update clickhouse query
func (a *NewsCQL) Updates(id any, keys map[string]any, where ...string) (err error) {

	if keys == nil {
		return errors.New("emptykeys")
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	var fields []string
	var values []any
	for k, v := range keys {
		if !NewsValidKey(k) {
			return errors.New(k)
		}
		in := NewsKeyIndex(k)
		fields = append(fields, fmt.Sprintf("%s = ?", in.ClickhouseName()))
		values = append(values, v)
	}
	list := strings.Join(fields, ", ")
	values = append(values, id)

	q := fmt.Sprintf("alter table news update %s where inc = ?", list)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	return a.conn.Exec(context.Background(), q, values...)
}

// update clickhouse query
func (a *NewsCQL) UpdatesWhere(keys map[string]any, where string, args ...any) (err error) {

	if keys == nil {
		return errors.New("emptykeys")
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	var fields []string
	var values []any
	for k, v := range keys {
		if !NewsValidKey(k) {
			return errors.New(k)
		}
		in := NewsKeyIndex(k)
		fields = append(fields, fmt.Sprintf("%s = ?", in.ClickhouseName()))
		values = append(values, v)
	}
	list := strings.Join(fields, ", ")

	where = fmt.Sprintf(where, args...)
	q := fmt.Sprintf("alter table news update %s where %s", list, where)
	return a.conn.Exec(context.Background(), q, values...)
}

// select custom clickhouse query
func (a *NewsCQL) Select(where string, fields ...NewsIndexType) (res []*News, err error) {

	if fields == nil {
		fields = a.Fields()
	}
	var list []string
	for _, x := range fields {
		list = append(list, x.ClickhouseName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s from news", fieldlist)
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
		var item News
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
func (a *NewsCQL) Delete(id any) (err error) {
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	q := "alter table news delete where inc = ?"
	return a.conn.Exec(context.Background(), q, id)
}

// delete clickhouse items
func (a *NewsCQL) DeleteWhere(where string) (err error) {
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	q := "alter table news delete where " + where
	return a.conn.Exec(context.Background(), q)
}

// has value in clickhouse
func (a *NewsCQL) Has(field NewsIndexType, v any) (has bool, err error) {
	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	q := fmt.Sprintf("select count() from news where %s = ? limit 1", field.ClickhouseName())
	var n uint64
	err = a.conn.QueryRow(context.Background(), q, v).Scan(&n)
	has = n > 0
	return
}

// stat by clickhouse field
func (a *NewsCQL) Stat(field NewsIndexType) (res map[string]int, err error) {
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
func (a *NewsCQL) CreateTable() (err error) {
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	q := `create table if not exists news (
	inc                        Int64,
	ints                       Int64,
	ints8                      Int8,
	ints16                     Int16,
	ints32                     Int32,
	ints64                     Int64,
	uints                      UInt64,
	uints8                     UInt8,
	uints16                    UInt16,
	uints32                    UInt32,
	uints64                    UInt64,
	floats32                   Float32,
	floats64                   Float64,
	bools                      Bool,
	byte1                      String,
	bytes                      String,
	list_ints                  Array(Int64),
	list_string                Array(String),
	list_float                 Array(Float64),
	map_string_string          Map(String, String),
	map_string_bytes           JSON,
	map_string_bool            Map(String, Bool),
	map_string_int             Map(String, Int64),
	map_string_float64         Map(String, Float64),
	map_string_any             JSON,
	map_int_string             Map(Int64, String),
	map_int_int                Map(Int64, Int64),
	map_int_bool               Map(Int64, Bool),
	renameSQL                  String,
	renameGO                   String,
	renameJS                   String,
	mast_upper_go              String,
	intToSmallInt              Int64,
	skip                       String,
	sql_unique_u1_1            Int64,
	sql_unique_u1_2            Int64,
	sql_index1_1               Int64,
	sql_index1_2               Int64,
	sql_index1_3               Int64,
	sql_keys_1                 Int64,
	sql_keys_2                 Int64,
	sql_keys_3                 Int64,
	sql_search                 String,
	sql_get                    String,
	sql_unique_x1              Int64,
	sql_unique_x2              Int64,
	sql_unique_x1_x2           Int64,
	sql_primary                Float64,
	sql_jsonb_index            JSON,
	time_duration              Int64,
	go_type_int_to_strings     Array(String),
	public_field1              Int64,
	public_field2              Int64,
	public_field3              Int64,
	public_field_me1           Int64,
	public_field_me2           Int64,
	public_field_me3           Int64
)
engine = MergeTree()
order by tuple()`
	return a.conn.Exec(context.Background(), q)
}

// parse clickhouse query
func (a *NewsCQL) Fields() (res []NewsIndexType) {
	return []NewsIndexType{IndexInc, IndexInts, IndexInts8, IndexInts16, IndexInts32, IndexInts64, IndexUints, IndexUints8, IndexUints16, IndexUints32, IndexUints64, IndexFloats32, IndexFloats64, IndexBools, IndexByte1, IndexBytes, IndexList_ints, IndexList_string, IndexList_float, IndexMap_string_string, IndexMap_string_bytes, IndexMap_string_bool, IndexMap_string_int, IndexMap_string_float64, IndexMap_string_any, IndexMap_int_string, IndexMap_int_int, IndexMap_int_bool, IndexRenameSQL, IndexRenameGO_OK, IndexRenameJS, IndexMAST_UPPER_GO, IndexIntToSmallInt, IndexSkip, IndexSql_unique_u1_1, IndexSql_unique_u1_2, IndexSql_index1_1, IndexSql_index1_2, IndexSql_index1_3, IndexSql_keys_1, IndexSql_keys_2, IndexSql_keys_3, IndexSql_search, IndexSql_get, IndexSql_unique_x1, IndexSql_unique_x2, IndexSql_unique_x1_x2, IndexSql_primary, IndexSql_jsonb_index, IndexTime_duration, IndexGo_type_int_to_strings, IndexPublic_field1, IndexPublic_field2, IndexPublic_field3, IndexPublic_field_me1, IndexPublic_field_me2, IndexPublic_field_me3}
}

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type NewsProto struct {
	state              protoimpl.MessageState     `protogen:"open.v1"`
	Inc                int64                      `protobuf:"varint,1,opt,name=inc,proto3" json:"inc,omitempty"`
	Ints               int64                      `protobuf:"varint,2,opt,name=ints,proto3" json:"ints,omitempty"`
	Ints8              int32                      `protobuf:"varint,3,opt,name=ints8,proto3" json:"ints8,omitempty"`
	Ints16             int32                      `protobuf:"varint,4,opt,name=ints16,proto3" json:"ints16,omitempty"`
	Ints32             int32                      `protobuf:"varint,5,opt,name=ints32,proto3" json:"ints32,omitempty"`
	Ints64             int64                      `protobuf:"varint,6,opt,name=ints64,proto3" json:"ints64,omitempty"`
	Uints              uint64                     `protobuf:"varint,7,opt,name=uints,proto3" json:"uints,omitempty"`
	Uints8             uint32                     `protobuf:"varint,8,opt,name=uints8,proto3" json:"uints8,omitempty"`
	Uints16            uint32                     `protobuf:"varint,9,opt,name=uints16,proto3" json:"uints16,omitempty"`
	Uints32            uint32                     `protobuf:"varint,10,opt,name=uints32,proto3" json:"uints32,omitempty"`
	Uints64            uint64                     `protobuf:"varint,11,opt,name=uints64,proto3" json:"uints64,omitempty"`
	Floats32           float32                    `protobuf:"fixed32,12,opt,name=floats32,proto3" json:"floats32,omitempty"`
	Floats64           float64                    `protobuf:"fixed64,13,opt,name=floats64,proto3" json:"floats64,omitempty"`
	Bools              bool                       `protobuf:"varint,14,opt,name=bools,proto3" json:"bools,omitempty"`
	Byte1              uint32                     `protobuf:"varint,15,opt,name=byte1,proto3" json:"byte1,omitempty"`
	Bytes              []byte                     `protobuf:"bytes,16,opt,name=bytes,proto3" json:"bytes,omitempty"`
	ListInts           []int64                    `protobuf:"varint,17,rep,packed,name=list_ints,json=listInts,proto3" json:"list_ints,omitempty"`
	ListString         []string                   `protobuf:"bytes,18,rep,name=list_string,json=listString,proto3" json:"list_string,omitempty"`
	ListFloat          []float64                  `protobuf:"fixed64,19,rep,packed,name=list_float,json=listFloat,proto3" json:"list_float,omitempty"`
	MapStringString    map[string]string          `protobuf:"bytes,20,rep,name=map_string_string,json=mapStringString,proto3" json:"map_string_string,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	MapStringBytes     map[string][]byte          `protobuf:"bytes,21,rep,name=map_string_bytes,json=mapStringBytes,proto3" json:"map_string_bytes,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	MapStringBool      map[string]bool            `protobuf:"bytes,22,rep,name=map_string_bool,json=mapStringBool,proto3" json:"map_string_bool,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	MapStringInt       map[string]int64           `protobuf:"bytes,23,rep,name=map_string_int,json=mapStringInt,proto3" json:"map_string_int,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	MapStringFloat64   map[string]float64         `protobuf:"bytes,24,rep,name=map_string_float64,json=mapStringFloat64,proto3" json:"map_string_float64,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"fixed64,2,opt,name=value"`
	MapStringAny       map[string]*structpb.Value `protobuf:"bytes,25,rep,name=map_string_any,json=mapStringAny,proto3" json:"map_string_any,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	MapIntString       map[int64]string           `protobuf:"bytes,26,rep,name=map_int_string,json=mapIntString,proto3" json:"map_int_string,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	MapIntInt          map[int64]int64            `protobuf:"bytes,27,rep,name=map_int_int,json=mapIntInt,proto3" json:"map_int_int,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	MapIntBool         map[int64]bool             `protobuf:"bytes,28,rep,name=map_int_bool,json=mapIntBool,proto3" json:"map_int_bool,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"varint,2,opt,name=value"`
	RenameSql          string                     `protobuf:"bytes,29,opt,name=rename_sql,json=renameSql,proto3" json:"rename_sql,omitempty"`
	RenameGo           string                     `protobuf:"bytes,30,opt,name=rename_go,json=renameGo,proto3" json:"rename_go,omitempty"`
	RenameJs           string                     `protobuf:"bytes,31,opt,name=rename_js,json=renameJs,proto3" json:"rename_js,omitempty"`
	MastUpperGo        string                     `protobuf:"bytes,32,opt,name=mast_upper_go,json=mastUpperGo,proto3" json:"mast_upper_go,omitempty"`
	IntToSmallInt      int64                      `protobuf:"varint,33,opt,name=int_to_small_int,json=intToSmallInt,proto3" json:"int_to_small_int,omitempty"`
	Skip               string                     `protobuf:"bytes,34,opt,name=skip,proto3" json:"skip,omitempty"`
	SqlUniqueU1_1      int64                      `protobuf:"varint,35,opt,name=sql_unique_u1_1,json=sqlUniqueU11,proto3" json:"sql_unique_u1_1,omitempty"`
	SqlUniqueU1_2      int64                      `protobuf:"varint,36,opt,name=sql_unique_u1_2,json=sqlUniqueU12,proto3" json:"sql_unique_u1_2,omitempty"`
	SqlIndex1_1        int64                      `protobuf:"varint,37,opt,name=sql_index1_1,json=sqlIndex11,proto3" json:"sql_index1_1,omitempty"`
	SqlIndex1_2        int64                      `protobuf:"varint,38,opt,name=sql_index1_2,json=sqlIndex12,proto3" json:"sql_index1_2,omitempty"`
	SqlIndex1_3        int64                      `protobuf:"varint,39,opt,name=sql_index1_3,json=sqlIndex13,proto3" json:"sql_index1_3,omitempty"`
	SqlKeys_1          int64                      `protobuf:"varint,40,opt,name=sql_keys_1,json=sqlKeys1,proto3" json:"sql_keys_1,omitempty"`
	SqlKeys_2          int64                      `protobuf:"varint,41,opt,name=sql_keys_2,json=sqlKeys2,proto3" json:"sql_keys_2,omitempty"`
	SqlKeys_3          int64                      `protobuf:"varint,42,opt,name=sql_keys_3,json=sqlKeys3,proto3" json:"sql_keys_3,omitempty"`
	SqlSearch          string                     `protobuf:"bytes,43,opt,name=sql_search,json=sqlSearch,proto3" json:"sql_search,omitempty"`
	SqlGet             string                     `protobuf:"bytes,44,opt,name=sql_get,json=sqlGet,proto3" json:"sql_get,omitempty"`
	SqlUniqueX1        int64                      `protobuf:"varint,45,opt,name=sql_unique_x1,json=sqlUniqueX1,proto3" json:"sql_unique_x1,omitempty"`
	SqlUniqueX2        int64                      `protobuf:"varint,46,opt,name=sql_unique_x2,json=sqlUniqueX2,proto3" json:"sql_unique_x2,omitempty"`
	SqlUniqueX1X2      int64                      `protobuf:"varint,47,opt,name=sql_unique_x1_x2,json=sqlUniqueX1X2,proto3" json:"sql_unique_x1_x2,omitempty"`
	SqlPrimary         float64                    `protobuf:"fixed64,48,opt,name=sql_primary,json=sqlPrimary,proto3" json:"sql_primary,omitempty"`
	SqlJsonbIndex      map[string]*structpb.Value `protobuf:"bytes,49,rep,name=sql_jsonb_index,json=sqlJsonbIndex,proto3" json:"sql_jsonb_index,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	TimeDuration       int64                      `protobuf:"varint,50,opt,name=time_duration,json=timeDuration,proto3" json:"time_duration,omitempty"`
	GoTypeIntToStrings []string                   `protobuf:"bytes,51,rep,name=go_type_int_to_strings,json=goTypeIntToStrings,proto3" json:"go_type_int_to_strings,omitempty"`
	PublicField1       int64                      `protobuf:"varint,52,opt,name=public_field1,json=publicField1,proto3" json:"public_field1,omitempty"`
	PublicField2       int64                      `protobuf:"varint,53,opt,name=public_field2,json=publicField2,proto3" json:"public_field2,omitempty"`
	PublicField3       int64                      `protobuf:"varint,54,opt,name=public_field3,json=publicField3,proto3" json:"public_field3,omitempty"`
	PublicFieldMe1     int64                      `protobuf:"varint,55,opt,name=public_field_me1,json=publicFieldMe1,proto3" json:"public_field_me1,omitempty"`
	PublicFieldMe2     int64                      `protobuf:"varint,56,opt,name=public_field_me2,json=publicFieldMe2,proto3" json:"public_field_me2,omitempty"`
	PublicFieldMe3     int64                      `protobuf:"varint,57,opt,name=public_field_me3,json=publicFieldMe3,proto3" json:"public_field_me3,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *NewsProto) Reset() {
	*x = NewsProto{}
	mi := &file_test_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NewsProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewsProto) ProtoMessage() {}

func (x *NewsProto) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewsProto.ProtoReflect.Descriptor instead.
func (*NewsProto) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{0}
}

func (x *NewsProto) GetInc() int64 {
	if x != nil {
		return x.Inc
	}
	return 0
}

func (x *NewsProto) GetInts() int64 {
	if x != nil {
		return x.Ints
	}
	return 0
}

func (x *NewsProto) GetInts8() int32 {
	if x != nil {
		return x.Ints8
	}
	return 0
}

func (x *NewsProto) GetInts16() int32 {
	if x != nil {
		return x.Ints16
	}
	return 0
}

func (x *NewsProto) GetInts32() int32 {
	if x != nil {
		return x.Ints32
	}
	return 0
}

func (x *NewsProto) GetInts64() int64 {
	if x != nil {
		return x.Ints64
	}
	return 0
}

func (x *NewsProto) GetUints() uint64 {
	if x != nil {
		return x.Uints
	}
	return 0
}

func (x *NewsProto) GetUints8() uint32 {
	if x != nil {
		return x.Uints8
	}
	return 0
}

func (x *NewsProto) GetUints16() uint32 {
	if x != nil {
		return x.Uints16
	}
	return 0
}

func (x *NewsProto) GetUints32() uint32 {
	if x != nil {
		return x.Uints32
	}
	return 0
}

func (x *NewsProto) GetUints64() uint64 {
	if x != nil {
		return x.Uints64
	}
	return 0
}

func (x *NewsProto) GetFloats32() float32 {
	if x != nil {
		return x.Floats32
	}
	return 0
}

func (x *NewsProto) GetFloats64() float64 {
	if x != nil {
		return x.Floats64
	}
	return 0
}

func (x *NewsProto) GetBools() bool {
	if x != nil {
		return x.Bools
	}
	return false
}

func (x *NewsProto) GetByte1() uint32 {
	if x != nil {
		return x.Byte1
	}
	return 0
}

func (x *NewsProto) GetBytes() []byte {
	if x != nil {
		return x.Bytes
	}
	return nil
}

func (x *NewsProto) GetListInts() []int64 {
	if x != nil {
		return x.ListInts
	}
	return nil
}

func (x *NewsProto) GetListString() []string {
	if x != nil {
		return x.ListString
	}
	return nil
}

func (x *NewsProto) GetListFloat() []float64 {
	if x != nil {
		return x.ListFloat
	}
	return nil
}

func (x *NewsProto) GetMapStringString() map[string]string {
	if x != nil {
		return x.MapStringString
	}
	return nil
}

func (x *NewsProto) GetMapStringBytes() map[string][]byte {
	if x != nil {
		return x.MapStringBytes
	}
	return nil
}

func (x *NewsProto) GetMapStringBool() map[string]bool {
	if x != nil {
		return x.MapStringBool
	}
	return nil
}

func (x *NewsProto) GetMapStringInt() map[string]int64 {
	if x != nil {
		return x.MapStringInt
	}
	return nil
}

func (x *NewsProto) GetMapStringFloat64() map[string]float64 {
	if x != nil {
		return x.MapStringFloat64
	}
	return nil
}

func (x *NewsProto) GetMapStringAny() map[string]*structpb.Value {
	if x != nil {
		return x.MapStringAny
	}
	return nil
}

func (x *NewsProto) GetMapIntString() map[int64]string {
	if x != nil {
		return x.MapIntString
	}
	return nil
}

func (x *NewsProto) GetMapIntInt() map[int64]int64 {
	if x != nil {
		return x.MapIntInt
	}
	return nil
}

func (x *NewsProto) GetMapIntBool() map[int64]bool {
	if x != nil {
		return x.MapIntBool
	}
	return nil
}

func (x *NewsProto) GetRenameSql() string {
	if x != nil {
		return x.RenameSql
	}
	return ""
}

func (x *NewsProto) GetRenameGo() string {
	if x != nil {
		return x.RenameGo
	}
	return ""
}

func (x *NewsProto) GetRenameJs() string {
	if x != nil {
		return x.RenameJs
	}
	return ""
}

func (x *NewsProto) GetMastUpperGo() string {
	if x != nil {
		return x.MastUpperGo
	}
	return ""
}

func (x *NewsProto) GetIntToSmallInt() int64 {
	if x != nil {
		return x.IntToSmallInt
	}
	return 0
}

func (x *NewsProto) GetSkip() string {
	if x != nil {
		return x.Skip
	}
	return ""
}

func (x *NewsProto) GetSqlUniqueU1_1() int64 {
	if x != nil {
		return x.SqlUniqueU1_1
	}
	return 0
}

func (x *NewsProto) GetSqlUniqueU1_2() int64 {
	if x != nil {
		return x.SqlUniqueU1_2
	}
	return 0
}

func (x *NewsProto) GetSqlIndex1_1() int64 {
	if x != nil {
		return x.SqlIndex1_1
	}
	return 0
}

func (x *NewsProto) GetSqlIndex1_2() int64 {
	if x != nil {
		return x.SqlIndex1_2
	}
	return 0
}

func (x *NewsProto) GetSqlIndex1_3() int64 {
	if x != nil {
		return x.SqlIndex1_3
	}
	return 0
}

func (x *NewsProto) GetSqlKeys_1() int64 {
	if x != nil {
		return x.SqlKeys_1
	}
	return 0
}

func (x *NewsProto) GetSqlKeys_2() int64 {
	if x != nil {
		return x.SqlKeys_2
	}
	return 0
}

func (x *NewsProto) GetSqlKeys_3() int64 {
	if x != nil {
		return x.SqlKeys_3
	}
	return 0
}

func (x *NewsProto) GetSqlSearch() string {
	if x != nil {
		return x.SqlSearch
	}
	return ""
}

func (x *NewsProto) GetSqlGet() string {
	if x != nil {
		return x.SqlGet
	}
	return ""
}

func (x *NewsProto) GetSqlUniqueX1() int64 {
	if x != nil {
		return x.SqlUniqueX1
	}
	return 0
}

func (x *NewsProto) GetSqlUniqueX2() int64 {
	if x != nil {
		return x.SqlUniqueX2
	}
	return 0
}

func (x *NewsProto) GetSqlUniqueX1X2() int64 {
	if x != nil {
		return x.SqlUniqueX1X2
	}
	return 0
}

func (x *NewsProto) GetSqlPrimary() float64 {
	if x != nil {
		return x.SqlPrimary
	}
	return 0
}

func (x *NewsProto) GetSqlJsonbIndex() map[string]*structpb.Value {
	if x != nil {
		return x.SqlJsonbIndex
	}
	return nil
}

func (x *NewsProto) GetTimeDuration() int64 {
	if x != nil {
		return x.TimeDuration
	}
	return 0
}

func (x *NewsProto) GetGoTypeIntToStrings() []string {
	if x != nil {
		return x.GoTypeIntToStrings
	}
	return nil
}

func (x *NewsProto) GetPublicField1() int64 {
	if x != nil {
		return x.PublicField1
	}
	return 0
}

func (x *NewsProto) GetPublicField2() int64 {
	if x != nil {
		return x.PublicField2
	}
	return 0
}

func (x *NewsProto) GetPublicField3() int64 {
	if x != nil {
		return x.PublicField3
	}
	return 0
}

func (x *NewsProto) GetPublicFieldMe1() int64 {
	if x != nil {
		return x.PublicFieldMe1
	}
	return 0
}

func (x *NewsProto) GetPublicFieldMe2() int64 {
	if x != nil {
		return x.PublicFieldMe2
	}
	return 0
}

func (x *NewsProto) GetPublicFieldMe3() int64 {
	if x != nil {
		return x.PublicFieldMe3
	}
	return 0
}

var File_test_proto protoreflect.FileDescriptor

const file_test_proto_rawDesc = "" +
	"\n" +
	"\n" +
	"test.proto\x12\x04news\x1a\x1cgoogle/protobuf/struct.proto\"\xe6\x16\n" +
	"\tNewsProto\x12\x10\n" +
	"\x03inc\x18\x01 \x01(\x03R\x03inc\x12\x12\n" +
	"\x04ints\x18\x02 \x01(\x03R\x04ints\x12\x14\n" +
	"\x05ints8\x18\x03 \x01(\x05R\x05ints8\x12\x16\n" +
	"\x06ints16\x18\x04 \x01(\x05R\x06ints16\x12\x16\n" +
	"\x06ints32\x18\x05 \x01(\x05R\x06ints32\x12\x16\n" +
	"\x06ints64\x18\x06 \x01(\x03R\x06ints64\x12\x14\n" +
	"\x05uints\x18\a \x01(\x04R\x05uints\x12\x16\n" +
	"\x06uints8\x18\b \x01(\rR\x06uints8\x12\x18\n" +
	"\auints16\x18\t \x01(\rR\auints16\x12\x18\n" +
	"\auints32\x18\n" +
	" \x01(\rR\auints32\x12\x18\n" +
	"\auints64\x18\v \x01(\x04R\auints64\x12\x1a\n" +
	"\bfloats32\x18\f \x01(\x02R\bfloats32\x12\x1a\n" +
	"\bfloats64\x18\r \x01(\x01R\bfloats64\x12\x14\n" +
	"\x05bools\x18\x0e \x01(\bR\x05bools\x12\x14\n" +
	"\x05byte1\x18\x0f \x01(\rR\x05byte1\x12\x14\n" +
	"\x05bytes\x18\x10 \x01(\fR\x05bytes\x12\x1b\n" +
	"\tlist_ints\x18\x11 \x03(\x03R\blistInts\x12\x1f\n" +
	"\vlist_string\x18\x12 \x03(\tR\n" +
	"listString\x12\x1d\n" +
	"\n" +
	"list_float\x18\x13 \x03(\x01R\tlistFloat\x12P\n" +
	"\x11map_string_string\x18\x14 \x03(\v2$.news.NewsProto.MapStringStringEntryR\x0fmapStringString\x12M\n" +
	"\x10map_string_bytes\x18\x15 \x03(\v2#.news.NewsProto.MapStringBytesEntryR\x0emapStringBytes\x12J\n" +
	"\x0fmap_string_bool\x18\x16 \x03(\v2\".news.NewsProto.MapStringBoolEntryR\rmapStringBool\x12G\n" +
	"\x0emap_string_int\x18\x17 \x03(\v2!.news.NewsProto.MapStringIntEntryR\fmapStringInt\x12S\n" +
	"\x12map_string_float64\x18\x18 \x03(\v2%.news.NewsProto.MapStringFloat64EntryR\x10mapStringFloat64\x12G\n" +
	"\x0emap_string_any\x18\x19 \x03(\v2!.news.NewsProto.MapStringAnyEntryR\fmapStringAny\x12G\n" +
	"\x0emap_int_string\x18\x1a \x03(\v2!.news.NewsProto.MapIntStringEntryR\fmapIntString\x12>\n" +
	"\vmap_int_int\x18\x1b \x03(\v2\x1e.news.NewsProto.MapIntIntEntryR\tmapIntInt\x12A\n" +
	"\fmap_int_bool\x18\x1c \x03(\v2\x1f.news.NewsProto.MapIntBoolEntryR\n" +
	"mapIntBool\x12\x1d\n" +
	"\n" +
	"rename_sql\x18\x1d \x01(\tR\trenameSql\x12\x1b\n" +
	"\trename_go\x18\x1e \x01(\tR\brenameGo\x12\x1b\n" +
	"\trename_js\x18\x1f \x01(\tR\brenameJs\x12\"\n" +
	"\rmast_upper_go\x18  \x01(\tR\vmastUpperGo\x12'\n" +
	"\x10int_to_small_int\x18! \x01(\x03R\rintToSmallInt\x12\x12\n" +
	"\x04skip\x18\" \x01(\tR\x04skip\x12%\n" +
	"\x0fsql_unique_u1_1\x18# \x01(\x03R\fsqlUniqueU11\x12%\n" +
	"\x0fsql_unique_u1_2\x18$ \x01(\x03R\fsqlUniqueU12\x12 \n" +
	"\fsql_index1_1\x18% \x01(\x03R\n" +
	"sqlIndex11\x12 \n" +
	"\fsql_index1_2\x18& \x01(\x03R\n" +
	"sqlIndex12\x12 \n" +
	"\fsql_index1_3\x18' \x01(\x03R\n" +
	"sqlIndex13\x12\x1c\n" +
	"\n" +
	"sql_keys_1\x18( \x01(\x03R\bsqlKeys1\x12\x1c\n" +
	"\n" +
	"sql_keys_2\x18) \x01(\x03R\bsqlKeys2\x12\x1c\n" +
	"\n" +
	"sql_keys_3\x18* \x01(\x03R\bsqlKeys3\x12\x1d\n" +
	"\n" +
	"sql_search\x18+ \x01(\tR\tsqlSearch\x12\x17\n" +
	"\asql_get\x18, \x01(\tR\x06sqlGet\x12\"\n" +
	"\rsql_unique_x1\x18- \x01(\x03R\vsqlUniqueX1\x12\"\n" +
	"\rsql_unique_x2\x18. \x01(\x03R\vsqlUniqueX2\x12'\n" +
	"\x10sql_unique_x1_x2\x18/ \x01(\x03R\rsqlUniqueX1X2\x12\x1f\n" +
	"\vsql_primary\x180 \x01(\x01R\n" +
	"sqlPrimary\x12J\n" +
	"\x0fsql_jsonb_index\x181 \x03(\v2\".news.NewsProto.SqlJsonbIndexEntryR\rsqlJsonbIndex\x12#\n" +
	"\rtime_duration\x182 \x01(\x03R\ftimeDuration\x122\n" +
	"\x16go_type_int_to_strings\x183 \x03(\tR\x12goTypeIntToStrings\x12#\n" +
	"\rpublic_field1\x184 \x01(\x03R\fpublicField1\x12#\n" +
	"\rpublic_field2\x185 \x01(\x03R\fpublicField2\x12#\n" +
	"\rpublic_field3\x186 \x01(\x03R\fpublicField3\x12(\n" +
	"\x10public_field_me1\x187 \x01(\x03R\x0epublicFieldMe1\x12(\n" +
	"\x10public_field_me2\x188 \x01(\x03R\x0epublicFieldMe2\x12(\n" +
	"\x10public_field_me3\x189 \x01(\x03R\x0epublicFieldMe3\x1aB\n" +
	"\x14MapStringStringEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01\x1aA\n" +
	"\x13MapStringBytesEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\fR\x05value:\x028\x01\x1a@\n" +
	"\x12MapStringBoolEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\bR\x05value:\x028\x01\x1a?\n" +
	"\x11MapStringIntEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\x03R\x05value:\x028\x01\x1aC\n" +
	"\x15MapStringFloat64Entry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\x01R\x05value:\x028\x01\x1aW\n" +
	"\x11MapStringAnyEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12,\n" +
	"\x05value\x18\x02 \x01(\v2\x16.google.protobuf.ValueR\x05value:\x028\x01\x1a?\n" +
	"\x11MapIntStringEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\x03R\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value:\x028\x01\x1a<\n" +
	"\x0eMapIntIntEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\x03R\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\x03R\x05value:\x028\x01\x1a=\n" +
	"\x0fMapIntBoolEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\x03R\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\bR\x05value:\x028\x01\x1aX\n" +
	"\x12SqlJsonbIndexEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12,\n" +
	"\x05value\x18\x02 \x01(\v2\x16.google.protobuf.ValueR\x05value:\x028\x01B\tZ\a./;newsb\x06proto3"

var (
	file_test_proto_rawDescOnce sync.Once
	file_test_proto_rawDescData []byte
)

func file_test_proto_rawDescGZIP() []byte {
	file_test_proto_rawDescOnce.Do(func() {
		file_test_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_test_proto_rawDesc), len(file_test_proto_rawDesc)))
	})
	return file_test_proto_rawDescData
}

var file_test_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_test_proto_goTypes = []any{
	(*NewsProto)(nil),      // 0: news.NewsProto
	nil,                    // 1: news.NewsProto.MapStringStringEntry
	nil,                    // 2: news.NewsProto.MapStringBytesEntry
	nil,                    // 3: news.NewsProto.MapStringBoolEntry
	nil,                    // 4: news.NewsProto.MapStringIntEntry
	nil,                    // 5: news.NewsProto.MapStringFloat64Entry
	nil,                    // 6: news.NewsProto.MapStringAnyEntry
	nil,                    // 7: news.NewsProto.MapIntStringEntry
	nil,                    // 8: news.NewsProto.MapIntIntEntry
	nil,                    // 9: news.NewsProto.MapIntBoolEntry
	nil,                    // 10: news.NewsProto.SqlJsonbIndexEntry
	(*structpb.Value)(nil), // 11: google.protobuf.Value
}
var file_test_proto_depIdxs = []int32{
	1,  // 0: news.NewsProto.map_string_string:type_name -> news.NewsProto.MapStringStringEntry
	2,  // 1: news.NewsProto.map_string_bytes:type_name -> news.NewsProto.MapStringBytesEntry
	3,  // 2: news.NewsProto.map_string_bool:type_name -> news.NewsProto.MapStringBoolEntry
	4,  // 3: news.NewsProto.map_string_int:type_name -> news.NewsProto.MapStringIntEntry
	5,  // 4: news.NewsProto.map_string_float64:type_name -> news.NewsProto.MapStringFloat64Entry
	6,  // 5: news.NewsProto.map_string_any:type_name -> news.NewsProto.MapStringAnyEntry
	7,  // 6: news.NewsProto.map_int_string:type_name -> news.NewsProto.MapIntStringEntry
	8,  // 7: news.NewsProto.map_int_int:type_name -> news.NewsProto.MapIntIntEntry
	9,  // 8: news.NewsProto.map_int_bool:type_name -> news.NewsProto.MapIntBoolEntry
	10, // 9: news.NewsProto.sql_jsonb_index:type_name -> news.NewsProto.SqlJsonbIndexEntry
	11, // 10: news.NewsProto.MapStringAnyEntry.value:type_name -> google.protobuf.Value
	11, // 11: news.NewsProto.SqlJsonbIndexEntry.value:type_name -> google.protobuf.Value
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

var _ = initNewsProto()

func initNewsProto() struct{} {
	file_test_proto_init()
	return struct{}{}
}
func file_test_proto_init() {
	if File_test_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_test_proto_rawDesc), len(file_test_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_test_proto_goTypes,
		DependencyIndexes: file_test_proto_depIdxs,
		MessageInfos:      file_test_proto_msgTypes,
	}.Build()
	File_test_proto = out.File
	file_test_proto_goTypes = nil
	file_test_proto_depIdxs = nil
}

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson34ac3833DecodeGithubComMonopollyJsonsgeneratorTest(in *jlexer.Lexer, out *News) {
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
		case "inc":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Inc = int(in.Int())
			}
		case "ints":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Ints = int(in.Int())
			}
		case "ints8":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Ints8 = int8(in.Int8())
			}
		case "ints16":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Ints16 = int16(in.Int16())
			}
		case "ints32":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Ints32 = int32(in.Int32())
			}
		case "ints64":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Ints64 = int64(in.Int64())
			}
		case "uints":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Uints = uint(in.Uint())
			}
		case "uints8":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Uints8 = uint8(in.Uint8())
			}
		case "uints16":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Uints16 = uint16(in.Uint16())
			}
		case "uints32":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Uints32 = uint32(in.Uint32())
			}
		case "uints64":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Uints64 = uint64(in.Uint64())
			}
		case "floats32":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Floats32 = float32(in.Float32())
			}
		case "floats64":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Floats64 = float64(in.Float64())
			}
		case "bools":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Bools = bool(in.Bool())
			}
		case "byte1":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Byte1 = uint8(in.Uint8())
			}
		case "bytes":
			if in.IsNull() {
				in.Skip()
				out.Bytes = nil
			} else {
				out.Bytes = in.Bytes()
			}
		case "list_ints":
			if in.IsNull() {
				in.Skip()
				out.List_ints = nil
			} else {
				in.Delim('[')
				if out.List_ints == nil {
					if !in.IsDelim(']') {
						out.List_ints = make([]int, 0, 8)
					} else {
						out.List_ints = []int{}
					}
				} else {
					out.List_ints = (out.List_ints)[:0]
				}
				for !in.IsDelim(']') {
					var v2 int
					if in.IsNull() {
						in.Skip()
					} else {
						v2 = int(in.Int())
					}
					out.List_ints = append(out.List_ints, v2)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "list_string":
			if in.IsNull() {
				in.Skip()
				out.List_string = nil
			} else {
				in.Delim('[')
				if out.List_string == nil {
					if !in.IsDelim(']') {
						out.List_string = make([]string, 0, 4)
					} else {
						out.List_string = []string{}
					}
				} else {
					out.List_string = (out.List_string)[:0]
				}
				for !in.IsDelim(']') {
					var v3 string
					if in.IsNull() {
						in.Skip()
					} else {
						v3 = string(in.String())
					}
					out.List_string = append(out.List_string, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "list_float":
			if in.IsNull() {
				in.Skip()
				out.List_float = nil
			} else {
				in.Delim('[')
				if out.List_float == nil {
					if !in.IsDelim(']') {
						out.List_float = make([]float64, 0, 8)
					} else {
						out.List_float = []float64{}
					}
				} else {
					out.List_float = (out.List_float)[:0]
				}
				for !in.IsDelim(']') {
					var v4 float64
					if in.IsNull() {
						in.Skip()
					} else {
						v4 = float64(in.Float64())
					}
					out.List_float = append(out.List_float, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "map_string_string":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Map_string_string = make(map[string]string)
				} else {
					out.Map_string_string = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v5 string
					if in.IsNull() {
						in.Skip()
					} else {
						v5 = string(in.String())
					}
					(out.Map_string_string)[key] = v5
					in.WantComma()
				}
				in.Delim('}')
			}
		case "map_string_bytes":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Map_string_bytes = make(map[string][]uint8)
				} else {
					out.Map_string_bytes = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v6 []uint8
					if in.IsNull() {
						in.Skip()
						v6 = nil
					} else {
						v6 = in.Bytes()
					}
					(out.Map_string_bytes)[key] = v6
					in.WantComma()
				}
				in.Delim('}')
			}
		case "map_string_bool":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Map_string_bool = make(map[string]bool)
				} else {
					out.Map_string_bool = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v8 bool
					if in.IsNull() {
						in.Skip()
					} else {
						v8 = bool(in.Bool())
					}
					(out.Map_string_bool)[key] = v8
					in.WantComma()
				}
				in.Delim('}')
			}
		case "map_string_int":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Map_string_int = make(map[string]int)
				} else {
					out.Map_string_int = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v9 int
					if in.IsNull() {
						in.Skip()
					} else {
						v9 = int(in.Int())
					}
					(out.Map_string_int)[key] = v9
					in.WantComma()
				}
				in.Delim('}')
			}
		case "map_string_float64":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Map_string_float64 = make(map[string]float64)
				} else {
					out.Map_string_float64 = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v10 float64
					if in.IsNull() {
						in.Skip()
					} else {
						v10 = float64(in.Float64())
					}
					(out.Map_string_float64)[key] = v10
					in.WantComma()
				}
				in.Delim('}')
			}
		case "map_string_any":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Map_string_any = make(map[string]interface{})
				} else {
					out.Map_string_any = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v11 interface{}
					if m, ok := v11.(easyjson.Unmarshaler); ok {
						m.UnmarshalEasyJSON(in)
					} else if m, ok := v11.(json.Unmarshaler); ok {
						_ = m.UnmarshalJSON(in.Raw())
					} else {
						v11 = in.Interface()
					}
					(out.Map_string_any)[key] = v11
					in.WantComma()
				}
				in.Delim('}')
			}
		case "map_int_string":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Map_int_string = make(map[int]string)
				} else {
					out.Map_int_string = nil
				}
				for !in.IsDelim('}') {
					key := int(in.IntStr())
					in.WantColon()
					var v12 string
					if in.IsNull() {
						in.Skip()
					} else {
						v12 = string(in.String())
					}
					(out.Map_int_string)[key] = v12
					in.WantComma()
				}
				in.Delim('}')
			}
		case "map_int_int":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Map_int_int = make(map[int]int)
				} else {
					out.Map_int_int = nil
				}
				for !in.IsDelim('}') {
					key := int(in.IntStr())
					in.WantColon()
					var v13 int
					if in.IsNull() {
						in.Skip()
					} else {
						v13 = int(in.Int())
					}
					(out.Map_int_int)[key] = v13
					in.WantComma()
				}
				in.Delim('}')
			}
		case "map_int_bool":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Map_int_bool = make(map[int]bool)
				} else {
					out.Map_int_bool = nil
				}
				for !in.IsDelim('}') {
					key := int(in.IntStr())
					in.WantColon()
					var v14 bool
					if in.IsNull() {
						in.Skip()
					} else {
						v14 = bool(in.Bool())
					}
					(out.Map_int_bool)[key] = v14
					in.WantComma()
				}
				in.Delim('}')
			}
		case "renameSQL":
			if in.IsNull() {
				in.Skip()
			} else {
				out.RenameSQL = string(in.String())
			}
		case "renameGO":
			if in.IsNull() {
				in.Skip()
			} else {
				out.RenameGO_OK = string(in.String())
			}
		case "renameJS_OK":
			if in.IsNull() {
				in.Skip()
			} else {
				out.RenameJS = string(in.String())
			}
		case "mast_upper_go":
			if in.IsNull() {
				in.Skip()
			} else {
				out.MAST_UPPER_GO = string(in.String())
			}
		case "intToSmallInt":
			if in.IsNull() {
				in.Skip()
			} else {
				out.IntToSmallInt = int(in.Int())
			}
		case "skip":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Skip = string(in.String())
			}
		case "sql_unique_u1_1":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Sql_unique_u1_1 = int(in.Int())
			}
		case "sql_unique_u1_2":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Sql_unique_u1_2 = int(in.Int())
			}
		case "sql_index1_1":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Sql_index1_1 = int(in.Int())
			}
		case "sql_index1_2":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Sql_index1_2 = int(in.Int())
			}
		case "sql_index1_3":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Sql_index1_3 = int(in.Int())
			}
		case "sql_keys_1":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Sql_keys_1 = int(in.Int())
			}
		case "sql_keys_2":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Sql_keys_2 = int(in.Int())
			}
		case "sql_keys_3":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Sql_keys_3 = int(in.Int())
			}
		case "sql_search":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Sql_search = string(in.String())
			}
		case "sql_get":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Sql_get = string(in.String())
			}
		case "sql_unique_x1":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Sql_unique_x1 = int(in.Int())
			}
		case "sql_unique_x2":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Sql_unique_x2 = int(in.Int())
			}
		case "sql_unique_x1_x2":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Sql_unique_x1_x2 = int(in.Int())
			}
		case "sql_primary":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Sql_primary = float64(in.Float64())
			}
		case "sql_jsonb_index":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				if !in.IsDelim('}') {
					out.Sql_jsonb_index = make(map[string]interface{})
				} else {
					out.Sql_jsonb_index = nil
				}
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v15 interface{}
					if m, ok := v15.(easyjson.Unmarshaler); ok {
						m.UnmarshalEasyJSON(in)
					} else if m, ok := v15.(json.Unmarshaler); ok {
						_ = m.UnmarshalJSON(in.Raw())
					} else {
						v15 = in.Interface()
					}
					(out.Sql_jsonb_index)[key] = v15
					in.WantComma()
				}
				in.Delim('}')
			}
		case "time_duration":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Time_duration = time.Duration(in.Int64())
			}
		case "go_type_int_to_strings":
			if in.IsNull() {
				in.Skip()
				out.Go_type_int_to_strings = nil
			} else {
				in.Delim('[')
				if out.Go_type_int_to_strings == nil {
					if !in.IsDelim(']') {
						out.Go_type_int_to_strings = make([]string, 0, 4)
					} else {
						out.Go_type_int_to_strings = []string{}
					}
				} else {
					out.Go_type_int_to_strings = (out.Go_type_int_to_strings)[:0]
				}
				for !in.IsDelim(']') {
					var v16 string
					if in.IsNull() {
						in.Skip()
					} else {
						v16 = string(in.String())
					}
					out.Go_type_int_to_strings = append(out.Go_type_int_to_strings, v16)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "public_field1":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Public_field1 = int(in.Int())
			}
		case "public_field2":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Public_field2 = int(in.Int())
			}
		case "public_field3":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Public_field3 = int(in.Int())
			}
		case "public_field_me1":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Public_field_me1 = int(in.Int())
			}
		case "public_field_me2":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Public_field_me2 = int(in.Int())
			}
		case "public_field_me3":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Public_field_me3 = int(in.Int())
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
func easyjson34ac3833EncodeGithubComMonopollyJsonsgeneratorTest(out *jwriter.Writer, in News) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Inc != 0 {
		const prefix string = ",\"inc\":"
		first = false
		out.RawString(prefix[1:])
		out.Int(int(in.Inc))
	}
	if in.Ints != 0 {
		const prefix string = ",\"ints\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Ints))
	}
	if in.Ints8 != 0 {
		const prefix string = ",\"ints8\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int8(int8(in.Ints8))
	}
	if in.Ints16 != 0 {
		const prefix string = ",\"ints16\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int16(int16(in.Ints16))
	}
	if in.Ints32 != 0 {
		const prefix string = ",\"ints32\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int32(int32(in.Ints32))
	}
	if in.Ints64 != 0 {
		const prefix string = ",\"ints64\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Ints64))
	}
	if in.Uints != 0 {
		const prefix string = ",\"uints\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint(uint(in.Uints))
	}
	if in.Uints8 != 0 {
		const prefix string = ",\"uints8\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint8(uint8(in.Uints8))
	}
	if in.Uints16 != 0 {
		const prefix string = ",\"uints16\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint16(uint16(in.Uints16))
	}
	if in.Uints32 != 0 {
		const prefix string = ",\"uints32\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint32(uint32(in.Uints32))
	}
	if in.Uints64 != 0 {
		const prefix string = ",\"uints64\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.Uints64))
	}
	if in.Floats32 != 0 {
		const prefix string = ",\"floats32\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Float32(float32(in.Floats32))
	}
	if in.Floats64 != 0 {
		const prefix string = ",\"floats64\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Float64(float64(in.Floats64))
	}
	if in.Bools {
		const prefix string = ",\"bools\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Bools))
	}
	if in.Byte1 != 0 {
		const prefix string = ",\"byte1\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint8(uint8(in.Byte1))
	}
	if len(in.Bytes) != 0 {
		const prefix string = ",\"bytes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Base64Bytes(in.Bytes)
	}
	if len(in.List_ints) != 0 {
		const prefix string = ",\"list_ints\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v19, v20 := range in.List_ints {
				if v19 > 0 {
					out.RawByte(',')
				}
				out.Int(int(v20))
			}
			out.RawByte(']')
		}
	}
	if len(in.List_string) != 0 {
		const prefix string = ",\"list_string\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v21, v22 := range in.List_string {
				if v21 > 0 {
					out.RawByte(',')
				}
				out.String(string(v22))
			}
			out.RawByte(']')
		}
	}
	if len(in.List_float) != 0 {
		const prefix string = ",\"list_float\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v23, v24 := range in.List_float {
				if v23 > 0 {
					out.RawByte(',')
				}
				out.Float64(float64(v24))
			}
			out.RawByte(']')
		}
	}
	if len(in.Map_string_string) != 0 {
		const prefix string = ",\"map_string_string\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v25First := true
			for v25Name, v25Value := range in.Map_string_string {
				if v25First {
					v25First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v25Name))
				out.RawByte(':')
				out.String(string(v25Value))
			}
			out.RawByte('}')
		}
	}
	if len(in.Map_string_bytes) != 0 {
		const prefix string = ",\"map_string_bytes\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v26First := true
			for v26Name, v26Value := range in.Map_string_bytes {
				if v26First {
					v26First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v26Name))
				out.RawByte(':')
				out.Base64Bytes(v26Value)
			}
			out.RawByte('}')
		}
	}
	if len(in.Map_string_bool) != 0 {
		const prefix string = ",\"map_string_bool\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v29First := true
			for v29Name, v29Value := range in.Map_string_bool {
				if v29First {
					v29First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v29Name))
				out.RawByte(':')
				out.Bool(bool(v29Value))
			}
			out.RawByte('}')
		}
	}
	if len(in.Map_string_int) != 0 {
		const prefix string = ",\"map_string_int\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v30First := true
			for v30Name, v30Value := range in.Map_string_int {
				if v30First {
					v30First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v30Name))
				out.RawByte(':')
				out.Int(int(v30Value))
			}
			out.RawByte('}')
		}
	}
	if len(in.Map_string_float64) != 0 {
		const prefix string = ",\"map_string_float64\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v31First := true
			for v31Name, v31Value := range in.Map_string_float64 {
				if v31First {
					v31First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v31Name))
				out.RawByte(':')
				out.Float64(float64(v31Value))
			}
			out.RawByte('}')
		}
	}
	if len(in.Map_string_any) != 0 {
		const prefix string = ",\"map_string_any\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v32First := true
			for v32Name, v32Value := range in.Map_string_any {
				if v32First {
					v32First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v32Name))
				out.RawByte(':')
				if m, ok := v32Value.(easyjson.Marshaler); ok {
					m.MarshalEasyJSON(out)
				} else if m, ok := v32Value.(json.Marshaler); ok {
					out.Raw(m.MarshalJSON())
				} else {
					out.Raw(json.Marshal(v32Value))
				}
			}
			out.RawByte('}')
		}
	}
	if len(in.Map_int_string) != 0 {
		const prefix string = ",\"map_int_string\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v33First := true
			for v33Name, v33Value := range in.Map_int_string {
				if v33First {
					v33First = false
				} else {
					out.RawByte(',')
				}
				out.IntStr(int(v33Name))
				out.RawByte(':')
				out.String(string(v33Value))
			}
			out.RawByte('}')
		}
	}
	if len(in.Map_int_int) != 0 {
		const prefix string = ",\"map_int_int\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v34First := true
			for v34Name, v34Value := range in.Map_int_int {
				if v34First {
					v34First = false
				} else {
					out.RawByte(',')
				}
				out.IntStr(int(v34Name))
				out.RawByte(':')
				out.Int(int(v34Value))
			}
			out.RawByte('}')
		}
	}
	if len(in.Map_int_bool) != 0 {
		const prefix string = ",\"map_int_bool\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v35First := true
			for v35Name, v35Value := range in.Map_int_bool {
				if v35First {
					v35First = false
				} else {
					out.RawByte(',')
				}
				out.IntStr(int(v35Name))
				out.RawByte(':')
				out.Bool(bool(v35Value))
			}
			out.RawByte('}')
		}
	}
	if in.RenameSQL != "" {
		const prefix string = ",\"renameSQL\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.RenameSQL))
	}
	if in.RenameGO_OK != "" {
		const prefix string = ",\"renameGO\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.RenameGO_OK))
	}
	if in.RenameJS != "" {
		const prefix string = ",\"renameJS_OK\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.RenameJS))
	}
	if in.MAST_UPPER_GO != "" {
		const prefix string = ",\"mast_upper_go\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.MAST_UPPER_GO))
	}
	if in.IntToSmallInt != 0 {
		const prefix string = ",\"intToSmallInt\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.IntToSmallInt))
	}
	if in.Skip != "" {
		const prefix string = ",\"skip\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Skip))
	}
	if in.Sql_unique_u1_1 != 0 {
		const prefix string = ",\"sql_unique_u1_1\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Sql_unique_u1_1))
	}
	if in.Sql_unique_u1_2 != 0 {
		const prefix string = ",\"sql_unique_u1_2\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Sql_unique_u1_2))
	}
	if in.Sql_index1_1 != 0 {
		const prefix string = ",\"sql_index1_1\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Sql_index1_1))
	}
	if in.Sql_index1_2 != 0 {
		const prefix string = ",\"sql_index1_2\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Sql_index1_2))
	}
	if in.Sql_index1_3 != 0 {
		const prefix string = ",\"sql_index1_3\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Sql_index1_3))
	}
	if in.Sql_keys_1 != 0 {
		const prefix string = ",\"sql_keys_1\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Sql_keys_1))
	}
	if in.Sql_keys_2 != 0 {
		const prefix string = ",\"sql_keys_2\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Sql_keys_2))
	}
	if in.Sql_keys_3 != 0 {
		const prefix string = ",\"sql_keys_3\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Sql_keys_3))
	}
	if in.Sql_search != "" {
		const prefix string = ",\"sql_search\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Sql_search))
	}
	if in.Sql_get != "" {
		const prefix string = ",\"sql_get\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Sql_get))
	}
	if in.Sql_unique_x1 != 0 {
		const prefix string = ",\"sql_unique_x1\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Sql_unique_x1))
	}
	if in.Sql_unique_x2 != 0 {
		const prefix string = ",\"sql_unique_x2\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Sql_unique_x2))
	}
	if in.Sql_unique_x1_x2 != 0 {
		const prefix string = ",\"sql_unique_x1_x2\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Sql_unique_x1_x2))
	}
	if in.Sql_primary != 0 {
		const prefix string = ",\"sql_primary\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Float64(float64(in.Sql_primary))
	}
	if len(in.Sql_jsonb_index) != 0 {
		const prefix string = ",\"sql_jsonb_index\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('{')
			v36First := true
			for v36Name, v36Value := range in.Sql_jsonb_index {
				if v36First {
					v36First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v36Name))
				out.RawByte(':')
				if m, ok := v36Value.(easyjson.Marshaler); ok {
					m.MarshalEasyJSON(out)
				} else if m, ok := v36Value.(json.Marshaler); ok {
					out.Raw(m.MarshalJSON())
				} else {
					out.Raw(json.Marshal(v36Value))
				}
			}
			out.RawByte('}')
		}
	}
	if in.Time_duration != 0 {
		const prefix string = ",\"time_duration\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Time_duration))
	}
	if len(in.Go_type_int_to_strings) != 0 {
		const prefix string = ",\"go_type_int_to_strings\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v37, v38 := range in.Go_type_int_to_strings {
				if v37 > 0 {
					out.RawByte(',')
				}
				out.String(string(v38))
			}
			out.RawByte(']')
		}
	}
	if in.Public_field1 != 0 {
		const prefix string = ",\"public_field1\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Public_field1))
	}
	if in.Public_field2 != 0 {
		const prefix string = ",\"public_field2\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Public_field2))
	}
	if in.Public_field3 != 0 {
		const prefix string = ",\"public_field3\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Public_field3))
	}
	if in.Public_field_me1 != 0 {
		const prefix string = ",\"public_field_me1\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Public_field_me1))
	}
	if in.Public_field_me2 != 0 {
		const prefix string = ",\"public_field_me2\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Public_field_me2))
	}
	if in.Public_field_me3 != 0 {
		const prefix string = ",\"public_field_me3\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Public_field_me3))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v News) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson34ac3833EncodeGithubComMonopollyJsonsgeneratorTest(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v News) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson34ac3833EncodeGithubComMonopollyJsonsgeneratorTest(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *News) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson34ac3833DecodeGithubComMonopollyJsonsgeneratorTest(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *News) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson34ac3833DecodeGithubComMonopollyJsonsgeneratorTest(l, v)
}

// easyjson marshal
func (a *News) Marshal() []byte {
	b, _ := easyjson.Marshal(a)
	return b
}

// easyjson unmarshal
func NewsUnmarshal(src []byte) (a *News) {
	a = new(News)
	err := easyjson.Unmarshal(src, a)
	if err != nil {
		return nil
	}
	return
}
