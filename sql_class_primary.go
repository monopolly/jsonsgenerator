package main

import (
	"fmt"
	"strings"
)

// primary key fields: sql{primary}
func (a *SQL) primaryKeys() (res []string) {
	for _, x := range fields {
		if x.SQL.Skip || !x.SQL.Primary {
			continue
		}
		res = append(res, x.SQL.Name)
	}
	return
}

// single primary: "asset bigint primary key"
func (a *SQL) primaryKeyField(x *Field, line []string) []string {
	if !x.SQL.Primary || len(a.primaryKeys()) != 1 {
		return line
	}
	if strings.Contains(strings.Join(line, " "), "primary key") {
		return line
	}
	return append(line, "primary key")
}

// composite primary: "primary key (asset, user)"
func (a *SQL) primaryKeySQL() string {
	keys := a.primaryKeys()
	if len(keys) < 2 {
		return ""
	}
	return fmt.Sprintf("\tprimary key (%s)", strings.Join(keys, ", "))
}
