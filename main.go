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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	_ "recipes-api/docs"
	handler "recipes-api/handler"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var RecipeHandler *handler.RecipesHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	RecipeHandler = handler.NewRecipesHandler(ctx, client.Database("demo").Collection("recipes"))
	log.Println("Connected to MongoDB")
}

func main() {
	router := gin.Default()
	router.POST("/recipes", RecipeHandler.NewRecipeHandler)
	router.GET("/recipes", RecipeHandler.ListRecipeHandler)
	router.PUT("/recipes/:id", RecipeHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", RecipeHandler.DeleteRecipeHandler)
	router.GET("/recipes/search", RecipeHandler.SearchRecipeHandler)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run()
}
