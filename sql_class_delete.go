package main

import (
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// get list
func (a *SQL) Delete() []byte {

	var list []string
	list = append(list, "\n//delete item")
	list = append(list, fmt.Sprintf(`func (a *%s) Delete(id any) (err error) {`, settings.SQL.Class))

	list = append(list, sqlConnOpen)

	incField := "id"
	if settings.Fields.IncField != "" {
		incField = settings.Fields.IncField
	}

	q := `_, err = conn.Exec(c, "delete from tableName where incField = $1", id)
	return
	`
	q = tools.Replace(q, "tableName", settings.SQL.Table)
	q = tools.Replace(q, "incField", incField)
	list = append(list, q, sqlConnClose)

	return []byte(strings.Join(list, "\n"))
}
