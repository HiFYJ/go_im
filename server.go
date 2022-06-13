package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//	在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//	消息广播的channel
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

/**
启动服务器接口
**/
func (this *Server) Start() {
	//创建一个server的接口
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("listener init error:", err)
		return
	}
	defer listener.Close()

	/*启动监听*/
	go this.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		go this.Handler(conn)
	}

}
func (this *Server) Handler(conn net.Conn) {
	//fmt.Println("连接建立成功")

	user := NewUser(conn, this)
	/*
		//用户上线，将用户加入onlineMap中
		this.mapLock.Lock()
		this.OnlineMap[user.Name] = user
		this.mapLock.Unlock()

		//广播当前用户上线消息
		this.BroadCast(user, "已上线")*/
	user.Online()

	//监听用户是否活跃的channel
	isLive := make(chan bool)

	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 40096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				//this.BroadCast(user, "下线了")
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read error")
				return
			}

			//	提取用户发送的信息(去除最后的'\n')
			msg := string(buf[:n-1])
			//广播消息
			//this.BroadCast(user, msg)
			user.DoMessage(msg)

			//	用户任意消息表示活跃
			isLive <- true

		}
	}()

	for {
		select {
		case <-isLive:
		//	表示当前用户是活跃的，要重置定时器
		//不做任何事情，为了激活select，更新下面的定时器

		case <-time.After(time.Minute * 10):
			//	已经超时，将当前的User强制的关闭、
			user.SendMessage("超时被强制下线\n")

			close(user.C) //销毁资源

			conn.Close() //关闭连接

			//退出当前handler
			return //使用runtime.Goexit()也可以

		}
	}
}

/*广播消息方法*/
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

/*监听Message广播消息channel的goroutine，一旦有消息就发送给全部在线User*/
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			fmt.Println("发送消息：【" + msg + "】给用户：" + cli.Name)
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}
