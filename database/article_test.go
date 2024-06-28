package database

import (
	"context"
	"slices"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestCreateArticle(t *testing.T) {
	db, closeFn := initTestDB(t)
	defer closeFn()

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

func TestGetArticles(t *testing.T) {
	db, closeFn := initTestDB(t)
	defer closeFn()

	repo := NewArticleRepository(db)
	articles := []Article{
		{
			Title:   "Machine Learning in 20 minutes",
			Content: "Learn machine learning in 20 mins",
			Tags:    Tags{"machine learning", "ai"},
		},
		{
			Title: "The art of cooking",
			Content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
			Quisque bibendum posuere dolor, euismod venenatis sem condimentum ut.
			Proin maximus tincidunt auctor. Quisque nunc urna, mollis.`,
			Tags: Tags{"food", "cooking"},
		},
		{
			Title: "Golang for dummies",
			Content: `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
			Quisque bibendum posuere dolor, euismod venenatis sem condimentum ut.
			Proin maximus tincidunt auctor. Quisque nunc urna, mollis.`,
			Tags: Tags{"programming", "golang", "go"},
		},
		{
			Title: "Go Error Handling Explained",
			Content: `Sed ipsum ante, dapibus eu placerat in, rhoncus vitae nisi.
			 Nam at pharetra enim. Aenean sollicitudin sed ante in vestibulum.`,
			Tags: Tags{"golang", "go"},
		},
		{
			Title: "How to make a slushy",
			Content: `Duis consectetur laoreet nisl, in malesuada lorem tristique sit amet.
			Cras in tortor porttitor, mollis ligula nec, elementum tellus. Maecenas arcu ante.`,
			Tags: Tags{"slushy", "food"},
		},
	}

	for _, article := range articles {
		_, err := repo.CreateArticle(context.Background(), &article)
		require.NoError(t, err)
	}

	t.Run("fetch with no filter", func(t *testing.T) {
		foundArticles, err := repo.GetArticles(context.Background(), ArticleFilter{})
		require.NoError(t, err)
		require.Len(t, foundArticles, len(articles))
	})

	t.Run("filter by multiple tags", func(t *testing.T) {
		foundArticles, err := repo.GetArticles(context.Background(), ArticleFilter{Tags: Tags{"golang", "food"}})
		require.NoError(t, err)
		require.Len(t, foundArticles, 4)

		for _, article := range foundArticles {
			require.True(t, slices.Contains(article.Tags, "golang") || slices.Contains(article.Tags, "food"))
		}
	})

	t.Run("filter by empty tags", func(t *testing.T) {
		foundArticles, err := repo.GetArticles(context.Background(), ArticleFilter{Tags: Tags{}})
		require.NoError(t, err)
		require.Len(t, foundArticles, len(articles))
	})
}

func TestGetArticleByID(t *testing.T) {
	// get test db
	db, closeFn := initTestDB(t)
	defer closeFn()

	// create repo
	repo := NewArticleRepository(db)

	t.Run("article successfully found", func(t *testing.T) {
		payload := Article{
			Title:   "How to write code efficiently",
			Content: "lorem ipsum and stuff",
			Tags:    Tags{"coding"},
		}

		// seed article
		createdArticle, err := repo.CreateArticle(context.Background(), &payload)
		require.NoError(t, err)

		foundArticle, err := repo.GetArticleByID(context.Background(), createdArticle.ID)
		require.NoError(t, err)
		require.Equal(t, createdArticle, foundArticle)
	})

	t.Run("article not found", func(t *testing.T) {
		nonExistentID := 0

		foundArticle, err := repo.GetArticleByID(context.Background(), nonExistentID)

		require.Nil(t, foundArticle)
		require.ErrorIs(t, err, ErrArticleNotFound)
	})
}

func TestUpdateArticle(t *testing.T) {
	// get test db
	db, closeFn := initTestDB(t)
	defer closeFn()

	// create repo
	repo := NewArticleRepository(db)
	newArticle := Article{
		Title:   "How to write code efficiently",
		Content: "lorem ipsum and stuff",
		Tags:    Tags{"coding"},
	}

	// seed article
	createdArticle, err := repo.CreateArticle(context.Background(), &newArticle)
	require.NoError(t, err)

	// update content
	createdArticle.Title = "How to cook oats"
	createdArticle.Content = "Just do it"
	createdArticle.Tags = Tags{"oats", "cooking", "food"}

	updatedArticle, err := repo.UpdateArticle(context.Background(), createdArticle)
	require.NoError(t, err)
	require.Equal(t, createdArticle.Title, updatedArticle.Title)
	require.Equal(t, createdArticle.Content, updatedArticle.Content)
	require.Equal(t, createdArticle.Tags, updatedArticle.Tags)
	require.NotEqual(t, createdArticle.UpdatedAt, updatedArticle.UpdatedAt)

}
