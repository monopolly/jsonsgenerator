package news

import (
	"context"
	json "encoding/json"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/jackc/pgx/v5/pgxpool"
	jsoniter "github.com/json-iterator/go"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
	"github.com/monopolly/cast"
	"github.com/monopolly/jsons"
	"github.com/niubaoshu/gotiny"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"
)

/*
	Auto-generate. Do not change!
	Help to auto-generated correct models for golang structures.
	Sergey Keplin (c) 2026
	a@senthy.com, lava.mobi@gmail.com, https://t.me/martinprestone
	github.com/monopolly/jsons
	Ex: //! [] go=News js=NewsJson sql=news noinit up...

	STRUCT:
	!omit                         No omit tag for json
	proto=NewsProto               Generate proto file and compile to current package
	enum                          Generate swift enum
	optimize                      Try to create golang padding optimized struct
	test                          Generate go tests for generated model
	noprefix                      Generate simple index IndexID instead IndexNewsID
	lock                          Lock model for generation. Can't change model.
	noinit                        No New() init function for struct
	js=NewsJson                   Change json struct names. Ex: js=NewsJson
	sql=news                      Set sql table name. Ex: sql=accounts
	ch=views                      Set clickhouse sql table name. Ex: ch=views
	chengine=MergeTree            Optional, mergeTree by default
	gotiny                        Create gotiny marshal/unmarshal
	demo                          Generate demo json file with default values
	go=News                       Set golang struct names. Ex: go=News1 > type News1 struct{}
	ts=news                       Set typescript struct names. Ex: ts=NewsJson
	swift                         Generate swift model
	debug                         Debug mode
	msgp                          Create message pack marshal/unmarshal

	Field:
	title{}                       Add custom title for index. title{Nice & Sweet}
	desc{}                        Add custom desc for index. desc{This field for success}

	GO:
	nofunc                        Do not create any jsons functions for fiels
	desc=""                       Add custom desc for index
	req{}                         Add user required fields
	desc{}                        Add field user title desc{Use it nice}
	up                            Make uppercase for functions
	must                          Create one validation function for all must fields
	type=""                       Replace golang type. Ex: type="[]*News"
	name=""                       Replace golang struct name. Ex: name="NewsList"
	title=""                      Add custom title for index
	title{}                       Add field user title title{Nice}
	#                             Add lists for fields. Ex: id int //#readonly #must...
	@                             Add custom structs for fields. Ex: id int //@public @team...

	SQL:
	type="jsonb"                  Rewrite sql type for field. Ex: type="jsonb"
	add="primary key"             Append sql for field. Ex: replace="primary key"
	unique="groupname"            Add unique fields constrains by group (you need set group name, then generator join fields). Ex: unique="group1" unique="group2"
	ver="v4"                      Create an alter table record in SQL file. Add new column. With new version in comment. Ex: ver="2"
	renames="oldname"             Create an alter table record in SQL file. Rename table column.
	altertable                    Add field line Alter table to SQL file with current time comment
	idx                           Add index fields by group name. Ex: idx="nameIndex" idx="credsIndex"
	defaults                      Add default value based on field type
	unix                          Add default value: extract(epoch from now())
	skip                          Skip field for sql queries
	replace="bigint primary key"  Rewrite sql for field. Ex: replace="bigint primary key"
	noinsert                      Do use field for insert function
	index                         Create simple default index or gin for jsonb
	search                        Add tsvector index by group. Ex: search="tsv": tsv tsvector GENERATED ALWAYS AS (to_tsvector('simple', title || ' ' || brand)). For search: SELECT brand, title FROM assets WHERE search @@ to_tsquery('english', 'f8');

	Clickhouse:
	low                           Wrap field type in LowCardinality(...)
	type=ip                       Use ClickHouse IPv4 type for IP address fields

	PROTO:
	type="string"                 Rewrite proto type for field. Ex: type="google.protobuf.Value"
	name="field"                  Rewrite proto field name
	skip                          Skip field for proto

	SWIFT:
	file=""                       Create another swift file for this field Swift model. file="f1", file="f2"
	must                          Swift field with required values
	skip                          Ignore field for Swift
	type=""                       Replace type for Swift model. Ex: type="string"

	JSON:
	raw                           Set Raw json function inside field
	name=""                       Replace json field name. Ex: name="sid"
	skip                          Skip json field
	inc                           Add inc jsons function for numbers fields
	bool                          Add set jsons function for bool fields
	time                          Create convert jsons function for unixtime fields
*/

// field type
type Stat2IndexType int

// int index
const (
	IndexStat2ID = Stat2IndexType(iota)
	IndexStat2UID
	IndexStat2Type
	IndexStat2Name
	IndexStat2Source
	IndexStat2IP
	IndexStat2Image
)

// string index
const (
	FieldStat2ID     = "id"     // int sql{inc}  #readonly
	FieldStat2UID    = "uid"    // int ch{primary}
	FieldStat2Type   = "type"   // int ch{low}
	FieldStat2Name   = "name"   // string
	FieldStat2Source = "source" // string
	FieldStat2IP     = "ip"     // string
	FieldStat2Image  = "image"  // bool
)

// index func
func Stat2Indexes() []Stat2IndexType {
	return []Stat2IndexType{IndexStat2ID, IndexStat2UID, IndexStat2Type, IndexStat2Name, IndexStat2Source, IndexStat2IP, IndexStat2Image}

}

// 80 bytes (go padding)
//
//easyjson:json
type Stat2 struct {
	ID     int    `json:"id,omitempty"`     // sql{inc}  #readonly
	UID    int    `json:"uid,omitempty"`    // ch{primary}
	Type   int    `json:"type,omitempty"`   // ch{low}
	Name   string `json:"name,omitempty"`   //
	Source string `json:"source,omitempty"` //
	IP     string `json:"ip,omitempty"`     //
	Image  bool   `json:"image,omitempty"`  //
}

