package model

type Account struct {
	Id       int64
	Name     string `xorm:"varchar(25) notnull unique"` //账户名称
	Password string //账户密码
}
