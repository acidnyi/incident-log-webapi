package incident_log

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implIncidentLogAPI struct {
}

func NewIncidentLogApi() IncidentLogAPI {
	return &implIncidentLogAPI{}
}

func (o implIncidentLogAPI) CreateIncident(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implIncidentLogAPI) DeleteIncident(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implIncidentLogAPI) GetIncident(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implIncidentLogAPI) GetIncidents(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implIncidentLogAPI) UpdateIncident(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