// Parse []any to ID struct
func ParseStat2ToStruct(r []any) (a *Stat2) {
	a = new(Stat2)

	for pos, x := range r {
		switch Stat2IndexType(pos) {
		case IndexStat2ID:
			cast.Convert(&a.ID, x) //int
		case IndexStat2UID:
			cast.Convert(&a.UID, x) //int
		case IndexStat2Type:
			cast.Convert(&a.Type, x) //int
		case IndexStat2Name:
			cast.Convert(&a.Name, x) //string
		case IndexStat2Source:
			cast.Convert(&a.Source, x) //string
		case IndexStat2IP:
			cast.Convert(&a.IP, x) //string
		case IndexStat2Image:
			cast.Convert(&a.Image, x) //bool
		}
	}
	return
}

// Tuple create an array from struct
func (a *Stat2) Tuple() (r []any) {
	return []any{a.ID, a.UID, a.Type, a.Name, a.Source, a.IP, a.Image}
}

// Tuple create an array from struct
func (a *Stat2) sqlTuple() (r []any) {
	return []any{a.UID, a.Type, a.Name, a.Source, a.IP, a.Image}
}

// Tuple create an array from struct
func (a *Stat2) sqlAllTuples() (r []any) {
	return []any{a.ID, a.UID, a.Type, a.Name, a.Source, a.IP, a.Image}
}

// Tuple create an array from struct for clickhouse
func (a *Stat2) clickhouseTuple() (r []any) {
	return []any{a.ID, a.UID, a.Type, a.Name, a.Source, net.ParseIP(a.IP), a.Image}
}

// update struct with function
func (a *Stat2) Update(k string, x any) {
	switch k {
	case "id":
		cast.Convert(&a.ID, x) //int
	case "uid":
		cast.Convert(&a.UID, x) //int
	case "type":
		cast.Convert(&a.Type, x) //int
	case "name":
		cast.Convert(&a.Name, x) //string
	case "source":
		cast.Convert(&a.Source, x) //string
	case "ip":
		cast.Convert(&a.IP, x) //string
	case "image":
		cast.Convert(&a.Image, x) //bool
	}
}

// get struct value with function
func (a *Stat2) Get(k string) (v any) {
	switch k {
	case "id":
		return a.ID //int
	case "uid":
		return a.UID //int
	case "type":
		return a.Type //int
	case "name":
		return a.Name //string
	case "source":
		return a.Source //string
	case "ip":
		return a.IP //string
	case "image":
		return a.Image //bool
	}
	return
}

// get any struct value as string
func (a *Stat2) String(k string) (v string) {
	switch k {
	case "id":
		return fmt.Sprint(a.ID) //int
	case "uid":
		return fmt.Sprint(a.UID) //int
	case "type":
		return fmt.Sprint(a.Type) //int
	case "name":
		return fmt.Sprint(a.Name) //string
	case "source":
		return fmt.Sprint(a.Source) //string
	case "ip":
		return fmt.Sprint(a.IP) //string
	case "image":
		return fmt.Sprint(a.Image) //bool
	}
	return
}

// Struct to json
func (a *Stat2) ToJson() (r []byte) {
	js := jsons.Create().
		Add(FieldStat2ID, a.ID).
		Add(FieldStat2UID, a.UID).
		Add(FieldStat2Type, a.Type).
		Add(FieldStat2Name, a.Name).
		Add(FieldStat2Source, a.Source).
		Add(FieldStat2IP, a.IP).
		Add(FieldStat2Image, a.Image)
	return js.Bytes()
}

func Stat2ReadonlyList() []Stat2IndexType {
	return []Stat2IndexType{IndexStat2ID}
}

func (a Stat2IndexType) Readonly() bool {
	switch a {
	case IndexStat2ID:
		return true
	default:
		return false
	}
}

// key index string
func (a Stat2IndexType) String() string {
	switch a {
	case IndexStat2ID:
		return "id"
	case IndexStat2UID:
		return "uid"
	case IndexStat2Type:
		return "type"
	case IndexStat2Name:
		return "name"
	case IndexStat2Source:
		return "source"
	case IndexStat2IP:
		return "ip"
	case IndexStat2Image:
		return "image"
	default:
		return ""
	}
}

// key index string
func (a Stat2IndexType) SQLName() string {
	switch a {
	case IndexStat2ID:
		return "id"
	case IndexStat2UID:
		return "uid"
	case IndexStat2Type:
		return "type"
	case IndexStat2Name:
		return "name"
	case IndexStat2Source:
		return "source"
	case IndexStat2IP:
		return "ip"
	case IndexStat2Image:
		return "image"
	default:
		return ""
	}
}

// key index clickhouse string
func (a Stat2IndexType) ClickhouseName() string {
	switch a {
	case IndexStat2ID:
		return "id"
	case IndexStat2UID:
		return "uid"
	case IndexStat2Type:
		return "type"
	case IndexStat2Name:
		return "name"
	case IndexStat2Source:
		return "source"
	case IndexStat2IP:
		return "ip"
	case IndexStat2Image:
		return "image"
	default:
		return ""
	}
}

// key index type
func (a Stat2IndexType) Type() string {
	switch a {
	case IndexStat2ID:
		return "int"
	case IndexStat2UID:
		return "int"
	case IndexStat2Type:
		return "int"
	case IndexStat2Name:
		return "string"
	case IndexStat2Source:
		return "string"
	case IndexStat2IP:
		return "string"
	case IndexStat2Image:
		return "bool"
	default:
		return ""
	}
}

