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
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
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
	select {}

}
