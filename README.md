# Personal Expense Tracker API

A backend REST API built using **Go** and the **Beego v2 framework** that helps users manage daily expenses efficiently.  
Users can register, log in, add expenses, filter them, and view spending summaries.  
All data is stored locally using **CSV files** (no database used).

---

##  Project Description

This is a **solo backend project** designed to simulate a real-world expense tracking system.

It focuses on:
- REST API design
- Authentication system
- CRUD operations
- File-based data persistence (CSV)
- Filtering, sorting, and analytics
- Testing and code quality

---

##  Tech Stack

- **Language:** Go (Golang)
- **Framework:** Beego v2
- **Storage:** CSV files
- **Testing:** Go testing package
- **Tools:**
  - Bee CLI
  - Go modules
  - Swagger
  - gofmt
  - go vet
  - Postman

---
## Code Coverage

> ### Total code coverage = 85.2% 
<img width="944" height="270" alt="codeCoverage" src="https://github.com/user-attachments/assets/b1c45a36-4450-498d-8892-94350c3765a4" />



##  Project Features

### Phase 1 — Setup & Authentication
- User registration API
- User login API
- Authentication using `X-User-ID`

### Phase 2 — Expense CRUD
- Create expense
- Get all expenses
- Get single expense
- Update expense
- Delete expense

### Phase 3 — Filtering, Sorting & Summary
- Filter by category
- Filter by date range
- Sort by amount/date
- Expense summary

### Phase 4 — Testing & Code Quality
- Unit testing
- Code coverage
- gofmt formatting
- go vet analysis

---

##  API Endpoints

### Register
```bash
curl -i -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Abcde","email":"w3@example.com","password":"1234586"}'
```

### Login
```bash
curl -i -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"w3@example.com","password":"1234586"}'
```

### Create Expense
```bash
curl -i -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"title":"lunch","amount":3440.50,"category":"Food","note":"Team lunch","expense_date":"2025-06-10"}'
```

### Get Expenses with limit
```bash
curl -i -X GET "http://localhost:8080/api/v1/expenses?limit=10" \
  -H "X-User-ID: 1"
```

### Get All Expenses 
```bash
curl -i -X GET "http://localhost:8080/api/v1/expenses" \
  -H "X-User-ID: 1"
```

### Get Single Expense
```bash
curl -i -X GET http://localhost:8080/api/v1/expenses/1 \
  -H "X-User-ID: 1"
```

### Update Expense
```bash
curl -i -X PUT http://localhost:8080/api/v1/expenses/1 \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{"title":"Dinner","amount":500.75,"category":"Food","note":"Updated note","expense_date":"2025-06-11"}'
```

### Delete Expense
```bash
curl -i -X DELETE http://localhost:8080/api/v1/expenses/1 \
  -H "X-User-ID: 1"
```

### Filter & Sort
```bash
curl -i "http://localhost:8080/api/v1/expenses?category=Food&date_from=2025-06-01&date_to=2025-06-30&sort_by=amount&sort_order=desc" \
  -H "X-User-ID: 1"
```
### Sort By Date
```bash
curl -i "http://localhost:8080/api/v1/expenses?sort_by=expense_date&sort_order=asc" \
  -H "X-User-ID: 1"
```
### Summary
```bash
curl -i "http://localhost:8080/api/v1/expenses/summary?date_from=2025-06-01&date_to=2025-06-30" \
  -H "X-User-ID: 1"
```

```bash
curl -i "http://localhost:8080/api/v1/expenses/summary" \
  -H "X-User-ID: 1"
```

---

### Some Failure Cases

```bash
curl -i -X GET http://localhost:8080/api/v1/expenses
```

```bash
curl -i "http://localhost:8080/api/v1/expenses?sort_by=title" -H "X-User-ID: 1"
```


```bash
curl -i "http://localhost:8080/api/v1/expenses?date_from=06-01-2025" -H "X-User-ID: 1"
```

```bash
curl -i "http://localhost:8080/api/v1/expenses/summary"
```


---

##  Setup & Run Instructions

This section explains how to set up and run the project step by step on a local machine.

---

### 1. Clone the Repository

First, download the project from the Git repository:

```bash
git clone https://github.com/Parisa-Reza/Expense-Tracker-API-Golang.git
```

Then navigate into the project folder:

```bash
cd Expense-Tracker-API-Golang
```

---


Run the following command for app.conf file:

```bash
cp conf/app.conf.example conf/app.conf
```
---
### 2. Install Go Dependencies

Before running the project, make sure all required Go modules are installed.

```bash
go mod tidy
```

This command will:
- Download all required dependencies
- Sync `go.mod` and `go.sum` files

---

### 3. Install Bee CLI (Beego Tool)

This project uses **Bee CLI** to run and manage the Beego server.

Install it globally:

```bash
go install github.com/beego/bee/v2@latest
```

Verify installation:

```bash
bee version
```

If it shows a version number, the installation is successful.

---

### 4. Run the Application

Start the development server using:

```bash
bee run
```

This will:
- Compile the Go project
- Start the HTTP server
- Automatically restart the server when code changes

---

### 5. Access the Application

Once the server is running, you can access:

- API Base URL:
```bash
http://localhost:8080
```

- Swagger Documentation:
```bash
http://localhost:8080/swagger/
```

---

##  Testing 



###  Run Unit Tests

Execute all test cases in the project:

```bash
go test ./...
```

This will:
- Run all package-level tests
- Show pass/fail results
- Verify application logic correctness

---

###  Run Tests with Coverage

Check how much of your code is covered by tests:

```bash
go test ./... -cover
```

For a detailed report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```


---

###  Static Code Analysis (go vet)

Analyze code for potential issues:

```bash
go vet ./...
```

This helps detect:
- Suspicious code patterns
- Unused variables
- Potential runtime bugs
- Incorrect API usage

---

## Lisence

This project is for assesment and learning purpose assigned by W3 Engineers Ltd.
