package cinema

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cinema"
	cinemaReq "github.com/flipped-aurora/gin-vue-admin/server/model/cinema/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type VideoApi struct{}

// CreateVideo 创建视频表
// @Tags Video
// @Summary 创建视频表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cinema.Video true "创建视频表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /video/createVideo [post]
func (videoApi *VideoApi) CreateVideo(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var video cinema.Video
	err := c.ShouldBindJSON(&video)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = videoService.CreateVideo(ctx, &video)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteVideo 删除视频表
// @Tags Video
// @Summary 删除视频表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cinema.Video true "删除视频表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /video/deleteVideo [delete]
func (videoApi *VideoApi) DeleteVideo(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	vodId := c.Query("vodId")
	err := videoService.DeleteVideo(ctx, vodId)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteVideoByIds 批量删除视频表
// @Tags Video
// @Summary 批量删除视频表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /video/deleteVideoByIds [delete]
func (videoApi *VideoApi) DeleteVideoByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	vodIds := c.QueryArray("vodIds[]")
	err := videoService.DeleteVideoByIds(ctx, vodIds)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateVideo 更新视频表
// @Tags Video
// @Summary 更新视频表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cinema.Video true "更新视频表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /video/updateVideo [put]
func (videoApi *VideoApi) UpdateVideo(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var video cinema.Video
	err := c.ShouldBindJSON(&video)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = videoService.UpdateVideo(ctx, video)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindVideo 用id查询视频表
// @Tags Video
// @Summary 用id查询视频表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param vodId query int true "用id查询视频表"
// @Success 200 {object} response.Response{data=cinema.Video,msg=string} "查询成功"
// @Router /video/findVideo [get]
func (videoApi *VideoApi) FindVideo(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	vodId := c.Query("vodId")
	revideo, err := videoService.GetVideo(ctx, vodId)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(revideo, c)
}

// GetVideoList 分页获取视频表列表
// @Tags Video
// @Summary 分页获取视频表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query cinemaReq.VideoSearch true "分页获取视频表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /video/getVideoList [get]
func (videoApi *VideoApi) GetVideoList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo cinemaReq.VideoSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := videoService.GetVideoInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetVideoPublic 不需要鉴权的视频表接口
// @Tags Video
// @Summary 不需要鉴权的视频表接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /video/getVideoPublic [get]
func (videoApi *VideoApi) GetVideoPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	videoService.GetVideoPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的视频表接口信息",
	}, "获取成功", c)
}

