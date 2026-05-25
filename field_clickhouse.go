package main

import (
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// sql options
func (a *Field) parseClickhouseOptions(v string) {

	a.Clickhouse.Type = tools.Between(v, `type="`, `"`)
	if a.Clickhouse.Type == "" {
		a.Clickhouse.Type = tools.Value(v, "type=")
	}
	a.Clickhouse.Skip = strings.Contains(v, "skip")
	a.Clickhouse.Low = strings.Contains(v, "low")
	a.Clickhouse.Primary = strings.Contains(v, "primary")
	if !a.Clickhouse.Primary {
		a.Clickhouse.Primary = strings.Contains(v, "!!")
	}
	a.Clickhouse.Unix = strings.Contains(v, "unix")

	name := tools.Between(v, `name="`, `"`)
	if name != "" {
		a.Clickhouse.Name = name
	}

	a.normalizeClickhouseIPType()

	// go type to postgres type
	if a.Clickhouse.Type == "" {
		a.Clickhouse.Type = clickhouse.Type(a)
	}

	// help
	// helps.Clickhouse[`type="uint256"`] = `Rewrite sql type for field. Ex: type="jsonb"`
	// helps.Clickhouse["skip"] = "Skip field for sql queries"
	// helps.Clickhouse[`primary or !!"`] = `Make primary key`
	// helps.Clickhouse["unix"] = "Set DateTime type"
	helps.Clickhouse["low"] = "Wrap field type in LowCardinality(...)"
	helps.Clickhouse[`type=ip`] = `Use ClickHouse IPv4 type for IP address fields`
	// helps.Clickhouse[`type=ipv4`] = `Use ClickHouse IPv4 type for IP address fields`

}

func (a *Field) normalizeClickhouseIPType() {
	t := strings.ToLower(strings.Trim(a.Clickhouse.Type, `"',`))
	switch t {
	case "ip", "ipv4", "ip4", "ipv6", "ip6":
		a.Clickhouse.IP = true
		a.Clickhouse.Type = "IPv4"
	case "":
		if strings.EqualFold(a.Name, "ip") || strings.EqualFold(a.Clickhouse.Name, "ip") {
			a.Clickhouse.IP = true
			a.Clickhouse.Type = "IPv4"
		}
	}
}

func hasClickhouseIPFields() bool {
	for _, x := range fields {
		if x.Clickhouse.Skip {
			continue
		}
		if x.Clickhouse.IP {
			return true
		}
	}
	return false
}
