package main

import (
	"bytes"
)

var golang Golang

type Golang int

func (a *Golang) Generate() []byte {
	var res bytes.Buffer

	res.Write(a.Index())
	res.Write(a.Struct())
	res.Write(a.StructCustom())
	res.Write(a.New())
	res.Write(a.TupleToStruct())
	res.Write(a.IDStructToTuple())
	res.Write(a.SQLStructToTuple())
	res.Write(a.SQLAllStructToTuple())
	res.Write(a.ClickhouseStructToTuple())
	res.Write(a.Update())
	res.Write(a.Get())
	res.Write(a.GetString())
	res.Write(a.ToJson())
	res.Write(a.Valid())
	res.Write(a.Tags())
	res.Write(a.IndexType())
	res.Write(a.IndexKeys())
	res.Write(a.Map())
	res.Write(a.Iterate())
	res.Write(a.Gotiny())
	res.Write(a.GotinyUnmarshall())
	res.Write(a.Msgp())
	res.Write(a.MsgpUnmarshall())
	res.Write(a.FastJson())
	res.Write(a.FastJsonUnmarshall())
	return res.Bytes()
}
