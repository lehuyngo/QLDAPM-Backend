package contact_attach_file

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/middlewares"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/normalize"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/services"
)

func (h Handler) ListContactAttachFile(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _, _ := middlewares.ParseToken(c)

	uuid := c.Param("uuid")
	contact, err := h.Contact.ReadByUUID(ctx, uuid)
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

	if user.OrganizationID != contact.OrganizationID {
		c.JSON(http.StatusForbidden, err)
		return
	}

	resp := &apis.ListFileResponse{}
	// Attach files
	files, err := h.ContactAttachFile.List(ctx, contact.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	for _, val := range files {
		fileName := normalize.URLEncode(val.GetOriginalName()) + val.Ext
		resp.Data = append(resp.Data, &apis.File{
			Name: val.GetOriginalName() + val.Ext,
			UUID: val.GetUUID(),
			URL: fmt.Sprintf("%sapi/v1/contacts/%s/downloaded-files/%s/%s", services.Config.Domain.Domain, contact.GetUUID(), val.GetUUID(), fileName),
		})
	}

	c.JSON(http.StatusOK, resp)
}
