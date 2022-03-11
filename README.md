# efficient-sshd



> 个人已经习惯了在苹果系统进行下软件开发工作，在调试其它平台时，例如`windows`等常规的操作有点繁琐，因此用当前工具，只为提高工作效率。本工具，只为提高工作效率，请不要用于生产的系统中，以免对生产系统产生严重漏洞。仅供个人学习使用，请不要用于商业目的。



# 生成证书

```bash
ssh-keygen -N "" -f ssh.rsa -t rsa
```



# Windows运行环境

如果要在windows下正常运行，windows必须先安装`msys`，安装教程: [https://www.msys2.org/#installation](https://www.msys2.org/#installation)



# 本工程引用第三方库说明

| 第三方包                        | 用途                    | 备注 |
| ------------------------------- | ----------------------- | ---- |
| golang.org/x/crypto/ssh         | 用于ssh协议             |      |
| github.com/runletapp/go-console | 用于终端操作            |      |
| github.com/kardianos/service    | 用于windows下的服务安装 |      |