// PublicVideo 将数据保存到资源站点
// @Tags Video
// @Summary 发布视频到资源站点
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param vodIds body []string true "视频ID数组"
// @Success 200 {object} response.Response{msg=string} "发布成功"
// @Router /video/publicVideo [post]
func (videoApi *VideoApi) PublicVideo(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 从JSON body中获取vodIds数组
	var req struct {
		VodIds []string `json:"vodIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("vodIds参数不能为空，且必须为数组格式", c)
		return
	}

	if len(req.VodIds) == 0 {
		response.FailWithMessage("vodIds数组不能为空", c)
		return
	}

	// 检查配置中的stream-provide-save路径
	provideSavePath := global.GVA_CONFIG.Adapter.StreamProvideSave
	if provideSavePath == "" {
		global.GVA_LOG.Error("adapter.stream-provide-save配置为空")
		response.FailWithMessage("adapter.stream-provide-save配置未设置", c)
		return
	}

	// 获取access token
	accessToken, err := utils.GetAccessToken()
	if err != nil {
		global.GVA_LOG.Error("获取access token失败!", zap.Error(err))
		response.FailWithMessage("获取access token失败:"+err.Error(), c)
		return
	}

	// 批量获取所有视频信息
	videos, err := videoService.GetVideosByIds(ctx, req.VodIds)
	if err != nil {
		global.GVA_LOG.Error("批量获取视频信息失败!", zap.Error(err))
		response.FailWithMessage("批量获取视频信息失败:"+err.Error(), c)
		return
	}

	// 检查是否有视频未找到
	foundVodIds := make(map[string]bool)
	for _, video := range videos {
		if video.VodId != nil {
			foundVodIds[fmt.Sprintf("%d", *video.VodId)] = true
		}
	}

	var notFoundVodIds []string
	for _, vodId := range req.VodIds {
		if !foundVodIds[vodId] {
			notFoundVodIds = append(notFoundVodIds, vodId)
		}
	}

	if len(videos) == 0 {
		response.FailWithMessage("未找到任何视频信息", c)
		return
	}

	// 设置请求头，包含 Authorization token
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"accept":        "application/json",
	}

	// 将视频数据转换为下划线格式的 JSON
	videoMaps := make([]map[string]interface{}, 0, len(videos))
	for _, video := range videos {
		// 先序列化为 JSON（驼峰格式）
		videoJSON, err := json.Marshal(video)
		if err != nil {
			global.GVA_LOG.Error("序列化视频数据失败!", zap.Error(err))
			response.FailWithMessage("序列化视频数据失败:"+err.Error(), c)
			return
		}

		// 反序列化为 map
		var videoMap map[string]interface{}
		if err := json.Unmarshal(videoJSON, &videoMap); err != nil {
			global.GVA_LOG.Error("反序列化视频数据失败!", zap.Error(err))
			response.FailWithMessage("反序列化视频数据失败:"+err.Error(), c)
			return
		}

		// 转换为下划线格式的 map
		snakeCaseMap := make(map[string]interface{})
		for key, value := range videoMap {
			// 将驼峰格式转换为下划线格式
			snakeKey := utils.HumpToUnderscore(key)
			snakeCaseMap[snakeKey] = value
		}
		videoMaps = append(videoMaps, snakeCaseMap)
	}

	/*/ 序列化请求体（下划线格式）用于打印 curl 命令和发送请求
	requestBody, err := json.Marshal(videoMaps)
	if err != nil {
		global.GVA_LOG.Error("序列化请求体失败!", zap.Error(err))
		response.FailWithMessage("序列化请求体失败:"+err.Error(), c)
		return
	}

	// 构建 curl 命令用于调试
	curlCmd := fmt.Sprintf("curl -X 'POST'   '%s' ", provideSavePath)
	for key, value := range headers {
		curlCmd += fmt.Sprintf("  -H '%s: %s' ", key, value)
	}
	curlCmd += fmt.Sprintf("  -d '%s'", string(requestBody))

	// 打印 curl 命令
	global.GVA_LOG.Info("发送请求到 adapter 接口", zap.String("curl", curlCmd))*/

	// 批量发送所有视频数据到adapter接口（下划线格式）
	resp, err := request.HttpRequest(
		provideSavePath,
		"POST",
		headers,
		nil,       // 查询参数已在配置的 URL 中包含
		videoMaps, // 发送下划线格式的视频数组
	)

	if err != nil {
		global.GVA_LOG.Error("调用adapter接口失败!", zap.Error(err))
		response.FailWithMessage("调用adapter接口失败:"+err.Error(), c)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		global.GVA_LOG.Error("读取adapter接口响应失败!", zap.Error(err))
		response.FailWithMessage("读取adapter接口响应失败:"+err.Error(), c)
		return
	}

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		global.GVA_LOG.Error("adapter接口返回错误",
			zap.Int("statusCode", resp.StatusCode),
			zap.String("body", string(bodyBytes)))
		response.FailWithMessage(fmt.Sprintf("adapter接口返回错误，状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes)), c)
		return
	}

	// 解析 JSON 响应
	var adapterResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Count  int   `json:"count"`
			VodIds []int `json:"vod_ids"`
		} `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &adapterResp); err != nil {
		global.GVA_LOG.Error("解析adapter接口响应失败!",
			zap.Error(err),
			zap.String("body", string(bodyBytes)))
		response.FailWithMessage("解析adapter接口响应失败:"+err.Error(), c)
		return
	}

	// 检查业务状态码（code: 0 表示成功）
	if adapterResp.Code != 0 {
		global.GVA_LOG.Error("adapter接口返回业务错误",
			zap.Int("code", adapterResp.Code),
			zap.String("message", adapterResp.Message))
		response.FailWithMessage(fmt.Sprintf("adapter接口返回错误: %s (code: %d)", adapterResp.Message, adapterResp.Code), c)
		return
	}

	// 构建响应信息
	message := fmt.Sprintf("发布成功，共处理 %d 个视频", adapterResp.Data.Count)
	if len(notFoundVodIds) > 0 {
		message += fmt.Sprintf("，其中 %d 个视频未找到: %v", len(notFoundVodIds), notFoundVodIds)
	}

	// 返回成功响应，包含 adapter 返回的数据
	responseData := gin.H{
		"count":         adapterResp.Data.Count,
		"vodIds":        adapterResp.Data.VodIds,
		"successCount":  len(videos),
		"notFoundCount": len(notFoundVodIds),
	}
	if len(notFoundVodIds) > 0 {
		responseData["notFoundVodIds"] = notFoundVodIds
	}

	response.OkWithDetailed(responseData, message, c)
}

