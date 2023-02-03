package dal

import "time"

// https://gorm.io/docs/models.html
type User struct {
	ID       int    `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(255)"`
	Username string `gorm:"type:varchar(32)"`
	Password string `gorm:"type:varchar(32)"`
}
type VideoMeta struct {
	ID         int `gorm:"primaryKey"`
	Author     int
	PlayUrl    string    `gorm:"type:varchar(255)"`
	CoverUrl   string    `gorm:"type:varchar(255)"`
	Title      string    `gorm:"type:varchar(255)"`
	CreateTime time.Time `gorm:"autoUpdateTime:milli"`
	UpdateTime time.Time `gorm:"autoCreateTime"`
}
type FavouriteRelation struct {
	UserID  int `gorm:"primaryKey;autoIncrement:false"`
	VideoID int `gorm:"primaryKey;autoIncrement:false"`
}
type Comment struct {
	Id      int `gorm:"primaryKey"`
	Content string
	UserID  int
	VideoID int
}
