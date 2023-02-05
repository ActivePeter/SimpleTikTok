package service

import (
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/dal/mysql"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"strconv"
)

const pre string = "video:favorite:"

func FavoriteAction(c *app.RequestContext, user model.User, videoId int64, actionType int32) {
	//利用redis中的set判断是否已经对该视频点赞

	//key  video:favorite: + videoId
	videoKey := pre + strconv.FormatInt(videoId, 10)

	var hasFavorite bool
	_, err := dal.RD.Get(videoKey).Result()
	//不存在该key，也就是一定没人点过赞
	if err == redis.Nil {
		hasFavorite = false
	} else {
		//判断useId是否在video:favorite:中
		if exists, _ := dal.RD.SIsMember(videoKey, user.Id).Result(); exists {
			hasFavorite = true
		} else {
			hasFavorite = false
		}
	}

	//查询数据库判断
	//hasFavorite = dal.HasFavorite(mysql.DB, user, videoId)

	success := func() {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 0,
		})
	}

	//如果没有点过赞且行为为点赞，则可以点赞
	if actionType == 1 && !hasFavorite {
		//给前台返回结果
		success()
		//进行点赞操作
		dal.RD.SAdd(videoKey, user.Id)
		//操作数据库
		if err := dal.FavoriteVideo(mysql.DB, user, videoId); err != nil {
			log.Default().Println(err.Error())
			c.JSON(http.StatusBadRequest, model.Response{
				StatusCode: 400,
				StatusMsg:  "服务器出现错误",
			})
			return
		}
	} else if actionType == 1 && hasFavorite {
		log.Default().Println("已经给该视频点过赞！")
		c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: 400,
			StatusMsg:  "不能给已经点过赞的视频点赞！",
		})
	} else if actionType == 2 && hasFavorite {

		success()
		//取消点赞
		dal.RD.SRem(videoKey, user.Id)

		if err := dal.UnFavoriteVideo(mysql.DB, user, videoId); err != nil {
			log.Default().Println(err.Error())
			c.JSON(http.StatusBadRequest, model.Response{
				StatusCode: 400,
				StatusMsg:  "服务器出现错误",
			})
			return
		}
	} else {
		log.Default().Println("已经取消点赞！")
		c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: 400,
			StatusMsg:  "不能给未点赞的视频取消点赞",
		})
	}
}
