package service

import (
	"errors"
	pb "github.com/sakiib/grpc-api/gen/go/book"
	"sync"
)

type BookStore interface {
	Set(*pb.Book) error
	Get(string) (*pb.Book, error)
	GetAll() []*pb.Book
	GetData() map[string]*pb.Book
	DeleteBook(string) ([]*pb.Book, error)
	UpdateBook(book *pb.Book) ([]*pb.Book, error)
}

type InMemoryStore struct {
	mutex sync.Mutex
	data  map[string]*pb.Book
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]*pb.Book),
	}
}

func (store *InMemoryStore) GetData() map[string]*pb.Book {
	return store.data
}

func (store *InMemoryStore) Set(book *pb.Book) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if book == nil {
		return errors.New("failed to set new book, request body is nil")
	}

	if _, exists := store.data[book.Id]; exists {
		return errors.New("failed to set new book, book with the given id already exists")
	}

	store.data[book.Id] = book
	return nil
}

func (store *InMemoryStore) Get(id string) (*pb.Book, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if id == "" {
		return nil, errors.New("failed to get book, request body is empty")
	}

	if _, exists := store.data[id]; !exists {
		return nil, errors.New("failed to get book, book with the given id not found")
	}

	return store.data[id], nil
}

func (store *InMemoryStore) GetAll() []*pb.Book {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	books := make([]*pb.Book, 0)
	for _, book := range store.data {
		books = append(books, book)
	}
	return books
}

func (store *InMemoryStore) DeleteBook(id string) ([]*pb.Book, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, exist := store.data[id]; !exist {
		return nil, errors.New("failed to delete book with given key, id not found")
	}
	delete(store.data, id)

	books := make([]*pb.Book, 0)
	for _, book := range store.data {
		books = append(books, book)
	}
	return books, nil
}

func (store *InMemoryStore) UpdateBook(book *pb.Book) ([]*pb.Book, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if book == nil {
		return nil, errors.New("failed update book, request body is nil")
	}

	if book.Id == "" {
		return nil, errors.New("failed update book, id not found in request body")
	}

	if _, exist := store.data[book.Id]; !exist {
		return nil, errors.New("failed to update book with given key, id not found")
	}

	store.data[book.Id] = book

	books := make([]*pb.Book, 0)
	for _, b := range store.data {
		books = append(books, b)
	}
	return books, nil
}
