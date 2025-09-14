package controllers

import (
    "net/http"
    // "fmt"
    contractsHttp "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"
    // "goravel/app/http/requests"
    // "strings"
    "errors"
    "github.com/goravel/framework/auth"
    "goravel/app/models"
    "time"
    "regexp"
)

type AuthController struct{}

func NewAuthController() *AuthController {
    return &AuthController{}
}

func (r *AuthController) Test(ctx contractsHttp.Context) contractsHttp.Response {
    return ctx.Response().Json(200, contractsHttp.Json{
        "message": "AuthController is working",
        "timestamp": time.Now(),
    })
}

func (r *AuthController) Register(ctx contractsHttp.Context) contractsHttp.Response {
    // Get input data
    name := ctx.Request().Input("name")
    email := ctx.Request().Input("email")
    password := ctx.Request().Input("password")
    passwordConfirm := ctx.Request().Input("password_confirm")
    
    // Validation errors map
    errors := make(map[string][]string)
    
    // Validate name
    if name == "" {
        errors["name"] = append(errors["name"], "Nama wajib diisi")
    } else if len(name) < 3 {
        errors["name"] = append(errors["name"], "Nama minimal 3 karakter")
    }
    
    // Validate email
    if email == "" {
        errors["email"] = append(errors["email"], "Email wajib diisi")
    } else if !isValidEmail(email) {
        errors["email"] = append(errors["email"], "Format email tidak valid")
    } else {
        // Check email uniqueness dengan proper error handling
        var existingUser models.User
        err := facades.Orm().Query().Where("email", email).First(&existingUser)
        
        // Hanya error jika benar-benar ada user (tidak error "record not found")
        if err == nil && existingUser.ID > 0 {
            errors["email"] = append(errors["email"], "Email sudah terdaftar")
        }
    }
    
    // Validate password
    if password == "" {
        errors["password"] = append(errors["password"], "Password wajib diisi")
    } else if len(password) < 6 {
        errors["password"] = append(errors["password"], "Password minimal 6 karakter")
    }
    
    // Validate password confirmation
    if passwordConfirm == "" {
        errors["password_confirm"] = append(errors["password_confirm"], "Konfirmasi password wajib diisi")
    } else if password != passwordConfirm {
        errors["password_confirm"] = append(errors["password_confirm"], "Konfirmasi password tidak sesuai")
    }
    
    // If there are validation errors
    if len(errors) > 0 {
        return ctx.Response().Json(400, contractsHttp.Json{
            "message": "Validation failed",
            "errors":  errors,
        })
    }
    
    // Create new user
    var user models.User
    user.Name = name
    user.Email = email
    
    // Hash password
    if err := user.HashPassword(password); err != nil {
        return ctx.Response().Json(500, contractsHttp.Json{
            "message": "Failed to hash password",
            "error": err.Error(),
        })
    }
    
    // Save user to database
    if err := facades.Orm().Query().Create(&user); err != nil {
        return ctx.Response().Json(500, contractsHttp.Json{
            "message": "Failed to create user",
            "error": err.Error(),
        })
    }
    
    // Generate JWT token
    token, err := facades.Auth(ctx).Login(&user)
    if err != nil {
        return ctx.Response().Json(500, contractsHttp.Json{
            "message": "Failed to generate token",
            "error": err.Error(),
        })
    }
    
    return ctx.Response().Json(201, contractsHttp.Json{
        "message": "User registered successfully",
        "data": contractsHttp.Json{
            "user":  user,
            "token": token,
        },
    })
}

// Helper function untuk validasi email
func isValidEmail(email string) bool {
    // Simple email regex
    emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    match, _ := regexp.MatchString(emailRegex, email)
    return match
}

