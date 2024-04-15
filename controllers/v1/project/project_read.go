package project

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
)

func (h Handler) ReadProject(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.Project.ReadByUUID(ctx, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit project of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.GetOrganization().GetID() != data.OrganizationID {
		c.JSON(http.StatusForbidden, err)
		return
	}

	resp := &apis.Project{
		UUID:           data.UUID,
		FullName:       data.FullName,
		ShortName:      data.ShortName,
		Code:           data.GetCode(),
		ProjectStatus:  data.ProjectStatus.Value(),
		LastActiveTime: data.LastActiveTime,
		CreatedTime:    data.GetCreatedAt().UnixMilli(),
	}
	if data.GetClient() != nil {
		resp.Client = &apis.Client{
			UUID:        data.GetClient().GetUUID(),
			FullName:    data.GetClient().GetFullName(),
			ShortName:   data.GetClient().GetShortName(),
			Code:        data.GetClient().GetCode(),
			Fax:         data.GetClient().GetFax(),
			Website:     data.GetClient().GetWebsite(),
			Phone:       data.GetClient().GetPhone(),
			Email:       data.GetClient().GetEmail(),
			Address:     data.GetClient().GetAddress(),
			CreatedTime: data.CreatedAt.UnixMilli(),
		}
	}

	for _, contact := range data.Contacts {
		resp.Contacts = append(resp.Contacts, &apis.Contact{
			UUID:      contact.GetContact().GetUUID(),
			FullName:  contact.GetContact().GetFullName(),
			ShortName: contact.GetContact().GetShortName(),
			Phone:     contact.GetContact().GetPhone(),
			Email:     contact.GetContact().GetEmail(),
		})
	}

	meetings, err := h.Meeting.List(ctx, data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// get lastest meeting created time
	if len(meetings) > 0 {
		sort.Slice(meetings, func(i, j int) bool {
			return meetings[i].GetCreatedAt().UnixMilli() > meetings[j].GetCreatedAt().UnixMilli()
		})
		resp.LastMeetingCreatedTime = meetings[0].GetCreatedAt().UnixMilli()
	}

	c.JSON(http.StatusOK, resp)
}
