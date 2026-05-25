package main

import (
	"bytes"
)

var clickhouse Clickhouse

type Clickhouse int

func (a *Clickhouse) Generate() []byte {
	var list [][]byte

	// sql class
	list = append(list, a.Class())
	list = append(list, a.TableName())
	list = append(list, a.Count())
	list = append(list, a.Add())
	list = append(list, a.Push())
	list = append(list, a.Deamon())
	list = append(list, a.Get())
	list = append(list, a.GetWhere())
	list = append(list, a.Update())
	list = append(list, a.UpdateWhere())
	list = append(list, a.Updates())
	list = append(list, a.UpdatesWhere())
	list = append(list, a.Select())
	list = append(list, a.Delete())
	list = append(list, a.DeleteWhere())
	list = append(list, a.Has())
	list = append(list, a.Stat())
	list = append(list, a.CreateTable())
	list = append(list, a.Fields())

	return bytes.Join(list, []byte("\n\n"))
}

/*
CREATE TABLE IF NOT EXISTS price (
    id UInt64,
    date DateTime,
	usd UInt64,
	created DateTime,
	source Text,
	comment Text
) ENGINE = MergeTree()
  PRIMARY KEY (id, date)
*/
