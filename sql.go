package main

import (
	"bytes"
	"fmt"
	"strings"
)

var sql SQL

type SQL int

func (a *SQL) Generate() []byte {
	var list [][]byte

	// sql class
	list = append(list, a.Class())
	list = append(list, a.Conn())
	list = append(list, a.TableName())
	list = append(list, a.Count())
	list = append(list, a.Insert())
	list = append(list, a.InsertNoConflict())
	list = append(list, a.InsertFull())
	list = append(list, a.ReindexSerial())
	list = append(list, a.Get())
	list = append(list, a.GetCustomStructs())
	list = append(list, a.GetWhere())
	list = append(list, a.Row())
	list = append(list, a.All())
	list = append(list, a.List())
	list = append(list, a.Update())
	list = append(list, a.UpdateWhere())
	list = append(list, a.Updates())
	list = append(list, a.UpdatesWhere())
	list = append(list, a.UpdateJsonb())
	list = append(list, a.Query())
	list = append(list, a.QueryInitFunction())
	list = append(list, a.Search())
	list = append(list, a.Select())
	list = append(list, a.Delete())
	list = append(list, a.DeleteWhere())

	list = append(list, a.Has())
	// list = append(list, a.Engine())
	list = append(list, a.CreateTable())
	list = append(list, a.Fields())

	return bytes.Join(list, []byte("\n\n"))
}

// type AccountsSQL int
func (a *SQL) Class() []byte {

	var list []string
	list = append(list, fmt.Sprintf("//sql %s class", settings.SQL.Class))
	// list = append(list, fmt.Sprintf("var %s %s", settings.SQL.ClassVarName, settings.SQL.Class))
	list = append(list, fmt.Sprintf("type %s struct {pool *pgxpool.Pool}", settings.SQL.Class))
	// Pool    *pgxpool.Pool
	sqlstructs := `
	func New%[1]s(pool *pgxpool.Pool) (a *%[1]s){
		a = new(%[1]s)
		a.pool = pool
		a.CreateTable()
		return
	}
	`

	/*
		conn, err := a.pool.Acquire(context.Background())
		if err != nil {
			return
		}
		defer conn.Release()
	*/
	list = append(list, fmt.Sprintf(sqlstructs, settings.SQL.Class))

	return []byte(strings.Join(list, "\n"))

}
