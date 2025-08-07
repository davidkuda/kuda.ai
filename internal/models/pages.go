package models

import (
	"database/sql"
	"fmt"
	"html/template"
	"time"

	"github.com/russross/blackfriday/v2"
)

type Pages []*Page

type Page struct {
	ID          int
	Name        string
	Version     int
	Content     string // Markdown
	HTMLContent template.HTML
	CreatedAt   time.Time
}

type PageModel struct {
	DB *sql.DB
}

func (m *PageModel) Insert(page *Page) error {
	var err error

	// get version first if exists:
	var version int

	vstmt := `
	SELECT COALESCE(MAX(version), 0)
	FROM website.pages
	WHERE name = $1
	`

	row := m.DB.QueryRow(vstmt, page.Name)
	err = row.Scan(&version)
	if err != nil {
		if err == sql.ErrNoRows {
			// continue
		} else {
			return fmt.Errorf("row.Scan(&page.Version): %v", err)
		}
	}
	page.Version = version + 1

	stmt := `
	INSERT INTO website.pages (name, version, content)
	VALUES ($1, $2, $3);
	`

	_, err = m.DB.Exec(
		stmt,
		page.Name,
		page.Version,
		page.Content,
	)

	if err != nil {
		return fmt.Errorf("failed executing insert sql: %v", err)
	}

	return nil
}

func (m *PageModel) Get(name string) (*Page, error) {
	var err error

	stmt := `
	SELECT id, name, version, content, created_at
	FROM website.pages
	WHERE name = $1
	ORDER BY version DESC
	LIMIT 1;
	`

	row := m.DB.QueryRow(stmt, name)

	page := Page{}

	err = row.Scan(
		&page.ID,
		&page.Name,
		&page.Version,
		&page.Content,
		&page.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	htmlBytes := blackfriday.Run([]byte(page.Content))
	page.HTMLContent = template.HTML(htmlBytes)

	return &page, nil
}
