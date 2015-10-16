package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"server/msg"
)

type AgentInfo struct {
	accID  int64
	userID int64
}

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
	skeleton.RegisterChanRPC("UserLogin", rpcUserLogin)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)

	a.SetUserData(new(AgentInfo)) //Agent的用户数据存放的是AgentInfo
}

func rpcUserLogin(args []interface{}) {
	a := args[0].(gate.Agent)
	accID := args[1].(int64)

	// network closed
	if a.UserData() == nil {
		return
	}

	// login repeated
	oldSession := accIDSessions[accID]
	if oldSession != nil {
		m := &msg.ClosePush{Code: msg.ClosePush_REPEATED.Enum()}
		a.WriteMsg(m)
		a.Close()
		oldSession.WriteMsg(m)
		oldSession.Close()
		log.Debug("acc %v login repeated", accID)
		return
	}

	log.Debug("acc %v login", accID)

	// login
	newSession := new(Session)
	newSession.Agent = a
	newSession.LinearContext = skeleton.NewLinearContext()
	newSession.state = Session_Login
	a.UserData().(*AgentInfo).accID = accID
	accIDSessions[accID] = newSession

	newSession.login(accID)
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)

	accID := a.UserData().(*AgentInfo).accID
	a.SetUserData(nil)

	session := accIDSessions[accID]
	if session == nil {
		return
	}

	log.Debug("acc %v logout", accID)

	// logout
	if session.state == Session_Login { //还在登录中
		session.state = Session_Logout //仅仅设置标志，待登录完成执行logout操作
	} else {
		session.state = Session_Logout
		session.logout(accID)
	}
}
