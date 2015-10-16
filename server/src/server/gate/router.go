package gate

import (
	"server/game"
	"server/login"
	"server/msg"
)

func init() {

	//在这里设置路由，将特定的消息路由到特定的的模块RPC内

	// login
	msg.ProtobufProcessor.SetRouter(&msg.RegisterReq{}, login.ChanRPC)
	msg.ProtobufProcessor.SetRouter(&msg.LoginReq{}, login.ChanRPC)

	// game
	msg.ProtobufProcessor.SetRouter(&msg.RoomListReq{}, game.ChanRPC)
	msg.ProtobufProcessor.SetRouter(&msg.EnterRoomReq{}, game.ChanRPC)
	msg.ProtobufProcessor.SetRouter(&msg.ExitRoomReq{}, game.ChanRPC)
}