// custom title
func (a Stat2IndexType) Title() string {
	switch a {
	case IndexStat2ID:
		return "ID"
	case IndexStat2UID:
		return "UID"
	case IndexStat2Type:
		return "Type"
	case IndexStat2Name:
		return "Name"
	case IndexStat2Source:
		return "Source"
	case IndexStat2IP:
		return "IP"
	case IndexStat2Image:
		return "Image"
	default:
		return ""
	}
}

// custom desc
func (a Stat2IndexType) Desc() string {
	switch a {
	default:
		return ""
	}
}

// struct key to index
func Stat2KeyIndex(key string) Stat2IndexType {
	switch key {
	case "id":
		return IndexStat2ID
	case "uid":
		return IndexStat2UID
	case "type":
		return IndexStat2Type
	case "name":
		return IndexStat2Name
	case "source":
		return IndexStat2Source
	case "ip":
		return IndexStat2IP
	case "image":
		return IndexStat2Image
	default:
		return 0
	}
}

// valid struct key check
func Stat2ValidKey(key string) bool {
	switch key {
	case "id", "uid", "type", "name", "source", "ip", "image":
		return true
	default:
		return false
	}
}

// struct to map
func (a *Stat2) Map() map[string]any {
	return map[string]any{
		"id":     a.ID,
		"uid":    a.UID,
		"type":   a.Type,
		"name":   a.Name,
		"source": a.Source,
		"ip":     a.IP,
		"image":  a.Image,
	}
}

// struct to map
func (a *Stat2) Iterate(f func(k Stat2IndexType, v any)) {
	for _, x := range Stat2Indexes() {
		f(x, a.Get(x.String()))
	}
}

// gotiny marshal
func (a *Stat2) MarshalGotiny() []byte {
	return gotiny.Marshal(&a.ID, &a.UID, &a.Type, &a.Name, &a.Source, &a.IP, &a.Image)
}

// parse gotiny
func UnmarshalStat2Gotiny(v []byte) (a Stat2) {
	gotiny.Unmarshal(v, &a.ID, &a.UID, &a.Type, &a.Name, &a.Source, &a.IP, &a.Image)
	return
}

// fast json marshal
func (a *Stat2) Pack() []byte {
	var jsoner = jsoniter.ConfigCompatibleWithStandardLibrary
	b, _ := jsoner.Marshal(a)
	return b
}

// fast json unmarshal
func ParseStat2(v []byte) (a Stat2, err error) {
	var jsoner = jsoniter.ConfigCompatibleWithStandardLibrary
	err = jsoner.Unmarshal(v, &a)
	return
}

// Parse []any to json
func ParseTupleToStat2Json(r []any) (a Stat2Json) {
	for pos, x := range r {
		switch Stat2IndexType(pos) {
		case IndexStat2ID:
			a.Set(FieldStat2ID, x)
		case IndexStat2UID:
			a.Set(FieldStat2UID, x)
		case IndexStat2Type:
			a.Set(FieldStat2Type, x)
		case IndexStat2Name:
			a.Set(FieldStat2Name, x)
		case IndexStat2Source:
			a.Set(FieldStat2Source, x)
		case IndexStat2IP:
			a.Set(FieldStat2IP, x)
		case IndexStat2Image:
			a.Set(FieldStat2Image, x)
		}
	}
	return
}

//NewStat2Json create struct
func NewStat2Json() Stat2Json {
	return []byte("{}")
}

//Stat2Json is a struct
type Stat2Json []byte

//Set value
func (a *Stat2Json) Set(k string, v any) *Stat2Json {
	(*a) = jsons.Set((*a), k, v)
	return a
}

//Get value
func (a *Stat2Json) Get(k string) jsons.Result {
	return jsons.Get((*a), k)
}

//Get value
func (a *Stat2Json) DeleteFields(fields ...string) {
	(*a) = jsons.Delete((*a), fields...)
}

//ID set or get value
func (a *Stat2Json) ID(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldStat2ID)
	}
	a.Set(FieldStat2ID, v[0])
	return
}

//UID set or get value
func (a *Stat2Json) UID(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldStat2UID)
	}
	a.Set(FieldStat2UID, v[0])
	return
}

//Type set or get value
func (a *Stat2Json) Type(v ...int) (res int) {
	if v == nil {
		return jsons.Int((*a), FieldStat2Type)
	}
	a.Set(FieldStat2Type, v[0])
	return
}

//Name set or get value
func (a *Stat2Json) Name(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldStat2Name)
	}
	a.Set(FieldStat2Name, v[0])
	return
}

//Source set or get value
func (a *Stat2Json) Source(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldStat2Source)
	}
	a.Set(FieldStat2Source, v[0])
	return
}

//IP set or get value
func (a *Stat2Json) IP(v ...string) (res string) {
	if v == nil {
		return jsons.String((*a), FieldStat2IP)
	}
	a.Set(FieldStat2IP, v[0])
	return
}

//Image set or get value
func (a *Stat2Json) Image(v ...bool) (res bool) {
	if v == nil {
		return jsons.Bool((*a), FieldStat2Image)
	}
	a.Set(FieldStat2Image, v[0])
	return
}

//sql Stat2SQL class
type Stat2SQL struct{ pool *pgxpool.Pool }

func NewStat2SQL(pool *pgxpool.Pool) (a *Stat2SQL) {
	a = new(Stat2SQL)
	a.pool = pool
	a.CreateTable()
	return
}

//delete item

