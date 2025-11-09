package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	IP             string           // IP 地址
	Port           int              // Port 端口号
	OnlineUserMap  map[string]*User // 在线用户列表数据
	mapLock        sync.RWMutex     // Map集合锁
	MessageChannel chan string      // 消息广播通道
}

// InitServer Server初始化函数
func InitServer(ip string, port int) *Server {
	// 初始化Server实例对象
	server := &Server{
		IP:             ip,
		Port:           port,
		OnlineUserMap:  make(map[string]*User),
		MessageChannel: make(chan string),
	}

	// 返回一个执行当前Server实例对象的指针
	return server
}

// Start 启动服务器对象
func (server *Server) Start() {
	// 使用net包中的TCP端口监听
	listener, err := net.Listen("tcp", server.IP+":"+strconv.Itoa(server.Port))

	if err != nil {
		fmt.Println("net.Listen err: ", err)
		return
	}
	// 延迟释放IP和端口监听
	defer listener.Close()

	// 启动监听Message的Goroutine
	go server.ListenMessage()

	for {
		// 监听端口的连接请求，如果有请求进入会创建一个connection连接并返回
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err: ", err)
			continue
		}

		// 处理连接请求
		go server.InitConnectHandler(connection)
	}
}

// BroadCast 广播消息给所有的在线用户
func (server *Server) BroadCast(user *User, message string) {
	messageInfo := "[" + user.Name + "] " + message

	// 消息管道中写入消息
	server.MessageChannel <- messageInfo
}

// ListenMessage 监听MessageChannel中的消息，启动一个Go协程当有消息的时候，广播给全部的在线User
func (server *Server) ListenMessage() {
	for {
		message := <-server.MessageChannel

		server.mapLock.Lock()
		for _, user := range server.OnlineUserMap {
			user.MessageChannel <- message
		}
		server.mapLock.Unlock()
	}
}

// InitConnectHandler 连接建立初始化时调用当前方法进行处理
func (server *Server) InitConnectHandler(connection net.Conn) {
	fmt.Println("Connect init：", connection.RemoteAddr().String())

	// 当用连接进入时，表示用户上线，需要记录用户的连接信息到OnlineUserMap中
	// 初始化用户并进行在线用户记录
	user := InitUser(connection, server)

	// 给OnlineUserMap加锁
	//server.mapLock.Lock()
	//server.OnlineUserMap[user.Name] = user
	//server.mapLock.Unlock()

	// 广播当前用户上线给其他用户
	// server.BroadCast(user, "已上线")

	// 用户上线
	user.Online()

	// 接收客户端发送的消息
	go handleReceive(connection, server, user)

	// 阻塞当前handler
	for {
		select {
		case <-user.IsLive:
			// 表示当前用户处于活跃状态，重制超时定时器
		case <-time.After(time.Minute * 10):
			// 进入当前作用域，活跃超时，强行对User进行处理
			user.SendMessage("长时间未活跃，连接关闭")

			// 移除用户在线状态
			user.Offline()

			// 资源回收
			close(user.MessageChannel)
			user.connection.Close()

			return
		}
	}
}

// 接收客户端发送的消息，然后进行消息广播
func handleReceive(connection net.Conn, server *Server, user *User) {
	buffer := make([]byte, 4096)

	for {
		// 从Connection中读取数据
		n, err := connection.Read(buffer)
		if n == 0 {
			server.BroadCast(user, "已下线")
			//user.Offline()
		}
		if err != nil && err == io.EOF {
			fmt.Println("Connection read err: ", err)
			return
		}

		// 提取用户消息（去除 \n）
		message := string(buffer[:n])
		// 将获取到的消息广播给所有用户
		// server.BroadCast(user, message)
		user.HandleMessage(message) // 处理消息

		// 更新用户的活跃状态
		user.IsLive <- true
	}
}

// GetOnlineUserByName 根据用户名查询在线用户
func (server *Server) GetOnlineUserByName(userName string) (*User, bool) {
	server.mapLock.RLock()
	defer server.mapLock.RUnlock()

	// 根据用户名查询在线用户
	user, ok := server.OnlineUserMap[userName]
	if ok {
		return user, ok
	} else {
		return nil, false
	}
}
