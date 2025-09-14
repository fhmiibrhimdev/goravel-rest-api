package requests

import (
    "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/contracts/validation"
)

type RegisterRequest struct {
    Name            string `form:"name" json:"name"`
    Email           string `form:"email" json:"email"`
    Password        string `form:"password" json:"password"`
    PasswordConfirm string `form:"password_confirm" json:"password_confirm"`
}

func (r *RegisterRequest) Authorize(ctx http.Context) error {
    return nil
}

func (r *RegisterRequest) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "name":             "required|string|min:3",
        "email":            "required|email",
        "password":         "required|string|min:6",
        "password_confirm": "required|same:password",
    }
}

func (r *RegisterRequest) Messages(ctx http.Context) map[string]string {
    return map[string]string{
        "name.required":             "Nama wajib diisi",
        "email.required":            "Email wajib diisi",
        "password.required":         "Password wajib diisi",
        "password_confirm.required": "Konfirmasi password wajib diisi",
    }
}

func (r *RegisterRequest) Attributes(ctx http.Context) map[string]string {
    return map[string]string{}
}

func (r *RegisterRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
    return nil
}