package main

import (
	"fmt"
	"jj/tools"
	"strings"
)

// get list
func (a *SQL) Delete() []byte {

	var list []string
	list = append(list, "\n//delete item")
	list = append(list, fmt.Sprintf(`func (a *%s) Delete(id any) (err error) {`, settings.SQL.Class))

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	incField := "id"
	if settings.Fields.IncField != "" {
		incField = settings.Fields.IncField
	}

	q := `_, err = conn.Exec(c, "delete from tableName where incField = $1", id)
	return
	`
	q = tools.Replace(q, "tableName", settings.SQL.Table)
	q = tools.Replace(q, "incField", incField)
	list = append(list, q)
	list = append(list, "}")

	return []byte(strings.Join(list, "\n"))
}
