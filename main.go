package main

import (
	"efficient-sshd/serve"
	"flag"
)

func main() {
	config := &serve.Config{}

	flag.IntVar(&config.Port, "port", 8022, "listen tcp port");
	flag.Parse()

	// 打印信息
	config.Dump()

	// 启动服务
	serve.ServerMain(config)
}
