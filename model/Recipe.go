package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name"`
	Ingredients  []string           `json:"ingredients"`
	Instructions []string           `json:"instructions"`
	Tags         []string           `json:"tags"`
	PublishedAt  time.Time          `json:"publishedAt"`
}
