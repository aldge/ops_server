package cinema

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type AmateurRouter struct{}

// InitAmateurRouter 初始化 自制视频 路由信息
func (s *AmateurRouter) InitAmateurRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	amateurRouter := Router.Group("amateur").Use(middleware.OperationRecord())
	amateurRouterWithoutRecord := Router.Group("amateur")
	amateurRouterWithoutAuth := PublicRouter.Group("amateur")
	{
		amateurRouter.POST("createAmateur", amateurApi.CreateAmateur)             // 新建自制视频
		amateurRouter.DELETE("deleteAmateur", amateurApi.DeleteAmateur)           // 删除自制视频
		amateurRouter.DELETE("deleteAmateurByIds", amateurApi.DeleteAmateurByIds) // 批量删除自制视频
		amateurRouter.PUT("updateAmateur", amateurApi.UpdateAmateur)              // 更新自制视频
		amateurRouter.GET("getVideoID", amateurApi.GetVideoID)                    // 获取视频ID
		amateurRouter.GET("publicChips", amateurApi.PublicChips)                  // 获取视频ID
	}
	{
		amateurRouterWithoutRecord.GET("findAmateur", amateurApi.FindAmateur)       // 根据ID获取自制视频
		amateurRouterWithoutRecord.GET("getAmateurList", amateurApi.GetAmateurList) // 获取自制视频列表
	}
	{
		amateurRouterWithoutAuth.GET("getAmateurPublic", amateurApi.GetAmateurPublic) // 自制视频开放接口
	}
}
