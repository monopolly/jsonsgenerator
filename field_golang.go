package main

import (
	"jj/tools"
	"strings"
)

func (a *Field) parseGolangOptions(v string) {

	a.Go.Nofunc = strings.Contains(v, "nofunc")
	a.Go.Must = strings.Contains(v, "must")
	a.Go.UpperCase = strings.Contains(v, "up")
	a.Go.Type = tools.Between(v, `type="`, `"`)
	a.Go.Title = tools.Between(v, `title="`, `"`)
	a.Go.Desc = tools.Between(v, `desc="`, `"`)

	name := tools.Between(v, `name="`, `"`)
	if name != "" {
		a.Go.Name = name
	}

	// help
	helps.Golang["nofunc"] = "Do not create any jsons functions for fiels"
	helps.Golang["up"] = "Make uppercase for functions"
	helps.Golang["must"] = "Create one validation function for all must fields"
	helps.Golang[`type=""`] = `Replace golang type. Ex: type="[]*News"`
	helps.Golang[`name=""`] = `Replace golang struct name. Ex: name="NewsList"`
	helps.Golang[`title=""`] = `Add custom title for index`
	helps.Golang[`desc=""`] = `Add custom desc for index`

}
