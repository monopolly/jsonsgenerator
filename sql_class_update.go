package main

import (
	"fmt"
	"strings"
)

// update field
func (a *SQL) Update() []byte {

	var list []string
	list = append(list, "\n//update sql query")
	list = append(list, fmt.Sprintf(`func (a *%s) Update(id any, k string, v any, where ...string) (err error) {`, settings.SQL.Class))

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	list = append(list, fmt.Sprintf(`if !%sValidKey(k) {
		return errors.New("invalid key")
	}`, settings.Go.StructName))

	incField := "id"
	if settings.Fields.IncField != "" {
		incField = settings.Fields.IncField
	}

	inits := fmt.Sprintf("update %s set", settings.SQL.Table)
	list = append(list, `q := fmt.Sprintf("`+inits+` %s = $1 where `+incField+` = $2", k)`)
	list = append(list, `if len(where) > 0 {q += " and "+ where[0]}`)
	list = append(list, `_,err = conn.Exec(c, q, v, id)`)
	list = append(list, "return")
	list = append(list, "}")
	return []byte(strings.Join(list, "\n"))
}

// update field
func (a *SQL) UpdateWhere() []byte {

	var list []string
	list = append(list, "\n//update sql query")
	list = append(list, fmt.Sprintf(`func (a *%s) UpdateWhere(k string, v any, where string) (err error) {`, settings.SQL.Class))

	list = append(list, fmt.Sprintf(`if !%sValidKey(k) {return errors.New("invalid key")}`, settings.Go.StructName))

	list = append(list, `return a.Conn(func(conn *pgxpool.Conn) (err error) {
		q := fmt.Sprintf("update %s set %s = $1 where %s", a.TableName(), k, where)
		_,err = conn.Exec(context.Background(), q, v)
		return
	})
		}	
	`)
	// list = append(list, "}")
	return []byte(strings.Join(list, "\n"))
}
