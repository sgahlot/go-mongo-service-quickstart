package mongo

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/sgahlot/go-mongo-service-quickstart/pkg/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

const (
	ALL_ROWS = "ALL"
)

type Service interface {
	InsertFruit(req *FruitRequest) FruitResponse
	GetFruits(req *FruitRequest) FruitResponse
	GetFruit(req *FruitRequest) Fruit
	DeleteFruits(req *FruitRequest) FruitResponse
}

type FruitService struct {
	fruit *Fruit
}

func (receiver *FruitService) InsertFruit(fruit *FruitRequest) FruitResponse {
	log.Printf("Inserting Fruit (%+v)\n", fruit)

	collection := GetMongoDbCollection(DEFAULT_COLLECTION)

	dbContext := GetContext()
	// collection.FindOneAndUpdate(dbContext, nil, fruit)

	fruit.Id = uuid.NewString() // generate the ID ourselves

	inserted, err := collection.InsertOne(dbContext, fruit)
	common.CheckErrorWithPanic(err, fmt.Sprintf("error while inserting Fruit. Data: %+v", fruit))

	fruitId := inserted.InsertedID
	log.Printf("Inserted Fruit (id=%s, %+v)\n", fruitId, fruit)

	return FruitResponse{
		Id:      fruitId,
		Message: common.RESPOSNE_SUCCESS,
		Err:     nil,
	}
}

func (receiver *FruitService) DeleteFruits(req *FruitRequest) FruitResponse {
	log.Printf("Deleting Fruit(s) (%+v)\n", req)

	searchQuery := req.GetDbSearchQuery()
	collection := GetMongoDbCollection(DEFAULT_COLLECTION)

	deletedData, err := collection.DeleteMany(GetContext(), searchQuery)
	common.CheckErrorWithPanic(err, fmt.Sprintf("error while deleting Fruit (%+v)", req))

	log.Printf("Deleted Fruit(s) (%+v)\n", deletedData)
	return FruitResponse{
		Message: fmt.Sprintf(common.RESPOSNE_SUCCESS+" in deleting %d fruits", deletedData.DeletedCount),
		Err:     nil,
	}
}

func (receiver *FruitService) GetFruit(req *FruitRequest) Fruit {
	fruitResponse := receiver.GetFruits(req)

	var fruit Fruit
	if fruitResponse.Fruits != nil && len(fruitResponse.Fruits) > 0 {
		fruit = fruitResponse.Fruits[0]
	}

	return fruit
}

func (receiver *FruitService) GetFruits(req *FruitRequest) FruitResponse {
	log.Printf("Retrieving Fruits (+%v)\n", req)

	dbContext := GetContext()

	cursor := receiver.searchDb(req, dbContext)

	defer cursor.Close(dbContext)
	var fruits []Fruit

	for cursor.Next(GetContext()) {
		fruit := Fruit{}
		err := cursor.Decode(&fruit)
		common.CheckErrorWithPanic(err, "error while trying to decode Fruit")
		fruits = append(fruits, fruit)
	}

	var message = fmt.Sprintf("Found %d fruits", len(fruits))

	if len(fruits) == 0 {
		// Create default "no-result" response
		message = "No fruits found for given query"
		fruits = nil
	}

	response := FruitResponse{
		Message: message,
		Fruits:  fruits,
	}

	return response
}

func (receiver *FruitService) searchDb(req *FruitRequest, ctx context.Context) *mongo.Cursor {
	collection := GetMongoDbCollection(DEFAULT_COLLECTION)

	var bsonMap interface{}
	if req.Name != ALL_ROWS {
		query := req.GetDbSearchQuery()
		bsonMap = bson.M{"$or": query}
	} else {
		bsonMap = bson.M{}
	}

	cursor, err := collection.Find(ctx, bsonMap)
	common.CheckErrorWithPanic(err, fmt.Sprintf("error while retrieving Fruits (+%v)", req))

	return cursor
}
