package main

import (
	"os"
	"os/signal"
	"robot/elog"
	"robot/robot"
	"server/msg"
	"strconv"
	"syscall"
	//"encoding/binary"
	//"bytes"
	//"io"
)

const defRobotNum = 1

const (
	//用户状态
	kStateIdle        = iota //空闲状态
	kStateReqLogining        //请求登陆状态
	kStateReqRegister        //请求注册状态
	kStateLogin              //
	kStateMax
)

type C2GLogon struct {
	name string
	pwd  string
}

type STATE_CB func(user *User) //状态回调

type User struct {
	name string
	pwd  string

	state int

	r *robot.Robot

	stateCbs [kStateMax]STATE_CB
}

func (u *User) RegisterCb() {

	u.stateCbs[kStateIdle] = onIdle               //空闲状态回调
	u.stateCbs[kStateReqLogining] = onReqLogining //请求登陆回调
	//u.stateCbs[kStateReqRegister] = onReqRegister
	//u.stateCbs[kStateLogin]       = onLogin
}

func onIdle(u *User) {

	// var msg robot.Message

	// //binary.Write( buf, binary.LittleEndian,  u.name )
	// //binary.Write( buf, binary.LittleEndian,  u.pwd )

	// msg.Data = make([]byte, 96)

	// userName := msg.Data[0:64]
	// pwd := msg.Data[64:]

	// copy(userName, []byte(u.name))
	// copy(pwd, []byte(u.pwd))
	// msg.SetId(1)
	// u.r.SendMsg(msg)

	// elog.LogInfo(" send msg :%d ", len(msg.Data))
	// u.state = kStateReqLogining
}

func onReqLogining(u *User) {

	// msg, err := u.r.RecvMsg()
	// if err != nil {
	// 	//elog.LogInfo(" get msg fail ")
	// 	return
	// } else {
	// 	elog.LogInfo(" i receieve msg :%d ", msg.GetId())
	// }
}

//更新回调函数定义
func process(r *robot.Robot) {

	if r.UserData == nil { //用户数据为空

		elog.LogInfoln(" user data ", r.UserData)
		r.UserData = &User{ //创建用户
			name:  strconv.Itoa(int(r.GetId())),
			pwd:   "123",
			state: kStateIdle, //初始状态为空闲
			r:     r,
		}
		user := r.UserData.(*User)
		user.RegisterCb() //注册状态回调
	}

	//elog.LogInfoln( " user data ", r.UserData )
	user := r.UserData.(*User)
	user.stateCbs[user.state](user) //根据用户状态执行回调
}

func main() {

	elog.InitLog(elog.INFO)
	elog.LogInfo(" --------------------- let go ------------------------")
	robotMng := robot.NewRobotMng() //创建机器人管理器
	robotMng.SetUpdateCb(process)   //设置更新回调
	robotMng.Run()                  //启动机器人管理器
	//系统退出信号处理
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	robotMng.Close() //关闭机器人管理器
}
