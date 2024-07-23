package database

import (
	"context"
	"fmt"
	"log"

	"github.com/Thinhhoagn0211/go-parser/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type DbConfig struct {
	Client         *mongo.Client
	Collection     *mongo.Collection
	CollectionName string
	Url            string
	DbName         string
}

func (db *DbConfig) InitMongo(uri, dbName, collectionName string) {
	// Set client options
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	db.Client = client
	db.Url = uri
	db.DbName = dbName
	db.CollectionName = collectionName
}
func (db *DbConfig) Connect() {
	// Check the connection
	err := db.Client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	// Get a handle for your collection
	collection := db.Client.Database(db.DbName).Collection(db.CollectionName)
	db.Collection = collection
}

func (db *DbConfig) GetAllAddresses() ([]string, error) {
	var addresses []string

	// Define a cursor to find all documents in the collection
	cursor, err := db.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor
	for cursor.Next(context.TODO()) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		// Assuming the address is stored as a string
		if address, ok := doc["address"].(string); ok {
			addresses = append(addresses, address)
		}
	}

	// Check for errors encountered during iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (db *DbConfig) Insert(log bson.M) error {
	// Check if a document with the same address already exists
	var result bson.M
	err := db.Collection.FindOne(context.TODO(), bson.M{"address": log["address"]}).Decode(&result)
	if err == nil {
		// If no error, a document with the same address exists
		err = fmt.Errorf("a document with the same address already exists")
		return err
	} else if err != mongo.ErrNoDocuments {
		return err
	}

	// If no document with the same address exists, insert the new document
	insertResult, err := db.Collection.InsertOne(context.TODO(), log)
	if err != nil {
		return err
	}
	fmt.Printf("Inserted log with ID: %v\n", insertResult.InsertedID)
	return nil
}

func (db *DbConfig) Update(filter, update bson.M) error {
	_, err := db.Collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func (db *DbConfig) GetTransactionsByAddress(address string) ([]models.Transaction, error) {
	var transactions []models.Transaction

	// Create a filter to query transactions by address
	filter := bson.M{"address": address}

	// Query the database
	cursor, err := db.Collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Iterate through the cursor and decode each document into a Transaction struct
	for cursor.Next(context.Background()) {
		var doc models.Document
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}
		transactions = doc.Transactions
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
