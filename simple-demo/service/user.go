package service

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/mw"
	"github.com/cloudwego/hertz/pkg/app"
)

func GetUserFromContext(c *app.RequestContext) model.User {
	user, _ := c.Get(mw.IdentityKey)
	return user.(model.User)
}
