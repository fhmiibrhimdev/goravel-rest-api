package requests

import (
    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/contracts/validation"
)

type LoginRequest struct {
    Email    string `form:"email" json:"email"`
    Password string `form:"password" json:"password"`
}

func (r *LoginRequest) Authorize(ctx http.Context) error {
    return nil
}

func (r *LoginRequest) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "email":    "required|email",
        "password": "required|string",
    }
}

func (r *LoginRequest) Messages(ctx http.Context) map[string]string {
    return map[string]string{
        "email.required":    "Email wajib diisi",
        "email.email":       "Format email tidak valid",
        "password.required": "Password wajib diisi",
    }
}

func (r *LoginRequest) Attributes(ctx http.Context) map[string]string {
    return map[string]string{}
}

func (r *LoginRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
    return nil
}