package daos

import (
	"book-project/helpers"
	"book-project/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DAO) CreateBook(art *models.Book) error {

	ctx, cancel := helpers.NewCtx()
	defer cancel()

	_, err := d.bookColl.InsertOne(ctx, art)
	return err
}
func (d *DAO) IsBook(condition bson.D) (bool, error) {
	// Sử dụng FindOne để kiểm tra sự tồn tại của một bản ghi dựa trên điều kiện
	ctx, cancel := helpers.NewCtx()
	defer cancel()

	result := d.bookColl.FindOne(ctx, condition)

	if result.Err() == nil {
		// Bản ghi tồn tại
		return true, nil
	} else if result.Err() == mongo.ErrNoDocuments {
		// Không tìm thấy bản ghi, tức là không tồn tại
		return false, nil
	} else {
		// Có lỗi xảy ra
		return false, result.Err()
	}
}

func (d *DAO) GetBook(queries bson.M, opts *options.FindOneOptions) (art *models.Book, err error) {
	ctx, cancel := helpers.NewCtx()
	defer cancel()

	if err = d.bookColl.FindOne(ctx, queries, opts).Decode(&art); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "MongoDB.FindOne")
	}

	return
}
func (d *DAO) DeleteBook(queries bson.M, opts *options.FindOneOptions) (art *models.Book, err error) {
	ctx, cancel := helpers.NewCtx()
	defer cancel()

	_, err = d.bookColl.DeleteOne(ctx, queries)
	if err != nil {
		return nil, errors.Wrap(err, "MongoDB.DeleteOne")
	}
	return
}
func (d *DAO) UpdateBook(id primitive.ObjectID, updatedData map[string]interface{}) (*models.Book, error) {
	ctx, cancel := helpers.NewCtx()
	defer cancel()

	var art *models.Book

	// Tạo một bản ghi cập nhật sử dụng phép toán $set
	update := bson.M{"$set": updatedData}

	returnDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &returnDoc,
	}

	if err := d.bookColl.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&art); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "MongoDB.FindOneAndUpdate")
	}

	return art, nil
}
