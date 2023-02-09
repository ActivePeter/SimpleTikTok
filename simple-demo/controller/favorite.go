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

	service.FavoriteAction(c, user, videoId, actionType)

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
