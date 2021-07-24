package service

import (
	"context"
	pb "github.com/sakiib/grpc-api/gen/go/book"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"sort"
)

type BookService struct {
	pb.UnimplementedBookServiceServer
	store BookStore
}

func NewBookService(store BookStore) *BookService {
	return &BookService{pb.UnimplementedBookServiceServer{}, store}
}

// CreateBook demonstrates an unary client -> server call that takes a Book item as request body
// and creates the book if the ID isn't already present, returns error otherwise
func (bs *BookService) CreateBook(ctx context.Context, in *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
	book := in.GetBook()
	if book == nil {
		return nil, status.Error(codes.InvalidArgument, "request body is nil")
	}

	if err := bs.store.Set(book); err != nil {
		return nil, status.Error(codes.AlreadyExists, err.Error())
	}

	return &pb.CreateBookResponse{
		Id: book.Id,
	}, nil
}

// GetBook demonstrates an unary client -> server call that takes an ID as request body
// and returns the Book as response if found, not found error otherwise
func (bs *BookService) GetBook(ctx context.Context, in *pb.GetBookRequest) (*pb.GetBookResponse, error) {
	id := in.GetId()

	book, err := bs.store.Get(id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetBookResponse{
		Book: book,
	}, nil
}

// GetBooks demonstrates an unary client -> server call, that takes an empty request body
// and returns all the Books stored in the storage, response is a repeated filed like an array
func (bs *BookService) GetBooks(ctx context.Context, in *pb.EmptyRequest) (*pb.GetBooksResponse, error) {
	books := bs.store.GetAll()
	res := &pb.GetBooksResponse{
		Book: books,
	}
	return res, nil
}

// ListBooks demonstrates a server streaming call, that takes an empty request body
// and streams all the book one by one from the storage
func (bs *BookService) ListBooks(req *pb.EmptyRequest, stream pb.BookService_ListBooksServer) error {
	storedData := bs.store.GetData()
	for key := range storedData {
		if err := stream.Send(&pb.GetBookResponse{
			Book: storedData[key],
		}); err != nil {
			return status.Error(codes.Unknown, "failed to stream data to client")
		}
	}
	return nil
}

// DeleteBook demonstrates an unary client -> server call that takes an ID as request body
// and deletes the book with this ID if present in storage and returns all the remaining Books
func (bs *BookService) DeleteBook(ctx context.Context, in *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	books, err := bs.store.DeleteBook(in.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	res := &pb.DeleteBookResponse{
		Book: books,
	}
	return res, nil
}

// UpdateBook demonstrates an unary client -> server call that takes an ID as request body
// and updates the book with this ID if present in storage and returns all the Books
func (bs *BookService) UpdateBook(ctx context.Context, in *pb.UpdateBookRequest) (*pb.UpdateBookResponse, error) {
	books, err := bs.store.UpdateBook(in.GetBook())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	res := &pb.UpdateBookResponse{
		Book: books,
	}
	return res, nil
}

// BooksCost demonstrates a client streaming call, client streams some Books ID
// and server receives the streams & returns the total cost of these books when client streaming ends
func (bs *BookService) BooksCost(stream pb.BookService_BooksCostServer) error {
	cost := 0.0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.BooksCostResponse{
				Cost: cost,
			})
		}

		if err != nil {
			return status.Error(codes.Unknown, err.Error())
		}

		id := req.GetId()
		if id == "" {
			return status.Error(codes.InvalidArgument, "failed to get id from request streaming")
		}

		book, err := bs.store.Get(req.GetId())
		if err != nil {
			return status.Error(codes.NotFound, err.Error())
		}

		cost += book.GetPrice()
	}
}

// TopRatedBook demonstrates and bi-directions client <==> server streaming call,
// client sends streams of Books & after each call server streams at most two books with the most ratings
func (bs *BookService) TopRatedBook(stream pb.BookService_TopRatedBookServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Error(codes.Unknown, err.Error())
		}

		book := req.GetBook()
		if err = bs.store.Set(book); err != nil {
			return status.Error(codes.AlreadyExists, err.Error())
		}

		books := bs.store.GetAll()
		sort.SliceStable(books, func(i, j int) bool {
			return books[i].GetRating() > books[j].GetRating()
		})

		for idx, b := range books {
			if idx > 1 {
				break
			}
			if err = stream.Send(&pb.TopRatedResponse{
				Book: b,
			}); err != nil {
				return status.Error(codes.Unknown, err.Error())
			}
		}
	}
}
