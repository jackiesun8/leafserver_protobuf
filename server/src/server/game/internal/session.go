package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/go"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/timer"
	"github.com/name5566/leaf/util"
	"server/model"
	"server/msg"
	"time"
)

var (
	accIDSessions  = make(map[int64]*Session)
	userIDSessions = make(map[int64]*Session)
)

//用户状态
const (
	Session_Login  = iota //登录中
	Session_Logout        //已登出，客户端主动断开连接
	Session_InGame        //游戏中
)

type Session struct {
	gate.Agent
	*g.LinearContext
	state       int
	data        *UserData
	saveDBTimer *timer.Timer
}

func (session *Session) login(accID int64) {
	userData := new(UserData)
	skeleton.Go(func() {

		user := new(model.User)
		user.AccId = accID
		has, err := model.Engine.Get(user)

		if err != nil {
			log.Error("load acc %v data error: %v", accID, err)
			userData = nil
			m := &msg.ClosePush{Code: msg.ClosePush_INNER.Enum()}
			session.WriteMsg(m)
			session.Close()
			return
		}
		if has {
			userData.User = user
		} else {
			err := userData.initValue(accID)
			if err != nil {
				log.Error("init acc %v data error: %v", accID, err)
				userData = nil
				m := &msg.ClosePush{Code: msg.ClosePush_INNER.Enum()}
				session.WriteMsg(m)
				session.Close()
				return
			}
		}
	}, func() {
		// network closed
		if session.state == Session_Logout {
			session.logout(accID)
			return
		}

		// db error
		session.state = Session_InGame
		if userData == nil {
			return
		}

		// ok
		session.data = userData
		userIDSessions[userData.User.Id] = session
		session.UserData().(*AgentInfo).userID = userData.User.Id
		session.onLogin()
		session.autoSaveDB()
	})
}

func (session *Session) logout(accID int64) {
	if session.data != nil {
		session.saveDBTimer.Stop()
		session.onLogout()
		delete(userIDSessions, session.data.User.Id)
	}

	// save
	data := util.DeepClone(session.data)
	session.Go(func() {
		if data != nil {
			//db := mongoDB.Ref()
			//defer mongoDB.UnRef(db)
			// userID := data.(*UserData).UserID
			// _, err := db.DB("game").C("users").
			// 	UpsertId(userID, data)

			//err := data.(*UserData).saveDB()
			err := session.data.saveDB()

			if err != nil {
				log.Error("save user %v data error: %v", session.data.User.Id, err)
			}
		}
	}, func() {
		delete(accIDSessions, accID)
	})
}

func (session *Session) autoSaveDB() {
	const duration = 5 * time.Minute

	// save
	session.saveDBTimer = skeleton.AfterFunc(duration, func() {
		data := util.DeepClone(session.data)
		session.Go(func() {
			//db := mongoDB.Ref()
			//defer mongoDB.UnRef(db)
			// userID := data.(*UserData).UserID
			// _, err := db.DB("game").C("users").
			// 	UpsertId(userID, data)
			err := data.(*UserData).saveDB()
			if err != nil {
				log.Error("save user %v data error: %v", session.data.User.Id, err)
			}
		}, func() {
			session.autoSaveDB()
		})
	})
}

func (session *Session) isOffline() bool {
	return session.state == Session_Logout
}

func (session *Session) onLogin() {

}

func (session *Session) onLogout() {

}
