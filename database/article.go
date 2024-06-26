package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type articleRepo struct {
	db *sqlx.DB
}

var (
	ErrArticleNotFound = errors.New("article not found")
)

const (
	createArticle = `
	INSERT INTO "articles" (title, content, tags)
	VALUES ($1, $2, $3) RETURNING *`

	getArticles = `
	SELECT
		id,
		title,
		content,
		tags,
		published_at,
		updated_at
	FROM "articles"
	WHERE tags ?| $1 OR $1 = '{}' OR $1 IS NULL
	ORDER BY published_at DESC;`

	getArticleByID = `
	SELECT
		id,
		title,
		content,
		tags,
		published_at,
		updated_at
	FROM "articles"
	WHERE id = $1;`
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

func (repo *articleRepo) GetArticles(ctx context.Context, filter ArticleFilter) ([]Article, error) {
	rows, err := repo.db.QueryxContext(ctx, getArticles, pq.Array(filter.Tags))
	if err != nil {
		return nil, err
	}

	var articles []Article
	for rows.Next() {
		var article Article
		err := rows.StructScan(&article)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func (repo *articleRepo) GetArticleByID(ctx context.Context, ID int) (*Article, error) {
	var article Article

	err := repo.db.QueryRowxContext(ctx, getArticleByID, ID).StructScan(&article)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrArticleNotFound
		}
		return nil, err
	}

	return &article, nil
}
