package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"reflect"
	"server/game"
	"server/model"
	"server/msg"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {

	//在这里注册模块的RPC函数，来处理特定的消息
	handleMsg(&msg.RegisterReq{}, handleRegister)
	handleMsg(&msg.LoginReq{}, handleLogin)
}

func handleLogin(args []interface{}) {
	m := args[0].(*msg.LoginReq)
	a := args[1].(gate.Agent)
	log.Debug("%v", m, a)
	account := &model.Account{Name: m.GetName()}
	has, err := model.Engine.Get(account)
	if err != nil {
		a.WriteMsg(&msg.LoginRes{Code: msg.LoginRes_NO.Enum()})
		return
	}
	log.Debug("%v", account)
	if has {
		if account.Password != m.GetPassword() {
			a.WriteMsg(&msg.LoginRes{Code: msg.LoginRes_PASS_ERR.Enum()})
		} else {
			a.WriteMsg(&msg.LoginRes{Code: msg.LoginRes_YES.Enum()})
			game.ChanRPC.Go("UserLogin", a, account.Id)
		}
	} else {
		a.WriteMsg(&msg.LoginRes{Code: msg.LoginRes_NOT_EXIST.Enum()})
	}
}

func handleRegister(args []interface{}) {
	m := args[0].(*msg.RegisterReq)
	a := args[1].(gate.Agent)
	log.Debug("%v", m, a)
	account := new(model.Account)
	account.Name = m.GetName()
	account.Password = m.GetPassword()
	affected, err := model.Engine.Insert(account)
	if err != nil {
		//所有错误都返回用户名已存在
		log.Error("%v", err)
		a.WriteMsg(&msg.RegisterRes{Code: msg.RegisterRes_EXIST.Enum()})
		return
	}
	log.Debug("%v", affected)
	if affected == 1 {
		a.WriteMsg(&msg.RegisterRes{Code: msg.RegisterRes_YES.Enum()})
	} else {
		a.WriteMsg(&msg.RegisterRes{Code: msg.RegisterRes_NO.Enum()})
	}
}
