-- requires createdb kuda_ai

SET ROLE NONE;

-- - - - - - - - - - - - - - - - - - - - - - - - - - -
-- Roles: Groups:

-- developer:

CREATE ROLE developer WITH nologin;

GRANT SELECT, INSERT, UPDATE, DELETE 
ON ALL TABLES IN SCHEMA public
TO developer;

GRANT CREATE ON SCHEMA public TO developer;

-- app:

CREATE ROLE app WITH nologin;

GRANT SELECT, INSERT, UPDATE, DELETE 
ON ALL TABLES IN SCHEMA public
TO app;


-- Roles: Users:

CREATE ROLE dev WITH login PASSWORD 'pa55word' INHERIT;
GRANT developer TO dev;

CREATE ROLE kuda_ai WITH login PASSWORD 'pa55word' INHERIT;
GRANT app TO kuda_ai;


ALTER DATABASE kuda_ai OWNER TO kuda_ai;
