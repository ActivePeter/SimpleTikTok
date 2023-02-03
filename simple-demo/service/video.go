package service

import (
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/model"
)

type video struct{}

var ServiceVideo = video{}

func (*video) GetFeedList(userid int, after int64) (error, []model.Video) {
	return dal.DBVideo.SelectVideo(userid, after)
}
