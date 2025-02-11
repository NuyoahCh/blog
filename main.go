// @Author Gopher
// @Date 2025/2/11 09:19:00
// @Desc 创建main.go并启动项目
package main

import (
	"blog/common/initialize"
)

func main() {
	print("hello")
	initialize.LoadConfig()
	initialize.Mysql()
	initialize.Router()
}
