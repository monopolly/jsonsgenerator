package main

import (
	"fmt"
	"strings"
)

// insert sql query
func DetectClickhouseType(x *Field) {

	// detect sql type
	switch x.Type {
	case "int", "time.Duration":
		x.Clickhouse.Type = "Int64"
	case "int64":
		switch x.Clickhouse.Unix {
		case true:
			x.Clickhouse.Type = "DateTime"
		default:
			x.Clickhouse.Type = "Int64"
		}
	case "uint", "uint64":
		x.Clickhouse.Type = "UInt64"
	case "uint256":
		x.Clickhouse.Type = "UInt256"
	case "float64":
		x.Clickhouse.Type = "Float64"
	case "bool":
		x.Clickhouse.Type = "Bool"
	case "string":
		x.Clickhouse.Type = "String"
	case "ip4":
		x.Clickhouse.Type = "IPv4"
	case "ip6":
		x.Clickhouse.Type = "IPv6"
	default:
		switch {
		case strings.HasPrefix(x.Type, "[]"):
			p := strings.ReplaceAll(x.Type, "[]", "")
			switch p {
			case "int":
				x.Clickhouse.Type = "Array(Int64)"
			case "string":
				x.Clickhouse.Type = "Array(String)"
			default:
				x.Clickhouse.Type = "JSON"
			}
		case strings.HasPrefix(x.Type, "map["):
			x.Clickhouse.Type = "JSON"
		default:
			panic(fmt.Sprintf("clickhouse type %s is not supported yet", x.Type))
		}

	}

}

/*
Int8	[-128 : 127]
Int16	[-32768 : 32767]
Int32	[-2147483648 : 2147483647]
Int64	[-9223372036854775808 : 9223372036854775807]
Int128	[-170141183460469231731687303715884105728 : 170141183460469231731687303715884105727]
Int256	[-57896044618658097711785492504343953926634992332820282019728792003956564819968 : 57896044618658097711785492504343953926634992332820282019728792003956564819967]

UInt8	[0 : 255]
UInt16	[0 : 65535]
UInt32	[0 : 4294967295]
UInt64	[0 : 18446744073709551615]
UInt128	[0 : 340282366920938463463374607431768211455]
UInt256	[0 : 115792089237316195423570985008687907853269984665640564039457584007913129639935]

Int8	TINYINT, INT1, BYTE, TINYINT SIGNED, INT1 SIGNED
Int16	SMALLINT, SMALLINT SIGNED
Int32	INT, INTEGER, MEDIUMINT, MEDIUMINT SIGNED, INT SIGNED, INTEGER SIGNED
Int64	BIGINT, SIGNED, BIGINT SIGNED, TIME

UInt8	TINYINT UNSIGNED, INT1 UNSIGNED
UInt16	SMALLINT UNSIGNED
UInt32	MEDIUMINT UNSIGNED, INT UNSIGNED, INTEGER UNSIGNED
UInt64	UNSIGNED, BIGINT UNSIGNED, BIT, SET

Float32 — FLOAT, REAL, SINGLE.
Float64 — DOUBLE, DOUBLE PRECISION

String — LONGTEXT, MEDIUMTEXT, TINYTEXT, TEXT, LONGBLOB, MEDIUMBLOB, TINYBLOB, BLOB, VARCHAR, CHAR, CHAR LARGE OBJECT, CHAR VARYING, CHARACTER LARGE OBJECT, CHARACTER VARYING, NCHAR LARGE OBJECT, NCHAR VARYING, NATIONAL CHARACTER LARGE OBJECT, NATIONAL CHARACTER VARYING, NATIONAL CHAR VARYING, NATIONAL CHARACTER, NATIONAL CHAR, BINARY LARGE OBJECT, BINARY VARYING

Date - unix
Time - unix from integer interpreted as number of seconds since 1970-01-01
DateTime - unix, DateTime('Asia/Istanbul')

IPv4
IPv6

Bool

Map - Map(String, UInt64) SELECT m['key2'] FROM tab INSERT INTO tab VALUES ({'key1':100}), ({});

JSON INSERT INTO test VALUES ('{"a" : {"b" : 42}, "c" : [1, 2, 3]}')

Point - is represented by its X and Y coordinates, stored as a Tuple(Float64, Float64)
Ring - is a simple polygon without holes stored as an array of points: Array(Point).
LineString - is a line stored as an array of points: Array(Point).
MultiLineString is multiple lines stored as an array of LineString: Array(LineString)
Polygon is a polygon with holes stored as an array of rings: Array(Ring). First element of outer array is the outer shape of polygon and all the following elements are holes.
MultiPolygon consists of multiple polygons and is stored as an array of polygons: Array(Polygon).

CREATE TABLE t
(
    column1 AggregateFunction(uniq, UInt64),
    column2 AggregateFunction(anyIf, String, UInt8),
    column3 AggregateFunction(quantiles(0.5, 0.9), UInt64)
)

any
any_respect_nulls
anyLast
anyLast_respect_nulls
min
max
sum
sumWithOverflow
groupBitAnd
groupBitOr
groupBitXor
groupArrayArray
groupUniqArrayArray
groupUniqArrayArrayMap
sumMap
minMap
maxMap
*/

// func (a *Clickhouse) ClickhouseType(v) (res string) {
// 	switch expression {
// 	case condition:

// 	}
// }
