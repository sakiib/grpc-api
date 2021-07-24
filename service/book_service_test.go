package service

import (
	"context"
	pb "github.com/sakiib/grpc-api/gen/go/book"
	"github.com/sakiib/grpc-api/utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"io"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestBookService_CreateBook(t *testing.T) {
	storage := NewInMemoryStore()
	serverAddress := runTestServer(t, storage)
	client := getTestClient(t, serverAddress)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testCreateBook(client, ctx, t)
}

func TestBookService_GetBook(t *testing.T) {
	storage := NewInMemoryStore()
	serverAddress := runTestServer(t, storage)
	client := getTestClient(t, serverAddress)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insertBooks(storage)
	testGetBook(client, ctx, t)
}

func TestBookService_GetBooks(t *testing.T) {
	storage := NewInMemoryStore()
	serverAddress := runTestServer(t, storage)
	client := getTestClient(t, serverAddress)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insertBooks(storage)
	testGetBooks(client, ctx, t)
}

func TestBookService_ListBooks(t *testing.T) {
	storage := NewInMemoryStore()
	serverAddress := runTestServer(t, storage)
	client := getTestClient(t, serverAddress)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insertBooks(storage)
	testListBooks(client, ctx, t)
}

func TestBookService_DeleteBooks(t *testing.T) {
	storage := NewInMemoryStore()
	serverAddress := runTestServer(t, storage)
	client := getTestClient(t, serverAddress)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insertBooks(storage)
	testDeleteBooks(client, ctx, t)
}

func TestBookService_UpdateBook(t *testing.T) {
	storage := NewInMemoryStore()
	serverAddress := runTestServer(t, storage)
	client := getTestClient(t, serverAddress)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insertBooks(storage)
	testUpdateBook(client, ctx, t)
}

func TestBookService_BooksCost(t *testing.T) {
	storage := NewInMemoryStore()
	serverAddress := runTestServer(t, storage)
	client := getTestClient(t, serverAddress)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insertBooks(storage)
	expCost := 0.0
	for key := range storage.data {
		expCost += storage.data[key].GetPrice()
	}
	testBooksCost(client, ctx, t, expCost)
}

func TestBookService_TopRatedBook(t *testing.T) {
	storage := NewInMemoryStore()
	serverAddress := runTestServer(t, storage)
	client := getTestClient(t, serverAddress)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.TopRatedBook(ctx)
	require.NoError(t, err)

	waitc := make(chan struct{})
	go func() {
		for {
			_, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			require.NoError(t, err)
		}
	}()

	for i := 0; i < 5; i++ {
		book := utils.RandomBookWithID(strconv.Itoa(i))
		err = stream.Send(&pb.TopRatedRequest{
			Book: book,
		})
		require.NoError(t, err)
	}
	stream.CloseSend()
	<-waitc
}

func testBooksCost(client pb.BookServiceClient, ctx context.Context, t *testing.T, expCost float64) {
	streams, err := client.BooksCost(ctx)
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		req := &pb.BooksCostRequest{
			Id: strconv.Itoa(i),
		}
		err = streams.Send(req)
		require.NoError(t, err)
	}

	cost, err := streams.CloseAndRecv()

	require.NoError(t, err)
	require.Equal(t, cost.GetCost(), expCost)
}

func testUpdateBook(client pb.BookServiceClient, ctx context.Context, t *testing.T) {
	for i := 0; i < 5; i++ {
		book := utils.RandomBookWithID(strconv.Itoa(i))
		res, err := client.UpdateBook(ctx, &pb.UpdateBookRequest{
			Book: book,
		})
		require.NoError(t, err)
		require.NotEmpty(t, res)
	}
}

func testDeleteBooks(client pb.BookServiceClient, ctx context.Context, t *testing.T) {
	for i := 0; i < 15; i++ {
		res, err := client.DeleteBook(ctx, &pb.DeleteBookRequest{
			Id: strconv.Itoa(i),
		})
		if i < 10 {
			require.NoError(t, err)
			require.NotEmpty(t, res)
		} else {
			require.Error(t, err)
			require.Empty(t, res)
		}
	}
}

func testListBooks(client pb.BookServiceClient, ctx context.Context, t *testing.T) {
	res, err := client.ListBooks(ctx, &pb.EmptyRequest{})
	require.NoError(t, err)
	require.NotEmpty(t, res)
}

func testGetBooks(client pb.BookServiceClient, ctx context.Context, t *testing.T) {
	res, err := client.GetBooks(ctx, &pb.EmptyRequest{})
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, len(res.GetBook()), 10)
}

func testGetBook(client pb.BookServiceClient, ctx context.Context, t *testing.T) {
	for i := 0; i < 15; i++ {
		_, err := client.GetBook(ctx, &pb.GetBookRequest{
			Id: strconv.Itoa(i),
		})
		if i < 10 {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func testCreateBook(client pb.BookServiceClient, ctx context.Context, t *testing.T) {
	for i := 0; i < 10; i++ {
		book := utils.RandomBookWithID(strconv.Itoa(i))
		_, err := client.CreateBook(ctx, &pb.CreateBookRequest{
			Book: book,
		})
		require.NoError(t, err)
	}
}

func insertBooks(storage *InMemoryStore) {
	for i := 0; i < 10; i++ {
		book := utils.RandomBookWithID(strconv.Itoa(i))
		storage.data[book.Id] = book
	}
}

func getTestClient(t *testing.T, serverAddress string) pb.BookServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	require.NoError(t, err)

	return pb.NewBookServiceClient(conn)
}

func runTestServer(t *testing.T, storage *InMemoryStore) string {
	grpcServer := grpc.NewServer()
	bookServer := NewBookService(storage)

	pb.RegisterBookServiceServer(grpcServer, bookServer)

	listener, err := net.Listen(utils.Network, ":0")
	require.NoError(t, err)

	go func() {
		err = grpcServer.Serve(listener)
		require.NoError(t, err)
	}()

	return listener.Addr().String()
}
