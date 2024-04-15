package task

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) CreateTask(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.CreateTaskRequest{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[create-task] query user failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	data := &entities.Task{
		Base: entities.Base{
			CreatedBy: userID,
			UpdatedBy: userID,
		},
		Title:          req.Title,
		Status:         entities.TaskStatus(req.Status),
		Priority:       entities.TaskPriority(req.Priority),
		Label:          entities.TaskLabel(req.Label),
		Deadline:       req.Deadline,
		EstimatedHours: req.EstimatedHours,
		Description:    req.Description,
		OrganizationID: user.OrganizationID,
	}

	if req.ProjectUUID == "" {
		log.For(c).Error("[create-project] project uuid is empty", log.Field("user_id", userID))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	project, err := h.Project.ReadByUUID(ctx, req.ProjectUUID)
	if err != nil {
		log.For(c).Error("[create-project] query project by uuid failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusBadRequest, err)
		return
	}
	data.ProjectID = project.GetID()

	userList, err := h.User.ListByUUIDs(ctx, strings.Split(req.AssigneeUUIDs, ","))
	if err != nil {
		log.For(c).Error("[create-project] query user by username failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// list assignees
	if len(userList) > 0 {
		for _, val := range userList {
			data.Assignees = append(data.Assignees, &entities.TaskAssignee{
				AssigneeID: val.ID,
			})
		}
	}

	// list attach files
	for i := 1; i <= 20; i++ {
		file, err := services.UploadFile(c, fmt.Sprintf("attach_file_%d", i))
		if err != nil {
			continue
		}

		data.AttachFiles = append(data.AttachFiles, &entities.TaskAttachFile{
			UUID:           file.UUID,
			OriginalName:   file.OriginalName,
			RelativePath:   file.RelativePathFile,
			Ext:            file.FileExt,
			OrganizationID: user.OrganizationID,
			CreatedBy:      user.ID,
		})
	}

	// create task
	_, err = h.Task.Create(ctx, data)
	if err != nil {
		log.For(c).Error("[create-task] create task failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &apis.CreateResponse{
		UUID: data.UUID,
	})
}
