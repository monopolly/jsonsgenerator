package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

type sqlIndex struct {
	Func string // CreateIndexTitle
	Name string // idx_news_title
	SQL  string // create index if not exists ...
}

// generated: (a *StructSQL) CreateIndexXxx() error for every index + CreateIndexes()
func (a *SQL) CreateIndexes() []byte {

	indexes := a.indexes()
	if len(indexes) == 0 {
		return nil
	}

	var list []string
	var calls []string
	for _, x := range indexes {
		f := `
		//create index #NAME
		func (a *#CLASS) #FUNC() (err error) {
			return a.Conn(func(conn *pgxpool.Conn) (err error) {
				_, err = conn.Exec(context.Background(), "#SQL")
				return
			})
		}`
		f = strings.ReplaceAll(f, "#NAME", x.Name)
		f = strings.ReplaceAll(f, "#CLASS", settings.SQL.Class)
		f = strings.ReplaceAll(f, "#FUNC", x.Func)
		f = strings.ReplaceAll(f, "#SQL", x.SQL)
		list = append(list, f)
		calls = append(calls, fmt.Sprintf("if err = a.%s(); err != nil {return}", x.Func))
	}

	all := `
	//create all indexes
	func (a *#CLASS) CreateIndexes() (err error) {
		#CALLS
		return
	}`
	all = strings.ReplaceAll(all, "#CLASS", settings.SQL.Class)
	all = strings.ReplaceAll(all, "#CALLS", strings.Join(calls, "\n"))
	list = append(list, all)

	return []byte(strings.Join(list, "\n"))
}

// indexes from tags: sql{index}, sql{idx="group"}, sql{search="group"}
func (a *SQL) indexes() (res []sqlIndex) {

	used := map[string]int{}
	add := func(name, idxname, sql string) {
		f := "CreateIndex" + name
		used[f]++
		if used[f] > 1 {
			f = fmt.Sprintf("%s%d", f, used[f])
		}
		res = append(res, sqlIndex{Func: f, Name: idxname, SQL: sql})
	}

	var groups, searches []string
	group := map[string][]string{}
	search := map[string][]string{}

	for _, x := range fields {
		if x.SQL.Skip {
			continue
		}

		// simple index
		if x.SQL.Index.Simple {
			switch x.SQL.Type {
			case "jsonb":
				name := fmt.Sprintf("gin_%s_%s", settings.Origin, x.SQL.Name)
				add(x.Go.Name, name, fmt.Sprintf("create index if not exists %s on %s using gin(%s)", name, settings.SQL.Table, x.SQL.Name))
			default:
				name := fmt.Sprintf("idx_%s_%s", settings.Origin, x.SQL.Name)
				add(x.Go.Name, name, fmt.Sprintf("create index if not exists %s on %s (%s)", name, settings.SQL.Table, x.SQL.Name))
			}
		}

		// fields order
		for _, g := range x.SQL.Index.Group {
			if group[g] == nil {
				groups = append(groups, g)
			}
			group[g] = append(group[g], x.SQL.Name)
		}
		for _, g := range x.SQL.Index.Search {
			if search[g] == nil {
				searches = append(searches, g)
			}
			search[g] = append(search[g], x.SQL.Name)
		}
	}

	sort.Strings(groups)
	sort.Strings(searches)

	// group index: idx="group"
	for _, g := range groups {
		name := fmt.Sprintf("idx_%s_%s", settings.Origin, g)
		add(sqlIndexFuncName(g), name, fmt.Sprintf("create index if not exists %s on %s (%s)", name, settings.SQL.Table, strings.Join(group[g], ", ")))
	}

	// tsvector index: search="group"
	for _, g := range searches {
		name := fmt.Sprintf("gin_%s_%s", settings.Origin, g)
		add(sqlIndexFuncName(g), name, fmt.Sprintf("create index if not exists %s on %s using gin(%s)", name, settings.SQL.Table, g))
	}

	return
}

// name_index -> NameIndex
func sqlIndexFuncName(v string) string {
	var res []string
	for _, x := range strings.FieldsFunc(v, func(r rune) bool { return r == '_' || r == '-' || r == ' ' }) {
		res = append(res, tools.Title(x))
	}
	return strings.Join(res, "")
}
