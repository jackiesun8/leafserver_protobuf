package model

type User struct {
	Id       int64
	AccId    int64  `xorm:"notnull"` //账户ID
	IconId   int    //头像ID
	Nickname string `xorm:"default 'hello,world'"` //昵称
	Gold     int    //金币
	Diamond  int    //钻石
}
