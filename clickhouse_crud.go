package main

import (
	"fmt"
	"strings"
)

func (a *Clickhouse) keyField() string {
	if primary := a.primaryFields(); len(primary) > 0 {
		return primary[0]
	}
	if settings.Fields.IncField != "" {
		for _, x := range fields {
			if x.SQL.Name == settings.Fields.IncField {
				return x.Clickhouse.Name
			}
		}
	}
	return "id"
}

func (a *Clickhouse) Add() []byte {
	var list []string
	list = append(list, "//Add items to clickhouse queue")
	list = append(list, fmt.Sprintf(`func (a *%s) Add(items ...%s) {`, settings.Clickhouse.Class, settings.Go.StructName))
	list = append(list, `a.Lock()
	defer a.Unlock()
	a.list = append(a.list, items...)
}`)
	return []byte(strings.Join(list, "\n"))
}

func (a *Clickhouse) Push() []byte {
	var list []string
	list = append(list, "//Push queued items to clickhouse")
	list = append(list, fmt.Sprintf(`func (a *%s) Push() (err error) {`, settings.Clickhouse.Class))
	list = append(list, `a.Lock()
	defer a.Unlock()
	if len(a.list) == 0 {
		return
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}`)
	list = append(list, fmt.Sprintf(`batch, err := a.conn.PrepareBatch(context.Background(), "insert into %s (%s)")`, settings.Clickhouse.Table, strings.Join(a.fieldNames(), ", ")))
	list = append(list, `if err != nil {
		return
	}
	for _, item := range a.list {
		err = batch.Append(item.clickhouseTuple()...)
		if err != nil {
			_ = batch.Abort()
			return
		}
	}
	err = batch.Send()
	if err != nil {
		_ = batch.Abort()
		return
	}
	a.list = a.list[:0]
	return
}`)
	return []byte(strings.Join(list, "\n"))
}

func (a *Clickhouse) Deamon() []byte {
	var list []string
	list = append(list, "//deamon pushes queued items every 10 minutes")
	list = append(list, fmt.Sprintf(`func (a *%s) deamon() {`, settings.Clickhouse.Class))
	list = append(list, `ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		a.Lock()
		hasItems := len(a.list) > 0
		a.Unlock()
		if hasItems {
			_ = a.Push()
		}
	}
}`)
	return []byte(strings.Join(list, "\n"))
}

// get count
func (a *Clickhouse) Count() []byte {
	var list []string
	list = append(list, fmt.Sprintf(`func (a *%s) Count(where ...string) (count int, err error) {`, settings.Clickhouse.Class))
	list = append(list, `
	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	var q string
	switch len(where) {
	case 0:
		q = "select count() from `+settings.Clickhouse.Table+`"
	default:
		q = "select count() from `+settings.Clickhouse.Table+` where " + strings.Join(where, " ")
	}
	var n uint64
	err = a.conn.QueryRow(context.Background(), q).Scan(&n)
	count = int(n)
	return
}`)
	return []byte(strings.Join(list, "\n"))
}

// get row by primary key
func (a *Clickhouse) Get() []byte {
	var list []string
	list = append(list, "\n//parse clickhouse query")
	list = append(list, fmt.Sprintf(`func (a *%s) Get(id any, fields ...%s) (res *%s, err error) {`, settings.Clickhouse.Class, settings.IndexTypeName, settings.Go.StructName))
	list = append(list, a.selectFieldsCode())
	list = append(list, fmt.Sprintf(`q := fmt.Sprintf("select %%s from %s where %s = ? limit 1", fieldlist)`, settings.Clickhouse.Table, a.keyField()))
	list = append(list, fmt.Sprintf(`res = new(%s)`, settings.Go.StructName))
	list = append(list, a.scanOneCode("q, id"))
	list = append(list, "return}")
	return []byte(strings.Join(list, "\n"))
}

// get row by custom where
func (a *Clickhouse) GetWhere() []byte {
	var list []string
	list = append(list, "\n//parse clickhouse query")
	list = append(list, fmt.Sprintf(`func (a *%s) GetWhere(where string, fields ...%s) (res *%s, err error) {`, settings.Clickhouse.Class, settings.IndexTypeName, settings.Go.StructName))
	list = append(list, a.selectFieldsCode())
	list = append(list, fmt.Sprintf(`q := fmt.Sprintf("select %%s from %s where %%s limit 1", fieldlist, where)`, settings.Clickhouse.Table))
	list = append(list, fmt.Sprintf(`res = new(%s)`, settings.Go.StructName))
	list = append(list, a.scanOneCode("q"))
	list = append(list, "return}")
	return []byte(strings.Join(list, "\n"))
}

