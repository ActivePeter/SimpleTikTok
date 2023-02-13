package model

type Response struct {
	StatusCode int32  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg,omitempty"`
}

type VideoId = int64
type Video struct {
	Id            VideoId `json:"id,omitempty"`
	Author        User    `json:"author"`
	PlayUrl       string  `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string  `json:"cover_url,omitempty"`
	FavoriteCount int64   `json:"favorite_count,omitempty"`
	CommentCount  int64   `json:"comment_count,omitempty"`
	IsFavorite    bool    `json:"is_favorite,omitempty"`
}

type CommentId = int64
type Comment struct {
	Id         CommentId `json:"id,omitempty"`
	User       User      `json:"user"`
	Content    string    `json:"content,omitempty"`
	CreateDate string    `json:"create_date,omitempty"`
}

type UserId = int64
type User struct {
	Id            UserId `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type MessageId = int64
type Message struct {
	Id         MessageId `json:"id,omitempty"`
	Content    string    `json:"content,omitempty"`
	CreateTime string    `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     UserId `json:"user_id,omitempty"`
	ToUserId   UserId `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId UserId `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type FavouriteRelation struct {
	UserID  int `gorm:"primaryKey;autoIncrement:false"`
	VideoID int `gorm:"primaryKey;autoIncrement:false"`
}

type FriendUser struct {
	User
	Avatar  string `json:"avatar,omitempty"`
	Message string `json:"message,omitempty"`
	MsgType int64  `json:"msgType,omitempty"`
}

type Msg struct {
	//Id         int64  `json:"id,omitempty"`
	ToUserId   UserId `json:"to_user_id,omitempty"`
	FromUserId UserId `json:"from_user_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}
