package main

import (
	"efficient-sshd/serve"
	"flag"
	"fmt"
	"github.com/kardianos/service"
	"log"
	"os"
	"os/exec"
	"strings"
)

type EfficientSshdSvr struct {
	Config *serve.Config
	cmd    *exec.Cmd
}

func (self *EfficientSshdSvr) Start(s service.Service) error {
	// 获取当前目录
	pwd, _ := os.Getwd()
	if "" != self.Config.Home {
		pwd = self.Config.Home
	}

	// 构建启动新进程
	exe, _ := exec.LookPath(os.Args[0])
	exe = strings.ReplaceAll(exe, "\\", "/")
	if false == strings.HasSuffix(strings.ToLower(exe), ".exe") {
		exe += ".exe"
	}

	// 添加参数
	exe += " shell"
	exe += fmt.Sprintf(" --port %d", self.Config.Port)
	exe += fmt.Sprintf(" --home \"%s\"", pwd)
	if "" != self.Config.Passwd {
		exe += fmt.Sprintf(" --passwd %s", self.Config.Passwd)
	}

	log.Println("start exe: ", exe)
	// 启动子进程
	self.cmd = exec.Command("C:\\msys64\\msys2_shell.cmd", "-msys", "-c", exe)
	self.cmd.Env = os.Environ()
	// 启动子进程
	return self.cmd.Start()
}

func (self *EfficientSshdSvr) Stop(s service.Service) error {
	return self.cmd.Process.Kill()
}

func main() {
	// 判断参数个数是否正确
	if len(os.Args) < 2 {
		fmt.Println("use create/run/delete")
		return
	}

	// 读取配置参数
	config := &serve.Config{}

	flag.IntVar(&config.Port, "port", 8022, "监听的端口")
	flag.StringVar(&config.Home, "home", "", "设置当前的主工程目录")
	flag.StringVar(&config.Passwd, "passwd", "", "登录密码")
	flag.CommandLine.Parse(os.Args[2:])

	// 生成随机密码
	if "" == config.Passwd {
		config.Passwd = serve.RandomPassword()
	}

	config.LoadPublicKey()
	// 打印信息
	config.Dump()

	// 配置服务
	svr, err := service.New(&EfficientSshdSvr{
		Config: config,
	}, &service.Config{
		Name:        "Efficient-Sshd",
		Description: "一个简单高效的服务",
		Arguments:   append([]string{"run"}, os.Args[2:]...),
	})
	if nil != err {
		log.Panic(err)
	}

	switch os.Args[1] {
	case "create":
		err := svr.Install()
		if nil != err {
			log.Panic(err)
		}

		log.Println("安装成功")
	case "delete":
		err := svr.Uninstall()
		if nil != err {
			log.Panic(err)
		}

		log.Println("删除成功")
	case "run":
		err := svr.Run()
		if nil != err {
			log.Panic(err)
		}
	case "shell":
		ShellMain(config)
	default:
		log.Panic("参数错误")
	}
}

func ShellMain(config *serve.Config) {
	// 启动服务
	serve.ServerMain(config)
}
