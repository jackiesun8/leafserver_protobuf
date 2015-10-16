package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine
var Engine *xorm.Engine

func init() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:1234qwer@/my_db?charset=utf8")
	if err != nil {
		panic(err)
	}

	engine.ShowSQL = true

	engine.ShowInfo = true
	engine.ShowErr = true
	engine.ShowDebug = true
	engine.ShowWarn = true

	err = engine.Sync2(new(Account), new(User))
	if err != nil {
		panic(err)
	}
	Engine = engine
}
