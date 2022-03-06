package serve

import (
	"fmt"
	"log"
	"net"
)

// 监听的服务器
func ServerMain(config *Config) {
	// 创建一个tcp监听
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", config.Port))
	if nil != err {
		log.Panic(err)
	}

	defer l.Close()

	log.Println("listen", config.Port, "success.")
	// 等待连接
	for {
		conn, err := l.Accept()
		if nil != err {
			log.Panic(err)
		}

		go NewSimpleSshd(config).Serve(conn)
	}
}
