package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

type Client struct {
	ServerIP   string
	ServerPort int
	Name       string
	connection net.Conn
	flag       int // 当前客户端模式
}

// InitClient 初始化Client
func InitClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", serverIP+":"+strconv.Itoa(serverPort))
	if err != nil {
		fmt.Println("net.Dial error: ", err)
		return nil
	}

	client.connection = conn
	return client
}

func (client *Client) Menu() bool {
	var flag int

	fmt.Println("1.群聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("4.获取当前在线用户列表")
	fmt.Println("0.退出")

	// 等待用户输入
	fmt.Scanln(&flag)

	if flag >= 1 && flag <= 4 {
		client.flag = flag
		return true
	} else {
		fmt.Println("--- 请输入合法的菜单选项 ---")
		return false
	}
}

func (client *Client) updateUserName() bool {
	fmt.Println("--- 请输入用户名 ---")
	fmt.Scanln(&client.Name)

	requestMessage := ":cn " + client.Name

	// 向连接发送数据
	_, err := client.connection.Write([]byte(requestMessage))
	if err != nil {
		fmt.Println("write error: ", err)
		return false
	}
	return true
}

// 获取在线用户列表
func (client *Client) listOnlineUser() bool {
	requestMessage := ":ul "
	_, err := client.connection.Write([]byte(requestMessage))
	if err != nil {
		fmt.Println("write error: ", err)
		return false
	}
	return true
}

// 群聊模式
func (client *Client) publishSendMessage() bool {
	var chatMessage string

	fmt.Println("--- 请输入聊天内容，输入【exit】退出 ---")
	fmt.Scanln(&chatMessage)

	for chatMessage != "exit" {
		// 发送给服务端
		_, err := client.connection.Write([]byte(chatMessage))
		if err != nil {
			fmt.Println("write error: ", err)
			return false
		}

		chatMessage = ""
		fmt.Println("--- 请输入聊天内容，输入【exit】退出 ---")
		fmt.Scanln(&chatMessage)
	}

	return true
}

// 私聊模式
func (client *Client) privateSendMessage() bool {
	var userName string
	var chatMessage string

	fmt.Println("--- 请输入聊天用户名，输入【exit】退出 ---")
	fmt.Scanln(&userName)

	if userName != "exit" {
		fmt.Println("--- 请输入聊天内容，输入【exit】退出 ---")
		fmt.Scanln(&chatMessage)

		for chatMessage != "exit" {
			requestMessage := ":to " + userName + " " + chatMessage

			// 发送给服务端
			_, err := client.connection.Write([]byte(requestMessage))
			if err != nil {
				fmt.Println("write error: ", err)
				return false
			}

			chatMessage = ""
			fmt.Println("--- 请输入聊天内容，输入【exit】退出 ---")
			fmt.Scanln(&chatMessage)
		}
	}
	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.Menu() != true {
		}

		// 根据客户端不同模式处理
		switch client.flag {
		case 1:
			// 群聊模式
			// fmt.Println("--- 群聊模式 ---")
			client.publishSendMessage()
		case 2:
			// 	私聊模式
			// fmt.Println("--- 私聊模式 ---")
			client.privateSendMessage()
		case 3:
			// 更新用户名
			// fmt.Println("--- 更新用户名 ---")
			client.updateUserName()
			break
		case 4:
			client.listOnlineUser()
			break
		}
	}
}

// 处理服务端回复消息
func (client *Client) handleResponse() {
	// 一旦client.connection有数据，就直接copy到stdout标准输出上，永久阻塞监听
	io.Copy(os.Stdout, client.connection)

	//for {
	//	var readBuffer = make([]byte, 4096)
	//	client.connection.Read(readBuffer)
	//	fmt.Println(string(readBuffer))
	//}
}

var serverIP string
var serverPort int

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "Server IP")
	flag.IntVar(&serverPort, "port", 8000, "Server Port")
}

// go build -o client client.go
func main() {
	// 解析命令行参数
	flag.Parse()

	client := InitClient(serverIP, serverPort)
	if client == nil {
		fmt.Println("--- 服务器连接失败 ---")
		return
	}

	// 单独开启Go协程，处理服务器回执消息
	go client.handleResponse()

	fmt.Println("--- 服务器连接成功 ---")

	// 启动客户端
	client.Run()
}
