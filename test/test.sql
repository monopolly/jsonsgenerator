--drop table news;
create table if not exists news (
	newid                                        bigserial primary key,
	createdNew                                   bigint default extract(epoch from now()),
	count                                        bigint,
	active                                       boolean,
	kyc                                          newtype,
	oid                                          bigint,
	type                                         bigint,
	verify                                       boolean,
	title                                        text,
	html                                         bytea,
	tags                                         jsonb default '[]'::jsonb,
	channels                                     jsonb default '[]'::jsonb,
	channels64                                   jsonb default '[]'::jsonb,
	floats                                       double precision,
	keys                                         jsonb default '{}'::jsonb,
	features                                     jsonb default '{}'::jsonb,
	likes                                        jsonb default '{}'::jsonb,
	providers                                    jsonb default '{}'::jsonb,
	stats                                        jsonb default '{}'::jsonb,
	price                                        jsonb default '{}'::jsonb,
	meta                                         jsonb default '{}'::jsonb,
	timeout                                      bigint,
	value                                        jsonb default '{}'::jsonb,
	rawbytes                                     bytea,
	search                                       tsvector generated always as (to_tsvector('simple', title)) stored,
	unique(oid,html,channels),
	unique(tags,channels),
	primary key (floats,keys,features)
);

--generated indexes
--drop index idx_news_type;
create index if not exists idx_news_type on news (type);

--drop index gin_news_meta;
--ex: select * from news where (meta->>'year')::int >= 2020;
create index if not exists gin_news_meta on news using gin(meta);

--drop index gin_news_search;
--ex: select * from news where search @@ to_tsquery('f8');
create index if not exists gin_news_search on news using gin(search);

--drop index idx_news_i1;
create index if not exists idx_news_i1 on news (oid,type);

--drop index idx_news_i2;
create index if not exists idx_news_i2 on news (type);
