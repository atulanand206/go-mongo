package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientUrl string

func ConfigureMongoClient(url string) {
	clientUrl = url
}

func GetMongoClient() *mongo.Client {
	if clientUrl == "" {
		log.Print("Client id not present")
		return nil
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(
		clientUrl,
	))
	if err != nil {
		log.Print(err)
		return nil
	}
	return client
}

func Write(dbName string, collection string, doc bson.D) (*mongo.InsertOneResult, error) {
	client := GetMongoClient()
	ctx := context.Background()
	client.Connect(ctx)
	defer client.Disconnect(ctx)
	db := client.Database(dbName)
	dbCollection := db.Collection(collection)
	result, err := dbCollection.InsertOne(ctx, doc)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return result, nil
}

func FindOne(dbName string, collection string, filter bson.M) *mongo.SingleResult {
	client := GetMongoClient()
	ctx := context.Background()
	client.Connect(ctx)
	defer client.Disconnect(ctx)
	db := client.Database(dbName)
	dbCollection := db.Collection(collection)
	singleResult := dbCollection.FindOne(ctx, filter)
	return singleResult
}

func Find(dbName string, collection string, filter bson.M, options *options.FindOptions) (*mongo.Cursor, error) {
	client := GetMongoClient()
	ctx := context.Background()
	client.Connect(ctx)
	defer client.Disconnect(ctx)
	db := client.Database(dbName)
	dbCollection := db.Collection(collection)
	cursor, err := dbCollection.Find(ctx, filter, options)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return cursor, err
}

func Update(dbName string, collection string, identifier bson.M, change bson.D) (*mongo.UpdateResult, error) {
	client := GetMongoClient()
	ctx := context.Background()
	client.Connect(ctx)
	defer client.Disconnect(ctx)
	db := client.Database(dbName)
	dbCollection := db.Collection(collection)
	result, err := dbCollection.UpdateMany(ctx, identifier, change)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return result, nil
}

func Delete(dbName string, collection string, identifier bson.M) (*mongo.DeleteResult, error) {
	client := GetMongoClient()
	ctx := context.Background()
	client.Connect(ctx)
	defer client.Disconnect(ctx)
	db := client.Database(dbName)
	dbCollection := db.Collection(collection)
	result, err := dbCollection.DeleteMany(ctx, identifier)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return result, nil
}
