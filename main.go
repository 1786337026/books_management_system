package main

import (
	"echoingsilence.cn/books-management-system/database"
	"echoingsilence.cn/books-management-system/routes"
)

// @title 图书管理系统
// @version 1.0
// @description 用户的基本功能 读者和管理员权限的区分 对图书的增删改查
// @contact.name Xuchao
// @contact.email chaoxu051103@gmail.com
// @host localhost
// @BasePath /api/v1

func main() {
	// 初始化数据库
	database.InitDatabase()

	// 配置路由
	//http.HandleFunc("/register", RegisterHandler)
	//http.HandleFunc("/login", LoginHandler)
	routes.NewRouter()
	// 启动 HTTP 服务
}
