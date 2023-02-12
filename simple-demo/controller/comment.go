package controller

import (
	"context"
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/dal/mysql"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
	"net/http"
	"strconv"
)

type CommentListResponse struct {
	model.Response
	CommentList []model.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	model.Response
	Comment model.Comment `json:"comment,omitempty"`
}

func CommentAction(ctx context.Context, c *app.RequestContext) {
	user, status := service.GetUserFromContext(c)
	if status == false {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	var commentAction struct {
		VideoId     int64  `form:"video_id" json:"video_id" query:"video_id"`
		ActionType  int32  `form:"action_type" json:"action_type" query:"action_type"`
		CommentText string `form:"comment_text" json:"comment_text" query:"comment_text"`
		CommentId   int64  `form:"comment_id" json:"comment_id" query:"comment_id"`
	}
	if err := c.BindAndValidate(&commentAction); err != nil {
		log.Default().Println("绑定参数错误！")
		return
	}
	if commentAction.ActionType == 1 { // 发布评论
		comment, err := dal.CreateComment(mysql.DB, commentAction.VideoId, user.Id, commentAction.CommentText)
		if err == nil {
			c.JSON(http.StatusOK, CommentActionResponse{Response: model.Response{StatusCode: 0, StatusMsg: "comment success"},
				Comment: model.Comment{
					Id:         comment.Id,
					User:       user,
					Content:    comment.Content,
					CreateDate: comment.CreateTime.Format("01-02"),
				}})
			return
		}
	} else { //删除评论
		err := dal.DeleteCommentByCommentId(mysql.DB, commentAction.CommentId)
		if err == nil {
			c.JSON(http.StatusOK, CommentActionResponse{Response: model.Response{StatusCode: 0, StatusMsg: "delete success"}})
			return
		}
	}
	c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "fail"})
}

// CommentList all videos have same demo comment list
func CommentList(ctx context.Context, c *app.RequestContext) {
	VideoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	comments, err := dal.GetCommentsByVideoId(mysql.DB, VideoId)
	if err == nil {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    model.Response{StatusCode: 0, StatusMsg: "success"},
			CommentList: comments,
		})
	} else {
		c.JSON(http.StatusOK, CommentListResponse{
			Response:    model.Response{StatusCode: 1, StatusMsg: "fail"},
			CommentList: comments,
		})
	}
}