func (a *Stat2SQL) Conn(f func(conn *pgxpool.Conn) (err error)) (err error) {

	conn, err := a.pool.Acquire(context.Background())
	if err != nil {
		return
	}
	defer conn.Release()
	return f(conn)
}

// parse sql query
func (a *Stat2SQL) TableName() (res string) {
	return "stat2"
}

func (a *Stat2SQL) Count(where ...string) (count int, err error) {

	var q string

	switch len(where) {
	case 0:
		q = "select count(*) from stat2"
	default:
		q = "select count(*) from stat2 where " + strings.Join(where, " ")
	}

	err = a.Conn(func(conn *pgxpool.Conn) (err error) {
		return conn.QueryRow(context.Background(), q).Scan(&count)
	})
	return
}

// Insert struct and return int id
func (a *Stat2SQL) Insert(v *Stat2) (id int, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "insert into stat2 (uid, type, name, source, ip, image) values ($1, $2, $3, $4, $5, $6) returning id"
	err = conn.QueryRow(c, q, v.sqlTuple()...).Scan(&id)
	return
}

// Insert struct and return int id
func (a *Stat2SQL) InsertNoConflict(v *Stat2) (id int, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "insert into stat2 (uid, type, name, source, ip, image) values ($1, $2, $3, $4, $5, $6) on conflict do nothing returning id"
	err = conn.QueryRow(c, q, v.sqlTuple()...).Scan(&id)
	return
}

// Insert full struct ID must be (ignore all skips)
func (a *Stat2SQL) InsertFull(v *Stat2) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := "insert into stat2 (id, uid, type, name, source, ip, image) values ($1, $2, $3, $4, $5, $6, $7)"
	_, err = conn.Exec(c, q, v.sqlAllTuples()...)
	return
}

// update serial counter (count all items, add 1 and plus custom int)
// useful if you insert with ID
func (a *Stat2SQL) ReindexIDSerialCounter(add ...int) (err error) {
	return a.Conn(func(conn *pgxpool.Conn) (err error) {
		plus := 1
		if len(add) > 0 {
			plus = add[0]
		}
		_, err = conn.Exec(context.Background(), "SELECT setval(pg_get_serial_sequence('stat2', 'id'), COALESCE((SELECT MAX(id) FROM stat2), 0) + $1, false)", plus)
		return
	})
}

// set serial counter
func (a *Stat2SQL) SetIDSerialCounter(value int) (err error) {
	return a.Conn(func(conn *pgxpool.Conn) (err error) {
		_, err = conn.Exec(context.Background(), "SELECT setval(pg_get_serial_sequence('stat2', 'id'), $1, false)", value)
		return
	})
}

// parse sql query
func (a *Stat2SQL) Get(id any, fields ...Stat2IndexType) (res *Stat2, err error) {

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
	q = q + "from stat2 where id = $1 limit 1"
	res = new(Stat2)
	rows, err := conn.Query(c, q, id)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("not found")
		return
	}
	v, err := rows.Values()
	if err != nil {
		return
	}
	if len(v) != len(fields) {
		err = errors.New("len")
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), v[pos])
	}
	err = rows.Err()
	return
}

// parse sql query
func (a *Stat2SQL) GetWhere(where string, fields ...Stat2IndexType) (res *Stat2, err error) {

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
	q = q + " from stat2"
	q = q + fmt.Sprintf(" where %s ", where)
	q = q + " limit 1"
	res = new(Stat2)
	rows, err := conn.Query(c, q)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("not found")
		return
	}
	v, err := rows.Values()
	if err != nil {
		return
	}
	if len(v) != len(fields) {
		err = errors.New("len")
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), v[pos])
	}
	err = rows.Err()
	return
}

// parse sql query
func (a *Stat2SQL) Row(eq map[string]any, fields ...Stat2IndexType) (res *Stat2, err error) {

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
		if !Stat2ValidKey(k) {
			err = errors.New(k)
			return
		}
		in := Stat2KeyIndex(k)
		count++
		andlist = append(andlist, fmt.Sprintf("%s = $%d", in.SQLName(), count))
		values = append(values, v)
	}
	if len(andlist) == 0 {
		err = errors.New("emptykeys")
		return
	}
	keys := strings.Join(andlist, " and ")

	q := fmt.Sprintf("select %s from %s where %s limit 1", fieldlist, a.TableName(), keys)
	res = new(Stat2)
	rows, err := conn.Query(c, q, values...)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("not found")
		return
	}
	v, err := rows.Values()
	if err != nil {
		return
	}
	if len(v) != len(fields) {
		err = errors.New("len")
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), v[pos])
	}
	err = rows.Err()
	return
}

// get simple list from db
func (a *Stat2SQL) All(fields ...Stat2IndexType) (res []*Stat2, err error) {

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
	q += " from stat2"
	rows, err := conn.Query(c, q)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item Stat2
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
	err = rows.Err()

	return
}

// get simple list from db
func (a *Stat2SQL) List(limit, offset int, fields ...Stat2IndexType) (res []*Stat2, err error) {

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
	q += " from stat2 limit $1 offset $2"
	rows, err := conn.Query(c, q, limit, offset)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item Stat2
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
	err = rows.Err()

	return
}

// update sql query
func (a *Stat2SQL) Update(id any, k string, v any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if !Stat2ValidKey(k) {
		return errors.New("invalid key")
	}
	in := Stat2KeyIndex(k)
	q := fmt.Sprintf("update stat2 set %s = $1 where id = $2", in.SQLName())
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, v, id)
	return
}

