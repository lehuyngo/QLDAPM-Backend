package account

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/apis"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/http_parser"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/log"
)

func (h Handler) Authenticate(c *gin.Context) {
	fmt.Printf("[Authenticate] start process ========== \n")
	ctx := c.Request.Context()
	req := &apis.AuthenticateRequest{}

	log.For(c).Debug("[Authenticate] start process")

	err := http_parser.BindAndValid(c, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		log.For(c).Error("[Authenticate] parse req failed")
		return
	}

	token, err := h.User.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		log.For(c).Error("[Authenticate] authen token failed")
		return
	}

	log.For(c).Info("[Authenticate] process success")
	c.JSON(http.StatusOK, &apis.AuthenticateResponse{
		Token: token,
	})
}
