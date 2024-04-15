package task_attach_file

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) CreateTaskAttachFile(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-task-attach-file] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	taskUUID := c.Param("uuid")
	task, err := h.Task.ReadByUUID(ctx, taskUUID)
	if err != nil {
		log.For(c).Error("[create-task-attach-file] query task info", log.Field("user_id", userID),
			log.Field("task_id", task.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != task.OrganizationID {
		c.JSON(http.StatusForbidden, err)
		return
	}

	file, err := services.UploadFile(c, "file")
	if err != nil {
		log.For(c).Error("[create-task-attach-file] write file failed", log.Field("user_id", userID),
			log.Field("task_id", task.ID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
	}

	data := &entities.TaskAttachFile{
		OriginalName:   file.OriginalName,
		RelativePath:   file.RelativePathFile,
		Ext:            file.FileExt,
		OrganizationID: user.OrganizationID,
		TaskID:         task.ID,
		CreatedBy:      userID,
	}
	_, err = h.TaskAttachFile.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-task-attach-file] update database failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: data.UUID,
	})
}
