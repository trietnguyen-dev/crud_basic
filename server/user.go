package server

import (
	pb "book-project/protobuf/gen/go"
	"book-project/service"
	"context"
)

// UserSrv : struct
type UserSrv struct {
	pb.UnimplementedUserSrvServer
	userSvc *service.UserScv
}

// NewUserSrvServer return UserSrv struct
func NewUserSrvServer(userSvc *service.UserScv) *UserSrv {
	return &UserSrv{
		userSvc: userSvc,
	}
}

func (s *UserSrv) RegisterUser(c context.Context, req *pb.RegisterUserReq) (*pb.RegisterUserRes, error) {

	// Gọi phương thức RegisterUser của service để xử lý đăng ký
	response, err := s.userSvc.RegisterUser(req)
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}

	return response, nil // Trả về phản hồi thành công
}

func (s *UserSrv) LoginUser(c context.Context, req *pb.LoginUserReq) (*pb.LoginUserRes, error) {
	response, err := s.userSvc.LoginUser(req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *UserSrv) GetInfoUser(c context.Context, req *pb.GetInfoUserReq) (*pb.GetInfoUserRes, error) {

	response, err := s.userSvc.GetInfoUser(req)
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}
	return response, nil
}

func (s *UserSrv) UpdateUser(c context.Context, req *pb.UpdateUserReq) (*pb.UpdateUserRes, error) {
	response, err := s.userSvc.UpdateUser(req)
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}
	return response, nil
}
