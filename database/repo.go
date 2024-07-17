package database

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ArticleFilter struct {
	Tags Tags
}

type PaginationData struct {
	CurrentPage int `json:"current_page" example:"1"`
	TotalPages  int `json:"total_pages" example:"2"`
	ItemCount   int `json:"item_count" example:"25"`
	TotalItems  int `json:"total_items" example:"40"`
	PerPage     int `json:"per_page" example:"25"`
}

func (p *PaginationData) Build(paging Paging, itemCount, totalItems int) {
	p.CurrentPage = paging.Page
	p.PerPage = paging.PerPage
	p.ItemCount = itemCount
	p.TotalItems = totalItems
	p.TotalPages = int(math.Ceil(float64(totalItems) / float64(paging.PerPage)))
}

type Paging struct {
	Page    int
	PerPage int
}

func (p Paging) Limit() int {
	return p.PerPage
}

func (p Paging) Offset() int {
	return (p.Page - 1) * p.PerPage
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
	ID          int       `json:"id" db:"id" example:"1"`
	Title       string    `json:"title" db:"title" example:"I love Golang"`
	Content     string    `json:"content" db:"content" example:"lorem ipsum lorem ipsum"`
	Tags        Tags      `json:"tags" db:"tags" example:"golang,go,tech"`
	PublishedAt time.Time `json:"published_at" db:"published_at" example:"2024-06-23T22:21:19.00199+01:00"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" example:"2024-06-23T22:21:19.00199+01:00"`
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
	GetArticles(ctx context.Context, filter ArticleFilter, pageable Paging) ([]Article, PaginationData, error)
	GetArticleByID(ctx context.Context, ID int) (*Article, error)
	CreateArticle(ctx context.Context, article *Article) (*Article, error)
	UpdateArticle(ctx context.Context, article *Article) (*Article, error)
	DeleteArticle(ctx context.Context, ID int) error
}
