package dal

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
	"strings"
	"time"
)

func AddMessage(tx *gorm.DB, msg model.Msg) bool {
	if err := tx.Table("messages").Select("ToUserId", "FromUserId", "Content", "CreateTime").Create(&msg).Error; err != nil {
		return false
	} else {
		return true
	}
}

type LatestMessage struct {
	Message string
	MsgType int64
}

// 获取最新一条消息并判断是收到的还是发送的
func GetLatestMessage(tx *gorm.DB, userId model.UserId, friendId model.UserId) (LatestMessage, error) {
	latestMessage := LatestMessage{
		Message: "",
		MsgType: -1,
	}
	send_message := model.Msg{}
	receive_message := model.Msg{}

	err := DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().Table("messages").
			Where("from_user_id = ? AND to_user_id = ?", userId, friendId).
			Order("create_time desc").First(&send_message)
		if err != nil {
			return err.Error
		}
		return nil
	})

	err = DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().Table("messages").
			Where("to_user_id = ? AND from_user_id = ?", userId, friendId).
			Order("create_time desc").First(&receive_message)
		if err != nil {
			return err.Error
		}
		return nil
	})

	if send_message.CreateTime == "" && receive_message.CreateTime == "" {
		return latestMessage, nil
	} else if send_message.CreateTime == "" {
		latestMessage.Message = receive_message.Content
		latestMessage.MsgType = 0
		return latestMessage, nil
	} else if receive_message.CreateTime == "" {
		latestMessage.Message = send_message.Content
		latestMessage.MsgType = 1
		return latestMessage, nil
	}

	//string转time
	layout := "2006-01-02 15:04:05.9999999 -0700 MST"
	sendTime := send_message.CreateTime
	receiveTime := receive_message.CreateTime
	sendTime = strings.Split(sendTime, " m=")[0]
	receiveTime = strings.Split(receiveTime, " m=")[0]
	sendTime_timeFmt, err := time.Parse(layout, sendTime)
	if err != nil {
		fmt.Println("Failed to parse time:", err)
		return latestMessage, nil
	}
	receiveTime_timeFmt, err := time.Parse(layout, receiveTime)
	if err != nil {
		fmt.Println("Failed to parse time:", err)
		return latestMessage, nil
	}

	//对比时间先后
	if sendTime_timeFmt.After(receiveTime_timeFmt) {
		latestMessage.Message = send_message.Content
		latestMessage.MsgType = 1
		return latestMessage, nil
	} else {
		latestMessage.Message = receive_message.Content
		latestMessage.MsgType = 0
		return latestMessage, nil
	}
}
func GetMassges(userId model.UserId, toUserId model.UserId) ([]model.Msg, error) {
	res := make([]model.Msg, 0)
	err := DB.Transaction(func(tx *gorm.DB) error {

		err := tx.Debug().Table("messages").
			Where("to_user_id = ? AND from_user_id = ?", toUserId, userId).Find(&res)
		if err != nil {
			return err.Error
		}
		return nil
	})
	return res, err
}
