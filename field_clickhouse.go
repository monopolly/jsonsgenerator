package main

import (
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// sql options
func (a *Field) parseClickhouseOptions(v string) {

	a.Clickhouse.Type = tools.Between(v, `type="`, `"`)
	a.Clickhouse.Skip = strings.Contains(v, "skip")
	a.Clickhouse.Primary = strings.Contains(v, "primary")
	if !a.Clickhouse.Primary {
		a.Clickhouse.Primary = strings.Contains(v, "!!")
	}
	a.Clickhouse.Unix = strings.Contains(v, "unix")

	name := tools.Between(v, `name="`, `"`)
	if name != "" {
		a.Clickhouse.Name = name
	}

	// go type to postgres type
	if a.Clickhouse.Type == "" {
		DetectClickhouseType(a)
	}

	// help
	helps.Clickhouse[`type="uint256"`] = `Rewrite sql type for field. Ex: type="jsonb"`
	helps.Clickhouse["skip"] = "Skip field for sql queries"
	helps.Clickhouse[`primary or !!"`] = `Make primary key`
	helps.Clickhouse["unix"] = "Set DateTime type"

}
