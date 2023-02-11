package dal

import "github.com/RaymondCode/simple-demo/model"

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
