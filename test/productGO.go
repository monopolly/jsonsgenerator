package news

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	jsoniter "github.com/json-iterator/go"
	"github.com/monopolly/cast"
	"github.com/monopolly/errors"
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
	enum                          Generate swift enum
	noprefix                      Generate simple index IndexID instead IndexNewsID
	debug                         Debug mode
	!omit                         No omit tag for json
	js=NewsJson                   Change json struct names. Ex: js=NewsJson
	swift                         Generate swift model
	demo                          Generate demo json file with default values
	sql=news                      Set sql table name. Ex: sql=accounts
	lock                          Lock model for generation. Can't change model.
	noinit                        No New() init function for struct
	gotiny                        Create gotiny marshal/unmarshal
	msgp                          Create message pack marshal/unmarshal
	go=News                       Set golang struct names. Ex: go=News1 > type News1 struct{}
	ts=news                       Set typescript struct names. Ex: ts=NewsJson

	Field:
	title{}                       Add custom title for index. title{Nice & Sweet}
	desc{}                        Add custom desc for index. desc{This field for success}

	GO:
	desc=""                       Add custom desc for index
	nofunc                        Do not create any jsons functions for fiels
	must                          Create one validation function for all must fields
	req{}                         Add user required fields
	title{}                       Add field user title title{Nice}
	desc{}                        Add field user title desc{Use it nice}
	#                             Add lists for fields. Ex: id int //#readonly #must...
	up                            Make uppercase for functions
	type=""                       Replace golang type. Ex: type="[]*News"
	name=""                       Replace golang struct name. Ex: name="NewsList"
	title=""                      Add custom title for index

	SQL:
	skip                          Skip field for sql queries
	replace="bigint primary key"  Rewrite sql for field. Ex: replace="bigint primary key"
	add="primary key"             Append sql for field. Ex: replace="primary key"
	unique="groupname"            Add unique fields constrains by group (you need set group name, then generator join fields). Ex: unique="group1" unique="group2"
	index                         Create simple default index or gin for jsonb
	idx                           Add index fields by group name. Ex: idx="nameIndex" idx="credsIndex"
	defaults                      Add default value based on field type
	type="jsonb"                  Rewrite sql type for field. Ex: type="jsonb"
	ver="v4"                      Create an alter table record in SQL file. Add new column. With new version in comment. Ex: ver="2"
	renames="oldname"             Create an alter table record in SQL file. Rename table column.
	noinsert                      Do use field for insert function
	altertable                    Add field line Alter table to SQL file with current time comment
	search                        Add tsvector index by group. Ex: search="tsv": tsv tsvector GENERATED ALWAYS AS (to_tsvector('simple', title || ' ' || brand)). For search: SELECT brand, title FROM assets WHERE search @@ to_tsquery('english', 'f8');
	unix                          Add default value: extract(epoch from now())

	SWIFT:
	type=""                       Replace type for Swift model. Ex: type="string"
	file=""                       Create another swift file for this field Swift model. file="f1", file="f2"
	must                          Swift field with required values
	skip                          Ignore field for Swift

	JSON:
	skip                          Skip json field
	inc                           Add inc jsons function for numbers fields
	bool                          Add set jsons function for bool fields
	time                          Create convert jsons function for unixtime fields
	raw                           Set Raw json function inside field
	name=""                       Replace json field name. Ex: name="sid"
*/

// field type
type ProductIndexType int

// int index
const (
	IndexProductID = ProductIndexType(iota)
	IndexProductCreated
	IndexProductUpdated
	IndexProductActive
	IndexProductRef
	IndexProductSku
	IndexProductCategory
	IndexProductBrand
	IndexProductModel
	IndexProductTitle
	IndexProductLine
	IndexProductAbout
	IndexProductImage
	IndexProductOriginal
	IndexProductCost
	IndexProductDiscount
	IndexProductSample
	IndexProductWidth
	IndexProductLength
	IndexProductSize
	IndexProductThickness
	IndexProductCount
	IndexProductPrice
	IndexProductBoxwidth
	IndexProductBoxheight
	IndexProductBoxsize
	IndexProductBoxweight
	IndexProductEdge
	IndexProductFlooring
	IndexProductInstallation
	IndexProductConstruction
	IndexProductGloss
	IndexProductWaste
	IndexProductWear
	IndexProductFeatures
	IndexProductLayers
	IndexProductMaterials
	IndexProductMeta
)

// string index
const (
	FieldProductID           = "id"           // int sql{inc}
	FieldProductCreated      = "created"      // int64 sql{unix}
	FieldProductUpdated      = "updated"      // int64 sql{unix}
	FieldProductActive       = "active"       // bool sql{defaults}
	FieldProductRef          = "ref"          // string
	FieldProductSku          = "sku"          // string
	FieldProductCategory     = "category"     // string vinyl
	FieldProductBrand        = "brand"        // string sql{altertable}
	FieldProductModel        = "model"        // string sql{ver="2"}
	FieldProductTitle        = "title"        // string James
	FieldProductLine         = "line"         // string A good guy
	FieldProductAbout        = "about"        // string md sql{renames="oldname"}
	FieldProductImage        = "image"        // string sql{ver="October"}
	FieldProductOriginal     = "original"     // string link to product
	FieldProductCost         = "cost"         // float64 1.15,
	FieldProductDiscount     = "discount"     // float64 5 percent
	FieldProductSample       = "sample"       // float64 1.5 cost usd
	FieldProductWidth        = "width"        // float64 7.1 in
	FieldProductLength       = "length"       // float64 48 in
	FieldProductSize         = "size"         // float64 540 sqft
	FieldProductThickness    = "thickness"    // float64 mm
	FieldProductCount        = "count"        // int 30 count
	FieldProductPrice        = "price"        // float64 98.4 usd
	FieldProductBoxwidth     = "boxwidth"     // float64 200 in
	FieldProductBoxheight    = "boxheight"    // float64 400 in
	FieldProductBoxsize      = "boxsize"      // float64 1400 sqft
	FieldProductBoxweight    = "boxweight"    // float64 30 kg
	FieldProductEdge         = "edge"         // string square, bevel
	FieldProductFlooring     = "flooring"     // string tile, plank
	FieldProductInstallation = "installation" // string glue, plank
	FieldProductConstruction = "construction" // string attached, attached pad, floating
	FieldProductGloss        = "gloss"        // string low,
	FieldProductWaste        = "waste"        // int Recommended Waste Factor percents
	FieldProductWear         = "wear"         // int wear layers, mil
	FieldProductFeatures     = "features"     // map[string]bool waterproof
	FieldProductLayers       = "layers"       // map[string]bool waterproof, vinyl, soundproof
	FieldProductMaterials    = "materials"    // map[string]bool vinyl, plastic
	FieldProductMeta         = "meta"         // map[string]any
)

// index func
func ProductIndexes() []ProductIndexType {
	return []ProductIndexType{IndexProductID, IndexProductCreated, IndexProductUpdated, IndexProductActive, IndexProductRef, IndexProductSku, IndexProductCategory, IndexProductBrand, IndexProductModel, IndexProductTitle, IndexProductLine, IndexProductAbout, IndexProductImage, IndexProductOriginal, IndexProductCost, IndexProductDiscount, IndexProductSample, IndexProductWidth, IndexProductLength, IndexProductSize, IndexProductThickness, IndexProductCount, IndexProductPrice, IndexProductBoxwidth, IndexProductBoxheight, IndexProductBoxsize, IndexProductBoxweight, IndexProductEdge, IndexProductFlooring, IndexProductInstallation, IndexProductConstruction, IndexProductGloss, IndexProductWaste, IndexProductWear, IndexProductFeatures, IndexProductLayers, IndexProductMaterials, IndexProductMeta}

}

