package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type articleRepo struct {
	db *sqlx.DB
}

const (
	createArticle = `
	INSERT INTO "articles" (title, content, tags)
	VALUES ($1, $2, $3) RETURNING *`
)

func NewArticleRepository(database Database) ArticleRepository {
	return &articleRepo{db: database.GetDB()}
}

func (repo *articleRepo) CreateArticle(ctx context.Context, article *Article) (*Article, error) {
	newArticle := &Article{}

	row := repo.db.QueryRowxContext(ctx, createArticle,
		article.Title,
		article.Content,
		article.Tags,
	)

	err := row.StructScan(newArticle)
	if err != nil {
		return nil, err
	}

	return newArticle, nil
}
