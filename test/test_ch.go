package news

// go ch sql js proto test gotiny
type stat struct {
	id      int // sql{inc}  #readonly
	uid     int // ch{primary}
	types   int // ch{low}
	name    string
	source  string
	ip      string
	image   bool
	stack   []string
	uints   []uint
	uints8  []uint8
	uints16 []uint16
	uints32 []uint32
	uints64 []uint64
	ints8   []int8
	ints16  []int16
	ints32  []int32
	ints64  []int64
}
