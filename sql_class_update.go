package main

import (
	"fmt"
	"strings"
)

// update field
func (a *SQL) Update() []byte {

	var list []string
	list = append(list, "\n//update sql query")
	list = append(list, fmt.Sprintf(`func (a *%s) Update(id any, k string, v any) (err error) {`, settings.SQL.Class))

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	list = append(list, fmt.Sprintf(`if !%sValidKey(k) {
		return fmt.Errorf("invalid key")
	}`, settings.Go.StructName))

	incField := "id"
	if settings.Fields.IncField != "" {
		incField = settings.Fields.IncField
	}

	inits := fmt.Sprintf("update %s set", settings.SQL.Table)
	list = append(list, `q := fmt.Sprintf("`+inits+` %s = $1 where `+incField+` = $2", k)`)
	list = append(list, fmt.Sprintf(`_,err = conn.Exec(c, q, v, id)`))
	list = append(list, "return")
	list = append(list, "}")
	return []byte(strings.Join(list, "\n"))
}
