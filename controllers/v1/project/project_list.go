package project

import (
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
)

func (h Handler) ListProject(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	data, err := h.Project.ListByOrgID(ctx, user.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListProjectResponse{}
	for _, val := range data {
		project := &apis.Project{
			UUID:           val.UUID,
			FullName:       val.FullName,
			ShortName:      val.ShortName,
			Code:           val.Code,
			ProjectStatus:  val.ProjectStatus.Value(),
			LastActiveTime: val.LastActiveTime,
			CreatedTime:    val.CreatedAt.UnixMilli(),
		}
		if val.GetClient() != nil {
			if val.GetClient().Available() {
				project.Client = &apis.Client{
					UUID:        val.GetClient().GetUUID(),
					FullName:    val.GetClient().GetFullName(),
					ShortName:   val.GetClient().GetShortName(),
					Code:        val.GetClient().GetCode(),
					Fax:         val.GetClient().GetFax(),
					Website:     val.GetClient().GetWebsite(),
					Phone:       val.GetClient().GetPhone(),
					Email:       val.GetClient().GetEmail(),
					Address:     val.GetClient().GetAddress(),
					CreatedTime: val.CreatedAt.UnixMilli(),
				}
			}
		}

		for _, contact := range val.Contacts {
			if contact.GetContact().Available() {
				project.Contacts = append(project.Contacts, &apis.Contact{
					UUID:      contact.GetContact().GetUUID(),
					FullName:  contact.GetContact().GetFullName(),
					ShortName: contact.GetContact().GetShortName(),
					Phone:     contact.GetContact().GetPhone(),
					Email:     contact.GetContact().GetEmail(),
				})
			}
		}

		meetings, err := h.Meeting.List(ctx, val.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		// get lastest meeting created time
		if len(meetings) > 0 {
			sort.Slice(meetings, func(i, j int) bool {
				return meetings[i].GetCreatedAt().UnixMilli() > meetings[j].GetCreatedAt().UnixMilli()
			})
			project.LastMeetingCreatedTime = meetings[0].GetCreatedAt().UnixMilli()
		}

		resp.Data = append(resp.Data, project)
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return strings.EqualFold(resp.Data[i].FullName, resp.Data[j].FullName)
	})

	c.JSON(http.StatusOK, resp)
}