type Product struct {
	ID           int             `json:"id,omitempty" msg:"id,omitempty"`                     // sql{inc}
	Created      int64           `json:"created,omitempty" msg:"created,omitempty"`           // sql{unix}
	Updated      int64           `json:"updated,omitempty" msg:"updated,omitempty"`           // sql{unix}
	Active       bool            `json:"active,omitempty" msg:"active,omitempty"`             // sql{defaults}
	Ref          string          `json:"ref,omitempty" msg:"ref,omitempty"`                   //
	Sku          string          `json:"sku,omitempty" msg:"sku,omitempty"`                   //
	Category     string          `json:"category,omitempty" msg:"category,omitempty"`         // vinyl
	Brand        string          `json:"brand,omitempty" msg:"brand,omitempty"`               // sql{altertable}
	Model        string          `json:"model,omitempty" msg:"model,omitempty"`               // sql{ver="2"}
	Title        string          `json:"title,omitempty" msg:"title,omitempty"`               // James
	Line         string          `json:"line,omitempty" msg:"line,omitempty"`                 // A good guy
	About        string          `json:"about,omitempty" msg:"about,omitempty"`               // md sql{renames="oldname"}
	Image        string          `json:"image,omitempty" msg:"image,omitempty"`               // sql{ver="October"}
	Original     string          `json:"original,omitempty" msg:"original,omitempty"`         // link to product
	Cost         float64         `json:"cost,omitempty" msg:"cost,omitempty"`                 // 1.15,
	Discount     float64         `json:"discount,omitempty" msg:"discount,omitempty"`         // 5 percent
	Sample       float64         `json:"sample,omitempty" msg:"sample,omitempty"`             // 1.5 cost usd
	Width        float64         `json:"width,omitempty" msg:"width,omitempty"`               // 7.1 in
	Length       float64         `json:"length,omitempty" msg:"length,omitempty"`             // 48 in
	Size         float64         `json:"size,omitempty" msg:"size,omitempty"`                 // 540 sqft
	Thickness    float64         `json:"thickness,omitempty" msg:"thickness,omitempty"`       // mm
	Count        int             `json:"count,omitempty" msg:"count,omitempty"`               // 30 count
	Price        float64         `json:"price,omitempty" msg:"price,omitempty"`               // 98.4 usd
	Boxwidth     float64         `json:"boxwidth,omitempty" msg:"boxwidth,omitempty"`         // 200 in
	Boxheight    float64         `json:"boxheight,omitempty" msg:"boxheight,omitempty"`       // 400 in
	Boxsize      float64         `json:"boxsize,omitempty" msg:"boxsize,omitempty"`           // 1400 sqft
	Boxweight    float64         `json:"boxweight,omitempty" msg:"boxweight,omitempty"`       // 30 kg
	Edge         string          `json:"edge,omitempty" msg:"edge,omitempty"`                 // square, bevel
	Flooring     string          `json:"flooring,omitempty" msg:"flooring,omitempty"`         // tile, plank
	Installation string          `json:"installation,omitempty" msg:"installation,omitempty"` // glue, plank
	Construction string          `json:"construction,omitempty" msg:"construction,omitempty"` // attached, attached pad, floating
	Gloss        string          `json:"gloss,omitempty" msg:"gloss,omitempty"`               // low,
	Waste        int             `json:"waste,omitempty" msg:"waste,omitempty"`               // Recommended Waste Factor percents
	Wear         int             `json:"wear,omitempty" msg:"wear,omitempty"`                 // wear layers, mil
	Features     map[string]bool `json:"features,omitempty" msg:"features,omitempty"`         // waterproof
	Layers       map[string]bool `json:"layers,omitempty" msg:"layers,omitempty"`             // waterproof, vinyl, soundproof
	Materials    map[string]bool `json:"materials,omitempty" msg:"materials,omitempty"`       // vinyl, plastic
	Meta         map[string]any  `json:"meta,omitempty" msg:"meta,omitempty"`                 //
}

// Parse []any to ID struct
func ParseProductToStruct(r []any) (a *Product) {
	a = new(Product)

	for pos, x := range r {
		switch ProductIndexType(pos) {
		case IndexProductID:
			a.ID = cast.Int(x) //int
		case IndexProductCreated:
			a.Created = cast.Int64(x) //int64
		case IndexProductUpdated:
			a.Updated = cast.Int64(x) //int64
		case IndexProductActive:
			a.Active = cast.Bool(x) //bool
		case IndexProductRef:
			a.Ref = cast.String(x) //string
		case IndexProductSku:
			a.Sku = cast.String(x) //string
		case IndexProductCategory:
			a.Category = cast.String(x) //string
		case IndexProductBrand:
			a.Brand = cast.String(x) //string
		case IndexProductModel:
			a.Model = cast.String(x) //string
		case IndexProductTitle:
			a.Title = cast.String(x) //string
		case IndexProductLine:
			a.Line = cast.String(x) //string
		case IndexProductAbout:
			a.About = cast.String(x) //string
		case IndexProductImage:
			a.Image = cast.String(x) //string
		case IndexProductOriginal:
			a.Original = cast.String(x) //string
		case IndexProductCost:
			a.Cost = cast.Float(x) //float64
		case IndexProductDiscount:
			a.Discount = cast.Float(x) //float64
		case IndexProductSample:
			a.Sample = cast.Float(x) //float64
		case IndexProductWidth:
			a.Width = cast.Float(x) //float64
		case IndexProductLength:
			a.Length = cast.Float(x) //float64
		case IndexProductSize:
			a.Size = cast.Float(x) //float64
		case IndexProductThickness:
			a.Thickness = cast.Float(x) //float64
		case IndexProductCount:
			a.Count = cast.Int(x) //int
		case IndexProductPrice:
			a.Price = cast.Float(x) //float64
		case IndexProductBoxwidth:
			a.Boxwidth = cast.Float(x) //float64
		case IndexProductBoxheight:
			a.Boxheight = cast.Float(x) //float64
		case IndexProductBoxsize:
			a.Boxsize = cast.Float(x) //float64
		case IndexProductBoxweight:
			a.Boxweight = cast.Float(x) //float64
		case IndexProductEdge:
			a.Edge = cast.String(x) //string
		case IndexProductFlooring:
			a.Flooring = cast.String(x) //string
		case IndexProductInstallation:
			a.Installation = cast.String(x) //string
		case IndexProductConstruction:
			a.Construction = cast.String(x) //string
		case IndexProductGloss:
			a.Gloss = cast.String(x) //string
		case IndexProductWaste:
			a.Waste = cast.Int(x) //int
		case IndexProductWear:
			a.Wear = cast.Int(x) //int
		case IndexProductFeatures:
			a.Features = cast.StringMapBool(x) //map[string]bool
		case IndexProductLayers:
			a.Layers = cast.StringMapBool(x) //map[string]bool
		case IndexProductMaterials:
			a.Materials = cast.StringMapBool(x) //map[string]bool
		case IndexProductMeta:
			a.Meta = cast.StringMap(x) //map[string]any
		}
	}
	return
}

// Tuple create an array from struct
func (a *Product) Tuple() (r []any) {
	return []any{a.ID, a.Created, a.Updated, a.Active, a.Ref, a.Sku, a.Category, a.Brand, a.Model, a.Title, a.Line, a.About, a.Image, a.Original, a.Cost, a.Discount, a.Sample, a.Width, a.Length, a.Size, a.Thickness, a.Count, a.Price, a.Boxwidth, a.Boxheight, a.Boxsize, a.Boxweight, a.Edge, a.Flooring, a.Installation, a.Construction, a.Gloss, a.Waste, a.Wear, a.Features, a.Layers, a.Materials, a.Meta}
}

// Tuple create an array from struct
func (a *Product) sqlTuple() (r []any) {
	return []any{a.Created, a.Updated, a.Active, a.Ref, a.Sku, a.Category, a.Brand, a.Model, a.Title, a.Line, a.About, a.Image, a.Original, a.Cost, a.Discount, a.Sample, a.Width, a.Length, a.Size, a.Thickness, a.Count, a.Price, a.Boxwidth, a.Boxheight, a.Boxsize, a.Boxweight, a.Edge, a.Flooring, a.Installation, a.Construction, a.Gloss, a.Waste, a.Wear, a.Features, a.Layers, a.Materials, a.Meta}
}

