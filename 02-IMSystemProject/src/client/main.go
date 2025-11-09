package main

import "fmt"

func main() {
	client := InitClient("127.0.0.1", 8000)
	if client == nil {
		fmt.Println("--- 服务器连接失败 ---")
		return
	}

	fmt.Println("--- 服务器连接成功 ---")
	select {}
}
