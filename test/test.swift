struct NewsModel: Identifiable, Decodable, Encodable {
	var id          : Int? = 0 // int sql{inc name="newid"} go{title="Account ID" name="AccountID" must}  #readonly
	var created     : Int? = 0 // int64 #readonly js{time} sql{unix name="createdNew"}
	var count       : Int? = 0 // int64 js{name="money"} title{NewTitleCount}
	var active      : Bool = false // bool go{must} js{bool} swift{must} sql{default}
	var kyc         : Bool? = false // bool go{up} sql{type="newtype"} desc{KYC use for account validation}
	var oid         : Int? = 0 // int sql{unique="1" idx="i1"}
	var type        : Int? = 0 // int sql{index idx="i1" idx="i2"}
	var verify      : Bool? = false // bool js{name="verified"}
	var title       : String? = "" // string sql{search="search"} go{must}
	var html        : Data? = Data() // []byte sql{unique="1"}
	var tags        : [String]? = [String]() // []string sql{unique="2"}
	var channels    : [Int]? = [Int]() // []int sql{unique="1", unique="2"}
	var channels64  : [Int]? = [Int]() // []int64 go{must}
	var floats      : Double? = 0 // float64 go{must} sql{primarykey}
	var keys        : [String:String]? = [String:String]() // map[string]string go{must} sql{primarykey}
	var features    : [String:Bool]? = [String:Bool]() // map[string]bool go{must} sql{primarykey}
	var likes       : [String:Int]? = [String:Int]() // map[string]int go{must}
	var providers   : [Int:String]? = [Int:String]() // map[int]string go{must}
	var stats       : [Int:Int]? = [Int:Int]() // map[int]int go{must}
	var price       : Data? = Data() // map[string]float64 go{must}
	var meta        : [String:Any]? = [String:Any]() // map[string]any sql{index}
	var timeout     : Int? = 0 // time.Duration go{must}
	var value       : Data? = Data() // any go{type="[]string"}
	var raw         : Data? = Data() // []byte go{raw} sql{name="rawbytes"}
}