// select custom query
func (a *Clickhouse) Select() []byte {
	var list []string
	list = append(list, "//select custom clickhouse query")
	list = append(list, fmt.Sprintf(`func (a *%s) Select(where string, fields ...%s) (res []*%s, err error) {`, settings.Clickhouse.Class, settings.IndexTypeName, settings.Go.StructName))
	list = append(list, a.selectFieldsCode())
	list = append(list, fmt.Sprintf(`q := fmt.Sprintf("select %%s from %s", fieldlist)`, settings.Clickhouse.Table))
	list = append(list, `if where != "" { q += " where " + where }`)
	list = append(list, a.scanListCode("q"))
	list = append(list, "return}")
	return []byte(strings.Join(list, "\n"))
}

// update field
func (a *Clickhouse) Update() []byte {
	var list []string
	list = append(list, "\n//update clickhouse query")
	list = append(list, fmt.Sprintf(`func (a *%s) Update(id any, k string, v any, where ...string) (err error) {`, settings.Clickhouse.Class))
	list = append(list, fmt.Sprintf(`if !%sValidKey(k) {
		return errors.New("invalid key")
	}`, settings.Go.StructName))
	list = append(list, `if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}`)
	list = append(list, fmt.Sprintf(`in := %sKeyIndex(k)`, settings.Go.StructName))
	list = append(list, fmt.Sprintf(`q := fmt.Sprintf("alter table %s update %%s = ? where %s = ?", in.ClickhouseName())`, settings.Clickhouse.Table, a.keyField()))
	list = append(list, `if len(where) > 0 {q += " and "+ where[0]}`)
	list = append(list, `return a.conn.Exec(context.Background(), q, v, id)
}`)
	return []byte(strings.Join(list, "\n"))
}

// update field by where
func (a *Clickhouse) UpdateWhere() []byte {
	var list []string
	list = append(list, "\n//update clickhouse query")
	list = append(list, fmt.Sprintf(`func (a *%s) UpdateWhere(k string, v any, where string) (err error) {`, settings.Clickhouse.Class))
	list = append(list, fmt.Sprintf(`if !%sValidKey(k) {return errors.New("invalid key")}`, settings.Go.StructName))
	list = append(list, `if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}`)
	list = append(list, fmt.Sprintf(`in := %sKeyIndex(k)`, settings.Go.StructName))
	list = append(list, `q := fmt.Sprintf("alter table %s update %s = ? where %s", a.TableName(), in.ClickhouseName(), where)
	return a.conn.Exec(context.Background(), q, v)
}`)
	return []byte(strings.Join(list, "\n"))
}

// update fields
func (a *Clickhouse) Updates() []byte {
	var list []string
	list = append(list, "\n//update clickhouse query")
	list = append(list, fmt.Sprintf(`func (a *%s) Updates(id any, keys map[string]any, where ...string) (err error) {`, settings.Clickhouse.Class))
	list = append(list, fmt.Sprintf(`
	if keys == nil {
		return errors.New("emptykeys")
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	var fields []string
	var values []any
	for k,v:=range keys{
		if !%[1]sValidKey(k){
			return errors.New(k)
		}
		in := %[1]sKeyIndex(k)
		fields = append(fields, fmt.Sprintf("%%s = ?", in.ClickhouseName()))
		values = append(values, v)
	}
	list := strings.Join(fields, ", ")
	values = append(values, id)
`, settings.Go.StructName))
	list = append(list, fmt.Sprintf(`q := fmt.Sprintf("alter table %s update %%s where %s = ?", list)`, settings.Clickhouse.Table, a.keyField()))
	list = append(list, `if len(where) > 0 {q += " and "+ where[0]}`)
	list = append(list, `return a.conn.Exec(context.Background(), q, values...)
}`)
	return []byte(strings.Join(list, "\n"))
}

// update fields by where
func (a *Clickhouse) UpdatesWhere() []byte {
	var list []string
	list = append(list, "\n//update clickhouse query")
	list = append(list, fmt.Sprintf(`func (a *%s) UpdatesWhere(keys map[string]any, where string, args ...any) (err error) {`, settings.Clickhouse.Class))
	list = append(list, fmt.Sprintf(`
	if keys == nil {
		return errors.New("emptykeys")
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	var fields []string
	var values []any
	for k,v:=range keys{
		if !%[1]sValidKey(k){
			return errors.New(k)
		}
		in := %[1]sKeyIndex(k)
		fields = append(fields, fmt.Sprintf("%%s = ?", in.ClickhouseName()))
		values = append(values, v)
	}
	list := strings.Join(fields, ", ")
`, settings.Go.StructName))
	list = append(list, `where = fmt.Sprintf(where, args...)`)
	list = append(list, fmt.Sprintf(`q := fmt.Sprintf("alter table %s update %%s where %%s", list, where)`, settings.Clickhouse.Table))
	list = append(list, `return a.conn.Exec(context.Background(), q, values...)
}`)
	return []byte(strings.Join(list, "\n"))
}

