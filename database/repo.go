package database

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type ArticleFilter struct {
	Tags        Tags
	PublishedAt time.Time
}

type Tags []string

func (t Tags) Value() (driver.Value, error) {
	if t == nil {
		return "[]", nil
	}

	return json.Marshal(t)
}

func (t *Tags) Scan(value interface{}) error {
	if value == nil {
		*t = Tags{}
		return nil
	}

	b, ok := value.([]byte)
	if !ok {
		return errors.New("unexpected value from driver")
	}

	return json.Unmarshal(b, t)
}

type Article struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	Tags        Tags      `json:"tags" db:"tags"`
	PublishedAt time.Time `json:"published_at" db:"published_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ArticleRepository interface {
	// GetArticles(filter ArticleFilter) ([]Article, error)
	// GetArticleByID(ID int) (*Article, error)
	CreateArticle(ctx context.Context, article *Article) (*Article, error)
	// UpdateArticle(article *Article) (*Article, error)
	// DeleteArticle(ID int) error
}
