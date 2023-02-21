package dal

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
	"log"
)

type dAOVideo struct{}

var DAOVideo = dAOVideo{}

const onceFeedCnt = 30

func (*dAOVideo) SelectVideo(userid model.UserId, afterTime int64) (error, []model.Video) {

	res := make([]model.Video, 0)
	err := DB.Transaction(func(tx *gorm.DB) error {
		//1. 获取视频
		rows, err := tx.Debug().Model(&VideoMeta{}).Limit(onceFeedCnt).
			Joins("inner join users on video_meta.author=users.id").
			Select("video_meta.id, video_meta.author, video_meta.play_url, video_meta.cover_url,"+ //video
				"users.id, users.name, "+
				"(?), (?),"+ //follow & follower count
				"exists (?),"+ //是否喜欢
				"(?),"+ //喜爱数
				"(?)", //评论数
				tx.Model(&FollowRelation{}).Select("COUNT(from_id)").Where("from_id=?", userid),
				tx.Model(&FollowRelation{}).Select("COUNT(from_id)").Where("to_id=?", userid),

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