// delete item
func (a *Clickhouse) Delete() []byte {
	var list []string
	list = append(list, "\n//delete clickhouse item")
	list = append(list, fmt.Sprintf(`func (a *%s) Delete(id any) (err error) {`, settings.Clickhouse.Class))
	list = append(list, `if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}`)
	list = append(list, fmt.Sprintf(`q := "alter table %s delete where %s = ?"`, settings.Clickhouse.Table, a.keyField()))
	list = append(list, `return a.conn.Exec(context.Background(), q, id)
}`)
	return []byte(strings.Join(list, "\n"))
}

// delete items
func (a *Clickhouse) DeleteWhere() []byte {
	var list []string
	list = append(list, "\n//delete clickhouse items")
	list = append(list, fmt.Sprintf(`func (a *%s) DeleteWhere(where string) (err error) {`, settings.Clickhouse.Class))
	list = append(list, `if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}`)
	list = append(list, fmt.Sprintf(`q := "alter table %s delete where " + where`, settings.Clickhouse.Table))
	list = append(list, `return a.conn.Exec(context.Background(), q)
}`)
	return []byte(strings.Join(list, "\n"))
}

// has value in db
func (a *Clickhouse) Has() []byte {
	var list []string
	list = append(list, "\n//has value in clickhouse")
	list = append(list, fmt.Sprintf(`func (a *%s) Has(field %s, v any) (has bool, err error) {`, settings.Clickhouse.Class, settings.IndexTypeName))
	list = append(list, `if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}`)
	list = append(list, fmt.Sprintf(`q := fmt.Sprintf("select count() from %s where %%s = ? limit 1", field.ClickhouseName())`, settings.Clickhouse.Table))
	list = append(list, `var n uint64
	err = a.conn.QueryRow(context.Background(), q, v).Scan(&n)
	has = n > 0
	return
}`)
	return []byte(strings.Join(list, "\n"))
}

func (a *Clickhouse) Stat() []byte {
	var list []string
	list = append(list, "\n//stat by clickhouse field")
	list = append(list, fmt.Sprintf(`func (a *%s) Stat(field %s) (res map[string]int, err error) {`, settings.Clickhouse.Class, settings.IndexTypeName))
	list = append(list, `res = make(map[string]int)
	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	name := field.ClickhouseName()
	q := fmt.Sprintf("select toString(%[1]s), count() from %[2]s group by %[1]s", name, a.TableName())
	rows, err := a.conn.Query(context.Background(), q)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var k string
		var n uint64
		err = rows.Scan(&k, &n)
		if err != nil {
			return
		}
		res[k] = int(n)
	}
	err = rows.Err()
	return
}`)
	return []byte(strings.Join(list, "\n"))
}

// Create table
func (a *Clickhouse) CreateTable() []byte {
	var list []string
	list = append(list, "//Create clickhouse table")
	list = append(list, fmt.Sprintf(`func (a *%s) CreateTable() (err error) {`, settings.Clickhouse.Class))
	list = append(list, `if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}`)
	list = append(list, fmt.Sprintf("q := `%s`", a.createTableSQL()))
	list = append(list, `return a.conn.Exec(context.Background(), q)
}`)
	return []byte(strings.Join(list, "\n"))
}

func (a *Clickhouse) selectFieldsCode() string {
	return `
	if fields == nil {
		fields = a.Fields()
	}
	var list []string
	for _,x := range fields{
		list = append(list, x.ClickhouseName())
	}
	fieldlist := strings.Join(list, ", ")
`
}

func (a *Clickhouse) scanOneCode(args string) string {
	return fmt.Sprintf(`
	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	rows, err := a.conn.Query(context.Background(), %s)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("not found")
		return
	}
	rowValues := make([]any, len(fields))
	scan := make([]any, len(fields))
	for i := range rowValues {
		scan[i] = &rowValues[i]
	}
	err = rows.Scan(scan...)
	if err != nil {
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), rowValues[pos])
	}
	err = rows.Err()
`, args)
}

func (a *Clickhouse) scanListCode(args string) string {
	return fmt.Sprintf(`
	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	rows, err := a.conn.Query(context.Background(), %s)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var item %s
		values := make([]any, len(fields))
		scan := make([]any, len(fields))
		for i := range values {
			scan[i] = &values[i]
		}
		err = rows.Scan(scan...)
		if err != nil {
			continue
		}
		for pos, x := range fields {
			item.Update(x.String(), values[pos])
		}
		res = append(res, &item)
	}
	err = rows.Err()
`, args, settings.Go.StructName)
}
