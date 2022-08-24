package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Run()
}

type Recipe struct {
	Name         string    `json:"name,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
	Ingredients  []string  `json:"ingredients,omitempty"`
	Instructions []string  `json:"instructions,omitempty"`
	PublishedAt  time.Time `json:"published_at,omitempty"`
}
