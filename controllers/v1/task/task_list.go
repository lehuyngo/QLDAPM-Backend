package task

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/normalize"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) ListTask(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	data, err := h.Task.ListByOrgID(ctx, user.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListTaskResponse{}
	for _, val := range data {
		task := &apis.Task{
			UUID:           val.UUID,
			Title:          val.Title,
			Status:         apis.TaskStatus(val.Status),
			Priority:       apis.TaskPriority(val.Priority),
			Label:          apis.TaskLabel(val.Label),
			Deadline:       val.Deadline,
			EstimatedHours: val.EstimatedHours,
			Description:    val.Description,
			CreatedTime:    val.GetCreatedAt().UnixMilli(),
		}

		if val.GetCreator() != nil {
			task.Creator = &apis.User{
				UUID:        val.GetCreator().GetUUID(),
				DisplayName: val.GetCreator().GetDisplayName(),
			}
		}

		if val.GetProject() != nil {
			task.Project = &apis.Project{
				UUID:      val.GetProject().GetUUID(),
				FullName:  val.GetProject().GetFullName(),
				ShortName: val.GetProject().GetShortName(),
			}
		}

		for _, assignee := range val.GetAssignees() {
			task.Assignees = append(task.Assignees, &apis.User{
				UUID:        assignee.GetAssignee().GetUUID(),
				DisplayName: assignee.GetAssignee().GetDisplayName(),
			})
		}

		for _, attachFile := range val.GetAttachFiles() {
			fileName := normalize.URLEncode(attachFile.GetOriginalName()) + attachFile.Ext
			task.AttachFiles = append(task.AttachFiles, &apis.File{
				UUID: attachFile.GetUUID(),
				Name: attachFile.GetOriginalName() + attachFile.Ext,
				URL:  fmt.Sprintf("%sapi/v1/tasks/%s/downloaded-files/%s/%s", services.Config.Domain.Domain, task.GetUUID(), attachFile.GetUUID(), fileName),
			})
		}

		resp.Data = append(resp.Data, task)
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return strings.EqualFold(resp.Data[i].Title, resp.Data[j].Title)
	})

	c.JSON(http.StatusOK, resp)
}
