package serve

import (
	"efficient-sshd/sshex"
	"fmt"
	"github.com/runletapp/go-console"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
)

type SimpleSshd struct {
	config    *Config
	ptyWidth  int
	ptyHeight int
	console   console.Console
	env       []string
}

func NewSimpleSshd() *SimpleSshd {
	return &SimpleSshd{
		ptyWidth:  80,
		ptyHeight: 24,
	}
}

func (self *SimpleSshd) Serve(conn net.Conn) {
	defer conn.Close()
	// 打印日志
	log.Printf("%s connected\n", conn.RemoteAddr())
	defer log.Printf("%s disconnect\n", conn.RemoteAddr())

	// 创建配置文件
	serverConfig := &ssh.ServerConfig{
		NoClientAuth: true, // 不需要用户认证
	}

	// 加载证书
	signer, err := ssh.ParsePrivateKey([]byte(rsaPrivate))
	if nil != err {
		log.Println("load signer error:", err)
		return
	}

	serverConfig.AddHostKey(signer)

	// 创建一个ssh服务端
	sc, channel, request, err := ssh.NewServerConn(conn, serverConfig)
	if nil != err {
		log.Println("new server conn error: ", err)
		return
	}

	defer sc.Close()

	go self.serveNewChannel(channel)
	go self.serveRequest(request)

	sc.Wait()
}

// 一个新的通道
func (self *SimpleSshd) serveNewChannel(channel <-chan ssh.NewChannel) {
	for {
		ch, ok := <-channel
		if false == ok {
			return
		}

		log.Println("channel type: ", ch.ChannelType())
		go self.serveChannel(ch)
	}
}

func (self *SimpleSshd) serveChannel(channel ssh.NewChannel) {
	// 接收输入
	ch, req, err := channel.Accept()
	if nil != err {
		log.Println("channel accept error: ", err)
		return
	}

	defer ch.Close()

	// 对channel进行处理
	switch channel.ChannelType() {
	case "session":
		self.serveSession(ch, req)
	default:
		fmt.Printf("unknown channel type: %s\n", channel.ChannelType())
	}
}

func (self *SimpleSshd) serveSession(ch ssh.Channel, request <-chan *ssh.Request) {
	// 启动请求
	for {
		req, ok := <-request
		if false == ok {
			return
		}

		switch req.Type {
		case "pty-req":
			self.requestPtyReq(req)
		case "env":
			self.requestEnv(req)
		case "shell":
			go self.startShell(ch, req)
		case "window-change":
			self.requestWindowChange(req)
		default:
			fmt.Println("request info: ", req.Type)
			fmt.Println("want reply: ", req.WantReply)
			fmt.Println("playload: ", req.Payload)
		}
	}

}

// 启动shell
func (self *SimpleSshd) startShell(ch ssh.Channel, req *ssh.Request) {
	defer ch.Close()

	// 创建一个新的console
	pty, err := console.New(self.ptyWidth, self.ptyHeight)
	if nil != err {
		log.Println("create console error:", err)
		return
	}

	// 启动命令
	// 设置环境变量
	env := os.Environ()
	env = append(env, self.getEnvHome())
	env = append(env, self.getEnvPs1())
	env = append(env, self.env...)
	pty.SetENV(env)

	// 启动终端
	err = pty.Start([]string{"/bin/bash"})
	if nil != err {
		log.Println("start pty error: ", err)
		return
	}

	log.Println("start pty success")
	// 回复客户端
	req.Reply(true, nil)

	self.console = pty
	// 交换数据
	go func() {
		defer pty.Close()
		io.Copy(pty, ch)
		self.console = nil
	}()

	io.Copy(ch, pty)
}

// 配置窗口改变
func (self *SimpleSshd) requestWindowChange(req *ssh.Request) {
	//      byte      SSH_MSG_CHANNEL_REQUEST
	//      uint32    recipient channel
	//      string    "window-change"
	//      boolean   FALSE
	//      uint32    terminal width, columns
	//      uint32    terminal height, rows
	//      uint32    terminal width, pixels
	//      uint32    terminal height, pixels
	pty := self.console
	if nil == pty {
		return
	}

	pl := sshex.NewPayload(req.Payload)
	width, e1 := pl.ReadUint32()
	height, e2 := pl.ReadUint32()
	if nil != e1 || nil != e2 {
		log.Println("'window-change' error: ", e1, e2)
		return
	}

	pty.SetSize(int(width), int(height))
	log.Printf("window-change: width=%d, height=%d\n", width, height)

}

// 重定向home目录
func (self *SimpleSshd) getEnvHome() string {
	dir, _ := os.Getwd()
	return fmt.Sprintf("HOME=%s", dir)
}

// 设置终端颜色
func (self *SimpleSshd) getEnvPs1() string {
	// F       B
	//30      40      黑色
	//31      41      红色
	//32      42      绿色
	//33      43      黄色
	//34      44      蓝色
	//35      45      紫色
	//36      46      青蓝色
	//37      47      白色
	return "PS1=\\033[32m\\h:\\W \\u > \\033[0m"
}

func (self *SimpleSshd) serveRequest(request <-chan *ssh.Request) {
	for {
		req, ok := <-request
		if false == ok {
			return
		}

		switch req.Type {
		default:
			fmt.Println("request info: ", req.Type)
			fmt.Println("want reply: ", req.WantReply)
			fmt.Println("playload: ", req.Payload)
		}
	}
}

func (self *SimpleSshd) requestPtyReq(req *ssh.Request) {
	pty, err := sshex.ParsePtyReqPayload(req.Payload)
	if nil != err {
		log.Println("parse pty-req error: ", err)
		return
	}

	self.ptyWidth = int(pty.CharWidth)
	self.ptyHeight = int(pty.CharHeight)
	log.Printf("pty-req: width=%d, height=%d\n", self.ptyWidth, self.ptyHeight)

	req.Reply(true, nil)
}

func (self *SimpleSshd) requestEnv(req *ssh.Request) {
	pl := sshex.NewPayload(req.Payload)

	key := ""
	for {
		str, err := pl.ReadString()
		if nil != err {
			break
		}

		if "" == key {
			key = str
		} else {
			env := key + "=" + str
			key = ""

			self.env = append(self.env, env)
			log.Println("env: ", env)
		}
	}
}
