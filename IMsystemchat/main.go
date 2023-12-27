package main

import (
	"IMsystemchat/router"
	"IMsystemchat/utils"
)

func main() {
	//初始化配置文件
	utils.InitConfig()
	//初始化数据操作
	utils.InitMYSQL()
	//初始化redis
	utils.InitRedis()
	r := router.Router()
	r.Run(":8081")
}

//学完了p34集了，接下来学P35集了
