package main

import (
	"github.com/monopolly/jsonsgenerator/tools"
	"strings"
)

// sql options
func (a *Field) parseSQLOptions(v string) {

	a.SQL.NoInsert = strings.Contains(v, "noinsert")
	a.SQL.Skip = strings.Contains(v, "skip")
	a.SQL.Primary = strings.Contains(v, "primarykey")
	a.SQL.Inc = strings.Contains(v, "inc")
	if a.SQL.Inc {
		a.SQL.NoInsert = true
	}
	a.SQL.Default.Simple = strings.Contains(v, "defaults")
	a.SQL.Default.Unixtime = strings.Contains(v, "unix")
	a.SQL.Index.Simple = strings.Contains(v, "index")
	if strings.Contains(v, "altertable") {
		a.SQL.AlterTable = "1"
	}
	// a.SQL.Index.Jsonb = strings.Contains(v, "gin")
	a.SQL.Replace = tools.Between(v, `replace="`, `"`)
	a.SQL.Append = tools.Between(v, `add="`, `"`)
	a.SQL.Unique = tools.Betweens(v, `unique="`, `"`)
	a.SQL.Rename = tools.Between(v, `renames="`, `"`)
	a.SQL.Index.Group = tools.Betweens(v, `idx="`, `"`)
	a.SQL.Index.Search = tools.Betweens(v, `search="`, `"`)
	a.SQL.Type = tools.Between(v, `type="`, `"`)
	alter := tools.Between(v, `ver="`, `"`)
	if alter != "" {
		a.SQL.AlterTable = alter
	}
	name := tools.Between(v, `name="`, `"`)
	if name != "" {
		a.SQL.Name = name
	}

	// go type to postgres type
	if a.SQL.Type == "" {
		DetectSQLType(a)
	}

	DetectSQLDefault(a)

	// help
	helps.SQL["noinsert"] = "Do use field for insert function"
	helps.SQL["index"] = "Create simple default index or gin for jsonb"
	helps.SQL["altertable"] = "Add field line Alter table to SQL file with current time comment"
	helps.SQL["idx"] = `Add index fields by group name. Ex: idx="nameIndex" idx="credsIndex"`
	helps.SQL["search"] = `Add tsvector index by group. Ex: search="tsv": tsv tsvector GENERATED ALWAYS AS (to_tsvector('simple', title || ' ' || brand)). For search: SELECT brand, title FROM assets WHERE search @@ to_tsquery('english', 'f8');`
	helps.SQL["defaults"] = "Add default value based on field type"
	helps.SQL["unix"] = "Add default value: extract(epoch from now())"
	helps.SQL["skip"] = "Skip field for sql queries"
	helps.SQL[`replace="bigint primary key"`] = `Rewrite sql for field. Ex: replace="bigint primary key"`
	helps.SQL[`type="jsonb"`] = `Rewrite sql type for field. Ex: type="jsonb"`
	helps.SQL[`add="primary key"`] = `Append sql for field. Ex: replace="primary key"`
	helps.SQL[`unique="groupname"`] = `Add unique fields constrains by group (you need set group name, then generator join fields). Ex: unique="group1" unique="group2"`
	helps.SQL[`ver="v4"`] = `Create an alter table record in SQL file. Add new column. With new version in comment. Ex: ver="2"`
	helps.SQL[`renames="oldname"`] = `Create an alter table record in SQL file. Rename table column.`
}
