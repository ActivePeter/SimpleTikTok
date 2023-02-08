package mw

import (
	"context"
	"errors"
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/model"
	myUtils "github.com/RaymondCode/simple-demo/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"
	"net/http"
	"strings"
	"time"
)

var (
	JwtMiddleware *jwt.HertzJWTMiddleware
	IdentityKey   = "identity"
)

func InitJwt() {
	var err error
	JwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:         "SimpleTikTok",
		Key:           []byte("secret key"),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		TokenLookup:   "header: Authorization, query: token, cookie: jwt, form: token",
		TokenHeadName: "Bearer",
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			if code == 200 {
				code = 0
			} else {
				code = -1
			}
			c.JSON(http.StatusOK, utils.H{
				"status_code": code,
				"status_msg":  "success",
				"token":       token,
				"expire":      expire.Format(time.RFC3339),
			})
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginStruct struct {
				Username string `form:"username" json:"username" query:"username" vd:"(len($) > 0 && len($) < 33); msg:'Illegal format'"`
				Password string `form:"password" json:"password" query:"password" vd:"(len($) > 0 && len($) < 33); msg:'Illegal format'"`
			}
			if err := c.BindAndValidate(&loginStruct); err != nil {
				return nil, errors.New("用户名和密码长度不能超过32位且不能为空")
			}

			if strings.Index(c.FullPath(), "login") != -1 { // 登陆
				users, err := dal.CheckUser(dal.DB, loginStruct.Username, myUtils.MD5(loginStruct.Password)) // 校验用户名和密码是否正确
				if err != nil {
					return nil, errors.New("用户名或密码错误！")
				}
				return users[0], nil
			} else { // 注册
				res, _ := dal.FindUserByUsername(dal.DB, loginStruct.Username) // 校验用户名是否已经存在
				if res > 0 {
					return nil, errors.New("用户名已存在！")
				}
				user, _ := dal.CreateUser(dal.DB, loginStruct.Username, myUtils.MD5(loginStruct.Password))
				return user, nil
			}
		},
		IdentityKey: IdentityKey,
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return model.User{
				Id:            int64(claims["Id"].(float64)),
				Name:          claims["Name"].(string),
				FollowCount:   int64(claims["FollowCount"].(float64)),
				FollowerCount: int64(claims["FollowerCount"].(float64)),
			}
		},
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					"Id":            v.Id,
					"Name":          v.Name,
					"FollowCount":   v.FollowCount,
					"FollowerCount": v.FollowerCount,
				}
			}
			return jwt.MapClaims{}
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			type UserResponse struct {
				model.Response
				User model.User `json:"user"`
			}
			c.JSON(http.StatusOK, UserResponse{
				Response: model.Response{StatusCode: 1, StatusMsg: message},
			})
		},
	})
	if err != nil {
		panic(err)
	}
}
