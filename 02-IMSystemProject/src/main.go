package main

func main() {
	// 初始化Server
	server := InitServer("127.0.0.1", 8000)
	// 启动Server，监听指定的IP地址和端口
	server.Start()
}
