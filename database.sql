/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE test (
	id serial PRIMARY KEY,
	name VARCHAR ( 50 ) UNIQUE NOT NULL
);

INSERT INTO test (name) VALUES ('test1');
INSERT INTO test (name) VALUES ('test2');

CREATE TABLE users (
  id serial primary key,
  phone_number VARCHAR(13) UNIQUE NOT NULL,
  full_name VARCHAR(60) NOT NULL,
  password VARCHAR(256) NOT NULL,
  total_login int not null default 0,
  created_at timestamptz default now(),
  updated_at timestamptz,
  updated_by int
);

create index user_phone_number on users using hash(phone_number);
