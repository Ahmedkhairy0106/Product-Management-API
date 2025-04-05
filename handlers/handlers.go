package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"Product-Management-API-v3/db"

)

type Product struct {
    ID          string  `json:"id"`
    Title       string  `json:"title"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Quantity    int     `json:"quantity"`
}

func AddProduct(w http.ResponseWriter, r *http.Request) {
    var p Product

    if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    stmt, err := db.DB.Prepare(`
        INSERT INTO Products
        (ID, Title, Description, Price, Quantity)
        VALUES (@ID, @Title, @Description, @Price, @Quantity)
    `)
    if err != nil {
        http.Error(w, "Failed to prepare ADD statement", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()
    
    _, err = stmt.Exec(
        sql.Named("ID", p.ID),
        sql.Named("Title", p.Title),
        sql.Named("Description", p.Description),
        sql.Named("Price", p.Price),
        sql.Named("Quantity", p.Quantity),
    )
    
    if err != nil {
        if strings.Contains(err.Error(), "Violation of PRIMARY KEY constraint") {
            http.Error(w, "ID already exists", http.StatusBadRequest)
            return
        }
    
        http.Error(w, "Failed to insert product", http.StatusInternalServerError)
        return
    }

    response := map[string]string{
        "message": "Product added successfully",
        "status":  "success",
    }

    render.JSON(w, r, response)
}

func GetAllProduct(w http.ResponseWriter, r *http.Request) {
    stmt, err := db.DB.Prepare(`
        SELECT ID, Title, Description, Price, Quantity FROM Products
    `)
    if err != nil {
        http.Error(w, "Failed to prepare Get all statement", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    rows, err := stmt.Query()
    if err != nil {
        http.Error(w, "Failed to query database", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var products []Product

    for rows.Next() {
        var product Product
        err := rows.Scan(&product.ID, &product.Title, &product.Description, &product.Price, &product.Quantity)
        if err != nil {
            http.Error(w, "Failed to scan row", http.StatusInternalServerError)
            return
        }
        products = append(products, product)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, "Error occurred during row iteration", http.StatusInternalServerError)
        return
    }

    if len(products) == 0 {
        response := map[string]interface{}{
            "message": "No products found",
            "status":  "success",
            "data":    products,
        }
        w.WriteHeader(http.StatusNotFound)
        render.JSON(w, r, response)
        return
    }

    response := map[string]interface{}{
        "message": "All products retrieved successfully",
        "status":  "success",
        "data":    products,
    }

    render.JSON(w, r, response)
}

func GetProductByID(w http.ResponseWriter, r *http.Request) {

    productID := chi.URLParam(r, "id")

    if productID == "" {
        http.Error(w, "Missing product ID", http.StatusBadRequest)
        return
    }

    _, eror := strconv.Atoi(productID)
    if eror != nil {
        http.Error(w, "Invalid product ID: must be an integer", http.StatusBadRequest)
        return
    }
    
    var product Product

    err := db.DB.QueryRow(`
        SELECT ID, Title, Description, Price, Quantity 
        FROM Products 
        WHERE ID = @ID
    `, sql.Named("ID", productID)).Scan(&product.ID, &product.Title, &product.Description, &product.Price, &product.Quantity)

    if err == sql.ErrNoRows {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "Failed to retrieve product: "+err.Error(), http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "message": "The product retrieved successfully",
        "status":  "success",
        "data":    product,
    }

    render.JSON(w, r, response)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
    productID := chi.URLParam(r, "id")

   if productID == "" {
        http.Error(w, "Missing product ID", http.StatusBadRequest)
        return
    }

    _, eror := strconv.Atoi(productID)
    if eror != nil {
        http.Error(w, "Invalid product ID: must be an integer", http.StatusBadRequest)
        return
    }
    
    var product Product

    // Decode the request body into the Product struct
    if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    if product.ID != "" && product.ID != productID {
        http.Error(w, "Product ID in body does not match URL", http.StatusBadRequest)
        return
    }

    stmt, err := db.DB.Prepare(`
        UPDATE Products
        SET Title = @Title, Description = @Description, Price = @Price, Quantity = @Quantity
        WHERE ID = @ID
    `)
    if err != nil {
        http.Error(w, "Failed to prepare update statement", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    result, err := stmt.Exec(
        sql.Named("Title", product.Title),
        sql.Named("Description", product.Description),
        sql.Named("Price", product.Price),
        sql.Named("Quantity", product.Quantity),
        sql.Named("ID", productID),
    )

    if err != nil {
        http.Error(w, "Failed to update product: "+err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, "Failed to check rows affected", http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    }

    response := map[string]string{
        "message": "Product updated successfully",
        "status":  "success",
    }

    render.JSON(w, r, response)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
    productID := chi.URLParam(r, "id")

    if productID == "" {
        http.Error(w, "Missing product ID", http.StatusBadRequest)
        return
    }

    _, eror := strconv.Atoi(productID)
    if eror != nil {
        http.Error(w, "Invalid product ID: must be an integer", http.StatusBadRequest)
        return
    }
    
    stmt, err := db.DB.Prepare(`
        DELETE FROM Products 
        WHERE ID = @ID
    `)
    if err != nil {
        http.Error(w, "Failed to prepare Delete statement", http.StatusInternalServerError)
        return
    }
    defer stmt.Close()

    result, err := stmt.Exec(sql.Named("ID", productID))
    if err != nil {
        http.Error(w, "Failed to delete product: "+err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, "Failed to check rows affected", http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "Product not found", http.StatusNotFound)
        return
    }

    response := map[string]string{
        "message": "Product deleted successfully",
        "status":  "success",
    }

    render.JSON(w, r, response)
}