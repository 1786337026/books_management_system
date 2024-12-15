CREATE DATABASE IF NOT EXISTS book_management_system;

USE book_management_system;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role ENUM('reader', 'admin') NOT NULL DEFAULT 'reader',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
DROP TABLE IF EXISTS books;
CREATE TABLE books (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    author VARCHAR(100) NOT NULL,
    publisher VARCHAR(100) NOT NULL,
    price VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
INSERT INTO users (username, password, role)
VALUES ('admin1', 'welcome-admin1', 'admin');
INSERT INTO books (title, author, publisher, price)
VALUES ('生死疲劳', '莫言', '浙江文艺出版社', '66.6元');
