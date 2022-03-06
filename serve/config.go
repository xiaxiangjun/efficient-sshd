package serve

import (
	"encoding/base64"
	"log"
	"math/rand"
	"strings"
	"time"
)

type Config struct {
	Port   int    // 监听的端口
	Home   string // 主目录，默认为当前目录
	Passwd string // 密码，默认为随机密码
}

// 打印相关的信息
func (self *Config) Dump() {
	log.Println("listen port: ", self.Port)
	log.Println("use home: ", self.Home)
	log.Println("password: ", self.Passwd)
}

// 生成随机密码
func RandomPassword() string {
	pwd := make([]byte, 6)

	rand.Seed(time.Now().UnixNano())
	rand.Read(pwd)
	return strings.ReplaceAll(base64.StdEncoding.EncodeToString(pwd), "=", "")
}
