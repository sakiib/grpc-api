package gateway

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sakiib/grpc-api/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	pb "github.com/sakiib/grpc-api/gen/go/book"
)

// Run establishes a connection with gRPC server that we're already running
func Run(dialAddr string) error {
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	// Creates a connection with the gRPC server
	connection, err := grpc.DialContext(
		context.Background(),
		dialAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server with: %s", err.Error())
	}

	mux := runtime.NewServeMux()
	err = pb.RegisterBookServiceHandler(context.Background(), mux, connection)
	if err != nil {
		return fmt.Errorf("failed to register gateway with: %s", err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = utils.ProxyPort
	}

	// Creates the gateway server - HTTP/1 server
	gwServer := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/v1") {
				mux.ServeHTTP(w, r)
				return
			}
		}),
	}

	return fmt.Errorf("serving gRPC-Gateway server: %v", gwServer.ListenAndServe())
}