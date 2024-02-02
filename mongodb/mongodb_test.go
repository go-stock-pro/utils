package mongodb_test

import (
	"context"
	"testing"

	env "github.com/joho/godotenv"
	"github.com/rohanraj7316/utils/mongodb"
)

const (
	envFile = ".env"
)

var loadEnv = env.Load

type Student struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

func TestMongoDbConnection(t *testing.T) {
	ctx := context.Background()

	// load env config
	err := loadEnv(envFile)
	if err != nil {
		t.Error(err)
	}

	_, err = mongodb.New(ctx)
	if err != nil {
		t.Errorf("failed to connect to mongodb: %s", err)
	}
}

func TestMultipleMongoDbConnection(t *testing.T) {
	ctx := context.Background()

	// load env config
	err := loadEnv(envFile)
	if err != nil {
		t.Error(err)
	}

	_, err = mongodb.New(ctx)
	if err != nil {
		t.Errorf("failed to connect to mongodb: %s", err)
	}
}

func TestFindAll(t *testing.T) {
	ctx := context.Background()

	// load env config
	err := loadEnv(envFile)
	if err != nil {
		t.Error(err)
	}

	conn, err := mongodb.New(ctx)
	if err != nil {
		t.Errorf("failed to connect to mongodb: %s", err)
	}

	filter := Student{
		Name: "Rohan Raj",
	}
	dst := &[]Student{}

	err = conn.Find(ctx, filter, dst)
	if err != nil {
		t.Errorf("failed to fetch from mongodb: %s", err)
	}

	t.Errorf("%+v", dst)
}

// func TestFindById(t *testing.T) {
// 	ctx := context.Background()
// 	err := loadEnv(envFile)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	conn, err := New(ctx)
// 	if err != nil {
// 		t.Errorf("failed to connect to mongodb: %s", err)
// 	}
// 	coll := conn.Client.Database("sample_mflix").Collection("movies")
// 	id := "634cebdf950badb537f62e22"
// 	objID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	// var result bson.D
// 	filter := bson.M{"_id": bson.M{"$eq": objID}}
// 	_, err = coll.Find(ctx, filter)
// 	// if cursor.All()

// }

// func TestUpdate(t *testing.T) {
// 	ctx := context.Background()
// 	err := loadEnv(envFile)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	conn, err := New()
// 	if err != nil {
// 		t.Errorf("failed to connect to mongodb: %s", err)
// 	}
// 	coll := conn.Client.Database("sample_mflix").Collection("movies")
// 	title := "The Favourite"
// 	filter := bson.D{{Key: "title", Value: title}}
// 	update := bson.D{{Key: "$set", Value: bson.D{{Key: "type", Value: "documentary"}}}}
// 	_, err = coll.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

// func TestUpdateById(t *testing.T) {
// 	ctx := context.Background()
// 	err := loadEnv(envFile)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	conn, err := New()
// 	if err != nil {
// 		t.Errorf("failed to connect to mongodb: %s", err)
// 	}
// 	coll := conn.Client.Database("sample_mflix").Collection("movies")
// 	id, _ := primitive.ObjectIDFromHex("634cebdf950badb537f62e22")
// 	filter := bson.D{{Key: "_id", Value: id}}
// 	update := bson.D{{Key: "$set", Value: bson.D{{Key: "year", Value: 2019}}}}
// 	_, err = coll.UpdateOne(ctx, filter, update)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestInsert(t *testing.T) {
	ctx := context.Background()

	err := loadEnv(envFile)
	if err != nil {
		t.Error(err)
	}

	conn, err := mongodb.New(ctx)
	if err != nil {
		t.Errorf("failed to connect to mongodb: %s", err)
	}

	stu := Student{
		Name: "Rohan Raj",
		Age:  25,
	}

	_, err = conn.Create(ctx, stu)
	if err != nil {
		t.Errorf("failed to insert into mongodb: %s", err)
	}
}

// func TestDelete(t *testing.T) {
// 	ctx := context.Background()
// 	err := loadEnv(envFile)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	conn, err := New()
// 	if err != nil {
// 		t.Errorf("failed to connect to mongodb: %s", err)
// 	}
// 	coll := conn.Client.Database("sample_mflix").Collection("movies")
// 	filter := bson.D{{Key: "title", Value: "Jurassic World: Fallen Kingdom"}}
// 	_, err = coll.DeleteOne(ctx, filter)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }
