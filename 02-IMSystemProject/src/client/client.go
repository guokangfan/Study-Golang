package main

import (
	"fmt"
	"net"
	"strconv"
)

type Client struct {
	ServerIP   string
	ServerPort int
	ServerName string
	connection net.Conn
}

// InitClient 初始化Client
func InitClient(serverIP string, serverPort int) *Client {
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
	}

	conn, err := net.Dial("tcp", serverIP+":"+strconv.Itoa(serverPort))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	client.connection = conn
	return client
}
