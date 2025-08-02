package main

import (
	"fmt"
	"strings"
)

// get list
func (a *SQL) QueryInitFunction() []byte {

	var list []string

	list = append(list, fmt.Sprintf("func New%[1]s() *%[1]s  {", settings.SQL.QueryName))

	list = append(list, fmt.Sprintf("a := new(%[1]s)", settings.SQL.QueryName))
	list = append(list, "a.EQ = make(map[string]any)")
	list = append(list, "a.GT = make(map[string]any)")
	list = append(list, "a.LT = make(map[string]any)")
	list = append(list, "a.NOT = make(map[string][]any)")
	list = append(list, "a.Like = make(map[string]string)")
	list = append(list, "a.IN = make(map[string][]int)")
	list = append(list, "a.INS = make(map[string][]string)")
	list = append(list, "a.INSQL = make(map[string]string)")
	list = append(list, "return a")

	list = append(list, "}")
	list = append(list, "\n\n")
	return []byte(strings.Join(list, "\n"))
}
