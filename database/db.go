package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

var DB *sql.DB

const (
	username = "root"
	password = "Xuchao20051103~~"
	host     = "localhost"
	port     = "3306"
	dbname   = "book_management_system"
)

// 初始化数据库连接
func InitDatabase() {
	var err error
	path := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", dbname, "?charset=utf8"}, "")
	DB, err = sql.Open("mysql", path)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 测试连接
	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connection established")
}
