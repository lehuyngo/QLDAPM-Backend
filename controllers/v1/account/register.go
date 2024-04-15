package account

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
)

func (h Handler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	req := &apis.Register{}

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if req.OrganizationName == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("organization name is empty"))
		return
	}

	_, err = h.User.ReadByUsername(ctx, req.Username)
	if err == nil {
		c.JSON(http.StatusConflict, &apis.Error{
			Message: &apis.ErrorMessage{
				VI: "User name đã tồn tại",
				EN: "User name is exist",
			},
		})
		return
	}

	data := &entities.User{
		DisplayName:	req.DisplayName,
		Username:		req.Username,
		PasswordHash:	req.Password,
		Email:			req.Email,
	}

	// TGL Solutions
	org, err := h.Organization.First(ctx)
	if err == nil {
		data.OrganizationID = org.ID
		data.Organization = org
	} else {
		data.Organization = &entities.Organization{
			DisplayName: req.OrganizationName,
		}
	}

	token, err := h.User.Register(ctx, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, &apis.AuthenticateResponse{
		Token: token,
	})
}
