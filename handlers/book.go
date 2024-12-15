package handlers

import (
	"echoingsilence.cn/books-management-system/database"
	"echoingsilence.cn/books-management-system/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// GetBooksHandler 获取图书列表
// @Summary 获取所有图书
// @Description 获取系统中所有的图书，读者可以访问
// @Tags books
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.BookInfo "Success"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 403 {string} string "Forbidden"  // 403 Forbidden 错误
// @Router /books [get]
func GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := database.DB.Query("SELECT id, title, author, publisher, price FROM books")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []model.BookInfo
	for rows.Next() {
		var book model.BookInfo
		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.Publisher, &book.Price)
		//fmt.Println(book.Id, book.Title, book.Author, "haha")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error reading database rows", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// AddBookHandler 添加图书
// @Summary 添加图书（管理员专用）
// @Description 只有管理员可以添加图书
// @Tags books
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 201 {string} string "Book added successfully"
// @Failure 500 {string} string "Internal Server Error"
// @Failure 403 {string} string "Forbidden"  // 403 Forbidden 错误
// @Router /add-book [post]
func AddBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println(r.Method)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Title     string `json:"title"`
		Author    string `json:"author"`
		Publisher string `json:"publisher"`
		Price     string `json:"price"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	_, err = database.DB.Exec("INSERT INTO books (title, author, publisher, price) VALUES (?, ?, ?, ?)", input.Title, input.Author, input.Publisher, input.Price)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to add book", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message":"Book added successfully"}`))
}

// GetBookHandler 获取单个图书的信息
// @Summary 获取单个图书的详细信息
// @Description 通过图书 ID 获取特定图书的详细信息
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "图书ID"
//
//	@Success 200 {object} model.BookInfo "成功获取图书信息"
//
// @Failure 404 {string} string "图书未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Failure 403 {string} string "权限不足"
// @Router /book/{id} [get]
func GetBookHandler(w http.ResponseWriter, r *http.Request) {
	// 从 URL 参数中获取图书 ID
	idStr := r.URL.Path[len("/book/"):]
	// 将 ID 转换为整数
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid book ID"}`, http.StatusBadRequest)
		return
	}

	// 查询数据库获取图书信息
	var book model.BookInfo
	err = database.DB.QueryRow("SELECT id, title, author, publisher, price FROM books WHERE id = ?", id).Scan(&book.Id, &book.Title, &book.Author, &book.Publisher, &book.Price)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, `{"error":"Book not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"Failed to fetch book"}`, http.StatusInternalServerError)
		}
		return
	}

	// 返回图书信息
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// DeleteBookHandler 删除图书
// @Summary 删除图书
// @Description 根据图书 ID 删除图书
// @Tags books
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "图书ID"
// @Success 200 {string} string "成功删除图书"
// @Failure 400 {string} string "请求格式错误"
// @Failure 404 {string} string "图书未找到"
// @Failure 500 {string} string "内部服务器错误"
// @Router /delete-book/{id} [delete]
func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	// 从 URL 参数中获取图书 ID
	idStr := r.URL.Path[len("/delete-book/"):]
	// 将 ID 转换为整数
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println(err)
		http.Error(w, `{"error":"Invalid book ID"}`, http.StatusBadRequest)
		return
	}

	// 删除图书
	result, err := database.DB.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		http.Error(w, `{"error":"Failed to delete book"}`, http.StatusInternalServerError)
		return
	}

	// 检查是否有记录被删除
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, `{"error":"Failed to check affected rows"}`, http.StatusInternalServerError)
		return
	}

	// 如果没有记录被删除，表示图书未找到
	if rowsAffected == 0 {
		http.Error(w, `{"error":"Book not found"}`, http.StatusNotFound)
		return
	}

	// 成功删除图书
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Book deleted successfully"}`))
}
