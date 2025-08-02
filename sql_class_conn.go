package main

import (
	"fmt"
	"strings"
)

// get list
func (a *SQL) Conn() []byte {

	var list []string
	list = append(list, "\n//delete item")
	list = append(list, fmt.Sprintf(`
	func (a *%s) Conn(f func(conn *pgxpool.Conn) (err error)) (err error) {
			
			conn, err := a.pool.Acquire(context.Background())
			if err != nil {		
				return
			}
			defer conn.Release()
			return f(conn)
		}
	`, settings.SQL.Class))

	return []byte(strings.Join(list, "\n"))
}
