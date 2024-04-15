package task

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/normalize"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) ReadTask(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.Task.ReadByUUID(ctx, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can read task of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != data.OrganizationID {
		c.JSON(http.StatusForbidden, err)
		return
	}

	resp := &apis.Task{
		UUID:           data.UUID,
		Title:          data.Title,
		Status:         apis.TaskStatus(data.Status),
		Priority:       apis.TaskPriority(data.Priority),
		Label:          apis.TaskLabel(data.Label),
		Deadline:       data.Deadline,
		EstimatedHours: data.EstimatedHours,
		Description:    data.Description,
		CreatedTime:    data.GetCreatedAt().UnixMilli(),
	}

	if data.GetCreator() != nil {
		resp.Creator = &apis.User{
			UUID:        data.GetCreator().GetUUID(),
			DisplayName: data.GetCreator().GetDisplayName(),
		}
	}

	if data.GetProject() != nil {
		resp.Project = &apis.Project{
			UUID:      data.GetProject().GetUUID(),
			FullName:  data.GetProject().GetFullName(),
			ShortName: data.GetProject().GetShortName(),
		}
	}

	for _, assignee := range data.GetAssignees() {
		resp.Assignees = append(resp.Assignees, &apis.User{
			UUID:        assignee.GetAssignee().GetUUID(),
			DisplayName: assignee.GetAssignee().GetDisplayName(),
		})
	}

	for _, val := range data.GetAttachFiles() {
		fileName := normalize.URLEncode(val.GetOriginalName()) + val.Ext
		resp.AttachFiles = append(resp.AttachFiles, &apis.File{
			UUID: val.GetUUID(),
			Name: val.GetOriginalName() + val.Ext,
			URL:  fmt.Sprintf("%sapi/v1/tasks/%s/downloaded-files/%s/%s", services.Config.Domain.Domain, resp.GetUUID(), val.GetUUID(), fileName),
		})
	}

	c.JSON(http.StatusOK, resp)
}
