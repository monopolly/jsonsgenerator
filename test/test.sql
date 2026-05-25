--drop table news;
create table if not exists news (
	inc                        bigserial,
	ints                       bigint,
	ints8                      smallint,
	ints16                     int,
	ints32                     int,
	ints64                     bigint,
	uints                      bigint,
	uints8                     smallint,
	uints16                    jsonb default '{}'::jsonb,
	uints32                    int,
	uints64                    bigint,
	floats32                   double precision,
	floats64                   double precision,
	bools                      boolean,
	byte1                      smallint,
	bytes                      bytea,
	list_ints                  jsonb default '[]'::jsonb,
	list_string                jsonb default '[]'::jsonb,
	list_float                 jsonb default '[]'::jsonb,
	map_string_string          jsonb default '{}'::jsonb,
	map_string_bytes           jsonb default '{}'::jsonb,
	map_string_bool            jsonb default '{}'::jsonb,
	map_string_int             jsonb default '{}'::jsonb,
	map_string_float64         jsonb default '{}'::jsonb,
	map_string_any             jsonb default '{}'::jsonb,
	map_int_string             jsonb default '{}'::jsonb,
	map_int_int                jsonb default '{}'::jsonb,
	map_int_bool               jsonb default '{}'::jsonb,
	renameSQL_OK               text,
	renameGO                   text,
	renameJS                   text,
	mast_upper_go              text,
	intToSmallInt              smallint,
	sql_unique_u1_1            bigint,
	sql_unique_u1_2            bigint,
	sql_index1_1               bigint,
	sql_index1_2               bigint,
	sql_index1_3               bigint,
	sql_keys_1                 bigint,
	sql_keys_2                 bigint,
	sql_keys_3                 bigint,
	sql_search                 text,
	sql_get                    text,
	sql_unique_x1              bigint,
	sql_unique_x2              bigint,
	sql_unique_x1_x2           bigint,
	sql_primary                double precision,
	sql_jsonb_index            jsonb default '{}'::jsonb,
	time_duration              bigint,
	go_type_int_to_strings     jsonb default '[]'::jsonb,
	public_field1              bigint,
	public_field2              bigint,
	public_field3              bigint,
	public_field_me1           bigint,
	public_field_me2           bigint,
	public_field_me3           bigint,
	search                     tsvector generated always as (to_tsvector('simple', sql_search)) stored,
	unique(sql_unique_x1,sql_unique_x1_x2),
	unique(sql_unique_x1_x2,sql_unique_x2),
	unique(sql_unique_u1_1,sql_unique_u1_2),
	primary key (sql_primary)
);

--generated indexes
--drop index idx_news_sql_index1_1;
create index if not exists idx_news_sql_index1_1 on news (sql_index1_1);

--drop index idx_news_sql_index1_2;
create index if not exists idx_news_sql_index1_2 on news (sql_index1_2);

--drop index idx_news_sql_index1_3;
create index if not exists idx_news_sql_index1_3 on news (sql_index1_3);

--drop index gin_news_sql_jsonb_index;
--ex: select * from news where (sql_jsonb_index->>'year')::int >= 2020;
create index if not exists gin_news_sql_jsonb_index on news using gin(sql_jsonb_index);

--drop index gin_news_search;
--ex: select * from news where search @@ to_tsquery('f8');
create index if not exists gin_news_search on news using gin(search);

--drop index idx_news_index1;
create index if not exists idx_news_index1 on news (sql_index1_1,sql_index1_2,sql_index1_3);
