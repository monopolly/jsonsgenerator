package news

import (
	"time"
)

// go=News sql=news ts=news swift demo enum=Category noprefix msgp gotiny
type news struct {
	inc int // sql{inc}  #readonly

	// int types
	ints   int
	ints8  int8
	ints16 int16
	ints32 int32
	ints64 int64

	// uint types
	uints   uint
	uints8  uint8
	uints16 uint16
	uints32 uint32
	uints64 uint64

	// float types
	floats32 float32
	floats64 float64

	// bool
	bools bool

	// bytes
	byte1 byte
	bytes []byte

	// ints types
	list_ints   []int
	list_string []string
	list_float  []float64

	// map string
	map_string_string  map[string]string
	map_string_bytes   map[string][]byte
	map_string_bool    map[string]bool
	map_string_int     map[string]int
	map_string_float64 map[string]float64
	map_string_any     map[string]any

	// map int
	map_int_string map[int]string
	map_int_int    map[int]int
	map_int_bool   map[int]bool

	// rename
	renameSQL string // sql{name="renameSQL_OK"}
	renameGO  string // go{name="RenameGO_OK"}
	renameJS  string // js{name="renameJS_OK"}

	// upper
	mast_upper_go string // go{up}

	// upper
	intToSmallInt int // sql{type="smallint"} rename

	// skipped
	skip string //sql{skip} swift{skip} js{skip}

	// sql unique
	sql_unique_u1_1 int // sql{unique="u1"}
	sql_unique_u1_2 int // sql{unique="u1"}

	// sql index
	sql_index1_1 int // sql{idx="index1"}
	sql_index1_2 int // sql{idx="index1"}
	sql_index1_3 int // sql{idx="index1"}

	// sql keys
	sql_keys_1 int // sql{keys="keys1"}
	sql_keys_2 int // sql{keys="keys1"}
	sql_keys_3 int // sql{keys="keys1"}

	// search
	sql_search string // sql{search="search" get}

	// get
	sql_get string // sql{get}

	sql_unique_x1    int // sql{unique="x1"}
	sql_unique_x2    int // sql{unique="x2"}
	sql_unique_x1_x2 int // sql{unique="x1", unique="x2"}

	sql_primary float64 //sql{primarykey}

	sql_jsonb_index map[string]any // sql{index}

	time_duration time.Duration

	go_type_int_to_strings int // go{type="[]string"}

	// add
	// sql{CREATE INDEX ON film USING GIN(mapAny)}
}
