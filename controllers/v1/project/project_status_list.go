package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
)

func (h Handler) ListProjectStatus(c *gin.Context) {
	resp := &apis.ListProjectStatusResponse{}
	for status := entities.Prospect; status <= entities.ProjectReceived; status++ {
		resp.Data = append(resp.Data, &apis.ProjectStatus{
			Value: status.Value(),
			Name: status.Name(),
			Note: status.Note(),
		})
	}

	c.JSON(http.StatusOK, resp)
}
