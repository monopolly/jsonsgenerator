package main

import (
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

// tags
func (a *Field) parseCustomStructsOptions() {
	helps.Golang["@"] = "Add custom structs for fields. Ex: id int //@public @team..."
	for x := range strings.FieldsSeq(a.Comment) {
		tag := tools.Value(x, "@")
		if tag != "" {
			tag = strings.TrimSpace(tag)
			a.Go.CustomStructs[tools.Title(tag)] = true
			a.Go.Tags[tools.Title(tag)] = true
		}
	}
}
