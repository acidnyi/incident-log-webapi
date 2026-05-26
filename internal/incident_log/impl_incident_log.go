package incident_log

import (
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type implIncidentLogAPI struct {
}

func NewIncidentLogApi() IncidentLogAPI {
	return &implIncidentLogAPI{}
}

func (o implIncidentLogAPI) CreateIncident(c *gin.Context) {
	updateIncidentLogFunc(c, func(
		c *gin.Context,
		incidentLog *IncidentLogDefinition,
	) (*IncidentLogDefinition, interface{}, int) {
		var entry Incident

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		if entry.Id == "" || entry.Id == "@new" {
			entry.Id = uuid.NewString()
		}

		conflictIndex := slices.IndexFunc(incidentLog.Incidents, func(existing Incident) bool {
			return entry.Id == existing.Id
		})

		if conflictIndex >= 0 {
			return nil, gin.H{
				"status":  http.StatusConflict,
				"message": "Incident already exists",
			}, http.StatusConflict
		}

		incidentLog.Incidents = append(incidentLog.Incidents, entry)
		incidentLog.reconcileIncidents()

		entryIndex := slices.IndexFunc(incidentLog.Incidents, func(existing Incident) bool {
			return entry.Id == existing.Id
		})

		if entryIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Failed to save incident",
			}, http.StatusInternalServerError
		}

		return incidentLog, incidentLog.Incidents[entryIndex], http.StatusOK
	})
}

// DeleteIncident - Deletes specific incident
func (o implIncidentLogAPI) DeleteIncident(c *gin.Context) {
	updateIncidentLogFunc(c, func(
		c *gin.Context,
		incidentLog *IncidentLogDefinition,
	) (*IncidentLogDefinition, interface{}, int) {
		incidentId := c.Param("incidentId")

		if incidentId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Incident ID is required",
			}, http.StatusBadRequest
		}

		entryIndex := slices.IndexFunc(incidentLog.Incidents, func(existing Incident) bool {
			return incidentId == existing.Id
		})

		if entryIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Incident not found",
			}, http.StatusNotFound
		}

		incidentLog.Incidents = append(
			incidentLog.Incidents[:entryIndex],
			incidentLog.Incidents[entryIndex+1:]...,
		)

		incidentLog.reconcileIncidents()

		return incidentLog, nil, http.StatusNoContent
	})
}

func (o implIncidentLogAPI) GetIncidents(c *gin.Context) {
	updateIncidentLogFunc(c, func(
		c *gin.Context,
		incidentLog *IncidentLogDefinition,
	) (*IncidentLogDefinition, interface{}, int) {
		result := incidentLog.Incidents
		if result == nil {
			result = []Incident{}
		}

		return nil, result, http.StatusOK
	})
}

// GetIncident - Provides details about incident
func (o implIncidentLogAPI) GetIncident(c *gin.Context) {
	updateIncidentLogFunc(c, func(
		c *gin.Context,
		incidentLog *IncidentLogDefinition,
	) (*IncidentLogDefinition, interface{}, int) {
		incidentId := c.Param("incidentId")

		if incidentId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Incident ID is required",
			}, http.StatusBadRequest
		}

		entryIndex := slices.IndexFunc(incidentLog.Incidents, func(existing Incident) bool {
			return incidentId == existing.Id
		})

		if entryIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Incident not found",
			}, http.StatusNotFound
		}

		return nil, incidentLog.Incidents[entryIndex], http.StatusOK
	})
}

// UpdateIncident - Updates specific incident
func (o implIncidentLogAPI) UpdateIncident(c *gin.Context) {
	updateIncidentLogFunc(c, func(
		c *gin.Context,
		incidentLog *IncidentLogDefinition,
	) (*IncidentLogDefinition, interface{}, int) {
		var entry Incident

		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}

		incidentId := c.Param("incidentId")

		if incidentId == "" {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Incident ID is required",
			}, http.StatusBadRequest
		}

		entryIndex := slices.IndexFunc(incidentLog.Incidents, func(existing Incident) bool {
			return incidentId == existing.Id
		})

		if entryIndex < 0 {
			return nil, gin.H{
				"status":  http.StatusNotFound,
				"message": "Incident not found",
			}, http.StatusNotFound
		}

		if entry.Id != "" {
			incidentLog.Incidents[entryIndex].Id = entry.Id
		}

		if entry.IncidentType != "" {
			incidentLog.Incidents[entryIndex].IncidentType = entry.IncidentType
		}

		if entry.Location != "" {
			incidentLog.Incidents[entryIndex].Location = entry.Location
		}

		if entry.OccurredAt.After(time.Time{}) {
			incidentLog.Incidents[entryIndex].OccurredAt = entry.OccurredAt
		}

		if entry.Description != "" {
			incidentLog.Incidents[entryIndex].Description = entry.Description
		}

		if entry.Severity != "" {
			incidentLog.Incidents[entryIndex].Severity = entry.Severity
		}

		if entry.Status != "" {
			incidentLog.Incidents[entryIndex].Status = entry.Status
		}

		if entry.Attachments != nil {
			incidentLog.Incidents[entryIndex].Attachments = entry.Attachments
		}

		if entry.InvestigationReport != "" {
			incidentLog.Incidents[entryIndex].InvestigationReport = entry.InvestigationReport
		}

		if entry.Notes != "" {
			incidentLog.Incidents[entryIndex].Notes = entry.Notes
		}

		incidentLog.reconcileIncidents()

		updatedIndex := slices.IndexFunc(incidentLog.Incidents, func(existing Incident) bool {
			return incidentLog.Incidents[entryIndex].Id == existing.Id
		})

		if updatedIndex < 0 {
			updatedIndex = entryIndex
		}

		return incidentLog, incidentLog.Incidents[updatedIndex], http.StatusOK
	})
}
