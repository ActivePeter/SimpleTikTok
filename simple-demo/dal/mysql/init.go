package mysql

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(config *utils.ServerConfig) {
	var err error
	fmt.Printf("config:%v\n", config)
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?parseTime=True",
		config.SqlUser, config.SqlPw, config.SqlAddr, config.Schema)
	println(dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//SkipDefaultTransaction: true, //跳过默认开启事务模式
		//PrepareStmt: true, //在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率
		//Logger:      logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		// todo 异常处理
		panic(err)
	}
}
