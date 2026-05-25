create table if not exists stat (
	id          Int64,
	uid         Int64,
	type        LowCardinality(Int64),
	name        String,
	source      String,
	ip          IPv4,
	image       Bool,
	stack       Array(String),
	uints       Array(UInt64),
	uints8      Array(UInt8),
	uints16     Array(UInt16),
	uints32     Array(UInt32),
	uints64     Array(UInt64),
	ints8       Array(Int8),
	ints16      Array(Int16),
	ints32      Array(Int32),
	ints64      Array(Int64)
)
engine = MergeTree()
order by (uid)
primary key (uid);
