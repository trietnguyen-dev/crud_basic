package server

import (
	pb "book-project/protobuf/gen/go"
	"book-project/service"
	"context"
)

type BookSrv struct {
	pb.UnimplementedBookServer
	bookSvc *service.BookScv
}

func NewBookSrvServer(bookSvc *service.BookScv) *BookSrv {
	return &BookSrv{
		bookSvc: bookSvc,
	}
}

func (s *BookSrv) CreateBook(c context.Context, req *pb.CreateBookReq) (*pb.CreateBookRes, error) {

	response, err := s.bookSvc.CreateBook(req)
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}

	return response, nil // Trả về phản hồi thành công
}
func (s *BookSrv) GetBook(c context.Context, req *pb.GetBookReq) (*pb.GetBookRes, error) {
	response, err := s.bookSvc.GetBook(req)
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}

	return response, nil // Trả về phản hồi thành công
}
func (s *BookSrv) DeleteBook(c context.Context, req *pb.DeleteBookReq) (*pb.DeleteBookRes, error) {
	response, err := s.bookSvc.DeleteBook(req)
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}

	return response, nil // Trả về phản hồi thành công
}
func (s *BookSrv) UpdateBook(c context.Context, req *pb.UpdateBookReq) (*pb.UpdateBookRes, error) {
	response, err := s.bookSvc.UpdateBook(req)
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}

	return response, nil // Trả về phản hồi thành công
}
