package main

import (
	"jj/tools"
	"strings"
)

func (a *Field) parseSwiftOptions(v string) {

	a.Swift.Files = tools.Betweens(v, `file="`, `"`)
	a.Swift.Type = tools.Between(v, `type="`, `"`)
	a.Swift.Skip = strings.Contains(v, "skip")
	a.Swift.Enum.Skip = strings.Contains(v, "skip")
	a.Swift.Must = strings.Contains(v, "must")

	// swift
	helps.Swift["must"] = "Swift field with required values"
	helps.Swift["skip"] = "Ignore field for Swift"
	helps.Swift[`type=""`] = `Replace type for Swift model. Ex: type="string"`
	helps.Swift[`file=""`] = `Create another swift file for this field Swift model. file="f1", file="f2"`

}
