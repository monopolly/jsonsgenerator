package main

import (
	"github.com/monopolly/jsonsgenerator/tools"
	"strings"
)

func (a *Field) parseJSOptions(v string) {

	a.Json.Name = tools.Between(v, `name="`, `"`)
	a.Json.Skip = strings.Contains(v, "skip")
	a.Json.Bool = strings.Contains(v, "bool")
	a.Json.Inc = strings.Contains(v, "inc")
	a.Json.Time = strings.Contains(v, "time")
	a.Json.Raw = strings.Contains(v, "raw")

	helps.JS["skip"] = `Skip json field`
	helps.JS["inc"] = "Add inc jsons function for numbers fields"
	helps.JS["bool"] = "Add set jsons function for bool fields"
	helps.JS["time"] = "Create convert jsons function for unixtime fields"
	helps.JS["raw"] = "Set Raw json function inside field"
	helps.JS[`name=""`] = `Replace json field name. Ex: name="sid"`
}
