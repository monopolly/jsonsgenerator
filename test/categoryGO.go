package news

import (
	"fmt"

	"github.com/monopolly/cast"
	"github.com/monopolly/jsons"
)

/*
	Model Generator
	Help to auto-generated correct models for golang structures.
	Martin Prestone (c) 2024
	github.com/monopolly
	Ex: //! [] go=News js=NewsJson sql=news noinit up...

	STRUCT:
	go=News        Set golang struct names. Ex: go=News1 > type News1 struct{}
	js=NewsJson    Change json struct names. Ex: js=NewsJson
	demo           Generate demo json file with default values
	lock!          Lock model for new jsons generation. Can't change model.
	noinit         No init function for struct
	!omit          No omit tag for json
	enum           Generate swift enum
	noprefix       Generate simple index IndexID instead IndexNewsID
	sql=news       Set sql table name. Ex: sql=accounts
	ts=news        Set typescript struct names. Ex: ts=NewsJson
	swift          Generate swift model

	GO:
	nofunc         Do not create any jsons functions for fiels
	up             Make uppercase for functions
	must           Create one validation function for all must fields
	type=""        Replace golang type. Ex: type="[]*News"
	name=""        Replace golang struct name. Ex: name="NewsList"
	#              Add lists for fields. Ex: id int //#readonly #must...

	SQL:
	skip           Skip field for sql queries
	type=""        Rewrite sql type for field. Ex: type="jsonb"
	add=""         Append sql for field. Ex: replace="primary key"
	unique=""      Add unique fields constrains by group (you need set group name, then generator join fields). Ex: unique="group1" unique="group2"
	noinsert       Do use field for insert function
	idx            Add index fields by group name. Ex: idx="nameIndex" idx="credsIndex"
	defaults       Add default value based on field type
	replace=""     Rewrite sql for field. Ex: replace="bigint primary key"
	index          Create simple default index or gin for jsonb
	search         Add tsvector index by group. Ex: search="tsv": tsv tsvector GENERATED ALWAYS AS (to_tsvector('simple', title || ' ' || brand)). For search: SELECT brand, title FROM assets WHERE search @@ to_tsquery('english', 'f8');
	unix           Add default value: extract(epoch from now())

	SWIFT:
	must           Swift field with required values
	skip           Ignore field for Swift
	type=""        Replace type for Swift model. Ex: type="string"
	file=""        Create another swift file for this field Swift model. file="f1", file="f2"

	JSON:
	skip           Skip json field
	inc            Add inc jsons function for numbers fields
	bool           Add set jsons function for bool fields
	time           Create convert jsons function for unixtime fields
	raw            Set Raw json function inside field
	name=""        Replace json field name. Ex: name="sid"
*/

// порядковые номера полей
type CategoryIndexType int

const (
	IndexCategoryArt = CategoryIndexType(iota)
	IndexCategoryCars
	IndexCategoryRealestate
	IndexCategoryWatch
	IndexCategoryJewelry
)

const (
	FieldCategoryArt        = "art"        // int swift{skip}
	FieldCategoryCars       = "cars"       // int
	FieldCategoryRealestate = "realestate" // int
	FieldCategoryWatch      = "watch"      // int
	FieldCategoryJewelry    = "jewelry"    // int
)

func CategoryIndexes() []CategoryIndexType {
	return []CategoryIndexType{IndexCategoryArt, IndexCategoryCars, IndexCategoryRealestate, IndexCategoryWatch, IndexCategoryJewelry}

}

// public struct with json tags
type Category struct {
	Art        int `json:"art,omitempty"`        // swift{skip}
	Cars       int `json:"cars,omitempty"`       //
	Realestate int `json:"realestate,omitempty"` //
	Watch      int `json:"watch,omitempty"`      //
	Jewelry    int `json:"jewelry,omitempty"`    //
}

// Unmarshal json to struct
func JsonToStructCategory(r []byte) (a Category, err error) {
	err = jsons.Unmarshal(r, &a)
	return
}

