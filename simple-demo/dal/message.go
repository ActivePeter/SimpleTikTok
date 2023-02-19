package dal

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
	"strconv"
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
func GetLatestMessage(userId model.UserId, friendId model.UserId) (LatestMessage, error) {
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

// 获取最新的消息
func GetMessages(userIdA model.UserId, UserIdB model.UserId, preTime string) ([]model.Msg, error) {
	rows := make([]model.Msg, 0)
	res := make([]model.Msg, 0)
	err := DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Debug().Table("messages").
			Where("to_user_id IN (?) AND from_user_id IN (?)", []int64{userIdA, UserIdB}, []int64{userIdA, UserIdB}).
			Order("create_time").Find(&rows)
		if err != nil {
			return err.Error
		}
		return nil
	})
	for _, row := range rows {
		//转换为UNIX时间戳
		if createTime_unix, err := TimeToUnix(row.CreateTime); err != nil {
			return nil, err
		} else {
			createTime, _ := strconv.Atoi(createTime_unix)
			preMsgTime, _ := strconv.Atoi(preTime)
			if createTime > preMsgTime {
				r := model.Msg{
					Id:         row.Id,
					ToUserId:   row.ToUserId,
					FromUserId: row.FromUserId,
					Content:    row.Content,
					CreateTime: createTime_unix,
				}
				res = append(res, r)
			}
		}
	}
	return res, err
}
func TimeToUnix(originTime string) (string, error) {
	layout := "2006-01-02 15:04:05.9999999 -0700 MST"
	//分割字符串，只取有用的时间信息
	splitTime := strings.Split(originTime, " m=")[0]
	//先转换为time格式
	if time_timeFmt, err := time.Parse(layout, splitTime); err != nil {
		return "", err
	} else {
		//转换为Unix时间戳
		unixTime := strconv.FormatInt(time_timeFmt.Unix(), 10)
		return unixTime, nil
	}
}
