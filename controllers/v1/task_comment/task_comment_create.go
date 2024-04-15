package task_comment

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) CreateTaskComment(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateTaskCommentRequest{}

	err := http_parser.BindJSONAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-task-comment] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	taskUUID := c.Param("uuid")
	task, err := h.Task.ReadByUUID(ctx, taskUUID)
	if err != nil {
		log.For(c).Error("[create-task-comment] query task info failed", log.Field("user_id", userID), log.Field("task_uuid", taskUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != task.OrganizationID {
		log.For(c).Error("[create-task-comment] query task id isnot match", log.Field("user_id", userID), log.Field("task_uuid", taskUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("task_organization_id", task.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	data := &entities.TaskComment{
		Content:   req.Content,
		TaskID:    task.ID,
		CreatedBy: user.ID,
	}
	_, err = h.TaskComment.Create(ctx, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, &apis.CreateResponse{UUID: data.UUID})
}
