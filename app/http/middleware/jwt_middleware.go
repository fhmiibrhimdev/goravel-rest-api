package middleware

import (
    "errors"
    "net/http"
    "strings"
    
    contractsHttp "github.com/goravel/framework/contracts/http"
    "github.com/goravel/framework/facades"
    "github.com/goravel/framework/auth"
)

type JWTMiddleware struct{}

func NewJWTMiddleware() *JWTMiddleware {
    return &JWTMiddleware{}
}

func (m *JWTMiddleware) Handle(ctx contractsHttp.Context) {
    token := m.getTokenFromRequest(ctx)
    if token == "" {
        ctx.Response().Json(http.StatusUnauthorized, contractsHttp.Json{
            "message": "Token is required",
        })
        return
    }

    // Parse token dan set user context
    _, err := facades.Auth(ctx).Parse(token)
    if err != nil {
        if errors.Is(err, auth.ErrorTokenExpired) {
            ctx.Response().Json(http.StatusUnauthorized, contractsHttp.Json{
                "message": "Token expired",
            })
        } else {
            ctx.Response().Json(http.StatusUnauthorized, contractsHttp.Json{
                "message": "Invalid token",
            })
        }
        return
    }

    // Lanjutkan ke next handler
    ctx.Request().Next()
}

func (m *JWTMiddleware) getTokenFromRequest(ctx contractsHttp.Context) string {
    // Get token from Authorization header
    authHeader := ctx.Request().Header("Authorization")
    if authHeader != "" {
        if strings.HasPrefix(authHeader, "Bearer ") {
            return strings.TrimPrefix(authHeader, "Bearer ")
        }
        return authHeader
    }

    // Fallback ke query parameter
    return ctx.Request().Query("token", "")
}