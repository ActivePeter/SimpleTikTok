package mw

import (
	"context"
	"errors"
	"github.com/RaymondCode/simple-demo/dal/mysql"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/hertz-contrib/jwt"
	"net/http"
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
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, utils.H{
				"code":    code,
				"token":   token,
				"expire":  expire.Format(time.RFC3339),
				"message": "success",
			})
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginStruct struct {
				Username string `form:"username" json:"username" query:"username" vd:"(len($) > 0 && len($) < 30); msg:'Illegal format'"`
				Password string `form:"password" json:"password" query:"password" vd:"(len($) > 0 && len($) < 30); msg:'Illegal format'"`
			}
			if err := c.BindAndValidate(&loginStruct); err != nil {
				return nil, err
			}
			users, err := mysql.CheckUser(loginStruct.Username, loginStruct.Password)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			if len(users) == 0 {
				return nil, errors.New("user already exists or wrong password")
			}
			return users[0], nil
		},
		IdentityKey: IdentityKey,
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return &model.User{
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
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			hlog.CtxErrorf(ctx, "jwt biz err = %+v", e.Error())
			return e.Error()
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			type UserResponse struct {
				model.Response
				User model.User `json:"user"`
			}
			c.JSON(http.StatusOK, UserResponse{
				Response: model.Response{StatusCode: 1, StatusMsg: "用户名或密码错误"},
			})
		},
	})
	if err != nil {
		panic(err)
	}
}
