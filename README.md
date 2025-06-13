
# 🏥 Clinic Portal golang

A full-stack web application for managing patient records with role-based access for **Receptionists** and **Doctors**, built using **Golang**, **Gin**, **PostgreSQL**, and a simple HTML frontend.

---

## 📦 Features

- 🛡️ JWT-based Authentication & Role-based Authorization
- 🩺 Doctors can view and update patient records
- 🧾 Receptionists can perform full CRUD operations on patients
- 🧪 Comprehensive unit and integration tests
- 📄 API documentation via Postman collection
- 🎨 Split frontend UI (HTML/CSS/JS)
- 🧱 Clean modular architecture with proper directory structure

---

## 🛠️ Tech Stack

- **Backend**: Go, Gin, GORM, JWT, Bcrypt
- **Database**: PostgreSQL
- **Frontend**: HTML, CSS, Vanilla JS 
- **Documentation**: Postman (v2.1)

---

## 🗂️ Directory Structure

```
/config         → DB setup and initialization
/controllers    → Business logic (Auth, Patient)
/middleware     → JWT auth middleware
/models         → User and Patient schema
/routes         → Grouped route definitions
/frontend       → HTML/CSS/JS files
```

---

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/jaichiranjeeva/Clinic_portal.git
cd Clinic_portal
```

### 2. Setup PostgreSQL

Create two databases: one for development and one for testing.

```bash
createdb hospital
createdb hospital_test
```

### 3. Configure DB connection in `config/database.go`

Update the DSN line as needed.

### 4. Run the server

```bash
go run main.go
```

Server runs at: [http://localhost:8088](http://localhost:8088)

---

## 🧪 Run Tests

```bash
go test -v
```

---

## 📬 API Documentation

-  Postman

---

## 💡 Extras

- Modular file structure using standard Go practices
- Follows Repository & Controller patterns
- JWT Auth implemented with token validation middleware
- Frontend pages split with localStorage-based JWT handling

---
