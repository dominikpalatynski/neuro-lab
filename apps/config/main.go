package main

import (
	"net/http"

	"config/server"
	"database"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	db := database.Connect()
	db.AutoMigrate(&database.Device{}, &database.TestSession{}, &database.Scenario{}, &database.ScenarioCondition{}, &database.ConditionValue{})

	srv := server.NewServer(db, r)
	srv.Start()

	http.ListenAndServe(":3002", r)
}
