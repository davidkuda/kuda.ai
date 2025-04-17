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
	Path        string
	Version     int
	Title       string
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
	WHERE path = $1
	`

	row := m.DB.QueryRow(vstmt, page.Path)
	err = row.Scan(&version)
	if err != nil {
		return fmt.Errorf("row.Scan(&page.Version): %v", err)
	}
	page.Version = version + 1

	stmt := `
	INSERT INTO website.pages (path, version, title, content)
	VALUES ($1, $2, $3, $4);
	`

	_, err = m.DB.Exec(
		stmt,
		page.Path,
		page.Version,
		page.Title,
		page.Content,
	)

	if err != nil {
		return fmt.Errorf("failed executing insert sql: %v", err)
	}

	return nil
}

func (m *PageModel) GetByPath(path string) (*Page, error) {
	var err error

	stmt := `
	SELECT *
	FROM website.pages
	WHERE path = $1
	ORDER BY version DESC
	LIMIT 1;
	`

	row := m.DB.QueryRow(stmt, path)

	page := Page{}

	err = row.Scan(
		&page.Path,
		&page.Version,
		&page.Title,
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
