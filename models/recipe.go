package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id" swaggerignore:"true"`
	Name         string             `json:"name,omitempty" bson:"name"`
	Tags         []string           `json:"tags,omitempty" bson:"tags"`
	Ingredients  []string           `json:"ingredients,omitempty" bson:"ingredients"`
	Instructions []string           `json:"instructions,omitempty" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt,omitempty" bson:"publishedAt"`
}
