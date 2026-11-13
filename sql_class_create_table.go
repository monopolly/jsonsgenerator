package main

import (
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// go struct: insert sql query
func (a *SQL) CreateTable() []byte {

	var list []string
	list = append(list, "//Create table")
	list = append(list, fmt.Sprintf(`func (a *%s) CreateTable() (err error) {`, settings.SQL.Class))

	list = append(list, `return a.Conn(func(conn *pgxpool.Conn) (err error) {`)

	q := "q := `sql`"
	q = tools.Replace(q, "sql", a.createTableSQL())
	list = append(list, q)
	list = append(list, `_, err = conn.Exec(context.Background(), q)`)
	list = append(list, "return")

	list = append(list, `})`)

	list = append(list, "}")

	return []byte(strings.Join(list, "\n"))
}

func (a *SQL) createTableSQL() (res string) {

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

	var list, fieldList []string
	unique := make(map[string]map[string]bool)
	hasPrimary := len(a.primaryKeys()) > 0

	list = append(list, fmt.Sprintf(`create table if not exists %s (`, settings.SQL.Table))

	indexgroup := map[string]map[string]bool{}
	indexsearch := map[string]map[string]bool{}

	var pad int
	for _, x := range fields {
		if len(x.SQL.Name) > pad {
			pad = len(x.SQL.Name)
		}
	}

	for _, x := range fields {
		if x.SQL.Skip {
			continue
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
			switch hasPrimary {
			case true:
				sqlFieldLine = append(sqlFieldLine, "bigserial")
			case false:
				sqlFieldLine = append(sqlFieldLine, "bigserial primary key")
			}
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

		// single primary key
		sqlFieldLine = a.primaryKeyField(x, sqlFieldLine)

		sqltype := fmt.Sprintf("\t%s%s%s", x.SQL.Name, strings.Repeat(" ", (pad+5)-len(x.SQL.Name)), strings.Join(sqlFieldLine, " "))
		fieldList = append(fieldList, sqltype)

	}

	// search index tsvector column
	for groupname, fieldslist := range indexsearch {
		var names []string
		for name := range fieldslist {
			names = append(names, name)
		}
		fieldList = append(fieldList, fmt.Sprintf("\t%s%s%s", groupname, strings.Repeat(" ", (pad+5)-len(groupname)), fmt.Sprintf("tsvector generated always as (to_tsvector('simple', %s)) stored", strings.Join(names, " || ' ' || "))))
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
	if primary := a.primaryKeySQL(); primary != "" {
		fieldList = append(fieldList, primary)
	}

	list = append(list, strings.Join(fieldList, ",\n"))
	list = append(list, ")\n")

	return strings.Join(list, "\n")
}
