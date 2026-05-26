package incident_log

import (
	"net/http"

	"github.com/acidnyi/incident-log-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type implIncidentLogsAPI struct {
}

func NewIncidentLogsApi() IncidentLogsAPI {
	return &implIncidentLogsAPI{}
}

// CreateIncidentLog - Saves new incident log definition
func (o *implIncidentLogsAPI) CreateIncidentLog(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db not found",
			"error":   "db not found",
		})
		return
	}

	db, ok := value.(db_service.DbService[IncidentLogDefinition])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db context is not of required type",
			"error":   "cannot cast db context to db_service.DbService",
		})
		return
	}

	incidentLog := IncidentLogDefinition{}
	err := c.BindJSON(&incidentLog)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad Request",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if incidentLog.Id == "" {
		incidentLog.Id = uuid.New().String()
	}

	err = db.CreateDocument(c.Request.Context(), incidentLog.Id, &incidentLog)

	switch err {
	case nil:
		c.JSON(http.StatusCreated, incidentLog)
	case db_service.ErrConflict:
		c.JSON(http.StatusConflict, gin.H{
			"status":  "Conflict",
			"message": "Incident log already exists",
			"error":   err.Error(),
		})
	default:
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  "Bad Gateway",
			"message": "Failed to create incident log in database",
			"error":   err.Error(),
		})
	}
}

func (o *implIncidentLogsAPI) DeleteIncidentLog(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service not found",
			"error":   "db_service not found",
		})
		return
	}

	db, ok := value.(db_service.DbService[IncidentLogDefinition])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Server Error",
			"message": "db_service context is not of type db_service.DbService",
			"error":   "cannot cast db_service context to db_service.DbService",
		})
		return
	}

	incidentLogId := c.Param("incidentLogId")
	err := db.DeleteDocument(c.Request.Context(), incidentLogId)

	switch err {
	case nil:
		c.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "Not Found",
			"message": "Incident log not found",
			"error":   err.Error(),
		})
	default:
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  "Bad Gateway",
			"message": "Failed to delete incident log from database",
			"error":   err.Error(),
		})
	}
}
