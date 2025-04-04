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
	stmt := `
	INSERT INTO website.til (path, title, category, summary, text)
	VALUES ($1, $2, $3, $4, $5);
	`

	_, err := m.DB.Exec(
		stmt,
		t.Path,
		t.Title,
		t.Category,
		t.Summary,
		t.Text,
	)
	if err != nil {
		log.Printf("failed executing insert sql: %v", err)
	}

	return nil
}

func (m *TILModel) UpdateExisting(t *TIL) error {
	stmt := `
	UPDATE website.til
	SET path = $2,
		title = $3,
		category = $4
		summary = $5,
		text = $6,
		updated_at = CURRENT_DATE
	WHERE id = $1;
	`

	_, err := m.DB.Exec(
		stmt,
		t.ID,
		t.Path,
		t.Title,
		t.Category,
		t.Summary,
		t.Text,
	)
	if err != nil {
		log.Printf("failed executing insert sql: %v", err)
	}

	return nil
}

func (m *TILModel) GetBy(TILPath string) (*TIL, error) {
	stmt := `
	SELECT id, path, title, category, summary, text, created_at, updated_at
	FROM til
	WHERE path = $1;
	`

	row := m.DB.QueryRow(stmt, TILPath)

	til := TIL{}

	err := row.Scan(
		&til.ID,
		&til.Path,
		&til.Title,
		&til.Category,
		&til.Summary,
		&til.Text,
		&til.CreatedAt,
		&til.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &til, nil
}
