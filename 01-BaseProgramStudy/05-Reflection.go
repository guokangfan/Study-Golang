package main

import "fmt"

// pair 变量类型（type value 对），反射能通过变量获取变量类型

func main() {
	// data = type（String） + Value（hello world）
	var data string = "hello world"
	var allType interface{}
	allType = data

	fmt.Println(allType.(string))
}
