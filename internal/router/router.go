package router

import (
	"go-crud-oapi/internal/controller"
	"go-crud-oapi/internal/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(userController *controller.UserController, authCtrl *controller.AuthController) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Post("/login", authCtrl.Login)

	// User routes
	r.Route("/users", func(r chi.Router) {
		r.Get("/", userController.ListUsers)
		r.Get("/{id}", userController.GetUser)

		r.With(middleware.JWTAuthMiddleware).Post("/", userController.CreateUser)
		r.With(middleware.JWTAuthMiddleware).Put("/{id}", userController.UpdateUser)
		r.With(middleware.JWTAuthMiddleware).Delete("/{id}", userController.DeleteUser)
	})

	return r
}
