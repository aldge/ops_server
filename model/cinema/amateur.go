// 自动生成模板Amateur
package cinema

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 自制视频 结构体  Amateur
type Amateur struct {
	global.GVA_MODEL
	App        *string `json:"app" form:"app" gorm:"column:app;"`                            //应用
	VideoID    *string `json:"video_id" form:"video_id" gorm:"primarykey;column:video_id;"`  //视频ID
	VodName    *string `json:"vod_name" form:"vod_name" gorm:"column:vod_name;"`             //剧名
	VodLang    *string `json:"vod_lang" form:"vod_lang" gorm:"column:vod_lang;"`             //语言
	VodTotal   *string `json:"vod_total" form:"vod_total" gorm:"column:vod_total;"`          //总集数
	VodSerial  *int64  `json:"vod_serial" form:"vod_serial" gorm:"column:vod_serial;"`       //序列号
	VodCutters *int64  `json:"vod_cutters" form:"vod_cutters" gorm:"column:vod_cutters;"`    //切片数量
	Status     *int64  `json:"status" form:"status" gorm:"column:status;"`                   //状态 0:删除 1:草稿 2:切片 3:发布
	VodTime    *int64  `json:"vod_time" form:"vod_time" gorm:"column:vod_time;"`             //更新时间
	VodTimeAdd *int64  `json:"vod_time_add" form:"vod_time_add" gorm:"column:vod_time_add;"` //添加时间
}

// TableName 自制视频 Amateur自定义表名 cine_amateur
func (Amateur) TableName() string {
	return "cine_amateur"
}

// 切片表 结构体  Cutter
// 状态 0:已发布, 1:队列中, 2:处理中, 3:发布失败, 4:已完成;
type Cutter struct {
	ID        int        `json:"id" form:"id" gorm:"column:id;"`                         //主键ID
	VideoID   *string    `json:"video_id" form:"video_id" gorm:"column:video_id;"`       //视频ID
	Code      string     `json:"code" form:"code" gorm:"column:code;"`                   //代码
	Tags      string     `json:"tags" form:"tags" gorm:"column:tags;"`                   //标签
	Title     string     `json:"title" form:"title" gorm:"column:title;"`                //标题
	Params    string     `json:"params" form:"params" gorm:"column:params"`              //参数
	Output    string     `json:"output" form:"output" gorm:"column:output;"`             //输出
	Status    int        `json:"status" form:"status" gorm:"column:status;"`             // 状态
	CreatedAt time.Time  `json:"created_at" form:"created_at" gorm:"column:created_at;"` //创建时间
	UpdatedAt *time.Time `json:"updated_at" form:"updated_at" gorm:"column:updated_at;"` //更新时间
}

// TableName 切片表 Cutter自定义表名 cutter
func (Cutter) TableName() string {
	return "cutter_videos"
}

// 视频ts文件表 结构体  Chips
type Chips struct {
	ID         uint64  `json:"id" form:"id" gorm:"column:id;"`                            //主键id
	VideoID    string  `json:"video_id" form:"video_id" gorm:"column:video_id;"`          //视频id
	TsSequence uint    `json:"ts_sequence" form:"ts_sequence" gorm:"column:ts_sequence;"` //TS序号
	TsPath     string  `json:"ts_path" form:"ts_path" gorm:"column:ts_path;"`             //TS文件存储路径
	Duration   float64 `json:"duration" form:"duration" gorm:"column:duration;"`          //TS片段时长(秒)
	Definition string  `json:"definition" form:"definition" gorm:"column:definition;"`    //清晰度
	CreateTime uint    `json:"create_time" form:"create_time" gorm:"column:create_time;"` //创建时间
}

// TableName 视频ts文件表 Chips自定义表名 chips
func (Chips) TableName() string {
	return "cutter_chips"
}

// 密钥表 结构体  Secrets
type Secrets struct {
	ID        int       `json:"id" form:"id" gorm:"column:id;"`                         //主键ID
	VideoID   string    `json:"video_id" form:"video_id" gorm:"column:video_id;"`       //视频ID
	Iv        string    `json:"iv" form:"iv" gorm:"column:iv;"`                         //初始化向量
	Key       string    `json:"key" form:"key" gorm:"column:key;"`                      //密钥
	CreatedAt time.Time `json:"created_at" form:"created_at" gorm:"column:created_at;"` //创建时间
}

// TableName 密钥表 Secrets自定义表名 secrets
func (Secrets) TableName() string {
	return "cutter_secrets"
}
