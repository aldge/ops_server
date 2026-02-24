package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/cinema"
)

func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(cinema.Video{}, cinema.Amateur{})
	if err != nil {
		return err
	}
	return nil
}