// Parse []any to ID struct
func ParseCategoryToStruct(r []any) (a *Category) {
	a = new(Category)

	for pos, x := range r {
		switch CategoryIndexType(pos) {
		case IndexCategoryArt:
			a.Art = cast.Int(x) //int
		case IndexCategoryCars:
			a.Cars = cast.Int(x) //int
		case IndexCategoryRealestate:
			a.Realestate = cast.Int(x) //int
		case IndexCategoryWatch:
			a.Watch = cast.Int(x) //int
		case IndexCategoryJewelry:
			a.Jewelry = cast.Int(x) //int
		}
	}
	return
}

// Tuple create an array from struct
func (a *Category) Tuple() (r []any) {
	return []any{a.Art, a.Cars, a.Realestate, a.Watch, a.Jewelry}
}

// update struct with function
func (a *Category) Update(k string, x any) {
	switch k {
	case "art":
		a.Art = cast.Int(x) //int
	case "cars":
		a.Cars = cast.Int(x) //int
	case "realestate":
		a.Realestate = cast.Int(x) //int
	case "watch":
		a.Watch = cast.Int(x) //int
	case "jewelry":
		a.Jewelry = cast.Int(x) //int
	}
	return
}

// get struct value with function
func (a *Category) Get(k string) (v any) {
	switch k {
	case "art":
		return a.Art //int
	case "cars":
		return a.Cars //int
	case "realestate":
		return a.Realestate //int
	case "watch":
		return a.Watch //int
	case "jewelry":
		return a.Jewelry //int
	}
	return
}

// get any struct value as string
func (a *Category) String(k string) (v string) {
	switch k {
	case "art":
		return fmt.Sprint(a.Art) //int
	case "cars":
		return fmt.Sprint(a.Cars) //int
	case "realestate":
		return fmt.Sprint(a.Realestate) //int
	case "watch":
		return fmt.Sprint(a.Watch) //int
	case "jewelry":
		return fmt.Sprint(a.Jewelry) //int
	}
	return
}

// Struct to json
func (a *Category) ToJson() (r []byte) {
	js := jsons.Create().
		Add(FieldCategoryArt, a.Art).
		Add(FieldCategoryCars, a.Cars).
		Add(FieldCategoryRealestate, a.Realestate).
		Add(FieldCategoryWatch, a.Watch).
		Add(FieldCategoryJewelry, a.Jewelry)
	return js.Bytes()
}

// key index string
func (a CategoryIndexType) String() string {
	switch a {
	case IndexCategoryArt:
		return "art"
	case IndexCategoryCars:
		return "cars"
	case IndexCategoryRealestate:
		return "realestate"
	case IndexCategoryWatch:
		return "watch"
	case IndexCategoryJewelry:
		return "jewelry"
	default:
		return ""
	}
}

// key index type
func (a CategoryIndexType) Type() string {
	switch a {
	case IndexCategoryArt:
		return "int"
	case IndexCategoryCars:
		return "int"
	case IndexCategoryRealestate:
		return "int"
	case IndexCategoryWatch:
		return "int"
	case IndexCategoryJewelry:
		return "int"
	default:
		return ""
	}
}

// key index title
func (a CategoryIndexType) Title() string {
	switch a {
	case IndexCategoryArt:
		return "Art"
	case IndexCategoryCars:
		return "Cars"
	case IndexCategoryRealestate:
		return "Realestate"
	case IndexCategoryWatch:
		return "Watch"
	case IndexCategoryJewelry:
		return "Jewelry"
	default:
		return ""
	}
}

// struct key to index
func CategoryKeyIndex(key string) CategoryIndexType {
	switch key {
	case "art":
		return IndexCategoryArt
	case "cars":
		return IndexCategoryCars
	case "realestate":
		return IndexCategoryRealestate
	case "watch":
		return IndexCategoryWatch
	case "jewelry":
		return IndexCategoryJewelry
	default:
		return 0
	}
}

// valid struct key check
func CategoryValidKey(key string) bool {
	switch key {
	case "art", "cars", "realestate", "watch", "jewelry":
		return true
	default:
		return false
	}
}

// struct to map
func (a *Category) Map() map[string]any {
	return map[string]any{
		"art":        a.Art,
		"cars":       a.Cars,
		"realestate": a.Realestate,
		"watch":      a.Watch,
		"jewelry":    a.Jewelry,
	}
}

// struct to map
func (a *Category) Iterate(f func(k CategoryIndexType, v any)) {
	for _, x := range CategoryIndexes() {
		f(x, a.Get(x.String()))
	}
}
