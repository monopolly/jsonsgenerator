package main

import (
	"fmt"
	"strings"
)

// sql file
func (a *Clickhouse) File() []byte {
	return []byte(a.createTableSQL() + ";\n")
}

func (a *Clickhouse) createTableSQL() (res string) {
	var list, fieldList []string
	primary := a.primaryFields()

	list = append(list, fmt.Sprintf(`create table if not exists %s (`, settings.Clickhouse.Table))

	var pad int
	for _, x := range fields {
		if len(x.Clickhouse.Name) > pad {
			pad = len(x.Clickhouse.Name)
		}
	}

	for _, x := range fields {
		if x.Clickhouse.Skip {
			continue
		}
		fieldList = append(fieldList, fmt.Sprintf("\t%s%s%s", x.Clickhouse.Name, strings.Repeat(" ", (pad+5)-len(x.Clickhouse.Name)), a.fieldType(x)))
	}

	list = append(list, strings.Join(fieldList, ",\n"))
	list = append(list, ")")
	list = append(list, a.engineSQL(primary))

	return strings.Join(list, "\n")
}

func (a *Clickhouse) fieldType(x *Field) string {
	if !x.Clickhouse.Low {
		return x.Clickhouse.Type
	}
	if strings.HasPrefix(x.Clickhouse.Type, "LowCardinality(") {
		return x.Clickhouse.Type
	}
	return fmt.Sprintf("LowCardinality(%s)", x.Clickhouse.Type)
}

func (a *Clickhouse) fieldNames() (res []string) {
	for _, x := range fields {
		if x.Clickhouse.Skip {
			continue
		}
		res = append(res, x.Clickhouse.Name)
	}
	return
}

func (a *Clickhouse) placeholders() (res []string) {
	for range a.fieldNames() {
		res = append(res, "?")
	}
	return
}

func (a *Clickhouse) primaryFields() (res []string) {
	for _, x := range fields {
		if x.Clickhouse.Skip || !x.Clickhouse.Primary {
			continue
		}
		res = append(res, x.Clickhouse.Name)
	}
	return
}

func (a *Clickhouse) engine() string {
	engine := settings.Clickhouse.Engine
	if engine == "" {
		return "MergeTree()"
	}
	if strings.Contains(engine, "(") {
		return engine
	}
	return engine + "()"
}

func (a *Clickhouse) engineSQL(primary []string) string {
	engine := a.engine()
	if strings.Contains(strings.ToUpper(engine), "ORDER BY") {
		return "engine = " + engine
	}

	if len(primary) == 0 {
		return "engine = " + engine + "\norder by tuple()"
	}

	keys := strings.Join(primary, ", ")
	return fmt.Sprintf("engine = %s\norder by (%s)\nprimary key (%s)", engine, keys, keys)
}
