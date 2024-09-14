create schema songbook;

create table if not exists songbook.songs (
    id text primary key unique,
    artist text,
    name text,
    lyrics text,
    chords text,
    copyright text
);

CREATE TABLE users (
    email VARCHAR(255) PRIMARY KEY,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);
