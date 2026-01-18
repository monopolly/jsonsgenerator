type news = {
	inc:                        number
	ints:                       number
	ints8:                      number
	ints16:                     number
	ints32:                     number
	ints64:                     number
	uints:                      number
	uints8:                     number
	uints16:                    any
	uints32:                    number
	uints64:                    number
	floats32:                   number
	floats64:                   number
	bools:                      boolean
	byte1:                      any
	bytes:                      string
	list_ints:                  number[]
	list_string:                string[]
	list_float:                 number[]
	map_string_string:          Record<string, string>
	map_string_bytes:           Record<string, string>
	map_string_bool:            Record<string, boolean>
	map_string_int:             Record<string, number>
	map_string_float64:         Record<string, number>
	map_string_any:             Record<string, any>
	map_int_string:             Record<string, number>
	map_int_int:                Record<string, number>
	map_int_bool:               Record<string, number>
	renameSQL:                  string
	renameGO:                   string
	renameJS:                   string
	mast_upper_go:              string
	intToSmallInt:              number
	skip:                       string
	sql_unique_u1_1:            number
	sql_unique_u1_2:            number
	sql_index1_1:               number
	sql_index1_2:               number
	sql_index1_3:               number
	sql_keys_1:                 number
	sql_keys_2:                 number
	sql_keys_3:                 number
	sql_search:                 string
	sql_get:                    string
	sql_unique_x1:              number
	sql_unique_x2:              number
	sql_unique_x1_x2:           number
	sql_primary:                number
	sql_jsonb_index:            Record<string, any>
	time_duration:              number
	go_type_int_to_strings:     number
}