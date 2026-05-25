create table if not exists stat2 (
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
primary key (uid);
