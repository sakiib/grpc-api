package main

import (
	"flag"
	"fmt"
	"github.com/sakiib/grpc-api/gateway"
	pb "github.com/sakiib/grpc-api/gen/go/book"
	"github.com/sakiib/grpc-api/service"
	"github.com/sakiib/grpc-api/utils"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	port := flag.String("port", utils.GRPCServerPort, "server port")
	flag.Parse()

	log.Printf("starting the server on port: %s", *port)

	grpcServer := grpc.NewServer()
	bookServer := service.NewBookService(service.NewInMemoryStore())

	pb.RegisterBookServiceServer(grpcServer, bookServer)

	listener, err := net.Listen(utils.Network, fmt.Sprintf("0.0.0.0:%s", *port))
	if err != nil {
		log.Fatalf("failed to listen with: %s", err.Error())
	}

	// Run the grpc server in a go routine as this one is a blocking function
	go func() {
		log.Fatalf("grpc server error: %s", grpcServer.Serve(listener))
	}()

	log.Fatalf("gateway error: %s", gateway.Run("dns:///"+fmt.Sprintf("0.0.0.0:%s", *port)))
}
