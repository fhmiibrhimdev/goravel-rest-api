package config

import (
    "github.com/goravel/framework/facades"
)

func init() {
    config := facades.Config()
    config.Add("auth", map[string]interface{}{
        "defaults": map[string]interface{}{
            "guard": "api",
        },
        "guards": map[string]interface{}{
            "api": map[string]interface{}{
                "driver": "jwt",
                "provider": "users",
                "ttl": 60, // 60 minutes untuk guard ini
            },
        },
        "providers": map[string]interface{}{
            "users": map[string]interface{}{
                "driver": "orm",
                "model": "models.User",
            },
        },
    })
}