// update struct with function
func (a *Product) Update(k string, x any) {
	switch k {
	case "id":
		a.ID = cast.Int(x) //int
	case "created":
		a.Created = cast.Int64(x) //int64
	case "updated":
		a.Updated = cast.Int64(x) //int64
	case "active":
		a.Active = cast.Bool(x) //bool
	case "ref":
		a.Ref = cast.String(x) //string
	case "sku":
		a.Sku = cast.String(x) //string
	case "category":
		a.Category = cast.String(x) //string
	case "brand":
		a.Brand = cast.String(x) //string
	case "model":
		a.Model = cast.String(x) //string
	case "title":
		a.Title = cast.String(x) //string
	case "line":
		a.Line = cast.String(x) //string
	case "about":
		a.About = cast.String(x) //string
	case "image":
		a.Image = cast.String(x) //string
	case "original":
		a.Original = cast.String(x) //string
	case "cost":
		a.Cost = cast.Float(x) //float64
	case "discount":
		a.Discount = cast.Float(x) //float64
	case "sample":
		a.Sample = cast.Float(x) //float64
	case "width":
		a.Width = cast.Float(x) //float64
	case "length":
		a.Length = cast.Float(x) //float64
	case "size":
		a.Size = cast.Float(x) //float64
	case "thickness":
		a.Thickness = cast.Float(x) //float64
	case "count":
		a.Count = cast.Int(x) //int
	case "price":
		a.Price = cast.Float(x) //float64
	case "boxwidth":
		a.Boxwidth = cast.Float(x) //float64
	case "boxheight":
		a.Boxheight = cast.Float(x) //float64
	case "boxsize":
		a.Boxsize = cast.Float(x) //float64
	case "boxweight":
		a.Boxweight = cast.Float(x) //float64
	case "edge":
		a.Edge = cast.String(x) //string
	case "flooring":
		a.Flooring = cast.String(x) //string
	case "installation":
		a.Installation = cast.String(x) //string
	case "construction":
		a.Construction = cast.String(x) //string
	case "gloss":
		a.Gloss = cast.String(x) //string
	case "waste":
		a.Waste = cast.Int(x) //int
	case "wear":
		a.Wear = cast.Int(x) //int
	case "features":
		a.Features = cast.StringMapBool(x) //map[string]bool
	case "layers":
		a.Layers = cast.StringMapBool(x) //map[string]bool
	case "materials":
		a.Materials = cast.StringMapBool(x) //map[string]bool
	case "meta":
		a.Meta = cast.StringMap(x) //map[string]any
	}
}

// get struct value with function
func (a *Product) Get(k string) (v any) {
	switch k {
	case "id":
		return a.ID //int
	case "created":
		return a.Created //int64
	case "updated":
		return a.Updated //int64
	case "active":
		return a.Active //bool
	case "ref":
		return a.Ref //string
	case "sku":
		return a.Sku //string
	case "category":
		return a.Category //string
	case "brand":
		return a.Brand //string
	case "model":
		return a.Model //string
	case "title":
		return a.Title //string
	case "line":
		return a.Line //string
	case "about":
		return a.About //string
	case "image":
		return a.Image //string
	case "original":
		return a.Original //string
	case "cost":
		return a.Cost //float64
	case "discount":
		return a.Discount //float64
	case "sample":
		return a.Sample //float64
	case "width":
		return a.Width //float64
	case "length":
		return a.Length //float64
	case "size":
		return a.Size //float64
	case "thickness":
		return a.Thickness //float64
	case "count":
		return a.Count //int
	case "price":
		return a.Price //float64
	case "boxwidth":
		return a.Boxwidth //float64
	case "boxheight":
		return a.Boxheight //float64
	case "boxsize":
		return a.Boxsize //float64
	case "boxweight":
		return a.Boxweight //float64
	case "edge":
		return a.Edge //string
	case "flooring":
		return a.Flooring //string
	case "installation":
		return a.Installation //string
	case "construction":
		return a.Construction //string
	case "gloss":
		return a.Gloss //string
	case "waste":
		return a.Waste //int
	case "wear":
		return a.Wear //int
	case "features":
		return a.Features //map[string]bool
	case "layers":
		return a.Layers //map[string]bool
	case "materials":
		return a.Materials //map[string]bool
	case "meta":
		return a.Meta //map[string]any
	}
	return
}

// get any struct value as string
func (a *Product) String(k string) (v string) {
	switch k {
	case "id":
		return fmt.Sprint(a.ID) //int
	case "created":
		return fmt.Sprint(a.Created) //int64
	case "updated":
		return fmt.Sprint(a.Updated) //int64
	case "active":
		return fmt.Sprint(a.Active) //bool
	case "ref":
		return fmt.Sprint(a.Ref) //string
	case "sku":
		return fmt.Sprint(a.Sku) //string
	case "category":
		return fmt.Sprint(a.Category) //string
	case "brand":
		return fmt.Sprint(a.Brand) //string
	case "model":
		return fmt.Sprint(a.Model) //string
	case "title":
		return fmt.Sprint(a.Title) //string
	case "line":
		return fmt.Sprint(a.Line) //string
	case "about":
		return fmt.Sprint(a.About) //string
	case "image":
		return fmt.Sprint(a.Image) //string
	case "original":
		return fmt.Sprint(a.Original) //string
	case "cost":
		return fmt.Sprint(a.Cost) //float64
	case "discount":
		return fmt.Sprint(a.Discount) //float64
	case "sample":
		return fmt.Sprint(a.Sample) //float64
	case "width":
		return fmt.Sprint(a.Width) //float64
	case "length":
		return fmt.Sprint(a.Length) //float64
	case "size":
		return fmt.Sprint(a.Size) //float64
	case "thickness":
		return fmt.Sprint(a.Thickness) //float64
	case "count":
		return fmt.Sprint(a.Count) //int
	case "price":
		return fmt.Sprint(a.Price) //float64
	case "boxwidth":
		return fmt.Sprint(a.Boxwidth) //float64
	case "boxheight":
		return fmt.Sprint(a.Boxheight) //float64
	case "boxsize":
		return fmt.Sprint(a.Boxsize) //float64
	case "boxweight":
		return fmt.Sprint(a.Boxweight) //float64
	case "edge":
		return fmt.Sprint(a.Edge) //string
	case "flooring":
		return fmt.Sprint(a.Flooring) //string
	case "installation":
		return fmt.Sprint(a.Installation) //string
	case "construction":
		return fmt.Sprint(a.Construction) //string
	case "gloss":
		return fmt.Sprint(a.Gloss) //string
	case "waste":
		return fmt.Sprint(a.Waste) //int
	case "wear":
		return fmt.Sprint(a.Wear) //int
	case "features":
		return fmt.Sprint(a.Features) //map[string]bool
	case "layers":
		return fmt.Sprint(a.Layers) //map[string]bool
	case "materials":
		return fmt.Sprint(a.Materials) //map[string]bool
	case "meta":
		return fmt.Sprint(a.Meta) //map[string]any
	}
	return
}

// Struct to json
func (a *Product) ToJson() (r []byte) {
	js := jsons.Create().
		Add(FieldProductID, a.ID).
		Add(FieldProductCreated, a.Created).
		Add(FieldProductUpdated, a.Updated).
		Add(FieldProductActive, a.Active).
		Add(FieldProductRef, a.Ref).
		Add(FieldProductSku, a.Sku).
		Add(FieldProductCategory, a.Category).
		Add(FieldProductBrand, a.Brand).
		Add(FieldProductModel, a.Model).
		Add(FieldProductTitle, a.Title).
		Add(FieldProductLine, a.Line).
		Add(FieldProductAbout, a.About).
		Add(FieldProductImage, a.Image).
		Add(FieldProductOriginal, a.Original).
		Add(FieldProductCost, a.Cost).
		Add(FieldProductDiscount, a.Discount).
		Add(FieldProductSample, a.Sample).
		Add(FieldProductWidth, a.Width).
		Add(FieldProductLength, a.Length).
		Add(FieldProductSize, a.Size).
		Add(FieldProductThickness, a.Thickness).
		Add(FieldProductCount, a.Count).
		Add(FieldProductPrice, a.Price).
		Add(FieldProductBoxwidth, a.Boxwidth).
		Add(FieldProductBoxheight, a.Boxheight).
		Add(FieldProductBoxsize, a.Boxsize).
		Add(FieldProductBoxweight, a.Boxweight).
		Add(FieldProductEdge, a.Edge).
		Add(FieldProductFlooring, a.Flooring).
		Add(FieldProductInstallation, a.Installation).
		Add(FieldProductConstruction, a.Construction).
		Add(FieldProductGloss, a.Gloss).
		Add(FieldProductWaste, a.Waste).
		Add(FieldProductWear, a.Wear).
		Add(FieldProductFeatures, a.Features).
		Add(FieldProductLayers, a.Layers).
		Add(FieldProductMaterials, a.Materials).
		Add(FieldProductMeta, a.Meta)
	return js.Bytes()
}

