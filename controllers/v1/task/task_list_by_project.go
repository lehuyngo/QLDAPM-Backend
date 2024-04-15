package task

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/normalize"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) ListTaskByProject(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	projectUUID := c.Param("uuid")
	project, err := h.Project.ReadByUUID(ctx, projectUUID)
	if err != nil {
		log.For(c).Error("[list-task-by-project] query meeting info failed", log.Field("user_id", userID), log.Field("project_uuid", projectUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != project.OrganizationID {
		log.For(c).Error("[list-task-by-project] organization id is not match", log.Field("user_id", userID), log.Field("project_uuid", projectUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("project_organization_id", project.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	data, err := h.Task.ListByProjectID(ctx, project.ID)
	if err != nil {
		log.For(c).Error("[list-task-by-project] query list task failed", log.Field("user_id", userID), log.Field("project_uuid", projectUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("project_organization_id", project.OrganizationID))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListTaskResponse{}
	for _, val := range data {
		if val.GetCreator() == nil || val.GetProject() == nil {
			continue
		}
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
			Project: &apis.Project{
				UUID:      val.GetProject().GetUUID(),
				FullName:  val.GetProject().GetFullName(),
				ShortName: val.GetProject().GetShortName(),
			},
			Creator: &apis.User{
				UUID:        val.GetCreator().GetUUID(),
				DisplayName: val.GetCreator().GetDisplayName(),
			},
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

	c.JSON(http.StatusOK, resp)
}
