package example

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/example"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cinema"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	exampleRes "github.com/flipped-aurora/gin-vue-admin/server/model/example/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/upload"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BreakpointContinue
// @Tags      ExaFileUploadAndDownload
// @Summary   断点续传到服务器
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file                           true  "an example for breakpoint resume, 断点续传示例"
// @Success   200   {object}  response.Response{msg=string}  "断点续传到服务器"
// @Router    /fileUploadAndDownload/breakpointContinue [post]
func (b *FileUploadAndDownloadApi) BreakpointContinue(c *gin.Context) {
	fileMd5 := c.Request.FormValue("fileMd5")
	fileName := c.Request.FormValue("fileName")
	chunkMd5 := c.Request.FormValue("chunkMd5")
	chunkNumber, _ := strconv.Atoi(c.Request.FormValue("chunkNumber"))
	chunkTotal, _ := strconv.Atoi(c.Request.FormValue("chunkTotal"))
	_, FileHeader, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	f, err := FileHeader.Open()
	if err != nil {
		global.GVA_LOG.Error("文件读取失败!", zap.Error(err))
		response.FailWithMessage("文件读取失败", c)
		return
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)
	cen, _ := io.ReadAll(f)
	if !utils.CheckMd5(cen, chunkMd5) {
		global.GVA_LOG.Error("检查md5失败!", zap.Error(err))
		response.FailWithMessage("检查md5失败", c)
		return
	}
	file, err := fileUploadAndDownloadService.FindOrCreateFile(fileMd5, fileName, chunkTotal)
	if err != nil {
		global.GVA_LOG.Error("查找或创建记录失败!", zap.Error(err))
		response.FailWithMessage("查找或创建记录失败", c)
		return
	}
	pathC, err := utils.BreakPointContinue(cen, fileName, chunkNumber, chunkTotal, fileMd5)
	if err != nil {
		global.GVA_LOG.Error("断点续传失败!", zap.Error(err))
		response.FailWithMessage("断点续传失败", c)
		return
	}

	if err = fileUploadAndDownloadService.CreateFileChunk(file.ID, pathC, chunkNumber); err != nil {
		global.GVA_LOG.Error("创建文件记录失败!", zap.Error(err))
		response.FailWithMessage("创建文件记录失败", c)
		return
	}
	response.OkWithMessage("切片创建成功", c)
}

// FindFile
// @Tags      ExaFileUploadAndDownload
// @Summary   查找文件
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file                                                        true  "Find the file, 查找文件"
// @Success   200   {object}  response.Response{data=exampleRes.FileResponse,msg=string}  "查找文件,返回包括文件详情"
// @Router    /fileUploadAndDownload/findFile [get]
func (b *FileUploadAndDownloadApi) FindFile(c *gin.Context) {
	fileMd5 := c.Query("fileMd5")
	fileName := c.Query("fileName")
	chunkTotal, _ := strconv.Atoi(c.Query("chunkTotal"))
	file, err := fileUploadAndDownloadService.FindOrCreateFile(fileMd5, fileName, chunkTotal)
	if err != nil {
		global.GVA_LOG.Error("查找失败!", zap.Error(err))
		response.FailWithMessage("查找失败", c)
	} else {
		response.OkWithDetailed(exampleRes.FileResponse{File: file}, "查找成功", c)
	}
}

