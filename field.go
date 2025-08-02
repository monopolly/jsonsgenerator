package main

type Field struct {
	// ID        string //оригинальное имя в структуре: types
	Name      string //имя в json структуре: type
	Type      string //int
	Comment   string //nice touch
	FieldName string //fieldType для индекса полей

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
		Raw   bool   //только если []byte json внутри
		Named string //JsonFieldNewsID генерить наименованные поля если [] но нужно передавать json
	}

	MessagePack struct {
		Name string
	}

	Go struct {
		Index string //IndexType для порядкового индекса IndexType = 0 и тд
		Name  string // Created, ID
		Type  string // []*List

		Nofunc    bool
		Func      string          //NewsAdd
		UpperCase bool            //sid > SID
		Must      bool            //special function for validation field valuer
		Tags      map[string]bool //#private #visible

		Title string // type title: dims > Dimension
		Desc  string // desc{Nice}

		UserLists    map[string][]string // list{title="Required", art, cars, any}
		UserRequired []string            // req{art, cars, any}

	}

	SQL struct {
		Name       string //sql name
		Replace    string //заменяет sql type
		Append     string //добавляет sql
		NoInsert   bool   //если не нужно добавлять в sql insert
		Type       string //new sql type jsonb
		JsonbArray bool   //new sql type

		Index struct {
			Simple       bool //создать индекс по полю
			Concurrently bool //конкурентный индекс
			// Jsonb        bool     //создать jsonb индекс
			Search []string //создать текстовый индекс чтобы искать текст
			Group  []string
		}
		Default struct {
			Value    string
			Simple   bool
			Unixtime bool
		}
		Inc        bool     //"bigserial primary key" noinsert
		Skip       bool     //если не нужно добавлять в sql таблицу
		Unique     []string //unique fields
		Primary    bool     //primary key fields
		AlterTable string   //version || 1
		Rename     string   //rename="oldname"
	}

	Clickhouse struct {
		Name       string //sql name
		Engine     string //sql engine ENGINE = MergeTree()
		Replace    string //заменяет sql type
		Append     string //добавляет sql
		NoInsert   bool   //если не нужно добавлять в sql insert
		Type       string //new sql type https://clickhouse.com/docs/sql-reference/data-types/int-uint
		JsonbArray bool   //new sql type

		Index struct {
			Simple       bool //создать индекс по полю
			Concurrently bool //конкурентный индекс
			// Jsonb        bool     //создать jsonb индекс
			Search []string //создать текстовый индекс чтобы искать текст
			Group  []string
		}
		Default struct {
			Value    string
			Simple   bool
			Unixtime bool
		}
		Inc        bool     //"bigserial primary key" noinsert
		Skip       bool     //если не нужно добавлять в sql таблицу
		Unique     []string //unique fields
		Primary    bool     //primary key fields PRIMARY KEY(category_id, sale_date)
		AlterTable string   //version || 1
		Rename     string   //rename="oldname"
	}
}
