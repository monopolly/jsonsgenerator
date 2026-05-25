package main

var JSONnamefunc = `
					//funcNameString get value as time
					func (a *structName) funcNameString() (res string) {
						return a.funcName().String()
					}

				`

var JSONany = `
			//funcName set or get value
			func (a *structName) funcName(v ...funcType) (res jsons.Result) {
				if v == nil {
					return jsons.Get((*a), fieldKey)
				}
				a.Set(fieldKey, v[0])
				return 
			}
			`

var JSONdateint = `
//funcName set or get value
func (a *structName) funcName(v ...date.Int) (res date.Int) {
	if v == nil {
		return date.Int(jsons.Int((*a), fieldKey))
	}
	a.Set(fieldKey, v[0])
	return 
}
`

var JSONconvertfunc = `
		 //funcName set or get value
		 func (a *structName) funcName(v ...funcType) (res funcType) {
			 if v == nil {
				 return funcType(jsons.jsonsFunc((*a), fieldKey))
			 }
			 a.Set(fieldKey, v[0])
			 return
		 }
		 `

var JSONgeneric = `
		 //funcName set or get value
		 func (a *structName) funcName(v ...funcType) (res funcType) {
			 if v == nil {
				 _ = jsoniter.Unmarshal([]byte(jsons.Get((*a), fieldKey).Raw), &res)
				 return
			 }
			 a.Set(fieldKey, v[0])
			 return
		 }
		 `

var JSONarray = `
	//funcName set or get value
	func (a *structName) funcName(v ...varType) (res funcType) {
		if v == nil {
			return jsons.jsonsFunc((*a), fieldKey)
		}
		a.Set(fieldKey, v)
		return 
	}

	//funcNameAdd add values
	func (a *structName) funcNameAdd(v ...varType) {
		a.Set(fieldKey, append(a.funcName(), v...))
	}

	//funcNamePrepend add values
	func (a *structName) funcNamePrepend(v ...varType) {
		a.Set(fieldKey, append(v, a.funcName()...))
	}

	//funcNameAddUnique add unique values only
	func (a *structName) funcNameAddUnique(v ...varType) {
		var list []varType
		un := map[varType]bool{}
		for _, x := range a.funcName(){
			if un[x]{
				continue
			}
			list = append(list, x)
			un[x] = true
		}
		for _, x :=range v{
			if un[x]{
				continue
			}
			list = append(list, x)
			un[x] = true
		}

		a.funcName(list...)
	}

	//funcNameDelete add unique values only
	func (a *structName) funcNameDelete(v varType) {
		var list []varType
		for _, x := range a.funcName(){
			if x == v {
				continue
			}
			list = append(list, x)
		}
		a.funcName(list...)
	}
`

var JSONarrayGeneric = `
	//funcName set or get value
	func (a *structName) funcName(v ...varType) (res funcType) {
		if v == nil {
			_ = jsoniter.Unmarshal([]byte(jsons.Get((*a), fieldKey).Raw), &res)
			return
		}
		a.Set(fieldKey, v)
		return
	}

	//funcNameAdd add values
	func (a *structName) funcNameAdd(v ...varType) {
		a.Set(fieldKey, append(a.funcName(), v...))
	}

	//funcNamePrepend add values
	func (a *structName) funcNamePrepend(v ...varType) {
		a.Set(fieldKey, append(v, a.funcName()...))
	}

	//funcNameAddUnique add unique values only
	func (a *structName) funcNameAddUnique(v ...varType) {
		var list []varType
		un := map[varType]bool{}
		for _, x := range a.funcName(){
			if un[x]{
				continue
			}
			list = append(list, x)
			un[x] = true
		}
		for _, x :=range v{
			if un[x]{
				continue
			}
			list = append(list, x)
			un[x] = true
		}

		a.funcName(list...)
	}

	//funcNameDelete add unique values only
	func (a *structName) funcNameDelete(v varType) {
		var list []varType
		for _, x := range a.funcName(){
			if x == v {
				continue
			}
			list = append(list, x)
		}
		a.funcName(list...)
	}
`

/* funcTypeMapKey = "string"
funcTypeMapValue = "string" */

var JSONmap = `
//funcName set or get value
func (a *structName) funcName(v ...funcType) (res funcType) {
	if v == nil {
		return jsons.jsonsFunc((*a), fieldKey)
	}
	a.Set(fieldKey, v[0])
	return 
}

//funcNameAdd add values
func (a *structName) funcNameAdd(k funcTypeMapKey, v funcTypeMapValue) {
	maps := a.funcName()
	maps[k] = v
	a.funcName(maps)
}

//funcNameDelete add unique values only
func (a *structName) funcNameDelete(k funcTypeMapKey) {
	maps := a.funcName()
	delete(maps, k)
	a.funcName(maps)
}

//funcNameHas check value
func (a *structName) funcNameHas(k funcTypeMapKey) bool {
	maps := a.funcName()
	return !reflect.ValueOf(maps[k]).IsZero()
}
`

var JSONinc = `
	//funcNameInc inc value
	func (a *structName) funcNameInc(v funcType) {a.Set(fieldKey, a.funcName()+v)}

	//funcNameDec dec value
	func (a *structName) funcNameDec(v funcType) {a.Set(fieldKey, a.funcName()-v)}

`

var JSONbool = `
	//funcNameTrue set true
	func (a *structName) funcNameTrue() { a.Set(fieldKey, true) }

	//funcNameFalse set false
	func (a *structName) funcNameFalse() { a.Set(fieldKey, false) }
`

var JSONdefaultfunc = `
		 //funcName set or get value
		 func (a *structName) funcName(v ...funcType) (res funcType) {
			 if v == nil {
				 return jsons.jsonsFunc((*a), fieldKey)
			 }
			 a.Set(fieldKey, v[0])
			 return 
		 }
		 `
