package dal

import (
	"github.com/RaymondCode/simple-demo/model"
	"time"
)

// https://gorm.io/docs/models.html
type User struct {
	ID       model.UserId `gorm:"primaryKey"`
	Name     string       `gorm:"type:varchar(255)"`
	Username string       `gorm:"type:varchar(32)"`
	Password string       `gorm:"type:varchar(32)"`
}
type VideoMeta struct {
	ID         model.UserId `gorm:"primaryKey"`
	Author     model.UserId
	PlayUrl    string    `gorm:"type:varchar(255)"`
	CoverUrl   string    `gorm:"type:varchar(255)"`
	Title      string    `gorm:"type:varchar(255)"`
	CreateTime time.Time `gorm:"autoUpdateTime:milli"`
	UpdateTime time.Time `gorm:"autoCreateTime"`
}
type FavouriteRelation struct {
	UserID  model.UserId  `gorm:"primaryKey;autoIncrement:false"`
	VideoID model.VideoId `gorm:"primaryKey;autoIncrement:false"`
}
type Comment struct {
	Id      model.CommentId `gorm:"primaryKey"`
	Content string
	UserID  model.UserId
	VideoID model.VideoId
}
type FollowRelation struct {
	FromID model.UserId `gorm:"primaryKey;autoIncrement:false"`
	ToID   model.UserId `gorm:"primaryKey;autoIncrement:false"`
}
