package controller

import (
	"context"
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/dal/mysql"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
	"net/http"
)

type VideosResponse struct {
	StatusCode int32             `json:"status_code"`
	VideoList  []dal.DetailVideo `json:"video_list,omitempty"`
}

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(ctx context.Context, c *app.RequestContext) {

	var favoriteAction struct {
		VideoId    int64 `form:"video_id" json:"video_id" query:"video_id"`
		ActionType int32 `form:"action_type" json:"action_type" query:"action_type"`
	}

	if err := c.BindAndValidate(&favoriteAction); err != nil {
		log.Default().Println("绑定参数错误！")
		return
	}

	//获取当前用户
	user, _ := service.GetUserFromContext(c)
	//获取被点赞的视频id
	videoId := favoriteAction.VideoId
	//获取行为
	actionType := favoriteAction.ActionType

	hasFavorite := dal.HasFavorite(mysql.DB, user, videoId)

	success := func() {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 0,
		})
	}

	//如果没有点过赞且行为为点赞，则可以点赞
	if actionType == 1 && !hasFavorite {
		//进行点赞操作
		if err := dal.FavoriteVideo(mysql.DB, user, videoId); err != nil {
			log.Default().Println(err.Error())
			c.JSON(http.StatusBadRequest, model.Response{
				StatusCode: 400,
				StatusMsg:  "服务器出现错误",
			})
			return
		} else {
			log.Default().Println("点赞成功")
			success()
		}
	} else if actionType == 1 && hasFavorite {
		log.Default().Println("已经给该视频点过赞！")
		c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: 400,
			StatusMsg:  "不能给已经点过赞的视频点赞！",
		})
	} else if actionType == 2 && hasFavorite {
		//取消点赞
		if err := dal.UnFavoriteVideo(mysql.DB, user, videoId); err != nil {
			log.Default().Println(err.Error())
			c.JSON(http.StatusBadRequest, model.Response{
				StatusCode: 400,
				StatusMsg:  "服务器出现错误",
			})
			return
		} else {
			log.Default().Println("成功取消点赞")
			success()
		}
	} else {
		log.Default().Println("已经取消点赞！")
		c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: 400,
			StatusMsg:  "不能给未点赞的视频取消点赞",
		})
	}

}

// 返回该用户的点赞列表
func FavoriteList(ctx context.Context, c *app.RequestContext) {

	user, _ := service.GetUserFromContext(c)

	if videos, err := dal.GetFavoriteVideos(mysql.DB, user); err != nil {
		log.Default().Println(err.Error())
		c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: 400,
			StatusMsg:  "服务器出现错误",
		})
	} else {
		c.JSON(http.StatusOK, VideosResponse{
			0,
			videos,
		})
	}
}