// update sql query
func (a *Stat2SQL) UpdateWhere(k string, v any, where string) (err error) {
	if !Stat2ValidKey(k) {
		return errors.New("invalid key")
	}
	in := Stat2KeyIndex(k)
	return a.Conn(func(conn *pgxpool.Conn) (err error) {
		q := fmt.Sprintf("update %s set %s = $1 where %s", a.TableName(), in.SQLName(), where)
		_, err = conn.Exec(context.Background(), q, v)
		return
	})
}

// update sql query
func (a *Stat2SQL) Updates(id any, keys map[string]any, where ...string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if keys == nil {
		return errors.New("emptykeys")
	}

	var fields []string
	var values []any
	var count int
	for k, v := range keys {
		if !Stat2ValidKey(k) {
			return errors.New(k)
		}
		in := Stat2KeyIndex(k)
		count++
		fields = append(fields, fmt.Sprintf("%s = $%d", in.SQLName(), count))
		values = append(values, v)
	}

	list := strings.Join(fields, ", ")
	count++
	values = append(values, id)

	q := fmt.Sprintf("update stat2 set %s where id = $%d", list, count)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	_, err = conn.Exec(c, q, values...)
	return
}

// update sql query
func (a *Stat2SQL) UpdatesWhere(keys map[string]any, where string, args ...any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	if keys == nil {
		return errors.New("emptykeys")
	}

	var fields []string
	var values []any
	var count int
	for k, v := range keys {
		if !Stat2ValidKey(k) {
			return errors.New(k)
		}
		in := Stat2KeyIndex(k)
		count++
		fields = append(fields, fmt.Sprintf("%s = $%d", in.SQLName(), count))
		values = append(values, v)
	}

	list := strings.Join(fields, ", ")
	count++

	where = fmt.Sprintf(where, args...)
	q := fmt.Sprintf("update stat2 set %s where %s", list, where)
	_, err = conn.Exec(c, q, values...)
	return
}

