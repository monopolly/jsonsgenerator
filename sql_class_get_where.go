package main

import (
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// get list
func (a *SQL) GetWhere() []byte {

	var list []string
	list = append(list, "\n//parse sql query")
	list = append(list, fmt.Sprintf(`func (a *%s) GetWhere(where string, fields ...%s) (res *%s, err error) {`, settings.SQL.Class, settings.IndexTypeName, settings.Go.StructName))

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

	list = append(list, `q := fmt.Sprintf("select %s", fieldlist)`)
	list = append(list, fmt.Sprintf(`q = q + " from %s"`, settings.SQL.Table))
	list = append(list, `q = q + fmt.Sprintf(" where %s ", where)`)
	list = append(list, `q = q + " limit 1"`)
	list = append(list, fmt.Sprintf(`res = new(%s)`, settings.Go.StructName))

	list = append(list, `rows, err := conn.Query(c, q)
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
		err = errors.New("len")
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), v[pos])
	}`)

	list = append(list, "return}")
	return []byte(strings.Join(list, "\n"))
}
