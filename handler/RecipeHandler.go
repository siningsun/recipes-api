package handler

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"net/http"
	"recipes-api/model"
	"time"
)

type RecipesHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewRecipesHandler(ctx context.Context, collection *mongo.
	Collection) *RecipesHandler {
	return &RecipesHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// NewRecipeHandler
// @Summary Create a new recipe
// @Description Create a new recipe with the provided details
// @Accept json
// @Produce json
// @Param recipe body Recipe true "Recipe details"
// @Success 200 {object} Recipe
// @Failure 400 {object} Response "Error response"
// @Router /recipes [post]
func (h *RecipesHandler) NewRecipeHandler(c *gin.Context) {
	var recipe model.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	// Save to MongoDB
	_, err := h.collection.InsertOne(h.ctx, recipe)
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
func (h *RecipesHandler) SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")
	var results []model.Recipe
	if tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tag query parameter is required"})
		return
	}
	filter := bson.M{"tags": tag}
	cursor, err := h.collection.Find(h.ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search recipes"})
		return
	}
	for cursor.Next(h.ctx) {
		var recipe model.Recipe
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
func (h *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}
	_, err = h.collection.DeleteOne(h.ctx, bson.M{"_id": objectID})
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
func (h *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe model.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}
	_, err = h.collection.UpdateOne(h.ctx, bson.M{
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
func (h *RecipesHandler) ListRecipeHandler(c *gin.Context) {
	cursor, err := h.collection.Find(h.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve recipes"})
		return
	}
	var recipes []model.Recipe
	for cursor.Next(h.ctx) {
		var recipe model.Recipe
		if err := cursor.Decode(&recipe); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode recipe"})
			return
		}
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, recipes)
}
