// 统一在此编写数据库访问语句

package mysql

import (
	//"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
)

// FindUserByUsername 用于注册账号时校验用户名是否唯一
func FindUserByUsername(username string) (int64, error) {
	tmp := model.User{}
	res := DB.Table("users").Where("username = ?", username).First(&tmp)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

// CreateUser 通过用户名和密码创建账号，用户名称默认和用户名相同，返回用户id
func CreateUser(username string, password string) (int64, error) {
	type registerInf struct {
		Id       int64
		Username string
		Name     string
		Password string
	}
	user := registerInf{Username: username, Name: username, Password: password}
	return user.Id, DB.Table("users").Create(&user).Error
}

func CheckUser(username, password string) ([]*model.User, error) {
	res := make([]*model.User, 0)
	if err := DB.Where("username = ? AND password = ?", username, password).First(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func UserFollowingCnt(tx *gorm.DB, uid int) {
	//var cnt int
	//tx.Model()
}

//func CreateUser(users []*model.User) error {
//	return DB.Create(users).Error
//}
//
//func CreateUsers(users []*model.User) error {
//	return DB.Create(users).Error
//}
//
//func FindUserByNameOrEmail(userName, email string) ([]*model.User, error) {
//	res := make([]*model.User, 0)
//	if err := DB.Where(DB.Or("user_name = ?", userName).
//		Or("email = ?", email)).
//		Find(&res).Error; err != nil {
//		return nil, err
//	}
//	return res, nil
//}
//
//func CheckUser(account, password string) ([]*model.User, error) {
//	res := make([]*model.User, 0)
//	if err := DB.Where(DB.Or("user_name = ?", account).
//		Or("email = ?", account)).Where("password = ?", password).
//		Find(&res).Error; err != nil {
//		return nil, err
//	}
//	return res, nil
//}
