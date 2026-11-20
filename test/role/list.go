package role

/*
Constants Generator
Help to auto-generated correct constants for lists
Martin Prestone (c) 2024
github.com/monopolly

Struct options (first // line)
name=Status   Named type for constants: type Status int. Package funcs get type prefix: StatusIndex(), StatusValid(), so many lists can live in one package
type=string   Constants type: int (default), int8..int64, uint..uint64, string
#simple       Constants only, no functions
#up           Uppercase golang names: NEW instead of New
#index        Additional int constants StatusIndexNew, StatusIndexActive... (from 1)
#swift        Create .swift enum next to .go file

Field options (after // on field line)
_             Skip value (reserves iota position)
iota=100      Start new const block from this value (int types)
name{}        Replace string value. Ex: art // name{artwork}
go{}          Replace golang const name. Ex: art // go{ArtList}
title{}       Title() for field. Ex: title{Real Estate}. Golang name if not set
desc{}        Desc() for field. Ex: desc{Asset is on sale}. Desc() is generated only if any field has desc{}
list{}        List() []string for field. Ex: list{geoid, country}
#tag          Group by tag. Ex: #virtual > IsVirtual() bool, Virtual() bool, StatusVirtualIndexes()
@tag          Same as #tag
key{int}      Custom func. Ex: code{401} > Code() int
key{text}     Custom func. Ex: label{Nice code} > Label() string
key{bool}     Custom func. Ex: paid{true} > Paid() bool
@key=value    Same as key{value}, value without spaces. Ex: @code=401

Custom funcs: key is capitalized (code{} > Code()), type is int if all values are numbers,
bool if all values are true/false, otherwise string. Missing values return 0 / "" / false.

Generated functions
(Role) String() string           String value
(Role) Title() string            title{}, golang name if not set
(Role) Int() int                 Value as int
RoleValidIndexType(v Role) bool  Is value known
(Role) Valid() bool              Is value known
RoleValid(v int) bool            Is raw int a known value
RoleIndexes() []Role             All values
RoleIndex(v string) int          String value to int, 0 if unknown
*/
type Role int

const (
	Owner = Role(iota) + 1 //iota=1
	Admin
	Manager
)

func (a Role) String() string {
	switch a {
	case Owner:
		return "owner"
	case Admin:
		return "admin"
	case Manager:
		return "manager"
	default:
		return ""
	}
}

func (a Role) Title() string {
	switch a {
	case Owner:
		return "Owner"
	case Admin:
		return "Admin"
	case Manager:
		return "Manager"
	default:
		return ""
	}
}

func (a Role) Int() int {
	return int(a)
}

func RoleValidIndexType(v Role) bool {
	switch v {
	case Owner, Admin, Manager:
		return true
	default:
		return false
	}
}

func (a Role) Valid() bool {
	return RoleValidIndexType(a)
}

func RoleValid(v int) bool {
	switch Role(v) {
	case Owner, Admin, Manager:
		return true
	default:
		return false
	}
}

func RoleIndexes() []Role {
	return []Role{Owner, Admin, Manager}
}

func RoleIndex(v string) int {
	switch v {
	case "owner":
		return int(Owner)
	case "admin":
		return int(Admin)
	case "manager":
		return int(Manager)
	default:
		return 0
	}
}
