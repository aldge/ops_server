// 自动生成模板Video
package cinema

// 视频表 结构体  Video
type Video struct {
	App              *string  `json:"app" form:"app" gorm:"column:app;"`                                                           //应用
	VodId            *int32   `json:"vodId" form:"vodId" gorm:"primarykey;comment:编号;column:vod_id;"`                              //编号
	EntId            *string  `json:"entId" form:"entId" gorm:"column:ent_id;"`                                                    //发行号
	TypeId           *int64   `json:"typeId" form:"typeId" gorm:"comment:类型编号;column:type_id;"`                                    //类型编号
	TypeId1          *int64   `json:"typeId1" form:"typeId1" gorm:"comment:子类型编号;column:type_id_1;"`                               //子类型编号
	GroupId          *int64   `json:"groupId" form:"groupId" gorm:"comment:分组编号;column:group_id;"`                                 //分组编号
	VodName          *string  `json:"vodName" form:"vodName" gorm:"comment:剧名;column:vod_name;size:255;"`                          //剧名
	VodSub           *string  `json:"vodSub" form:"vodSub" gorm:"comment:副标题;column:vod_sub;size:255;"`                            //副标题
	VodEn            *string  `json:"vodEn" form:"vodEn" gorm:"comment:英文名;column:vod_en;size:255;"`                               //英文名
	VodStatus        *bool    `json:"vodStatus" form:"vodStatus" gorm:"comment:状态;column:vod_status;"`                             //状态
	VodLetter        *string  `json:"vodLetter" form:"vodLetter" gorm:"comment:首字母;column:vod_letter;"`                            //首字母
	VodColor         *string  `json:"vodColor" form:"vodColor" gorm:"comment:颜色;column:vod_color;size:6;"`                         //颜色
	VodTag           *string  `json:"vodTag" form:"vodTag" gorm:"comment:标签;column:vod_tag;size:100;"`                             //标签
	VodClass         *string  `json:"vodClass" form:"vodClass" gorm:"comment:分类;column:vod_class;size:255;"`                       //分类
	VodPic           *string  `json:"vodPic" form:"vodPic" gorm:"comment:封面;column:vod_pic;size:1024;"`                            //封面
	VodPicThumb      *string  `json:"vodPicThumb" form:"vodPicThumb" gorm:"comment:缩略图;column:vod_pic_thumb;size:1024;"`           //缩略图
	VodPicSlide      *string  `json:"vodPicSlide" form:"vodPicSlide" gorm:"comment:海报;column:vod_pic_slide;size:1024;"`            //海报
	VodPicScreenshot *string  `json:"vodPicScreenshot" form:"vodPicScreenshot" gorm:"comment:截图;column:vod_pic_screenshot;"`       //截图
	VodActor         *string  `json:"vodActor" form:"vodActor" gorm:"comment:演员;column:vod_actor;size:255;"`                       //演员
	VodDirector      *string  `json:"vodDirector" form:"vodDirector" gorm:"comment:导演;column:vod_director;size:255;"`              //导演
	VodWriter        *string  `json:"vodWriter" form:"vodWriter" gorm:"comment:编剧;column:vod_writer;size:100;"`                    //编剧
	VodBehind        *string  `json:"vodBehind" form:"vodBehind" gorm:"comment:幕后人员;column:vod_behind;size:100;"`                  //幕后人员
	VodBlurb         *string  `json:"vodBlurb" form:"vodBlurb" gorm:"comment:剧情简介;column:vod_blurb;size:255;type:text;"`           //剧情简介
	VodRemarks       *string  `json:"vodRemarks" form:"vodRemarks" gorm:"comment:备注;column:vod_remarks;size:100;"`                 //备注
	VodPubdate       *string  `json:"vodPubdate" form:"vodPubdate" gorm:"comment:发行日期;column:vod_pubdate;size:100;"`               //发行日期
	VodTotal         *int64   `json:"vodTotal" form:"vodTotal" gorm:"default:1;comment:总集数;column:vod_total;"`                     //总集数
	VodSerial        *string  `json:"vodSerial" form:"vodSerial" gorm:"comment:连载状态;column:vod_serial;size:20;"`                   //连载状态
	VodTv            *string  `json:"vodTv" form:"vodTv" gorm:"comment:电视台;column:vod_tv;size:30;"`                                //电视台
	VodWeekday       *string  `json:"vodWeekday" form:"vodWeekday" gorm:"comment:更新日;column:vod_weekday;size:30;"`                 //更新日
	VodArea          *string  `json:"vodArea" form:"vodArea" gorm:"comment:地区;column:vod_area;size:20;"`                           //地区
	VodLang          *string  `json:"vodLang" form:"vodLang" gorm:"comment:语言;column:vod_lang;size:10;"`                           //语言
	VodYear          *string  `json:"vodYear" form:"vodYear" gorm:"comment:年份;column:vod_year;size:10;"`                           //年份
	VodVersion       *string  `json:"vodVersion" form:"vodVersion" gorm:"comment:版本;column:vod_version;size:30;"`                  //版本
	VodState         *string  `json:"vodState" form:"vodState" gorm:"comment:状态;column:vod_state;size:30;"`                        //状态
	VodAuthor        *string  `json:"vodAuthor" form:"vodAuthor" gorm:"comment:作者;column:vod_author;size:60;"`                     //作者
	VodJumpurl       *string  `json:"vodJumpurl" form:"vodJumpurl" gorm:"comment:跳转;column:vod_jumpurl;size:150;"`                 //跳转
	VodTpl           *string  `json:"vodTpl" form:"vodTpl" gorm:"comment:模板;column:vod_tpl;size:30;"`                              //模板
	VodTplPlay       *string  `json:"vodTplPlay" form:"vodTplPlay" gorm:"comment:播放模板;column:vod_tpl_play;size:30;"`               //播放模板
	VodTplDown       *string  `json:"vodTplDown" form:"vodTplDown" gorm:"comment:下载模板;column:vod_tpl_down;size:30;"`               //下载模板
	VodIsend         *bool    `json:"vodIsend" form:"vodIsend" gorm:"comment:是否完结;column:vod_isend;"`                              //是否完结
	VodLock          *bool    `json:"vodLock" form:"vodLock" gorm:"comment:是否锁定;column:vod_lock;"`                                 //是否锁定
	VodLevel         *bool    `json:"vodLevel" form:"vodLevel" gorm:"comment:等级;column:vod_level;"`                                //等级
	VodCopyright     *bool    `json:"vodCopyright" form:"vodCopyright" gorm:"comment:版权;column:vod_copyright;"`                    //版权
	VodPoints        *int64   `json:"vodPoints" form:"vodPoints" gorm:"default:0;comment:积分;column:vod_points;"`                   //积分
	VodPointsPlay    *int64   `json:"vodPointsPlay" form:"vodPointsPlay" gorm:"default:0;comment:播放积分;column:vod_points_play;"`    //播放积分
	VodPointsDown    *int64   `json:"vodPointsDown" form:"vodPointsDown" gorm:"default:0;comment:下载积分;column:vod_points_down;"`    //下载积分
	VodHits          *int64   `json:"vodHits" form:"vodHits" gorm:"default:0;comment:总点击量;column:vod_hits;"`                       //总点击量
	VodHitsDay       *int64   `json:"vodHitsDay" form:"vodHitsDay" gorm:"default:0;comment:日点击量;column:vod_hits_day;"`             //日点击量
	VodHitsWeek      *int64   `json:"vodHitsWeek" form:"vodHitsWeek" gorm:"default:0;comment:周点击量;column:vod_hits_week;"`          //周点击量
	VodHitsMonth     *int64   `json:"vodHitsMonth" form:"vodHitsMonth" gorm:"default:0;comment:月点击量;column:vod_hits_month;"`       //月点击量
	VodDuration      *string  `json:"vodDuration" form:"vodDuration" gorm:"comment:单集时长;column:vod_duration;size:10;"`             //单集时长
	VodUp            *int64   `json:"vodUp" form:"vodUp" gorm:"default:0;comment:顶数;column:vod_up;"`                               //顶数
	VodDown          *int64   `json:"vodDown" form:"vodDown" gorm:"default:0;comment:踩数;column:vod_down;"`                         //踩数
	VodScore         *float64 `json:"vodScore" form:"vodScore" gorm:"comment:当前评分;column:vod_score;size:3;"`                       //当前评分
	VodScoreAll      *int64   `json:"vodScoreAll" form:"vodScoreAll" gorm:"default:0;comment:总评分;column:vod_score_all;"`           //总评分
	VodScoreNum      *int64   `json:"vodScoreNum" form:"vodScoreNum" gorm:"default:0;comment:评分人数;column:vod_score_num;"`          //评分人数
	VodTime          *int32   `json:"vodTime" form:"vodTime" gorm:"default:0;comment:更新时间;column:vod_time;"`                       //更新时间
	VodTimeAdd       *int32   `json:"vodTimeAdd" form:"vodTimeAdd" gorm:"default:0;comment:添加时间;column:vod_time_add;"`             //添加时间
	VodTimeHits      *int32   `json:"vodTimeHits" form:"vodTimeHits" gorm:"default:0;comment:点击时间;column:vod_time_hits;"`          //点击时间
	VodTimeMake      *int32   `json:"vodTimeMake" form:"vodTimeMake" gorm:"default:0;comment:制作时间;column:vod_time_make;"`          //制作时间
	VodTrysee        *int64   `json:"vodTrysee" form:"vodTrysee" gorm:"default:0;comment:试看集数;column:vod_trysee;"`                 //试看集数
	VodDoubanId      *int32   `json:"vodDoubanId" form:"vodDoubanId" gorm:"default:0;comment:豆瓣ID;column:vod_douban_id;"`          //豆瓣ID
	VodDoubanScore   *float64 `json:"vodDoubanScore" form:"vodDoubanScore" gorm:"comment:豆瓣评分;column:vod_douban_score;size:3;"`    //豆瓣评分
	VodReurl         *string  `json:"vodReurl" form:"vodReurl" gorm:"comment:重定向URL;column:vod_reurl;size:255;"`                   //重定向URL
	VodRelVod        *string  `json:"vodRelVod" form:"vodRelVod" gorm:"comment:相关影视;column:vod_rel_vod;size:255;"`                 //相关影视
	VodRelArt        *string  `json:"vodRelArt" form:"vodRelArt" gorm:"comment:相关文章;column:vod_rel_art;size:255;"`                 //相关文章
	VodPwd           *string  `json:"vodPwd" form:"vodPwd" gorm:"comment:密码;column:vod_pwd;size:10;"`                              //密码
	VodPwdUrl        *string  `json:"vodPwdUrl" form:"vodPwdUrl" gorm:"comment:密码页URL;column:vod_pwd_url;size:255;"`               //密码页URL
	VodPwdPlay       *string  `json:"vodPwdPlay" form:"vodPwdPlay" gorm:"comment:播放密码;column:vod_pwd_play;size:10;"`               //播放密码
	VodPwdPlayUrl    *string  `json:"vodPwdPlayUrl" form:"vodPwdPlayUrl" gorm:"comment:播放页URL;column:vod_pwd_play_url;size:255;"`  //播放页URL
	VodPwdDown       *string  `json:"vodPwdDown" form:"vodPwdDown" gorm:"comment:下载密码;column:vod_pwd_down;size:10;"`               //下载密码
	VodPwdDownUrl    *string  `json:"vodPwdDownUrl" form:"vodPwdDownUrl" gorm:"comment:下载密码URL;column:vod_pwd_down_url;size:255;"` //下载密码URL
	VodContent       *string  `json:"vodContent" form:"vodContent" gorm:"comment:剧情介绍;column:vod_content;type:text;"`              //剧情介绍
	VodPlayFrom      *string  `json:"vodPlayFrom" form:"vodPlayFrom" gorm:"comment:播放来源;column:vod_play_from;size:255;"`           //播放来源
	VodPlayServer    *string  `json:"vodPlayServer" form:"vodPlayServer" gorm:"comment:播放服务器;column:vod_play_server;size:255;"`    //播放服务器
	VodPlayNote      *string  `json:"vodPlayNote" form:"vodPlayNote" gorm:"comment:播放说明;column:vod_play_note;size:255;"`           //播放说明
	VodPlayUrl       *string  `json:"vodPlayUrl" form:"vodPlayUrl" gorm:"comment:播放地址;column:vod_play_url;type:text;"`             //播放地址
	VodDownFrom      *string  `json:"vodDownFrom" form:"vodDownFrom" gorm:"comment:下载来源;column:vod_down_from;size:255;"`           //下载来源
	VodDownServer    *string  `json:"vodDownServer" form:"vodDownServer" gorm:"comment:下载服务器;column:vod_down_server;size:255;"`    //下载服务器
	VodDownNote      *string  `json:"vodDownNote" form:"vodDownNote" gorm:"comment:下载说明;column:vod_down_note;size:255;"`           //下载说明
	VodDownUrl       *string  `json:"vodDownUrl" form:"vodDownUrl" gorm:"comment:下载地址;column:vod_down_url;"`                       //下载地址
	VodPlot          *bool    `json:"vodPlot" form:"vodPlot" gorm:"comment:是否有分集剧情;column:vod_plot;"`                              //是否有分集剧情
	VodPlotName      *string  `json:"vodPlotName" form:"vodPlotName" gorm:"comment:分集剧情名称;column:vod_plot_name;"`                  //分集剧情名称
	VodPlotDetail    *string  `json:"vodPlotDetail" form:"vodPlotDetail" gorm:"comment:分集剧情详情;column:vod_plot_detail;type:text;"`  //分集剧情详情
}

// TableName 视频表 Video自定义表名 cine_vod
func (Video) TableName() string {
	return "cine_vod"
}
