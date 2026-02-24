package cinema

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cinema"
	cinemaReq "github.com/flipped-aurora/gin-vue-admin/server/model/cinema/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AmateurApi struct{}

// CreateAmateur 创建自制视频
// @Tags Amateur
// @Summary 创建自制视频
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cinema.Amateur true "创建自制视频"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /amateur/createAmateur [post]
func (amateurApi *AmateurApi) CreateAmateur(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var amateur cinema.Amateur
	err := c.ShouldBindJSON(&amateur)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 判断前端是否传入了 videoID，如果没有则生成
	if amateur.VideoID == nil || *amateur.VideoID == "" {
		// 生成 videoID
		videoID := generateVideoID(c, amateur.VodName)
		amateur.VideoID = &videoID
	}

	err = amateurService.CreateAmateur(ctx, &amateur)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func generateVideoID(ctx *gin.Context, vodName *string) string {
	var videoID string
	if vodName != nil && *vodName != "" {
		// 根据 VodName 查询数据库中是否有相同的名字
		existingAmateur, err := amateurService.GetAmateurByVodName(ctx.Request.Context(), *vodName)
		if err == nil && existingAmateur.VideoID != nil && len(*existingAmateur.VideoID) >= 6 {
			// 如果找到同名视频，取前6位 + 随机生成后2位
			videoID = (*existingAmateur.VideoID)[:6] + utils.RandomString(2)
		} else {
			// 如果数据库中没有同名的视频，全新生成一个8位ID
			videoID = utils.RandomString(8)
		}
	} else {
		// 如果没有 VodName，直接生成8位ID
		videoID = utils.RandomString(8)
	}
	return videoID
}

// GetVideoID 根据 vodName 生成 VideoID
// @Tags Amateur
// @Summary 根据 vodName 生成 VideoID
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param vod_name query string false "剧名"
// @Success 200 {object} response.Response{data=object,msg=string} "生成成功"
// @Router /amateur/getVideoID [get]
func (amateurApi *AmateurApi) GetVideoID(c *gin.Context) {
	// 从 query 参数获取 vodName
	vodNameStr := c.Query("vod_name")
	var vodName *string
	if vodNameStr != "" {
		vodName = &vodNameStr
	}

	// 生成 videoID
	videoID := generateVideoID(c, vodName)

	// 返回生成的 videoID
	response.OkWithDetailed(gin.H{
		"video_id": videoID,
	}, "生成成功", c)
}

// DeleteAmateur 删除自制视频
// @Tags Amateur
// @Summary 删除自制视频
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cinema.Amateur true "删除自制视频"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /amateur/deleteAmateur [delete]
func (amateurApi *AmateurApi) DeleteAmateur(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := amateurService.DeleteAmateur(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteAmateurByIds 批量删除自制视频
// @Tags Amateur
// @Summary 批量删除自制视频
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /amateur/deleteAmateurByIds [delete]
func (amateurApi *AmateurApi) DeleteAmateurByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := amateurService.DeleteAmateurByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateAmateur 更新自制视频
// @Tags Amateur
// @Summary 更新自制视频
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body cinema.Amateur true "更新自制视频"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /amateur/updateAmateur [put]
func (amateurApi *AmateurApi) UpdateAmateur(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var amateur cinema.Amateur
	err := c.ShouldBindJSON(&amateur)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = amateurService.UpdateAmateur(ctx, amateur)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindAmateur 用id查询自制视频
// @Tags Amateur
// @Summary 用id查询自制视频
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询自制视频"
// @Success 200 {object} response.Response{data=cinema.Amateur,msg=string} "查询成功"
// @Router /amateur/findAmateur [get]
func (amateurApi *AmateurApi) FindAmateur(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	reamateur, err := amateurService.GetAmateur(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(reamateur, c)
}

// GetAmateurList 分页获取自制视频列表
// @Tags Amateur
// @Summary 分页获取自制视频列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query cinemaReq.AmateurSearch true "分页获取自制视频列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /amateur/getAmateurList [get]
func (amateurApi *AmateurApi) GetAmateurList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo cinemaReq.AmateurSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := amateurService.GetAmateurInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:          list,
		Total:         total,
		Page:          pageInfo.Page,
		PageSize:      pageInfo.PageSize,
		CutterPlayUrl: global.GVA_CONFIG.System.CutterPlay,
		StreamPlayUrl: global.GVA_CONFIG.System.StreamPlay,
	}, "获取成功", c)
}

// GetAmateurPublic 不需要鉴权的自制视频接口
// @Tags Amateur
// @Summary 不需要鉴权的自制视频接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /amateur/getAmateurPublic [get]
func (amateurApi *AmateurApi) GetAmateurPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	amateurService.GetAmateurPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的自制视频接口信息",
	}, "获取成功", c)
}