// key index string
func (a ProductIndexType) String() string {
	switch a {
	case IndexProductID:
		return "id"
	case IndexProductCreated:
		return "created"
	case IndexProductUpdated:
		return "updated"
	case IndexProductActive:
		return "active"
	case IndexProductRef:
		return "ref"
	case IndexProductSku:
		return "sku"
	case IndexProductCategory:
		return "category"
	case IndexProductBrand:
		return "brand"
	case IndexProductModel:
		return "model"
	case IndexProductTitle:
		return "title"
	case IndexProductLine:
		return "line"
	case IndexProductAbout:
		return "about"
	case IndexProductImage:
		return "image"
	case IndexProductOriginal:
		return "original"
	case IndexProductCost:
		return "cost"
	case IndexProductDiscount:
		return "discount"
	case IndexProductSample:
		return "sample"
	case IndexProductWidth:
		return "width"
	case IndexProductLength:
		return "length"
	case IndexProductSize:
		return "size"
	case IndexProductThickness:
		return "thickness"
	case IndexProductCount:
		return "count"
	case IndexProductPrice:
		return "price"
	case IndexProductBoxwidth:
		return "boxwidth"
	case IndexProductBoxheight:
		return "boxheight"
	case IndexProductBoxsize:
		return "boxsize"
	case IndexProductBoxweight:
		return "boxweight"
	case IndexProductEdge:
		return "edge"
	case IndexProductFlooring:
		return "flooring"
	case IndexProductInstallation:
		return "installation"
	case IndexProductConstruction:
		return "construction"
	case IndexProductGloss:
		return "gloss"
	case IndexProductWaste:
		return "waste"
	case IndexProductWear:
		return "wear"
	case IndexProductFeatures:
		return "features"
	case IndexProductLayers:
		return "layers"
	case IndexProductMaterials:
		return "materials"
	case IndexProductMeta:
		return "meta"
	default:
		return ""
	}
}

// key index string
func (a ProductIndexType) SQLName() string {
	switch a {
	case IndexProductID:
		return "id"
	case IndexProductCreated:
		return "created"
	case IndexProductUpdated:
		return "updated"
	case IndexProductActive:
		return "active"
	case IndexProductRef:
		return "ref"
	case IndexProductSku:
		return "sku"
	case IndexProductCategory:
		return "category"
	case IndexProductBrand:
		return "brand"
	case IndexProductModel:
		return "model"
	case IndexProductTitle:
		return "title"
	case IndexProductLine:
		return "line"
	case IndexProductAbout:
		return "about"
	case IndexProductImage:
		return "image"
	case IndexProductOriginal:
		return "original"
	case IndexProductCost:
		return "cost"
	case IndexProductDiscount:
		return "discount"
	case IndexProductSample:
		return "sample"
	case IndexProductWidth:
		return "width"
	case IndexProductLength:
		return "length"
	case IndexProductSize:
		return "size"
	case IndexProductThickness:
		return "thickness"
	case IndexProductCount:
		return "count"
	case IndexProductPrice:
		return "price"
	case IndexProductBoxwidth:
		return "boxwidth"
	case IndexProductBoxheight:
		return "boxheight"
	case IndexProductBoxsize:
		return "boxsize"
	case IndexProductBoxweight:
		return "boxweight"
	case IndexProductEdge:
		return "edge"
	case IndexProductFlooring:
		return "flooring"
	case IndexProductInstallation:
		return "installation"
	case IndexProductConstruction:
		return "construction"
	case IndexProductGloss:
		return "gloss"
	case IndexProductWaste:
		return "waste"
	case IndexProductWear:
		return "wear"
	case IndexProductFeatures:
		return "features"
	case IndexProductLayers:
		return "layers"
	case IndexProductMaterials:
		return "materials"
	case IndexProductMeta:
		return "meta"
	default:
		return ""
	}
}

// key index type
func (a ProductIndexType) Type() string {
	switch a {
	case IndexProductID:
		return "int"
	case IndexProductCreated:
		return "int64"
	case IndexProductUpdated:
		return "int64"
	case IndexProductActive:
		return "bool"
	case IndexProductRef:
		return "string"
	case IndexProductSku:
		return "string"
	case IndexProductCategory:
		return "string"
	case IndexProductBrand:
		return "string"
	case IndexProductModel:
		return "string"
	case IndexProductTitle:
		return "string"
	case IndexProductLine:
		return "string"
	case IndexProductAbout:
		return "string"
	case IndexProductImage:
		return "string"
	case IndexProductOriginal:
		return "string"
	case IndexProductCost:
		return "float64"
	case IndexProductDiscount:
		return "float64"
	case IndexProductSample:
		return "float64"
	case IndexProductWidth:
		return "float64"
	case IndexProductLength:
		return "float64"
	case IndexProductSize:
		return "float64"
	case IndexProductThickness:
		return "float64"
	case IndexProductCount:
		return "int"
	case IndexProductPrice:
		return "float64"
	case IndexProductBoxwidth:
		return "float64"
	case IndexProductBoxheight:
		return "float64"
	case IndexProductBoxsize:
		return "float64"
	case IndexProductBoxweight:
		return "float64"
	case IndexProductEdge:
		return "string"
	case IndexProductFlooring:
		return "string"
	case IndexProductInstallation:
		return "string"
	case IndexProductConstruction:
		return "string"
	case IndexProductGloss:
		return "string"
	case IndexProductWaste:
		return "int"
	case IndexProductWear:
		return "int"
	case IndexProductFeatures:
		return "map[string]bool"
	case IndexProductLayers:
		return "map[string]bool"
	case IndexProductMaterials:
		return "map[string]bool"
	case IndexProductMeta:
		return "map[string]any"
	default:
		return ""
	}
}

// custom title
func (a ProductIndexType) Title() string {
	switch a {
	case IndexProductID:
		return "ID"
	case IndexProductCreated:
		return "Created"
	case IndexProductUpdated:
		return "Updated"
	case IndexProductActive:
		return "Active"
	case IndexProductRef:
		return "Ref"
	case IndexProductSku:
		return "Sku"
	case IndexProductCategory:
		return "Category"
	case IndexProductBrand:
		return "Brand"
	case IndexProductModel:
		return "Model"
	case IndexProductTitle:
		return "Title"
	case IndexProductLine:
		return "Line"
	case IndexProductAbout:
		return "About"
	case IndexProductImage:
		return "Image"
	case IndexProductOriginal:
		return "Original"
	case IndexProductCost:
		return "Cost"
	case IndexProductDiscount:
		return "Discount"
	case IndexProductSample:
		return "Sample"
	case IndexProductWidth:
		return "Width"
	case IndexProductLength:
		return "Length"
	case IndexProductSize:
		return "Size"
	case IndexProductThickness:
		return "Thickness"
	case IndexProductCount:
		return "Count"
	case IndexProductPrice:
		return "Price"
	case IndexProductBoxwidth:
		return "Boxwidth"
	case IndexProductBoxheight:
		return "Boxheight"
	case IndexProductBoxsize:
		return "Boxsize"
	case IndexProductBoxweight:
		return "Boxweight"
	case IndexProductEdge:
		return "Edge"
	case IndexProductFlooring:
		return "Flooring"
	case IndexProductInstallation:
		return "Installation"
	case IndexProductConstruction:
		return "Construction"
	case IndexProductGloss:
		return "Gloss"
	case IndexProductWaste:
		return "Waste"
	case IndexProductWear:
		return "Wear"
	case IndexProductFeatures:
		return "Features"
	case IndexProductLayers:
		return "Layers"
	case IndexProductMaterials:
		return "Materials"
	case IndexProductMeta:
		return "Meta"
	default:
		return ""
	}
}

// custom desc
func (a ProductIndexType) Desc() string {
	switch a {
	default:
		return ""
	}
}

// struct key to index
func ProductKeyIndex(key string) ProductIndexType {
	switch key {
	case "id":
		return IndexProductID
	case "created":
		return IndexProductCreated
	case "updated":
		return IndexProductUpdated
	case "active":
		return IndexProductActive
	case "ref":
		return IndexProductRef
	case "sku":
		return IndexProductSku
	case "category":
		return IndexProductCategory
	case "brand":
		return IndexProductBrand
	case "model":
		return IndexProductModel
	case "title":
		return IndexProductTitle
	case "line":
		return IndexProductLine
	case "about":
		return IndexProductAbout
	case "image":
		return IndexProductImage
	case "original":
		return IndexProductOriginal
	case "cost":
		return IndexProductCost
	case "discount":
		return IndexProductDiscount
	case "sample":
		return IndexProductSample
	case "width":
		return IndexProductWidth
	case "length":
		return IndexProductLength
	case "size":
		return IndexProductSize
	case "thickness":
		return IndexProductThickness
	case "count":
		return IndexProductCount
	case "price":
		return IndexProductPrice
	case "boxwidth":
		return IndexProductBoxwidth
	case "boxheight":
		return IndexProductBoxheight
	case "boxsize":
		return IndexProductBoxsize
	case "boxweight":
		return IndexProductBoxweight
	case "edge":
		return IndexProductEdge
	case "flooring":
		return IndexProductFlooring
	case "installation":
		return IndexProductInstallation
	case "construction":
		return IndexProductConstruction
	case "gloss":
		return IndexProductGloss
	case "waste":
		return IndexProductWaste
	case "wear":
		return IndexProductWear
	case "features":
		return IndexProductFeatures
	case "layers":
		return IndexProductLayers
	case "materials":
		return IndexProductMaterials
	case "meta":
		return IndexProductMeta
	default:
		return 0
	}
}

