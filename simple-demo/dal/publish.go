package dal

import (
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
)

func GetViedoList(tx *gorm.DB, user model.User) ([]model.Video, error) {
	user_id := user.Id
	//获取该用户发表的视频id
	res := make([]model.Video, 0)
	tx.Where("author = ?", user_id).Find(&res)
	return res, nil
}
