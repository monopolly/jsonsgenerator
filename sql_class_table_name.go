package main

import (
	"fmt"
	"strings"
)

// get list
func (a *SQL) TableName() []byte {

	var list []string
	list = append(list, "\n//parse sql query")
	list = append(list, fmt.Sprintf(`func (a *%s) TableName() (res string) {	
		return "%s"
	}
	`, settings.SQL.Class, settings.SQL.Table))
	return []byte(strings.Join(list, "\n"))
}
