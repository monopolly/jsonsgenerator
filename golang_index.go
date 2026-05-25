package main

// index stores json field numbers
// must be first when it is used elsewhere
func (a *Golang) Index() []byte {

	index := a.IndexInt()
	fields := a.IndexStrings()
	function := a.IndexFunction()

	res := append(index, fields...)
	res = append(res, function...)
	return res
}
