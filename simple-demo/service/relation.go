package service

import (
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/model"
	"log"
)

type relation struct{}

var Relation = relation{}

func (*relation) SetFollow(from model.UserId, to model.UserId, follow bool) error {
	return dal.DAORelation.SetFollow(from, to, follow)
}

func (*relation) GetFriendList(from model.UserId) ([]model.FriendUser, error) {
	var friends []model.FriendUser
	//根据userid获取User列表
	if users, err := dal.GetFriendList(dal.DB, from); err != nil {
		log.Default().Println("获取好友列表失败!")
		return nil, err
	} else {
		url := "https://tse4-mm.cn.bing.net/th/id/OIP-C.kHtLyqTcn4yBGCUcWQxvcwHaHa?pid=ImgDet&rs=1"
		//将User封装为FriendUser
		for _, user := range users {
			fu := model.FriendUser{
				User:    user,
				Avatar:  url,
				MsgType: 0,
			}
			friends = append(friends, fu)
		}
		return friends, nil
	}
}