type Stat2Query struct {
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

func (a *Stat2Query) Render() (sql string, fields []Stat2IndexType, values []any) {

	switch len(a.Fields) == 0 {
	case true:
		fields = []Stat2IndexType{IndexStat2ID, IndexStat2UID, IndexStat2Type, IndexStat2Name, IndexStat2Source, IndexStat2IP, IndexStat2Image}
	default:
		for _, x := range a.Fields {
			if !Stat2ValidKey(x) {
				continue
			}
			fields = append(fields, Stat2KeyIndex(x))
		}
	}

	var fieldsStrings []string
	for _, x := range fields {
		fieldsStrings = append(fieldsStrings, x.SQLName())
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
	list = append(list, "from stat2")

	var andlist []string

	// EQ where
	for k, v := range a.EQ {
		if !Stat2ValidKey(k) {
			continue
		}
		p := Stat2KeyIndex(k)
		sqlName := p.SQLName()
		switch p.Type() {
		case "bool":
			switch cast.Bool(v) {
			case true:
				andlist = append(andlist, sqlName)
			case false:
				andlist = append(andlist, fmt.Sprintf("(not %[1]s or %[1]s is null)", sqlName))
			}

		default:
			count++
			andlist = append(andlist, fmt.Sprintf("%s = $%d", sqlName, count))
			values = append(values, v)
		}
	}

	// GT where
	for k, v := range a.GT {

		if !Stat2ValidKey(k) {
			continue
		}
		sqlName := Stat2KeyIndex(k).SQLName()
		count++
		andlist = append(andlist, fmt.Sprintf("%s > $%d", sqlName, count))
		values = append(values, v)
	}

	// LT where
	for k, v := range a.LT {
		if !Stat2ValidKey(k) {
			continue
		}
		sqlName := Stat2KeyIndex(k).SQLName()
		count++
		andlist = append(andlist, fmt.Sprintf("%s < $%d", sqlName, count))
		values = append(values, v)
	}

	// NOT where
	for k, v := range a.NOT {
		if !Stat2ValidKey(k) {
			continue
		}
		sqlName := Stat2KeyIndex(k).SQLName()

		for _, x := range v {
			count++
			andlist = append(andlist, fmt.Sprintf("%s != $%d", sqlName, count))
			values = append(values, x)
		}

	}

	// LIKE where
	for k, v := range a.Like {
		if !Stat2ValidKey(k) {
			continue
		}
		sqlName := Stat2KeyIndex(k).SQLName()
		count++
		andlist = append(andlist, fmt.Sprintf("%s ilike $%d", sqlName, count))
		values = append(values, "%"+v+"%")
	}

	// IN where
	if a.IN != nil {
		for k, v := range a.IN {
			if !Stat2ValidKey(k) {
				continue
			}
			sqlName := Stat2KeyIndex(k).SQLName()
			var inlist []string
			for _, num := range v {
				count++
				inlist = append(inlist, fmt.Sprintf("$%d", count))
				values = append(values, num)
			}
			if len(inlist) > 0 {
				andlist = append(andlist, fmt.Sprintf("%s in (%s)", sqlName, strings.Join(inlist, ",")))
			}
		}
	}

	// INS where
	if a.INS != nil {
		for k, v := range a.INS {
			if !Stat2ValidKey(k) {
				continue
			}
			sqlName := Stat2KeyIndex(k).SQLName()
			var inlist []string
			for _, str := range v {
				count++
				inlist = append(inlist, fmt.Sprintf("$%d", count))
				values = append(values, str)
			}
			if len(inlist) > 0 {
				andlist = append(andlist, fmt.Sprintf("%s in (%s)", sqlName, strings.Join(inlist, ",")))
			}
		}
	}

	// INSQ where
	if a.INSQL != nil {
		for k, v := range a.INSQL {
			if !Stat2ValidKey(k) {
				continue
			}
			sqlName := Stat2KeyIndex(k).SQLName()
			count++
			andlist = append(andlist, fmt.Sprintf("%s in (%s)", sqlName, v))
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
	if a.Sort != "" && Stat2ValidKey(a.Sort) {
		list = append(list, fmt.Sprintf("order by %s", Stat2KeyIndex(a.Sort).SQLName()))
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

func NewStat2Query() *Stat2Query {
	a := new(Stat2Query)
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

func (a *Stat2SQL) Search(q *Stat2Query) (res []*Stat2, err error) {

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
		var item Stat2
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
	err = rows.Err()

	return
}

// select custom sql query
func (a *Stat2SQL) Select(where string, fields ...Stat2IndexType) (res []*Stat2, err error) {

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
	q = q + "from stat2"
	if where != "" {
		q += " where " + where
	}
	rows, err := conn.Query(c, q)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var item Stat2
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
	err = rows.Err()
	return
}

// delete item
func (a *Stat2SQL) Delete(id any) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	_, err = conn.Exec(c, "delete from stat2 where id = $1", id)
	return

}

// delete item where: i = 1 and w = 'nice'
func (a *Stat2SQL) DeleteWhere(where string) (err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	_, err = conn.Exec(c, fmt.Sprintf("delete from stat2 where %s", where))
	return

}

// has value in db
func (a *Stat2SQL) Has(field Stat2IndexType, v any) (has bool, err error) {

	c := context.Background()
	conn, err := a.pool.Acquire(c)
	if err != nil {
		return
	}
	defer conn.Release()

	q := fmt.Sprintf("select exists (select 1 from stat2 where %s = $1 limit 1)", field.SQLName())
	err = conn.QueryRow(c, q, v).Scan(&has)
	return
}

// Create table
func (a *Stat2SQL) CreateTable() (err error) {
	return a.Conn(func(conn *pgxpool.Conn) (err error) {
		q := `create table if not exists stat2 (
	id         bigserial primary key,
	uid        bigint,
	type       bigint,
	name       text,
	source     text,
	ip         text,
	image      boolean
)
`
		_, err = conn.Exec(context.Background(), q)
		return
	})
}

// parse sql query
func (a *Stat2SQL) Fields() (res []Stat2IndexType) {
	return []Stat2IndexType{IndexStat2ID, IndexStat2UID, IndexStat2Type, IndexStat2Name, IndexStat2Source, IndexStat2IP, IndexStat2Image}
} //clickhouse Stat2CQL class
type Stat2CQL struct {
	conn driver.Conn
	list []Stat2
	sync.Mutex
}

func NewStat2CQL(conn driver.Conn) (a *Stat2CQL) {
	a = new(Stat2CQL)
	a.conn = conn
	a.CreateTable()
	go a.deamon()
	return
}

// parse clickhouse query
func (a *Stat2CQL) TableName() (res string) {
	return "stat2"
}

func (a *Stat2CQL) Count(where ...string) (count int, err error) {

	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	var q string
	switch len(where) {
	case 0:
		q = "select count() from stat2"
	default:
		q = "select count() from stat2 where " + strings.Join(where, " ")
	}
	var n uint64
	err = a.conn.QueryRow(context.Background(), q).Scan(&n)
	count = int(n)
	return
}

// Add items to clickhouse queue
func (a *Stat2CQL) Add(items ...Stat2) {
	a.Lock()
	defer a.Unlock()
	a.list = append(a.list, items...)
}

// Push queued items to clickhouse
func (a *Stat2CQL) Push() (err error) {
	a.Lock()
	defer a.Unlock()
	if len(a.list) == 0 {
		return
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	batch, err := a.conn.PrepareBatch(context.Background(), "insert into stat2 (id, uid, type, name, source, ip, image)")
	if err != nil {
		return
	}
	for _, item := range a.list {
		err = batch.Append(item.clickhouseTuple()...)
		if err != nil {
			_ = batch.Abort()
			return
		}
	}
	err = batch.Send()
	if err != nil {
		_ = batch.Abort()
		return
	}
	a.list = a.list[:0]
	return
}

// deamon pushes queued items every 10 minutes
func (a *Stat2CQL) deamon() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		a.Lock()
		hasItems := len(a.list) > 0
		a.Unlock()
		if hasItems {
			_ = a.Push()
		}
	}
}

// parse clickhouse query
func (a *Stat2CQL) Get(id any, fields ...Stat2IndexType) (res *Stat2, err error) {

	if fields == nil {
		fields = a.Fields()
	}
	var list []string
	for _, x := range fields {
		list = append(list, x.ClickhouseName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s from stat2 where uid = ? limit 1", fieldlist)
	res = new(Stat2)

	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	rows, err := a.conn.Query(context.Background(), q, id)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("not found")
		return
	}
	rowValues := make([]any, len(fields))
	scan := make([]any, len(fields))
	for i := range rowValues {
		scan[i] = &rowValues[i]
	}
	err = rows.Scan(scan...)
	if err != nil {
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), rowValues[pos])
	}
	err = rows.Err()

	return
}

// parse clickhouse query
func (a *Stat2CQL) GetWhere(where string, fields ...Stat2IndexType) (res *Stat2, err error) {

	if fields == nil {
		fields = a.Fields()
	}
	var list []string
	for _, x := range fields {
		list = append(list, x.ClickhouseName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s from stat2 where %s limit 1", fieldlist, where)
	res = new(Stat2)

	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	rows, err := a.conn.Query(context.Background(), q)
	if err != nil {
		return
	}
	defer rows.Close()
	if !rows.Next() {
		err = errors.New("not found")
		return
	}
	rowValues := make([]any, len(fields))
	scan := make([]any, len(fields))
	for i := range rowValues {
		scan[i] = &rowValues[i]
	}
	err = rows.Scan(scan...)
	if err != nil {
		return
	}
	for pos, x := range fields {
		res.Update(x.String(), rowValues[pos])
	}
	err = rows.Err()

	return
}

// update clickhouse query
func (a *Stat2CQL) Update(id any, k string, v any, where ...string) (err error) {
	if !Stat2ValidKey(k) {
		return errors.New("invalid key")
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	in := Stat2KeyIndex(k)
	q := fmt.Sprintf("alter table stat2 update %s = ? where uid = ?", in.ClickhouseName())
	if len(where) > 0 {
		q += " and " + where[0]
	}
	return a.conn.Exec(context.Background(), q, v, id)
}

// update clickhouse query
func (a *Stat2CQL) UpdateWhere(k string, v any, where string) (err error) {
	if !Stat2ValidKey(k) {
		return errors.New("invalid key")
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	in := Stat2KeyIndex(k)
	q := fmt.Sprintf("alter table %s update %s = ? where %s", a.TableName(), in.ClickhouseName(), where)
	return a.conn.Exec(context.Background(), q, v)
}

// update clickhouse query
func (a *Stat2CQL) Updates(id any, keys map[string]any, where ...string) (err error) {

	if keys == nil {
		return errors.New("emptykeys")
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	var fields []string
	var values []any
	for k, v := range keys {
		if !Stat2ValidKey(k) {
			return errors.New(k)
		}
		in := Stat2KeyIndex(k)
		fields = append(fields, fmt.Sprintf("%s = ?", in.ClickhouseName()))
		values = append(values, v)
	}
	list := strings.Join(fields, ", ")
	values = append(values, id)

	q := fmt.Sprintf("alter table stat2 update %s where uid = ?", list)
	if len(where) > 0 {
		q += " and " + where[0]
	}
	return a.conn.Exec(context.Background(), q, values...)
}

// update clickhouse query
func (a *Stat2CQL) UpdatesWhere(keys map[string]any, where string, args ...any) (err error) {

	if keys == nil {
		return errors.New("emptykeys")
	}
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	var fields []string
	var values []any
	for k, v := range keys {
		if !Stat2ValidKey(k) {
			return errors.New(k)
		}
		in := Stat2KeyIndex(k)
		fields = append(fields, fmt.Sprintf("%s = ?", in.ClickhouseName()))
		values = append(values, v)
	}
	list := strings.Join(fields, ", ")

	where = fmt.Sprintf(where, args...)
	q := fmt.Sprintf("alter table stat2 update %s where %s", list, where)
	return a.conn.Exec(context.Background(), q, values...)
}

// select custom clickhouse query
func (a *Stat2CQL) Select(where string, fields ...Stat2IndexType) (res []*Stat2, err error) {

	if fields == nil {
		fields = a.Fields()
	}
	var list []string
	for _, x := range fields {
		list = append(list, x.ClickhouseName())
	}
	fieldlist := strings.Join(list, ", ")

	q := fmt.Sprintf("select %s from stat2", fieldlist)
	if where != "" {
		q += " where " + where
	}

	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	rows, err := a.conn.Query(context.Background(), q)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var item Stat2
		values := make([]any, len(fields))
		scan := make([]any, len(fields))
		for i := range values {
			scan[i] = &values[i]
		}
		err = rows.Scan(scan...)
		if err != nil {
			continue
		}
		for pos, x := range fields {
			item.Update(x.String(), values[pos])
		}
		res = append(res, &item)
	}
	err = rows.Err()

	return
}

// delete clickhouse item
func (a *Stat2CQL) Delete(id any) (err error) {
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	q := "alter table stat2 delete where uid = ?"
	return a.conn.Exec(context.Background(), q, id)
}

// delete clickhouse items
func (a *Stat2CQL) DeleteWhere(where string) (err error) {
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	q := "alter table stat2 delete where " + where
	return a.conn.Exec(context.Background(), q)
}

// has value in clickhouse
func (a *Stat2CQL) Has(field Stat2IndexType, v any) (has bool, err error) {
	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	q := fmt.Sprintf("select count() from stat2 where %s = ? limit 1", field.ClickhouseName())
	var n uint64
	err = a.conn.QueryRow(context.Background(), q, v).Scan(&n)
	has = n > 0
	return
}

// stat by clickhouse field
func (a *Stat2CQL) Stat(field Stat2IndexType) (res map[string]int, err error) {
	res = make(map[string]int)
	if a.conn == nil {
		err = errors.New("clickhouse connection is nil")
		return
	}
	name := field.ClickhouseName()
	q := fmt.Sprintf("select toString(%[1]s), count() from %[2]s group by %[1]s", name, a.TableName())
	rows, err := a.conn.Query(context.Background(), q)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var k string
		var n uint64
		err = rows.Scan(&k, &n)
		if err != nil {
			return
		}
		res[k] = int(n)
	}
	err = rows.Err()
	return
}

// Create clickhouse table
func (a *Stat2CQL) CreateTable() (err error) {
	if a.conn == nil {
		return errors.New("clickhouse connection is nil")
	}
	q := `create table if not exists stat2 (
	id         Int64,
	uid        Int64,
	type       LowCardinality(Int64),
	name       String,
	source     String,
	ip         IPv4,
	image      Bool
)
engine = MergeTree()
order by (uid)
primary key (uid)`
	return a.conn.Exec(context.Background(), q)
}

// parse clickhouse query
func (a *Stat2CQL) Fields() (res []Stat2IndexType) {
	return []Stat2IndexType{IndexStat2ID, IndexStat2UID, IndexStat2Type, IndexStat2Name, IndexStat2Source, IndexStat2IP, IndexStat2Image}
}

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Stat2Proto struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Uid           int64                  `protobuf:"varint,2,opt,name=uid,proto3" json:"uid,omitempty"`
	Type          int64                  `protobuf:"varint,3,opt,name=type,proto3" json:"type,omitempty"`
	Name          string                 `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Source        string                 `protobuf:"bytes,5,opt,name=source,proto3" json:"source,omitempty"`
	Ip            string                 `protobuf:"bytes,6,opt,name=ip,proto3" json:"ip,omitempty"`
	Image         bool                   `protobuf:"varint,7,opt,name=image,proto3" json:"image,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Stat2Proto) Reset() {
	*x = Stat2Proto{}
	mi := &file_test_ch2_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Stat2Proto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stat2Proto) ProtoMessage() {}

func (x *Stat2Proto) ProtoReflect() protoreflect.Message {
	mi := &file_test_ch2_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stat2Proto.ProtoReflect.Descriptor instead.
func (*Stat2Proto) Descriptor() ([]byte, []int) {
	return file_test_ch2_proto_rawDescGZIP(), []int{0}
}

func (x *Stat2Proto) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Stat2Proto) GetUid() int64 {
	if x != nil {
		return x.Uid
	}
	return 0
}

func (x *Stat2Proto) GetType() int64 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *Stat2Proto) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Stat2Proto) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *Stat2Proto) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *Stat2Proto) GetImage() bool {
	if x != nil {
		return x.Image
	}
	return false
}

var File_test_ch2_proto protoreflect.FileDescriptor

const file_test_ch2_proto_rawDesc = "" +
	"\n" +
	"\x0etest_ch2.proto\x12\x04news\"\x94\x01\n" +
	"\n" +
	"Stat2Proto\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x10\n" +
	"\x03uid\x18\x02 \x01(\x03R\x03uid\x12\x12\n" +
	"\x04type\x18\x03 \x01(\x03R\x04type\x12\x12\n" +
	"\x04name\x18\x04 \x01(\tR\x04name\x12\x16\n" +
	"\x06source\x18\x05 \x01(\tR\x06source\x12\x0e\n" +
	"\x02ip\x18\x06 \x01(\tR\x02ip\x12\x14\n" +
	"\x05image\x18\a \x01(\bR\x05imageB\tZ\a./;newsb\x06proto3"

var (
	file_test_ch2_proto_rawDescOnce sync.Once
	file_test_ch2_proto_rawDescData []byte
)

func file_test_ch2_proto_rawDescGZIP() []byte {
	file_test_ch2_proto_rawDescOnce.Do(func() {
		file_test_ch2_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_test_ch2_proto_rawDesc), len(file_test_ch2_proto_rawDesc)))
	})
	return file_test_ch2_proto_rawDescData
}

