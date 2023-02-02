package controller

import (
	"context"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/cloudwego/hertz/pkg/app"
	"net/http"
)

type CommentListResponse struct {
	model.Response
	CommentList []model.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	model.Response
	Comment model.Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	actionType := c.Query("action_type")

	if user, exist := usersLoginInfo[token]; exist {
		if actionType == "1" {
			text := c.Query("comment_text")
			c.JSON(http.StatusOK, CommentActionResponse{Response: model.Response{StatusCode: 0},
				Comment: model.Comment{
					Id:         1,
					User:       user,
					Content:    text,
					CreateDate: "05-01",
				}})
			return
		}
		c.JSON(http.StatusOK, model.Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    model.Response{StatusCode: 0},
		CommentList: DemoComments,
	})
}
