package models

import (
	"database/sql"
	"time"
	"log"
)

type TILs []*TIL

type TIL struct {
	ID      string
	Date    time.Time
	Title   string
	Teaser  string
	Content string
}

type TILModel struct {
	DB *sql.DB
}

func (m *TILModel) GetAll() (TILs, error) {
	stmt := "select id, date, title, teaser, content from til;"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tils TILs

	for rows.Next() {
		var til *TIL
		err = rows.Scan(
			til.ID,
			til.Date,
			til.Title,
			til.Teaser,
			til.Content,
		)
		if err != nil {
			return nil, err
		}
		tils = append(tils, til)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tils, nil
}

func (m *TILModel) Insert(t *TIL) error {
	stmt := `insert into til (
				id, date, title, teaser, content
			) VALUES (
				$1, $2, $3, $4, $5
			);`

	_, err := m.DB.Exec(
		stmt,
		t.ID,
		t.Date,
		t.Title,
		t.Teaser,
		t.Content,
	)
	if err != nil {
		log.Printf("failed executing insert sql: %v", err)
	}

	return nil
}

func (m *TILModel) Get(tilID string) (*TIL, error) {
	stmt := `select id, date, title, teaser, content
	from til
	where id = $1;`

	row := m.DB.QueryRow(stmt, tilID)

	til := TIL{}

	err := row.Scan(
		til.ID,
		til.Date,
		til.Title,
		til.Teaser,
		til.Content,
	)

	if err != nil {
		return nil, err
	}

	return &til, nil
}
