package cinema

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	VideoApi
	AmateurApi
}

var (
	videoService   = service.ServiceGroupApp.CinemaServiceGroup.VideoService
	amateurService = service.ServiceGroupApp.CinemaServiceGroup.AmateurService
)
