package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"io/ioutil"
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
	// Pre-populate with some sample recipes
	file, _ := ioutil.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.Run()
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	for i, recipe := range recipes {
		if recipe.ID == id {
			recipes = append(recipes[:i], recipes[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var updatedRecipe Recipe
	if err := c.ShouldBindJSON(&updatedRecipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, recipe := range recipes {
		if recipe.ID == id {
			updatedRecipe.ID = id
			updatedRecipe.PublishedAt = recipe.PublishedAt // Keep original published date
			recipes[i] = updatedRecipe
			c.JSON(http.StatusOK, updatedRecipe)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
}

func ListRecipeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}
