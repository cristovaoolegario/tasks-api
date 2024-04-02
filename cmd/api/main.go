package main

import (
	"fmt"
	"github.com/cristovaoolegario/tasks-api/internal/infra"
	"github.com/heptiolabs/healthcheck"
	"net/http"
	"time"
)

func main() {

	fmt.Println("Starting app")
	dbConnect := infra.InitDB()
	db, err := dbConnect.DB()
	if err != nil {
		panic(err)
	}
	// Create a new health check handler
	health := healthcheck.NewHandler()

	// Register health checks for any dependencies
	// Serve http://0.0.0.0:3000/live and http://0.0.0.0:3000/ready endpoints.
	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(db, 1*time.Second))

	fmt.Println("Ready...")

	http.ListenAndServe(":3000", health)
}
