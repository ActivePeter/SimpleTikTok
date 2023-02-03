package mysql

import (
	//"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
)

// 判断是否已经点赞
func HasFavorite(tx *gorm.DB, user model.User, videoId int64) bool {
	tmp := model.FavouriteRelation{}
	//查询
	if res := tx.Where("user_id = ? And video_id = ?", user.Id, videoId).First(&tmp).RowsAffected; res == 0 {
		//未查到则没有点赞
		return false
	} else {
		//查到，则已经点赞
		return true
	}
}

func FavoriteVideo(tx *gorm.DB, user model.User, videoId int64) error {
	relation := model.FavouriteRelation{
		VideoID: int(videoId),
		UserID:  int(user.Id),
	}
	if err := tx.Create(&relation).Error; err != nil {
		return err
	} else {
		return nil
	}
}

func UnFavoriteVideo(tx *gorm.DB, user model.User, videoId int64) error {
	relation := model.FavouriteRelation{
		VideoID: int(videoId),
		UserID:  int(user.Id),
	}
	if err := tx.Where("user_id = ? And video_id = ?", user.Id, videoId).Delete(&relation).Error; err != nil {
		return err
	} else {
		return nil
	}
}
