package routes

import (
	"echoingsilence.cn/books-management-system/handlers"
	"echoingsilence.cn/books-management-system/middleware"
	"log"
	"net/http"
)

// NewRouter 初始化路由
func NewRouter() *http.ServeMux {
	r := http.NewServeMux()

	// 包装路由统一支持 CORS
	withCORS := func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			middleware.CORS(http.HandlerFunc(handler)).ServeHTTP(w, r)
		}
	}

	// 公共接口（无需认证）
	r.HandleFunc("/register", withCORS(handlers.RegisterHandler)) // 用户注册
	r.HandleFunc("/login", withCORS(handlers.LoginHandler))       // 用户登录
	r.HandleFunc("/books", withCORS(handlers.GetBooksHandler))    // 公共接口获取图书列表
	r.HandleFunc("/book/", withCORS(handlers.GetBookHandler))
	// 需要认证的接口
	r.HandleFunc("/add-book",
		withCORS(middleware.TokenAuthMiddleware(
			middleware.AdminAuthMiddleware(handlers.AddBookHandler),
		)),
	) // 管理员添加图书

	r.HandleFunc("/delete-book/",
		withCORS(middleware.TokenAuthMiddleware(
			middleware.AdminAuthMiddleware(handlers.DeleteBookHandler),
		)),
	) // 管理员删除图书
	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
	return r
}
