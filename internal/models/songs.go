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
	return nil, nil
}

func (m *SongModel) Insert(s *Song) error {
	stmt := `insert into songbook.songs (
				id, artist, name, lyrics, chords, copyright
			) VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := m.DB.Exec(stmt, s.ID, s.Artist, s.Name, s.Lyrics, s.Chords, s.Copyright)
	if err != nil {
		log.Printf("failed executing insert sql: %v", err)
	}

	return nil
}

func (m *SongModel) Get(songID string) (Song, error) {
	stmt := `select artist, name, lyrics, chords
	from songbook.songs
	where id = $1;`

	row := m.DB.QueryRow(stmt, songID)

	s := Song{}

	err := row.Scan(&s.Artist, &s.Name, &s.Lyrics, &s.Chords)
	if err != nil {
		return Song{}, err
	}

	return s, nil
}

func (m *SongModel) Latest() ([]Song, error) {
	return nil, nil
}
