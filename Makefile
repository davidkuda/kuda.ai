PG_DSN_ADMIN = postgres://davidkuda:@${DB_ADDRESS}/${DB_NAME}?sslmode=disable
PG_DSN_APP = postgres://${DB_USER}:${DB_PASSWORD}@${DB_ADDRESS}/${DB_NAME}?sslmode=disable

db/backup:
	pg_dump \
	--data-only \
	--column-inserts \
	--no-privileges \
	--no-owner \
	--table=songs \
	> ./data/postgres/2025-03-28--songs


db/init:
	createdb kuda_ai

db/drop:
	dropdb kuda_ai

db/restore:
	psql -X kuda_ai --single-transaction < ./data/postgres/dumpfile--data-only

db/migrate/newsql:
	migrate create \
	-seq \
	-ext=.sql \
	-dir=./migrations \
	songszzz

db/migrate/up/roles:
	migrate \
	-path=./migrations \
	-database=${PG_DSN_ADMIN} \
	up

db/migrate/version:
	migrate \
	-URL=${PG_DSN_ADMIN} \
	version

user ?= dev
psql/dev:
	psql \
	--host localhost \
	--username ${user} \
	--port 5432 \
	--dbname kuda_ai