// PublicChips 发布视频切片数据到adapter接口
// @Tags Amateur
// @Summary 发布视频切片数据到adapter接口
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param video_id query string true "视频ID"
// @Success 200 {object} response.Response{msg=string} "发布成功"
// @Router /amateur/publicChips [post]
func (amateurApi *AmateurApi) PublicChips(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	videoId := c.Query("video_id")
	if videoId == "" {
		response.FailWithMessage("video_id参数不能为空", c)
		return
	}

	// 检查配置中的stream-ts-save路径
	tsSavePath := global.GVA_CONFIG.Adapter.StreamTsSave
	if tsSavePath == "" {
		global.GVA_LOG.Error("adapter.stream-ts-save配置为空")
		response.FailWithMessage("adapter.stream-ts-save配置未设置", c)
		return
	}

	// 获取access token
	accessToken, err := utils.GetAccessToken()
	if err != nil {
		global.GVA_LOG.Error("获取access token失败!", zap.Error(err))
		response.FailWithMessage("获取access token失败:"+err.Error(), c)
		return
	}

	// 获取 key 与 iv
	secret, err := amateurService.GetSecret(ctx, videoId)
	if err != nil {
		global.GVA_LOG.Error("获取密钥信息失败!", zap.Error(err), zap.String("videoId", videoId))
		response.FailWithMessage("获取密钥信息失败:"+err.Error(), c)
		return
	}

	// 根据 videoId 获取所有切片
	chips, err := amateurService.GetChips(ctx, videoId)
	if err != nil {
		global.GVA_LOG.Error("获取切片信息失败!", zap.Error(err), zap.String("videoId", videoId))
		response.FailWithMessage("获取切片信息失败:"+err.Error(), c)
		return
	}

	if len(chips) == 0 {
		response.FailWithMessage("未找到任何切片数据", c)
		return
	}

	// 构建ts_data数组，使用JSON序列化来获取正确的字段名（已经是下划线格式）
	type TsDataItem struct {
		TsSequence uint    `json:"ts_sequence"`
		TsPath     string  `json:"ts_path"`
		Duration   float64 `json:"duration"`
		Definition string  `json:"definition"`
	}

	tsDataItems := make([]TsDataItem, 0, len(chips))
	for _, chip := range chips {
		tsDataItems = append(tsDataItems, TsDataItem{
			TsSequence: chip.TsSequence,
			TsPath:     chip.TsPath,
			Duration:   chip.Duration,
			Definition: chip.Definition,
		})
	}

	// 构建请求体结构（字段名通过JSON标签已经是下划线格式）
	type RequestBody struct {
		VideoID string       `json:"video_id"`
		Key     string       `json:"key"`
		Iv      string       `json:"iv"`
		TsData  []TsDataItem `json:"ts_data"`
	}

	requestBody := RequestBody{
		VideoID: videoId,
		Key:     secret.Key,
		Iv:      secret.Iv,
		TsData:  tsDataItems,
	}

	// 序列化为JSON（字段名已经是下划线格式）
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		global.GVA_LOG.Error("序列化请求体失败!", zap.Error(err))
		response.FailWithMessage("序列化请求体失败:"+err.Error(), c)
		return
	}

	// 反序列化为map，字段名已经是下划线格式
	var snakeCaseMap map[string]interface{}
	if err := json.Unmarshal(requestJSON, &snakeCaseMap); err != nil {
		global.GVA_LOG.Error("反序列化请求体失败!", zap.Error(err))
		response.FailWithMessage("反序列化请求体失败:"+err.Error(), c)
		return
	}

	// 设置请求头，包含 Authorization token
	headers := map[string]string{
		"Authorization": "Bearer " + accessToken,
		"accept":        "application/json",
	}

	// 发送数据到adapter接口（下划线格式）
	resp, err := request.HttpRequest(
		tsSavePath,
		"POST",
		headers,
		nil,          // 查询参数已在配置的 URL 中包含
		snakeCaseMap, // 发送下划线格式的数据
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
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
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

	// 更新视频状态为3（发布）
	status := int64(3)
	err = amateurService.UpdateAmateurStatusByVideoID(ctx, videoId, status)
	if err != nil {
		global.GVA_LOG.Error("更新视频状态失败!", zap.Error(err), zap.String("videoId", videoId))
		// 注意：这里不返回错误，因为adapter接口已经成功，状态更新失败不应该影响整体流程
		// 但会记录日志以便后续排查
	}

	// 返回成功响应
	response.OkWithDetailed(gin.H{
		"data": adapterResp.Data,
	}, "发布成功", c)
}
