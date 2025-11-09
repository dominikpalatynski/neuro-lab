package server

import (
	"config/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

type Server struct {
	db            *gorm.DB
	router        *chi.Mux
	deviceHandler *handlers.DeviceHandler
}

func NewServer(db *gorm.DB, r *chi.Mux) *Server {
	deviceHandler := handlers.NewDeviceHandler(db)
	return &Server{db: db, router: r, deviceHandler: deviceHandler}
}

func (s *Server) Start() {
	s.router.Use(middleware.Logger)

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Route("/device", func(r chi.Router) {
			r.Post("/", s.deviceHandler.CreateDevice)
			r.Put("/{id}", s.deviceHandler.UpdateDevice)
			r.Delete("/{id}", s.deviceHandler.DeleteDevice)
			r.Get("/{id}", s.deviceHandler.GetDevice)
		})
	})
}
