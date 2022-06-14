package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
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
	fmt.Println("4.在线用户")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag < 5 {
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
			//fmt.Println("进入公聊模式...")
			client.PublicChat()
			break
		case 2:
			//fmt.Println("私聊模式...")
			client.PrivateChat()
			break
		case 3:
			//fmt.Println("更新用户名...")
			client.UpdateName()
			break
		case 4:
			client.SelectUser()
			break

		}
	}
}

//更新用户名功能
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>请输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client.conn.Write error:", err)
		return false
	}
	return true
}

//公聊功能
func (client *Client) PublicChat() {

	//	提示消息
	var msg string
	fmt.Println(">>>>>请输入消息内容，exit退出。")
	fmt.Scanln(&msg)

	for msg != "exit" {
		//	发送给服务器
		//消息不为空则发送
		if len(msg) != 0 {
			sendMsg := msg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn write error:", err)
				break
			}
		}

		//	置空
		msg = ""
		fmt.Println(">>>>>请输入消息内容，exit退出。")
		fmt.Scanln(&msg)
	}
}

//选择在线用户
func (client *Client) SelectUser() {
	msg := "who\n"
	_, err := client.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("client conn error:", err)
		return
	}
}

//私聊功能
func (client *Client) PrivateChat() {
	var remoteName string
	var msg string

	client.SelectUser()
	fmt.Println(">>>>>请输入聊天对象，exit退出。")
	fmt.Scanln(&remoteName)
	for remoteName != "exit" {
		fmt.Println(">>>>>请输入消息内容，exit退出。")
		fmt.Scanln(&msg)
		for msg != "exit" {
			if len(msg) != 0 {
				sendMsg := "to|" + remoteName + "|" + msg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn write error:", err)
					break
				}
			}
			//	置空
			msg = ""
			fmt.Println(">>>>>请输入消息内容，exit退出。")
			fmt.Scanln(&msg)
		}
		client.SelectUser()
		fmt.Println(">>>>>请输入聊天对象，exit退出。")
		fmt.Scanln(&remoteName)
	}
}

//处理server回应的消息，直接显示到标准输出
func (client *Client) DealResponse() {
	//一旦client.conn有数据，就直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.conn)

	//等价于
	/*for{
		buf := make()
		client.conn.Read(buf)
		fmt.Print(buf)
	}*/
}

func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(ip, port)
	if client == nil {
		fmt.Println(">>>>>>>链接服务器失败....")
		return
	}

	//单独开启一个goroutine处理server的回执消息
	go client.DealResponse()
	fmt.Println(">>>>>>>>>>>链接服务器成功....")

	//启动客户端的业务
	//select {}
	client.Run()
}
