package main

import (
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// goStruct, sqlClassName

// get list
func (a *SQL) Search() []byte {

	var list []string

	s := `func (a *sqlClassName) Search(q *sqlQueryName) (res []*goStruct, err error) {
	
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
	defer conn.Release()
	

	sql, fields, values := q.Render()

	rows, err := conn.Query(c, sql, values...)
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
	err = rows.Err()

	return
}`

	s = tools.Replace(s, "sqlQueryName", settings.SQL.QueryName)
	s = tools.Replace(s, "goStruct", settings.Go.StructName)
	s = tools.Replace(s, "sqlClassName", settings.SQL.Class)

	list = append(list, s)
	list = append(list, "\n\n")

	return []byte(strings.Join(list, "\n"))
}
