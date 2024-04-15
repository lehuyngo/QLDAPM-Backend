package task_comment

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) ListTaskComment(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	// Only can edit client of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[task-comment-list] query user info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	taskUUID := c.Param("uuid")
	task, err := h.Task.ReadByUUID(ctx, taskUUID)
	if err != nil {
		log.For(c).Error("[task-comment-list] query task info failed", log.Field("user_id", userID), log.Field("task_uuid", taskUUID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.OrganizationID != task.OrganizationID {
		log.For(c).Error("[task-comment-list] query task id isnot match", log.Field("user_id", userID), log.Field("task_uuid", taskUUID),
			log.Field("user_organization_id", user.OrganizationID), log.Field("task_organization_id", task.OrganizationID))

		c.JSON(http.StatusForbidden, err)
		return
	}

	data, err := h.TaskComment.List(ctx, task.ID)
	if err != nil {
		log.For(c).Error("[list-task-comment] query database info failed", log.Field("user_id", userID), log.Err(err))
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListTaskComment{}
	for _, val := range data {
		resp.Data = append(resp.Data, &apis.TaskComment{
			UUID:    val.UUID,
			Content: val.Content,
			Creator: &apis.User{
				UUID:        val.GetCreator().GetUUID(),
				DisplayName: val.GetCreator().GetDisplayName(),
			},
			CreatedTime: val.CreatedAt.UnixMilli(),
		})
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return resp.Data[i].CreatedTime < resp.Data[j].CreatedTime
	})

	c.JSON(http.StatusOK, resp)
}
