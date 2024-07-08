package main

import (
	"database/sql"

	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
)

type Product struct{
	
	Name string `json:"name"`
	Price string `json:"price"`
	Availability bool `json:"availability"`
}

var db *sql.DB

func main() {

	var err error
    db, err = sql.Open("mysql", "root:Piyush@5336@tcp(localhost:3306)/store")
    if err != nil {
        panic(err)
    }
    defer db.Close()

	router := gin.Default();

	router.POST("/store-products", createProducts);

	router.GET("/list-products", getProducts);


	router.Run(":8000")

}

func createProducts(c *gin.Context) {
	var product Product

	if err := c.BindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	insertQuery := "INSERT INTO products(name,price,availability) VALUES (?,?,?)"

	_ , err := db.Exec(insertQuery, product.Name, product.Price,product.Availability)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
		"message":"Success",
	})
}

func getProducts(c *gin.Context){
	selectQuery := "SELECT * FROM products"

	rows, err := db.Query(selectQuery)

	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	defer rows.Close()
	

	var products []Product

    // Iterate through rows
    for rows.Next() {
        var product Product
   
        err := rows.Scan(&product.Name, &product.Price, &product.Availability)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
      
        products = append(products, product)
    }

    if err := rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

 
    c.JSON(http.StatusOK, gin.H{"products": products})

	

}

