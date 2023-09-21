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

func (d *DAO) CreateUser(art *models.User) error {

	ctx, cancel := helpers.NewCtx()
	defer cancel()

	_, err := d.userColl.InsertOne(ctx, art)
	return err
}

func (d *DAO) IsExist(condition bson.D) (bool, error) {
	// Sử dụng FindOne để kiểm tra sự tồn tại của một bản ghi dựa trên điều kiện
	ctx, cancel := helpers.NewCtx()
	defer cancel()

	result := d.userColl.FindOne(ctx, condition)

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
func (d *DAO) GetHashedPassword(condition bson.D) (string, error) {
	// Sử dụng FindOne để lấy mật khẩu đã băm từ cơ sở dữ liệu
	ctx, cancel := helpers.NewCtx()
	defer cancel()

	var user struct {
		Password string `bson:"password"`
	}

	result := d.userColl.FindOne(ctx, condition).Decode(&user)
	if result != nil {
		if result == mongo.ErrNoDocuments {
			// Không tìm thấy bản ghi, tức là không tồn tại
			return "", nil
		}
		// Xử lý các trường hợp lỗi khác
		return "", result
	}

	return user.Password, nil
}

func (d *DAO) GetInfoUser(queries bson.M, opts *options.FindOneOptions) (art *models.User, err error) {
	ctx, cancel := helpers.NewCtx()
	defer cancel()

	if err = d.userColl.FindOne(ctx, queries, opts).Decode(&art); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "MongoDB.FindOne")
	}

	return
}
func (d *DAO) UpdateUser(id primitive.ObjectID, updatedData map[string]interface{}) (*models.User, error) {
	ctx, cancel := helpers.NewCtx()
	defer cancel()

	var art *models.User

	// Tạo một bản ghi cập nhật sử dụng phép toán $set
	update := bson.M{"$set": updatedData}

	returnDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &returnDoc,
	}

	if err := d.userColl.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&art); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "MongoDB.FindOneAndUpdate")
	}

	return art, nil
}
