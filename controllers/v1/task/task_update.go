package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) UpdateTask(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.UpdateTaskRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.Task.ReadByUUID(ctx, uuid)
	if err != nil {
		log.For(c).Error("[update-task] query task info failed", log.Field("user_id", userID), log.Field("uuid", uuid), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit task of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[update-task] query user info failed", log.Field("user_id", userID), log.Field("uuid", uuid), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != data.OrganizationID {
		log.For(c).Error("[update-task] organization id not match", log.Field("user_id", userID),
			log.Field("user_organizationID", user.OrganizationID), log.Field("data_organizationID", data.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	project, err := h.Project.ReadByUUID(ctx, req.ProjectUUID)
	if err != nil {
		log.For(c).Error("[update-task] query project info failed", log.Field("user_id", userID), log.Field("uuid", uuid), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if project.OrganizationID != user.OrganizationID {
		log.For(c).Error("[update-task] organization id not match", log.Field("user_id", userID),
			log.Field("user_organizationID", user.OrganizationID), log.Field("project_organizationID", project.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	data.UpdatedBy = userID
	data.Title = req.Title
	data.Status = entities.TaskStatus(req.Status)
	data.Priority = entities.TaskPriority(req.Priority)
	data.Label = entities.TaskLabel(req.Label)
	data.Deadline = req.Deadline
	data.EstimatedHours = req.EstimatedHours
	data.Description = req.Description
	data.ProjectID = project.ID

	err = h.Task.Update(ctx, data)
	if err != nil {
		log.For(c).Error("[update-project] update database failed", log.Field("user_id", userID), log.Field("task_id", data.GetID()), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
