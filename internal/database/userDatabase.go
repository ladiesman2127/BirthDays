package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userCollection = "users"

type UserDB struct {
	client *mongo.Client
	name   *string
}

func New(dbName *string) *UserDB {

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", DB_HOST, DB_PORT)))
	if err != nil {
		log.Fatal(err)
	}

	return &UserDB{db, dbName}
}

func (db *UserDB) UsersCollection() *mongo.Collection {
	return db.client.Database(*db.name).Collection(userCollection)
}
