package cinema

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct {
	VideoRouter
	AmateurRouter
}

var (
	videoApi   = api.ApiGroupApp.CinemaApiGroup.VideoApi
	amateurApi = api.ApiGroupApp.CinemaApiGroup.AmateurApi
)
