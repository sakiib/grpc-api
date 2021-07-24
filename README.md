### gRPC API

#### Project Structure
```bash
├── buf.gen.yaml
├── buf.yaml
├── Dockerfile
├── gateway
│   └── gateway.go
├── gen
│   └── go
│       └── book
│           ├── author.pb.go
│           ├── book.pb.go
│           ├── book_service_grpc.pb.go
│           ├── book_service.pb.go
│           └── book_service.pb.gw.go
├── go.mod
├── go.sum
├── LICENSE
├── Makefile
├── proto
│   ├── book
│   │   ├── author.proto
│   │   ├── book.proto
│   │   └── book_service.proto
│   └── google
│       └── api
│           ├── annotations.proto
│           └── http.proto
├── README.md
├── server
│   └── main.go
├── service
│   ├── book_service.go
│   ├── book_service_test.go
│   └── book_store.go
├── testdata
│   └── testdata.json
├── utils
│   ├── constants.go
│   ├── generator.go
│   └── generator_test.go
└── vendor

```

---
#### Clone the repo & run locally:
```bash
# install the dependencies
$ make install
# generate pb.go, pb.gw.go, grpc.pb.go
$ make gen
# run the gRPC server & the gateway
# format go code
$ make fmt
# gRPC server running on port: 8080 & the gateway server running on port: 8000
$ make run
# run the unit-tests
$ make unit-tests
# remove the generated pb.go, pb.gw.go, grpc.pb.go
$ make clean
```

#### or Use the docker image:
```bash
$ docker pull sakibalamin/grpc-api:v0.0.1
$ docker run -p 8000:8000 -it sakibalamin/grpc-api:v0.0.1
```