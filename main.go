// @title Recipes API
//
// @description This is a sample recipes API. You can find out more about the API at https://github.com/francknouama/recipes-api
//
// @schemes http
// @host localhost:8080
// @BasePath /
// @version 1.0.0
// @contact.name Franck Nouama
// @contact.email franck.nouama@gmail.com
//
// @Accept json
// @Produce json
package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/francknouama/recipes-api/docs"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
	Ingredients  []string  `json:"ingredients,omitempty"`
	Instructions []string  `json:"instructions,omitempty"`
	PublishedAt  time.Time `json:"published_at,omitempty"`
}

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// NewRecipeHandler godoc
// @Summary Create a new recipe
// @Description POST /recipes recipes newRecipe
// @Tags recipes
// @Accept json
// @Produce json
// @Param recipe body Recipe true "Recipe payload"
// @Success 200 {object} Recipe "Successful operation"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Router /recipes [post]
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)
}

// ListRecipesHandler godoc
// @Summary Return all recipes in the repository.
// @Description GET /recipes recipes listRecipes.
// @Tags recipes
// @Produce json
// @Success 200 {array} Recipe
// @Router /recipes [get]
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

// UpdateRecipeHandler godoc
// @Summary Update an existing recipe
// @Description PUT /recipes/{id} recipes updateRecipe
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200 {object} Recipe "Successful operation"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 404 {object} ErrorResponse "Invalid recipe ID"
// @Router /recipes/{id} [put]
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipe.ID = id
	recipe.PublishedAt = time.Now()
	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
}

// DeleteRecipeHandler godoc
// @Summary Delete an existing recipe
// @Description DELETE /recipes/{id} recipes deleteRecipe
// @Tags recipes
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200 {object} Recipe "Successful operation"
// @Failure 404 {object} ErrorResponse "Invalid recipe ID"
// @Router /recipes/{id} [delete]
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Recipe not found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted",
	})
}

// SearchRecipesHandler godoc
// @Summary Return recipes matching our search criteria
// @Description GET /recipes/search recipes searchRecipes.
// @Tags recipes
// @Produce json
// @Success 200 {array} Recipe "Operation successfully ran"
// @Router /recipes/search [get]
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}

	c.JSON(http.StatusOK, listOfRecipes)
}

func main() {

	docs.SwaggerInfo.Title = "Recipes API"

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)

	router.Run()
}
