package main

import (
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/mw"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	go service.RunMessageServer()
	dal.Init()
	mw.InitJwt()
	h := server.Default()
	initRouter(h)
	h.Spin()
}
