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

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	q := "q := `sql`"
	q = tools.Replace(q, "sql", a.createTableSQL())

	list = append(list, q)
	list = append(list, `_, err = conn.Exec(c, q)`)
	list = append(list, "return")
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

	var list, fieldList, primary []string
	unique := make(map[string]map[string]bool)

	list = append(list, fmt.Sprintf(`create table if not exists %s (`, settings.SQL.Table))

	indexgroup := map[string]map[string]bool{}
	indexsearch := map[string]map[string]bool{}

	for _, x := range fields {
		if x.SQL.Skip {
			continue
		}

		if x.SQL.Primary {
			primary = append(primary, x.SQL.Name)
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
	list = append(list, ")\n")

	return strings.Join(list, "\n")
}
