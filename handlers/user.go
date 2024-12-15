package handlers

import (
	"echoingsilence.cn/books-management-system/auth"
	"echoingsilence.cn/books-management-system/database"
	"encoding/json"
	"net/http"
)

// RegisterHandler 用户注册
// @Summary 用户注册，创建新用户并返回 JWT Token
// @Description 用户提供用户名、密码等信息，完成注册后返回 JWT Token，用于后续的身份验证
// @Tags auth
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 201 {object} map[string]string "成功注册，返回 token"
// @Failure 400 {string} string "请求格式错误"
// @Failure 409 {string} string "用户名已存在"
// @Failure 500 {string} string "内部服务器错误"
// @Router /register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// 检查用户名是否已存在
	var existingUsername string
	err = database.DB.QueryRow("SELECT username FROM users WHERE username = ?", input.Username).Scan(&existingUsername)
	if err == nil {
		http.Error(w, `{"error":"Username already exists"}`, http.StatusBadRequest)
		return
	}

	// 插入用户数据
	_, err = database.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", input.Username, input.Password)
	if err != nil {
		http.Error(w, `{"error":"Failed to register user"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"User registered successfully"}`))
}

// LoginHandler 用户登录
// @Summary 用户登录，生成并返回 JWT Token
// @Description 提供用户名和密码进行登录，登录成功返回 JWT Token，用于后续的身份验证
// @Tags auth
// @Accept json
// @Produce json
// @Param username body string true "用户名"
// @Param password body string true "密码"
// @Success 200 {object} map[string]string "成功登录，返回 token"
// @Failure 400 {string} string "请求格式错误"
// @Failure 401 {string} string "用户名或密码错误"
// @Failure 500 {string} string "内部服务器错误"
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// 验证用户名和密码
	var storedPassword string
	var role string
	err = database.DB.QueryRow("SELECT password, role FROM users WHERE username = ?", input.Username).Scan(&storedPassword, &role)
	if err != nil {
		http.Error(w, `{"error":"Invalid username or password"}`, http.StatusUnauthorized)
		return
	}

	if storedPassword != input.Password {
		http.Error(w, `{"error":"Invalid username or password"}`, http.StatusUnauthorized)
		return
	}

	// 生成 Token（带上角色信息）
	token, err := auth.GenerateToken(input.Username, role)
	if err != nil {
		http.Error(w, `{"error":"Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	// 返回 Token
	response := map[string]string{
		"token": token,
		"role":  role,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
