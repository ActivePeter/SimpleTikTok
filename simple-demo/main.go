package main

import (
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/mw"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	err, config := utils.LoadConfig()
	if err != nil {
		return
	}
	go service.RunMessageServer()
	dal.Init(config)
	service.ServerDomain = config.ServerDomain
	mw.InitJwt()
	h := server.New(server.WithStreamBody(true))

	service.Video.InitVideoFs(h)

	initRouter(h)
	h.Spin()
}