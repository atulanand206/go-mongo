package mongo

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	mg "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockDB struct {
	Data map[string]map[interface{}]interface{}
}

func NewMockDb() DBConn {
	return &MockDB{}
}

func (db *MockDB) CreateCollection(collection string) error {
	if _, ok := db.Data[collection]; !ok {
		db.Data[collection] = make(map[interface{}]interface{})
	}
	return nil
}

func (db *MockDB) Create(request interface{}, collection string) (err error) {
	doc, err := Document(request)
	if err != nil {
		return
	}
	db.Data[collection][doc["_id"]] = doc
	return
}

func (db *MockDB) CreateMany(request []interface{}, collection string) (err error) {
	for _, v := range request {
		db.Create(v, collection)
	}
	return
}

func Document(request interface{}) (doc bson.M, err error) {
	data, err := bson.Marshal(request)
	if err != nil {
		return
	}
	err = bson.Unmarshal(data, &doc)
	return
}

func (db *MockDB) FindOne(collection string, filters bson.M, findOptions *options.FindOneOptions) (result bson.Raw, err error) {
	for _, entry := range db.Data[collection] {
		if db.isFilterMatch(entry, filters) {
			result, _ = bson.Marshal(entry)
			return
		}
	}
	err = fmt.Errorf("no document found")
	return
}

func (db *MockDB) Find(collection string, filters bson.M, findOptions *options.FindOptions) (results []bson.Raw, err error) {
	results = make([]bson.Raw, 0)
	for _, entry := range db.Data[collection] {
		if db.isFilterMatch(entry, filters) {
			result, _ := bson.Marshal(entry)
			results = append(results, result)
		}
	}
	return
}

func (db *MockDB) Delete(collection string, identifier bson.M) (result *mg.DeleteResult, err error) {
	for key, entry := range db.Data[collection] {
		if db.isFilterMatch(entry, identifier) {
			delete(db.Data[collection], key)
			return
		}
	}
	err = fmt.Errorf("no document found")
	return
}

func (db *MockDB) Update(collection string, identifier bson.M, v interface{}) (result *mg.UpdateResult, err error) {
	for _, entry := range db.Data[collection] {
		if db.isFilterMatch(entry, identifier) {
			db.Create(v, collection)
			return
		}
	}
	err = fmt.Errorf("no document found")
	return
}

func (db *MockDB) isFilterMatch(entry interface{}, filters bson.M) (bool bool) {
	var x = 0
	for k, filter := range filters {
		val := reflect.ValueOf(filter)
		if val.Kind() == reflect.Map {
			for kind, values := range val.MapKeys() {
				strct := val.MapIndex(values)
				if kind == 0 {
					for _, v := range strct.Interface().([]string) {
						if entry.(bson.M)[k] == v {
							x++
							break
						}
					}
				}
			}
		} else {
			if entry.(bson.M)[k] == filter {
				x++
			}
		}
	}
	return x == len(filters)
}