// setDefaultVideoValues 为视频设置默认值（仅在创建新记录时使用）
func setDefaultVideoValues(v *cinema.Video) {
	// 默认值常量
	var emptyString string = ""
	var zeroInt32 int32 = 0
	var zeroInt64 int64 = 0
	var oneInt64 int64 = 1
	var zeroFloat64 float64 = 0.0
	var falseBool bool = false

	// 字符串类型字段默认值
	if v.App == nil {
		v.App = &emptyString
	}
	if v.EntId == nil {
		v.EntId = &emptyString
	}
	if v.VodName == nil {
		v.VodName = &emptyString
	}
	if v.VodSub == nil {
		v.VodSub = &emptyString
	}
	if v.VodEn == nil {
		v.VodEn = &emptyString
	}
	if v.VodLetter == nil {
		v.VodLetter = &emptyString
	}
	if v.VodColor == nil {
		v.VodColor = &emptyString
	}
	if v.VodTag == nil {
		v.VodTag = &emptyString
	}
	if v.VodClass == nil {
		v.VodClass = &emptyString
	}
	if v.VodPic == nil {
		v.VodPic = &emptyString
	}
	if v.VodPicThumb == nil {
		v.VodPicThumb = &emptyString
	}
	if v.VodPicSlide == nil {
		v.VodPicSlide = &emptyString
	}
	if v.VodPicScreenshot == nil {
		v.VodPicScreenshot = &emptyString
	}
	if v.VodActor == nil {
		v.VodActor = &emptyString
	}
	if v.VodDirector == nil {
		v.VodDirector = &emptyString
	}
	if v.VodWriter == nil {
		v.VodWriter = &emptyString
	}
	if v.VodBehind == nil {
		v.VodBehind = &emptyString
	}
	if v.VodBlurb == nil {
		v.VodBlurb = &emptyString
	}
	if v.VodRemarks == nil {
		v.VodRemarks = &emptyString
	}
	if v.VodPubdate == nil {
		v.VodPubdate = &emptyString
	}
	if v.VodSerial == nil {
		v.VodSerial = &emptyString
	}
	if v.VodTv == nil {
		v.VodTv = &emptyString
	}
	if v.VodWeekday == nil {
		v.VodWeekday = &emptyString
	}
	if v.VodArea == nil {
		v.VodArea = &emptyString
	}
	if v.VodLang == nil {
		v.VodLang = &emptyString
	}
	if v.VodYear == nil {
		v.VodYear = &emptyString
	}
	if v.VodVersion == nil {
		v.VodVersion = &emptyString
	}
	if v.VodState == nil {
		v.VodState = &emptyString
	}
	if v.VodAuthor == nil {
		v.VodAuthor = &emptyString
	}
	if v.VodJumpurl == nil {
		v.VodJumpurl = &emptyString
	}
	if v.VodTpl == nil {
		v.VodTpl = &emptyString
	}
	if v.VodTplPlay == nil {
		v.VodTplPlay = &emptyString
	}
	if v.VodTplDown == nil {
		v.VodTplDown = &emptyString
	}
	if v.VodDuration == nil {
		v.VodDuration = &emptyString
	}
	if v.VodReurl == nil {
		v.VodReurl = &emptyString
	}
	if v.VodRelVod == nil {
		v.VodRelVod = &emptyString
	}
	if v.VodRelArt == nil {
		v.VodRelArt = &emptyString
	}
	if v.VodPwd == nil {
		v.VodPwd = &emptyString
	}
	if v.VodPwdUrl == nil {
		v.VodPwdUrl = &emptyString
	}
	if v.VodPwdPlay == nil {
		v.VodPwdPlay = &emptyString
	}
	if v.VodPwdPlayUrl == nil {
		v.VodPwdPlayUrl = &emptyString
	}
	if v.VodPwdDown == nil {
		v.VodPwdDown = &emptyString
	}
	if v.VodPwdDownUrl == nil {
		v.VodPwdDownUrl = &emptyString
	}
	if v.VodContent == nil {
		v.VodContent = &emptyString
	}
	if v.VodPlayFrom == nil {
		v.VodPlayFrom = &emptyString
	}
	if v.VodPlayServer == nil {
		v.VodPlayServer = &emptyString
	}
	if v.VodPlayNote == nil {
		v.VodPlayNote = &emptyString
	}
	if v.VodPlayUrl == nil {
		v.VodPlayUrl = &emptyString
	}
	if v.VodDownFrom == nil {
		v.VodDownFrom = &emptyString
	}
	if v.VodDownServer == nil {
		v.VodDownServer = &emptyString
	}
	if v.VodDownNote == nil {
		v.VodDownNote = &emptyString
	}
	if v.VodDownUrl == nil {
		v.VodDownUrl = &emptyString
	}
	if v.VodPlotName == nil {
		v.VodPlotName = &emptyString
	}
	if v.VodPlotDetail == nil {
		v.VodPlotDetail = &emptyString
	}

	// 整数类型字段默认值
	if v.TypeId == nil {
		v.TypeId = &zeroInt64
	}
	if v.TypeId1 == nil {
		v.TypeId1 = &zeroInt64
	}
	if v.GroupId == nil {
		v.GroupId = &zeroInt64
	}
	if v.VodTotal == nil {
		v.VodTotal = &oneInt64 // 总集数默认为1
	}
	if v.VodPoints == nil {
		v.VodPoints = &zeroInt64
	}
	if v.VodPointsPlay == nil {
		v.VodPointsPlay = &zeroInt64
	}
	if v.VodPointsDown == nil {
		v.VodPointsDown = &zeroInt64
	}
	if v.VodHits == nil {
		v.VodHits = &zeroInt64
	}
	if v.VodHitsDay == nil {
		v.VodHitsDay = &zeroInt64
	}
	if v.VodHitsWeek == nil {
		v.VodHitsWeek = &zeroInt64
	}
	if v.VodHitsMonth == nil {
		v.VodHitsMonth = &zeroInt64
	}
	if v.VodUp == nil {
		v.VodUp = &zeroInt64
	}
	if v.VodDown == nil {
		v.VodDown = &zeroInt64
	}
	if v.VodScoreAll == nil {
		v.VodScoreAll = &zeroInt64
	}
	if v.VodScoreNum == nil {
		v.VodScoreNum = &zeroInt64
	}
	if v.VodTrysee == nil {
		v.VodTrysee = &zeroInt64
	}
	if v.VodTimeHits == nil {
		v.VodTimeHits = &zeroInt32
	}
	if v.VodTimeMake == nil {
		v.VodTimeMake = &zeroInt32
	}
	if v.VodDoubanId == nil {
		v.VodDoubanId = &zeroInt32
	}

	// 浮点数类型字段默认值
	if v.VodScore == nil {
		v.VodScore = &zeroFloat64
	}
	if v.VodDoubanScore == nil {
		v.VodDoubanScore = &zeroFloat64
	}

	// 布尔类型字段默认值
	if v.VodStatus == nil {
		v.VodStatus = &falseBool
	}
	if v.VodIsend == nil {
		v.VodIsend = &falseBool
	}
	if v.VodLock == nil {
		v.VodLock = &falseBool
	}
	if v.VodLevel == nil {
		v.VodLevel = &falseBool
	}
	if v.VodCopyright == nil {
		v.VodCopyright = &falseBool
	}
	if v.VodPlot == nil {
		v.VodPlot = &falseBool
	}

	// 注意：VodId 是主键，不设置默认值，由数据库或业务逻辑控制
}

