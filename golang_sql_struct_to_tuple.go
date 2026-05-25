package main

import (
	"fmt"
	"strings"
)

// parser from the original struct to []interface{}
func (a *Golang) SQLStructToTuple() []byte {

	lines := []string{
		"\n",
		"//Tuple create an array from struct",
		fmt.Sprintf("func (a *%s) sqlTuple() (r []any){ ", settings.Go.StructName),
	}
	var fieldlist []string
	for _, x := range fields {
		if x.SQL.NoInsert || x.SQL.Skip {
			continue
		}
		fieldlist = append(fieldlist, fmt.Sprintf(`a.%s`, x.Go.Name))
	}
	lines = append(lines, fmt.Sprintf("return []any{%s}", strings.Join(fieldlist, ",")))
	lines = append(lines, "}")

	return []byte(strings.Join(lines, "\n"))
}

// parser from the original struct to []interface{}
func (a *Golang) SQLAllStructToTuple() []byte {

	lines := []string{
		"\n",
		"//Tuple create an array from struct",
		fmt.Sprintf("func (a *%s) sqlAllTuples() (r []any){ ", settings.Go.StructName),
	}
	var fieldlist []string
	for _, x := range fields {
		if x.SQL.Skip {
			continue
		}
		fieldlist = append(fieldlist, fmt.Sprintf(`a.%s`, x.Go.Name))
	}
	lines = append(lines, fmt.Sprintf("return []any{%s}", strings.Join(fieldlist, ",")))
	lines = append(lines, "}")

	return []byte(strings.Join(lines, "\n"))
}

// parser from the original struct to []interface{}
func (a *Golang) ClickhouseStructToTuple() []byte {

	lines := []string{
		"\n",
		"//Tuple create an array from struct for clickhouse",
		fmt.Sprintf("func (a *%s) clickhouseTuple() (r []any){ ", settings.Go.StructName),
	}
	var fieldlist []string
	for _, x := range fields {
		if x.Clickhouse.Skip {
			continue
		}
		switch x.Clickhouse.IP {
		case true:
			fieldlist = append(fieldlist, fmt.Sprintf(`net.ParseIP(a.%s)`, x.Go.Name))
		default:
			fieldlist = append(fieldlist, fmt.Sprintf(`a.%s`, x.Go.Name))
		}
	}
	lines = append(lines, fmt.Sprintf("return []any{%s}", strings.Join(fieldlist, ",")))
	lines = append(lines, "}")

	return []byte(strings.Join(lines, "\n"))
}
