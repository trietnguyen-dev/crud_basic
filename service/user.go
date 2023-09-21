package service

import (
	"book-project/config"
	"book-project/daos"
	"book-project/helpers"
	"book-project/models"
	pb "book-project/protobuf/gen/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/mail"
	"time"
)

type UserScv struct {
	dao  *daos.DAO
	conf *config.Config
}

// NewUserSvc returns new UserSvc struct
func NewUserSvc(dao *daos.DAO, conf *config.Config) *UserScv {
	return &UserScv{dao: dao, conf: conf}
}
func (s *UserScv) RegisterUser(req *pb.RegisterUserReq) (*pb.RegisterUserRes, error) {
	formatEmail, err := mail.ParseAddress(req.GetEmail())
	if err != nil {
		return nil, err
	}

	email := formatEmail.Address // Lấy địa chỉ email từ *mail.Address

	emailCondition := bson.D{{"email", email}}
	exists, err := s.dao.IsExist(emailCondition)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, status.Errorf(codes.InvalidArgument, "email already exists")
	}

	// Kiểm tra mật khẩu
	if req.GetPassword() != req.GetRePassword() {
		return nil, status.Errorf(codes.InvalidArgument, "passwords do not match")
	}
	isPhoneNumber := helpers.IsValidPhoneNumber(req.GetPhoneNumber())

	if !isPhoneNumber {
		return nil, status.Errorf(codes.InvalidArgument, "phone number invalid")
	}

	hash, _ := helpers.HashPassword(req.GetPassword())

	// Tạo người dùng
	user := models.User{
		FullName:    req.GetFullName(),
		Email:       req.GetEmail(),
		Password:    hash,
		PhoneNumber: req.GetPhoneNumber(),
		CreatedAt:   time.Now(),
	}

	err = s.dao.CreateUser(&user)
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}

	// Trả về phản hồi thành công
	return &pb.RegisterUserRes{Success: true}, nil
}

func (s *UserScv) LoginUser(req *pb.LoginUserReq) (*pb.LoginUserRes, error) {
	// Kiểm tra sự tồn tại của tài khoản dựa trên email
	emailCondition := bson.D{{"email", req.GetEmail()}}
	exists, err := s.dao.IsExist(emailCondition)
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}
	if !exists {
		return nil, status.Errorf(codes.InvalidArgument, "email does not exist")
	}

	hashedPassword, err := s.dao.GetHashedPassword(emailCondition)
	// Kiểm tra mật khẩu
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.GetPassword()))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "error password ")
	}

	// Đăng nhập thành công
	return &pb.LoginUserRes{Success: true}, nil
}
func (s *UserScv) GetInfoUser(req *pb.GetInfoUserReq) (*pb.GetInfoUserRes, error) {
	emailCondition := bson.D{{"email", req.Filter.GetEmail()}}
	exists, err := s.dao.IsExist(emailCondition)
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}
	if !exists {
		return nil, status.Errorf(codes.InvalidArgument, "email does not exist")
	}
	// Sử dụng DAO để lấy thông tin người dùng từ cơ sở dữ liệu dựa trên điều kiện truy vấn
	user, err := s.dao.GetInfoUser(bson.M{"email": req.GetFilter().GetEmail()}, nil)

	if err != nil {
		// Xử lý trường hợp lỗi và trả về lỗi gRPC
		return nil, err
	}

	// Tạo phản hồi với thông tin người dùng nếu tìm thấy
	response := &pb.GetInfoUserRes{
		UserInfo: &pb.UserInfo{
			FullName:    user.FullName,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
		},
	}

	// Trả về phản hồi thành công
	return response, nil
}
func (s *UserScv) UpdateUser(req *pb.UpdateUserReq) (*pb.UpdateUserRes, error) {
	userIdStr := req.GetId()
	// Chuyển đổi chuỗi userIdStr thành ObjectID
	userId, err := primitive.ObjectIDFromHex(userIdStr)

	if err != nil {
		return nil, err // Xử lý lỗi nếu userIdStr không hợp lệ
	}
	updatedUser := map[string]interface{}{
		"full_name":    req.GetFullName(),
		"phone_number": req.GetPhoneNumber(),
		"email":        req.GetEmail(),
		"updated_at":   time.Now(),
	}

	user, err := s.dao.UpdateUser(userId, updatedUser)
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}
	if user == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user does not exist")
	}

	// Tạo phản hồi với thông tin người dùng cập nhật
	response := &pb.UpdateUserRes{
		UserInfo: &pb.UserInfo{
			FullName:    user.FullName,
			PhoneNumber: user.PhoneNumber,
			Email:       user.Email,
		},
	}

	return response, nil
}
