package contact

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) ReadContact(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	data, err := h.Contact.ReadByUUID(ctx, uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	// Only can edit contact of org
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if user.GetOrganization().GetID() != data.OrganizationID {
		c.JSON(http.StatusForbidden, err)
		return
	}

	// Avatar
	var avatar *apis.File
	if data.Avatar != nil {
		avatar = &apis.File {
			UUID: data.GetAvatar().GetUUID(),
			URL: fmt.Sprintf("%sapi/v1/contacts/%s/avatar/%s", services.Config.Domain.Domain, uuid, data.GetAvatar().GetUUID()),
		}
	}

	// Name Card
	var nameCard *apis.File
	if data.NameCard != nil {
		nameCard = &apis.File {
			UUID: data.GetNameCard().GetUUID(),
			URL: fmt.Sprintf("%sapi/v1/contacts/%s/namecard/%s", services.Config.Domain.Domain, uuid, data.GetNameCard().GetUUID()),
		}
	}

	resp := &apis.Contact{
		UUID: data.UUID,
		FullName: data.FullName,
		ShortName: data.ShortName,
		JobTitle: data.JobTitle,
		Gender: data.Gender.Value(),
		Phone: data.Phone,
		Email: data.Email,
		Avatar: avatar,
		NameCard: nameCard,
		LastActiveTime: data.LastActiveTime,
		CreatedTime: data.GetCreatedAt().UnixMilli(),
		BirthDay: data.Birthday,
	}

	// Notes
	for _, val := range data.Notes {
		resp.Notes = append(resp.Notes, &apis.ContactNote{
			UUID: val.GetUUID(),
			Title: val.Title,
			Content: val.Content,
			Color: val.Color,
		})
	}

	// Tags
	for _, val := range data.Tags {
		resp.Tags = append(resp.Tags, &apis.ContactTag{
			UUID: val.GetTag().GetUUID(),
			Name: val.GetTag().GetName(),
			CreatedTime: val.GetCreatedAt().UnixMilli(),
			Color: val.GetTag().GetColor(),
		})
	}

	c.JSON(http.StatusOK, resp)
}
