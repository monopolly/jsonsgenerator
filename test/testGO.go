package news

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	jsoniter "github.com/json-iterator/go"
	"github.com/monopolly/cast"
	"github.com/monopolly/jsons"
	"github.com/niubaoshu/gotiny"
	"github.com/vmihailenco/msgpack/v5"
)

/*
	Model Generator
	Help to auto-generated correct models for golang structures.
	Martin Prestone (c) 2024
	github.com/monopolly
	Ex: //! [] go=News js=NewsJson sql=news noinit up...

	STRUCT:
	debug                         Debug mode
	noinit                        No New() init function for struct
	gotiny                        Create gotiny marshal/unmarshal
	msgp                          Create message pack marshal/unmarshal
	!omit                         No omit tag for json
	go=News                       Set golang struct names. Ex: go=News1 > type News1 struct{}
	demo                          Generate demo json file with default values
	chengine=MergeTree            Optional, mergeTree by default
	noprefix                      Generate simple index IndexID instead IndexNewsID
	sql=news                      Set sql table name. Ex: sql=accounts
	lock                          Lock model for generation. Can't change model.
	js=NewsJson                   Change json struct names. Ex: js=NewsJson
	swift                         Generate swift model
	ts=news                       Set typescript struct names. Ex: ts=NewsJson
	enum                          Generate swift enum
	ch=views                      Set clickhouse sql table name. Ex: ch=views

	Field:
	title{}                       Add custom title for index. title{Nice & Sweet}
	desc{}                        Add custom desc for index. desc{This field for success}

	GO:
	#                             Add lists for fields. Ex: id int //#readonly #must...
	up                            Make uppercase for functions
	type=""                       Replace golang type. Ex: type="[]*News"
	name=""                       Replace golang struct name. Ex: name="NewsList"
	title=""                      Add custom title for index
	desc=""                       Add custom desc for index
	req{}                         Add user required fields
	title{}                       Add field user title title{Nice}
	nofunc                        Do not create any jsons functions for fiels
	must                          Create one validation function for all must fields
	desc{}                        Add field user title desc{Use it nice}

	SQL:
	defaults                      Add default value based on field type
	skip                          Skip field for sql queries
	type="jsonb"                  Rewrite sql type for field. Ex: type="jsonb"
	ver="v4"                      Create an alter table record in SQL file. Add new column. With new version in comment. Ex: ver="2"
	renames="oldname"             Create an alter table record in SQL file. Rename table column.
	index                         Create simple default index or gin for jsonb
	idx                           Add index fields by group name. Ex: idx="nameIndex" idx="credsIndex"
	search                        Add tsvector index by group. Ex: search="tsv": tsv tsvector GENERATED ALWAYS AS (to_tsvector('simple', title || ' ' || brand)). For search: SELECT brand, title FROM assets WHERE search @@ to_tsquery('english', 'f8');
	unix                          Add default value: extract(epoch from now())
	replace="bigint primary key"  Rewrite sql for field. Ex: replace="bigint primary key"
	add="primary key"             Append sql for field. Ex: replace="primary key"
	unique="groupname"            Add unique fields constrains by group (you need set group name, then generator join fields). Ex: unique="group1" unique="group2"
	noinsert                      Do use field for insert function
	altertable                    Add field line Alter table to SQL file with current time comment

	SWIFT:
	skip                          Ignore field for Swift
	type=""                       Replace type for Swift model. Ex: type="string"
	file=""                       Create another swift file for this field Swift model. file="f1", file="f2"
	must                          Swift field with required values

	JSON:
	bool                          Add set jsons function for bool fields
	time                          Create convert jsons function for unixtime fields
	raw                           Set Raw json function inside field
	name=""                       Replace json field name. Ex: name="sid"
	skip                          Skip json field
	inc                           Add inc jsons function for numbers fields
*/

// field type
type NewsIndexType int

// int index
const (
	IndexAccountID = NewsIndexType(iota)
	IndexCreated
	IndexCount
	IndexActive
	IndexKYC
	IndexBID
	IndexOID
	IndexType
	IndexVerify
	IndexTitle
	IndexHtml
	IndexTags
	IndexChannels
	IndexChannels64
	IndexFloats
	IndexKeys
	IndexFeatures
	IndexLikes
	IndexProviders
	IndexStats
	IndexPrice
	IndexMeta
	IndexTimeout
	IndexValue
	IndexRaw
)

// string index
const (
	FieldAccountID  = "id"         // int sql{inc name="newid"} go{title="Account ID" name="AccountID" must}  #readonly
	FieldCreated    = "created"    // int64 #readonly js{time} sql{unix name="createdNew"}
	FieldCount      = "money"      // int64 js{name="money"} title{NewTitleCount}
	FieldActive     = "active"     // bool go{must} js{bool} swift{must} sql{default}
	FieldKYC        = "kyc"        // bool go{up} sql{type="newtype"} desc{KYC use for account validation}
	FieldBID        = "bid"        // uint64 sql{skip} swift{skip} js{skip}
	FieldOID        = "oid"        // int sql{unique="1" idx="i1"}
	FieldType       = "type"       // int sql{index idx="i1" idx="i2"}
	FieldVerify     = "verified"   // bool js{name="verified"}
	FieldTitle      = "title"      // string sql{search="search"} go{must}
	FieldHtml       = "html"       // []byte sql{unique="1"}
	FieldTags       = "tags"       // []string sql{unique="2"}
	FieldChannels   = "channels"   // []int sql{unique="1", unique="2"}
	FieldChannels64 = "channels64" // []int64 go{must}
	FieldFloats     = "floats"     // float64 go{must} sql{primarykey}
	FieldKeys       = "keys"       // map[string]string go{must} sql{primarykey}
	FieldFeatures   = "features"   // map[string]bool go{must} sql{primarykey}
	FieldLikes      = "likes"      // map[string]int go{must}
	FieldProviders  = "providers"  // map[int]string go{must}
	FieldStats      = "stats"      // map[int]int go{must}
	FieldPrice      = "price"      // map[string]float64 go{must}
	FieldMeta       = "meta"       // map[string]any sql{index}
	FieldTimeout    = "timeout"    // time.Duration go{must}
	FieldValue      = "value"      // any go{type="[]string"}
	FieldRaw        = "raw"        // []byte go{raw} sql{name="rawbytes"}
)

