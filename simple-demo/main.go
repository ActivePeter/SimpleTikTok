package main

import (
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/mw"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	err, config := utils.loadConfig()
	if err != nil {
		return
	}
	go service.RunMessageServer()
	dal.Init(config)
	mw.InitJwt()
	h := server.New(server.WithStreamBody(true))
	initRouter(h)
	h.Spin()
}
