package main

import "fmt"

// go struct: insert sql query
func (a *SQL) Insert() []byte {
	switch settings.Fields.IncField != "" {
	case true:
		return a.insertID()
	default:
		return a.insertNoID()
	}
}

// go struct: insert sql query
func (a *SQL) InsertNoConflict() []byte {
	switch settings.Fields.IncField != "" {
	case true:
		return a.insertIDNoConflict()
	default:
		return a.insertNoIDNoConflict()
	}
}

// parser from the original struct to []interface{}
func (a *SQL) insertTuple() (list []string) {

	for _, x := range fields {
		if x.SQL.NoInsert || x.SQL.Skip {
			continue
		}
		list = append(list, fmt.Sprintf(`a.%s`, x.Go.Name))
		//lines = append(lines, fmt.Sprintf(`r[%s] = a.%s`, x.Go.Index, x.Title))
	}

	return

}
