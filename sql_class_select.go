package main

import (
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// get list
func (a *SQL) Select() []byte {

	var list []string
	list = append(list, "\n//select custom sql query")
	list = append(list, fmt.Sprintf(`func (a *%s) Select(where string, fields ...%s) (res []*%s, err error) {`, settings.SQL.Class, settings.IndexTypeName, settings.Go.StructName))

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	// if fields == nil
	ifNilFields := `if fields == nil { 
		fields = a.Fields()
	}`
	// ifNilFields = tools.Replace(ifNilFields, "allFields", strings.Join(settings.Go.Indexes, ","))
	ifNilFields = tools.Replace(ifNilFields, "indexTypeName", settings.IndexTypeName)
	list = append(list, ifNilFields)

	list = append(list, `
		var list []string
		for _,x := range fields{
			list = append(list, x.SQLName())
		}
		fieldlist := strings.Join(list, ", ")
	`)

	list = append(list, `q := fmt.Sprintf("select %s ", fieldlist)`)
	list = append(list, fmt.Sprintf(`q = q + "from %s"`, settings.SQL.Table))
	list = append(list, `if where != "" { q += " where " + where }`)
	// list = append(list, fmt.Sprintf(`res = new(%s)`, settings.Go.StructName))

	s := `rows, err := conn.Query(c, q)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item goStruct
		v, err := rows.Values()
		if err != nil {
			continue
		}
		if len(v) != len(fields){
			continue
		}

		for pos, x := range fields {
			item.Update(x.String(), v[pos])
		}
		res = append(res, &item)
	}
	err = rows.Err()`

	// s = tools.Replace(s, "sqlQueryName", settings.SQL.QueryName)
	s = tools.Replace(s, "goStruct", settings.Go.StructName)
	// s = tools.Replace(s, "sqlClassName", settings.SQL.Class)
	list = append(list, s)

	list = append(list, "return}")
	return []byte(strings.Join(list, "\n"))
}

// // get list
// func (a *SQL) Search() []byte {

// 	var list []string

// 	s := `func (a *sqlClassName) Search(q *sqlQueryName) (res []*goStruct, err error) {

// 	c := context.Background()
// 	conn, err := a.pool.Acquire(c)
// 		if err != nil {
// 			return
// 		}
// 	defer conn.Release()

// 	sql, fields, values := q.Render()

// 	rows, err := conn.Query(c, sql, values...)
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var item goStruct
// 		v, err := rows.Values()
// 		if err != nil {
// 			continue
// 		}

// 		for pos, x := range fields {
// 			item.Update(x.String(), v[pos])
// 		}
// 		res = append(res, &item)
// 	}

// 	return
// }`

// 	s = tools.Replace(s, "sqlQueryName", settings.SQL.QueryName)
// 	s = tools.Replace(s, "goStruct", settings.Go.StructName)
// 	s = tools.Replace(s, "sqlClassName", settings.SQL.Class)

// 	list = append(list, s)
// 	list = append(list, "\n\n")

// 	return []byte(strings.Join(list, "\n"))
// }
