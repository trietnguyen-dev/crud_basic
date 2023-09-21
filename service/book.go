package service

import (
	"book-project/config"
	"book-project/daos"
	"book-project/models"
	pb "book-project/protobuf/gen/go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type BookScv struct {
	dao  *daos.DAO
	conf *config.Config
}

func NewBookSvc(dao *daos.DAO, conf *config.Config) *BookScv {
	return &BookScv{dao: dao, conf: conf}
}
func (s *BookScv) CreateBook(req *pb.CreateBookReq) (*pb.CreateBookRes, error) {

	book := models.Book{
		Name:        req.GetName(),
		Category:    req.GetCategory(),
		Description: req.GetCategory(),
		Author:      req.GetAuthor(),
		Quality:     req.GetQuality(),
		Language:    req.GetLanguage(),
		Price:       req.GetPrice(),
		CreatedAt:   time.Now(),
	}

	err := s.dao.CreateBook(&book)
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}

	return &pb.CreateBookRes{Success: true}, nil

}
func (s *BookScv) GetBook(req *pb.GetBookReq) (*pb.GetBookRes, error) {
	bookCondition := bson.D{{"name", req.Filter.GetName()}}
	exists, err := s.dao.IsBook(bookCondition)
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}
	if !exists {
		return nil, status.Errorf(codes.InvalidArgument, "book does not exist")
	}

	book, err := s.dao.GetBook(bson.M{"name": req.GetFilter().GetName()}, nil)

	if err != nil {
		// Xử lý trường hợp lỗi và trả về lỗi gRPC
		return nil, err
	}
	response := &pb.GetBookRes{
		BookInfo: &pb.BookInfo{
			Id:          book.ID.Hex(),
			Name:        book.Name,
			Category:    book.Category,
			Author:      book.Author,
			Language:    book.Language,
			Description: book.Description,
			Price:       book.Price,
			Quality:     book.Quality,
		},
	}
	// Trả về phản hồi thành công
	return response, nil
}
func (s *BookScv) DeleteBook(req *pb.DeleteBookReq) (*pb.DeleteBookRes, error) {
	bookID, err := primitive.ObjectIDFromHex(req.GetFilter().GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid book ID: %v", err)
	}

	query := bson.M{"_id": bookID}

	_, err = s.dao.DeleteBook(query, nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete book: %v", err)
	}

	return &pb.DeleteBookRes{Success: true}, nil
}
func (s *BookScv) UpdateBook(req *pb.UpdateBookReq) (*pb.UpdateBookRes, error) {
	bookID, err := primitive.ObjectIDFromHex(req.GetBookInfo().GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid book ID: %v", err)
	}
	updatedBook := map[string]interface{}{
		"name":        req.GetBookInfo().GetName(),
		"author":      req.GetBookInfo().GetAuthor(),
		"category":    req.GetBookInfo().GetCategory(),
		"description": req.GetBookInfo().GetDescription(),
		"quality":     req.GetBookInfo().GetQuality(),
		"language":    req.GetBookInfo().GetLanguage(),
		"price":       req.GetBookInfo().GetPrice(),
		"updated_at":  time.Now(),
	}

	book, err := s.dao.UpdateBook(bookID, updatedBook)
	if err != nil {
		return nil, err // Xử lý lỗi và trả về lỗi gRPC
	}
	if book == nil {
		return nil, status.Errorf(codes.InvalidArgument, "book does not exist")
	}

	response := &pb.UpdateBookRes{
		SuccessBook: &pb.BookInfo{
			Name:        book.Name,
			Category:    book.Category,
			Author:      book.Author,
			Language:    book.Language,
			Description: book.Description,
			Price:       book.Price,
			Quality:     book.Quality,
		},
	}
	return response, nil
}
