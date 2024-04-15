package contact

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) ListContact(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)
	user, err := h.User.Read(ctx, userID)
	if err != nil {
		log.For(c).Error("[list-contact] query user info failed", log.Field("user_id", userID), log.Err(err))

		c.JSON(http.StatusInternalServerError, err)
		return
	}

	log.For(c).Info("[list-contact] start process", log.Field("user_id", userID))

	data, err := h.Contact.ListByOrgID(ctx, user.OrganizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	resp := &apis.ListContactResponse{}
	for _, val := range data {
		var avatar *apis.File
		if val.Avatar != nil {
			avatar = &apis.File {
				UUID: val.GetAvatar().GetUUID(),
				URL: fmt.Sprintf("%sapi/v1/contacts/%s/avatar/%s", services.Config.Domain.Domain, val.UUID, val.GetAvatar().GetUUID()),
			}
		}

		var nameCard *apis.File
		if val.NameCard != nil {
			nameCard = &apis.File {
				UUID: val.GetNameCard().GetUUID(),
				URL: fmt.Sprintf("%sapi/v1/contacts/%s/namecard/%s", services.Config.Domain.Domain, val.UUID, val.GetNameCard().GetUUID()),
			}
		}

		contact := &apis.Contact{
			UUID: val.UUID,
			FullName: val.FullName,
			ShortName: val.ShortName,
			Phone: val.Phone,
			Email: val.Email,
			JobTitle: val.JobTitle,
			Gender: val.Gender.Value(),
			Avatar: avatar,
			NameCard: nameCard,
			BirthDay: val.Birthday,
			LastActiveTime: val.LastActiveTime,
		}

		for _, tag := range val.Tags {
			contact.Tags = append(contact.Tags, &apis.ContactTag{
				UUID: tag.GetTag().GetUUID(),
				Name: tag.GetTag().GetName(),
				Color: tag.GetTag().GetColor(),
			})
		}

		for _, client := range val.Clients {
			if client.GetClient() == nil {
				log.For(c).Warn("[list-contact] client is nil", log.Field("user_id", userID))
				continue
			}
			contact.Clients = append(contact.Clients, &apis.Client{
				UUID: client.GetClient().GetUUID(),
				FullName: client.GetClient().GetFullName(),
				ShortName: client.GetClient().GetShortName(),
				Code: client.GetClient().GetCode(),
			})
		}

		resp.Data = append(resp.Data, contact)
	}
	sort.Slice(resp.Data, func(i, j int) bool {
		return strings.EqualFold(resp.Data[i].FullName, resp.Data[j].FullName)
	})

	log.For(c).Info("[list-contact] process success", log.Field("user_id", userID), log.Field("resp_len", len(resp.Data)))
	c.JSON(http.StatusOK, resp)
}
