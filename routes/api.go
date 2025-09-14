package routes

import (
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"goravel/app/http/controllers"
    "goravel/app/http/middleware"
)

func Api() {
	authController := controllers.NewAuthController()
    jwtMiddleware := middleware.NewJWTMiddleware()
	userController := controllers.NewUserController()
	postController := controllers.NewPostController()

	facades.Route().Group(func(router route.Router) {
        // Public routes
        router.Post("/register", authController.Register)
        router.Post("/login", authController.Login)

		router.Get("/test", authController.Test)
        
        // Protected routes
        router.Group(func(router route.Router) {
            router.Middleware(jwtMiddleware.Handle)
            router.Get("/me", authController.Me)
            router.Post("/refresh", authController.Refresh)
            router.Post("/logout", authController.Logout)
			router.Put("/profile", authController.UpdateProfile)
    		router.Put("/password", authController.UpdatePassword)
			router.Put("/users/{id}/active", userController.SetActive)
        })
    })
	
	facades.Route().Get("/users/{id}", userController.Show)

	facades.Route().Get("/posts", postController.Index)
	facades.Route().Get("/posts/{id}", postController.Show)
	facades.Route().Post("/posts", postController.Store)
	facades.Route().Put("/posts/{id}", postController.Update)
	facades.Route().Delete("/posts/{id}", postController.Delete)
}
