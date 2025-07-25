BEGIN;

SET ROLE dev;

-- TODO: Do I really need to not null?? especially with default 0?
-- NOTE: Prices are stored in Rappen, not a fraction of CHF.
--       (in Go, post-process: CHF => float64(total_price) / 100.0)
create table website.bellevue_activities (
	id              SERIAL primary key,
	user_id         INT references auth.users(id),
	activity_date   DATE not null,
	breakfast_count INT default 0 not null,
	lunch_count     INT default 0 not null,
	dinner_count    INT default 0 not null,
	coffee_count    INT default 0 not null,
	sauna_count     INT default 0 not null,
	lecture_count   INT default 0 not null,
	snacks_chf      INT default 0 not null,
	comment         TEXT,
	total_price     INT not null,
	created_at      TIMESTAMPTZ default now() not null
);

GRANT SELECT, INSERT, UPDATE, DELETE
ON ALL TABLES IN SCHEMA website
TO app;

GRANT USAGE, SELECT ON SEQUENCE website.bellevue_activities_id_seq TO app;

ALTER TABLE website.bellevue_activities OWNER TO dev;

COMMIT;
