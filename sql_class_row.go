package main

import (
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// get list
func (a *SQL) Row() []byte {

	var list []string
	list = append(list, "\n//parse sql query")
	list = append(list, fmt.Sprintf(`func (a *%s) Row(eq map[string]any, fields ...%s) (res *%s, err error) {`, settings.SQL.Class, settings.IndexTypeName, settings.Go.StructName))

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

	list = append(list, `
		var andlist []string
		var values []any
		var count int
		for k,v := range eq{
			count++
			andlist = append(andlist, fmt.Sprintf("%s = $%d", k, count))
			values = append(values, v)
		}
		keys := strings.Join(andlist, " and ")
	`)

	list = append(list, `q := fmt.Sprintf("select %s from %s where %s limit 1", fieldlist, a.TableName(), keys)`)
	list = append(list, fmt.Sprintf(`res = new(%s)`, settings.Go.StructName))

	list = append(list, `rows, err := conn.Query(c, q, values...)
	if err != nil {
		return
	}
	defer rows.Close()
	rows.Next()
	v, err := rows.Values()
	if err != nil {
		return
	}
	if len(v) != len(fields){
		err = fmt.Errorf("len")
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), v[pos])
	}`)

	list = append(list, "return}")
	return []byte(strings.Join(list, "\n"))
}
