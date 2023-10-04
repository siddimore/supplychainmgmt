package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"supplychain-service/pkg/models"
)

type MongoDBService struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoDBService(connectionString, dbName string) (*MongoDBService, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)

	return &MongoDBService{
		client: client,
		db:     db,
	}, nil
}

func (m *MongoDBService) Write(product *models.CoffeeProduct) error {
	collection := m.db.Collection("products")
	_, err := collection.InsertOne(context.Background(), product)
	return err
}

func (m *MongoDBService) Read(productID string) (*models.CoffeeProduct, error) {
	collection := m.db.Collection("products")
	var product models.CoffeeProduct
	err := collection.FindOne(context.Background(), bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