// valid struct key check
func ProductValidKey(key string) bool {
	switch key {
	case "id", "created", "updated", "active", "ref", "sku", "category", "brand", "model", "title", "line", "about", "image", "original", "cost", "discount", "sample", "width", "length", "size", "thickness", "count", "price", "boxwidth", "boxheight", "boxsize", "boxweight", "edge", "flooring", "installation", "construction", "gloss", "waste", "wear", "features", "layers", "materials", "meta":
		return true
	default:
		return false
	}
}

// struct to map
func (a *Product) Map() map[string]any {
	return map[string]any{
		"id":           a.ID,
		"created":      a.Created,
		"updated":      a.Updated,
		"active":       a.Active,
		"ref":          a.Ref,
		"sku":          a.Sku,
		"category":     a.Category,
		"brand":        a.Brand,
		"model":        a.Model,
		"title":        a.Title,
		"line":         a.Line,
		"about":        a.About,
		"image":        a.Image,
		"original":     a.Original,
		"cost":         a.Cost,
		"discount":     a.Discount,
		"sample":       a.Sample,
		"width":        a.Width,
		"length":       a.Length,
		"size":         a.Size,
		"thickness":    a.Thickness,
		"count":        a.Count,
		"price":        a.Price,
		"boxwidth":     a.Boxwidth,
		"boxheight":    a.Boxheight,
		"boxsize":      a.Boxsize,
		"boxweight":    a.Boxweight,
		"edge":         a.Edge,
		"flooring":     a.Flooring,
		"installation": a.Installation,
		"construction": a.Construction,
		"gloss":        a.Gloss,
		"waste":        a.Waste,
		"wear":         a.Wear,
		"features":     a.Features,
		"layers":       a.Layers,
		"materials":    a.Materials,
		"meta":         a.Meta,
	}
}

// struct to map
func (a *Product) Iterate(f func(k ProductIndexType, v any)) {
	for _, x := range ProductIndexes() {
		f(x, a.Get(x.String()))
	}
}

// gotiny marshal
func (a *Product) Gotiny() []byte {
	return gotiny.Marshal(&a.ID, &a.Created, &a.Updated, &a.Active, &a.Ref, &a.Sku, &a.Category, &a.Brand, &a.Model, &a.Title, &a.Line, &a.About, &a.Image, &a.Original, &a.Cost, &a.Discount, &a.Sample, &a.Width, &a.Length, &a.Size, &a.Thickness, &a.Count, &a.Price, &a.Boxwidth, &a.Boxheight, &a.Boxsize, &a.Boxweight, &a.Edge, &a.Flooring, &a.Installation, &a.Construction, &a.Gloss, &a.Waste, &a.Wear, &a.Features, &a.Layers, &a.Materials, &a.Meta)
}

// parse gotiny
func ParseProductGotiny(v []byte) (a Product) {
	gotiny.Unmarshal(v, &a.ID, &a.Created, &a.Updated, &a.Active, &a.Ref, &a.Sku, &a.Category, &a.Brand, &a.Model, &a.Title, &a.Line, &a.About, &a.Image, &a.Original, &a.Cost, &a.Discount, &a.Sample, &a.Width, &a.Length, &a.Size, &a.Thickness, &a.Count, &a.Price, &a.Boxwidth, &a.Boxheight, &a.Boxsize, &a.Boxweight, &a.Edge, &a.Flooring, &a.Installation, &a.Construction, &a.Gloss, &a.Waste, &a.Wear, &a.Features, &a.Layers, &a.Materials, &a.Meta)
	return
}

// msgp marshal
func (a *Product) MessagePack() []byte {
	b, _ := msgpack.Marshal(a)
	return b
}

// msgp unmarshal
func ParseProductMessagePack(v []byte) (a Product, err error) {
	err = msgpack.Unmarshal(v, &a)
	return
}

// fast json marshal
func (a *Product) Pack() []byte {
	var jsoner = jsoniter.ConfigCompatibleWithStandardLibrary
	b, _ := jsoner.Marshal(a)
	return b
}

// fast json unmarshal
func ParseProduct(v []byte) (a Product, err error) {
	var jsoner = jsoniter.ConfigCompatibleWithStandardLibrary
	err = jsoner.Unmarshal(v, &a)
	return
}

// Parse []any to json
func ParseTupleToProductJson(r []any) (a ProductJson) {
	for pos, x := range r {
		switch ProductIndexType(pos) {
		case IndexProductID:
			a.Set(FieldProductID, x)
		case IndexProductCreated:
			a.Set(FieldProductCreated, x)
		case IndexProductUpdated:
			a.Set(FieldProductUpdated, x)
		case IndexProductActive:
			a.Set(FieldProductActive, x)
		case IndexProductRef:
			a.Set(FieldProductRef, x)
		case IndexProductSku:
			a.Set(FieldProductSku, x)
		case IndexProductCategory:
			a.Set(FieldProductCategory, x)
		case IndexProductBrand:
			a.Set(FieldProductBrand, x)
		case IndexProductModel:
			a.Set(FieldProductModel, x)
		case IndexProductTitle:
			a.Set(FieldProductTitle, x)
		case IndexProductLine:
			a.Set(FieldProductLine, x)
		case IndexProductAbout:
			a.Set(FieldProductAbout, x)
		case IndexProductImage:
			a.Set(FieldProductImage, x)
		case IndexProductOriginal:
			a.Set(FieldProductOriginal, x)
		case IndexProductCost:
			a.Set(FieldProductCost, x)
		case IndexProductDiscount:
			a.Set(FieldProductDiscount, x)
		case IndexProductSample:
			a.Set(FieldProductSample, x)
		case IndexProductWidth:
			a.Set(FieldProductWidth, x)
		case IndexProductLength:
			a.Set(FieldProductLength, x)
		case IndexProductSize:
			a.Set(FieldProductSize, x)
		case IndexProductThickness:
			a.Set(FieldProductThickness, x)
		case IndexProductCount:
			a.Set(FieldProductCount, x)
		case IndexProductPrice:
			a.Set(FieldProductPrice, x)
		case IndexProductBoxwidth:
			a.Set(FieldProductBoxwidth, x)
		case IndexProductBoxheight:
			a.Set(FieldProductBoxheight, x)
		case IndexProductBoxsize:
			a.Set(FieldProductBoxsize, x)
		case IndexProductBoxweight:
			a.Set(FieldProductBoxweight, x)
		case IndexProductEdge:
			a.Set(FieldProductEdge, x)
		case IndexProductFlooring:
			a.Set(FieldProductFlooring, x)
		case IndexProductInstallation:
			a.Set(FieldProductInstallation, x)
		case IndexProductConstruction:
			a.Set(FieldProductConstruction, x)
		case IndexProductGloss:
			a.Set(FieldProductGloss, x)
		case IndexProductWaste:
			a.Set(FieldProductWaste, x)
		case IndexProductWear:
			a.Set(FieldProductWear, x)
		case IndexProductFeatures:
			a.Set(FieldProductFeatures, x)
		case IndexProductLayers:
			a.Set(FieldProductLayers, x)
		case IndexProductMaterials:
			a.Set(FieldProductMaterials, x)
		case IndexProductMeta:
			a.Set(FieldProductMeta, x)
		}
	}
	return
}

// NewProductJson create struct
func NewProductJson() ProductJson {
	return []byte("{}")
}

// ProductJson is a struct
type ProductJson []byte

// Set value
func (a *ProductJson) Set(k string, v any) *ProductJson {
	(*a) = jsons.Set((*a), k, v)
	return a
}

// Get value
func (a *ProductJson) Get(k string) jsons.Result {
	return jsons.Get((*a), k)
}

// Get value
func (a *ProductJson) DeleteFields(fields ...string) {
	(*a) = jsons.Delete((*a), fields...)
}

// ID set or get value
func (a *ProductJson) ID(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldProductID)
	}
	a.Set(FieldProductID, v[0])
	return
}

// Created set or get value
func (a *ProductJson) Created(v ...int64) (res int64) {
	if v == nil {
		return jsons.Int64((*a), FieldProductCreated)
	}
	a.Set(FieldProductCreated, v[0])
	return
}

// CreatedTime get value as time
func (a *ProductJson) CreatedTime() (res time.Time) {
	return time.Unix(a.Created(), 0)
}

