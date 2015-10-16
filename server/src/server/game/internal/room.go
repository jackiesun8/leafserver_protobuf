package internal

import (
	"server/gamedata"
)

var (
	roomIdRooms = make(map[int32]*Room)
)

type Room struct {
	*gamedata.RoomConf
	players map[int64]*Player //key:userId
	desks   map[int]*Desk     //key:deskId
}

//初始化房间列表
func initRoomList() {
	for _, v := range gamedata.RoomIdRoomConfs {
		room := new(Room)
		room.RoomConf = v
		room.players = make(map[int64]*Player)
		room.desks = make(map[int]*Desk)

		var i int32 = 0
		for ; i < room.DeskNum; i++ {
			desk := NewDesk()
			desk.Id = i
		}
		roomIdRooms[room.Id] = room
	}
}

func (room *Room) EnterRoom() {

}
