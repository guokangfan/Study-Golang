package main // 设置当前程序的包名

//import "fmt"
import (
	"fmt"
	"time"
)

// 定义main函数
func main() {
	// golang中的表达式，加分号和不加分号都可以，建议是不加分号
	fmt.Println("Hello World")

	time.Sleep(1 * time.Second)
}