// index func
func NewsIndexes() []NewsIndexType {
	return []NewsIndexType{IndexAccountID, IndexCreated, IndexCount, IndexActive, IndexKYC, IndexBID, IndexOID, IndexType, IndexVerify, IndexTitle, IndexHtml, IndexTags, IndexChannels, IndexChannels64, IndexFloats, IndexKeys, IndexFeatures, IndexLikes, IndexProviders, IndexStats, IndexPrice, IndexMeta, IndexTimeout, IndexValue, IndexRaw}

}

type News struct {
	AccountID  int                `json:"id,omitempty" msg:"id,omitempty"`                 // sql{inc name="newid"} go{title="Account ID" name="AccountID" must}  #readonly
	Created    int64              `json:"created,omitempty" msg:"created,omitempty"`       // #readonly js{time} sql{unix name="createdNew"}
	Count      int64              `json:"money,omitempty" msg:"money,omitempty"`           // js{name="money"} title{NewTitleCount}
	Active     bool               `json:"active,omitempty" msg:"active,omitempty"`         // go{must} js{bool} swift{must} sql{default}
	KYC        bool               `json:"kyc,omitempty" msg:"kyc,omitempty"`               // go{up} sql{type="newtype"} desc{KYC use for account validation}
	BID        uint64             `json:"bid,omitempty" msg:"bid,omitempty"`               // sql{skip} swift{skip} js{skip}
	OID        int                `json:"oid,omitempty" msg:"oid,omitempty"`               // sql{unique="1" idx="i1"}
	Type       int                `json:"type,omitempty" msg:"type,omitempty"`             // sql{index idx="i1" idx="i2"}
	Verify     bool               `json:"verified,omitempty" msg:"verified,omitempty"`     // js{name="verified"}
	Title      string             `json:"title,omitempty" msg:"title,omitempty"`           // sql{search="search"} go{must}
	Html       []byte             `json:"html,omitempty" msg:"html,omitempty"`             // sql{unique="1"}
	Tags       []string           `json:"tags,omitempty" msg:"tags,omitempty"`             // sql{unique="2"}
	Channels   []int              `json:"channels,omitempty" msg:"channels,omitempty"`     // sql{unique="1", unique="2"}
	Channels64 []int64            `json:"channels64,omitempty" msg:"channels64,omitempty"` // go{must}
	Floats     float64            `json:"floats,omitempty" msg:"floats,omitempty"`         // go{must} sql{primarykey}
	Keys       map[string]string  `json:"keys,omitempty" msg:"keys,omitempty"`             // go{must} sql{primarykey}
	Features   map[string]bool    `json:"features,omitempty" msg:"features,omitempty"`     // go{must} sql{primarykey}
	Likes      map[string]int     `json:"likes,omitempty" msg:"likes,omitempty"`           // go{must}
	Providers  map[int]string     `json:"providers,omitempty" msg:"providers,omitempty"`   // go{must}
	Stats      map[int]int        `json:"stats,omitempty" msg:"stats,omitempty"`           // go{must}
	Price      map[string]float64 `json:"price,omitempty" msg:"price,omitempty"`           // go{must}
	Meta       map[string]any     `json:"meta,omitempty" msg:"meta,omitempty"`             // sql{index}
	Timeout    time.Duration      `json:"timeout,omitempty" msg:"timeout,omitempty"`       // go{must}
	Value      []string           `json:"value,omitempty" msg:"value,omitempty"`           // go{type="[]string"}
	Raw        []byte             `json:"raw,omitempty" msg:"raw,omitempty"`               // go{raw} sql{name="rawbytes"}
}

// Parse []any to ID struct
func ParseNewsToStruct(r []any) (a *News) {
	a = new(News)

	for pos, x := range r {
		switch NewsIndexType(pos) {
		case IndexAccountID:
			a.AccountID = cast.Int(x) //int
		case IndexCreated:
			a.Created = cast.Int64(x) //int64
		case IndexCount:
			a.Count = cast.Int64(x) //int64
		case IndexActive:
			a.Active = cast.Bool(x) //bool
		case IndexKYC:
			a.KYC = cast.Bool(x) //bool
		case IndexBID:
			a.BID = cast.Uint64(x) //uint64
		case IndexOID:
			a.OID = cast.Int(x) //int
		case IndexType:
			a.Type = cast.Int(x) //int
		case IndexVerify:
			a.Verify = cast.Bool(x) //bool
		case IndexTitle:
			a.Title = cast.String(x) //string
		case IndexHtml:
			a.Html = cast.Bytes(x) //[]byte
		case IndexTags:
			a.Tags = cast.SliceString(x) //[]string
		case IndexChannels:
			a.Channels = cast.SliceInt(x) //[]int
		case IndexChannels64:
			a.Channels64 = cast.SliceInt64(x) //[]int64
		case IndexFloats:
			a.Floats = cast.Float(x) //float64
		case IndexKeys:
			a.Keys = cast.StringMapString(x) //map[string]string
		case IndexFeatures:
			a.Features = cast.StringMapBool(x) //map[string]bool
		case IndexLikes:
			a.Likes = cast.MapStringInt(x) //map[string]int
		case IndexProviders:
			a.Providers = cast.MapIntString(x) //map[int]string
		case IndexStats:
			a.Stats = cast.MapIntInt(x) //map[int]int
		case IndexPrice:
			a.Price = cast.MapStringFloats(x) //map[string]float64
		case IndexMeta:
			a.Meta = cast.StringMap(x) //map[string]any
		case IndexTimeout:
			a.Timeout = cast.Duration(x) //time.Duration
		case IndexRaw:
			a.Raw = cast.Bytes(x) //[]byte
		}
	}
	return
}

