package main

import (
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// goStruct, sqlClassName

/* // count
func (a *AccountSQL) Count() (count int, err error) {
	err = a.Conn(func(conn *pgxpool.Conn) (err error) {
		return conn.QueryRow(context.Background(), "select count(id) from accounts").Scan(&count)
	})
	return
}
*/
// get list
func (a *SQL) Count() []byte {

	var list []string

	s := `func (a *sqlClassName) Count(where ...string) (count int, err error) {
	
	var q string
	
	switch len(where) {
	case 0:
		q = "select count(*) from @@tableName"
	default:
		q = "select count(*) from @@tableName where "+ strings.Join(where, " ") 
	}

	err = a.Conn(func(conn *pgxpool.Conn) (err error) {
		return conn.QueryRow(context.Background(), q).Scan(&count)
	})	
	return	
}`

	s = tools.Replace(s, "@@tableName", settings.SQL.Table)
	s = tools.Replace(s, "sqlClassName", settings.SQL.Class)

	list = append(list, s)
	list = append(list, "\n\n")

	return []byte(strings.Join(list, "\n"))
}
