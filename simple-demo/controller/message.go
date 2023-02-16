package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/cloudwego/hertz/pkg/app"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

var tempChat = map[string][]model.Msg{}

var messageIdSequence = int64(1)

type ChatResponse struct {
	model.Response
	MessageList []model.Msg `json:"message_list"`
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

func MessageChat(ctx context.Context, c *app.RequestContext) {
	log.Default().Println("MessageChat")
	//验证用户
	user, exists := service.GetUserFromContext(c)
	if !exists {
		log.Default().Println("用户不存在!")
		return
	}
	//连接消息服务器
	conn, err := net.Dial("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Println("Dial error: ", err)
		return
	}
	defer conn.Close()

	userId := user.Id
	toUserId := c.Query("to_user_id")
	userIdB, _ := strconv.Atoi(toUserId)
	chatKey := genChatKey(userId, int64(userIdB))

	// 读取消息
	if online, err := dal.GetMessages(userId, int64(userIdB)); err != nil {
		return
	} else {
		for _, o := range online {
			onlineEvent := model.MessageSendEvent{
				UserId:     o.FromUserId,
				ToUserId:   o.ToUserId,
				MsgContent: o.Content,
			}
			onlineData, _ := json.Marshal(onlineEvent)
			_, err = conn.Write(onlineData)
			if err != nil {
				fmt.Println("Error writing:", err.Error())
				return
			}
		}
	}
	//
	//// 读取用户输入并发送私信
	//scanner := bufio.NewScanner(os.Stdin)
	//for scanner.Scan() {
	//	msg := scanner.Text()
	//	if msg == "" {
	//		continue
	//	}
	//	sendEvent := model.MessageSendEvent{
	//		UserId:     userId,
	//		ToUserId:   int64(userIdB),
	//		MsgContent: msg,
	//	}
	//	sendData, _ := json.Marshal(sendEvent)
	//	_, err = conn.Write(sendData)
	//	if err != nil {
	//		fmt.Println("Error writing:", err.Error())
	//		return
	//	}
	//}
	c.JSON(http.StatusOK, ChatResponse{Response: model.Response{StatusCode: 0}, MessageList: tempChat[chatKey]})
}

// MessageChat all users have same follow list
//func MessageChat(ctx context.Context, c *app.RequestContext) {
//	log.Default().Println("MessageChat")
//	//验证用户
//	user, exists := service.GetUserFromContext(c)
//	if !exists {
//		log.Default().Println("用户不存在!")
//		return
//	}
//	userId := user.Id
//	to_user_id := c.Query("to_user_id")
//	toUserId, err := strconv.ParseInt(to_user_id, 10, 64)
//	if err != nil {
//		c.JSON(http.StatusOK, model.Response{
//			StatusCode: 1,
//			StatusMsg:  err.Error(),
//		})
//	}
//	if messages, err := dal.GetMassges(userId, toUserId); err != nil {
//		fmt.Println(err)
//		c.JSON(http.StatusBadRequest, model.Response{
//			StatusCode: 1,
//			StatusMsg:  err.Error(),
//		})
//	} else {
//		c.JSON(http.StatusOK, ChatResponse{
//			Response: model.Response{
//				StatusCode: 0,
//			},
//			MessageList: messages,
//		})
//	}
//}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
