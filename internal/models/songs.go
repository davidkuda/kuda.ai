package models

import (
	"database/sql"
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
}

type SongModel struct {
	DB *sql.DB
}

func (m *SongModel) GetAllSongs() (Songs, error) {
	return nil, nil
}

func (m *SongModel) Insert(s *Song) (int, error) {
	return 0, nil
}

func (m *SongModel) Get(id int) (Song, error) {
	return Song{}, nil
}

func (m *SongModel) Latest() ([]Song, error) {
	return nil, nil
}
