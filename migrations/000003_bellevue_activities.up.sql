BEGIN;

SET ROLE dev;

-- NOTE: Prices are stored in Rappen, not a fraction of CHF.
--       (CHF => float64(total_price) / 100.0)
create table website.bellevue_activities (
	id                 SERIAL primary key,
	user_id            INT references auth.users(id),
	activity_date      DATE not null,
	breakfast_count    INT default 0 not null,
	lunch_dinner_count INT default 0 not null,
	coffee_count       INT default 0 not null,
	sauna_count        INT default 0 not null,
	lecture_count      INT default 0 not null,
	comment            TEXT,
	total_price        INT,
	created_at         TIMESTAMPTZ default now() not null
);

GRANT SELECT, INSERT, UPDATE, DELETE
ON ALL TABLES IN SCHEMA website
TO app;

GRANT USAGE, SELECT ON SEQUENCE website.bellevue_activities_id_seq TO app;

ALTER TABLE website.bellevue_activities OWNER TO dev;

COMMIT;
