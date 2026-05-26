package incident_log

import (
	"net/http"

	"github.com/acidnyi/incident-log-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
)

type incidentLogUpdater = func(
	ctx *gin.Context,
	incidentLog *IncidentLogDefinition,
) (updatedIncidentLog *IncidentLogDefinition, responseContent interface{}, status int)

func updateIncidentLogFunc(ctx *gin.Context, updater incidentLogUpdater) {
	value, exists := ctx.Get("db_service")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service not found",
			"error":   "db_service not found",
		})
		return
	}

	db, ok := value.(db_service.DbService[IncidentLogDefinition])
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service context is not of type db_service.DbService",
			"error":   "cannot cast db_service context to db_service.DbService",
		})
		return
	}

	incidentLogId := ctx.Param("incidentLogId")
	if incidentLogId == "" {
		incidentLogId = "incident-log"
	}

	incidentLog, err := db.FindDocument(ctx.Request.Context(), incidentLogId)

	switch err {
	case nil:
	case db_service.ErrNotFound:
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Incident log not found",
			"error":   err.Error(),
		})
		return
	default:
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "Bad Gateway",
			"message": "Failed to load incident log from database",
			"error":   err.Error(),
		})
		return
	}

	updatedIncidentLog, responseObject, status := updater(ctx, incidentLog)

	if updatedIncidentLog != nil {
		err = db.UpdateDocument(ctx.Request.Context(), incidentLogId, updatedIncidentLog)
	} else {
		err = nil
	}

	switch err {
	case nil:
		if responseObject != nil {
			ctx.JSON(status, responseObject)
		} else {
			ctx.AbortWithStatus(status)
		}
	case db_service.ErrNotFound:
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Incident log was deleted while processing the request",
			"error":   err.Error(),
		})
	default:
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "Bad Gateway",
			"message": "Failed to update incident log in database",
			"error":   err.Error(),
		})
	}
}
