## SQL

### Dynamic, normalized bellevue activites

```sql
BEGIN;

SET ROLE dev;

drop table if exists bellevue.activities;
drop table if exists bellevue.activity_kinds;
drop table if exists bellevue.prices;
drop table if exists bellevue.activity_items;


create schema if not exists bellevue;

-- - - - - - - - - - - - - - - - - - - - - - - - - - -
-- ACTIVITIES:
CREATE TABLE bellevue.activities (
	id                     SERIAL  primary key,
	user_id                INTEGER not null references auth.users(id),
	activity_date          DATE    not null,
	comment                TEXT,
	snacks_and_drinks_chf  INTEGER, -- in Rappen
	created_at             TIMESTAMPTZ not null default now()
);

CREATE INDEX idx_activity_user_date
ON bellevue.activities (user_id, activity_date DESC);


INSERT INTO bellevue.activities (
	user_id,
	activity_date
) VALUES (
	1,
	now()
);

SELECT * FROM bellevue.activities;

-- - - - - - - - - - - - - - - - - - - - - - - - - - -
-- ACTIVITIY KINDS: enums:
CREATE TABLE bellevue.activity_kinds (
	code         TEXT primary key,
	display_name TEXT not null,
	price        INTEGER not null
);

-- alter table bellevue.activity_kinds add column price INTEGER;
-- update bellevue.activity_kinds set price = 800;
-- select * from bellevue.activity_kinds;
INSERT INTO bellevue.activity_kinds (
	code, display_name
) VALUES (
	'breakfast', 'Breakfast', 800
);

-- dimensions to activities
CREATE TABLE bellevue.activity_items (
	id           SERIAL   primary key,
	activity_id  INTEGER  not null references bellevue.activities(id) on delete cascade,
	kind         TEXT     not null references bellevue.activity_kinds(code),
	quantity     INTEGER  not null,
	unit_price   INTEGER  not null -- Rappen at booking time
);

INSERT INTO bellevue.activity_items (
	activity_id,
	kind,
	quantity,
	unit_price
) VALUES (
	1,
	'breakfast',
	1,
	SELECT price
	FROM bellevue.activity_kinds
	WHERE code = 'breakfast'
);

CREATE TABLE bellevue.prices (
	id          SERIAL  primary key,
	service_id  INTEGER not null references services(service_id) on delete cascade,
	amount      INTEGER not null,
	valid_from  DATE    not null,
	valid_to    DATE    DEFAULT 'infinity'::date
);

COMMIT;
```


### generated column that calculates all prices

```sql
/* generated column; integer = SUM(qty × unit_price) in centimes */
CREATE TABLE bellevue.activities (
	id            SERIAL primary key,
	user_id       INTEGER not null references auth.users(id),
	activity_date DATE not null,
	comment       TEXT,
	snacks_and_drinks_chf INTEGER, -- in Rappen
	created_at    TIMESTAMPTZ not null default now(),
	total_price   INTEGER generated always as (
		COALESCE((
			SELECT SUM(i.qty * i.unit_price)
			FROM   website.activity_items i
			WHERE  i.activity_id = id
		), 0)
	) STORED
);

```

### Numeric for Currency:

```sql
amount NUMERIC(10,2) not null,
```


### Checks:

```sql
qty INTEGER not null check (qty > 0),
```



## JavaScript

### HTMX scroll after swap:

```js
document.body.addEventListener('htmx:afterSwap', () => {
  window.scrollTo({ top: 0});
});

document.body.addEventListener("htmx:afterSwap", () => {
	const main = document.querySelector("main");
	if (main) {
		main.scrollIntoView();
	}
});

document.body.addEventListener("htmx:afterSwap", (e) => {
	const main = document.querySelector("main");
	// Only scroll when the swap target is <main>
	if (e.detail.target === main) {
		main.scrollIntoView();
	}
});
```
