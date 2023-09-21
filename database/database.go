package database

import (
	"book-project/config"
	"book-project/helpers"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	UserColl string = "users"
	BookColl string = "books"
)

// InitDB :
func InitDB(conf *config.Config) (*mongo.Database, error) {
	clientOptions := options.Client().SetMaxConnIdleTime(5 * time.Second).SetMaxPoolSize(250).ApplyURI(conf.DB).SetRegistry(bigIntMongoRegistry)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(conf.DBName)

	if err = CreateIndexes(db); err != nil {
		return nil, err
	}

	return db, nil
}

// CreateIndexes :
func CreateIndexes(db *mongo.Database) error {
	userCollection := db.Collection(UserColl)
	bookCollection := db.Collection(BookColl)

	userIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
				{Key: "phone_number", Value: 1},
			},
		},
	}

	bookIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
			},
			// Các tùy chọn index cho book collection
		},
	}

	ctx, cancel := helpers.NewCtx()
	defer cancel()

	_, err := userCollection.Indexes().CreateMany(ctx, userIndexes)

	if err != nil {
		return err
	}
	_, err = bookCollection.Indexes().CreateMany(ctx, bookIndexes)

	return nil
}
