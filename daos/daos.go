package daos

import (
	"book-project/config"
	"book-project/database"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type DAO struct {
	Client   *mongo.Client
	userColl *mongo.Collection
	bookColl *mongo.Collection
}

func NewDAO(conf *config.Config) (*DAO, error) {
	var err error

	db, err := database.InitDB(conf)
	if err != nil {
		return nil, errors.Wrap(err, "database.Init")
	}
	return &DAO{
		Client:   db.Client(),
		userColl: db.Collection(database.UserColl),
		bookColl: db.Collection(database.BookColl),
	}, nil
}
