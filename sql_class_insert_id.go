package main

import (
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// go struct: insert sql query
func (a *SQL) insertID() []byte {

	var list []string
	list = append(list, "//Insert struct and return int id")
	list = append(list, fmt.Sprintf(`func (a *%s) Insert(v *%s) (id int, err error) {`, settings.SQL.Class, settings.Go.StructName))

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)
	//`insert into users (id, login) values ($1, $2)`

	var flist []string
	var vars []string
	var pos int
	for _, x := range fields {
		if x.SQL.NoInsert || x.SQL.Skip {
			continue
		}
		pos++
		vars = append(vars, fmt.Sprintf("$%d", pos))
		flist = append(flist, x.SQL.Name)
	}

	q := `q := "insert into tableName (fieldsList) values (valuesList) returning incField"`
	q = tools.Replace(q, "tableName", settings.SQL.Table)
	q = tools.Replace(q, "incField", settings.Fields.IncField)
	q = tools.Replace(q, "fieldsList", strings.Join(flist, ", "))
	q = tools.Replace(q, "valuesList", strings.Join(vars, ", "))

	list = append(list, q)
	list = append(list, `err = conn.QueryRow(c, q, v.sqlTuple()...).Scan(&id)`)

	list = append(list, "return")
	list = append(list, "}")

	return []byte(strings.Join(list, "\n"))
}

// go struct: insert sql query
func (a *SQL) insertIDNoConflict() []byte {

	var list []string
	list = append(list, "//Insert struct and return int id")
	list = append(list, fmt.Sprintf(`func (a *%s) InsertNoConflict(v *%s) (id int, err error) {`, settings.SQL.Class, settings.Go.StructName))

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)
	//`insert into users (id, login) values ($1, $2)`

	var flist []string
	var vars []string
	var pos int
	for _, x := range fields {
		if x.SQL.NoInsert || x.SQL.Skip {
			continue
		}
		pos++
		vars = append(vars, fmt.Sprintf("$%d", pos))
		flist = append(flist, x.SQL.Name)
	}

	q := `q := "insert into tableName (fieldsList) values (valuesList) on conflict do nothing returning incField"`
	q = tools.Replace(q, "tableName", settings.SQL.Table)
	q = tools.Replace(q, "incField", settings.Fields.IncField)
	q = tools.Replace(q, "fieldsList", strings.Join(flist, ", "))
	q = tools.Replace(q, "valuesList", strings.Join(vars, ", "))

	list = append(list, q)
	list = append(list, `err = conn.QueryRow(c, q, v.sqlTuple()...).Scan(&id)`)

	list = append(list, "return")
	list = append(list, "}")

	return []byte(strings.Join(list, "\n"))
}
