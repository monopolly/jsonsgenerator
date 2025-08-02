type news = {
	id                                            :number // int sql{inc name="newid"} go{title="Account ID" name="AccountID" must}  #readonly
	created                                       :number // int64 #readonly js{time} sql{unix name="createdNew"}
	count                                         :number // int64 js{name="money"} title{NewTitleCount}
	active                                        :boolean // bool go{must} js{bool} swift{must} sql{default}
	kyc                                           :boolean // bool go{up} sql{type="newtype"} desc{KYC use for account validation}
	bid                                           :number // uint64 sql{skip} swift{skip} js{skip}
	oid                                           :number // int sql{unique="1" idx="i1"}
	type                                          :number // int sql{index idx="i1" idx="i2"}
	verify                                        :boolean // bool js{name="verified"}
	title                                         :string // string sql{search="search"} go{must}
	html                                          :string // []byte sql{unique="1"}
	tags                                          :string[] // []string sql{unique="2"}
	channels                                      :number[] // []int sql{unique="1", unique="2"}
	channels64                                    :number[] // []int64 go{must}
	floats                                        :number // float64 go{must} sql{primarykey}
	keys                                          :Record<string, string> // map[string]string go{must} sql{primarykey}
	features                                      :Record<string, boolean> // map[string]bool go{must} sql{primarykey}
	likes                                         :Record<string, number> // map[string]int go{must}
	providers                                     :Record<string, number> // map[int]string go{must}
	stats                                         :Record<string, number> // map[int]int go{must}
	price                                         :Record<string, number> // map[string]float64 go{must}
	meta                                          :Record<string, any> // map[string]any sql{index}
	timeout                                       :number // time.Duration go{must}
	value                                         :any // any go{type="[]string"}
	raw                                           :string // []byte go{raw} sql{name="rawbytes"}
}