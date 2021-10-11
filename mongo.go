package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Variable to store the client url of the mongo db database.
// Must be configured before calling the API methods.
var ClientUrl string

// Configures the client url of the mongo db database.
// Must be configured before calling the API methods.
func ConfigureMongoClient(url string) {
	ClientUrl = url
}

// Returns the db client of the attached client url.
func GetMongoClient() *mongo.Client {
	if ClientUrl == "" {
		log.Print("Client id not present")
		return nil
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(
		ClientUrl,
	))
	if err != nil {
		log.Print(err)
		return nil
	}
	return client
}

// Drops collections with given names
func DropCollections(dbName string, collections []string) error {
	client := GetMongoClient()
	ctx := context.Background()
	client.Connect(ctx)
	defer client.Disconnect(ctx)
	db := client.Database(dbName)
	for _, index := range collections {
		err := db.Collection(index).Drop(ctx)
		if err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

// Creates collections with given names
func CreateCollections(dbName string, collections []string) error {
	client := GetMongoClient()
	ctx := context.Background()
	client.Connect(ctx)
	defer client.Disconnect(ctx)
	db := client.Database(dbName)
	for _, index := range collections {
		err := db.CreateCollection(ctx, index)
		if err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

func CreateCollection(dbName string, collection string, options *options.CreateCollectionOptions) error {
	client := GetMongoClient()
	ctx := context.Background()
	client.Connect(ctx)
	defer client.Disconnect(ctx)
	db := client.Database(dbName)
	err := db.CreateCollection(ctx, collection, options)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// Inserts a document to the specified database and collection.
// Returns the id of the document upon creation and error otherwise.
// Client must be configured to use this endpoint.
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

// Converts a struct to the bson document used for inserting document.
func Document(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		log.Panic(err)
	}
	err = bson.Unmarshal(data, &doc)
	return
}

// Inserts many documents to the specified database and collection.
// Returns the id of the documents upon creation and error otherwise.
// Client must be configured to use this endpoint.
func WriteMany(dbName string, collection string, doc []interface{}) (*mongo.InsertManyResult, error) {
	client := GetMongoClient()
	ctx := context.Background()
	client.Connect(ctx)
	defer client.Disconnect(ctx)
	db := client.Database(dbName)
	dbCollection := db.Collection(collection)
	result, err := dbCollection.InsertMany(ctx, doc)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return result, nil
}

// Finds a document based on the filter from the specified database and collection.
// Client must be configured to use this endpoint.
func FindOne(dbName string, collection string, filter bson.M, options *options.FindOneOptions) *mongo.SingleResult {
	client := GetMongoClient()
	ctx := context.Background()
	client.Connect(ctx)
	defer client.Disconnect(ctx)
	db := client.Database(dbName)
	dbCollection := db.Collection(collection)
	singleResult := dbCollection.FindOne(ctx, filter, options)
	return singleResult
}

// Finds a cursor of the documents based on the filter and filter options
// from the specified database and collection.
// Returns the documents upon search and error otherwise.
// Client must be configured to use this endpoint.
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

// Updates the documents to the specified database and collection based on identifier filter.
// Returns the id of the document upon update if one entry or the count information if many.
// If update fails, returns an error otherwise.
// Client must be configured to use this endpoint.
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

// Deletes the documents from the specified database and collection based on identifier filter.
// Returns the id of the document upon delete if one entry or the count information if many.
// If delete fails, returns an error otherwise.
// Client must be configured to use this endpoint.
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
