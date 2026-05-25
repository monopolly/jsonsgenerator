package main

import (
	"fmt"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

type NewsQuery struct {
	Limit  int               `json:"limit,omitempty"`
	Offset int               `json:"offset,omitempty"`
	Sort   string            `json:"sort,omitempty"`
	Desc   bool              `json:"desc,omitempty"`
	EQ     map[string]any    `json:"eq,omitempty"`
	GT     map[string]any    `json:"gt,omitempty"`
	LT     map[string]any    `json:"lt,omitempty"`
	Like   map[string]string `json:"like,omitempty"`
	JSONB  []*JSONB
}
type JSONB struct {
	Column  string `json:"column,omitempty"` //"attr"
	Field   string `json:"field,omitempty"`  //"material"
	Command string // >, <, = , ! - has,
	Value   any
}

type Query struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`

	Sort string `json:"sort,omitempty"`
	Desc bool   `json:"desc,omitempty"`

	EQ   map[string]any    `json:"eq,omitempty"`
	GT   map[string]any    `json:"gt,omitempty"`
	LT   map[string]any    `json:"lt,omitempty"`
	Like map[string]string `json:"like,omitempty"`

	IN  map[string][]int    `json:"in,omitempty"`  //ids list
	INS map[string][]string `json:"ins,omitempty"` //ids list
}

// get list
func (a *SQL) Query() []byte {

	var list []string

	list = append(list, fmt.Sprintf("type %s struct {", settings.SQL.QueryName))
	list = append(list, fmt.Sprintf("Limit  int `%s`", `json:"limit,omitempty"`))
	list = append(list, fmt.Sprintf("Offset  int `%s`", `json:"offset,omitempty"`))
	list = append(list, fmt.Sprintf("Sort string `%s`", `json:"sort,omitempty"`))
	list = append(list, fmt.Sprintf("Desc bool `%s`", `json:"desc,omitempty"`))

	list = append(list, fmt.Sprintf("EQ   map[string]any `%s` //equal", `json:"eq,omitempty"`))
	list = append(list, fmt.Sprintf("GT   map[string]any `%s` //greater then...", `json:"gt,omitempty"`))
	list = append(list, fmt.Sprintf("LT   map[string]any `%s` //less then", `json:"lt,omitempty"`))
	list = append(list, fmt.Sprintf("NOT  map[string][]any `%s` //less then", `json:"not,omitempty"`))
	list = append(list, fmt.Sprintf("Like map[string]string `%s` //full text search", `json:"like,omitempty"`))
	list = append(list, fmt.Sprintf("Custom string`%s` //append unsafe where condition", `json:"-"`))

	list = append(list, fmt.Sprintf("Fields []string `%s`", `json:"fields,omitempty"`))
	list = append(list, fmt.Sprintf("IN map[string][]int `%s`", `json:"in,omitempty"`))
	list = append(list, fmt.Sprintf("INS map[string][]string `%s`", `json:"ins,omitempty"`))
	list = append(list, fmt.Sprintf("INSQL map[string]string `%s` //unsafe in (condition)", `json:"insql,omitempty"`))
	list = append(list, "}")
	list = append(list, "\n\n")

	// render
	s :=
		`func (a *sqlQueryName) Render() (sql string, fields []indexTypeName, values []any) {
	
			switch len(a.Fields) == 0{
			case true:
				fields = []indexTypeName{indexList}
			default:
				for _,x:=range a.Fields{						
					if !GoStructNameValidKey(x) {
						continue
					}
					fields = append(fields, GoStructNameKeyIndex(x))					
				}
			}
		
	
		var fieldsStrings []string
		for _, x := range fields {
			fieldsStrings = append(fieldsStrings, x.SQLName())
		}
	
		for k, v := range a.EQ {
			if fmt.Sprint(v) == "" {
				delete(a.EQ, k)
			}
		}
		for k, v := range a.GT {
			if fmt.Sprint(v) == "" {
				delete(a.GT, k)
			}
		}
		for k, v := range a.LT {
			if fmt.Sprint(v) == "" {
				delete(a.LT, k)
			}
		}
		for k, v := range a.NOT {
			if v == nil {
				delete(a.NOT, k)
			}
		}
		for k, v := range a.Like {
			if v == "" {
				delete(a.Like, k)
			}
		}
	
		// sql
		var list []string
	
		var count int
	
		// select
		list = append(list, "select")
		list = append(list, strings.Join(fieldsStrings, ", "))
	
		// from
		list = append(list, "from {{tablename}}")
	
		var andlist []string
	
		// EQ where
		for k, v := range a.EQ {
			if !GoStructNameValidKey(k) {
				continue
			}
			p := GoStructNameKeyIndex(k)
			sqlName := p.SQLName()
			switch p.Type() {
			case "bool":
				switch cast.Bool(v){
				case true:
					andlist = append(andlist, sqlName)
				case false:
					andlist = append(andlist, fmt.Sprintf("(not %[1]s or %[1]s is null)", sqlName))
				}
				
			default:
				count++
				andlist = append(andlist, fmt.Sprintf("%s = $%d", sqlName, count))
				values = append(values, v)
			}
		}
	
		// GT where
		for k, v := range a.GT {
	
			if !GoStructNameValidKey(k) {
				continue
			}
			sqlName := GoStructNameKeyIndex(k).SQLName()
			count++
			andlist = append(andlist, fmt.Sprintf("%s > $%d", sqlName, count))
			values = append(values, v)
		}
	
		// LT where
		for k, v := range a.LT {
			if !GoStructNameValidKey(k) {
				continue
			}
			sqlName := GoStructNameKeyIndex(k).SQLName()
			count++
			andlist = append(andlist, fmt.Sprintf("%s < $%d", sqlName, count))
			values = append(values, v)
		}

		// NOT where
		for k, v := range a.NOT {
			if !GoStructNameValidKey(k) {
				continue
			}
			sqlName := GoStructNameKeyIndex(k).SQLName()

			for _,x:= range v{
				count++
				andlist = append(andlist, fmt.Sprintf("%s != $%d", sqlName, count))
				values = append(values, x)
			}
			
		}
	
		// LIKE where
		for k, v := range a.Like {
			if !GoStructNameValidKey(k) {
				continue
			}
			sqlName := GoStructNameKeyIndex(k).SQLName()
			count++			
			andlist = append(andlist, fmt.Sprintf("%s ilike $%d", sqlName, count))			
			values = append(values, "%"+v+"%")
		}
		
		// IN where
		if a.IN != nil {
			for k, v := range a.IN {
				if !GoStructNameValidKey(k) {
					continue
				}
					sqlName := GoStructNameKeyIndex(k).SQLName()
					var inlist []string
					for _,num :=range v{
						count++
						inlist = append(inlist, fmt.Sprintf("$%d", count))
						values = append(values, num)
					}
				if len(inlist) > 0 {
					andlist = append(andlist, fmt.Sprintf("%s in (%s)", sqlName, strings.Join(inlist,",")))
				}
			}
		}

		// INS where
		if a.INS != nil {
			for k, v := range a.INS {
				if !GoStructNameValidKey(k) {
					continue
				}
				sqlName := GoStructNameKeyIndex(k).SQLName()
				var inlist []string
				for _,str :=range v{
					count++
					inlist = append(inlist, fmt.Sprintf("$%d", count))
					values = append(values, str)
				}
				if len(inlist) > 0 {
					andlist = append(andlist, fmt.Sprintf("%s in (%s)", sqlName, strings.Join(inlist,",")))
				}
			}
		}
		
		
		// INSQ where
		if a.INSQL != nil {
			for k, v := range a.INSQL {
				if !GoStructNameValidKey(k) {
					continue
				}
				sqlName := GoStructNameKeyIndex(k).SQLName()
				count++			
				andlist = append(andlist, fmt.Sprintf("%s in (%s)", sqlName, v))							
			}
		}

		if a.Custom != ""{
			andlist = append(andlist, a.Custom)
		}
	
		// render where
		if andlist != nil {
			list = append(list, "where")
			list = append(list, strings.Join(andlist, " and "))
		}
	
		// sort by
		if a.Sort != "" && GoStructNameValidKey(a.Sort) {
			list = append(list, fmt.Sprintf("order by %s", GoStructNameKeyIndex(a.Sort).SQLName()))
			if a.Desc {
				list = append(list, "desc")
			}
		}
	
		// limit, offset
		if a.Limit > 0 {
		  list = append(list, fmt.Sprintf("limit %d", a.Limit))
		}
		if a.Offset > 0 {
		  list = append(list, fmt.Sprintf("offset %d", a.Offset))
		}
		
	
		// render sql
		sql = strings.Join(list, " ")
		return
	}`

	s = tools.Replace(s, "sqlQueryName", settings.SQL.QueryName)
	s = tools.Replace(s, "GoStructName", settings.Go.StructName)
	s = tools.Replace(s, "indexTypeName", settings.IndexTypeName)
	s = tools.Replace(s, "indexList", strings.Join(a.sqlIndexes(), ", "))
	s = tools.Replace(s, "{{tablename}}", settings.SQL.Table)

	list = append(list, s)
	list = append(list, "\n\n")

	return []byte(strings.Join(list, "\n"))
}

func (a *SQL) sqlIndexes() (res []string) {
	for _, x := range fields {
		if x.SQL.Skip {
			continue
		}
		res = append(res, x.Go.Index)
	}
	return
}
