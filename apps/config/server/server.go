package server

import (
	"config/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

type Server struct {
	db                         *gorm.DB
	router                     *chi.Mux
	deviceHandler              *handlers.DeviceHandler
	testSessionHandler         *handlers.TestSessionHandler
	conditionHandler           *handlers.ConditionHandler
	conditionValueHandler      *handlers.ConditionValueHandler
	scenarioHandler            *handlers.ScenarioHandler
	scenarioConditionHandler   *handlers.ScenarioConditionHandler
	scenarioValidationHandler  *handlers.ScenarioValidationHandler
	discoveryHandler           *handlers.DiscoveryHandler
}

func NewServer(db *gorm.DB, r *chi.Mux) *Server {
	deviceHandler := handlers.NewDeviceHandler(db)
	testSessionHandler := handlers.NewTestSessionHandler(db)
	conditionHandler := handlers.NewConditionHandler(db)
	conditionValueHandler := handlers.NewConditionValueHandler(db)
	scenarioHandler := handlers.NewScenarioHandler(db)
	scenarioConditionHandler := handlers.NewScenarioConditionHandler(db)
	scenarioValidationHandler := handlers.NewScenarioValidationHandler(db)
	discoveryHandler := handlers.NewDiscoveryHandler(db)

	return &Server{
		db:                        db,
		router:                    r,
		deviceHandler:             deviceHandler,
		testSessionHandler:        testSessionHandler,
		conditionHandler:          conditionHandler,
		conditionValueHandler:     conditionValueHandler,
		scenarioHandler:           scenarioHandler,
		scenarioConditionHandler:  scenarioConditionHandler,
		scenarioValidationHandler: scenarioValidationHandler,
		discoveryHandler:          discoveryHandler,
	}
}

func (s *Server) Start() {
	// Add error handling and request tracking middleware
	s.router.Use(middleware.Logger)

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Get("/", s.discoveryHandler.GetAPIResources)
		})

		r.Post("/scenario-validation", s.scenarioValidationHandler.ValidateScenario)

		r.Route("/device", func(r chi.Router) {
			r.Post("/", s.deviceHandler.CreateDevice)
			r.Put("/{id}", s.deviceHandler.UpdateDevice)
			r.Delete("/{id}", s.deviceHandler.DeleteDevice)
			r.Get("/{id}", s.deviceHandler.GetDevice)
			r.Get("/", s.deviceHandler.GetDevices)
		})

		r.Route("/test-session", func(r chi.Router) {
			r.Post("/", s.testSessionHandler.CreateTestSession)
			r.Put("/{id}", s.testSessionHandler.UpdateTestSession)
			r.Delete("/{id}", s.testSessionHandler.DeleteTestSession)
			r.Get("/{id}", s.testSessionHandler.GetTestSession)
			r.Get("/list/{deviceID}", s.testSessionHandler.GetTestSessionsByDevice)
		})

		r.Route("/condition", func(r chi.Router) {
			r.Post("/", s.conditionHandler.CreateCondition)
			r.Put("/{id}", s.conditionHandler.UpdateCondition)
			r.Delete("/{id}", s.conditionHandler.DeleteCondition)
			r.Get("/{id}", s.conditionHandler.GetCondition)
			r.Get("/", s.conditionHandler.GetConditions)
		})

		r.Route("/condition-value", func(r chi.Router) {
			r.Post("/", s.conditionValueHandler.CreateConditionValue)
			r.Put("/{id}", s.conditionValueHandler.UpdateConditionValue)
			r.Delete("/{id}", s.conditionValueHandler.DeleteConditionValue)
			r.Get("/{id}", s.conditionValueHandler.GetConditionValue)
			r.Get("/", s.conditionValueHandler.GetConditionValues)
			r.Get("/list/{conditionID}", s.conditionValueHandler.GetConditionValuesByCondition)
		})

		r.Route("/scenario", func(r chi.Router) {
			r.Post("/", s.scenarioHandler.CreateScenario)
			r.Post("/with-condition-values", s.scenarioHandler.CreateScenarioWithConditionValues)
			r.Put("/{id}", s.scenarioHandler.UpdateScenario)
			r.Delete("/{id}", s.scenarioHandler.DeleteScenario)
			r.Get("/{id}", s.scenarioHandler.GetScenario)
			r.Get("/list/{testSessionID}", s.scenarioHandler.GetScenariosByTestSession)
		})

		r.Route("/scenario-condition", func(r chi.Router) {
			r.Post("/", s.scenarioConditionHandler.CreateScenarioCondition)
			r.Put("/{id}", s.scenarioConditionHandler.UpdateScenarioCondition)
			r.Delete("/{id}", s.scenarioConditionHandler.DeleteScenarioCondition)
			r.Get("/{id}", s.scenarioConditionHandler.GetScenarioCondition)
		})
	})
}
