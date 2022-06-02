package main

import "net"

/**
登录的用户的结构体
**/
type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn //连接
}

//创建一个用户
func NewUser(conn net.Conn) *User {
	userAdd := conn.RemoteAddr().String()
	user := &User{
		Name: userAdd,
		Addr: userAdd,
		C:    make(chan string),
		conn: conn,
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
