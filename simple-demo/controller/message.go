package controller

import (
	"context"
	"fmt"
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
	"net/http"
	"strconv"
	"time"
)

var tempChat = map[string][]model.Message{}

var messageIdSequence = int64(1)

type ChatResponse struct {
	model.Response
	MessageList []model.Message `json:"message_list"`
}

// MessageAction no practical effect, just check if token is valid
func MessageAction(ctx context.Context, c *app.RequestContext) {
	log.Default().Println("MessageAction")
	//验证用户
	user, exists := service.GetUserFromContext(c)
	if !exists {
		log.Default().Println("用户不存在!")
		return
	}

	var messageAction struct {
		ToUserId   int64  `form:"to_user_id" json:"to_user_id" query:"to_user_id"`
		ActionType int32  `form:"action_type" json:"action_type" query:"action_type"`
		Content    string `form:"content" json:"content" query:"content"`
	}

	if err := c.BindAndValidate(&messageAction); err != nil {
		log.Default().Println("参数绑定错误!")
		return
	}

	//发送方用户id
	fromUserId := user.Id
	toUserId := messageAction.ToUserId
	actionType := messageAction.ActionType
	content := messageAction.Content

	if actionType != 1 {
		log.Default().Println("功能暂未完善")
		c.JSON(http.StatusNotFound, model.Response{
			StatusCode: http.StatusNotFound,
			StatusMsg:  "功能尚未完善",
		})
		return
	}

	msg := model.Msg{
		ToUserId:   toUserId,
		FromUserId: fromUserId,
		Content:    content,
		CreateTime: time.Now().String(),
	}

	if success := dal.AddMessage(dal.DB, msg); !success {
		log.Default().Println("发送消息失败")
		c.JSON(http.StatusInternalServerError, model.Response{
			StatusCode: http.StatusInternalServerError,
			StatusMsg:  "发送消息失败",
		})
	} else {
		c.JSON(http.StatusOK, model.Response{
			StatusCode: http.StatusOK,
		})
	}

}

// MessageChat all users have same follow list
func MessageChat(ctx context.Context, c *app.RequestContext) {
	log.Default().Println("MessageChat")
	token := c.Query("token")
	toUserId := c.Query("to_user_id")

	if user, exist := usersLoginInfo[token]; exist {
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(user.Id, int64(userIdB))

		c.JSON(http.StatusOK, ChatResponse{Response: model.Response{StatusCode: 0}, MessageList: tempChat[chatKey]})
	} else {
		c.JSON(http.StatusOK, model.Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
