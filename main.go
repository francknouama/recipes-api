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
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/francknouama/recipes-api/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id" swaggerignore:"true"`
	Name         string             `json:"name,omitempty" bson:"name"`
	Tags         []string           `json:"tags,omitempty" bson:"tags"`
	Ingredients  []string           `json:"ingredients,omitempty" bson:"ingredients"`
	Instructions []string           `json:"instructions,omitempty" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt,omitempty" bson:"publishedAt"`
}

var recipes []Recipe

var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection

func init() {
	// recipes = make([]Recipe, 0)
	// file, _ := os.ReadFile("recipes.json")
	// _ = json.Unmarshal([]byte(file), &recipes)

	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")

	// var listOfRecipes []interface{}
	// for _, recipe := range recipes {
	// listOfRecipes = append(listOfRecipes, recipe)
	// }
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	// insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	// if err != nil {
	// log.Fatal(err)
	// }
	// log.Println("Inserted recipes: ", len(insertManyResult.InsertedIDs))
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

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err = collection.InsertOne(ctx, recipe)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Error while inserting a new recipe",
		})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// ListRecipesHandler godoc
// @Summary Return all recipes in the repository.
// @Description GET /recipes recipes listRecipes.
// @Tags recipes
// @Produce json
// @Success 200 {array} Recipe "Find all the recipes"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /recipes [get]
func ListRecipesHandler(c *gin.Context) {
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
	}
	defer cur.Close(ctx)

	recipes := make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}

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

	objectID, _ := primitive.ObjectIDFromHex(id)
	_, err = collection.UpdateOne(ctx, bson.M{
		"_id": objectID,
	}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: recipe.Name},
		{Key: "instructions", Value: recipe.Instructions},
		{Key: "ingredients", Value: recipe.Ingredients},
		{Key: "tags", Value: recipe.Tags},
	}}})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
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
	objectId, _ := primitive.ObjectIDFromHex(id)

	_, err := collection.DeleteOne(ctx, bson.M{
		"_id": objectId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}
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
// func SearchRecipesHandler(c *gin.Context) {
// 	tag := c.Query("tag")
// 	listOfRecipes := make([]Recipe, 0)

// 	for i := 0; i < len(recipes); i++ {
// 		found := false
// 		for _, t := range recipes[i].Tags {
// 			if strings.EqualFold(t, tag) {
// 				found = true
// 			}
// 		}
// 		if found {
// 			listOfRecipes = append(listOfRecipes, recipes[i])
// 		}
// 	}

// 	c.JSON(http.StatusOK, listOfRecipes)
// }

func main() {

	docs.SwaggerInfo.Title = "Recipes API"

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	// router.GET("/recipes/search", SearchRecipesHandler)

	router.Run()
}