// Tuple create an array from struct
func (a *News) Tuple() (r []any) {
	return []any{a.AccountID, a.Created, a.Count, a.Active, a.KYC, a.BID, a.OID, a.Type, a.Verify, a.Title, a.Html, a.Tags, a.Channels, a.Channels64, a.Floats, a.Keys, a.Features, a.Likes, a.Providers, a.Stats, a.Price, a.Meta, a.Timeout, a.Value, a.Raw}
}

// Tuple create an array from struct
func (a *News) sqlTuple() (r []any) {
	return []any{a.Created, a.Count, a.Active, a.KYC, a.OID, a.Type, a.Verify, a.Title, a.Html, a.Tags, a.Channels, a.Channels64, a.Floats, a.Keys, a.Features, a.Likes, a.Providers, a.Stats, a.Price, a.Meta, a.Timeout, a.Value, a.Raw}
}

// update struct with function
func (a *News) Update(k string, x any) {
	switch k {
	case "id":
		a.AccountID = cast.Int(x) //int
	case "created":
		a.Created = cast.Int64(x) //int64
	case "count":
		a.Count = cast.Int64(x) //int64
	case "active":
		a.Active = cast.Bool(x) //bool
	case "kyc":
		a.KYC = cast.Bool(x) //bool
	case "bid":
		a.BID = cast.Uint64(x) //uint64
	case "oid":
		a.OID = cast.Int(x) //int
	case "type":
		a.Type = cast.Int(x) //int
	case "verify":
		a.Verify = cast.Bool(x) //bool
	case "title":
		a.Title = cast.String(x) //string
	case "html":
		a.Html = cast.Bytes(x) //[]byte
	case "tags":
		a.Tags = cast.SliceString(x) //[]string
	case "channels":
		a.Channels = cast.SliceInt(x) //[]int
	case "channels64":
		a.Channels64 = cast.SliceInt64(x) //[]int64
	case "floats":
		a.Floats = cast.Float(x) //float64
	case "keys":
		a.Keys = cast.StringMapString(x) //map[string]string
	case "features":
		a.Features = cast.StringMapBool(x) //map[string]bool
	case "likes":
		a.Likes = cast.MapStringInt(x) //map[string]int
	case "providers":
		a.Providers = cast.MapIntString(x) //map[int]string
	case "stats":
		a.Stats = cast.MapIntInt(x) //map[int]int
	case "price":
		a.Price = cast.MapStringFloats(x) //map[string]float64
	case "meta":
		a.Meta = cast.StringMap(x) //map[string]any
	case "timeout":
		a.Timeout = cast.Duration(x) //time.Duration
	case "raw":
		a.Raw = cast.Bytes(x) //[]byte
	}
}

// get struct value with function
func (a *News) Get(k string) (v any) {
	switch k {
	case "id":
		return a.AccountID //int
	case "created":
		return a.Created //int64
	case "count":
		return a.Count //int64
	case "active":
		return a.Active //bool
	case "kyc":
		return a.KYC //bool
	case "bid":
		return a.BID //uint64
	case "oid":
		return a.OID //int
	case "type":
		return a.Type //int
	case "verify":
		return a.Verify //bool
	case "title":
		return a.Title //string
	case "html":
		return a.Html //[]byte
	case "tags":
		return a.Tags //[]string
	case "channels":
		return a.Channels //[]int
	case "channels64":
		return a.Channels64 //[]int64
	case "floats":
		return a.Floats //float64
	case "keys":
		return a.Keys //map[string]string
	case "features":
		return a.Features //map[string]bool
	case "likes":
		return a.Likes //map[string]int
	case "providers":
		return a.Providers //map[int]string
	case "stats":
		return a.Stats //map[int]int
	case "price":
		return a.Price //map[string]float64
	case "meta":
		return a.Meta //map[string]any
	case "timeout":
		return a.Timeout //time.Duration
	case "value":
		return a.Value //any
	case "raw":
		return a.Raw //[]byte
	}
	return
}

// get any struct value as string
func (a *News) String(k string) (v string) {
	switch k {
	case "id":
		return fmt.Sprint(a.AccountID) //int
	case "created":
		return fmt.Sprint(a.Created) //int64
	case "count":
		return fmt.Sprint(a.Count) //int64
	case "active":
		return fmt.Sprint(a.Active) //bool
	case "kyc":
		return fmt.Sprint(a.KYC) //bool
	case "bid":
		return fmt.Sprint(a.BID) //uint64
	case "oid":
		return fmt.Sprint(a.OID) //int
	case "type":
		return fmt.Sprint(a.Type) //int
	case "verify":
		return fmt.Sprint(a.Verify) //bool
	case "title":
		return fmt.Sprint(a.Title) //string
	case "html":
		return fmt.Sprint(a.Html) //[]byte
	case "tags":
		return fmt.Sprint(a.Tags) //[]string
	case "channels":
		return fmt.Sprint(a.Channels) //[]int
	case "channels64":
		return fmt.Sprint(a.Channels64) //[]int64
	case "floats":
		return fmt.Sprint(a.Floats) //float64
	case "keys":
		return fmt.Sprint(a.Keys) //map[string]string
	case "features":
		return fmt.Sprint(a.Features) //map[string]bool
	case "likes":
		return fmt.Sprint(a.Likes) //map[string]int
	case "providers":
		return fmt.Sprint(a.Providers) //map[int]string
	case "stats":
		return fmt.Sprint(a.Stats) //map[int]int
	case "price":
		return fmt.Sprint(a.Price) //map[string]float64
	case "meta":
		return fmt.Sprint(a.Meta) //map[string]any
	case "timeout":
		return fmt.Sprint(a.Timeout) //time.Duration
	case "value":
		return fmt.Sprint(a.Value) //any
	case "raw":
		return fmt.Sprint(a.Raw) //[]byte
	}
	return
}

