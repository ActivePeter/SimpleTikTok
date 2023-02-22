package service

import (
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/model"
	"log"
)

type video struct{}

var Video = video{}

func (*video) GetFeedList(userid model.UserId, after int64) (error, []model.Video) {
	err, videos := dal.DAOVideo.SelectVideo(userid, after)
	for i, _ := range videos {
		videos[i].PlayUrl = "http://" + ServerDomain + videos[i].PlayUrl
		videos[i].CoverUrl = "http://" + ServerDomain + videos[i].CoverUrl
	}
	log.Printf("%v", videos)
	return err, videos
}
