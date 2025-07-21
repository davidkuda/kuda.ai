BEGIN;

-- 1. Drop the primary key on id
ALTER TABLE auth.users DROP CONSTRAINT users_pkey;

-- 2. Drop the unique constraint on email
ALTER TABLE auth.users DROP CONSTRAINT IF EXISTS users_email_key;

-- 3. Drop the id column
ALTER TABLE auth.users DROP COLUMN id;

-- 4. Make email the primary key again
ALTER TABLE auth.users ADD PRIMARY KEY (email);

COMMIT;

