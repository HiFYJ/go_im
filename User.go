package main

import (
	"net"
	"strings"
)

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
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//	处理的消息格式 rename|张三
		newName := strings.Split(msg, "|")[1]
		//判断name是否存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMessage("当前用户名已被占用\n")
		} else {
			this.server.mapLock.Lock()
			//删除原有key-value
			delete(this.server.OnlineMap, this.Name)
			//新增
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMessage("已更新用户名:" + this.Name + "\n")
		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		//	处理的消息格式 to|username|content

		//	1、获取对方的用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMessage("消息格式不正确，请使用\"to|username|content\"格式。\n")
			return
		}

		//	2、根据用户名得到对方user对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMessage("该用户不存在\n")
			return
		}

		//	3获取消息内容，通过对方的user对象将消息内容发送出去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMessage("消息不能为空\n")
			return
		}
		remoteUser.SendMessage(this.Name + "向您发送消息：" + content + "\n")
	} else {
		this.server.BroadCast(this, msg)
	}
}

//给指定用户发送消息
func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}
