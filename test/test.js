export default{

	news: {
		id: undefined, //int sql{inc name="newid"} go{title="Account ID" name="AccountID" must}  #readonly
		created: undefined, //int64 #readonly js{time} sql{unix name="createdNew"}
		count: undefined, //int64 js{name="money"} title{NewTitleCount}
		active: false, //bool go{must} js{bool} swift{must} sql{default}
		kyc: false, //bool go{up} sql{type="newtype"} desc{KYC use for account validation}
		bid: undefined, //uint64 sql{skip} swift{skip} js{skip}
		oid: undefined, //int sql{unique="1" idx="i1"}
		type: undefined, //int sql{index idx="i1" idx="i2"}
		verify: false, //bool js{name="verified"}
		title: undefined, //string sql{search="search"} go{must}
		html: [], //[]byte sql{unique="1"}
		tags: [], //[]string sql{unique="2"}
		channels: [], //[]int sql{unique="1", unique="2"}
		channels64: [], //[]int64 go{must}
		floats: undefined, //float64 go{must} sql{primarykey}
		keys: {}, //map[string]string go{must} sql{primarykey}
		features: {}, //map[string]bool go{must} sql{primarykey}
		likes: {}, //map[string]int go{must}
		providers: {}, //map[int]string go{must}
		stats: {}, //map[int]int go{must}
		price: {}, //map[string]float64 go{must}
		meta: {}, //map[string]any sql{index}
		timeout: undefined, //time.Duration go{must}
		value: undefined, //any go{type="[]string"}
		raw: [], //[]byte go{raw} sql{name="rawbytes"}
	},

}
