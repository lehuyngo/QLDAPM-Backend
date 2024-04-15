package task_assignee

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) CreateTaskAssignee(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateTaskAssigneeRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-task-assignee] query user failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// get task
	taskUUID := c.Param("uuid")
	if taskUUID == "" {
		log.For(c).Error("[create-task-assignee] task uuid is empty", log.Field("user_id", userID))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	task, err := h.Task.ReadByUUID(ctx, taskUUID)
	if err != nil {
		log.For(c).Error("[create-task-assignee] query task by uuid failed", log.Field("user_id", userID), log.Field("task_uuid", taskUUID))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if user.OrganizationID != task.OrganizationID {
		log.For(c).Error("[create-task-assignee] user and task organization is not match", log.Field("user_id", userID), log.Field("task_uuid", taskUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("task_organization_id", task.OrganizationID))
		c.JSON(http.StatusForbidden, nil)
		return
	}

	// get assignee
	assigneeUUID := req.AssigneeUUID
	if assigneeUUID == "" {
		log.For(c).Error("[create-task-assignee] assignee uuid is empty", log.Field("user_id", userID))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	assignee, err := h.User.ReadByUUID(ctx, assigneeUUID)
	if err != nil {
		log.For(c).Error("[create-task-assignee] query assignee by uuid failed", log.Field("user_id", userID), log.Field("assignee_uuid", assigneeUUID))
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if user.OrganizationID != assignee.OrganizationID {
		log.For(c).Error("[create-task-assignee] user and assignee organization is not match", log.Field("user_id", userID), log.Field("assignee_uuid", assigneeUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("assignee_organization_id", assignee.OrganizationID))
		c.JSON(http.StatusForbidden, nil)
		return
	}

	// create task assignee
	err = h.TaskAssignee.Create(ctx, &entities.TaskAssignee{
		TaskID:     task.GetID(),
		AssigneeID: assignee.GetID(),
		CreatedBy:  user.ID,
	})
	if err != nil {
		log.For(c).Error("[create-task-assignee] create task assignee failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}
