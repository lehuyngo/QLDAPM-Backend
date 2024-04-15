package task_attach_file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) DeleteTaskAttachFile(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit task of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[delete-task-attach-file] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	taskUUID := c.Param("uuid")
	task, err := h.Task.ReadByUUID(ctx, taskUUID)
	if err != nil {
		log.For(c).Error("[delete-task-attach-file] query task info failed", log.Field("user_id", userID), log.Field("task_uuid", taskUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	fileUUID := c.Param("file_uuid")
	file, err := h.TaskAttachFile.Read(ctx, fileUUID)
	if err != nil {
		log.For(c).Error("[delete-task-attach-file] query task info failed", log.Field("user_id", userID), log.Field("ctask_uuid", taskUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != file.OrganizationID {
		log.For(c).Error("[delete-task-attach-file] query task id is not match", log.Field("user_id", userID), log.Field("task_uuid", taskUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("task_organization_id", task.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	err = h.TaskAttachFile.Delete(ctx, task.ID, file.ID)
	if err != nil {
		log.For(c).Error("[delete-task-task] udpate database failed", log.Field("user_id", userID),
			log.Field("task_id", task.ID), log.Field("file_uuid", file.UUID),
			log.Field("file_id", file.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[update-task-task] process success", log.Field("user_id", userID),
		log.Field("task_id", task.ID), log.Field("file_uuid", file.UUID),
		log.Field("file_id", file.ID))

	c.JSON(http.StatusOK, nil)
}
