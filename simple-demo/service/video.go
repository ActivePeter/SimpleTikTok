package service

import (
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"log"
	"os"
	"strings"
)

type video struct{}

var Video = video{}

var VideoStoreRoot = "public"
var VideoRelateFsDir = VideoStoreRoot + "/video/"
var VideoCoverRelateFsDir = VideoStoreRoot + "/photo/"
var VideoReachRoot = "public"

func WrapVideoAndCover(videos *[]model.Video) {
	for i, _ := range *videos {
		(*videos)[i].PlayUrl = "http://" + ServerDomain + "/" + (*videos)[i].PlayUrl
		(*videos)[i].CoverUrl = "http://" + ServerDomain + "/" + (*videos)[i].CoverUrl
	}
}
func WrapVideoAndCover2(videos *[]dal.DetailVideo) {
	for i, _ := range *videos {
		(*videos)[i].PlayUrl = "http://" + ServerDomain + "/" + (*videos)[i].PlayUrl
		(*videos)[i].CoverUrl = "http://" + ServerDomain + "/" + (*videos)[i].CoverUrl
	}
}

func (*video) GetFeedList(userid model.UserId, after int64) (error, []model.Video) {
	err, videos := dal.DAOVideo.SelectVideo(userid, after)
	WrapVideoAndCover(&videos)
	//for i, _ := range videos {
	//	videos[i].PlayUrl = "http://" + ServerDomain + videos[i].PlayUrl
	//	videos[i].CoverUrl = "http://" + ServerDomain + videos[i].CoverUrl
	//}
	log.Printf("%v", videos)
	return err, videos
}

func (*video) InitVideoFs(hertz *server.Hertz) {
	os.MkdirAll("./public/video", os.ModePerm)
	os.MkdirAll("./public/photo", os.ModePerm)
	fs := &app.FS{Root: "./" + VideoStoreRoot, PathRewrite: getPathRewriter("/" + VideoReachRoot)}
	hertz.StaticFS("/"+VideoReachRoot, fs)
}

func getPathRewriter(prefix string) app.PathRewriteFunc {
	// Cannot have an empty prefix
	if prefix == "" {
		prefix = "/"
	}
	// Prefix always start with a '/' or '*'
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}

	// Is prefix a direct wildcard?
	isStar := prefix == "/*"
	// Is prefix a partial wildcard?
	if strings.Contains(prefix, "*") {
		isStar = true
		prefix = strings.Split(prefix, "*")[0]
		// Fix this later
	}
	prefixLen := len(prefix)
	if prefixLen > 1 && prefix[prefixLen-1:] == "/" {
		// /john/ -> /john
		prefixLen--
		prefix = prefix[:prefixLen]
	}
	return func(ctx *app.RequestContext) []byte {
		path := ctx.Path()
		if len(path) >= prefixLen {
			if isStar && string(path[0:prefixLen]) == prefix {
				path = append(path[0:0], '/')
			} else {
				path = path[prefixLen:]
				if len(path) == 0 || path[len(path)-1] != '/' {
					path = append(path, '/')
				}
			}
		}
		if len(path) > 0 && path[0] != '/' {
			path = append([]byte("/"), path...)
		}
		return path
	}
}
