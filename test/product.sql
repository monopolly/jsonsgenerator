--drop table products;
create table if not exists products (
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
);

--altertable
--1
alter table products add column brand text;
--2
alter table products add column model text;
--October
alter table products add column image text;
--rename
alter table products rename column oldname to about;