// ComposeVideo 将多个自制视频合成到视频表
// @Tags Video
// @Summary 根据 ent_id 创建或更新视频
// @Description 根据传入的 ent_id 查找是否存在视频记录，如果不存在则创建新记录，如果存在则更新记录
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cinema.Video true "视频信息"
// @Success 200 {object} response.Response{msg=string} "操作成功"
// @Failure 400 {object} response.Response{msg=string} "参数错误"
// @Failure 500 {object} response.Response{msg=string} "操作失败"
// @Router /video/composeVideo [post]
func (videoApi *VideoApi) ComposeVideo(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var video cinema.Video
	err := c.ShouldBindJSON(&video)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 检查 EntId 是否为 nil
	if video.EntId == nil || *video.EntId == "" {
		response.FailWithMessage("entId 不能为空", c)
		return
	}

	// 查找是否存在当前 EntId 的视频
	existingVideo, err := videoService.GetVideoByEntID(ctx, *video.EntId)

	// 设置当前时间戳
	now := int32(time.Now().Unix())

	// 判断是否找到记录
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有找到记录，创建新视频
		// 创建时设置添加时间和更新时间
		video.VodTime = &now
		video.VodTimeAdd = &now

		// 设置所有未赋值字段的默认值（仅在创建时使用）
		setDefaultVideoValues(&video)

		err = videoService.CreateVideo(ctx, &video)
		if err != nil {
			global.GVA_LOG.Error("创建视频失败!", zap.Error(err))
			response.FailWithMessage("创建视频失败:"+err.Error(), c)
			return
		}
		response.OkWithMessage("创建视频成功", c)
		return
	} else if err != nil {
		// 查询过程中发生其他错误
		global.GVA_LOG.Error("查询视频失败!", zap.Error(err))
		response.FailWithMessage("查询视频失败:"+err.Error(), c)
		return
	}

	// 找到了记录，进行更新
	// 将找到的 VodId 设置到要更新的 video 中，确保更新正确的记录
	if existingVideo.VodId != nil {
		video.VodId = existingVideo.VodId
	}
	// 更新时只设置更新时间，不改变添加时间
	video.VodTime = &now
	// 保留原有的添加时间
	if existingVideo.VodTimeAdd != nil {
		video.VodTimeAdd = existingVideo.VodTimeAdd
	}

	err = videoService.UpdateVideo(ctx, video)
	if err != nil {
		global.GVA_LOG.Error("更新视频失败!", zap.Error(err))
		response.FailWithMessage("更新视频失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新视频成功", c)
}
