package task_attach_file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) DownloadTaskAttachFile(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	taskUUID := c.Param("uuid")
	task, err := h.Task.ReadByUUID(ctx, taskUUID)
	if err != nil {
		log.For(c).Error("[task-attach-file-download] query task info", log.Field("task_uuid", taskUUID),
			log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit task of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[task-attach-file-download] query user info", log.Field("task_uuid", taskUUID), log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != task.OrganizationID {
		log.For(c).Error("[task-attach-file-download] organization id not match",
			log.Field("user_id", userID), log.Field("task_uuid", taskUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("task_organization_id", task.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	fileUUID := c.Param("file_uuid")
	file, err := h.TaskAttachFile.Read(ctx, fileUUID)
	if err != nil {
		log.For(c).Error("[task-attach-file-download] query file info failed", log.Field("task_uuid", taskUUID),
			log.Field("user_id", userID), log.Field("file_uuid", fileUUID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}
	if file.TaskID != task.ID {
		log.For(c).Error("[task-attach-file-download] task id not match",
			log.Field("user_id", userID), log.Field("task_id", task.ID), log.Field("file_id", file.ID),
			log.Field("task_id", file.TaskID), log.Field("file_task_id", file.TaskID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	filePath := services.Config.FileStorage.Folder + file.GetRelativePath()
	fileName := file.GetOriginalName() + file.GetExt()
	log.For(c).Info("[task-attach-file-download] download",
		log.Field("user_id", userID), log.Field("task_uuid", taskUUID), log.Field("file_uuid", file.UUID),
		log.Field("file_path", filePath), log.Field("file_name", fileName))

	c.FileAttachment(filePath, fileName)
}
