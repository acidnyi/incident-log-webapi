package incident_log

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implIncidentTypesAPI struct {
}

func NewIncidentTypesApi() IncidentTypesAPI {
	return &implIncidentTypesAPI{}
}

func (o *implIncidentTypesAPI) GetIncidentTypes(c *gin.Context) {
	updateIncidentLogFunc(c, func(
		c *gin.Context,
		incidentLog *IncidentLogDefinition,
	) (updatedIncidentLog *IncidentLogDefinition, responseContent interface{}, status int) {
		result := incidentLog.PredefinedIncidentTypes
		if result == nil {
			result = []IncidentType{}
		}

		return nil, result, http.StatusOK
	})
}
