package main

import (
	"log"
	"os"
	"strings"

	"github.com/acidnyi/incident-log-webapi/api"
	"github.com/acidnyi/incident-log-webapi/internal/incident_log"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Printf("IncidentLog API server started")

	port := os.Getenv("INCIDENT_LOG_API_PORT")
	if port == "" {
		port = "8080"
	}

	environment := os.Getenv("INCIDENT_LOG_API_ENVIRONMENT")

	if !strings.EqualFold(environment, "production") {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	// request routings
	handleFunctions := &incident_log.ApiHandleFunctions{
		IncidentLogAPI:   incident_log.NewIncidentLogApi(),
		IncidentTypesAPI: incident_log.NewIncidentTypesApi(),
	}

	incident_log.NewRouterWithGinEngine(engine, *handleFunctions)

	engine.GET("/openapi", api.HandleOpenApi)

	engine.Run(":" + port)
}
