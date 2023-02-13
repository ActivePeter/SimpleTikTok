package controller

import (
	"bytes"
	"context"
	"fmt"
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/dal/mysql"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type VideoListResponse struct {
	model.Response
	VideoList []model.Video `json:"video_list"`
}

// 检测文件是否是视频文件
func isVideo(videoPath string) bool {
	file_bytes, err := os.ReadFile(videoPath)
	if err != nil {
		log.Fatal(err)
	}
	FileType := http.DetectContentType(file_bytes)
	fmt.Println(FileType)
	if strings.Contains(FileType, "video") {
		return true
	} else {
		return false
	}

}

// 截视频的第一帧作为封面
func GetSnapshot(videoPath, snapshotPath string, frameNum int) (snapshotName string, err error) {
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	err = imaging.Save(img, snapshotPath+".png")
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
		return "", err
	}

	names := strings.Split(snapshotPath, "\\")
	snapshotName = names[len(names)-1] + ".png"
	return
}

// Publish check token then save upload file to public directory
func Publish(ctx context.Context, c *app.RequestContext) {
	user, status := service.GetUserFromContext(c)
	if status == false {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	//接收文件和标题
	data, err := c.FormFile("data")
	title := c.PostForm("title")
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	filename := filepath.Base(data.Filename)

	//获取文件名后缀
	fileSuffix := path.Ext(filename)
	laststVideoId, err := dal.GetLatestVideoId(mysql.DB)
	if err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//文件名为用户id_视频id.后缀
	finalName := fmt.Sprintf("%d_%d%s", user.Id, laststVideoId+1, fileSuffix)
	photoName := fmt.Sprintf("%d_%d", user.Id, laststVideoId+1)
	fmt.Println(filename)
	saveFile := filepath.Join("./public/video/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//检测“视频“是否真的是视频文件
	if isVideo("./public/video/"+finalName) != true {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  "This file is not video!",
		})
		return
	}

	//截取视频第一帧作为封面
	GetSnapshot("./public/video/"+finalName, "./public/photo/"+photoName, 1)

	videoMeta := dal.VideoMeta{
		Author:     user.Id,
		PlayUrl:    "./public/video/" + finalName,
		CoverUrl:   "./public/photo/" + photoName + ".png",
		Title:      title,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	//上传视频信息至数据库
	dal.UploadVideo(mysql.DB, videoMeta)

	c.JSON(http.StatusOK, model.Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(ctx context.Context, c *app.RequestContext) {
	user, _ := service.GetUserFromContext(c)
	user_id := user.Id
	if vedios, err := dal.GetViedoList(user_id); err != nil {
		c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: model.Response{
				StatusCode: 0,
			},
			VideoList: vedios,
		})
	}
}