var file_test_ch2_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_test_ch2_proto_goTypes = []any{
	(*Stat2Proto)(nil), // 0: news.Stat2Proto
}
var file_test_ch2_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

var _ = initStat2Proto()

func initStat2Proto() struct{} {
	file_test_ch2_proto_init()
	return struct{}{}
}
func file_test_ch2_proto_init() {
	if File_test_ch2_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_test_ch2_proto_rawDesc), len(file_test_ch2_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_test_ch2_proto_goTypes,
		DependencyIndexes: file_test_ch2_proto_depIdxs,
		MessageInfos:      file_test_ch2_proto_msgTypes,
	}.Build()
	File_test_ch2_proto = out.File
	file_test_ch2_proto_goTypes = nil
	file_test_ch2_proto_depIdxs = nil
}

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonC97316b8DecodeGithubComMonopollyJsonsgeneratorTest(in *jlexer.Lexer, out *Stat2) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		switch key {
		case "id":
			if in.IsNull() {
				in.Skip()
			} else {
				out.ID = int(in.Int())
			}
		case "uid":
			if in.IsNull() {
				in.Skip()
			} else {
				out.UID = int(in.Int())
			}
		case "type":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Type = int(in.Int())
			}
		case "name":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Name = string(in.String())
			}
		case "source":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Source = string(in.String())
			}
		case "ip":
			if in.IsNull() {
				in.Skip()
			} else {
				out.IP = string(in.String())
			}
		case "image":
			if in.IsNull() {
				in.Skip()
			} else {
				out.Image = bool(in.Bool())
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonC97316b8EncodeGithubComMonopollyJsonsgeneratorTest(out *jwriter.Writer, in Stat2) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != 0 {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Int(int(in.ID))
	}
	if in.UID != 0 {
		const prefix string = ",\"uid\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.UID))
	}
	if in.Type != 0 {
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Type))
	}
	if in.Name != "" {
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	if in.Source != "" {
		const prefix string = ",\"source\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Source))
	}
	if in.IP != "" {
		const prefix string = ",\"ip\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.IP))
	}
	if in.Image {
		const prefix string = ",\"image\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Bool(bool(in.Image))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Stat2) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonC97316b8EncodeGithubComMonopollyJsonsgeneratorTest(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Stat2) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonC97316b8EncodeGithubComMonopollyJsonsgeneratorTest(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Stat2) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonC97316b8DecodeGithubComMonopollyJsonsgeneratorTest(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Stat2) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonC97316b8DecodeGithubComMonopollyJsonsgeneratorTest(l, v)
}

// easyjson marshal
func (a *Stat2) Marshal() []byte {
	b, _ := easyjson.Marshal(a)
	return b
}

// easyjson unmarshal
func Stat2Unmarshal(src []byte) (a *Stat2) {
	a = new(Stat2)
	err := easyjson.Unmarshal(src, a)
	if err != nil {
		return nil
	}
	return
}
