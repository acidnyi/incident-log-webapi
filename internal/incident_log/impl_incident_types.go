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

func (o implIncidentTypesAPI) GetIncidentTypes(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
