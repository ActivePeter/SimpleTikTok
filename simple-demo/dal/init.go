package dal

import (
	"github.com/RaymondCode/simple-demo/dal/mysql"
	"github.com/RaymondCode/simple-demo/dal/redisConfig"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var DB *gorm.DB = nil

var RD *redis.Client

func createTables() error {
	err := mysql.DB.AutoMigrate(
		&Comment{},
		&User{},
		&FavouriteRelation{},
		&VideoMeta{},
		&FollowRelation{})
	return err
}

func Init() {
	mysql.Init()
	redisConfig.Init()
	err := createTables()
	if err != nil {
		return
	}
	DB = mysql.DB
	RD = redisConfig.RD
}
