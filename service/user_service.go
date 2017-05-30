package service

import (
	"github.com/solarknight/simple_golang_web/common"
	"github.com/solarknight/simple_golang_web/dao"
)

func ConnectDemo() {
	dao.ConnectDB()
}

func QueryByID(id int) (*common.User, error) {
	db, err := dao.OpenDB()
	if err != nil {
		return nil, err
	}
	return dao.GetByIDSimpler(db, id)
}
