package database

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ArticleFilter struct {
	Tags Tags
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

func (a *Article) Validate() error {
	a.clean()
	return validation.ValidateStruct(a,
		validation.Field(&a.Title, validation.Required, validation.Length(5, 255)),
		validation.Field(&a.Content, validation.Required, validation.Length(5, 0)),
		validation.Field(&a.Tags, validation.Each(validation.Length(2, 0), is.LowerCase)),
	)
}

func (a *Article) clean() {
	a.Title = strings.TrimSpace(a.Title)
	a.Content = strings.TrimSpace(a.Content)

	for i, tag := range a.Tags {
		trimmed := strings.TrimSpace(tag)
		a.Tags[i] = strings.ToLower(trimmed)
	}
}

type ArticleRepository interface {
	GetArticles(ctx context.Context, filter ArticleFilter) ([]Article, error)
	GetArticleByID(ctx context.Context, ID int) (*Article, error)
	CreateArticle(ctx context.Context, article *Article) (*Article, error)
	UpdateArticle(ctx context.Context, article *Article) (*Article, error)
	DeleteArticle(ctx context.Context, ID int) error
}
