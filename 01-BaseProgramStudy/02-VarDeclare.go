package main

import "fmt"

func main() {
	// 方式一：声明变量不进行初始化，默认值
	var number int
	fmt.Println("Number is", number)

	// 方式二：根据值自行判定变量类型
	var str = "String Value"
	fmt.Println("String is", str)

	// 方法三：省略var关键字，使用:=进行变量类型推导（PS：不支持全局变量声明）
	data := "this is a data without var"
	fmt.Println("data is", data)
}
