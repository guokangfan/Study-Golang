package main

import (
	"fmt"
)

// 定义函数
func function(param string) string {
	fmt.Println("Function called with param", param)
	return param
}

// Golang中的函数支持返回多个值
func multiBackFunc(param string) (string, string) {
	fmt.Println("param is", param)
	return "first return param", "second return param"
}

func main() {
	data := function("this is a data without var")
	fmt.Println("data is", data)

	returnData1, returnData2 := multiBackFunc("this is a data for multi back func")
	fmt.Println("returnData1 is", returnData1)
	fmt.Println("returnData2 is", returnData2)

	// Function.FunctionTest("hello world!")
}
