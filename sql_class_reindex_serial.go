package main

import (
	"fmt"
	"strings"
)

/*

func (a *Asset2SQL) UpdateSerial(add ...int) (err error) {
	plus := 1
	if len(add) > 0 {
		plus = add[0]
	}
	return a.Conn(func(conn *pgxpool.Conn) (err error) {
		_, err = conn.Exec(context.Background(), "SELECT setval(pg_get_serial_sequence('assets2', 'id'), COALESCE((SELECT MAX(id) FROM assets2), 0) + $1, false)", plus)
		return
	})
}

*/

// update field
func (a *SQL) ReindexSerial() (res []byte) {

	for _, x := range fields {
		if !x.SQL.Inc {
			continue
		}
		res = append(res, a.reindexSerial(x.SQL.Name, x.Go.Name)...)
		res = append(res, []byte("\n\n")...)
		res = append(res, a.reindexSerialSet(x.SQL.Name, x.Go.Name)...)
		res = append(res, []byte("\n\n")...)
	}

	return
}

// update field
func (a *SQL) reindexSerial(sqlfield, goname string) []byte {

	var list []string
	list = append(list, "\n")
	list = append(list, "//update serial counter (count all items, add 1 and plus custom int)")
	list = append(list, "//useful if you insert with ID")
	list = append(list, fmt.Sprintf(`func (a *%s) Reindex%sSerialCounter(add ...int) (err error) {`, settings.SQL.Class, goname))

	// conn start
	list = append(list, `return a.Conn(func(conn *pgxpool.Conn) (err error) {`)

	list = append(list, `plus := 1
	if len(add) > 0 {
		plus = add[0]
	}`)

	line := `_, err = conn.Exec(context.Background(), "SELECT setval(pg_get_serial_sequence('@@tableName', '@@sqlfield'), COALESCE((SELECT MAX(@@sqlfield) FROM @@tableName), 0) + $1, false)", plus)`
	line = strings.ReplaceAll(line, "@@tableName", settings.SQL.Table)
	line = strings.ReplaceAll(line, "@@sqlfield", sqlfield)
	list = append(list, line)

	// conn end
	list = append(list, `return})`)

	// end func
	list = append(list, `}`)

	return []byte(strings.Join(list, "\n"))
}

// update field
func (a *SQL) reindexSerialSet(sqlfield, goname string) []byte {

	var list []string
	list = append(list, "\n")
	list = append(list, "//set serial counter")
	list = append(list, fmt.Sprintf(`func (a *%s) Set%sSerialCounter(value int) (err error) {`, settings.SQL.Class, goname))

	// conn start
	list = append(list, `return a.Conn(func(conn *pgxpool.Conn) (err error) {`)

	line := `_, err = conn.Exec(context.Background(), "SELECT setval(pg_get_serial_sequence('@@tableName', '@@sqlfield'), $1, false)", value)`
	line = strings.ReplaceAll(line, "@@tableName", settings.SQL.Table)
	line = strings.ReplaceAll(line, "@@sqlfield", sqlfield)
	list = append(list, line)

	// conn end
	list = append(list, `return})`)

	// end func
	list = append(list, `}`)

	return []byte(strings.Join(list, "\n"))
}
