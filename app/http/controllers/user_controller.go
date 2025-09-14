package controllers

import (
	// "net/http"
	contractsHttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
    "goravel/app/models"
)

type UserController struct {
	// Dependent services
}

func NewUserController() *UserController {
	return &UserController{
		// Inject services
	}
}

func (r *UserController) Show(ctx contractsHttp.Context) contractsHttp.Response {
	return ctx.Response().Success().Json(contractsHttp.Json{
		"Hello": "Goravel",
	})
}

func (r *UserController) SetActive(ctx contractsHttp.Context) contractsHttp.Response {
    var req struct {
        Active int `json:"active"`
    }
    if err := ctx.Request().Bind(&req); err != nil {
        return ctx.Response().Json(400, contractsHttp.Json{"message": "Invalid request"})
    }

    id := ctx.Request().Input("id")
    var user models.User
    if err := facades.Orm().Query().Where("id", id).First(&user); err != nil {
        return ctx.Response().Json(404, contractsHttp.Json{"message": "User not found"})
    }

    user.Active = req.Active
    if err := facades.Orm().Query().Save(&user); err != nil {
        return ctx.Response().Json(500, contractsHttp.Json{"message": "Failed to update status"})
    }

    return ctx.Response().Json(200, contractsHttp.Json{"message": "User status updated", "user": user})
}