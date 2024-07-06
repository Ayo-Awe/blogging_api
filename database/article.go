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
	ORDER BY published_at DESC
	LIMIT $2
	OFFSET $3;`

	countArticles = `
	SELECT
		count(*)
	FROM "articles"
	WHERE tags ?| $1 OR $1 = '{}' OR $1 IS NULL;`

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

	updateArticle = `
	UPDATE "articles"
	SET
		title = $2,
		content = $3,
		tags = $4,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $1
	RETURNING *;`

	deleteArticle = `
	DELETE FROM "articles"
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

func (repo *articleRepo) GetArticles(ctx context.Context, filter ArticleFilter, paging Paging) ([]Article, PaginationData, error) {
	rows, err := repo.db.QueryxContext(ctx, getArticles, pq.Array(filter.Tags), paging.Limit(), paging.Offset())
	if err != nil {
		return []Article{}, PaginationData{}, err
	}

	articles := []Article{}
	for rows.Next() {
		var article Article
		err := rows.StructScan(&article)
		if err != nil {
			return []Article{}, PaginationData{}, err
		}

		articles = append(articles, article)
	}

	var articleCount int
	if err = repo.db.QueryRowContext(ctx, countArticles, pq.Array(filter.Tags)).Scan(&articleCount); err != nil {
		return []Article{}, PaginationData{}, err
	}

	paginationData := PaginationData{}
	paginationData.Build(paging, len(articles), articleCount)

	return articles, paginationData, nil
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

func (repo *articleRepo) UpdateArticle(ctx context.Context, article *Article) (*Article, error) {
	var updatedArticle Article

	row := repo.db.QueryRowxContext(ctx, updateArticle,
		article.ID,
		article.Title,
		article.Content,
		article.Tags)

	if err := row.StructScan(&updatedArticle); err != nil {
		return nil, err
	}

	return &updatedArticle, nil
}

func (repo *articleRepo) DeleteArticle(ctx context.Context, ID int) error {
	_, err := repo.db.ExecContext(ctx, deleteArticle, ID)
	if err != nil {
		return err
	}

	return nil
}
