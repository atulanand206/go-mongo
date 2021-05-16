# go-mongo

[![Go Reference](https://pkg.go.dev/badge/github.com/atulanand206/go-mongo.svg)](https://pkg.go.dev/github.com/atulanand206/go-mongo)

A library exposing an implementation of [mongodb/mongo-go-driver](https://github.com/mongodb/mongo-go-driver). 

You can configure the client URL and start interacting with the mongo db database without worrying about managing database connections. Every endpoint instantiates a database client and is independent of each other. Every endpoint can also be worked on a pool of go routines bringing concurrency benefits.

## What is MongoDB?

MongoDB is a document based database service. It has a cloud-oriented product called MongoDB Atlas which handles the replication strategy without external effort. The schema is not required to be defined as the documents can contain anything.

## Installation

The recommended way to get started using the MongoDB Go driver is by using go modules to install the dependency in your project. This can be done either by importing packages from [go-mongo](github.com/atulanand206/go-mongo) and having the build step install the dependency or by explicitly running
```go
go get github.com/atulanand206/go-mongo
```

## Pre-requisites

- Mongo Client URL should be configured for the cluster of your choice.
- Database and Collection must be created on the database.

## How to implement

- Define the required variables in the environment
```go
mongoClientId := os.Getenv("MONGO_CLIENT_ID")
database = os.Getenv("DATABASE")
collection = os.Getenv("MONGO_COLLECTION")
```
- Configure the client before using any of the endpoints.
```go
mongo.ConfigureMongoClient(mongoClientId)
```
    
#### Examples of bson documents

- *bson.D*: Convert an object to an ordered representation a bson document.
```go
document, err := mongo.Document(&object)
if err != nil {
    log.Print(err)
    return
}
```

- *bson.M*: An unordered representation of a bson document.
```go
filter := bson.M{"name": name}
```

- *FindOptions*: Object used for specifying sort and additional criteria.
```go
opts := options.Find()
opts.SetSort(bson.D{primitive.E{Key: "rating", Value: -1}})
```

- Please refer to the [bson documentation](https://docs.mongodb.com/manual/core/document/) for information on the usage of the documents. 

#### Examples of endpoints

- Write a document to the mongo collection.
```go
mongo.Write(database, collection, *document)
```

- Update a document in the mongo collection based on filter criteria.
```go
mongo.Update(database, collection, filter, *document)
```

- Delete a document from the mongo collection based on filter criteria.
```go
mongo.Delete(database, collection, filter, *document)
```

- Find documents based on filters and options. The decoding depends on the actual definition of the object and hence required with the implementation.
```go
cursor, err := mongo.Find(database, collection, filter, opts)
for cursor.Next(context.Background()) {
    var object Object
    err := cursor.Decode(&object)
    if err != nil {
        log.Print(err)
        return
    }
    response = append(response, object)
}
```

- Find a document based on filters and decode to object. First item is returned if multiple documents match the criteria.
```go
response := mongo.FindOne(database, collection, filter)
var object Object
if err = response.Decode(&object); err != nil {
    log.Print(err)
    return
}
```

## Author

- Atul Anand

## License

The MongoDB Go Driver is licensed under the [MIT License](LICENSE).