// Struct to json
func (a *News) ToJson() (r []byte) {
	js := jsons.Create().
		Add(FieldAccountID, a.AccountID).
		Add(FieldCreated, a.Created).
		Add(FieldCount, a.Count).
		Add(FieldActive, a.Active).
		Add(FieldKYC, a.KYC).
		Add(FieldBID, a.BID).
		Add(FieldOID, a.OID).
		Add(FieldType, a.Type).
		Add(FieldVerify, a.Verify).
		Add(FieldTitle, a.Title).
		Add(FieldHtml, a.Html).
		Add(FieldTags, a.Tags).
		Add(FieldChannels, a.Channels).
		Add(FieldChannels64, a.Channels64).
		Add(FieldFloats, a.Floats).
		Add(FieldKeys, a.Keys).
		Add(FieldFeatures, a.Features).
		Add(FieldLikes, a.Likes).
		Add(FieldProviders, a.Providers).
		Add(FieldStats, a.Stats).
		Add(FieldPrice, a.Price).
		Add(FieldMeta, a.Meta).
		Add(FieldTimeout, a.Timeout).
		Add(FieldValue, a.Value).
		Add(FieldRaw, a.Raw)
	return js.Bytes()
}

// Valid empty values for selected fields by go{must}
func (a *News) Must() (emptyfields []string) {
	if a.AccountID < 1 {
		emptyfields = append(emptyfields, "id")
	}
	if a.Active == false {
		emptyfields = append(emptyfields, "active")
	}
	if a.Title == "" {
		emptyfields = append(emptyfields, "title")
	}
	if a.Channels64 == nil {
		emptyfields = append(emptyfields, "channels64")
	}
	if a.Floats < 1 {
		emptyfields = append(emptyfields, "floats")
	}
	if a.Keys == nil {
		emptyfields = append(emptyfields, "keys")
	}
	if a.Features == nil {
		emptyfields = append(emptyfields, "features")
	}
	if a.Likes == nil {
		emptyfields = append(emptyfields, "likes")
	}
	if a.Providers == nil {
		emptyfields = append(emptyfields, "providers")
	}
	if a.Stats == nil {
		emptyfields = append(emptyfields, "stats")
	}
	if a.Price == nil {
		emptyfields = append(emptyfields, "price")
	}
	if a.Timeout < 1 {
		emptyfields = append(emptyfields, "timeout")
	}
	return
}

func NewsReadonlyList() []NewsIndexType {
	return []NewsIndexType{IndexAccountID, IndexCreated}
}

func (a NewsIndexType) Readonly() bool {
	switch a {
	case IndexAccountID, IndexCreated:
		return true
	default:
		return false
	}
}

// key index string
func (a NewsIndexType) String() string {
	switch a {
	case IndexAccountID:
		return "id"
	case IndexCreated:
		return "created"
	case IndexCount:
		return "count"
	case IndexActive:
		return "active"
	case IndexKYC:
		return "kyc"
	case IndexBID:
		return "bid"
	case IndexOID:
		return "oid"
	case IndexType:
		return "type"
	case IndexVerify:
		return "verify"
	case IndexTitle:
		return "title"
	case IndexHtml:
		return "html"
	case IndexTags:
		return "tags"
	case IndexChannels:
		return "channels"
	case IndexChannels64:
		return "channels64"
	case IndexFloats:
		return "floats"
	case IndexKeys:
		return "keys"
	case IndexFeatures:
		return "features"
	case IndexLikes:
		return "likes"
	case IndexProviders:
		return "providers"
	case IndexStats:
		return "stats"
	case IndexPrice:
		return "price"
	case IndexMeta:
		return "meta"
	case IndexTimeout:
		return "timeout"
	case IndexValue:
		return "value"
	case IndexRaw:
		return "raw"
	default:
		return ""
	}
}

// key index string
func (a NewsIndexType) SQLName() string {
	switch a {
	case IndexAccountID:
		return "newid"
	case IndexCreated:
		return "createdNew"
	case IndexCount:
		return "count"
	case IndexActive:
		return "active"
	case IndexKYC:
		return "kyc"
	case IndexBID:
		return "bid"
	case IndexOID:
		return "oid"
	case IndexType:
		return "type"
	case IndexVerify:
		return "verify"
	case IndexTitle:
		return "title"
	case IndexHtml:
		return "html"
	case IndexTags:
		return "tags"
	case IndexChannels:
		return "channels"
	case IndexChannels64:
		return "channels64"
	case IndexFloats:
		return "floats"
	case IndexKeys:
		return "keys"
	case IndexFeatures:
		return "features"
	case IndexLikes:
		return "likes"
	case IndexProviders:
		return "providers"
	case IndexStats:
		return "stats"
	case IndexPrice:
		return "price"
	case IndexMeta:
		return "meta"
	case IndexTimeout:
		return "timeout"
	case IndexValue:
		return "value"
	case IndexRaw:
		return "rawbytes"
	default:
		return ""
	}
}

// key index type
func (a NewsIndexType) Type() string {
	switch a {
	case IndexAccountID:
		return "int"
	case IndexCreated:
		return "int64"
	case IndexCount:
		return "int64"
	case IndexActive:
		return "bool"
	case IndexKYC:
		return "bool"
	case IndexBID:
		return "uint64"
	case IndexOID:
		return "int"
	case IndexType:
		return "int"
	case IndexVerify:
		return "bool"
	case IndexTitle:
		return "string"
	case IndexHtml:
		return "[]byte"
	case IndexTags:
		return "[]string"
	case IndexChannels:
		return "[]int"
	case IndexChannels64:
		return "[]int64"
	case IndexFloats:
		return "float64"
	case IndexKeys:
		return "map[string]string"
	case IndexFeatures:
		return "map[string]bool"
	case IndexLikes:
		return "map[string]int"
	case IndexProviders:
		return "map[int]string"
	case IndexStats:
		return "map[int]int"
	case IndexPrice:
		return "map[string]float64"
	case IndexMeta:
		return "map[string]any"
	case IndexTimeout:
		return "time.Duration"
	case IndexValue:
		return "any"
	case IndexRaw:
		return "[]byte"
	default:
		return ""
	}
}

