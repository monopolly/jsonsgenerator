package main

type Field struct {
	// ID        string // original struct name: types
	Name      string // json struct name: type
	Type      string //int
	Comment   string //nice touch
	FieldName string // fieldType for field index

	Swift struct {
		Must  bool     //field required {{swift/must}}
		Skip  bool     //not use field {{swift/not}}
		Type  string   //new swift type {{swift/type/string}}
		Files []string //several swift file names
		Enum  struct {
			Skip bool //skip field in swift enum
		}
	}

	Json struct {
		Skip  bool
		Inc   bool
		Bool  bool
		Time  bool
		Name  string //id
		Raw   bool   // only when []byte contains json
		Named string // JsonFieldNewsID generates named fields when [] is used, but json must be passed
	}

	MessagePack struct {
		Name string
	}

	Proto struct {
		Name string
		Type string
		Skip bool
	}

	Go struct {
		Index string // IndexType for ordinal index: IndexType = 0, etc.
		Name  string // Created, ID
		Type  string // []*List

		Nofunc        bool
		Func          string          //NewsAdd
		UpperCase     bool            //sid > SID
		Must          bool            //special function for validation field valuer
		Tags          map[string]bool //#private #visible
		CustomStructs map[string]bool //@public > AccountPublic{}

		Title string // type title: dims > Dimension
		Desc  string // desc{Nice}

		UserLists    map[string][]string // list{title="Required", art, cars, any}
		UserRequired []string            // req{art, cars, any}

		Bits  int //1,2,8...
		Align int //1,2,8...

	}

	SQL struct {
		Name       string //sql name
		Replace    string // replaces sql type
		Append     string // appends sql
		NoInsert   bool   // skip in sql insert
		Type       string //new sql type jsonb
		JsonbArray bool   //new sql type
		Get        bool   //need golang sql methods GetName(v type)
		Keys       string //need golang sql methods GetNameSecond(v1 type, v2 type...)

		Index struct {
			Simple       bool // create index for field
			Concurrently bool // concurrent index
			// Jsonb        bool     // create jsonb index
			Search []string // create text index for search
			Group  []string
		}
		Default struct {
			Value    string
			Simple   bool
			Unixtime bool
		}
		Inc        bool     //"bigserial primary key" noinsert
		Skip       bool     // skip in sql table
		Unique     []string //unique fields
		Primary    bool     //primary key fields
		AlterTable string   //version || 1
		Rename     string   //rename="oldname"
	}

	Clickhouse struct {
		Name    string //sql name id, name...
		Type    string //new sql type https://clickhouse.com/docs/sql-reference/data-types/int-uint
		Skip    bool   // skip in sql table
		Primary bool   //primary key fields PRIMARY KEY(category_id, sale_date)
		Unix    bool   //primary key fields PRIMARY KEY(category_id, sale_date)
		Low     bool   //low cardinality
		IP      bool   //IPv4/IPv6 field
	}
}
