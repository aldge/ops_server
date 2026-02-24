package cinema

import (
	"context"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cinema"
	cinemaReq "github.com/flipped-aurora/gin-vue-admin/server/model/cinema/request"
)

type AmateurService struct{}

// CreateAmateur 创建自制视频记录
// Author [yourname](https://github.com/yourname)
func (amateurService *AmateurService) CreateAmateur(ctx context.Context, amateur *cinema.Amateur) (err error) {
	err = global.GVA_DB.Create(amateur).Error
	return err
}

// DeleteAmateur 删除自制视频记录
// Author [yourname](https://github.com/yourname)
func (amateurService *AmateurService) DeleteAmateur(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&cinema.Amateur{}, "id = ?", ID).Error
	return err
}

// DeleteAmateurByIds 批量删除自制视频记录
// Author [yourname](https://github.com/yourname)
func (amateurService *AmateurService) DeleteAmateurByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]cinema.Amateur{}, "id in ?", IDs).Error
	return err
}

// UpdateAmateur 更新自制视频记录
// Author [yourname](https://github.com/yourname)
func (amateurService *AmateurService) UpdateAmateur(ctx context.Context, amateur cinema.Amateur) (err error) {
	err = global.GVA_DB.Model(&cinema.Amateur{}).Where("id = ?", amateur.ID).Updates(&amateur).Error
	return err
}

// GetAmateur 根据ID获取自制视频记录
// Author [yourname](https://github.com/yourname)
func (amateurService *AmateurService) GetAmateur(ctx context.Context, ID string) (amateur cinema.Amateur, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&amateur).Error
	return
}

// GetAmateurByVodName 根据VodName获取第一个匹配的自制视频记录
// Author [yourname](https://github.com/yourname)
func (amateurService *AmateurService) GetAmateurByVodName(ctx context.Context, vodName string) (amateur cinema.Amateur, err error) {
	err = global.GVA_DB.Where("vod_name = ?", vodName).First(&amateur).Error
	return
}

// GetAmateurInfoList 分页获取自制视频记录
// Author [yourname](https://github.com/yourname)
func (amateurService *AmateurService) GetAmateurInfoList(ctx context.Context, info cinemaReq.AmateurSearch) (list []cinema.Amateur, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&cinema.Amateur{})
	var amateurs []cinema.Amateur
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}
	if info.App != nil && *info.App != "" {
		db = db.Where("app = ?", *info.App)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&amateurs).Error
	return amateurs, total, err
}
func (amateurService *AmateurService) GetAmateurPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}

func (amateurService *AmateurService) GetChips(ctx context.Context, VideoId string) (chips []cinema.Chips, err error) {
	err = global.GVA_DB.Where("video_id = ?", VideoId).Order("ts_sequence ASC").Find(&chips).Error
	return
}

func (amateurService *AmateurService) GetSecret(ctx context.Context, VideoId string) (secret cinema.Secrets, err error) {
	err = global.GVA_DB.Where("video_id = ?", VideoId).First(&secret).Error
	return
}

// UpdateAmateurStatusByVideoID 根据video_id更新视频状态
func (amateurService *AmateurService) UpdateAmateurStatusByVideoID(ctx context.Context, videoId string, status int64) (err error) {
	statusPtr := &status
	err = global.GVA_DB.Model(&cinema.Amateur{}).Where("video_id = ?", videoId).Update("status", statusPtr).Error
	return
}
