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
	"context"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	_ "recipes-api/docs"
)

type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name"`
	Ingredients  []string           `json:"ingredients"`
	Instructions []string           `json:"instructions"`
	Tags         []string           `json:"tags"`
	PublishedAt  string             `json:"publishedAt"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var ctx context.Context
var client *mongo.Client
var err error

// NewRecipeHandler
// @Summary Create a new recipe
// @Description Create a new recipe with the provided details
// @Accept json
// @Produce json
// @Param recipe body Recipe true "Recipe details"
// @Success 200 {object} Recipe
// @Failure 400 {object} Response "Error response"
// @Router /recipes [post]
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now().Format(time.RFC3339)
	// Save to MongoDB
	collection := client.Database("demo").Collection("recipes")
	_, err := collection.InsertOne(ctx, recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe"})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// SearchRecipeHandler
// @Summary Search recipes by tag
// @Description Search for recipes that contain the specified tag
// @Accept json
// @Produce json
// @Param tag query string true "Tag to search for"
// @Success 200 {array} Recipe
// @Failure 404 {object} Response "No recipes found"
// @Router /recipes/search [get]
func SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")
	var results []Recipe
	collection := client.Database("demo").Collection("recipes")
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag query parameter is required"})
		return
	}
	filter := bson.M{"tags": tag}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search recipes"})
		return
	}
	for cursor.Next(ctx) {
		var recipe Recipe
		if err := cursor.Decode(&recipe); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode recipe"})
			return
		}
		results = append(results, recipe)
	}
	if len(results) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No recipes found with the specified tag"})
		return
	}
	c.JSON(http.StatusOK, results)
}

// DeleteRecipeHandler
// @Summary Delete a recipe
// @Description Delete a recipe by its ID
// @Accept json
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200 {object} Response "Success message"
// @Failure 404 {object} Response "Recipe not found"
// @Router /recipes/{id} [delete]
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	collection := client.Database("demo").Collection("recipes")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete recipe"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted successfully"})
}

// UpdateRecipeHandler
// @Summary Update a recipe
// @Description Update a recipe by its ID with the provided details
// @Accept json
// @Produce json
// @Param id path string true "Recipe ID"
// @Param recipe body Recipe true "Updated recipe details"
// @Success 200 {object} Recipe
// @Failure 400 {object} Response "Error response"
// @Failure 404 {object} Response "Recipe not found"
// @Router /recipes/{id} [put]
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	collection := client.Database("demo").Collection("recipes")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}
	_, err = collection.UpdateOne(ctx, bson.M{
		"_id": objectID,
	}, bson.D{{"$set", bson.D{
		{"name", recipe.Name},
		{"instructions", recipe.Instructions},
		{"ingredients", recipe.Ingredients},
		{"tags", recipe.Tags},
	}}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipe"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe updated"})
}

// ListRecipeHandler
// @Summary List all recipes
// @Description Retrieve a list of all recipes
// @Accept json
// @Produce json
// @Success 200 {array} Recipe
// @Router /recipes [get]
func ListRecipeHandler(c *gin.Context) {
	collection := client.Database("demo").Collection("recipes")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve recipes"})
		return
	}
	var recipes []Recipe
	for cursor.Next(ctx) {
		var recipe Recipe
		if err := cursor.Decode(&recipe); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode recipe"})
			return
		}
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, recipes)
}

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
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
