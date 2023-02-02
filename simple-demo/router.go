package main

import (
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/RaymondCode/simple-demo/mw"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func initRouter(r *server.Hertz) {
	r.Static("/static", "./public")
	apiRouter := r.Group("/douyin")
	// basic apis
	apiRouter.GET("/feed/", controller.Feed)

	apiRouter.GET("/user/", mw.JwtMiddleware.MiddlewareFunc(), controller.UserInfo)
	apiRouter.POST("/user/register/", mw.JwtMiddleware.LoginHandler)
	apiRouter.POST("/user/login/", mw.JwtMiddleware.LoginHandler)
	apiRouter.POST("/publish/action/", mw.JwtMiddleware.MiddlewareFunc(), controller.Publish)
	apiRouter.GET("/publish/list/", mw.JwtMiddleware.MiddlewareFunc(), controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", mw.JwtMiddleware.MiddlewareFunc(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", mw.JwtMiddleware.MiddlewareFunc(), controller.FavoriteList)
	apiRouter.POST("/comment/action/", mw.JwtMiddleware.MiddlewareFunc(), controller.CommentAction)
	apiRouter.GET("/comment/list/", mw.JwtMiddleware.MiddlewareFunc(), controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", mw.JwtMiddleware.MiddlewareFunc(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", mw.JwtMiddleware.MiddlewareFunc(), controller.FollowList)
	apiRouter.GET("/relation/follower/list/", mw.JwtMiddleware.MiddlewareFunc(), controller.FollowerList)
	apiRouter.GET("/relation/friend/list/", mw.JwtMiddleware.MiddlewareFunc(), controller.FriendList)
	apiRouter.GET("/message/chat/", mw.JwtMiddleware.MiddlewareFunc(), controller.MessageChat)
	apiRouter.POST("/message/action/", mw.JwtMiddleware.MiddlewareFunc(), controller.MessageAction)

}
