package main

import (
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// get list
func (a *SQL) All() []byte {

	var list []string
	list = append(list, "//get simple list from db")
	list = append(list, fmt.Sprintf(`func (a *%s) All(fields ...%s) (res []*%s, err error) {`,
		settings.SQL.Class,
		settings.IndexTypeName,
		settings.Go.StructName,
	))

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
		for _,x:=range fields{
			list = append(list, x.SQLName())
		}
		fieldlist := strings.Join(list, ", ")
	`)
	list = append(list, `q := fmt.Sprintf("select %s", fieldlist)`)
	list = append(list, fmt.Sprintf(`q += " from %s"`, settings.SQL.Table))
	// list = append(list, `q = q + fmt.Sprintf(" ") `)
	// list = append(list, fmt.Sprintf(`a = new(%s)`, settings.Go.Name))

	list = append(list, fmt.Sprintf(`rows, err := conn.Query(c, q)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item %s
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
	
	`, settings.Go.StructName),
	)

	list = append(list, "return}")
	return []byte(strings.Join(list, "\n"))
}
