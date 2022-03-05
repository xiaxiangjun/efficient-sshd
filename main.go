package main

import (
	"efficient-sshd/serve"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	config := &serve.Config{}

	flag.IntVar(&config.Port, "port", 8022, "listen tcp port")
	flag.Parse()

	// 判断是否在msys下运行
	if runtime.GOOS == "windows" {
		if os.Getenv("MSYSTEM") != "MSYS" {
			exe, _ := filepath.Abs(os.Args[0])
			exe = strings.ReplaceAll(exe, "\\", "/")
			// 添加参数
			for i := 1; i < len(os.Args); i++ {
				exe += "\x20" + os.Args[i]
			}

			log.Println("start exe: ", exe)
			// 启动子进程
			cmd := exec.Command("C:\\msys64\\msys2_shell.cmd", "-msys", "-c", exe)
			cmd.Env = os.Environ()
			// 启动子进程
			cmd.Start()
			return
		}
	}

	// 打印信息
	config.Dump()

	// 启动服务
	serve.ServerMain(config)
}
