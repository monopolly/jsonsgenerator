package main

import (
	"fmt"
	"strings"
)

// update field
func (a *SQL) Has() []byte {

	var list []string
	list = append(list, "\n//has value in db")
	list = append(list, fmt.Sprintf(`func (a *%s) Has(field %s, v any) (has bool, err error) {`, settings.SQL.Class, settings.IndexTypeName))

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	q := `q := fmt.Sprintf("select exists (select 1 from TableName where %s = $1 limit 1)", field.SQLName())`
	q = strings.ReplaceAll(q, "TableName", settings.SQL.Table)
	list = append(list, q)

	list = append(list, `err = conn.QueryRow(c, q, v).Scan(&has)`)
	list = append(list, "return")
	list = append(list, "}")
	return []byte(strings.Join(list, "\n"))
}