func (r *AuthController) Login(ctx contractsHttp.Context) contractsHttp.Response {
    // Get input data
    email := ctx.Request().Input("email")
    password := ctx.Request().Input("password")
    
    // Validation errors map
    errors := make(map[string][]string)
    
    // Validate email
    if email == "" {
        errors["email"] = append(errors["email"], "Email wajib diisi")
    } else if !isValidEmail(email) {
        errors["email"] = append(errors["email"], "Format email tidak valid")
    }
    
    // Validate password
    if password == "" {
        errors["password"] = append(errors["password"], "Password wajib diisi")
    }
    
    // If there are validation errors
    if len(errors) > 0 {
        return ctx.Response().Json(400, contractsHttp.Json{
            "message": "Validation failed",
            "errors":  errors,
        })
    }
    
    // Find user by email
    var user models.User
    err := facades.Orm().Query().Where("email", email).First(&user)
    if err != nil || user.ID == 0 {
        return ctx.Response().Json(401, contractsHttp.Json{
            "message": "Invalid credentials",
        })
    }

    // Tambahkan pengecekan active
    if user.Active != 1 {
        return ctx.Response().Json(403, contractsHttp.Json{
            "message": "Akun Anda tidak aktif",
        })
    }

    // Check password
    if !user.CheckPassword(password) {
        return ctx.Response().Json(401, contractsHttp.Json{
            "message": "Invalid credentials",
        })
    }
    
    // Generate JWT token
    token, err := facades.Auth(ctx).Login(&user)
    if err != nil {
        return ctx.Response().Json(500, contractsHttp.Json{
            "message": "Failed to generate token",
            "error": err.Error(),
        })
    }
    
    return ctx.Response().Json(200, contractsHttp.Json{
        "message": "Login successful",
        "data": contractsHttp.Json{
            "user":  user,
            "token": token,
        },
    })
}

func (r *AuthController) Me(ctx contractsHttp.Context) contractsHttp.Response {
    authHeader := ctx.Request().Header("Authorization")
    token := authHeader
    if len(token) > 7 && token[:7] == "Bearer " {
        token = token[7:]
    }

    payload, err := facades.Auth(ctx).Parse(token)
    if err != nil && !errors.Is(err, auth.ErrorTokenExpired) {
        return ctx.Response().Json(401, contractsHttp.Json{
            "message": "Unauthorized",
        })
    }

    var user models.User
    if err := facades.Auth(ctx).User(&user); err != nil {
        return ctx.Response().Json(401, contractsHttp.Json{
            "message": "User not found",
        })
    }

    return ctx.Response().Json(200, contractsHttp.Json{
        "user": user,
        "payload": payload,
    })
}

func (r *AuthController) Refresh(ctx contractsHttp.Context) contractsHttp.Response {
    authHeader := ctx.Request().Header("Authorization")
    token := authHeader
    if len(token) > 7 && token[:7] == "Bearer " {
        token = token[7:]
    }

    // Parse token dulu
    _, err := facades.Auth(ctx).Parse(token)
    if err != nil {
        return ctx.Response().Json(401, contractsHttp.Json{
            "message": "Token invalid or expired",
        })
    }

    // Refresh token
    newToken, err := facades.Auth(ctx).Refresh()
    if err != nil {
        return ctx.Response().Json(400, contractsHttp.Json{
            "message": "Failed to refresh token",
        })
    }

    return ctx.Response().Json(200, contractsHttp.Json{
        "token": newToken,
    })
}

func (r *AuthController) Logout(ctx contractsHttp.Context) contractsHttp.Response {
    authHeader := ctx.Request().Header("Authorization")
    token := authHeader
    if len(token) > 7 && token[:7] == "Bearer " {
        token = token[7:]
    }

    // Parse token dulu (opsional, untuk validasi)
    _, err := facades.Auth(ctx).Parse(token)
    if err != nil {
        return ctx.Response().Json(http.StatusUnauthorized, contractsHttp.Json{
            "message": "Token invalid or expired",
        })
    }

    if err := facades.Auth(ctx).Logout(); err != nil {
        return ctx.Response().Json(http.StatusInternalServerError, contractsHttp.Json{
            "message": "Failed to logout",
        })
    }

    return ctx.Response().Json(http.StatusOK, contractsHttp.Json{
        "message": "Logout successful",
    })
}

