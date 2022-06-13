package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	//链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn

	//	返回对象
	return client
}

var ip string
var port int

//通过启动命令后面加-ip 127.0.0.1 -port 18888形式启动
func init() {
	//第一个值：设置变量，第二个值参数名称，第三个值：默认值，第四个值：help时的提示
	flag.StringVar(&ip, "ip", "127.0.0.1", "设置ip地址（默认127.0.0.1）")
	flag.IntVar(&port, "port", 18888, "设置端口（默认18888）")
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag < 4 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>请输入合法范围内的数字<<<<<<")
		return false
	}
}

func (client *Client) Run() {

	for client.flag != 0 {
		for client.menu() != true {

		}
		//	根据flag处理业务
		switch client.flag {
		case 1:
			fmt.Println("进入公聊模式...")
			break
		case 2:
			fmt.Println("私聊模式...")
			break
		case 3:
			fmt.Println("更新用户名...")
			break

		}
	}
}
func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(ip, port)
	if client == nil {
		fmt.Println(">>>>>>>链接服务器失败....")
		return
	}

	fmt.Println(">>>>>>>>>>>链接服务器成功....")

	//启动客户端的业务
	//select {}
	client.Run()
}
