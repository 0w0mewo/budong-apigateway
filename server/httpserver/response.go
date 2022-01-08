package httpserver

import (
	"net/http"

	"github.com/0w0mewo/budong-apigateway/server"
	"github.com/gin-gonic/gin"
)

type ErrResp struct {
	ErrMsg string `json:"error"`
}

type CountResp struct {
	Count uint64 `json:"amount"`
}

type InventoryResp struct {
	ErrResp
	Infos []*server.SetuInfo `json:"inventory"`
}

func tryOrSendErr(c *gin.Context, try func() error) {
	ersp := &ErrResp{}

	err := try()
	if err != nil {
		ersp.ErrMsg = err.Error()
		c.JSON(http.StatusBadRequest, ersp)

		return
	}

}