func (r *AuthController) UpdateProfile(ctx contractsHttp.Context) contractsHttp.Response {
    authHeader := ctx.Request().Header("Authorization")
    token := authHeader
    if len(token) > 7 && token[:7] == "Bearer " {
        token = token[7:]
    }

    // Parse token dulu untuk validasi
    _, err := facades.Auth(ctx).Parse(token)
    if err != nil {
        return ctx.Response().Json(401, contractsHttp.Json{"message": "Token invalid or expired"})
    }

    var req struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    if err := ctx.Request().Bind(&req); err != nil {
        return ctx.Response().Json(400, contractsHttp.Json{"message": "Invalid request"})
    }

    var user models.User
    if err := facades.Auth(ctx).User(&user); err != nil {
        return ctx.Response().Json(401, contractsHttp.Json{"message": "Unauthorized"})
    }

    // Update by struct
    updateData := models.User{
        Name:  req.Name,
        Email: req.Email,
    }
    _, err = facades.DB().Table("users").Where("id", user.ID).Update(updateData)
    if err != nil {
        return ctx.Response().Json(500, contractsHttp.Json{"message": "Failed to update profile"})
    }

    // Get updated user
    facades.Orm().Query().Where("id", user.ID).First(&user)

    return ctx.Response().Json(200, contractsHttp.Json{"message": "Profile updated", "user": user})
}

func (r *AuthController) UpdatePassword(ctx contractsHttp.Context) contractsHttp.Response {
    authHeader := ctx.Request().Header("Authorization")
    token := authHeader
    if len(token) > 7 && token[:7] == "Bearer " {
        token = token[7:]
    }

    // Parse token dulu untuk validasi
    _, err := facades.Auth(ctx).Parse(token)
    if err != nil {
        return ctx.Response().Json(401, contractsHttp.Json{"message": "Token invalid or expired"})
    }

    var req struct {
        OldPassword        string `json:"old_password"`
        Password           string `json:"password"`
        ConfirmationPassword string `json:"confirmation_password"`
    }
    if err := ctx.Request().Bind(&req); err != nil {
        return ctx.Response().Json(400, contractsHttp.Json{"message": "Invalid request"})
    }

    // Validation
    errors := make(map[string][]string)
    if req.OldPassword == "" {
        errors["old_password"] = append(errors["old_password"], "Password lama wajib diisi")
    }
    if req.Password == "" {
        errors["password"] = append(errors["password"], "Password baru wajib diisi")
    } else if len(req.Password) < 6 {
        errors["password"] = append(errors["password"], "Password baru minimal 6 karakter")
    }
    if req.ConfirmationPassword == "" {
        errors["confirmation_password"] = append(errors["confirmation_password"], "Konfirmasi password wajib diisi")
    } else if req.Password != req.ConfirmationPassword {
        errors["confirmation_password"] = append(errors["confirmation_password"], "Konfirmasi password tidak sesuai")
    }
    if len(errors) > 0 {
        return ctx.Response().Json(400, contractsHttp.Json{"message": "Validation failed", "errors": errors})
    }

    var user models.User
    if err := facades.Auth(ctx).User(&user); err != nil {
        return ctx.Response().Json(401, contractsHttp.Json{"message": "Unauthorized"})
    }

    // Check old password
    if !user.CheckPassword(req.OldPassword) {
        return ctx.Response().Json(400, contractsHttp.Json{"message": "Password lama salah"})
    }

    if err := user.HashPassword(req.Password); err != nil {
        return ctx.Response().Json(500, contractsHttp.Json{"message": "Failed to hash password"})
    }

    // Update by struct
    updateData := models.User{
        Password: user.Password,
    }
    _, err = facades.DB().Table("users").Where("id", user.ID).Update(updateData)
    if err != nil {
        return ctx.Response().Json(500, contractsHttp.Json{"message": "Failed to update password"})
    }

    return ctx.Response().Json(200, contractsHttp.Json{"message": "Password updated"})
}