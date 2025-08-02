package main

import (
	"fmt"
	"jj/tools"
	"strings"
)

// func (a *SQL) StructToPackage() []byte {

// 	var list []string
// 	list = append(list, "\n//get table sql struct")
// 	list = append(list, fmt.Sprintf(`func (a *%s) SQLTableStructure() (sql string) {`, settings.Go.StructName))
// 	list = append(list, fmt.Sprintf("return `%s`", string(a.Model())))
// 	list = append(list, "\n}")
// 	return []byte(strings.Join(list, "\n"))
// }

// insert sql query
func (a *SQL) File() []byte {

	/* 	id      	bigint primary key,
		login   	text unique,
		created 	bigint,
		updated 	bigint,
		name    	text,
		line    	text,
		body  		text, --description
		image   	text,
		lang    	text,
		country    	text,
		city    	text,
		sex 		int, -- 0 female, 1 male, null - nothing
		geoid    	bigint,
		tz			text,
		verify  	boolean default false,
	  	meta    	jsonb default '{}'::jsonb */

	var list, fieldList, primary []string
	unique := make(map[string]map[string]bool)

	//list = append(list, fmt.Sprintf(`drop table if exists %s;`, table))
	list = append(list, fmt.Sprintf(`--drop table %s;`, settings.SQL.Table))
	list = append(list, fmt.Sprintf(`create table if not exists %s (`, settings.SQL.Table))

	var indexes []string
	altertable := map[string][]string{}
	renames := map[string]string{}
	indexgroup := map[string]map[string]bool{}
	indexsearch := map[string]map[string]bool{}

	for _, x := range fields {
		if x.SQL.Skip {
			continue
		}

		if x.SQL.Primary {
			primary = append(primary, x.SQL.Name)
		}

		switch x.SQL.Type {
		case "jsonb":
			// jsonb indexes
			// create index on film using gin(special_features);
			if x.SQL.Index.Simple {
				indexname := fmt.Sprintf("gin_%s_%s", settings.Origin, x.SQL.Name)
				indexes = append(indexes, fmt.Sprintf("--drop index %s;", indexname))
				indexes = append(indexes, fmt.Sprintf("--ex: select * from %s where (%s->>'year')::int >= 2020;", settings.SQL.Table, x.SQL.Name))
				indexes = append(indexes, fmt.Sprintf("create index if not exists %s on %s using gin(%s);\n", indexname, settings.SQL.Table, x.SQL.Name))
			}
		default:
			// simple indexes
			if x.SQL.Index.Simple {
				indexname := fmt.Sprintf("idx_%s_%s", settings.Origin, x.SQL.Name)
				indexes = append(indexes, fmt.Sprintf("--drop index %s;", indexname))
				indexes = append(indexes, fmt.Sprintf("create index if not exists %s on %s (%s);\n", indexname, settings.SQL.Table, x.SQL.Name))
			}
		}

		if x.SQL.AlterTable != "" {
			altertable[x.SQL.AlterTable] = append(altertable[x.SQL.AlterTable], fmt.Sprintf("alter table %s add column %s %s;", settings.SQL.Table, x.SQL.Name, x.SQL.Type))
		}

		if x.SQL.Rename != "" {
			// fmt.Println("rename line", x.SQL.Rename, x.SQL.Name)
			renames[x.SQL.Rename] = x.SQL.Name
		}

		for _, groupname := range x.SQL.Index.Group {
			if indexgroup[groupname] == nil {
				indexgroup[groupname] = make(map[string]bool)
			}
			indexgroup[groupname][x.SQL.Name] = true
		}

		for _, groupname := range x.SQL.Index.Search {
			if indexsearch[groupname] == nil {
				indexsearch[groupname] = make(map[string]bool)
			}
			indexsearch[groupname][x.SQL.Name] = true
		}

		// unique constrains tag
		for _, k := range x.SQL.Unique {
			if unique[k] == nil {
				unique[k] = map[string]bool{}
			}
			unique[k][x.SQL.Name] = true
		}

		var sqlFieldLine []string

		switch x.SQL.Inc {
		case true:
			sqlFieldLine = append(sqlFieldLine, "bigserial primary key")
		case false:
			switch x.SQL.Replace != "" {
			case true:
				sqlFieldLine = append(sqlFieldLine, x.SQL.Replace)
			case false:
				sqlFieldLine = append(sqlFieldLine, x.SQL.Type)
			}
		}

		// append sql
		if x.SQL.Append != "" {
			sqlFieldLine = append(sqlFieldLine, x.SQL.Append)
		}

		// default value
		if x.SQL.Default.Value != "" {
			sqlFieldLine = append(sqlFieldLine, x.SQL.Default.Value)
		}

		sqltype := tools.Formatline("\t"+x.SQL.Name, strings.Join(sqlFieldLine, " "))
		fieldList = append(fieldList, sqltype)

	}

	// search index tsvector
	for groupname, fieldslist := range indexsearch {
		var names []string
		for name := range fieldslist {
			names = append(names, name)
		}
		// fieldList = append(fieldList, fmt.Sprintf("\t%s \ttsvector generated always as (to_tsvector('simple', %s)) stored", groupname, strings.Join(names, " || ' ' || ")))
		fieldList = append(fieldList, tools.Formatline("\t"+groupname, fmt.Sprintf("tsvector generated always as (to_tsvector('simple', %s)) stored", strings.Join(names, " || ' ' || "))))
		indexname := fmt.Sprintf("gin_%s_%s", settings.Origin, groupname)
		indexes = append(indexes, fmt.Sprintf("--drop index %s;", indexname))
		indexes = append(indexes, fmt.Sprintf("--ex: select * from %s where %s @@ to_tsquery('f8');", settings.SQL.Table, groupname))
		// CREATE INDEX ts_idx ON certificate USING GIN (to_tsvector('english', text));
		indexes = append(indexes, fmt.Sprintf("create index if not exists %s on %s using gin(%s);\n", indexname, settings.SQL.Table, groupname))
	}

	// unique constrains
	for _, fieldnames := range unique {
		var untag []string
		for fieldname := range fieldnames {
			untag = append(untag, fieldname)
		}
		fieldList = append(fieldList, fmt.Sprintf("\tunique(%s)", strings.Join(untag, ",")))
	}

	// primary keys
	if primary != nil {
		fieldList = append(fieldList, fmt.Sprintf("\tprimary key (%s)", strings.Join(primary, ",")))
	}

	list = append(list, strings.Join(fieldList, ",\n"))
	list = append(list, ");\n")

	if indexes != nil {
		list = append(list, "--generated indexes")
		list = append(list, indexes...)
	}

	if len(altertable) > 0 {
		list = append(list, "--altertable")
		for k, addlist := range altertable {
			list = append(list, fmt.Sprintf("--%s", k))
			list = append(list, addlist...)
		}
	}

	if len(renames) > 0 {
		// fmt.Println("renames", renames)
		list = append(list, "--rename")
		for from, to := range renames {
			list = append(list, fmt.Sprintf(`alter table %s rename column %s to %s;`, settings.SQL.Table, from, to))
		}
	}

	// simple index
	for groupname, fieldlist := range indexgroup {
		var names []string
		for name := range fieldlist {
			names = append(names, name)
		}
		indexname := fmt.Sprintf("idx_%s_%s", settings.Origin, groupname)
		list = append(list, fmt.Sprintf("--drop index %s;", indexname))
		list = append(list, fmt.Sprintf("create index if not exists %s on %s (%s);\n", indexname, settings.SQL.Table, strings.Join(names, ",")))
	}

	return []byte(strings.Join(list, "\n"))
}
