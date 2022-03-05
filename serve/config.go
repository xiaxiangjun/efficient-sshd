package serve

import "log"

type Config struct {
	Port int
}

// 打印相关的信息
func (self *Config) Dump() {
	log.Println("listen port: ", self.Port)
}