// custom title
func (a NewsIndexType) Title() string {
	switch a {
	case IndexAccountID:
		return "Account ID"
	case IndexCreated:
		return "Created"
	case IndexCount:
		return "Count"
	case IndexActive:
		return "Active"
	case IndexKYC:
		return "KYC"
	case IndexBID:
		return "BID"
	case IndexOID:
		return "OID"
	case IndexType:
		return "Type"
	case IndexVerify:
		return "Verify"
	case IndexTitle:
		return "Title"
	case IndexHtml:
		return "Html"
	case IndexTags:
		return "Tags"
	case IndexChannels:
		return "Channels"
	case IndexChannels64:
		return "Channels64"
	case IndexFloats:
		return "Floats"
	case IndexKeys:
		return "Keys"
	case IndexFeatures:
		return "Features"
	case IndexLikes:
		return "Likes"
	case IndexProviders:
		return "Providers"
	case IndexStats:
		return "Stats"
	case IndexPrice:
		return "Price"
	case IndexMeta:
		return "Meta"
	case IndexTimeout:
		return "Timeout"
	case IndexValue:
		return "Value"
	case IndexRaw:
		return "Raw"
	default:
		return ""
	}
}

// custom desc
func (a NewsIndexType) Desc() string {
	switch a {
	case IndexKYC:
		return "KYC use for account validation"
	default:
		return ""
	}
}

// struct key to index
func NewsKeyIndex(key string) NewsIndexType {
	switch key {
	case "id":
		return IndexAccountID
	case "created":
		return IndexCreated
	case "count":
		return IndexCount
	case "active":
		return IndexActive
	case "kyc":
		return IndexKYC
	case "bid":
		return IndexBID
	case "oid":
		return IndexOID
	case "type":
		return IndexType
	case "verify":
		return IndexVerify
	case "title":
		return IndexTitle
	case "html":
		return IndexHtml
	case "tags":
		return IndexTags
	case "channels":
		return IndexChannels
	case "channels64":
		return IndexChannels64
	case "floats":
		return IndexFloats
	case "keys":
		return IndexKeys
	case "features":
		return IndexFeatures
	case "likes":
		return IndexLikes
	case "providers":
		return IndexProviders
	case "stats":
		return IndexStats
	case "price":
		return IndexPrice
	case "meta":
		return IndexMeta
	case "timeout":
		return IndexTimeout
	case "value":
		return IndexValue
	case "raw":
		return IndexRaw
	default:
		return 0
	}
}

// valid struct key check
func NewsValidKey(key string) bool {
	switch key {
	case "id", "created", "count", "active", "kyc", "bid", "oid", "type", "verify", "title", "html", "tags", "channels", "channels64", "floats", "keys", "features", "likes", "providers", "stats", "price", "meta", "timeout", "value", "raw":
		return true
	default:
		return false
	}
}

// struct to map
func (a *News) Map() map[string]any {
	return map[string]any{
		"id":         a.AccountID,
		"created":    a.Created,
		"money":      a.Count,
		"active":     a.Active,
		"kyc":        a.KYC,
		"bid":        a.BID,
		"oid":        a.OID,
		"type":       a.Type,
		"verified":   a.Verify,
		"title":      a.Title,
		"html":       a.Html,
		"tags":       a.Tags,
		"channels":   a.Channels,
		"channels64": a.Channels64,
		"floats":     a.Floats,
		"keys":       a.Keys,
		"features":   a.Features,
		"likes":      a.Likes,
		"providers":  a.Providers,
		"stats":      a.Stats,
		"price":      a.Price,
		"meta":       a.Meta,
		"timeout":    a.Timeout,
		"value":      a.Value,
		"raw":        a.Raw,
	}
}

// struct to map
func (a *News) Iterate(f func(k NewsIndexType, v any)) {
	for _, x := range NewsIndexes() {
		f(x, a.Get(x.String()))
	}
}

// gotiny marshal
func (a *News) Gotiny() []byte {
	return gotiny.Marshal(&a.AccountID, &a.Created, &a.Count, &a.Active, &a.KYC, &a.BID, &a.OID, &a.Type, &a.Verify, &a.Title, &a.Html, &a.Tags, &a.Channels, &a.Channels64, &a.Floats, &a.Keys, &a.Features, &a.Likes, &a.Providers, &a.Stats, &a.Price, &a.Meta, &a.Timeout, &a.Value, &a.Raw)
}

// parse gotiny
func ParseNewsGotiny(v []byte) (a News) {
	gotiny.Unmarshal(v, &a.AccountID, &a.Created, &a.Count, &a.Active, &a.KYC, &a.BID, &a.OID, &a.Type, &a.Verify, &a.Title, &a.Html, &a.Tags, &a.Channels, &a.Channels64, &a.Floats, &a.Keys, &a.Features, &a.Likes, &a.Providers, &a.Stats, &a.Price, &a.Meta, &a.Timeout, &a.Value, &a.Raw)
	return
}

// msgp marshal
func (a *News) MessagePack() []byte {
	b, _ := msgpack.Marshal(a)
	return b
}

// msgp unmarshal
func ParseNewsMessagePack(v []byte) (a News, err error) {
	err = msgpack.Unmarshal(v, &a)
	return
}

// fast json marshal
func (a *News) Pack() []byte {
	var jsoner = jsoniter.ConfigCompatibleWithStandardLibrary
	b, _ := jsoner.Marshal(a)
	return b
}

