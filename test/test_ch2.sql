--drop table stat2;
create table if not exists stat2 (
	id         bigserial primary key,
	uid        bigint,
	type       bigint,
	name       text,
	source     text,
	ip         text,
	image      boolean
);
