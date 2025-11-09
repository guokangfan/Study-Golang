package main

import (
	"net"
)

type User struct {
	Name           string      // 当前用户名
	Addr           string      // 当前用户IP地址
	MessageChannel chan string // 消息管道
	connection     net.Conn    // TCP连接信息
	server         *Server     // 服务实例
}

// InitUser 初始化用户信息
func InitUser(connection net.Conn, server *Server) *User {
	// 从TCP连接中获取用户地址信息
	userAddr := connection.RemoteAddr().String()

	// 初始化一个用户信息
	user := &User{
		Name:           userAddr,
		Addr:           userAddr,
		MessageChannel: make(chan string),
		connection:     connection,
		server:         server,
	}

	// 启动一个Go协程，监听当前用户的Message Channel中是否有消息进入
	go user.ListenMessage()

	// 返回用户信息指针
	return user
}

// ListenMessage 监听管道消息
func (user *User) ListenMessage() {
	for {
		message := <-user.MessageChannel
		// 当监听到消息用户消息队列中有消息之后，把消息数据写入到当前用户的TCP Socket连接
		// 向用户发送消息数据
		user.connection.Write([]byte(message + "\n"))
	}
}

// Online 用户上线
func (user *User) Online() {
	user.server.mapLock.Lock()
	user.server.OnlineUserMap[user.Name] = user
	user.server.mapLock.Unlock()

	user.server.BroadCast(user, "已上线")
}

// Offline 用户下线
func (user *User) Offline() {
	user.server.mapLock.Lock()
	delete(user.server.OnlineUserMap, user.Name)
	user.server.mapLock.Unlock()

	user.server.BroadCast(user, "已下线")
}

// SendMessage 给当前用户发送消息
func (user *User) SendMessage(message string) {
	user.connection.Write([]byte(message))
}

// HandleMessage 处理消息广播发送
func (user *User) HandleMessage(message string) {
	if message == ":user list" {
		// 查询当前在线用户列表
		user.server.mapLock.RLock()
		for _, online := range user.server.OnlineUserMap {
			onlineMessage := "[" + online.Addr + "]" + online.Name + ": 在线\n"
			user.SendMessage(onlineMessage)
		}
		user.server.mapLock.RUnlock()
	} else {
		user.server.BroadCast(user, message)
	}
}
