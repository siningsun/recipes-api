// @title Recipes API
// @version 1.0
// @description This is a sample recipes API. You can find out more about
// @description the API at https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin
// @termsOfService https://github.com/PacktPublishing/Building-Distributed-Applications-in-Gin

// @contact.name Mohamed Labouardy
// @contact.email mohamed@labouardy.com
// @contact.url https://labouardy.com

// @host localhost:8080
// @BasePath /
// @schemes http
package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	_ "recipes-api/docs"
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

func SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")
	var results []Recipe
	for _, recipe := range recipes {
		for _, rTag := range recipe.Tags {
			if strings.EqualFold(rTag, tag) {
				results = append(results, recipe)
				break
			}
		}
	}
	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No recipes found with the specified tag"})
		return
	}
	c.JSON(http.StatusOK, results)
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

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipeHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run()
}
