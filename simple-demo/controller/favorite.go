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
	"strconv"
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
	//判断点赞合法性
	//获取当前用户
	user, _ := service.GetUserFromContext(c)
	//获取被点赞的视频id
	//tmp := c.FormValue("video_id")
	//videoId, _ := strconv.ParseInt(string(tmp[:]), 10, 64)
	videoId := favoriteAction.VideoId
	actionType := favoriteAction.ActionType
	//获取行为
	//tmp = c.FormValue("action_type")
	//actionType, _ := strconv.ParseInt(string(tmp[:]), 10, 64)

	hasFavorite := dal.HasFavorite(mysql.DB, user, videoId)

	success := func() {
		log.Default().Println("操作成功！")
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
			success()
		}
	} else if actionType == 1 && hasFavorite {
		c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: 400,
			StatusMsg:  "已点赞！",
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
			success()
		}
	} else {
		c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: 400,
			StatusMsg:  "已取消点赞！",
		})
	}

}

// 返回该用户的点赞列表
func FavoriteList(ctx context.Context, c *app.RequestContext) {
	//c.JSON(http.StatusOK, VideoListResponse{
	//	Response: model.Response{
	//		StatusCode: 0,
	//	},
	//	VideoList: DemoVideos,
	//})
	user, _ := service.GetUserFromContext(c)
	userId := c.Query("user_id")
	u, _ := strconv.ParseInt(userId, 10, 64)
	if user.Id != u {
		log.Default().Println("invalid request param")
		c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: 400,
			StatusMsg:  "invalid request param",
		})
	} else {
		//获取该用户点赞列表
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
}
