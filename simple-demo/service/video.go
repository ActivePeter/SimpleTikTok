package service

import (
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/model"
)

type video struct{}

var Video = video{}

func (*video) GetFeedList(userid model.UserId, after int64) (error, []model.Video) {
	return dal.DBVideo.SelectVideo(userid, after)
}
