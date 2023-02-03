package controller

import (
	"context"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
	"net/http"
)

type FeedResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list,omitempty"`
	NextTime  int64         `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(ctx context.Context, c *app.RequestContext) {
	//log.Default().Printf("feed %v \n", c.Request.QueryString())
	type FeedRequest struct {
		LatestTime int64  `json:"latest_time,omitempty"`
		Token      string `json:"token,omitempty"`
	}
	Fail := func() {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  model.Response{StatusCode: 1},
			VideoList: make([]model.Video, 0),
		})
	}
	req := FeedRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		log.Default().Printf("req invalid %v\n", err)
		Fail()
		return
	}
	service.GetUserFromContext(c)
	err, videos := service.Video.GetFeedList(-1, 0)
	if err != nil {
		log.Default().Printf("select video failed %v\n", err)
		Fail()
		return
	}
	// 搜索latest之前的video
	c.JSON(http.StatusOK, FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: videos,
		NextTime:  utils.TimeStamp(),
	})
}
