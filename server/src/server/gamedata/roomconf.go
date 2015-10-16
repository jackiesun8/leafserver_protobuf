package gamedata

import ()

var (
	RoomIdRoomConfs = make(map[int32]*RoomConf)
)

type RoomConf struct {
	Id            int32
	Name          string
	Remark        string
	ConditionGold int32
	DeskNum       int32
}

func init() {
	rf := readRf(RoomConf{})
	for i := 0; i < rf.NumRecord(); i++ {
		r := rf.Record(i).(*RoomConf)
		RoomIdRoomConfs[r.Id] = r
	}
}
