package controller

import (
	"context"
	"github.com/RaymondCode/simple-demo/dal/mysql"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/mw"
	"github.com/cloudwego/hertz/pkg/app"
	_ "gorm.io/gorm"
	"net/http"
	"strings"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]model.User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	model.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	model.Response
	User model.User `json:"user"`
}

func Register(ctx context.Context, c *app.RequestContext) {
	username := c.Query("username")
	password := c.Query("password")

	// 校验用户名和密码的合法性，均不能大于32位
	if strings.Count(username, "") > 32 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "用户名不能大于32位！"},
		})
		return
	}
	if strings.Count(password, "") > 32 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "密码不能大于32位！"},
		})
		return
	}

	// 校验用户名的唯一性
	res, _ := mysql.FindUserByUsername(username)
	if res > 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: model.Response{StatusCode: 1, StatusMsg: "用户名已存在！"},
		})
		return
	}

	// 校验通过，注册信息放入数据库，name默认和username相同
	id, _ := mysql.CreateUser(username, password)
	//token := username + password
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: model.Response{StatusCode: 0},
		UserId:   id,
		//Token:    token,
	})
}

func UserInfo(ctx context.Context, c *app.RequestContext) {
	user, _ := c.Get(mw.IdentityKey)
	c.JSON(http.StatusOK, UserResponse{
		Response: model.Response{StatusCode: 0},
		User:     *user.(*model.User),
	})
}
