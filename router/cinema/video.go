package cinema

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type VideoRouter struct{}

// InitVideoRouter 初始化 视频表 路由信息
func (s *VideoRouter) InitVideoRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	videoRouter := Router.Group("video").Use(middleware.OperationRecord())
	videoRouterWithoutRecord := Router.Group("video")
	videoRouterWithoutAuth := PublicRouter.Group("video")
	{
		videoRouter.POST("createVideo", videoApi.CreateVideo)             // 新建视频表
		videoRouter.DELETE("deleteVideo", videoApi.DeleteVideo)           // 删除视频表
		videoRouter.DELETE("deleteVideoByIds", videoApi.DeleteVideoByIds) // 批量删除视频表
		videoRouter.PUT("updateVideo", videoApi.UpdateVideo)              // 更新视频表
		videoRouter.POST("publicVideo", videoApi.PublicVideo)             // 发布视频到资源站点
		videoRouter.GET("composeVideo", videoApi.ComposeVideo)            // 将多个自制视频合成到视频表
	}
	{
		videoRouterWithoutRecord.GET("findVideo", videoApi.FindVideo)       // 根据ID获取视频表
		videoRouterWithoutRecord.GET("getVideoList", videoApi.GetVideoList) // 获取视频表列表
	}
	{
		videoRouterWithoutAuth.GET("getVideoPublic", videoApi.GetVideoPublic) // 视频表开放接口
	}
}
