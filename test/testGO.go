package news

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	jsoniter "github.com/json-iterator/go"
	"github.com/monopolly/cast"
	"github.com/monopolly/jsons"
	"github.com/niubaoshu/gotiny"
	"github.com/vmihailenco/msgpack/v5"
)

/*
	Auto-generate. Do not change!
	Help to auto-generated correct models for golang structures.
	Sergey Keplin (c) 2026
	a@senthy.com, lava.mobi@gmail.com, https://t.me/martinprestone
	github.com/monopolly/jsons
	Ex: //! [] go=News js=NewsJson sql=news noinit up...

	STRUCT:
	demo                          Generate demo json file with default values
	chengine=MergeTree            Optional, mergeTree by default
	debug                         Debug mode
	gotiny                        Create gotiny marshal/unmarshal
	noinit                        No New() init function for struct
	js=NewsJson                   Change json struct names. Ex: js=NewsJson
	enum                          Generate swift enum
	lock                          Lock model for generation. Can't change model.
	ts=news                       Set typescript struct names. Ex: ts=NewsJson
	noprefix                      Generate simple index IndexID instead IndexNewsID
	ch=views                      Set clickhouse sql table name. Ex: ch=views
	sql=news                      Set sql table name. Ex: sql=accounts
	msgp                          Create message pack marshal/unmarshal
	!omit                         No omit tag for json
	go=News                       Set golang struct names. Ex: go=News1 > type News1 struct{}
	swift                         Generate swift model

	Field:
	title{}                       Add custom title for index. title{Nice & Sweet}
	desc{}                        Add custom desc for index. desc{This field for success}

	GO:
	desc{}                        Add field user title desc{Use it nice}
	#                             Add lists for fields. Ex: id int //#readonly #must...
	nofunc                        Do not create any jsons functions for fiels
	up                            Make uppercase for functions
	must                          Create one validation function for all must fields
	type=""                       Replace golang type. Ex: type="[]*News"
	name=""                       Replace golang struct name. Ex: name="NewsList"
	title=""                      Add custom title for index
	req{}                         Add user required fields
	title{}                       Add field user title title{Nice}
	desc=""                       Add custom desc for index

	SQL:
	replace="bigint primary key"  Rewrite sql for field. Ex: replace="bigint primary key"
	add="primary key"             Append sql for field. Ex: replace="primary key"
	unique="groupname"            Add unique fields constrains by group (you need set group name, then generator join fields). Ex: unique="group1" unique="group2"
	renames="oldname"             Create an alter table record in SQL file. Rename table column.
	noinsert                      Do use field for insert function
	index                         Create simple default index or gin for jsonb
	altertable                    Add field line Alter table to SQL file with current time comment
	defaults                      Add default value based on field type
	unix                          Add default value: extract(epoch from now())
	type="jsonb"                  Rewrite sql type for field. Ex: type="jsonb"
	ver="v4"                      Create an alter table record in SQL file. Add new column. With new version in comment. Ex: ver="2"
	idx                           Add index fields by group name. Ex: idx="nameIndex" idx="credsIndex"
	search                        Add tsvector index by group. Ex: search="tsv": tsv tsvector GENERATED ALWAYS AS (to_tsvector('simple', title || ' ' || brand)). For search: SELECT brand, title FROM assets WHERE search @@ to_tsquery('english', 'f8');
	skip                          Skip field for sql queries

	Clickhouse:

	SWIFT:
	must                          Swift field with required values
	skip                          Ignore field for Swift
	type=""                       Replace type for Swift model. Ex: type="string"
	file=""                       Create another swift file for this field Swift model. file="f1", file="f2"

	JSON:
	name=""                       Replace json field name. Ex: name="sid"
	skip                          Skip json field
	inc                           Add inc jsons function for numbers fields
	bool                          Add set jsons function for bool fields
	time                          Create convert jsons function for unixtime fields
	raw                           Set Raw json function inside field
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
)

// index func
func NewsIndexes() []NewsIndexType {
	return []NewsIndexType{IndexInc, IndexInts, IndexInts8, IndexInts16, IndexInts32, IndexInts64, IndexUints, IndexUints8, IndexUints16, IndexUints32, IndexUints64, IndexFloats32, IndexFloats64, IndexBools, IndexByte1, IndexBytes, IndexList_ints, IndexList_string, IndexList_float, IndexMap_string_string, IndexMap_string_bytes, IndexMap_string_bool, IndexMap_string_int, IndexMap_string_float64, IndexMap_string_any, IndexMap_int_string, IndexMap_int_int, IndexMap_int_bool, IndexRenameSQL, IndexRenameGO_OK, IndexRenameJS, IndexMAST_UPPER_GO, IndexIntToSmallInt, IndexSkip, IndexSql_unique_u1_1, IndexSql_unique_u1_2, IndexSql_index1_1, IndexSql_index1_2, IndexSql_index1_3, IndexSql_keys_1, IndexSql_keys_2, IndexSql_keys_3, IndexSql_search, IndexSql_get, IndexSql_unique_x1, IndexSql_unique_x2, IndexSql_unique_x1_x2, IndexSql_primary, IndexSql_jsonb_index, IndexTime_duration, IndexGo_type_int_to_strings}

}

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
		}
	}
	return
}

// Tuple create an array from struct
func (a *News) Tuple() (r []any) {
	return []any{a.Inc, a.Ints, a.Ints8, a.Ints16, a.Ints32, a.Ints64, a.Uints, a.Uints8, a.Uints16, a.Uints32, a.Uints64, a.Floats32, a.Floats64, a.Bools, a.Byte1, a.Bytes, a.List_ints, a.List_string, a.List_float, a.Map_string_string, a.Map_string_bytes, a.Map_string_bool, a.Map_string_int, a.Map_string_float64, a.Map_string_any, a.Map_int_string, a.Map_int_int, a.Map_int_bool, a.RenameSQL, a.RenameGO_OK, a.RenameJS, a.MAST_UPPER_GO, a.IntToSmallInt, a.Skip, a.Sql_unique_u1_1, a.Sql_unique_u1_2, a.Sql_index1_1, a.Sql_index1_2, a.Sql_index1_3, a.Sql_keys_1, a.Sql_keys_2, a.Sql_keys_3, a.Sql_search, a.Sql_get, a.Sql_unique_x1, a.Sql_unique_x2, a.Sql_unique_x1_x2, a.Sql_primary, a.Sql_jsonb_index, a.Time_duration, a.Go_type_int_to_strings}
}

// Tuple create an array from struct
func (a *News) sqlTuple() (r []any) {
	return []any{a.Ints, a.Ints8, a.Ints16, a.Ints32, a.Ints64, a.Uints, a.Uints8, a.Uints16, a.Uints32, a.Uints64, a.Floats32, a.Floats64, a.Bools, a.Byte1, a.Bytes, a.List_ints, a.List_string, a.List_float, a.Map_string_string, a.Map_string_bytes, a.Map_string_bool, a.Map_string_int, a.Map_string_float64, a.Map_string_any, a.Map_int_string, a.Map_int_int, a.Map_int_bool, a.RenameSQL, a.RenameGO_OK, a.RenameJS, a.MAST_UPPER_GO, a.IntToSmallInt, a.Sql_unique_u1_1, a.Sql_unique_u1_2, a.Sql_index1_1, a.Sql_index1_2, a.Sql_index1_3, a.Sql_keys_1, a.Sql_keys_2, a.Sql_keys_3, a.Sql_search, a.Sql_get, a.Sql_unique_x1, a.Sql_unique_x2, a.Sql_unique_x1_x2, a.Sql_primary, a.Sql_jsonb_index, a.Time_duration, a.Go_type_int_to_strings}
}

// Tuple create an array from struct
func (a *News) sqlAllTuples() (r []any) {
	return []any{a.Inc, a.Ints, a.Ints8, a.Ints16, a.Ints32, a.Ints64, a.Uints, a.Uints8, a.Uints16, a.Uints32, a.Uints64, a.Floats32, a.Floats64, a.Bools, a.Byte1, a.Bytes, a.List_ints, a.List_string, a.List_float, a.Map_string_string, a.Map_string_bytes, a.Map_string_bool, a.Map_string_int, a.Map_string_float64, a.Map_string_any, a.Map_int_string, a.Map_int_int, a.Map_int_bool, a.RenameSQL, a.RenameGO_OK, a.RenameJS, a.MAST_UPPER_GO, a.IntToSmallInt, a.Sql_unique_u1_1, a.Sql_unique_u1_2, a.Sql_index1_1, a.Sql_index1_2, a.Sql_index1_3, a.Sql_keys_1, a.Sql_keys_2, a.Sql_keys_3, a.Sql_search, a.Sql_get, a.Sql_unique_x1, a.Sql_unique_x2, a.Sql_unique_x1_x2, a.Sql_primary, a.Sql_jsonb_index, a.Time_duration, a.Go_type_int_to_strings}
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
		Add(FieldGo_type_int_to_strings, a.Go_type_int_to_strings)
	return js.Bytes()
}

func NewsReadonlyList() []NewsIndexType {
	return []NewsIndexType{IndexInc}
}

func (a NewsIndexType) Readonly() bool {
	switch a {
	case IndexInc:
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
	default:
		return 0
	}
}

// valid struct key check
func NewsValidKey(key string) bool {
	switch key {
	case "inc", "ints", "ints8", "ints16", "ints32", "ints64", "uints", "uints8", "uints16", "uints32", "uints64", "floats32", "floats64", "bools", "byte1", "bytes", "list_ints", "list_string", "list_float", "map_string_string", "map_string_bytes", "map_string_bool", "map_string_int", "map_string_float64", "map_string_any", "map_int_string", "map_int_int", "map_int_bool", "renameSQL", "renameGO", "renameJS", "mast_upper_go", "intToSmallInt", "skip", "sql_unique_u1_1", "sql_unique_u1_2", "sql_index1_1", "sql_index1_2", "sql_index1_3", "sql_keys_1", "sql_keys_2", "sql_keys_3", "sql_search", "sql_get", "sql_unique_x1", "sql_unique_x2", "sql_unique_x1_x2", "sql_primary", "sql_jsonb_index", "time_duration", "go_type_int_to_strings":
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
	}
}

// struct to map
func (a *News) Iterate(f func(k NewsIndexType, v any)) {
	for _, x := range NewsIndexes() {
		f(x, a.Get(x.String()))
	}
}

// gotiny marshal
func (a *News) Gotiny() []byte {
	return gotiny.Marshal(&a.Inc, &a.Ints, &a.Ints8, &a.Ints16, &a.Ints32, &a.Ints64, &a.Uints, &a.Uints8, &a.Uints16, &a.Uints32, &a.Uints64, &a.Floats32, &a.Floats64, &a.Bools, &a.Byte1, &a.Bytes, &a.List_ints, &a.List_string, &a.List_float, &a.Map_string_string, &a.Map_string_bytes, &a.Map_string_bool, &a.Map_string_int, &a.Map_string_float64, &a.Map_string_any, &a.Map_int_string, &a.Map_int_int, &a.Map_int_bool, &a.RenameSQL, &a.RenameGO_OK, &a.RenameJS, &a.MAST_UPPER_GO, &a.IntToSmallInt, &a.Skip, &a.Sql_unique_u1_1, &a.Sql_unique_u1_2, &a.Sql_index1_1, &a.Sql_index1_2, &a.Sql_index1_3, &a.Sql_keys_1, &a.Sql_keys_2, &a.Sql_keys_3, &a.Sql_search, &a.Sql_get, &a.Sql_unique_x1, &a.Sql_unique_x2, &a.Sql_unique_x1_x2, &a.Sql_primary, &a.Sql_jsonb_index, &a.Time_duration, &a.Go_type_int_to_strings)
}

// parse gotiny
func ParseNewsGotiny(v []byte) (a News) {
	gotiny.Unmarshal(v, &a.Inc, &a.Ints, &a.Ints8, &a.Ints16, &a.Ints32, &a.Ints64, &a.Uints, &a.Uints8, &a.Uints16, &a.Uints32, &a.Uints64, &a.Floats32, &a.Floats64, &a.Bools, &a.Byte1, &a.Bytes, &a.List_ints, &a.List_string, &a.List_float, &a.Map_string_string, &a.Map_string_bytes, &a.Map_string_bool, &a.Map_string_int, &a.Map_string_float64, &a.Map_string_any, &a.Map_int_string, &a.Map_int_int, &a.Map_int_bool, &a.RenameSQL, &a.RenameGO_OK, &a.RenameJS, &a.MAST_UPPER_GO, &a.IntToSmallInt, &a.Skip, &a.Sql_unique_u1_1, &a.Sql_unique_u1_2, &a.Sql_index1_1, &a.Sql_index1_2, &a.Sql_index1_3, &a.Sql_keys_1, &a.Sql_keys_2, &a.Sql_keys_3, &a.Sql_search, &a.Sql_get, &a.Sql_unique_x1, &a.Sql_unique_x2, &a.Sql_unique_x1_x2, &a.Sql_primary, &a.Sql_jsonb_index, &a.Time_duration, &a.Go_type_int_to_strings)
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

// sql NewsSQL class
type NewsSQL struct{ pool *pgxpool.Pool }

func NewNewsSQL(pool *pgxpool.Pool) (a *NewsSQL) {
	a = new(NewsSQL)
	a.pool = pool
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

	q := "insert into news (ints, ints8, ints16, ints32, ints64, uints, uints8, uints16, uints32, uints64, floats32, floats64, bools, byte1, bytes, list_ints, list_string, list_float, map_string_string, map_string_bytes, map_string_bool, map_string_int, map_string_float64, map_string_any, map_int_string, map_int_int, map_int_bool, renameSQL_OK, renameGO, renameJS, mast_upper_go, intToSmallInt, sql_unique_u1_1, sql_unique_u1_2, sql_index1_1, sql_index1_2, sql_index1_3, sql_keys_1, sql_keys_2, sql_keys_3, sql_search, sql_get, sql_unique_x1, sql_unique_x2, sql_unique_x1_x2, sql_primary, sql_jsonb_index, time_duration, go_type_int_to_strings) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49) returning inc"
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

	q := "insert into news (ints, ints8, ints16, ints32, ints64, uints, uints8, uints16, uints32, uints64, floats32, floats64, bools, byte1, bytes, list_ints, list_string, list_float, map_string_string, map_string_bytes, map_string_bool, map_string_int, map_string_float64, map_string_any, map_int_string, map_int_int, map_int_bool, renameSQL_OK, renameGO, renameJS, mast_upper_go, intToSmallInt, sql_unique_u1_1, sql_unique_u1_2, sql_index1_1, sql_index1_2, sql_index1_3, sql_keys_1, sql_keys_2, sql_keys_3, sql_search, sql_get, sql_unique_x1, sql_unique_x2, sql_unique_x1_x2, sql_primary, sql_jsonb_index, time_duration, go_type_int_to_strings) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49) on conflict do nothing returning inc"
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

	q := "insert into news (inc, ints, ints8, ints16, ints32, ints64, uints, uints8, uints16, uints32, uints64, floats32, floats64, bools, byte1, bytes, list_ints, list_string, list_float, map_string_string, map_string_bytes, map_string_bool, map_string_int, map_string_float64, map_string_any, map_int_string, map_int_int, map_int_bool, renameSQL_OK, renameGO, renameJS, mast_upper_go, intToSmallInt, sql_unique_u1_1, sql_unique_u1_2, sql_index1_1, sql_index1_2, sql_index1_3, sql_keys_1, sql_keys_2, sql_keys_3, sql_search, sql_get, sql_unique_x1, sql_unique_x2, sql_unique_x1_x2, sql_primary, sql_jsonb_index, time_duration, go_type_int_to_strings) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50)"
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
	rows.Next()
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
	rows.Next()
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
		count++
		andlist = append(andlist, fmt.Sprintf("%s = $%d", k, count))
		values = append(values, v)
	}
	keys := strings.Join(andlist, " and ")

	q := fmt.Sprintf("select %s from %s where %s limit 1", fieldlist, a.TableName(), keys)
	res = new(News)
	rows, err := conn.Query(c, q, values...)
	if err != nil {
		return
	}
	defer rows.Close()
	rows.Next()
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
	q := fmt.Sprintf("update news set %s = $1 where inc = $2", k)
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
	return a.Conn(func(conn *pgxpool.Conn) (err error) {
		q := fmt.Sprintf("update %s set %s = $1 where %s", a.TableName(), k, where)
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
	q := fmt.Sprintf("update news set uints16 = uints16 || $$%s$$::jsonb where id = $1", res)
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

	q := "update news set uints16 = uints16 - $1 where id = $2"
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

	q := "update news set uints16 = uints16 - $1 || jsonb_build_object($2, uints16->$1) where id = $3"
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

	q := fmt.Sprintf("update news set list_ints = list_ints - $1 where inc = $2")
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

	q := fmt.Sprintf("update news set list_string = list_string - $1 where inc = $2")
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

	q := fmt.Sprintf("update news set list_float = list_float - $1 where inc = $2")
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
	q := fmt.Sprintf("update news set map_string_string = map_string_string || $$%s$$::jsonb where id = $1", res)
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

	q := "update news set map_string_string = map_string_string - $1 where id = $2"
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

	q := "update news set map_string_string = map_string_string - $1 || jsonb_build_object($2, map_string_string->$1) where id = $3"
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
	q := fmt.Sprintf("update news set map_string_bytes = map_string_bytes || $$%s$$::jsonb where id = $1", res)
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

	q := "update news set map_string_bytes = map_string_bytes - $1 where id = $2"
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

	q := "update news set map_string_bytes = map_string_bytes - $1 || jsonb_build_object($2, map_string_bytes->$1) where id = $3"
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
	q := fmt.Sprintf("update news set map_string_bool = map_string_bool || $$%s$$::jsonb where id = $1", res)
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

	q := "update news set map_string_bool = map_string_bool - $1 where id = $2"
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

	q := "update news set map_string_bool = map_string_bool - $1 || jsonb_build_object($2, map_string_bool->$1) where id = $3"
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
	q := fmt.Sprintf("update news set map_string_int = map_string_int || $$%s$$::jsonb where id = $1", res)
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

	q := "update news set map_string_int = map_string_int - $1 where id = $2"
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

	q := "update news set map_string_int = map_string_int - $1 || jsonb_build_object($2, map_string_int->$1) where id = $3"
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
	q := fmt.Sprintf("update news set map_string_float64 = map_string_float64 || $$%s$$::jsonb where id = $1", res)
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

	q := "update news set map_string_float64 = map_string_float64 - $1 where id = $2"
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

	q := "update news set map_string_float64 = map_string_float64 - $1 || jsonb_build_object($2, map_string_float64->$1) where id = $3"
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
	q := fmt.Sprintf("update news set map_string_any = map_string_any || $$%s$$::jsonb where id = $1", res)
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

	q := "update news set map_string_any = map_string_any - $1 where id = $2"
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

	q := "update news set map_string_any = map_string_any - $1 || jsonb_build_object($2, map_string_any->$1) where id = $3"
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
	q := fmt.Sprintf("update news set map_int_string = map_int_string || $$%s$$::jsonb where id = $1", res)
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

	q := "update news set map_int_string = map_int_string - $1 where id = $2"
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

	q := "update news set map_int_string = map_int_string - $1 || jsonb_build_object($2, map_int_string->$1) where id = $3"
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
	q := fmt.Sprintf("update news set map_int_int = map_int_int || $$%s$$::jsonb where id = $1", res)
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

	q := "update news set map_int_int = map_int_int - $1 where id = $2"
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

	q := "update news set map_int_int = map_int_int - $1 || jsonb_build_object($2, map_int_int->$1) where id = $3"
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
	q := fmt.Sprintf("update news set map_int_bool = map_int_bool || $$%s$$::jsonb where id = $1", res)
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

	q := "update news set map_int_bool = map_int_bool - $1 where id = $2"
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

	q := "update news set map_int_bool = map_int_bool - $1 || jsonb_build_object($2, map_int_bool->$1) where id = $3"
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
	q := fmt.Sprintf("update news set sql_jsonb_index = sql_jsonb_index || $$%s$$::jsonb where id = $1", res)
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

	q := "update news set sql_jsonb_index = sql_jsonb_index - $1 where id = $2"
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

	q := "update news set sql_jsonb_index = sql_jsonb_index - $1 || jsonb_build_object($2, sql_jsonb_index->$1) where id = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
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
		fields = []NewsIndexType{IndexInc, IndexInts, IndexInts8, IndexInts16, IndexInts32, IndexInts64, IndexUints, IndexUints8, IndexUints16, IndexUints32, IndexUints64, IndexFloats32, IndexFloats64, IndexBools, IndexByte1, IndexBytes, IndexList_ints, IndexList_string, IndexList_float, IndexMap_string_string, IndexMap_string_bytes, IndexMap_string_bool, IndexMap_string_int, IndexMap_string_float64, IndexMap_string_any, IndexMap_int_string, IndexMap_int_int, IndexMap_int_bool, IndexRenameSQL, IndexRenameGO_OK, IndexRenameJS, IndexMAST_UPPER_GO, IndexIntToSmallInt, IndexSkip, IndexSql_unique_u1_1, IndexSql_unique_u1_2, IndexSql_index1_1, IndexSql_index1_2, IndexSql_index1_3, IndexSql_keys_1, IndexSql_keys_2, IndexSql_keys_3, IndexSql_search, IndexSql_get, IndexSql_unique_x1, IndexSql_unique_x2, IndexSql_unique_x1_x2, IndexSql_primary, IndexSql_jsonb_index, IndexTime_duration, IndexGo_type_int_to_strings}
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
		fieldsStrings = append(fieldsStrings, x.String())
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
		switch p.Type() {
		case "bool":
			switch cast.Bool(v) {
			case true:
				andlist = append(andlist, k)
			case false:
				andlist = append(andlist, fmt.Sprintf("(not %[1]s or %[1]s is null)", k))
			}

		default:
			count++
			andlist = append(andlist, fmt.Sprintf("%s = $%d", k, count))
			values = append(values, v)
		}
	}

	// GT where
	for k, v := range a.GT {

		if !NewsValidKey(k) {
			continue
		}
		count++
		andlist = append(andlist, fmt.Sprintf("%s > $%d", k, count))
		values = append(values, v)
	}

	// LT where
	for k, v := range a.LT {
		if !NewsValidKey(k) {
			continue
		}
		count++
		andlist = append(andlist, fmt.Sprintf("%s < $%d", k, count))
		values = append(values, v)
	}

	// NOT where
	for k, v := range a.NOT {
		if !NewsValidKey(k) {
			continue
		}

		for _, x := range v {
			count++
			andlist = append(andlist, fmt.Sprintf("%s != $%d", k, count))
			values = append(values, x)
		}

	}

	// LIKE where
	for k, v := range a.Like {
		if !NewsValidKey(k) {
			continue
		}
		count++
		andlist = append(andlist, fmt.Sprintf("%s ilike $%d", k, count))
		values = append(values, "%"+v+"%")
	}

	// IN where
	if a.IN != nil {
		for k, v := range a.IN {
			if !NewsValidKey(k) {
				continue
			}
			var inlist []string
			for _, num := range v {
				inlist = append(inlist, fmt.Sprint(num))
			}
			count++
			andlist = append(andlist, fmt.Sprintf("%s in (%s)", k, strings.Join(inlist, ",")))
			// values = append(values, string)
		}
	}

	// INSQ where
	if a.INSQL != nil {
		for k, v := range a.INSQL {
			if !NewsValidKey(k) {
				continue
			}
			count++
			andlist = append(andlist, fmt.Sprintf("%s in (%s)", k, v))
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
		list = append(list, fmt.Sprintf("order by %s", a.Sort))
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

		for pos, x := range fields {
			item.Update(x.String(), v[pos])
		}
		res = append(res, &item)
	}

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
	q = q + "from news where " + where
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

		for pos, x := range fields {
			item.Update(x.String(), v[pos])
		}
		res = append(res, &item)
	}
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

	q := fmt.Sprintf("select exists (select id from news where %s = $1 limit 1)", field.SQLName())
	err = conn.QueryRow(c, q, v).Scan(&has)
	return
}

// Create table
func (a *NewsSQL) CreateTable() (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := `create table if not exists news (
	inc                        bigserial primary key,
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
	go_type_int_to_strings     bigint,
	unique(sql_unique_x1,sql_unique_x1_x2),
	unique(sql_unique_x2,sql_unique_x1_x2),
	unique(sql_unique_u1_2,sql_unique_u1_1),
	primary key (sql_primary)
)
`
	_, err = conn.Exec(c, q)
	return
}

// parse sql query
func (a *NewsSQL) Fields() (res []NewsIndexType) {
	return []NewsIndexType{IndexInc, IndexInts, IndexInts8, IndexInts16, IndexInts32, IndexInts64, IndexUints, IndexUints8, IndexUints16, IndexUints32, IndexUints64, IndexFloats32, IndexFloats64, IndexBools, IndexByte1, IndexBytes, IndexList_ints, IndexList_string, IndexList_float, IndexMap_string_string, IndexMap_string_bytes, IndexMap_string_bool, IndexMap_string_int, IndexMap_string_float64, IndexMap_string_any, IndexMap_int_string, IndexMap_int_int, IndexMap_int_bool, IndexRenameSQL, IndexRenameGO_OK, IndexRenameJS, IndexMAST_UPPER_GO, IndexIntToSmallInt, IndexSkip, IndexSql_unique_u1_1, IndexSql_unique_u1_2, IndexSql_index1_1, IndexSql_index1_2, IndexSql_index1_3, IndexSql_keys_1, IndexSql_keys_2, IndexSql_keys_3, IndexSql_search, IndexSql_get, IndexSql_unique_x1, IndexSql_unique_x2, IndexSql_unique_x1_x2, IndexSql_primary, IndexSql_jsonb_index, IndexTime_duration, IndexGo_type_int_to_strings}
}
