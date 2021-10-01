package mongo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type Fruit struct {
	Id          string `json:"id" bson:"_id"`
	Description string `json:"description"`
	Name        string `json:"name"`
}

type FruitRequest = Fruit

func (fruit *Fruit) String() string {
	return fmt.Sprintf("id=%s, name=%s, description=%s", fruit.Id, fruit.Name, fruit.Description)
}

type FruitResponse struct {
	Id      interface{} `json:"id,omitempty" bson:"_id"`
	Message string      `json:"message,omitempty"`
	Fruits  []Fruit     `json:"fruits,omitempty"` // Do not display if empty
	Err     *error      `json:"error,omitempty"`  // Do not display if empty
}

func (req *FruitRequest) GetDbSearchQuery() []bson.M {
	var query []bson.M

	if req.Id != "" {
		query = append(query, bson.M{"_id": req.Id})
	}
	if req.Name != "" {
		query = append(query, bson.M{"name": req.Name})
	}
	if req.Description != "" {
		query = append(query, bson.M{"description": req.Description})
	}

	return query
}