// Updated set or get value
func (a *ProductJson) Updated(v ...int64) (res int64) {
	if v == nil {
		return jsons.Int64((*a), FieldProductUpdated)
	}
	a.Set(FieldProductUpdated, v[0])
	return
}

// UpdatedTime get value as time
func (a *ProductJson) UpdatedTime() (res time.Time) {
	return time.Unix(a.Updated(), 0)
}

// Active set or get value
func (a *ProductJson) Active(v ...bool) (res bool) {
	if v == nil {
		return jsons.Bool((*a), FieldProductActive)
	}
	a.Set(FieldProductActive, v[0])
	return
}

// Ref set or get value
func (a *ProductJson) Ref(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductRef)
	}
	a.Set(FieldProductRef, v[0])
	return
}

// Sku set or get value
func (a *ProductJson) Sku(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductSku)
	}
	a.Set(FieldProductSku, v[0])
	return
}

// Category set or get value
func (a *ProductJson) Category(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductCategory)
	}
	a.Set(FieldProductCategory, v[0])
	return
}

// Brand set or get value
func (a *ProductJson) Brand(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductBrand)
	}
	a.Set(FieldProductBrand, v[0])
	return
}

// Model set or get value
func (a *ProductJson) Model(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductModel)
	}
	a.Set(FieldProductModel, v[0])
	return
}

// Title set or get value
func (a *ProductJson) Title(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductTitle)
	}
	a.Set(FieldProductTitle, v[0])
	return
}

// Line set or get value
func (a *ProductJson) Line(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductLine)
	}
	a.Set(FieldProductLine, v[0])
	return
}

// About set or get value
func (a *ProductJson) About(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductAbout)
	}
	a.Set(FieldProductAbout, v[0])
	return
}

// Image set or get value
func (a *ProductJson) Image(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductImage)
	}
	a.Set(FieldProductImage, v[0])
	return
}

// Original set or get value
func (a *ProductJson) Original(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductOriginal)
	}
	a.Set(FieldProductOriginal, v[0])
	return
}

// Cost set or get value
func (a *ProductJson) Cost(v ...float64) (res float64) {
	if v == nil {
		return jsons.Float64((*a), FieldProductCost)
	}
	a.Set(FieldProductCost, v[0])
	return
}

// Discount set or get value
func (a *ProductJson) Discount(v ...float64) (res float64) {
	if v == nil {
		return jsons.Float64((*a), FieldProductDiscount)
	}
	a.Set(FieldProductDiscount, v[0])
	return
}

// Sample set or get value
func (a *ProductJson) Sample(v ...float64) (res float64) {
	if v == nil {
		return jsons.Float64((*a), FieldProductSample)
	}
	a.Set(FieldProductSample, v[0])
	return
}

// Width set or get value
func (a *ProductJson) Width(v ...float64) (res float64) {
	if v == nil {
		return jsons.Float64((*a), FieldProductWidth)
	}
	a.Set(FieldProductWidth, v[0])
	return
}

// Length set or get value
func (a *ProductJson) Length(v ...float64) (res float64) {
	if v == nil {
		return jsons.Float64((*a), FieldProductLength)
	}
	a.Set(FieldProductLength, v[0])
	return
}

// Size set or get value
func (a *ProductJson) Size(v ...float64) (res float64) {
	if v == nil {
		return jsons.Float64((*a), FieldProductSize)
	}
	a.Set(FieldProductSize, v[0])
	return
}

// Thickness set or get value
func (a *ProductJson) Thickness(v ...float64) (res float64) {
	if v == nil {
		return jsons.Float64((*a), FieldProductThickness)
	}
	a.Set(FieldProductThickness, v[0])
	return
}

// Count set or get value
func (a *ProductJson) Count(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldProductCount)
	}
	a.Set(FieldProductCount, v[0])
	return
}

// Price set or get value
func (a *ProductJson) Price(v ...float64) (res float64) {
	if v == nil {
		return jsons.Float64((*a), FieldProductPrice)
	}
	a.Set(FieldProductPrice, v[0])
	return
}

// Boxwidth set or get value
func (a *ProductJson) Boxwidth(v ...float64) (res float64) {
	if v == nil {
		return jsons.Float64((*a), FieldProductBoxwidth)
	}
	a.Set(FieldProductBoxwidth, v[0])
	return
}

// Boxheight set or get value
func (a *ProductJson) Boxheight(v ...float64) (res float64) {
	if v == nil {
		return jsons.Float64((*a), FieldProductBoxheight)
	}
	a.Set(FieldProductBoxheight, v[0])
	return
}

// Boxsize set or get value
func (a *ProductJson) Boxsize(v ...float64) (res float64) {
	if v == nil {
		return jsons.Float64((*a), FieldProductBoxsize)
	}
	a.Set(FieldProductBoxsize, v[0])
	return
}

// Boxweight set or get value
func (a *ProductJson) Boxweight(v ...float64) (res float64) {
	if v == nil {
		return jsons.Float64((*a), FieldProductBoxweight)
	}
	a.Set(FieldProductBoxweight, v[0])
	return
}

// Edge set or get value
func (a *ProductJson) Edge(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductEdge)
	}
	a.Set(FieldProductEdge, v[0])
	return
}

// Flooring set or get value
func (a *ProductJson) Flooring(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductFlooring)
	}
	a.Set(FieldProductFlooring, v[0])
	return
}

// Installation set or get value
func (a *ProductJson) Installation(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductInstallation)
	}
	a.Set(FieldProductInstallation, v[0])
	return
}

// Construction set or get value
func (a *ProductJson) Construction(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductConstruction)
	}
	a.Set(FieldProductConstruction, v[0])
	return
}

// Gloss set or get value
func (a *ProductJson) Gloss(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldProductGloss)
	}
	a.Set(FieldProductGloss, v[0])
	return
}

// Waste set or get value
func (a *ProductJson) Waste(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldProductWaste)
	}
	a.Set(FieldProductWaste, v[0])
	return
}

// Wear set or get value
func (a *ProductJson) Wear(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldProductWear)
	}
	a.Set(FieldProductWear, v[0])
	return
}

// Features set or get value
func (a *ProductJson) Features(v ...map[string]bool) (res map[string]bool) {
	if v == nil {
		return jsons.MapBool((*a), FieldProductFeatures)
	}
	a.Set(FieldProductFeatures, v[0])
	return
}

// FeaturesAdd add values
func (a *ProductJson) FeaturesAdd(k string, v bool) {
	maps := a.Features()
	maps[k] = v
	a.Features(maps)
}

// FeaturesDelete add unique values only
func (a *ProductJson) FeaturesDelete(k string) {
	maps := a.Features()
	delete(maps, k)
	a.Features(maps)
}

// FeaturesHas check value
func (a *ProductJson) FeaturesHas(k string) bool {
	maps := a.Features()
	return !reflect.ValueOf(maps[k]).IsZero()
}

// Layers set or get value
func (a *ProductJson) Layers(v ...map[string]bool) (res map[string]bool) {
	if v == nil {
		return jsons.MapBool((*a), FieldProductLayers)
	}
	a.Set(FieldProductLayers, v[0])
	return
}

// LayersAdd add values
func (a *ProductJson) LayersAdd(k string, v bool) {
	maps := a.Layers()
	maps[k] = v
	a.Layers(maps)
}

// LayersDelete add unique values only
func (a *ProductJson) LayersDelete(k string) {
	maps := a.Layers()
	delete(maps, k)
	a.Layers(maps)
}

// LayersHas check value
func (a *ProductJson) LayersHas(k string) bool {
	maps := a.Layers()
	return !reflect.ValueOf(maps[k]).IsZero()
}

// Materials set or get value
func (a *ProductJson) Materials(v ...map[string]bool) (res map[string]bool) {
	if v == nil {
		return jsons.MapBool((*a), FieldProductMaterials)
	}
	a.Set(FieldProductMaterials, v[0])
	return
}

// MaterialsAdd add values
func (a *ProductJson) MaterialsAdd(k string, v bool) {
	maps := a.Materials()
	maps[k] = v
	a.Materials(maps)
}

// MaterialsDelete add unique values only
func (a *ProductJson) MaterialsDelete(k string) {
	maps := a.Materials()
	delete(maps, k)
	a.Materials(maps)
}

// MaterialsHas check value
func (a *ProductJson) MaterialsHas(k string) bool {
	maps := a.Materials()
	return !reflect.ValueOf(maps[k]).IsZero()
}

