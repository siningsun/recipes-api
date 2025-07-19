package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"net/http"
	"time"
)

type Recipe struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Tags         []string `json:"tags"`
	PublishedAt  string   `json:"publishedAt"`
}

var recipes []Recipe

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now().Format(time.RFC3339)
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}
func init() {
	recipes = make([]Recipe, 0)
}
func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.Run()
}
