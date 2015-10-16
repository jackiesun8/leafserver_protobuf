package internal

import (
	"github.com/golang/protobuf/proto"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"reflect"
	"server/msg"
)

func handleMsg(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), func(args []interface{}) {
		// get session
		a := args[1].(gate.Agent)
		session := userIDSessions[a.UserData().(*AgentInfo).userID]
		if session == nil {
			return
		}
		// reset agent to session
		args[1] = session
		h.(func([]interface{}))(args)
	})
}

func init() {
	handleMsg(&msg.RoomListReq{}, handleRoomList)
}

func handleRoomList(args []interface{}) {
	m := args[0].(*msg.RoomListReq)
	a := args[1].(*Session)
	log.Debug("%v", m, a)
	res := &msg.RoomListRes{}
	for _, v := range roomIdRooms {
		res.Rooms = append(res.Rooms, &msg.RoomListRes_RoomInfo{Id: proto.Int32(v.Id), PlayerNum: proto.Int32(int32(len(v.players)))})
	}

}
func handleEnterRoom(args []interface{}) {
	m := args[0].(*msg.EnterRoomReq)
	a := args[1].(*Session)
	log.Debug("%v", m, a)
}
