package service

import (
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/mw"
	"github.com/cloudwego/hertz/pkg/app"
)

func GetUserFromContext(c *app.RequestContext) (model.User, bool) {
	user, ok := c.Get(mw.IdentityKey)
	if ok {
		return user.(model.User), ok
	}
	return model.User{}, ok
}
