#!/bin/bash
set -e

PGPASSWORD=$POSTGRESQL_PASSWORD psql -v ON_ERROR_STOP=1 --username "$POSTGRESQL_USERNAME" --dbname "$POSTGRESQL_DATABASE" <<-EOSQL
	CREATE TABLE url (
  long_url varchar NOT NULL,
  ip varchar ,
  is_enable int4 NOT NULL,
  reg_date timestamp NOT NULL,
  url_id serial4 NOT NULL,
  CONSTRAINT url_pkey PRIMARY KEY (url_id)
  );
  CREATE INDEX idx_url_long_url ON url USING btree (long_url);
  ALTER SEQUENCE url_url_id_seq restart with 1;

  CREATE TABLE contents (
  content_id serial4 NOT NULL,
  url_id serial4,
  long_url varchar NOT NULL,
  size int4 ,
  type varchar ,
  hash text NOT NULL,
  is_enable int4 NOT NULL,
  reg_date timestamp NOT NULL,
  CONSTRAINT contents_pkey PRIMARY KEY (content_id)
  );
EOSQL