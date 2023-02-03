package controller

import (
	"context"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/cloudwego/hertz/pkg/app"
	_ "gorm.io/gorm"
	"net/http"
)

// usersLoginInfo 请不要使用该字典，该字典来自demo，token的有效性已经在jwt中间件进行了验证
// 要想获取user请使用utils.GetUserFromContext(c)
var usersLoginInfo = map[string]model.User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	model.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	model.Response
	User model.User `json:"user"`
}

func UserInfo(ctx context.Context, c *app.RequestContext) {
	user, _ := service.GetUserFromContext(c)
	c.JSON(http.StatusOK, UserResponse{
		Response: model.Response{StatusCode: 0},
		User:     user,
	})
}
