package dal

import (
	"github.com/RaymondCode/simple-demo/dal/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB = nil

func createTables() error {
	err := mysql.DB.AutoMigrate(
		&Comment{},
		&User{},
		&FavouriteRelation{},
		&VideoMeta{})
	return err
}

func Init() {
	mysql.Init()
	err := createTables()
	if err != nil {
		return
	}
	DB = mysql.DB
}
