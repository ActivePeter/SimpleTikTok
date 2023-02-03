package dal

import (
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
)

func UserBasicInfo(tx *gorm.DB, uid int) (model.User, error) {
	var res model.User
	tx.Model(&User{}).Where("id=?", uid).First(&res)
	return res, DB.Error
}
