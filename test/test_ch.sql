--drop table stat;
create table if not exists stat (
	id          bigserial primary key,
	uid         bigint,
	type        bigint,
	name        text,
	source      text,
	ip          text,
	image       boolean,
	stack       jsonb default '[]'::jsonb,
	uints       jsonb default '[]'::jsonb,
	uints8      bytea,
	uints16     jsonb default '[]'::jsonb,
	uints32     jsonb default '[]'::jsonb,
	uints64     jsonb default '[]'::jsonb,
	ints8       bytea,
	ints16      jsonb default '[]'::jsonb,
	ints32      jsonb default '[]'::jsonb,
	ints64      jsonb default '[]'::jsonb
);
