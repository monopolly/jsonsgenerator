package main

import (
	"strings"

	"github.com/monopolly/jsonsgenerator/tools"
)

func (a *Field) parseProtoOptions(v string) {
	a.Proto.Name = a.Name
	a.Proto.Type = tools.Between(v, `type="`, `"`)
	a.Proto.Skip = strings.Contains(v, "skip")

	name := tools.Between(v, `name="`, `"`)
	if name != "" {
		a.Proto.Name = name
	}

	helps.Proto[`type="string"`] = `Rewrite proto type for field. Ex: type="google.protobuf.Value"`
	helps.Proto[`name="field"`] = `Rewrite proto field name`
	helps.Proto["skip"] = "Skip field for proto"
}
