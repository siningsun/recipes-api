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
	"fmt"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	_ "recipes-api/docs"
	"recipes-api/handler"
	"recipes-api/middleware"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var RecipeHandler *handler.RecipesHandler
var AuthHandler *handler.AuthHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := redisClient.Ping().Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(redisClient.Ping())
	RecipeHandler = handler.NewRecipesHandler(ctx,
		client.Database("demo").Collection("recipes"),
		redisClient)
	AuthHandler = handler.NewAuthHandler(ctx,
		client.Database("demo").Collection("users"))
}

func main() {
	router := gin.Default()
	store, _ := redisStore.NewStore(10, "tcp", "localhost:6379", "", "", []byte("secret"))
	router.Use(sessions.Sessions("recipes_api", store))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/signin", AuthHandler.SignInHandler)
	router.POST("/refresh", AuthHandler.RefreshHandler)
	router.POST("/signout", AuthHandler.SignOutHandler)
	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddleware())

	{
		authorized.POST("/recipes", RecipeHandler.NewRecipeHandler)
		authorized.PUT("/recipes/:id", RecipeHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", RecipeHandler.DeleteRecipeHandler)
		authorized.GET("/recipes/search", RecipeHandler.SearchRecipeHandler)
		authorized.GET("/recipes", RecipeHandler.ListRecipeHandler)
	}
	err := router.Run()
	if err != nil {
		return
	}
}
