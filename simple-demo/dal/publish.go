package dal

import (
	"bytes"
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
)

//上传视频至数据库
func UploadVideo(tx *gorm.DB, Video VideoMeta) error {

	if err := tx.Create(&Video); err != nil {
		return err.Error
	}
	return nil
}

// 获取最新的视频id
func GetLatestVideoId(tx *gorm.DB) (model.VideoId, error) {
	//res := model.Video{}
	res_meta := VideoMeta{}

	//按id倒叙找到第一个id，即最大的id，再将这个vedio_meta型的id赋值给video型
	err := DB.Transaction(func(tx *gorm.DB) error {
		tx.Debug().Model(&VideoMeta{}).
			Select("id").
			Order("id desc").First(&res_meta)
		return nil
	})
	res := model.Video{
		Id: res_meta.ID,
	}

	fmt.Println(res.Id)
	return res.Id, err

}

// 截视频的第一帧作为封面
func GetSnapshot(videoPath, snapshotPath string, frameNum int) (snapshotName string, err error) {
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	err = imaging.Save(img, snapshotPath+".png")
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	names := strings.Split(snapshotPath, "\\")
	snapshotName = names[len(names)-1] + ".png"
	return
}

// 获取个人视频列表
func GetViedoList(userid model.UserId) ([]model.Video, error) {
	res := make([]model.Video, 0)
	err := DB.Transaction(func(tx *gorm.DB) error {
		//1. 获取视频
		rows, err := tx.Debug().Model(&VideoMeta{}).
			Joins("left join users on video_meta.author=users.id").
			Where("video_meta.author=?", userid).
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
			).Rows()
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
	return res, err
}
