package controllers

import (
	// "fmt"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
)

type PostController struct{}

func NewPostController() *PostController {
	return &PostController{}
}

// List all posts
func (r *PostController) Index(ctx http.Context) http.Response {
    var posts []models.Post
    if err := facades.Orm().Query().Find(&posts); err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "status":  "error",
            "message": "Something went wrong on our server",
        })
    }
    return ctx.Response().Json(http.StatusOK, http.Json{
        "status":  "success",
        "message": "List of data",
        "data":    posts,
    })
}

// Show single post
func (r *PostController) Show(ctx http.Context) http.Response {
    id := ctx.Request().Route("id")
    var post models.Post
    if err := facades.Orm().Query().Where("id", id).First(&post); err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "status":  "error",
            "message": "Something went wrong on our server",
        })
    }
    if post.ID == 0 {
        return ctx.Response().Json(http.StatusNotFound, http.Json{
            "status":  "error",
            "message": "Data not found",
        })
    }
    return ctx.Response().Json(http.StatusOK, http.Json{
        "status":  "success",
        "message": "Show post detail",
        "data":    post,
    })
}

// Create new post with validation
func (r *PostController) Store(ctx http.Context) http.Response {
    var input struct {
        Title       string `json:"title"`
        Body        string `json:"body"`
        PublishDate string `json:"publish_date"`
    }
    if err := ctx.Request().Bind(&input); err != nil {
        return ctx.Response().Json(http.StatusBadRequest, http.Json{
            "status":  "error",
            "message": "Invalid request payload",
        })
    }

    validator, err := facades.Validation().Make(map[string]any{
        "title":        input.Title,
        "body":         input.Body,
        "publish_date": input.PublishDate,
    }, map[string]string{
        "title":        "required",
        "body":         "required",
        "publish_date": "required",
    })
    if err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "status":  "error",
            "message": "Something went wrong on our server",
        })
    }
    if validator.Fails() {
        return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
            "status":  "error",
            "message": "Validation failed",
            "errors":  validator.Errors().All(),
        })
    }

    post := models.Post{
        Title:       input.Title,
        Body:        input.Body,
        PublishDate: input.PublishDate,
    }
    if err := facades.Orm().Query().Create(&post); err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "status":  "error",
            "message": "Something went wrong on our server",
        })
    }
    return ctx.Response().Json(http.StatusCreated, http.Json{
        "status":  "success",
        "message": "Data created successfully",
        "data":    post,
    })
}

// Update post with validation
func (r *PostController) Update(ctx http.Context) http.Response {
    id := ctx.Request().Route("id")
    var post models.Post
    if err := facades.Orm().Query().Where("id", id).First(&post); err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "status":  "error",
            "message": "Something went wrong on our server",
        })
    }
    if post.ID == 0 {
        return ctx.Response().Json(http.StatusNotFound, http.Json{
            "status":  "error",
            "message": "Data not found",
        })
    }

    // Bind input ke struct
    var input struct {
        Title       string `json:"title"`
        Body        string `json:"body"`
        PublishDate string `json:"publish_date"`
    }
    if err := ctx.Request().Bind(&input); err != nil {
        return ctx.Response().Json(http.StatusBadRequest, http.Json{
            "status":  "error",
            "message": "Invalid request payload",
        })
    }

    // Validasi dengan validator
    validator, err := facades.Validation().Make(map[string]any{
        "title":        input.Title,
        "body":         input.Body,
        "publish_date": input.PublishDate,
    }, map[string]string{
        "title":        "",
        "body":         "",
        "publish_date": "",
    })
    if err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "status":  "error",
            "message": "Something went wrong on our server",
        })
    }
    if validator.Fails() {
        return ctx.Response().Json(http.StatusUnprocessableEntity, http.Json{
            "status":  "error",
            "message": "Validation failed",
            "errors":  validator.Errors().All(),
        })
    }

    // Update dengan struct sesuai dokumentasi Goravel
    _, err = facades.Orm().Query().Where("id", id).Update(models.Post{
        Title:       input.Title,
        Body:        input.Body,
        PublishDate: input.PublishDate,
    })
    if err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "status":  "error",
            "message": "Something went wrong on our server",
        })
    }

    // Ambil data terbaru untuk response
    if err := facades.Orm().Query().Where("id", id).First(&post); err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "status":  "error",
            "message": "Something went wrong on our server",
        })
    }

    return ctx.Response().Json(http.StatusOK, http.Json{
        "status":  "success",
        "message": "Data updated successfully",
        "data":    post,
    })
}

// Delete post
func (r *PostController) Delete(ctx http.Context) http.Response {
    id := ctx.Request().Route("id")
    result, err := facades.Orm().Query().Where("id", id).Delete(&models.Post{})
    if err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, http.Json{
            "status":  "error",
            "message": "Something went wrong on our server",
        })
    }
    if result.RowsAffected == 0 {
        return ctx.Response().Json(http.StatusNotFound, http.Json{
            "status":  "error",
            "message": "Data not found",
        })
    }
    return ctx.Response().Json(http.StatusOK, http.Json{
        "status":  "success",
        "message": "Data deleted successfully",
    })
}