// fast json unmarshal
func ParseNews(v []byte) (a News, err error) {
	var jsoner = jsoniter.ConfigCompatibleWithStandardLibrary
	err = jsoner.Unmarshal(v, &a)
	return
}

// sql NewsSQL class
type NewsSQL struct{ pool *pgxpool.Pool }

func NewNewsSQL(pool *pgxpool.Pool) (a *NewsSQL) {
	a = new(NewsSQL)
	a.pool = pool
	return
}

//delete item

func (a *NewsSQL) Conn(f func(conn *pgxpool.Conn) (err error)) (err error) {

	conn, err := a.pool.Acquire(context.Background())
	if err != nil {
		return
	}
	defer conn.Release()
	return f(conn)
}

// parse sql query
func (a *NewsSQL) TableName() (res string) {
	return "news"
}

// parse sql query
func (a *NewsSQL) Get(id any, fields ...NewsIndexType) (res *News, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if fields == nil {
		fields = a.Fields()
	}

	var list []string
	for _, x := range fields {
		list = append(list, x.SQLName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s ", fieldlist)
	q = q + "from news where newid = $1 limit 1"
	res = new(News)
	rows, err := conn.Query(c, q, id)
	if err != nil {
		return
	}
	defer rows.Close()
	rows.Next()
	v, err := rows.Values()
	if err != nil {
		return
	}
	if len(v) != len(fields) {
		err = fmt.Errorf("len")
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), v[pos])
	}
	return
}

// parse sql query
func (a *NewsSQL) Row(eq map[string]any, fields ...NewsIndexType) (res *News, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if fields == nil {
		fields = a.Fields()
	}

	var list []string
	for _, x := range fields {
		list = append(list, x.SQLName())
	}
	fieldlist := strings.Join(list, ", ")

	var andlist []string
	var values []any
	var count int
	for k, v := range eq {
		count++
		andlist = append(andlist, fmt.Sprintf("%s = $%d", k, count))
		values = append(values, v)
	}
	keys := strings.Join(andlist, " and ")

	q := fmt.Sprintf("select %s from %s where %s limit 1", fieldlist, a.TableName(), keys)
	res = new(News)
	rows, err := conn.Query(c, q, values...)
	if err != nil {
		return
	}
	defer rows.Close()
	rows.Next()
	v, err := rows.Values()
	if err != nil {
		return
	}
	if len(v) != len(fields) {
		err = fmt.Errorf("len")
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), v[pos])
	}
	return
}

