package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn = "simple_tiktok:123456@tcp(43.143.166.162:3306)/simple_tiktok?parseTime=True"

var DB *gorm.DB

func Init() {
	var err error
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
