package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

func (a *SQL) UpdateJsonb() []byte {

	var res [][]byte
	for _, x := range fields {
		switch x.SQL.Type {
		case "jsonb":
			switch x.SQL.JsonbArray {
			case true:
				res = append(res, a.updateJsonbArray(x))
				res = append(res, a.updateJsonbArrayWhere(x))
			default:
				res = append(res, a.updateJsonbMap(x))
				res = append(res, a.updatesJsonbMap(x))
				res = append(res, a.deleteJsonbMapKey(x))
				res = append(res, a.renameJsonbMapKey(x))
			}
		}
	}

	return bytes.Join(res, []byte("\n\n"))
}

// update jsob map
func (a *SQL) updateJsonbMap(x *Field) []byte {

	ff := `
		func (a *sql_class) UpdateGo_name(id any, k string, v any, where ...string) (err error) {

			var custom string
			if len(where) > 0 {
				custom = " and " where[0]
			}

			q := "update table_name set #SQLFIELDNAME = #SQLFIELDNAME || $1::jsonb where #INCFIELD"

			return a.Conn(func(conn *pool.Conn)(err error){
				_,err = conn.Exec(c, #SQL+custom, id)
			})
		}

	`

	ff = strings.ReplaceAll(ff, "#SQLCLASS", settings.SQL.Class)
	ff = strings.ReplaceAll(ff, "#GONAME", x.Go.Name)

	q := `update #TABLENAME set #SQLFIELDNAME = #SQLFIELDNAME || %s::jsonb where #INCFIELD`
	q = strings.ReplaceAll(q, "#TABLENAME", settings.SQL.Table)
	q = strings.ReplaceAll(q, "#SQLFIELDNAME", x.SQL.Name)

	ff = strings.ReplaceAll(ff, "#SQL", q)

	var list []string
	list = append(list, "\n//update sql query")
	list = append(list, fmt.Sprintf(`func (a *%s) Update%s(id any, k string, v any, where ...string) (err error) {`, settings.SQL.Class, x.Go.Name))
	incField := "id"
	if settings.Fields.IncField != "" {
		incField = settings.Fields.IncField
	}

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	sql := "update tableName set sqlFieldName = sqlFieldName"
	sql = tools.Replace(sql, "tableName", settings.SQL.Table)
	sql = tools.Replace(sql, "sqlFieldName", x.SQL.Name)

	list = append(list, `//json escape`)
	list = append(list, `res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")`)
	list = append(list, `q := fmt.Sprintf("`+sql+` || $$%s$$::jsonb where `+incField+` = $1",  res)`)
	list = append(list, `if len(where) > 0 {q += " and "+ where[0]}`)
	list = append(list, `_,err = conn.Exec(c, q, id)`)
	list = append(list, "return")
	list = append(list, "}")
	return []byte(strings.Join(list, "\n"))
}

// update jsob map
func (a *SQL) deleteJsonbMapKey(x *Field) []byte {
	// update test set slots = slots - '1'

	var list []string
	list = append(list, "\n//delete key from jsonb")
	list = append(list, fmt.Sprintf(`func (a *%s) DeleteKey%s(id any, k string, where ...string) (err error) {`, settings.SQL.Class, x.Go.Name))
	incField := "id"
	if settings.Fields.IncField != "" {
		incField = settings.Fields.IncField
	}

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	sql := `"update tableName set sqlFieldName = sqlFieldName - $1 where ` + incField + ` = $2"`
	sql = tools.Replace(sql, "tableName", settings.SQL.Table)
	sql = tools.Replace(sql, "sqlFieldName", x.SQL.Name)

	list = append(list, fmt.Sprintf("q := %s", sql))
	list = append(list, `if len(where) > 0 {q += " and "+ where[0]}`)
	list = append(list, `_,err = conn.Exec(c, q, k, id)`)
	list = append(list, "return")
	list = append(list, "}")
	return []byte(strings.Join(list, "\n"))
}

