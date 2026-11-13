package main

import (
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// get list
func (a *SQL) DeleteWhere() []byte {

	var list []string
	list = append(list, "\n//delete item where: i = 1 and w = 'nice'")
	list = append(list, fmt.Sprintf(`func (a *%s) DeleteWhere(where string) (err error) {`, settings.SQL.Class))

	list = append(list, sqlConnOpen)

	q := `_, err = conn.Exec(c, fmt.Sprintf("delete from tableName where %s", where))
	return
	`
	q = tools.Replace(q, "tableName", settings.SQL.Table)
	list = append(list, q, sqlConnClose)

	return []byte(strings.Join(list, "\n"))
}