// BreakpointContinueFinish
// @Tags      ExaFileUploadAndDownload
// @Summary   创建文件
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file                                                            true  "上传文件完成"
// @Success   200   {object}  response.Response{data=exampleRes.FilePathResponse,msg=string}  "创建文件,返回包括文件路径"
// @Router    /fileUploadAndDownload/findFile [post]
func (b *FileUploadAndDownloadApi) BreakpointContinueFinish(c *gin.Context) {
	fileMd5 := c.Query("fileMd5")
	fileName := c.Query("fileName")
	videoID := c.Query("video_id")
	vodName := c.Query("vod_name")
	vodSerial := c.Query("vod_serial")

	// 先在本地合成文件
	localFilePath, err := utils.MakeFile(fileName, fileMd5)
	if err != nil {
		global.GVA_LOG.Error("文件创建失败!", zap.Error(err))
		response.FailWithDetailed(exampleRes.FilePathResponse{FilePath: localFilePath}, "文件创建失败", c)
		return
	}

	// 如果配置使用 minio，则上传到 minio
	finalFilePath := localFilePath
	if global.GVA_CONFIG.System.OssType == "minio" {
		// 获取 Minio 客户端实例
		minioClient, err := upload.GetMinio(
			global.GVA_CONFIG.Minio.Endpoint,
			global.GVA_CONFIG.Minio.AccessKeyId,
			global.GVA_CONFIG.Minio.AccessKeySecret,
			global.GVA_CONFIG.Minio.BucketName,
			global.GVA_CONFIG.Minio.UseSSL,
		)
		if err != nil {
			global.GVA_LOG.Error("获取Minio客户端失败!", zap.Error(err))
			response.FailWithDetailed(exampleRes.FilePathResponse{FilePath: localFilePath}, "获取Minio客户端失败", c)
			return
		}
		// 上传到 minio
		url, key, uploadErr := minioClient.UploadFileFromPath(localFilePath, fileName, videoID)
		if uploadErr != nil {
			global.GVA_LOG.Error("上传文件到minio失败!", zap.Error(uploadErr))
			response.FailWithDetailed(exampleRes.FilePathResponse{FilePath: localFilePath}, "文件上传到minio失败", c)
			return
		}
		finalFilePath = url
		global.GVA_LOG.Info("文件已上传到minio", zap.String("url", url), zap.String("key", key))

		// 上传成功后删除本地文件
		if err := os.Remove(localFilePath); err != nil {
			global.GVA_LOG.Warn("删除本地文件失败", zap.String("filePath", localFilePath), zap.Error(err))
			// 不返回错误，因为文件已经成功上传到 minio
		} else {
			global.GVA_LOG.Info("本地文件已删除", zap.String("filePath", localFilePath))
		}
	}

	title := fmt.Sprintf("%s第%s集", vodName, vodSerial)
	// 在 cutter 表中添加一条记录
	if videoID != "" {
		now := time.Now()
		cutter := cinema.Cutter{
			VideoID:   &videoID,
			Status:    1,
			Code:      "",
			Tags:      "",
			Title:     title,
			Params:    "",
			Output:    "",
			CreatedAt: now,
			UpdatedAt: &now,
		}
		if err := global.GVA_DB.Create(&cutter).Error; err != nil {
			global.GVA_LOG.Error("创建cutter记录失败!", zap.Error(err))
			// 不返回错误，因为文件已经成功创建
		} else {
			global.GVA_LOG.Info("cutter记录创建成功", zap.String("videoID", videoID))
		}
	}

	response.OkWithDetailed(exampleRes.FilePathResponse{FilePath: finalFilePath}, "文件创建成功", c)
}

// RemoveChunk
// @Tags      ExaFileUploadAndDownload
// @Summary   删除切片
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     file  formData  file                           true  "删除缓存切片"
// @Success   200   {object}  response.Response{msg=string}  "删除切片"
// @Router    /fileUploadAndDownload/removeChunk [post]
func (b *FileUploadAndDownloadApi) RemoveChunk(c *gin.Context) {
	var file example.ExaFile
	err := c.ShouldBindJSON(&file)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 路径穿越拦截：拒绝包含 "../" 或 "..\\" 的路径（真正的路径穿越攻击）
	if strings.Contains(file.FilePath, "../") || strings.Contains(file.FilePath, "..\\") {
		response.FailWithMessage("非法路径，禁止删除", c)
		return
	}

	// 检查是否是 URL 路径（minio 或其他 OSS）
	isURL := strings.HasPrefix(file.FilePath, "http://") || strings.HasPrefix(file.FilePath, "https://")

	if isURL {
		// 如果是 URL，验证是否是配置的 minio bucket URL
		if global.GVA_CONFIG.System.OssType == "minio" && global.GVA_CONFIG.Minio.BucketUrl != "" {
			// 验证 URL 是否以配置的 bucket URL 开头
			if !strings.HasPrefix(file.FilePath, global.GVA_CONFIG.Minio.BucketUrl) {
				response.FailWithMessage("非法路径，不能删除", c)
				return
			}
		} else {
			// 如果配置了 minio 但 URL 不匹配，或者未配置 minio 但提供了 URL，拒绝
			response.FailWithMessage("非法路径，不能删除", c)
			return
		}
	} else {
		// 本地文件路径验证：确保路径在允许的目录范围内（只允许 fileDir 目录下的文件）
		if !strings.HasPrefix(file.FilePath, "./fileDir/") && !strings.HasPrefix(file.FilePath, "fileDir/") {
			response.FailWithMessage("非法路径，不能删除", c)
			return
		}
	}
	err = utils.RemoveChunk(file.FileMd5)
	if err != nil {
		global.GVA_LOG.Error("缓存切片删除失败!", zap.Error(err))
		return
	}
	err = fileUploadAndDownloadService.DeleteFileChunk(file.FileMd5, file.FilePath)
	if err != nil {
		global.GVA_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("缓存切片删除成功", c)
}