// update jsob map
func (a *SQL) renameJsonbMapKey(x *Field) []byte {
	// update assets set meta = meta - 'ip' || jsonb_build_object('newip', meta->'ip') where meta ->> 'feat' = 'true';

	var list []string
	list = append(list, "\n//rename map key jsonb")
	list = append(list, fmt.Sprintf(`func (a *%s) RenameKey%s(id any, k, newkey string) (err error) {`, settings.SQL.Class, x.Go.Name))
	incField := "id"
	if settings.Fields.IncField != "" {
		incField = settings.Fields.IncField
	}

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	sql := `"update tableName set sqlFieldName = sqlFieldName - $1 || jsonb_build_object($2, sqlFieldName->$1) where ` + incField + ` = $3"`
	sql = tools.Replace(sql, "tableName", settings.SQL.Table)
	sql = tools.Replace(sql, "sqlFieldName", x.SQL.Name)

	list = append(list, fmt.Sprintf("q := %s", sql))

	list = append(list, "_,err = conn.Exec(c, q, k, newkey, id)")
	list = append(list, "return")
	list = append(list, "}")
	return []byte(strings.Join(list, "\n"))
}

// update jsob map
func (a *SQL) updatesJsonbMap(x *Field) []byte {

	var list []string
	list = append(list, "\n//update sql query")

	incField := "id"
	if settings.Fields.IncField != "" {
		incField = settings.Fields.IncField
	}

	list = append(list, fmt.Sprintf(`func (a *%s) Updates%s(id any, keys map[string]any, where ...string) (err error) {`, settings.SQL.Class, x.Go.Name))

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	sql := "update tableName set sqlFieldName = sqlFieldName"
	sql = tools.Replace(sql, "tableName", settings.SQL.Table)
	sql = tools.Replace(sql, "sqlFieldName", x.SQL.Name)

	list = append(list, `b, _ := jsoniter.Marshal(keys) `)

	list = append(list, `//json escape`)
	list = append(list, `res := strings.ReplaceAll(string(b), "$$", "$ $")`)

	list = append(list, `q := fmt.Sprintf("`+sql+` || $$%s$$::jsonb where `+incField+` = $1",  res)`)
	list = append(list, `if len(where) > 0 {q += " and "+ where[0]}`)
	list = append(list, `_,err = conn.Exec(c, q, id)`)
	list = append(list, "return")
	list = append(list, "}")
	return []byte(strings.Join(list, "\n"))
}

// add/delete jsonb array
func (a *SQL) updateJsonbArray(x *Field) []byte {

	/*
		add jsonb value
		UPDATE jsontesting
		SET jsondata = jsondata || '["newString"]'::jsonb
		WHERE id = 7;
	*/

	incField := "id"
	if settings.Fields.IncField != "" {
		incField = settings.Fields.IncField
	}

	var list []string
	list = append(list, "\n//add array value to jsonb array")
	list = append(list, fmt.Sprintf(`func (a *%s) Add%s(id any, v ...%s) (err error) {`, settings.SQL.Class, x.Go.Name, x.Go.ArrayType))

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	sql := "update tableName set sqlFieldName = sqlFieldName"
	sql = tools.Replace(sql, "tableName", settings.SQL.Table)
	sql = tools.Replace(sql, "sqlFieldName", x.SQL.Name)

	list = append(list, `//json escape`)
	list = append(list, "vvv, _ := jsoniter.Marshal(v)")
	list = append(list, "vv := string(vvv)")
	list = append(list, `q := "`+sql+` || $1::jsonb where `+incField+` = $2"`)

	list = append(list, "_,err = conn.Exec(c, q, vv, id)")
	list = append(list, "return")
	list = append(list, "}")

	/*
			remove jsonb value
			UPDATE jsontesting
		SET jsondata = jsondata - 'newString'
		WHERE id = 7;  */

	list = append(list, "\n//delete array value from jsonb array")
	list = append(list, fmt.Sprintf(`func (a *%s) Delete%s(id any, v %s, where ...string) (err error) {`, settings.SQL.Class, x.Go.Name, x.Go.ArrayType))

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	sql = "update tableName set sqlFieldName = sqlFieldName"
	sql = tools.Replace(sql, "tableName", settings.SQL.Table)
	sql = tools.Replace(sql, "sqlFieldName", x.SQL.Name)

	list = append(list, `//json escape`)
	list = append(list, fmt.Sprintf("vvv, _ := jsoniter.Marshal(%s{v})", x.Type))
	list = append(list, "vv := string(vvv)")
	list = append(list, `q := "`+sql+` - $1::jsonb where `+incField+` = $2"`)

	// list = append(list, `q := "`+sql+` - $1 where `+incField+` = $2"`)
	list = append(list, `if len(where) > 0 {q += " and "+ where[0]}`)
	list = append(list, "_,err = conn.Exec(c, q, vv, id)")
	list = append(list, "return")
	list = append(list, "}")

	return []byte(strings.Join(list, "\n"))
}

