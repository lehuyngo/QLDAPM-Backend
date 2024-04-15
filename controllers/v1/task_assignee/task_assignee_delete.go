package task_assignee

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) DeleteTaskAssignee(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[delete-task-assignee] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// get task
	taskUUID := c.Param("uuid")
	task, err := h.Task.ReadByUUID(ctx, taskUUID)
	if err != nil {
		log.For(c).Error("[delete-task-assignee] query task info failed", log.Field("user_id", userID), log.Field("task_uuid", taskUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// organization check
	if user.OrganizationID != task.OrganizationID {
		log.For(c).Error("[delete-task-assignee] query organization id is not match", log.Field("user_id", userID), log.Field("task_uuid", taskUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("task_organization_id", task.OrganizationID))

		c.JSON(http.StatusForbidden, nil)
		return
	}

	// get assignee
	assigneeUUID := c.Param("assignee_uuid")
	assignee, err := h.User.ReadByUUID(ctx, assigneeUUID)
	if err != nil {
		log.For(c).Error("[delete-task-assignee] query assignee info failed", log.Field("user_id", userID), log.Field("assignee_uuid", assigneeUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// organization check
	if user.OrganizationID != assignee.OrganizationID {
		log.For(c).Error("[delete-task-assignee] query organization id is not match", log.Field("user_id", userID), log.Field("task_uuid", taskUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("assignee_organization_id", assignee.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	// delete task assignee
	err = h.TaskAssignee.Delete(ctx, task.ID, assignee.ID)
	if err != nil {
		log.For(c).Error("[delete-task-assignee] udpate database failed", log.Field("user_id", userID),
			log.Field("task_uuid", task.UUID), log.Field("assignee_uuid", assignee.UUID),
			log.Field("task_id", task.ID), log.Field("assignee_id", assignee.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[delete-task-assignee] process success", log.Field("user_id", userID),
		log.Field("task_uuid", task.UUID), log.Field("assignee_uuid", assignee.UUID))

	c.JSON(http.StatusOK, nil)
}
