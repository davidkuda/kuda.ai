-- prerequisite: createdb kuda_ai

create schema website;
create schema auth;


SET ROLE NONE;

-- - - - - - - - - - - - - - - - - - - - - - - - - - -
-- Roles: Groups:

-- developer:

CREATE ROLE developer WITH nologin;

-- how developer is different to app: CREATE vs USAGE:
GRANT CREATE ON SCHEMA website, auth TO developer;

GRANT SELECT, INSERT, UPDATE, DELETE 
ON ALL TABLES IN SCHEMA website, auth
TO developer;

-- app:

CREATE ROLE app WITH nologin;

-- how developer is different to app: CREATE vs USAGE:
GRANT USAGE ON SCHEMA auth TO app;

GRANT SELECT, INSERT, UPDATE, DELETE 
ON ALL TABLES IN SCHEMA website, auth
TO app;


-- Roles: Users:

CREATE ROLE dev WITH login PASSWORD 'pa55word' INHERIT;
GRANT developer TO dev;

CREATE ROLE kuda_ai WITH login PASSWORD 'pa55word' INHERIT;
GRANT app TO kuda_ai;


ALTER DATABASE kuda_ai OWNER TO dev;
