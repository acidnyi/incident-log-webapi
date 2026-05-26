package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/acidnyi/incident-log-webapi/api"
	"github.com/acidnyi/incident-log-webapi/internal/db_service"
	"github.com/acidnyi/incident-log-webapi/internal/incident_log"
	"github.com/gin-contrib/cors"
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

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	dbService := db_service.NewMongoService[incident_log.IncidentLogDefinition](
		db_service.MongoServiceConfig{},
	)
	defer dbService.Disconnect(context.Background())

	engine.Use(func(ctx *gin.Context) {
		ctx.Set("db_service", dbService)
		ctx.Next()
	})

	handleFunctions := &incident_log.ApiHandleFunctions{
		IncidentLogAPI:   incident_log.NewIncidentLogApi(),
		IncidentTypesAPI: incident_log.NewIncidentTypesApi(),
		IncidentLogsAPI:  incident_log.NewIncidentLogsApi(),
	}

	incident_log.NewRouterWithGinEngine(engine, *handleFunctions)

	engine.GET("/openapi", api.HandleOpenApi)

	engine.Run(":" + port)
}
