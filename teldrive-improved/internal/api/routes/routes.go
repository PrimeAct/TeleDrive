package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/tgdrive/teldrive/internal/api/handlers"
	"github.com/tgdrive/teldrive/internal/api/middleware"
	"github.com/tgdrive/teldrive/internal/config"
	"github.com/tgdrive/teldrive/pkg/services"
)

func Register(r chi.Router, svc *services.Services, cfg *config.ServerConfig) {
	r.Group(func(r chi.Router) {
		r.Get("/health", handlers.HealthCheck)
		r.Post("/auth/login", handlers.Login(svc.Auth))
		r.Post("/auth/refresh", handlers.RefreshToken(svc.Auth))
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuth(&cfg.JWT))

		r.Route("/files", func(r chi.Router) {
			r.Get("/", handlers.ListFiles(svc.File))
			r.Post("/", handlers.CreateFile(svc.File))
			r.Get("/{id}", handlers.GetFile(svc.File))
			r.Put("/{id}", handlers.UpdateFile(svc.File))
			r.Delete("/{id}", handlers.DeleteFile(svc.File))
			r.Post("/{id}/share", handlers.ShareFile(svc.Share))
		})

		r.Route("/uploads", func(r chi.Router) {
			r.Post("/", handlers.StartUpload(svc.Upload))
			r.Post("/{id}/chunk", handlers.UploadChunk(svc.Upload))
			r.Post("/{id}/complete", handlers.CompleteUpload(svc.Upload))
		})

		r.Get("/user/me", handlers.GetCurrentUser(svc.User))
		r.Get("/user/stats", handlers.GetUserStats(svc.User))
	})
}
