package models

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/russross/blackfriday/v2"
)

type Blogs []*Blog

type Blog struct {
	ID      int
	Path    string
	Title   string
	Summary string
	Content string
	HTML    struct {
		Summary template.HTML
		Content template.HTML
	}
	CreatedAt time.Time
	UpdatedAt time.Time
}

type BlogModel struct {
	DB *sql.DB
}

func (m *BlogModel) GetAll() (Blogs, error) {
	stmt := `
	SELECT id, path, title, summary, content, created_at, updated_at
	FROM website.blogs
	ORDER BY created_at DESC;`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var blogs Blogs

	for rows.Next() {
		var blog Blog
		err = rows.Scan(
			&blog.ID,
			&blog.Path,
			&blog.Title,
			&blog.Summary,
			&blog.Content,
			&blog.CreatedAt,
			&blog.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, &blog)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}

func (m *BlogModel) Insert(b *Blog) error {
	stmt := `
	INSERT INTO website.blogs (
		path, title, summary, content
	) VALUES (
		$1, $2, $3, $4
	);
	`

	_, err := m.DB.Exec(
		stmt,
		b.Path,
		b.Title,
		b.Summary,
		b.Content,
	)
	if err != nil {
		return fmt.Errorf("failed executing insert sql: %v", err)
	}

	return nil
}

func (m *BlogModel) UpdateExisting(b *Blog) error {
	stmt := `
	UPDATE website.blogs
	SET path = $2,
		title = $3,
		summary = $4,
		content = $5,
		updated_at = current_timestamp
	WHERE id = $1;
	`

	_, err := m.DB.Exec(
		stmt,
		b.ID,
		b.Path,
		b.Title,
		b.Summary,
		b.Content,
	)
	if err != nil {
		log.Printf("failed executing insert sql: %v", err)
	}

	return nil
}

func (m *BlogModel) GetByPath(blogPath string) (*Blog, error) {
	stmt := `
	SELECT id, path, title, summary, content, created_at, updated_at
	FROM website.blogs
	WHERE path = $1;
	`

	row := m.DB.QueryRow(stmt, blogPath)

	blog := Blog{}

	err := row.Scan(
		&blog.ID,
		&blog.Path,
		&blog.Title,
		&blog.Summary,
		&blog.Content,
		&blog.CreatedAt,
		&blog.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	htmlBytes := blackfriday.Run([]byte(blog.Content))
	blog.HTML.Content = template.HTML(htmlBytes)

	return &blog, nil
}

func (m *BlogModel) PathIsUnique(path string) (bool, error) {
	stmt := `
	SELECT count(path)
	FROM website.blogs
	WHERE path = $1;
	`

	row := m.DB.QueryRow(stmt, path)

	var count int
	err := row.Scan(
		&count,
	)

	if err != nil {
		return false, err
	}

	if count > 0 {
		return false, nil
	}

	return true, nil
}