// Meta set or get value
func (a *ProductJson) Meta(v ...map[string]any) (res map[string]any) {
	if v == nil {
		return jsons.MapAny((*a), FieldProductMeta)
	}
	a.Set(FieldProductMeta, v[0])
	return
}

// MetaAdd add values
func (a *ProductJson) MetaAdd(k string, v any) {
	maps := a.Meta()
	maps[k] = v
	a.Meta(maps)
}

// MetaDelete add unique values only
func (a *ProductJson) MetaDelete(k string) {
	maps := a.Meta()
	delete(maps, k)
	a.Meta(maps)
}

// MetaHas check value
func (a *ProductJson) MetaHas(k string) bool {
	maps := a.Meta()
	return !reflect.ValueOf(maps[k]).IsZero()
}

// sql sqlProduct class
var ProductSQL sqlProduct

type sqlProduct int

// parse sql query
func (a *sqlProduct) TableName() (res string) {
	return "products"
}

// parse sql query
func (a *sqlProduct) Get(conn *pgxpool.Conn, c context.Context, id any, fields ...ProductIndexType) (res *Product, err error) {
	if fields == nil {
		fields = []ProductIndexType{IndexProductID, IndexProductCreated, IndexProductUpdated, IndexProductActive, IndexProductRef, IndexProductSku, IndexProductCategory, IndexProductBrand, IndexProductModel, IndexProductTitle, IndexProductLine, IndexProductAbout, IndexProductImage, IndexProductOriginal, IndexProductCost, IndexProductDiscount, IndexProductSample, IndexProductWidth, IndexProductLength, IndexProductSize, IndexProductThickness, IndexProductCount, IndexProductPrice, IndexProductBoxwidth, IndexProductBoxheight, IndexProductBoxsize, IndexProductBoxweight, IndexProductEdge, IndexProductFlooring, IndexProductInstallation, IndexProductConstruction, IndexProductGloss, IndexProductWaste, IndexProductWear, IndexProductFeatures, IndexProductLayers, IndexProductMaterials, IndexProductMeta}
	}

	var list []string
	for _, x := range fields {
		list = append(list, x.SQLName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s ", fieldlist)
	q = q + "from products where id = $1 limit 1"
	res = new(Product)
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
func (a *sqlProduct) Row(conn *pgxpool.Conn, c context.Context, eq map[string]any, fields ...ProductIndexType) (res *Product, err error) {
	if fields == nil {
		fields = []ProductIndexType{IndexProductID, IndexProductCreated, IndexProductUpdated, IndexProductActive, IndexProductRef, IndexProductSku, IndexProductCategory, IndexProductBrand, IndexProductModel, IndexProductTitle, IndexProductLine, IndexProductAbout, IndexProductImage, IndexProductOriginal, IndexProductCost, IndexProductDiscount, IndexProductSample, IndexProductWidth, IndexProductLength, IndexProductSize, IndexProductThickness, IndexProductCount, IndexProductPrice, IndexProductBoxwidth, IndexProductBoxheight, IndexProductBoxsize, IndexProductBoxweight, IndexProductEdge, IndexProductFlooring, IndexProductInstallation, IndexProductConstruction, IndexProductGloss, IndexProductWaste, IndexProductWear, IndexProductFeatures, IndexProductLayers, IndexProductMaterials, IndexProductMeta}
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
	res = new(Product)
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
func (a *sqlProduct) All(conn *pgxpool.Conn, c context.Context, fields ...ProductIndexType) (res []*Product, err error) {
	if fields == nil {
		fields = []ProductIndexType{IndexProductID, IndexProductCreated, IndexProductUpdated, IndexProductActive, IndexProductRef, IndexProductSku, IndexProductCategory, IndexProductBrand, IndexProductModel, IndexProductTitle, IndexProductLine, IndexProductAbout, IndexProductImage, IndexProductOriginal, IndexProductCost, IndexProductDiscount, IndexProductSample, IndexProductWidth, IndexProductLength, IndexProductSize, IndexProductThickness, IndexProductCount, IndexProductPrice, IndexProductBoxwidth, IndexProductBoxheight, IndexProductBoxsize, IndexProductBoxweight, IndexProductEdge, IndexProductFlooring, IndexProductInstallation, IndexProductConstruction, IndexProductGloss, IndexProductWaste, IndexProductWear, IndexProductFeatures, IndexProductLayers, IndexProductMaterials, IndexProductMeta}
	}

	var list []string
	for _, x := range fields {
		list = append(list, x.SQLName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s", fieldlist)
	q += " from products"
	rows, err := conn.Query(c, q)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item Product
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
func (a *sqlProduct) List(conn *pgxpool.Conn, c context.Context, limit, offset int, fields ...ProductIndexType) (res []*Product, err error) {
	if fields == nil {
		fields = []ProductIndexType{IndexProductID, IndexProductCreated, IndexProductUpdated, IndexProductActive, IndexProductRef, IndexProductSku, IndexProductCategory, IndexProductBrand, IndexProductModel, IndexProductTitle, IndexProductLine, IndexProductAbout, IndexProductImage, IndexProductOriginal, IndexProductCost, IndexProductDiscount, IndexProductSample, IndexProductWidth, IndexProductLength, IndexProductSize, IndexProductThickness, IndexProductCount, IndexProductPrice, IndexProductBoxwidth, IndexProductBoxheight, IndexProductBoxsize, IndexProductBoxweight, IndexProductEdge, IndexProductFlooring, IndexProductInstallation, IndexProductConstruction, IndexProductGloss, IndexProductWaste, IndexProductWear, IndexProductFeatures, IndexProductLayers, IndexProductMaterials, IndexProductMeta}
	}

	var list []string
	for _, x := range fields {
		list = append(list, x.SQLName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s", fieldlist)
	q += " from products limit $1 offset $2"
	rows, err := conn.Query(c, q, limit, offset)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item Product
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
func (a *sqlProduct) Update(conn *pgxpool.Conn, c context.Context, id any, k string, v any) (err error) {
	if !ProductValidKey(k) {
		return fmt.Errorf("invalid key")
	}
	q := fmt.Sprintf("update products set %s = $1 where id = $2", k)
	_, err = conn.Exec(c, q, v, id)
	return
}

// update sql query
func (a *sqlProduct) Updates(conn *pgxpool.Conn, c context.Context, id any, keys map[string]any) (err error) {

	if keys == nil {
		return fmt.Errorf("emptykeys")
	}

	var fields []string
	var values []any
	var count int
	for k, v := range keys {
		if !ProductValidKey(k) {
			return fmt.Errorf(k)
		}
		in := ProductKeyIndex(k)
		count++
		fields = append(fields, fmt.Sprintf("%s = $%d", in.SQLName(), count))
		values = append(values, v)
	}

	list := strings.Join(fields, ", ")
	count++
	values = append(values, id)

	q := fmt.Sprintf("update products set %s where id = $%d", list, count)
	_, err = conn.Exec(c, q, values...)
	return
}

// update sql query
func (a *sqlProduct) UpdateFeatures(conn *pgxpool.Conn, c context.Context, id any, k string, v any) (err error) {
	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update products set features = features || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *sqlProduct) UpdatesFeatures(conn *pgxpool.Conn, c context.Context, id any, keys map[string]any) (err error) {
	b, _ := json.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update products set features = features || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *sqlProduct) DeleteKeyFeatures(conn *pgxpool.Conn, c context.Context, id any, k string) (err error) {
	q := "update products set features = features - $1 where id = $2"
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *sqlProduct) RenameKeyFeatures(conn *pgxpool.Conn, c context.Context, id any, k, newkey string) (err error) {
	q := "update products set features = features - $1 || jsonb_build_object($2, features->$1) where id = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *sqlProduct) UpdateLayers(conn *pgxpool.Conn, c context.Context, id any, k string, v any) (err error) {
	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update products set layers = layers || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *sqlProduct) UpdatesLayers(conn *pgxpool.Conn, c context.Context, id any, keys map[string]any) (err error) {
	b, _ := json.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update products set layers = layers || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *sqlProduct) DeleteKeyLayers(conn *pgxpool.Conn, c context.Context, id any, k string) (err error) {
	q := "update products set layers = layers - $1 where id = $2"
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *sqlProduct) RenameKeyLayers(conn *pgxpool.Conn, c context.Context, id any, k, newkey string) (err error) {
	q := "update products set layers = layers - $1 || jsonb_build_object($2, layers->$1) where id = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *sqlProduct) UpdateMaterials(conn *pgxpool.Conn, c context.Context, id any, k string, v any) (err error) {
	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update products set materials = materials || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *sqlProduct) UpdatesMaterials(conn *pgxpool.Conn, c context.Context, id any, keys map[string]any) (err error) {
	b, _ := json.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update products set materials = materials || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *sqlProduct) DeleteKeyMaterials(conn *pgxpool.Conn, c context.Context, id any, k string) (err error) {
	q := "update products set materials = materials - $1 where id = $2"
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *sqlProduct) RenameKeyMaterials(conn *pgxpool.Conn, c context.Context, id any, k, newkey string) (err error) {
	q := "update products set materials = materials - $1 || jsonb_build_object($2, materials->$1) where id = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

// update sql query
func (a *sqlProduct) UpdateMeta(conn *pgxpool.Conn, c context.Context, id any, k string, v any) (err error) {
	// json escape
	res := strings.ReplaceAll(jsons.Creates(k, v).String(), "$$", "$ $")
	q := fmt.Sprintf("update products set meta = meta || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// update sql query
func (a *sqlProduct) UpdatesMeta(conn *pgxpool.Conn, c context.Context, id any, keys map[string]any) (err error) {
	b, _ := json.Marshal(keys)
	// json escape
	res := strings.ReplaceAll(string(b), "$$", "$ $")
	q := fmt.Sprintf("update products set meta = meta || $$%s$$::jsonb where id = $1", res)
	_, err = conn.Exec(c, q, id)
	return
}

// delete key from jsonb
func (a *sqlProduct) DeleteKeyMeta(conn *pgxpool.Conn, c context.Context, id any, k string) (err error) {
	q := "update products set meta = meta - $1 where id = $2"
	_, err = conn.Exec(c, q, k, id)
	return
}

// rename map key jsonb
func (a *sqlProduct) RenameKeyMeta(conn *pgxpool.Conn, c context.Context, id any, k, newkey string) (err error) {
	q := "update products set meta = meta - $1 || jsonb_build_object($2, meta->$1) where id = $3"
	_, err = conn.Exec(c, q, k, newkey, id)
	return
}

type ProductQuery struct {
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

func (a *ProductQuery) Render() (sql string, fields []ProductIndexType, values []any) {

	switch a.Fields == nil {
	case true:
		fields = []ProductIndexType{IndexProductID, IndexProductCreated, IndexProductUpdated, IndexProductActive, IndexProductRef, IndexProductSku, IndexProductCategory, IndexProductBrand, IndexProductModel, IndexProductTitle, IndexProductLine, IndexProductAbout, IndexProductImage, IndexProductOriginal, IndexProductCost, IndexProductDiscount, IndexProductSample, IndexProductWidth, IndexProductLength, IndexProductSize, IndexProductThickness, IndexProductCount, IndexProductPrice, IndexProductBoxwidth, IndexProductBoxheight, IndexProductBoxsize, IndexProductBoxweight, IndexProductEdge, IndexProductFlooring, IndexProductInstallation, IndexProductConstruction, IndexProductGloss, IndexProductWaste, IndexProductWear, IndexProductFeatures, IndexProductLayers, IndexProductMaterials, IndexProductMeta}
	default:
		for _, x := range a.Fields {
			if !ProductValidKey(x) {
				continue
			}
			fields = append(fields, ProductKeyIndex(x))
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

	// limit
	if a.Limit < 1 {
		a.Limit = 20
	}

	// offset
	if a.Offset < 0 {
		a.Offset = 0
	}

	// sql
	var list []string

	var count int

	// select
	list = append(list, "select")
	list = append(list, strings.Join(fieldsStrings, ", "))

	// from
	list = append(list, "from products")

	var andlist []string

	// EQ where
	for k, v := range a.EQ {
		if !ProductValidKey(k) {
			continue
		}
		p := ProductKeyIndex(k)
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

		if !ProductValidKey(k) {
			continue
		}
		count++
		andlist = append(andlist, fmt.Sprintf("%s > $%d", k, count))
		values = append(values, v)
	}

	// LT where
	for k, v := range a.LT {
		if !ProductValidKey(k) {
			continue
		}
		count++
		andlist = append(andlist, fmt.Sprintf("%s < $%d", k, count))
		values = append(values, v)
	}

	// NOT where
	for k, v := range a.NOT {
		if !ProductValidKey(k) {
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
		if !ProductValidKey(k) {
			continue
		}
		count++
		andlist = append(andlist, fmt.Sprintf("%s ilike $%d", k, count))
		values = append(values, "%"+v+"%")
	}

	// IN where
	if a.IN != nil {
		for k, v := range a.IN {
			if !ProductValidKey(k) {
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
			if !ProductValidKey(k) {
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
	if a.Sort != "" && ProductValidKey(a.Sort) {
		list = append(list, fmt.Sprintf("order by %s", a.Sort))
		if a.Desc {
			list = append(list, "desc")
		}
	}

	// limit, offset
	list = append(list, fmt.Sprintf("limit %d offset %d", a.Limit, a.Offset))

	// render sql
	sql = strings.Join(list, " ")
	return
}

func NewProductQuery() *ProductQuery {
	a := new(ProductQuery)
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

func (a *sqlProduct) Search(conn *pgxpool.Conn, c context.Context, q *ProductQuery) (res []*Product, err error) {

	sql, fields, values := q.Render()

	rows, err := conn.Query(c, sql, values...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item Product
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
func (a *sqlProduct) Delete(conn *pgxpool.Conn, c context.Context, id any) (err error) {
	_, err = conn.Exec(c, "delete from products where id = $1", id)
	return

}

// delete item where i = 1 and w = 'nice'
func (a *sqlProduct) DeleteWhere(conn *pgxpool.Conn, c context.Context, where string) (err error) {
	_, err = conn.Exec(c, fmt.Sprintf("delete from products where %s", where))
	return

}

// Insert struct and return int id
func (a *sqlProduct) Insert(conn *pgxpool.Conn, c context.Context, v *Product) (id int, err error) {
	q := "insert into products (created, updated, active, ref, sku, category, brand, model, title, line, about, image, original, cost, discount, sample, width, length, size, thickness, count, price, boxwidth, boxheight, boxsize, boxweight, edge, flooring, installation, construction, gloss, waste, wear, features, layers, materials, meta) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37) returning id"
	err = conn.QueryRow(c, q, v.sqlTuple()...).Scan(&id)
	return
}

// has value in db
func (a *sqlProduct) Has(conn *pgxpool.Conn, c context.Context, field ProductIndexType, v any) (has bool, err error) {
	q := fmt.Sprintf("select exists (select id from products where %s = $1 limit 1)", field.SQLName())
	err = conn.QueryRow(c, q, v).Scan(&has)
	return
}

//has value in db

// new db
func NewProductDB(pool *pgxpool.Pool, timeout time.Duration) (a *ProductDB) {
	a = new(ProductDB)
	a.Pool = pool
	a.Timeout = timeout
	return
}

// read
type ProductDB struct {
	Pool    *pgxpool.Pool
	Timeout time.Duration
}

// conn
func (a *ProductDB) Conn(f func(conn *pgxpool.Conn, c context.Context) (err errors.E)) (err errors.E) {
	c, cancel := context.WithTimeout(context.Background(), a.Timeout)
	defer cancel()

	conn, er := a.Pool.Acquire(c)
	if er != nil {
		err = errors.PGX(er)
		return
	}
	defer conn.Release()
	return f(conn, c)
}

// Create table
func (a *sqlProduct) CreateTable(conn *pgxpool.Conn, c context.Context) (err error) {
	q := `create table if not exists products (
	id                                           bigserial primary key,
	created                                      bigint default extract(epoch from now()),
	updated                                      bigint default extract(epoch from now()),
	active                                       boolean default false,
	ref                                          text,
	sku                                          text,
	category                                     text,
	brand                                        text,
	model                                        text,
	title                                        text,
	line                                         text,
	about                                        text,
	image                                        text,
	original                                     text,
	cost                                         double precision,
	discount                                     double precision,
	sample                                       double precision,
	width                                        double precision,
	length                                       double precision,
	size                                         double precision,
	thickness                                    double precision,
	count                                        bigint,
	price                                        double precision,
	boxwidth                                     double precision,
	boxheight                                    double precision,
	boxsize                                      double precision,
	boxweight                                    double precision,
	edge                                         text,
	flooring                                     text,
	installation                                 text,
	construction                                 text,
	gloss                                        text,
	waste                                        bigint,
	wear                                         bigint,
	features                                     jsonb default '{}'::jsonb,
	layers                                       jsonb default '{}'::jsonb,
	materials                                    jsonb default '{}'::jsonb,
	meta                                         jsonb default '{}'::jsonb
)
`
	_, err = conn.Exec(c, q)
	return
}
