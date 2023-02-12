package controller

import (
	"context"
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
	"net/http"
)

type UserListResponse struct {
	model.Response
	UserList []model.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(ctx context.Context, c *app.RequestContext) {
	type RelationActionRequest struct {
		Token      string
		ToUserId   model.UserId `query:"to_user_id"`
		ActionType int32        `query:"action_type"`
	}
	Fail := func(msg string) {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: 1,
			StatusMsg:  msg,
		})
	}
	req := RelationActionRequest{}
	// 1。read req
	err := c.BindAndValidate(&req)
	if err != nil {
		log.Default().Printf("req invalid %v\n", err)
		Fail("invalid request arg")
		return
	}
	// 2。user token
	user, ok := service.GetUserFromContext(c)
	if !ok {
		Fail("not logged in")
	}
	follow := true
	if req.ActionType == 2 {
		follow = false
	}
	// 3. ope sql
	err = service.Relation.SetFollow(user.Id, req.ToUserId, follow)
	if err != nil {
		Fail("state disagree")
		return
	}
	c.JSON(http.StatusOK, model.Response{StatusCode: 0})
}

func FollowList(ctx context.Context, c *app.RequestContext) {
	user, status := service.GetUserFromContext(c)
	if status == false {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	follows, err := dal.SelectFollows(dal.DB, user.Id)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "fail",
			},
		})
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserList: follows,
	})
}

func FollowerList(ctx context.Context, c *app.RequestContext) {
	user, status := service.GetUserFromContext(c)
	if status == false {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	followers, err := dal.SelectFollowers(dal.DB, user.Id)
	if err != nil {
		c.JSON(http.StatusOK, UserListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "fail",
			},
		})
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		UserList: followers,
	})
}

// FriendList all users have same friend list
func FriendList(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, UserListResponse{
		Response: model.Response{
			StatusCode: 0,
		},
		UserList: []model.User{DemoUser},
	})
}
