package api

import (
	"strings"

	"github.com/ayo-awe/blogging_api/database"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type SuccessReponse struct {
	Status   string      `json:"status" example:"success"`
	Data     interface{} `json:"data"`
	Metadata interface{} `json:"metadata,omitempty" swaggerignore:"true"`
}

type ErrorResponse struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Invalid Request Body"`
}

type CreateArticleResponse struct {
	Article database.Article `json:"article"`
}

type GetArticleByIDResponse struct {
	Article database.Article `json:"article"`
}

type UpdateArticleResponse struct {
	Article database.Article `json:"article"`
}

type GetArticlesResponse struct {
	Articles []database.Article `json:"articles"`
}

func NewSuccessResponse(data interface{}, metadata interface{}) *SuccessReponse {
	return &SuccessReponse{
		Status:   "success",
		Data:     data,
		Metadata: metadata,
	}
}

func NewErrResponse(msg string) *ErrorResponse {
	return &ErrorResponse{
		Status:  "error",
		Message: msg,
	}
}

type CreateArticleRequest struct {
	Title   string        `json:"title" example:"I love Golang"`
	Content string        `json:"content" example:"lorem ipsum lorem ipsum lorem ipsum"`
	Tags    database.Tags `json:"tags" example:"golang,tech"`
}

func (c *CreateArticleRequest) toArticle() *database.Article {
	return &database.Article{
		Title:   c.Title,
		Content: c.Content,
		Tags:    c.Tags,
	}
}

func (c *CreateArticleRequest) Validate() error {
	c.clean()
	return validation.ValidateStruct(c,
		validation.Field(&c.Title, validation.Length(5, 255)),
		validation.Field(&c.Content, validation.Length(5, 0)),
		validation.Field(&c.Tags, validation.Each(validation.Length(2, 0), is.LowerCase)),
	)
}

func (c *CreateArticleRequest) clean() {
	c.Title = strings.TrimSpace(c.Title)
	c.Content = strings.TrimSpace(c.Content)

	for i, tag := range c.Tags {
		trimmed := strings.TrimSpace(tag)
		c.Tags[i] = strings.ToLower(trimmed)
	}
}

type UpdateArticleRequest struct {
	Title   string        `json:"title" example:"I love Golang"`
	Content string        `json:"content" example:"lorem ipsum lorem ipsum lorem ipsum"`
	Tags    database.Tags `json:"tags" example:"golang,tech"`
}

func (u *UpdateArticleRequest) Validate() error {
	u.clean()
	return validation.ValidateStruct(u,
		validation.Field(&u.Title, validation.Length(5, 255)),
		validation.Field(&u.Content, validation.Length(5, 0)),
		validation.Field(&u.Tags, validation.Each(validation.Length(2, 0), is.LowerCase)),
	)
}

func (u *UpdateArticleRequest) clean() {
	u.Title = strings.TrimSpace(u.Title)
	u.Content = strings.TrimSpace(u.Content)

	for i, tag := range u.Tags {
		trimmed := strings.TrimSpace(tag)
		u.Tags[i] = strings.ToLower(trimmed)
	}
}
