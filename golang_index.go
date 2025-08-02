package main

// индекс это номера полей для json
// должен быть первым если где то используется
func (a *Golang) Index() []byte {

	index := a.IndexInt()
	fields := a.IndexStrings()
	function := a.IndexFunction()

	res := append(index, fields...)
	res = append(res, function...)
	return res
}
