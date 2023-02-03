package dal

import (
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
)

type videoDBAPI struct{}

var DBVideo = videoDBAPI{}

const onceFeedCnt = 30

func (*videoDBAPI) videoFavCnt(vid int, tx *gorm.DB) (error, int) {
	var cnt int
	tx.Model(&FavouriteRelation{}).
		Select("COUNT(user_id)").
		Where("video_id= ?", vid).First(&cnt)
	return tx.Error, cnt
}
func (*videoDBAPI) videoCommentsCnt(vid int, tx *gorm.DB) (error, int) {
	var cnt int
	tx.Model(&Comment{}).
		Select("COUNT(id)").
		Where("video_id= ?", vid).First(&cnt)
	return tx.Error, cnt
}
func (*videoDBAPI) isUserFav(vid int, uid int, tx *gorm.DB) (error, bool) {
	var cnt int
	tx.Model(&FavouriteRelation{}).Select("CNT(user_id)").
		Where("user_id=? AND video_id=?", uid, vid).First(&cnt)
	return tx.Error, cnt == 1
}

func (v *videoDBAPI) SelectVideo(userid int, afterTime int64) (error, []model.Video) {
	res := make([]model.Video, 0)
	err := DB.Transaction(func(tx *gorm.DB) error {
		//1. 获取视频
		rows, err := tx.Model(&VideoMeta{}).Limit(onceFeedCnt).Rows()
		if err != nil {
			return err
		}
		if rows.Next() {
			videoMeta := VideoMeta{}
			err = rows.Scan(&videoMeta)
			if err != nil {
				return err
			}
			err, favcnt := v.videoFavCnt(videoMeta.ID, tx)
			if err != nil {
				return err
			}
			err, cmtcnt := v.videoCommentsCnt(videoMeta.ID, tx)
			if err != nil {
				return err
			}
			isfav := false
			// user 为登录状态，需要检查是否喜爱
			if userid > -1 {
				err, isfav_ := v.isUserFav(videoMeta.ID, userid, tx)
				if err != nil {
					return err
				}
				isfav = isfav_
			}
			author, err := UserBasicInfo(tx, videoMeta.Author)
			if err != nil {
				return err
			}
			res = append(res, model.Video{
				Id:            int64(videoMeta.ID),
				Author:        author,
				PlayUrl:       videoMeta.PlayUrl,
				CoverUrl:      videoMeta.CoverUrl,
				FavoriteCount: int64(favcnt),
				CommentCount:  int64(cmtcnt),
				IsFavorite:    isfav,
			})
		}
		return nil
	})
	return err, res
}
