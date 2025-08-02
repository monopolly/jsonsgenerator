package main

import (
	"strings"
)

// insert sql query
func DetectSQLType(x *Field) {

	// detect sql type
	switch x.Type {
	case "int", "int64", "uint", "uint64", "time.Duration":
		x.SQL.Type = "bigint"
	case "int16", "int32", "uint32":
		x.SQL.Type = "int"
	case "int8", "uint8":
		x.SQL.Type = "smallint"
	case "float64", "float32":
		x.SQL.Type = "double precision"
	case "bool":
		x.SQL.Type = "boolean"
	case "string", "byte":
		x.SQL.Type = "text"
	case "[]byte", "[]int8", "[]uint8":
		x.SQL.Type = "bytea"
	default:
		switch {
		case strings.HasPrefix(x.Type, "[]"):
			x.SQL.Type = "jsonb"
			x.SQL.JsonbArray = true
			x.SQL.Default.Value = "default '[]'::jsonb" //by default for later updates
		case strings.HasPrefix(x.Type, "map["):
			x.SQL.Type = "jsonb"
			x.SQL.Default.Value = "default '{}'::jsonb" //by default for later updates
		default:
			x.SQL.Type = "jsonb"
			x.SQL.Default.Value = "default '{}'::jsonb" //by default for later updates
		}
	}

}

// insert sql query
func DetectSQLDefault(x *Field) {

	switch x.SQL.Default.Unixtime {
	case true:
		x.SQL.Default.Value = "default extract(epoch from now())"
	case false:
		switch x.SQL.Default.Simple {
		case true:
			switch x.SQL.Type {
			case "bigint", "int", "smallint":
				x.SQL.Default.Value = "default 0"
			case "boolean":
				x.SQL.Default.Value = "default false"
			case "text":
				x.SQL.Default.Value = "default ''"
			}
		}
	}
}
