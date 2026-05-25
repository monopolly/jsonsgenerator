# jsonsgenerator

`jsonsgenerator` is a beta Go code generator driven by compact struct annotations. It reads a regular `.go` file with a model declaration and generates production-oriented code around that model: JSON helpers, PostgreSQL CRUD, ClickHouse ingestion helpers, protobuf, gotiny, TypeScript, Swift, tests, and benchmarks.

The main idea is simple: describe the model once, then let the generator create the repetitive plumbing around it.

## What It Generates

Depending on the options in the comment above the struct, the generator can create:

- `NameGO.go` - the main Go file with the generated struct, field indexes, JSON helpers, easyjson code, gotiny/protobuf code.
- `NameGO_test.go` - tests and benchmarks for a specific generated struct.
- `Name.sql` - PostgreSQL DDL.
- `Name.cql` - ClickHouse DDL.
- `Name.proto` - protobuf schema.
- `Name.ts`, `Name.js` - TypeScript/JS model.
- `Name.swift`, `NameEnum.swift` - Swift model and enum.
- `Name.json` - demo JSON.

## Installation

```bash
go install github.com/monopolly/jsonsgenerator@latest
```

For local development from this repository:

```bash
go install .
```

Go 1.25+ is required.

The generator uses `easyjson`, so the `easyjson` binary must be available in `PATH`:

```bash
go install github.com/mailru/easyjson/...@latest
```

If you use the `proto` option, you also need `protoc` and the Go protobuf plugin:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

## Quick Start

Create `news.go`:

```go
package model

// go sql js proto test gotiny
type news struct {
	id     int    // sql{inc} #readonly
	title  string
	lang   string // ch{low}
	ip     string
	active bool
	tags   []string
	views  uint64
}
```

Run the generator:

```bash
jsonsgenerator news.go
```

The generated files will be written next to the source file:

```text
newsGO.go
newsGO_test.go
news.sql
news.proto
news.js
```

If you add the `ch` option, the generator will also create `news.cql` and a ClickHouse class.

## Struct Options

Options are written in the comment above the struct:

```go
// go=News sql=news ch=news js=NewsJson proto=NewsProto test gotiny
type news struct {
	// ...
}
```

Short option names are supported too:

```go
// go sql ch js proto
type news struct {
	// ...
}
```

In short form, names are derived from the original struct name:

- `go` -> `News`
- `sql` -> table `news`, class `NewsSQL`
- `ch` -> table `news`, class `NewsCQL`
- `js` -> `NewsJson`
- `proto` -> `NewsProto`

Common struct options:

| Option | Description |
| --- | --- |
| `go` / `go=News` | Generate the main Go type. |
| `sql` / `sql=news` | Generate PostgreSQL class and `.sql` file. |
| `ch` / `ch=news` | Generate ClickHouse class and `.cql` file. |
| `js` / `js=NewsJson` | Generate a JSON wrapper based on `jsons`. |
| `proto` / `proto=NewsProto` | Generate `.proto` and embed protobuf Go code into `GO.go`. |
| `test` | Generate `GO_test.go` with tests and benchmarks. |
| `gotiny` | Add `MarshalGotiny` and `UnmarshalNewsGotiny`. |
| `msgp` | Add MessagePack support. |
| `ts=name` | Generate a TypeScript model. |
| `swift` | Generate a Swift model. |
| `enum=Name` | Generate a Swift enum. |
| `demo` | Generate a demo JSON file. |
| `optimize` | Try to reorder fields to reduce Go struct padding. |
| `noinit` | Do not generate the `New...()` init function. |
| `!omit` | Remove `omitempty` from JSON tags. |
| `noprefix` | Use shorter field index names: `IndexID` instead of `IndexNewsID`. |

## Field Options

Field options are written in the comment after a field:

```go
type news struct {
	id       int            // sql{inc} #readonly
	title    string         // sql{index} json{name="headline"}
	lang     string         // ch{low}
	ip       string         // ch{type=ip}
	payload  map[string]any // sql{type="jsonb"} proto{type="google.protobuf.Value"}
	internal string         // sql{skip} proto{skip} js{skip}
}
```

Common field options:

