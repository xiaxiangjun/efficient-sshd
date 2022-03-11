package main

import (
	"bufio"
	"efficient-sshd/serve"
	"efficient-sshd/system"
	"flag"
	"fmt"
	"github.com/kardianos/service"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

type EfficientSshdSvr struct {
	Config *serve.Config
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
	return system.LaunchProcessWithUser("C:\\msys64\\msys2_shell.cmd", "-msys", "-c", exe)
	//cmd := exec.Command("C:\\msys64\\msys2_shell.cmd", "-msys", "-c", exe)
	//cmd.Env = os.Environ()
	//return cmd.Start()
}

func (self *EfficientSshdSvr) Stop(s service.Service) error {
	c, err := net.Dial("tcp", "127.0.0.1:24")
	if nil != err {
		return err
	}

	c.Write([]byte("exit\r\n\r\n"))
	return nil
}

func initLogger() {
	// 初始化日志
	exe, _ := exec.LookPath(os.Args[0])
	fp, _ := os.OpenFile(exe+".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	log.SetOutput(fp)
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
		// 初始化日志
		initLogger()
		// 启动服务
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
	// 监听退出事件
	go waitForClose()

	// 启动服务
	serve.ServerMain(config)
}

func waitForClose() {
	l, err := net.Listen("tcp", "127.0.0.1:24")
	if nil != err {
		log.Panic(err)
	}

	for {
		c, err := l.Accept()
		if nil != err {
			log.Panic(err)
		}

		go waitOneForClose(c)
	}
}

func waitOneForClose(c net.Conn) {
	defer c.Close()

	line, _, err := bufio.NewReader(c).ReadLine()
	if nil != err {
		return
	}

	if strings.HasPrefix(string(line), "exit") {
		os.Exit(-1)
	}
}
