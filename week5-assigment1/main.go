package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Drink struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Size     string  `json:"size"`
	Ordered  bool    `json:"ordered"`
}

var drinks = []Drink{
	{ID: 1, Name: "Americano", Size: "Medium", Ordered: false},
	{ID: 2, Name: "Latte", Size: "Large", Ordered: false},
	{ID: 3, Name: "Cappuccino", Size: "Medium", Ordered: false},
	{ID: 4, Name: "Mocha", Size: "Small",Ordered: true},
}

func getDrinks(c *gin.Context) {
	nameQuery := c.Query("name")

	filter := []Drink{}
	for _, d := range drinks {
		if nameQuery == "" || strings.Contains(strings.ToLower(d.Name), strings.ToLower(nameQuery)) {
			filter = append(filter, d)
		}
	}

	c.JSON(http.StatusOK, filter)
}

func orderDrink(c *gin.Context) {
	var req struct {
		ID int `json:"id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, d := range drinks {
		if d.ID == req.ID {
			if d.Ordered {
				c.JSON(http.StatusConflict, gin.H{"message": "Drink already ordered"})
				return
			}
			drinks[i].Ordered = true
			c.JSON(http.StatusOK, drinks[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Drink not found"})
}


func main() {
	r := gin.Default()

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "healthy"})
	})

	api := r.Group("/api/v1")
	{
		api.GET("/drinks", getDrinks)
		api.POST("/order", orderDrink)
	}

	r.Run(":8080")
}