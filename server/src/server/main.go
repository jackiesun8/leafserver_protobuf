package main

import (
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"server/conf"
	"server/game"
	_ "server/gamedata"
	"server/gate"
	"server/login"
)

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.ConsolePort = conf.Server.ConsolePort //开启控制台

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