// get simple list from db
func (a *NewsSQL) All(fields ...NewsIndexType) (res []*News, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if fields == nil {
		fields = a.Fields()
	}

	var list []string
	for _, x := range fields {
		list = append(list, x.SQLName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s", fieldlist)
	q += " from news"
	rows, err := conn.Query(c, q)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item News
		v, err := rows.Values()
		if err != nil {
			continue
		}
		if len(v) != len(fields) {
			continue
		}
		for pos, x := range fields {
			item.Update(x.String(), v[pos])
		}
		res = append(res, &item)
	}

	return
}

// get simple list from db
func (a *NewsSQL) List(limit, offset int, fields ...NewsIndexType) (res []*News, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if fields == nil {
		fields = []NewsIndexType{IndexAccountID, IndexCreated, IndexCount, IndexActive, IndexKYC, IndexBID, IndexOID, IndexType, IndexVerify, IndexTitle, IndexHtml, IndexTags, IndexChannels, IndexChannels64, IndexFloats, IndexKeys, IndexFeatures, IndexLikes, IndexProviders, IndexStats, IndexPrice, IndexMeta, IndexTimeout, IndexValue, IndexRaw}
	}

	var list []string
	for _, x := range fields {
		list = append(list, x.SQLName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s", fieldlist)
	q += " from news limit $1 offset $2"
	rows, err := conn.Query(c, q, limit, offset)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item News
		v, err := rows.Values()
		if err != nil {
			continue
		}
		if len(v) != len(fields) {
			continue
		}
		for pos, x := range fields {
			item.Update(x.String(), v[pos])
		}
		res = append(res, &item)
	}

	return
}

// update sql query
func (a *NewsSQL) Update(id any, k string, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if !NewsValidKey(k) {
		return fmt.Errorf("invalid key")
	}
	q := fmt.Sprintf("update news set %s = $1 where newid = $2", k)
	_, err = conn.Exec(c, q, v, id)
	return
}

// update sql query
func (a *NewsSQL) Updates(id any, keys map[string]any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if keys == nil {
		return fmt.Errorf("emptykeys")
	}

	var fields []string
	var values []any
	var count int
	for k, v := range keys {
		if !NewsValidKey(k) {
			return fmt.Errorf(k)
		}
		in := NewsKeyIndex(k)
		count++
		fields = append(fields, fmt.Sprintf("%s = $%d", in.SQLName(), count))
		values = append(values, v)
	}

	list := strings.Join(fields, ", ")
	count++
	values = append(values, id)

	q := fmt.Sprintf("update news set %s where newid = $%d", list, count)
	_, err = conn.Exec(c, q, values...)
	return
}

// add array value to jsonb array
func (a *NewsSQL) AddTags(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set tags = tags || '%s'::jsonb where newid = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *NewsSQL) DeleteTags(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update news set tags = tags - $1 where newid = $2")
	_, err = conn.Exec(c, q, v, id)
	return
}

// add array value to jsonb array
func (a *NewsSQL) AddChannels(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set channels = channels || '%s'::jsonb where newid = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *NewsSQL) DeleteChannels(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update news set channels = channels - $1 where newid = $2")
	_, err = conn.Exec(c, q, v, id)
	return
}

// add array value to jsonb array
func (a *NewsSQL) AddChannels64(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Create().Array(v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set channels64 = channels64 || '%s'::jsonb where newid = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete array value from jsonb array
func (a *NewsSQL) DeleteChannels64(id any, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("update news set channels64 = channels64 - $1 where newid = $2")
	_, err = conn.Exec(c, q, v, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateKeys(id any, k string, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set keys = keys || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesKeys(id any, keys map[string]any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := json.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set keys = keys || $$%s$$::jsonb where newid = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyKeys(id any, k string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set keys = keys - $1 where id = $2"
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyKeys(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set keys = keys - $1 || jsonb_build_object($2, keys->$1) where id = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateFeatures(id any, k string, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set features = features || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesFeatures(id any, keys map[string]any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := json.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set features = features || $$%s$$::jsonb where newid = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyFeatures(id any, k string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set features = features - $1 where id = $2"
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyFeatures(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set features = features - $1 || jsonb_build_object($2, features->$1) where id = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateLikes(id any, k string, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set likes = likes || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesLikes(id any, keys map[string]any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := json.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set likes = likes || $$%s$$::jsonb where newid = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyLikes(id any, k string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set likes = likes - $1 where id = $2"
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyLikes(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set likes = likes - $1 || jsonb_build_object($2, likes->$1) where id = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateProviders(id any, k string, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set providers = providers || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesProviders(id any, keys map[string]any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := json.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set providers = providers || $$%s$$::jsonb where newid = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyProviders(id any, k string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set providers = providers - $1 where id = $2"
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyProviders(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set providers = providers - $1 || jsonb_build_object($2, providers->$1) where id = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateStats(id any, k string, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set stats = stats || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesStats(id any, keys map[string]any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := json.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set stats = stats || $$%s$$::jsonb where newid = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyStats(id any, k string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set stats = stats - $1 where id = $2"
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyStats(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set stats = stats - $1 || jsonb_build_object($2, stats->$1) where id = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatePrice(id any, k string, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set price = price || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesPrice(id any, keys map[string]any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := json.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set price = price || $$%s$$::jsonb where newid = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyPrice(id any, k string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set price = price - $1 where id = $2"
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyPrice(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set price = price - $1 || jsonb_build_object($2, price->$1) where id = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateMeta(id any, k string, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set meta = meta || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesMeta(id any, keys map[string]any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := json.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set meta = meta || $$%s$$::jsonb where newid = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyMeta(id any, k string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set meta = meta - $1 where id = $2"
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyMeta(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set meta = meta - $1 || jsonb_build_object($2, meta->$1) where id = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *NewsSQL) UpdateValue(id any, k string, v any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update news set value = value || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *NewsSQL) UpdatesValue(id any, keys map[string]any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	b, _ := json.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update news set value = value || $$%s$$::jsonb where newid = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *NewsSQL) DeleteKeyValue(id any, k string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set value = value - $1 where id = $2"
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *NewsSQL) RenameKeyValue(id any, k, newkey string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "update news set value = value - $1 || jsonb_build_object($2, value->$1) where id = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

type NewsQuery struct {
	Limit  int                 `json:"limit,omitempty"`
	Offset int                 `json:"offset,omitempty"`
	Sort   string              `json:"sort,omitempty"`
	Desc   bool                `json:"desc,omitempty"`
	EQ     map[string]any      `json:"eq,omitempty"`   //equal
	GT     map[string]any      `json:"gt,omitempty"`   //greater then...
	LT     map[string]any      `json:"lt,omitempty"`   //less then
	NOT    map[string][]any    `json:"not,omitempty"`  //less then
	Like   map[string]string   `json:"like,omitempty"` //full text search
	Custom string              `json:"-"`              //append unsafe where condition
	Fields []string            `json:"fields,omitempty"`
	IN     map[string][]int    `json:"in,omitempty"`
	INS    map[string][]string `json:"ins,omitempty"`
	INSQL  map[string]string   `json:"insql,omitempty"` //unsafe in (condition)
}

func (a *NewsQuery) Render() (sql string, fields []NewsIndexType, values []any) {

	switch a.Fields == nil {
	case true:
		fields = []NewsIndexType{IndexAccountID, IndexCreated, IndexCount, IndexActive, IndexKYC, IndexBID, IndexOID, IndexType, IndexVerify, IndexTitle, IndexHtml, IndexTags, IndexChannels, IndexChannels64, IndexFloats, IndexKeys, IndexFeatures, IndexLikes, IndexProviders, IndexStats, IndexPrice, IndexMeta, IndexTimeout, IndexValue, IndexRaw}
	default:
		for _, x := range a.Fields {
			if !NewsValidKey(x) {
				continue
			}
			fields = append(fields, NewsKeyIndex(x))
		}
	}

	var fieldsStrings []string
	for _, x := range fields {
		fieldsStrings = append(fieldsStrings, x.String())
	}

	for k, v := range a.EQ {
		if fmt.Sprint(v) == "" {
			delete(a.EQ, k)
		}
	}
	for k, v := range a.GT {
		if fmt.Sprint(v) == "" {
			delete(a.GT, k)
		}
	}
	for k, v := range a.LT {
		if fmt.Sprint(v) == "" {
			delete(a.LT, k)
		}
	}
	for k, v := range a.NOT {
		if v == nil {
			delete(a.NOT, k)
		}
	}
	for k, v := range a.Like {
		if v == "" {
			delete(a.Like, k)
		}
	}

	// sql
	var list []string

	var count int

	// select
	list = append(list, "select")
	list = append(list, strings.Join(fieldsStrings, ", "))

	// from
	list = append(list, "from news")

	var andlist []string

	// EQ where
	for k, v := range a.EQ {
		if !NewsValidKey(k) {
			continue
		}
		p := NewsKeyIndex(k)
		switch p.Type() {
		case "bool":
			switch cast.Bool(v) {
			case true:
				andlist = append(andlist, k)
			case false:
				andlist = append(andlist, fmt.Sprintf("(not %[1]s or %[1]s is null)", k))
			}

		default:
			count++
			andlist = append(andlist, fmt.Sprintf("%s = $%d", k, count))
			values = append(values, v)
		}
	}

	// GT where
	for k, v := range a.GT {

		if !NewsValidKey(k) {
			continue
		}
		count++
		andlist = append(andlist, fmt.Sprintf("%s > $%d", k, count))
		values = append(values, v)
	}

	// LT where
	for k, v := range a.LT {
		if !NewsValidKey(k) {
			continue
		}
		count++
		andlist = append(andlist, fmt.Sprintf("%s < $%d", k, count))
		values = append(values, v)
	}

	// NOT where
	for k, v := range a.NOT {
		if !NewsValidKey(k) {
			continue
		}

		for _, x := range v {
			count++
			andlist = append(andlist, fmt.Sprintf("%s != $%d", k, count))
			values = append(values, x)
		}

	}

	// LIKE where
	for k, v := range a.Like {
		if !NewsValidKey(k) {
			continue
		}
		count++
		andlist = append(andlist, fmt.Sprintf("%s ilike $%d", k, count))
		values = append(values, "%"+v+"%")
	}

	// IN where
	if a.IN != nil {
		for k, v := range a.IN {
			if !NewsValidKey(k) {
				continue
			}
			var inlist []string
			for _, num := range v {
				inlist = append(inlist, fmt.Sprint(num))
			}
			count++
			andlist = append(andlist, fmt.Sprintf("%s in (%s)", k, strings.Join(inlist, ",")))
			// values = append(values, string)
		}
	}

	// INSQ where
	if a.INSQL != nil {
		for k, v := range a.INSQL {
			if !NewsValidKey(k) {
				continue
			}
			count++
			andlist = append(andlist, fmt.Sprintf("%s in (%s)", k, v))
		}
	}

	if a.Custom != "" {
		andlist = append(andlist, a.Custom)
	}

	// render where
	if andlist != nil {
		list = append(list, "where")
		list = append(list, strings.Join(andlist, " and "))
	}

	// sort by
	if a.Sort != "" && NewsValidKey(a.Sort) {
		list = append(list, fmt.Sprintf("order by %s", a.Sort))
		if a.Desc {
			list = append(list, "desc")
		}
	}

	// limit, offset
	if a.Limit > 0 {
		list = append(list, fmt.Sprintf("limit %d", a.Limit))
	}
	if a.Offset > 0 {
		list = append(list, fmt.Sprintf("offset %d", a.Offset))
	}

	// render sql
	sql = strings.Join(list, " ")
	return
}

func NewNewsQuery() *NewsQuery {
	a := new(NewsQuery)
	a.EQ = make(map[string]any)
	a.GT = make(map[string]any)
	a.LT = make(map[string]any)
	a.NOT = make(map[string][]any)
	a.Like = make(map[string]string)
	a.IN = make(map[string][]int)
	a.INS = make(map[string][]string)
	a.INSQL = make(map[string]string)
	return a
}

func (a *NewsSQL) Search(q *NewsQuery) (res []*News, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	sql, fields, values := q.Render()

	rows, err := conn.Query(c, sql, values...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item News
		v, err := rows.Values()
		if err != nil {
			continue
		}

		for pos, x := range fields {
			item.Update(x.String(), v[pos])
		}
		res = append(res, &item)
	}

	return
}

// delete item
func (a *NewsSQL) Delete(id any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	_, err = conn.Exec(c, "delete from news where newid = $1", id)
	return

}

// delete item where i = 1 and w = 'nice'
func (a *NewsSQL) DeleteWhere(where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	_, err = conn.Exec(c, fmt.Sprintf("delete from news where %s", where))
	return

}

// Insert struct and return int id
func (a *NewsSQL) Insert(v *News) (id int, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "insert into news (createdNew, count, active, kyc, oid, type, verify, title, html, tags, channels, channels64, floats, keys, features, likes, providers, stats, price, meta, timeout, value, rawbytes) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23) returning newid"
	err = conn.QueryRow(c, q, v.sqlTuple()...).Scan(&id)
	return
}

// has value in db
func (a *NewsSQL) Has(field NewsIndexType, v any) (has bool, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("select exists (select id from news where %s = $1 limit 1)", field.SQLName())
	err = conn.QueryRow(c, q, v).Scan(&has)
	return
}

// Create table
func (a *NewsSQL) CreateTable() (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := `create table if not exists news (
	newid                                        bigserial primary key,
	createdNew                                   bigint default extract(epoch from now()),
	count                                        bigint,
	active                                       boolean,
	kyc                                          newtype,
	oid                                          bigint,
	type                                         bigint,
	verify                                       boolean,
	title                                        text,
	html                                         bytea,
	tags                                         jsonb default '[]'::jsonb,
	channels                                     jsonb default '[]'::jsonb,
	channels64                                   jsonb default '[]'::jsonb,
	floats                                       double precision,
	keys                                         jsonb default '{}'::jsonb,
	features                                     jsonb default '{}'::jsonb,
	likes                                        jsonb default '{}'::jsonb,
	providers                                    jsonb default '{}'::jsonb,
	stats                                        jsonb default '{}'::jsonb,
	price                                        jsonb default '{}'::jsonb,
	meta                                         jsonb default '{}'::jsonb,
	timeout                                      bigint,
	value                                        jsonb default '{}'::jsonb,
	rawbytes                                     bytea,
	unique(tags,channels),
	unique(html,channels,oid),
	primary key (floats,keys,features)
)
`
	_, err = conn.Exec(c, q)
	return
}

// parse sql query
func (a *NewsSQL) Fields() (res []NewsIndexType) {
	return []NewsIndexType{IndexAccountID, IndexCreated, IndexCount, IndexActive, IndexKYC, IndexBID, IndexOID, IndexType, IndexVerify, IndexTitle, IndexHtml, IndexTags, IndexChannels, IndexChannels64, IndexFloats, IndexKeys, IndexFeatures, IndexLikes, IndexProviders, IndexStats, IndexPrice, IndexMeta, IndexTimeout, IndexValue, IndexRaw}
}
