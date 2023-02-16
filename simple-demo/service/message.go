package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/RaymondCode/simple-demo/model"
	"net"
	"sync"
)

var chatConnMap = sync.Map{}
var connections = make(map[int64]net.Conn)

func RunMessageServer() {
	listen, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		fmt.Printf("Run message sever failed: %v\n", err)
		return
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("Accept conn failed: %v\n", err)
			continue
		}

		go process(conn)
	}
}

func process(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		//读取消息
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error reading", err.Error())
			break
		}

		// 解析消息
		event := model.MessageSendEvent{}
		err = json.Unmarshal(buf, &event)
		if err != nil {
			fmt.Println("Error parsing message", err.Error())
			continue
		}

		// 消息内容为空表示用户上线，将连接存入map中
		if event.MsgContent == "" {
			connections[event.UserId] = conn
			fmt.Printf("User %d online\n", event.UserId)
			continue
		}

		// 从map中取出对应的连接并发送消息
		if toConn, ok := connections[event.ToUserId]; ok {
			pushEvent := model.MessagePushEvent{
				FromUserId: event.UserId,
				MsgContent: event.MsgContent,
			}
			pushData, _ := json.Marshal(pushEvent)
			_, err = toConn.Write(pushData)
			if err != nil {
				fmt.Println("Error writing", err.Error())
				continue
			}
			fmt.Printf("User %d send message to user %d: %s\n", event.UserId, event.ToUserId, event.MsgContent)
		} else {
			fmt.Printf("User %d offline\n", event.ToUserId)
		}
	}
}

//func process(conn net.Conn) {
//	defer conn.Close()
//
//	var buf [256]byte
//	for {
//		n, err := conn.Read(buf[:])
//		if n == 0 {
//			if err == io.EOF {
//				break
//			}
//			fmt.Printf("Read message failed: %v\n", err)
//			continue
//		}
//
//		var event = model.MessageSendEvent{}
//		_ = json.Unmarshal(buf[:n], &event)
//		fmt.Printf("Receive Message：%+v\n", event)
//
//		fromChatKey := fmt.Sprintf("%d_%d", event.UserId, event.ToUserId)
//		if len(event.MsgContent) == 0 {
//			chatConnMap.Store(fromChatKey, conn)
//			continue
//		}
//
//		toChatKey := fmt.Sprintf("%d_%d", event.ToUserId, event.UserId)
//		writeConn, exist := chatConnMap.Load(toChatKey)
//		if !exist {
//			fmt.Printf("User %d offline\n", event.ToUserId)
//			continue
//		}
//
//		pushEvent := model.MessagePushEvent{
//			FromUserId: event.UserId,
//			MsgContent: event.MsgContent,
//		}
//		pushData, _ := json.Marshal(pushEvent)
//		_, err = writeConn.(net.Conn).Write(pushData)
//		if err != nil {
//			fmt.Printf("Push message failed: %v\n", err)
//		}
//	}
//}
