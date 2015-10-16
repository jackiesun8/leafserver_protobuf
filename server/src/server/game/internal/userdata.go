package internal

import (
	"fmt"
	"github.com/name5566/leaf/log"
	"server/model"
)

type UserData struct {
	*model.User
}

func (data *UserData) initValue(accID int64) error {

	user := new(model.User)
	user.AccId = accID
	affected, err := model.Engine.Insert(user)
	if err != nil {
		return fmt.Errorf("insert user error: %v", err)
	}
	if affected == 1 {
		data.User = user
	} else {
		return fmt.Errorf("insert user failed")
	}
	return nil
}

func (data *UserData) saveDB() error {
	//save User

	id := data.User.Id
	affected, err := model.Engine.Id(id).AllCols().Omit("Id", "AccId").Update(data.User)
	if err != nil {
		return fmt.Errorf("save db error: %v", err)
	}
	if affected != 1 {
		log.Error("%v", "save db no affect,may be an error")
	} else {
		log.Release("%v", "save db ok")
	}
	return nil
}
