package database

import (
	"context"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateArticle(t *testing.T) {
	db, err := NewDatabase("postgresql://aweayo:aweayo@localhost:5432/blogging_api_test?sslmode=disable")
	require.NoError(t, err)

	defer func() {
		_, err := db.GetDB().Exec("TRUNCATE TABLE articles;")
		require.NoError(t, err)
		db.GetDB().Close()
	}()

	repo := NewArticleRepository(db)

	t.Run("create with tags", func(t *testing.T) {
		payload := &Article{
			Title:   "Deep learning for dummies",
			Content: "Learn deep learning",
			Tags:    Tags{"artificial intelligence", "technology"},
		}

		article, err := repo.CreateArticle(context.Background(), payload)
		require.NoError(t, err)

		require.NotEmpty(t, article)
		require.Equal(t, payload.Title, article.Title)
		require.Equal(t, payload.Content, article.Content)
		require.Equal(t, payload.Tags, article.Tags)
	})

	t.Run("create without tags", func(t *testing.T) {
		payload := &Article{
			Title:   "Deep learning for dummies",
			Content: "Learn deep learning",
		}

		article, err := repo.CreateArticle(context.Background(), payload)
		require.NoError(t, err)

		require.NotEmpty(t, article)
		require.Equal(t, payload.Title, article.Title)
		require.Equal(t, payload.Content, article.Content)
		require.NotNil(t, article.Tags)
	})

}
