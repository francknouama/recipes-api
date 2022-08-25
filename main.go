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
	"os"

	"github.com/francknouama/recipes-api/docs"
	"github.com/francknouama/recipes-api/handlers"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection

var recipesHandler *handlers.RecipesHandler

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
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	log.Println(status)
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
	// insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	// if err != nil {
	// log.Fatal(err)
	// }
	// log.Println("Inserted recipes: ", len(insertManyResult.InsertedIDs))
}

func main() {

	docs.SwaggerInfo.Title = "Recipes API"

	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	// router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)

	router.Run()
}
