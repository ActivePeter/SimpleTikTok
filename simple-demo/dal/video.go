package dal

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
	"log"
)

type videoDBAPI struct{}

var DBVideo = videoDBAPI{}

const onceFeedCnt = 30

//func (*videoDBAPI) videoFavCnt(vid int, tx *gorm.DB) (error, int) {
//	var cnt int
//	tx.Model(&FavouriteRelation{}).
//		Select("COUNT(user_id)").
//		Where("video_id= ?", vid).First(&cnt)
//	return tx.Error, cnt
//}
//func (*videoDBAPI) videoCommentsCnt(vid int, tx *gorm.DB) (error, int) {
//	var cnt int
//	tx.Model(&Comment{}).
//		Select("COUNT(id)").
//		Where("video_id= ?", vid).First(&cnt)
//	return tx.Error, cnt
//}
//func (*videoDBAPI) isUserFav(vid int, uid int, tx *gorm.DB) (error, bool) {
//	var cnt int
//	tx.Model(&FavouriteRelation{}).Select("CNT(user_id)").
//		//Where("user_id=? AND video_id=?", uid, vid).
//		First(&cnt)
//	return tx.Error, cnt == 1
//}

func (*videoDBAPI) SelectVideo(userid int, afterTime int64) (error, []model.Video) {
	res := make([]model.Video, 0)
	err := DB.Transaction(func(tx *gorm.DB) error {
		//1. 获取视频
		rows, err := tx.Debug().Model(&VideoMeta{}).Limit(onceFeedCnt).
			Joins("left join users on video_meta.author=users.id").
			Select("video_meta.id, video_meta.author, video_meta.play_url, video_meta.cover_url,"+ //video
				"users.id, users.name, users.follow_count, users.follower_count,"+ //author
				"exists (?),"+ //是否喜欢
				"(?),"+ //喜爱数
				"(?)", //评论数
				tx.Model(&FavouriteRelation{}).
					Where("user_id=? AND video_id=video_meta.id", userid),
				tx.Model(&FavouriteRelation{}).
					Where("video_id=video_meta.id").Select("COUNT(user_id)"),
				tx.Model(&Comment{}).
					Where("video_id=video_meta.id").Select("COUNT(id)"),
			).Limit(onceFeedCnt).Rows()
		if err != nil {
			return err
		}
		for rows.Next() {
			v := model.Video{}
			var author_id int
			err := rows.Scan(&v.Id, &author_id, &v.PlayUrl, &v.CoverUrl,
				&v.Author.Id, &v.Author.Name, &v.Author.FollowCount, &v.Author.FollowerCount,
				&v.IsFavorite,
				&v.FavoriteCount,
				&v.CommentCount,
			)
			if err != nil {
				log.Default().Printf("read video meta failed %v\n", err)
				return err
			}
			res = append(res, v)
		}
		return nil
	})
	fmt.Printf("%+v\n", res)
	return err, res
}
