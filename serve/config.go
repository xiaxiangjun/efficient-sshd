package serve

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Config struct {
	Port   int    // 监听的端口
	Home   string // 主目录，默认为当前目录
	Passwd string // 密码，默认为随机密码
	allKey []ssh.PublicKey
}

// 生成随机密码
func RandomPassword() string {
	pwd := make([]byte, 6)

	rand.Seed(time.Now().UnixNano())
	rand.Read(pwd)
	return strings.ReplaceAll(base64.StdEncoding.EncodeToString(pwd), "=", "")
}

// 打印相关的信息
func (self *Config) Dump() {
	log.Println("listen port: ", self.Port)
	log.Println("use home: ", self.Home)
	log.Println("password: ", self.Passwd)
}

func (self *Config) LoadPublicKey() {
	// 获取主目录
	home := self.Home
	if "" == home {
		home, _ = os.Getwd()
	}

	// 读取`~/.ssh/authorized_keys`文件
	buf, err := ioutil.ReadFile(filepath.Join(home, ".ssh", "authorized_keys"))
	if nil != err {
		log.Println("load ~/.ssh/authorized_keys error: ", err)
		return
	}

	read := bufio.NewReader(bytes.NewReader(buf))
	for {
		line, _, err := read.ReadLine()
		if nil != err {
			break
		}

		key, commet, _, _, err := ssh.ParseAuthorizedKey(line)
		log.Println("load public-key commet: ", commet)
		self.allKey = append(self.allKey, key)
	}
}
