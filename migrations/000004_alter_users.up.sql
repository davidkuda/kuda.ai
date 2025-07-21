-- before:
-- CREATE TABLE auth.users (
--     email VARCHAR(255) PRIMARY KEY,
--     hashed_password CHAR(60) NOT NULL,
--     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
-- );

-- after:
-- create table auth.users (
-- 	id              SERIAL primary key,
-- 	email           VARCHAR(255) unique,
-- 	hashed_password CHAR(60) not null,
-- 	created_at      TIMESTAMPTZ(0) default now() not null
-- );

BEGIN;

-- 1. Drop the primary key on email
ALTER TABLE auth.users DROP CONSTRAINT users_pkey;

-- 2. Add the new id column
ALTER TABLE auth.users ADD COLUMN id SERIAL;

-- 3. Make id the primary key
ALTER TABLE auth.users ADD PRIMARY KEY (id);

-- 4. Ensure emails remain unique
ALTER TABLE auth.users ADD CONSTRAINT users_email_key UNIQUE (email);

COMMIT;

