package main

import "github.com/gin-gonic/gin"

type Recipe struct {
	Name         string   `json:"name"`
	Ingredients  []string `json:"ingredients"`
	Instructions string   `json:"instructions"`
	Tags         []string `json:"tags"`
	PublishedAt  string   `json:"publishedAt"`
}

func main() {
	router := gin.Default()
	router.Run()
}
