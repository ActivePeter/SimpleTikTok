package dal

import (
	"github.com/RaymondCode/simple-demo/model"
	"gorm.io/gorm"
	"log"
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

func SelectFollowsNum(tx *gorm.DB, fromId model.UserId) (int64, error) {
	var count int64
	err := tx.Table("follow_relations").
		Select("users.id,users.name,users.follow_count,users.follower_count").
		Where("follow_relations.from_id = ?", fromId).
		Joins("left join users on to_id = users.id").Count(&count).Error
	return count, err
}

func SelectFollowers(tx *gorm.DB, toId model.UserId) ([]model.User, error) {
	var users []model.User
	err := tx.Table("follow_relations").
		Select("users.id,users.name,users.follow_count,users.follower_count").
		Where("follow_relations.to_id = ?", toId).
		Joins("left join users on from_id = users.id").Find(&users).Error
	for i := 0; i < len(users); i++ {
		res, _ := isFollow(DB, toId, users[i].Id)
		if res == true {
			users[i].IsFollow = true
		}
	}
	return users, err
}

func SelectFollowersNum(tx *gorm.DB, toId model.UserId) (int64, error) {
	var count int64
	err := tx.Table("follow_relations").
		Select("users.id,users.name,users.follow_count,users.follower_count").
		Where("follow_relations.to_id = ?", toId).
		Joins("left join users on from_id = users.id").Count(&count).Error
	return count, err
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

func GetFriendList(tx *gorm.DB, fromId model.UserId) ([]model.User, error) {
	log.Default().Println("GetFriendList")
	var users []model.User
	sql := "select id,username name,follow_count,follower_count " +
		"from users " +
		"where id in " +
		"(" +
		"select a.to_id " +
		"from follow_relations a " +
		"join follow_relations b " +
		"on a.from_id = b.to_id " +
		"and a.to_id = b.from_id " +
		"where a.from_id = ? " +
		")"
	if rows, err := tx.Raw(sql, fromId).Rows(); err != nil {
		return nil, err
	} else {
		for rows.Next() {
			tmp := model.User{}
			rows.Scan(&tmp.Id, &tmp.Name, &tmp.FollowCount, &tmp.FollowerCount)
			tmp.IsFollow = true
			users = append(users, tmp)
		}
		return users, nil
	}

	//tx.Model(&model.User{}).
	//	Select("id,username,follow_count,follower_count").
	//	Where("id in ?",
	//		tx.Table("follow_relations a").
	//		Select("to_id").
	//		Joins("join follow_relations b on a.from_id = b.to_id and a.to_id = b.from_id").
	//		Where("a.from_id = ?",formId).
	//	).Find(users)
}
