package mongo

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConn interface {
	CreateCollection(collection string) error
	Create(request interface{}, collection string) (err error)
	CreateMany(request []interface{}, collection string) (err error)
	FindOne(collection string, filters bson.M, findOptions *options.FindOneOptions) (result bson.Raw, err error)
	Find(collection string, filters bson.M, findOptions *options.FindOptions) (result []bson.Raw, err error)
	Delete(collection string, identifier bson.M) (result *mg.DeleteResult, err error)
	Update(collection string, identifier bson.M, doc interface{}) (result *mg.UpdateResult, err error)
}

type DB struct {
	client *MongoClient
}

func NewDb(dbName string) DBConn {
	return &DB{
		client: NewMongoClient(os.Getenv("MONGO_CLIENT_ID"), dbName),
	}
}

func (db *DB) CreateCollection(collection string) error {
	return db.client.CreateCollection(collection, &options.CreateCollectionOptions{})
}

func (db *DB) Create(request interface{}, collection string) (err error) {
	requestDto, err := db.client.Document(request)
	if err != nil {
		return
	}
	_, err = db.client.Write(collection, *requestDto)
	return
}

func (db *DB) CreateMany(request []interface{}, collection string) (err error) {
	_, err = db.client.WriteMany(collection, request)
	return
}

func (db *DB) FindOne(collection string, filters bson.M, findOptions *options.FindOneOptions) (result bson.Raw, err error) {
	res := db.client.FindOne(collection, filters, findOptions)
	if err = res.Err(); err != nil {
		return
	}
	return res.DecodeBytes()
}

func (db *DB) Find(collection string, filters bson.M, findOptions *options.FindOptions) (result []bson.Raw, err error) {
	cursor, err := db.client.Find(collection, filters, findOptions)
	if err != nil {
		return
	}
	result, err = db.DecodeRaw(cursor)
	return
}

func (db *DB) Delete(collection string, identifier bson.M) (result *mg.DeleteResult, err error) {
	return db.client.Delete(collection, identifier)
}

func (db *DB) Update(collection string, identifier bson.M, doc interface{}) (result *mg.UpdateResult, err error) {
	requestDto, err := db.client.Document(doc)
	if err != nil {
		return
	}
	return db.client.Update(collection, identifier, bson.D{primitive.E{Key: "$set", Value: *requestDto}})
}

func (db *DB) DecodeRaw(cursor *mg.Cursor) (documents []bson.Raw, err error) {
	for cursor.Next(context.Background()) {
		documents = append(documents, cursor.Current)
	}
	return
}
