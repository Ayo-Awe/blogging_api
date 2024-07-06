package database

import (
	"context"
	"math"
	"slices"
	"testing"

	"github.com/ayo-awe/blogging_api/utils"
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
		foundArticles, _, err := repo.GetArticles(context.Background(), ArticleFilter{}, Paging{
			Page:    1,
			PerPage: 5,
		})
		require.NoError(t, err)
		require.Len(t, foundArticles, len(articles))
	})

	t.Run("filter by multiple tags", func(t *testing.T) {
		foundArticles, _, err := repo.GetArticles(context.Background(), ArticleFilter{Tags: Tags{"golang", "food"}}, Paging{
			Page:    1,
			PerPage: 20,
		})
		require.NoError(t, err)
		require.Len(t, foundArticles, 4)

		for _, article := range foundArticles {
			require.True(t, slices.Contains(article.Tags, "golang") || slices.Contains(article.Tags, "food"))
		}
	})

	t.Run("filter by empty tags", func(t *testing.T) {
		foundArticles, _, err := repo.GetArticles(context.Background(), ArticleFilter{Tags: Tags{}}, Paging{
			Page:    1,
			PerPage: 20,
		})
		require.NoError(t, err)
		require.Len(t, foundArticles, len(articles))
	})

	paginationTestCases := []struct {
		name               string
		expectedTotalItems int
		filter             ArticleFilter
		perPage            int
	}{
		{
			name:               "page with no filter",
			expectedTotalItems: len(articles),
			perPage:            2,
		},
		{
			name:               "page with filter",
			expectedTotalItems: 2,
			filter:             ArticleFilter{Tags{"go"}},
			perPage:            2,
		},
	}

	for _, tc := range paginationTestCases {
		t.Run(tc.name, func(t *testing.T) {
			perPage := 2
			expectedTotalPages := math.Ceil(float64(tc.expectedTotalItems) / float64(perPage))

			currentPage := 1
			for currentPage <= int(expectedTotalPages) {
				paging := Paging{
					Page:    currentPage,
					PerPage: perPage,
				}

				foundArticles, paginationData, err := repo.GetArticles(context.Background(), tc.filter, paging)
				require.NoError(t, err)

				require.Equal(t, paging.PerPage, paginationData.PerPage)
				require.Equal(t, paging.Page, paginationData.CurrentPage)
				require.Equal(t, len(foundArticles), paginationData.ItemCount)

				expectedItemCount := utils.ClampInt(tc.expectedTotalItems-paging.PerPage*(currentPage-1), 0, paging.PerPage)
				require.Equal(t, expectedItemCount, paginationData.ItemCount)

				require.Equal(t, paginationData.TotalPages, int(expectedTotalPages))
				require.Equal(t, paginationData.TotalItems, tc.expectedTotalItems)
				currentPage++
			}
		})
	}

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

func TestDeleteArticle(t *testing.T) {
	db, closeFn := initTestDB(t)
	defer closeFn()

	repo := NewArticleRepository(db)

	payload := Article{
		Title:   "How to bake bread",
		Content: "Just do it",
		Tags:    Tags{"baking"},
	}

	article, err := repo.CreateArticle(context.Background(), &payload)
	require.NoError(t, err)

	require.NoError(t, repo.DeleteArticle(context.Background(), article.ID))

	_, err = repo.GetArticleByID(context.Background(), article.ID)
	require.ErrorIs(t, err, ErrArticleNotFound)
}