| Option | Description |
| --- | --- |
| `sql{inc}` | Auto-increment field. Usually `id int`. |
| `sql{index}` | Create an SQL index. |
| `sql{unique="group"}` | Add a unique constraint for a group of fields. |
| `sql{type="jsonb"}` | Override the SQL field type. |
| `sql{skip}` | Exclude the field from SQL generation. |
| `ch{primary}` | Use the field as the ClickHouse key/primary field. |
| `ch{low}` | Wrap the ClickHouse type in `LowCardinality(...)`. Useful for languages, countries, statuses, and similar fields. |
| `ch{type=ip}` | Use ClickHouse `IPv4`. A field named `ip` is also treated as IPv4 automatically. |
| `json{name="..."}` | Override the JSON field name. |
| `json{skip}` | Exclude the field from JSON helpers. |
| `proto{name="..."}` | Override the protobuf field name. |
| `proto{type="..."}` | Override the protobuf field type. |
| `proto{skip}` | Exclude the field from protobuf generation. |
| `go{type="..."}` | Override the generated Go field type. |
| `#readonly`, `#must` | Custom field tags used by generated helpers. |

## Example: PostgreSQL Model

```go
package model

// go=Account sql=accounts js=AccountJson test
type account struct {
	id      int    // sql{inc} #readonly
	email   string // sql{unique="email"}
	name    string // sql{index}
	active  bool
	balance int64
}
```

Generated output:

- `accountGO.go` with `Account`, `AccountSQL`, and `AccountQuery`.
- `accountGO_test.go` with `TestAccountMarshal`, `BenchmarkAccountMarshal`, `BenchmarkAccountUnmarshalStd`, and related tests.
- `account.sql` with table and index definitions.

Usage:

```go
pool, _ := pgxpool.New(ctx, dsn)
db := NewAccountSQL(pool)

item := Account{
	Email:  "hello@example.com",
	Name:   "Sergey",
	Active: true,
}

id, err := db.Insert(&item)
_ = id
_ = err

list, err := db.Search(&AccountQuery{
	Limit: 20,
	EQ: map[string]any{
		"active": true,
	},
})
_ = list
_ = err
```

## Example: ClickHouse Buffer

```go
package model

// go ch=news_events test
type event struct {
	id      uint64
	lang    string // ch{low}
	country string // ch{low}
	ip      string // ch{type=ip}
	path    string
}
```

ClickHouse generation creates `EventCQL`. The class has a local `list` buffer, `Add` for fast in-memory accumulation, and `Push` for batch insertion.

```go
ch := NewEventCQL(conn)

ch.Add(Event{
	ID:      1,
	Lang:    "ru",
	Country: "RU",
	IP:      "127.0.0.1",
	Path:    "/",
})

err := ch.Push()
_ = err
```

`Stat(field)` returns aggregate counts as `map[string]int`:

```go
stats, err := ch.Stat(IndexEventLang)
_ = stats // map[string]int{"ru": 425, "en": 134}
_ = err
```

## Example: Protobuf And Gotiny

```go
package model

// go proto gotiny test
type stat struct {
	id     int
	name   string
	values []uint64
}
```

Generated output:

- `statGO.go` with `Stat`, `StatProto`, protobuf runtime code, and gotiny helpers.
- `stat.proto`.
- `statGO_test.go` with JSON, protobuf, and gotiny tests/benchmarks.

Usage:

```go
item := Stat{ID: 1, Name: "views", Values: []uint64{1, 2, 3}}

jsonBytes := item.Marshal()
fromJSON := StatUnmarshal(jsonBytes)

tinyBytes := item.MarshalGotiny()
fromTiny := UnmarshalStatGotiny(tinyBytes)

_ = fromJSON
_ = fromTiny
```

## Tests And Benchmarks

The `test` option generates tests for the specific struct, so multiple generated models can live in the same package without function name collisions:

```go
func TestStatMarshal(t *testing.T)
func TestStat2Marshal(t *testing.T)
func BenchmarkStatMarshal(b *testing.B)
func BenchmarkStat2Marshal(b *testing.B)
```

Benchmarks use the modern `testing.B.Loop` API:

```go
for b.Loop() {
	_ = item.Marshal()
}
```

## Supported Types

The generator supports common Go scalar types, strings, booleans, numbers, `time.Time`, `time.Duration`, `[]byte`, arrays/slices of strings and numbers, and several `map[...]...` shapes for JSON/SQL/protobuf scenarios.

For ClickHouse, numeric arrays, string arrays, `LowCardinality`, and IPv4 fields are supported. IP fields are always generated as ClickHouse `IPv4`.

## Recommendations

- Keep source models separate from generated files, for example `news.go` and `newsGO.go`.
- Do not edit `*GO.go`, `.sql`, `.cql`, or `.proto` manually. Change the source struct and rerun the generator instead.
- For multiple models in the same package, use separate source files. Test names and protobuf init helpers are generated from the struct name.
- Commit generated files if you want reproducible builds without running the generator in CI.

