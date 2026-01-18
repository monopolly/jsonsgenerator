package main

import (
	"bytes"
	"fmt"
	"strings"
)

var clickhouse Clickhouse

type Clickhouse int

func (a *Clickhouse) Generate() []byte {
	var list [][]byte

	// sql class
	list = append(list, a.Class())
	// list = append(list, a.TableName())
	// list = append(list, a.Get())
	// list = append(list, a.Row())
	// list = append(list, a.All())
	// list = append(list, a.List())
	// list = append(list, a.Update())
	// list = append(list, a.Updates())
	// list = append(list, a.UpdateJsonb())
	// list = append(list, a.Query())
	// list = append(list, a.QueryInitFunction())
	// list = append(list, a.Search())
	// list = append(list, a.Delete())
	// list = append(list, a.DeleteWhere())
	// list = append(list, a.Insert())
	// list = append(list, a.Has())
	// list = append(list, a.Engine())
	// list = append(list, a.CreateTable())

	return bytes.Join(list, []byte("\n\n"))
}

// type AccountsSQL int
func (a *Clickhouse) Class() []byte {

	// type MailSQL struct{ pool *pgxpool.Pool }
	var list []string
	list = append(list, fmt.Sprintf("//clickhouse %s class", settings.Clickhouse.Class))
	list = append(list, fmt.Sprintf("type %s struct {conn driver.Conn}", settings.Clickhouse.Class))

	// Pool    *pgxpool.Pool
	sqlstructs := `
	func New%[1]s(conn driver.Conn) (a *%[1]s){
		a = new(%[1]s)
		a.conn = conn
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
	list = append(list, fmt.Sprintf(sqlstructs, settings.Clickhouse.Class))

	return []byte(strings.Join(list, "\n"))

}

// sql file
func (a *Clickhouse) File() []byte {

	return nil

}

/*
CREATE TABLE IF NOT EXISTS price (
    id UInt64,
    date DateTime,
	usd UInt64,
	created DateTime,
	source Text,
	comment Text
) ENGINE = MergeTree()
  PRIMARY KEY (id, date)
*/
