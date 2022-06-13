package main

import "net"

/**
登录的用户的结构体
**/
type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn //连接
	server *Server  //所属server
}

//创建一个用户
func NewUser(conn net.Conn, server *Server) *User {
	userAdd := conn.RemoteAddr().String()
	user := &User{
		Name:   userAdd,
		Addr:   userAdd,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	//启动监听当前user channel消息的goroutine
	go user.ListenMessage()
	return user
}

//监听当前User channel的方法，一旦有消息，就直接发送给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}

//用户上线
func (this *User) Online() {

	//用户上线，将用户加入onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//广播当前用户上线消息
	this.server.BroadCast(this, "已上线")
}

//用户下线
func (this *User) Offline() {
	//用户下线，将用户从onlineMap中去除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//广播当前用户下线消息
	this.server.BroadCast(this, "已下线")
}

//用户发送消息
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		//处理消息为who的请求
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线……\n"
			this.SendMessage(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else {
		this.server.BroadCast(this, msg)
	}
}

//给指定用户发送消息
func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}
