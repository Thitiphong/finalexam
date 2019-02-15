package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/Thitiphong/finalexam/database"

	"github.com/Thitiphong/finalexam/model"
	"github.com/gin-gonic/gin"
)

func createCustomerHandler(c *gin.Context) {
	var t model.Customer

	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t.Status = "active"

	row := database.InsertCustomer(t.Name, t.Email, t.Status)

	var id int
	if err := row.Scan(&id); err != nil {
		log.Println("cannot scan id", err)
		return
	}
	t.ID = id
	c.JSON(http.StatusCreated, t)
}

func getCustomerByIDHandler(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	row, err := database.SelectByKeyCustomer(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "not found"})
	}

	t := model.Customer{}
	if err := row.Scan(&t.ID, &t.Name, &t.Email, &t.Status); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "not found"})
	} else {
		c.JSON(http.StatusOK, t)
		return
	}
}
func getCustomerHandler(c *gin.Context) {
	status := c.Query("status")

	s := "select id, name, email, status FROM customers"
	if status != "" {
		s = s + " where status = $1"
	}
	stmt, err := database.Conn().Prepare(s)
	if err != nil {
		log.Fatal("can't prepare query all customers statment", err)
	}
	var rows *sql.Rows
	if status != "" {
		rows, err = stmt.Query(status)
	} else {
		rows, err = stmt.Query()
	}
	if err != nil {
		log.Fatal("can't query all customers", err)
	}

	var cs = []model.Customer{}
	for rows.Next() {
		t := model.Customer{}
		err := rows.Scan(&t.ID, &t.Name, &t.Email, &t.Status)
		if err != nil {
			log.Fatal("can't Scan row into variable", err)
		}
		cs = append(cs, t)
	}
	c.JSON(http.StatusOK, cs)
}
func updateCustomerHandler(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	var t model.Customer
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := database.UpdateCustomer(id, t); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		log.Println(err)

		return
	}
	c.JSON(http.StatusOK, t)
}
func deleteCustomerHandler(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	if _, err := database.DeleteCustomer(id); err != nil {
		log.Fatal("error execute delete ", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}

func authMiddleware(c *gin.Context) {
	log.Println("start middleware")
	authKey := c.GetHeader("Authorization")
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		c.Abort()
		return
	}

	c.Next()
}

func setupRoute() *gin.Engine {
	r := gin.Default()
	r.Use(authMiddleware)
	grp := r.Group("/")

	grp.POST("/customers", createCustomerHandler)
	grp.GET("/customers/:id", getCustomerByIDHandler)
	grp.GET("/customers", getCustomerHandler)
	grp.PUT("/customers/:id", updateCustomerHandler)
	grp.DELETE("/customers/:id", deleteCustomerHandler)
	return r
}

func main() {

	database.Conn()
	database.CreateTable()

	r := setupRoute()

	r.Run(":2019")
}
