package main

import (
	"context"
	"crypto/sha256"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	users := map[string]string{
		"admin":   "fCRmh4Q2JRseqkz",
		"packt":   "RE4zfHB35VPtTkbT",
		"fnouama": "L3nSFRcZzNQ67bcc",
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	h := sha256.New()

	for username, password := range users {
		collection.InsertOne(ctx, bson.M{
			"username": username,
			"password": string(h.Sum([]byte(password))),
		})
	}
}
