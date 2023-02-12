package dal

import (
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
)

// FindUserByUsername 用于注册账号时校验用户名是否唯一
func FindUserByUsername(tx *gorm.DB, username string) (int64, error) {
	tmp := model.User{}
	res := tx.Table("users").Where("username = ?", username).First(&tmp)
	if res.Error != nil {
		return 0, res.Error
	}
	return res.RowsAffected, nil
}

// CreateUser 通过用户名和密码创建账号，用户名称默认和用户名相同，返回用户id
func CreateUser(tx *gorm.DB, username string, password string) (*model.User, error) {
	type registerInf struct {
		Id       int64
		Username string
		Name     string
		Password string
	}
	regInf := registerInf{Username: username, Name: username, Password: password}
	err := tx.Table("users").Create(&regInf).Error

	user := &model.User{Id: regInf.Id, Name: username, FollowCount: 0, FollowerCount: 0}
	return user, err
}

func CheckUser(tx *gorm.DB, username, password string) ([]*model.User, error) {
	res := make([]*model.User, 0)
	if err := tx.Where("username = ? AND password = ?", username, password).First(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func FindUserById(tx *gorm.DB, uid int64) model.User {
	var user model.User
	tx.Table("users").Where("id = ?", uid).First(&user)
	return user
}
