struct NewsModel: Identifiable, Decodable, Encodable {
	var inc                     : Int? = 0 // int sql{inc}  #readonly
	var ints                    : Int? = 0 // int
	var ints8                   : Int? = 0 // int8
	var ints16                  : Int? = 0 // int16
	var ints32                  : Int? = 0 // int32
	var ints64                  : Int? = 0 // int64
	var uints                   : Int? = 0 // uint
	var uints8                  : Int? = 0 // uint8
	var uints16                 : Data? = Data() // uint16
	var uints32                 : Int? = 0 // uint32
	var uints64                 : Int? = 0 // uint64
	var floats32                : Double? = 0 // float32
	var floats64                : Double? = 0 // float64
	var bools                   : Bool? = false // bool
	var byte1                   : String? = "" // byte
	var bytes                   : Data? = Data() // []byte
	var list_ints               : [Int]? = [Int]() // []int
	var list_string             : [String]? = [String]() // []string
	var list_float              : Data? = Data() // []float64
	var map_string_string       : [String:String]? = [String:String]() // map[string]string
	var map_string_bytes        : Data? = Data() // map[string][]byte
	var map_string_bool         : [String:Bool]? = [String:Bool]() // map[string]bool
	var map_string_int          : [String:Int]? = [String:Int]() // map[string]int
	var map_string_float64      : Data? = Data() // map[string]float64
	var map_string_any          : [String:Any]? = [String:Any]() // map[string]any
	var map_int_string          : [Int:String]? = [Int:String]() // map[int]string
	var map_int_int             : [Int:Int]? = [Int:Int]() // map[int]int
	var map_int_bool            : Data? = Data() // map[int]bool
	var renameSQL               : String? = "" // string sql{name="renameSQL_OK"}
	var renameGO                : String? = "" // string go{name="RenameGO_OK"}
	var renameJS                : String? = "" // string js{name="renameJS_OK"}
	var mast_upper_go           : String? = "" // string go{up}
	var intToSmallInt           : Int? = 0 // int sql{type="smallint"} rename
	var sql_unique_u1_1         : Int? = 0 // int sql{unique="u1"}
	var sql_unique_u1_2         : Int? = 0 // int sql{unique="u1"}
	var sql_index1_1            : Int? = 0 // int sql{idx="index1"}
	var sql_index1_2            : Int? = 0 // int sql{idx="index1"}
	var sql_index1_3            : Int? = 0 // int sql{idx="index1"}
	var sql_keys_1              : Int? = 0 // int sql{keys="keys1"}
	var sql_keys_2              : Int? = 0 // int sql{keys="keys1"}
	var sql_keys_3              : Int? = 0 // int sql{keys="keys1"}
	var sql_search              : String? = "" // string sql{search="search" get}
	var sql_get                 : String? = "" // string sql{get}
	var sql_unique_x1           : Int? = 0 // int sql{unique="x1"}
	var sql_unique_x2           : Int? = 0 // int sql{unique="x2"}
	var sql_unique_x1_x2        : Int? = 0 // int sql{unique="x1", unique="x2"}
	var sql_primary             : Double? = 0 // float64 sql{primarykey}
	var sql_jsonb_index         : [String:Any]? = [String:Any]() // map[string]any sql{index}
	var time_duration           : Int? = 0 // time.Duration
	var go_type_int_to_strings  : Int? = 0 // int go{type="[]string"}
}