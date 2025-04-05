# ![Product Management API Banner](https://img.shields.io/badge/Product%20Management%20API-v3.0-blueviolet?style=for-the-badge&logo=go)  
![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)  
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)

## Project Description

**Product Management API** is a lightweight RESTful API built with Go to manage product data efficiently. It supports CRUD operations on a product catalog stored in Microsoft SQL Server, with features like rate limiting for security and scalability. Ideal for e-commerce, inventory systems, or backend experimentation.

## Features

- **CRUD Operations**: Create, read, update, and delete products via intuitive endpoints.
- **Rate Limiting**: IP-based limiter (1 req/sec, burst of 5) to prevent abuse.
- **JSON Responses**: Clear, structured responses with status and data.
- **MSSQL Integration**: Secure database storage with parameterized queries.
- **Middleware**: Built with `go-chi` for routing, logging, and rate limiting.
- **Error Handling**: Detailed responses for invalid inputs or errors.

## Requirements

- **Go**: 1.22+
- **Microsoft SQL Server**: Running instance
- **Dependencies** (via `go.mod`):
  - `github.com/denisenkom/go-mssqldb` v0.12.3
  - `github.com/go-chi/chi` v1.5.5
  - `github.com/go-chi/render` v1.0.3
  - `golang.org/x/time` v0.11.0
  - Indirect: `github.com/ajg/form`, `github.com/golang-sql/civil`, `github.com/golang-sql/sqlexp`, `golang.org/x/crypto`

## Setup Instructions

1. **Clone the Repo**  
   ```bash
   git clone <repository-url>
   cd Product-Management-API-v3

2. **Install Dependencies**
   ```bash
   go mod tidy

3. **Set Up SQL Server**
  Install and configure SQL Server. Click on https://medium.com/@analyticscodeexplained/running-microsoft-sql-server-in-docker-a8dfdd246e45.

4. **Create Products Table**
    ```sql
      CREATE TABLE Products (
          ID VARCHAR(255) PRIMARY KEY,
          Title VARCHAR(255),
          Description TEXT,
          Price FLOAT,
          Quantity INT
      )
     ```
5. **go run main.go**
    ```bash
    go run main.go
    ```
  The server will start at http://localhost:8081.

## Usage
### Endpoints:
- `GET /api/products` - List all products
- `GET /api/products/{id}` - Get product by ID
- `POST /api/products` - Add a product (JSON body)
- `PUT /api/products/{id}` - Update a product (JSON body)
- `DELETE /api/products/{id}` - Delete a product

## Testing with Postman:
- Download and install Postman.
- Set the request URL to `http://localhost:8081/api/<endpoint>`.
- For POST/PUT, select the method, go to the "Body" tab, choose "raw" and "JSON", add the body, and send.

### Example POST request body:
    ```json
    {
        "id": "1",
        "title": "Laptop",
        "description": "High-end laptop",
        "price": 999.99,
        "quantity": 10
    }
    ```
- In the "Headers" tab, set Content-Type: application/json.

## Conclusion
The Product Management API combines simplicity, performance, and modern features like rate limiting, making it a great starting point for developers. Whether you're learning Go or building a production-ready system, this project offers a solid base to extend and customize. Test it with Postman, tweak it, and contribute—your improvements are welcome!
