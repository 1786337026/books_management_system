package model

// Book 结构体表示图书
type BookInfo struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
	Price     string `json:"price"`
}
