create schema songbook;

create table if not exists songbook.songs (
id text primary key unique,
    artist text,
    name text,
    lyrics text,
    chords text,
    copyright text
);

