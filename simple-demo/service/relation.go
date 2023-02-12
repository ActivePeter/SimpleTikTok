package service

import (
	"github.com/RaymondCode/simple-demo/dal"
	"github.com/RaymondCode/simple-demo/model"
)

type relation struct{}

var Relation = relation{}

func (*relation) SetFollow(from model.UserId, to model.UserId, follow bool) error {
	return dal.DAORelation.SetFollow(from, to, follow)
}
