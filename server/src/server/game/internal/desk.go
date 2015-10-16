package internal

import ()

var (
	playerNum = 5 //牌桌允许最大玩家数
)

type Desk struct {
	Id    int32
	poker *Poker
}

//新建一个桌子
func NewDesk() *Desk {
	desk := new(Desk)
	desk.poker = new(Poker)
	return desk
}
