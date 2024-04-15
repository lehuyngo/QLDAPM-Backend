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

func (h Handler) UpdateTaskStatus(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.UpdateTaskStatusRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.Task.ReadByUUID(ctx, uuid)
	if err != nil {
		log.For(c).Error("[update-task-status] query task info failed", log.Field("user_id", userID), log.Field("uuid", uuid), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit task of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[update-task-status] query user info failed", log.Field("user_id", userID), log.Field("uuid", uuid), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != data.OrganizationID {
		log.For(c).Error("[update-task-status] organization id not match", log.Field("user_id", userID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("task_organization_id", data.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	data.UpdatedBy = userID
	data.Status = entities.TaskStatus(req.Status)

	err = h.Task.Update(ctx, data)
	if err != nil {
		log.For(c).Error("[update-task-status] update task status failed", log.Field("user_id", userID), log.Field("task_id", data.GetID()), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
