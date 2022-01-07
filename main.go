package main

import (
	"log"
	"net/http"
	"os"
	"ranking/config"
	"ranking/router"
)

func main() {
	// 加载配置
	loadConfiguration()

	// 注册路由
	router.RegisterRoutes()
	log.Println("starting ~~~")

	// 启动http-server
	_ = http.ListenAndServe(":8090", nil)
}

// 加载配置（注意命令行参数第一个固定是进程名）
func loadConfiguration() {
	args := os.Args
	var env string
	if len(args) >= 2 {
		env = args[1]
	}
	err := config.LoadConfiguration(env)
	if err != nil {
		log.Fatalln("load configuration error : ", err)
	}
}
