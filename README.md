# 📝 TodoList API – GoFiber + MongoDB

Simple REST API untuk manajemen **Todo List** menggunakan [GoFiber](https://gofiber.io/) sebagai web framework, [MongoDB](https://www.mongodb.com/) sebagai database, dan [JWT](https://jwt.io/) untuk autentikasi.

---

## 🚀 Tech Stack

- [Golang](https://go.dev/) – Backend
- [Fiber](https://gofiber.io/) – Web Framework
- [MongoDB](https://www.mongodb.com/) – Database
- [JWT](https://jwt.io/) – Authentication
- [Air](https://github.com/air-verse/air) – Live reload saat development

---

## 📂 Project Structure

.
├── config/ # Konfigurasi database
├── controllers/ # Handler / logic untuk routes
├── middleware/ # Middleware (JWT, role, dll)
├── models/ # Schema data (User, Todo, dll)
├── main.go # Entry point aplikasi
└── go.mod # Module & dependencies

---

## ⚙️ Setup

### 1️⃣ Clone repo

```bash
git clone https://github.com/username/todolist-gofiber.git
cd todolist-gofiber

2️⃣ Install dependencies

go mod tidy

3️⃣ Buat file .env

4️⃣ Run app (development)

air

📌 API Endpoints
🔑 Auth
Method	Endpoint	Description
POST	/register	Register user baru
POST	/login	Login, dapatkan JWT
```
