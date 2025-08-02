package news

import (
	"time"
)

// go=News sql=news ts=news swift demo enum=Category noprefix msgp gotiny ch=news
type news struct {
	id         int                    // sql{inc name="newid"} go{title="Account ID" name="AccountID" must}  #readonly
	created    int64                  // #readonly js{time} sql{unix name="createdNew"}
	count      int64                  // js{name="money"} title{NewTitleCount}
	active     bool                   // go{must} js{bool} swift{must} sql{default}
	kyc        bool                   // go{up} sql{type="newtype"} desc{KYC use for account validation}
	bid        uint64                 // sql{skip} swift{skip} js{skip}
	oid        int                    // sql{unique="1" idx="i1"}
	types      int                    // sql{index idx="i1" idx="i2"}
	verify     bool                   // js{name="verified"}
	title      string                 // sql{search="search"} go{must}
	html       []byte                 // sql{unique="1"}
	tags       []string               // sql{unique="2"}
	channels   []int                  // sql{unique="1", unique="2"}
	channels64 []int64                // go{must}
	floats     float64                //go{must} sql{primarykey}
	keys       map[string]string      //go{must} sql{primarykey}
	features   map[string]bool        //go{must} sql{primarykey}
	likes      map[string]int         //go{must}
	providers  map[int]string         //go{must}
	stats      map[int]int            //go{must}
	price      map[string]float64     //go{must}
	meta       map[string]interface{} // sql{index}
	timeout    time.Duration          // go{must}
	value      any                    // go{type="[]string"}
	raw        []byte                 // go{raw} sql{name="rawbytes"}

	// sql{CREATE INDEX ON film USING GIN(mapAny)}
}
