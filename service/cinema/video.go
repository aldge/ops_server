package cinema

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cinema"
	cinemaReq "github.com/flipped-aurora/gin-vue-admin/server/model/cinema/request"
)

type VideoService struct{}

// CreateVideo 创建视频表记录
// Author [yourname](https://github.com/yourname)
func (videoService *VideoService) CreateVideo(ctx context.Context, video *cinema.Video) (err error) {
	err = global.GVA_DB.Create(video).Error
	return err
}

// DeleteVideo 删除视频表记录
// Author [yourname](https://github.com/yourname)
func (videoService *VideoService) DeleteVideo(ctx context.Context, vodId string) (err error) {
	err = global.GVA_DB.Delete(&cinema.Video{}, "vod_id = ?", vodId).Error
	return err
}

// DeleteVideoByIds 批量删除视频表记录
// Author [yourname](https://github.com/yourname)
func (videoService *VideoService) DeleteVideoByIds(ctx context.Context, vodIds []string) (err error) {
	err = global.GVA_DB.Delete(&[]cinema.Video{}, "vod_id in ?", vodIds).Error
	return err
}

// UpdateVideo 更新视频表记录
// Author [yourname](https://github.com/yourname)
func (videoService *VideoService) UpdateVideo(ctx context.Context, video cinema.Video) (err error) {
	err = global.GVA_DB.Model(&cinema.Video{}).Where("vod_id = ?", video.VodId).Updates(&video).Error
	return err
}

// GetVideo 根据vodId获取视频表记录
// Author [yourname](https://github.com/yourname)
func (videoService *VideoService) GetVideo(ctx context.Context, vodId string) (video cinema.Video, err error) {
	err = global.GVA_DB.Where("vod_id = ?", vodId).First(&video).Error
	return
}

func (videoService *VideoService) GetVideoByEntID(ctx context.Context, entID string) (video cinema.Video, err error) {
	err = global.GVA_DB.Where("ent_id = ?", entID).First(&video).Error
	return
}

// GetVideosByIds 根据vodIds批量获取视频表记录
// Author [yourname](https://github.com/yourname)
func (videoService *VideoService) GetVideosByIds(ctx context.Context, vodIds []string) (videos []cinema.Video, err error) {
	err = global.GVA_DB.Where("vod_id IN ?", vodIds).Find(&videos).Error
	return
}

// GetVideoInfoList 分页获取视频表记录
// Author [yourname](https://github.com/yourname)
func (videoService *VideoService) GetVideoInfoList(ctx context.Context, info cinemaReq.VideoSearch) (list []cinema.Video, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&cinema.Video{})
	var videos []cinema.Video
	// 如果有条件搜索 下方会自动创建搜索语句

	if info.VodId != nil {
		db = db.Where("vod_id = ?", *info.VodId)
	}
	if info.VodName != nil && *info.VodName != "" {
		db = db.Where("vod_name LIKE ?", "%"+*info.VodName+"%")
	}
	if info.VodClass != nil && *info.VodClass != "" {
		db = db.Where("vod_class LIKE ?", "%"+*info.VodClass+"%")
	}
	if info.App != nil && *info.App != "" {
		db = db.Where("app = ?", *info.App)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 按照 vod_time 倒序排序
	db = db.Order("vod_time DESC")

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&videos).Error
	return videos, total, err
}
func (videoService *VideoService) GetVideoPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

func (videoService *VideoService) CreateCutter(ctx context.Context, cutter *cinema.Cutter) (err error) {
	err = global.GVA_DB.Create(cutter).Error
	return err
}
