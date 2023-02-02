package controller

import (
	"context"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
	"time"
)

type FeedResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list,omitempty"`
	NextTime  int64         `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, FeedResponse{
		Response:  model.Response{StatusCode: 0},
		VideoList: DemoVideos,
		NextTime:  time.Now().Unix(),
	})
}