package news

// go ch sql js proto test gotiny
type stat2 struct {
	id     int // sql{inc}  #readonly
	uid    int // ch{primary}
	types  int // ch{low}
	name   string
	source string
	ip     string
	image  bool
}
