package main

import (
	"jj/tools"
	"strings"
)

// tags
func (a *Field) parseTagsOptions() {
	helps.Golang["#"] = "Add lists for fields. Ex: id int //#readonly #must..."
	for _, x := range strings.Fields(a.Comment) {
		tag := tools.Value(x, "#")
		if tag != "" {
			tag = strings.TrimSpace(tag)
			a.Go.Tags[tools.Title(tag)] = true
		}
	}
}
