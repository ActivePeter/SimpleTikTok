package dal

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

type User2 struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"username,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
}

type Author struct {
	User2
	IsFollow bool
}

type DetailVideo struct {
	Id            int64  `json:"id,omitempty"`
	Author        Author `json:"author"`
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
}

// 获取点赞视频列表
func GetFavoriteVideos(tx *gorm.DB, user model.User) ([]DetailVideo, error) {
	relations := make([]model.FavouriteRelation, 0)
	//获取该用户所有点赞过的视频id
	tx.Where("user_id = ?", user.Id).Find(&relations)

	videos := make([]DetailVideo, 0)
	//循环获取每一个视频id对应的信息
	for _, relation := range relations {
		videoMeta := new(VideoMeta)

		//获取视频id对应的视频元数据信息
		tx.Where("id = ?", relation.VideoID).First(videoMeta)
		authorId := videoMeta.Author

		//通过视频作者id获取作者信息
		tmpUser := new(User2)
		tmp := new(Author)

		tx.Table("users").Where("id = ?", authorId).First(tmpUser)

		tmp.User2 = *tmpUser

		//暂未实现
		tmp.IsFollow = true

		video := DetailVideo{
			Id:            videoMeta.ID,
			Author:        *tmp,
			PlayUrl:       videoMeta.PlayUrl,
			CoverUrl:      videoMeta.CoverUrl,
			FavoriteCount: GetFavoriteCount(tx, videoMeta.ID),
			CommentCount:  0,
			IsFavorite:    true,
			Title:         videoMeta.Title,
		}
		videos = append(videos, video)
	}
	return videos, nil
}

// 获取视频点赞次数
func GetFavoriteCount(tx *gorm.DB, videoId int64) int64 {
	var count int64
	tx.Table("favourite_relations").Where("video_id = ?", videoId).Count(&count)
	return count
}
