package main

import (
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// update field
func (a *SQL) Updates() []byte {

	var list []string
	list = append(list, "\n//update sql query")
	list = append(list, fmt.Sprintf(`func (a *%s) Updates(id any, keys map[string]any) (err error) {`, settings.SQL.Class))

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	m := `
		if keys == nil {
			return fmt.Errorf("emptykeys")
		}


		var fields []string
		var values []any
		var count int
		for k,v:=range keys{
			if !GoStructNameValidKey(k){
				return fmt.Errorf(k)
			}
			in := GoStructNameKeyIndex(k)			
			count++			
			fields = append(fields, fmt.Sprintf("%s = $%d", in.SQLName(), count))
			values = append(values, v)
		}

		list := strings.Join(fields, ", ")
		count++
		values = append(values, id)
	`
	m = strings.ReplaceAll(m, "GoStructName", settings.Go.StructName)
	list = append(list, m)

	incField := "id"
	if settings.Fields.IncField != "" {
		incField = settings.Fields.IncField
	}

	list = append(list, tools.Replace(`q := fmt.Sprintf("update tableName set %s where `+incField+` = $%d", list, count)`, "tableName", settings.SQL.Table))
	list = append(list, fmt.Sprintf(`_,err = conn.Exec(c, q, values...)`))
	list = append(list, "return")
	list = append(list, "}")
	return []byte(strings.Join(list, "\n"))
}
