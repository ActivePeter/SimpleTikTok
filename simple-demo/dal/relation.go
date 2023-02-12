package dal

import (
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
)

type dAORelation struct{}

var DAORelation = dAORelation{}

func (*dAORelation) SetFollow(fromid model.UserId, toid model.UserId, follow bool) error {
	relation := FollowRelation{
		FromID: fromid,
		ToID:   toid,
	}
	if follow {
		err := DB.Create(&relation).Error
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		err := DB.Delete(&relation).Error
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

func SelectFollows(tx *gorm.DB, fromId model.UserId) ([]model.User, error) {
	var users []model.User
	err := tx.Table("follow_relations").
		Select("users.id,users.name,users.follow_count,users.follower_count").
		Where("follow_relations.from_id = ?", fromId).
		Joins("left join users on to_id = users.id").Find(&users).Error
	for i := 0; i < len(users); i++ {
		users[i].IsFollow = true
	}
	return users, err
}

func SelectFollowers(tx *gorm.DB, toId model.UserId) ([]model.User, error) {
	var users []model.User
	err := tx.Table("follow_relations").
		Select("users.id,users.name,users.follow_count,users.follower_count").
		Where("follow_relations.to_id = ?", toId).
		Joins("left join users on from_id = users.id").Find(&users).Error
	for i := 0; i < len(users); i++ {
		res, _ := isFollow(DB, users[i].Id, toId)
		if res == true {
			users[i].IsFollow = true
		}
	}
	return users, err
}

func isFollow(tx *gorm.DB, fromId model.UserId, toId model.UserId) (bool, error) {
	tmp := FollowRelation{}
	res := tx.Table("follow_relations").Where("from_id = ? AND to_id = ?", fromId, toId).First(&tmp)
	if res.Error != nil {
		return false, res.Error
	}
	if res.RowsAffected == 1 {
		return true, nil
	} else {
		return false, nil
	}
}