// add/delete jsonb array
func (a *SQL) updateJsonbArrayWhere(x *Field) []byte {

	/*
		add jsonb value
		UPDATE jsontesting
		SET jsondata = jsondata || '["newString"]'::jsonb
		WHERE id = 7;
	*/

	var list []string
	list = append(list, "\n//add array value to jsonb array")
	list = append(list, fmt.Sprintf(`func (a *%s) Add%sWhere(v %s, where string) (err error) {`, settings.SQL.Class, x.Go.Name, x.Go.ArrayType))

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	sql := "update tableName set sqlFieldName = sqlFieldName"
	sql = tools.Replace(sql, "tableName", settings.SQL.Table)
	sql = tools.Replace(sql, "sqlFieldName", x.SQL.Name)

	// list = append(list, `//json escape`)
	// list = append(list, `res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")`)
	// list = append(list, `q := fmt.Sprintf("`+sql+` || '%s'::jsonb where %s",  res, where)`)

	list = append(list, `//json escape`)
	list = append(list, fmt.Sprintf("vvv, _ := jsoniter.Marshal(%s{v})", x.Type))
	list = append(list, "vv := string(vvv)")
	list = append(list, `q := "`+sql+` || $1::jsonb where " + where`)

	list = append(list, "_,err = conn.Exec(c, vv, q)")
	list = append(list, "return")
	list = append(list, "}")

	/*
			remove jsonb value
			UPDATE jsontesting
		SET jsondata = jsondata - 'newString'
		WHERE id = 7;  */

	list = append(list, "\n//delete array value from jsonb array")
	list = append(list, fmt.Sprintf(`func (a *%s) Delete%sWhere(v %s, where string) (err error) {`, settings.SQL.Class, x.Go.Name, x.Go.ArrayType))

	list = append(list, `
	c := context.Background()
	conn, err := a.pool.Acquire(c)
		if err != nil {
			return
		}
		defer conn.Release()
	`)

	sql = "update tableName set sqlFieldName = sqlFieldName"
	sql = tools.Replace(sql, "tableName", settings.SQL.Table)
	sql = tools.Replace(sql, "sqlFieldName", x.SQL.Name)

	list = append(list, `//json escape`)
	list = append(list, fmt.Sprintf("vvv, _ := jsoniter.Marshal(%s{v})", x.Type))
	list = append(list, "vv := string(vvv)")
	list = append(list, `q := "`+sql+` - $1::jsonb where " + where`)

	// list = append(list, `q := fmt.Sprintf("`+sql+` - $1 where %s",where)`)
	list = append(list, "_,err = conn.Exec(c, q, vv)")
	list = append(list, "return")
	list = append(list, "}")

	return []byte(strings.Join(list, "\n"))
}
