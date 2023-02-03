package controller

import (
	"context"
	"github.com/RaymondCode/simple-demo/dal/mysql"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
	"net/http"
)

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
	user := utils.GetUserFromContext(c)
	//获取被点赞的视频id
	//tmp := c.FormValue("video_id")
	//videoId, _ := strconv.ParseInt(string(tmp[:]), 10, 64)
	videoId := favoriteAction.VideoId
	actionType := favoriteAction.ActionType
	//获取行为
	//tmp = c.FormValue("action_type")
	//actionType, _ := strconv.ParseInt(string(tmp[:]), 10, 64)

	hasFavorite := mysql.HasFavorite(mysql.DB, user, videoId)

	success := func() {
		log.Default().Println("操作成功！")
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 0,
		})
	}

	//如果没有点过赞且行为为点赞，则可以点赞
	if actionType == 1 && !hasFavorite {
		//进行点赞操作
		if err := mysql.FavoriteVideo(mysql.DB, user, videoId); err != nil {
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
		if err := mysql.UnFavoriteVideo(mysql.DB, user, videoId); err != nil {
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

// FavoriteList all users have same favorite video list
func FavoriteList(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
