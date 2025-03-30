package models

import (
	"database/sql"
	"html/template"
	"log"
)

type Songs []Song

type Song struct {
	ID        string
	Artist    string
	Name      string
	Lyrics    string
	Chords    string
	Copyright string
	MyCover   string
	Covers    []string

	HTML struct {
		Lyrics template.HTML
		Chords template.HTML
	}
}

type SongModel struct {
	DB *sql.DB
}

func (m *SongModel) GetAllSongs() (Songs, error) {
	stmt := "select id, artist, name from website.songs order by artist"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var songs Songs

	for rows.Next() {
		var song Song
		err = rows.Scan(&song.ID, &song.Artist, &song.Name)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}

func (m *SongModel) Insert(s *Song) error {
	stmt := `insert into website.songs (
				id, artist, name, lyrics, chords, copyright
			) VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO UPDATE SET
				artist = EXCLUDED.artist,
				name = EXCLUDED.name,
				lyrics = EXCLUDED.lyrics,
				chords = EXCLUDED.chords,
				copyright = EXCLUDED.copyright;`


	_, err := m.DB.Exec(stmt, s.ID, s.Artist, s.Name, s.Lyrics, s.Chords, s.Copyright)
	if err != nil {
		log.Printf("failed executing insert sql: %v", err)
	}

	return nil
}

func (m *SongModel) Get(songID string) (*Song, error) {
	stmt := `select artist, name, lyrics, chords, copyright
	from website.songs
	where id = $1;`

	row := m.DB.QueryRow(stmt, songID)

	s := Song{}

	err := row.Scan(&s.Artist, &s.Name, &s.Lyrics, &s.Chords, &s.Copyright)
	if err != nil {
		return nil, err
	}

	s.ID = songID

	return &s, nil
}

func (m *SongModel) Latest() ([]Song, error) {
	return nil, nil
}
