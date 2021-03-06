package main

import (
	"efficient-sshd/serve"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	config := &serve.Config{}

	flag.IntVar(&config.Port, "port", 8022, "监听的端口")
	flag.StringVar(&config.Home, "home", "", "设置当前的主工程目录")
	flag.StringVar(&config.Passwd, "passwd", "", "登录密码")
	flag.Parse()

	// 判断是否在msys下运行
	if runtime.GOOS == "windows" {
		if os.Getenv("MSYSTEM") != "MSYS" {
			msysShell(config)
			return
		}
	}

	// 生成随机密码
	if "" == config.Passwd {
		config.Passwd = serve.RandomPassword()
	}

	config.LoadPublicKey()
	// 打印信息
	config.Dump()

	// 启动服务
	serve.ServerMain(config)
}

// 启动msys-shell进程
func msysShell(config *serve.Config) {
	// 获取当前目录
	pwd, _ := os.Getwd()
	if "" != config.Home {
		pwd = config.Home
	}

	// 构建启动新进程
	exe, _ := exec.LookPath(os.Args[0])
	exe = strings.ReplaceAll(exe, "\\", "/")
	if false == strings.HasSuffix(strings.ToLower(exe), ".exe") {
		exe += ".exe"
	}

	// 添加参数
	exe += fmt.Sprintf(" --port %d", config.Port)
	exe += fmt.Sprintf(" --home \"%s\"", pwd)
	if "" != config.Passwd {
		exe += fmt.Sprintf(" --passwd %s", config.Passwd)
	}

	log.Println("start exe: ", exe)
	// 启动子进程
	cmd := exec.Command("C:\\msys64\\msys2_shell.cmd", "-msys", "-c", exe)
	cmd.Env = os.Environ()
	// 启动子进程
	cmd.Start()
}
