package main

import (
	"fmt"
	"net"
	"strconv"
)

type Server struct {
	IP   string // IP 地址
	Port int    // Port 端口号
}

// InitServer Server初始化函数
func InitServer(ip string, port int) *Server {
	// 初始化Server实例对象
	server := &Server{
		IP:   ip,
		Port: port,
	}

	// 返回一个执行当前Server实例对象的指针
	return server
}

// InitConnectHandler 连接建立初始化时调用当前方法进行处理
func (server *Server) InitConnectHandler(connection net.Conn) {
	fmt.Println("Connect init finish：", connection.RemoteAddr().String())
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
