package main

import (
	"fmt"
	"strings"
)

// indexes IndexNewsID = 0
func (a *Golang) IndexInt() []byte {
	prefix := "Index"

	var list []string
	list = append(list, "\n//field type")
	list = append(list, fmt.Sprintf("type %s int", settings.IndexTypeName))

	list = append(list, "//int index")
	list = append(list, "const (")
	for pos, x := range fields {
		switch settings.Go.NoPrefix {
		case true:
			x.Go.Index = fmt.Sprintf("%s%s", prefix, x.Go.Name) //NewsID
		case false:
			x.Go.Index = fmt.Sprintf("%s%s%s", prefix, settings.Go.StructName, x.Go.Name) //IndexNewsID
		}

		settings.Go.Indexes = append(settings.Go.Indexes, x.Go.Index)

		var line string
		if pos == 0 {
			line = fmt.Sprintf(`%s = %s(iota)`, x.Go.Index, settings.IndexTypeName) //IndexID = 0
		} else {
			//line = fmt.Sprintf(`%s = %d`, x.Index, pos) //IndexID = 0
			line = x.Go.Index
		}

		list = append(list, line)
	}
	list = append(list, ")\n\n")

	return []byte(strings.Join(list, "\n"))
}
