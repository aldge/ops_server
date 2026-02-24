package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type VideoSearch struct {
	VodId    *int    `json:"vodId" form:"vodId"`
	VodName  *string `json:"vodName" form:"vodName"`
	VodClass *string `json:"vodClass" form:"vodClass"`
	App      *string `json:"app" form:"app"`
	request.PageInfo
}
