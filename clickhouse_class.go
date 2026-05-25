package main

import (
	"fmt"
	"strings"
)

// type AccountsSQL int
func (a *Clickhouse) Class() []byte {

	// type MailSQL struct{ pool *pgxpool.Pool }
	var list []string
	list = append(list, fmt.Sprintf("//clickhouse %s class", settings.Clickhouse.Class))
	list = append(list, fmt.Sprintf(`type %s struct {
	conn driver.Conn
	list []%s
	sync.Mutex
}`, settings.Clickhouse.Class, settings.Go.StructName))

	// Pool    *pgxpool.Pool
	sqlstructs := `
	func New%[1]s(conn driver.Conn) (a *%[1]s){
		a = new(%[1]s)
		a.conn = conn
		a.CreateTable()
		go a.deamon()
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

func (a *Clickhouse) TableName() []byte {

	var list []string
	list = append(list, "\n//parse clickhouse query")
	list = append(list, fmt.Sprintf(`func (a *%s) TableName() (res string) {
		return "%s"
	}
	`, settings.Clickhouse.Class, settings.Clickhouse.Table))
	return []byte(strings.Join(list, "\n"))
}

func (a *Clickhouse) Fields() []byte {

	var list []string
	list = append(list, "\n//parse clickhouse query")
	list = append(list, fmt.Sprintf(`func (a *%s) Fields() (res []%s) {`, settings.Clickhouse.Class, settings.IndexTypeName))
	list = append(list, fmt.Sprintf(`return []%s { %s }`, settings.IndexTypeName, strings.Join(settings.Go.Indexes, ",")))
	list = append(list, "}")
	return []byte(strings.Join(list, "\n"))
}
