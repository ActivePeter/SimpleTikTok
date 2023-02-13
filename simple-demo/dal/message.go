package dal

import (
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
)

func AddMessage(tx *gorm.DB, msg model.Msg) bool {
	if err := tx.Table("messages").Select("ToUserId", "FromUserId", "Content", "CreateTime").Create(msg).Error; err != nil {
		return false
	} else {
		return true
	}
}
