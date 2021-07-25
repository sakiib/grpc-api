### gRPC API

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
# generate tls certificates
$ make gen-certs
# gRPC server running on port: 8080 & the gateway server running on port: 8000
$ make run
# run the unit-tests
$ make unit-tests
# remove the generated pb.go, pb.gw.go, grpc.pb.go
$ make clean
# set environment variable for enabling tls
$ export tls=true
```

#### or Use the docker image:
```bash
$ docker pull sakibalamin/grpc-api:v0.0.1
$ docker run -p 8000:8000 -it sakibalamin/grpc-api:v0.0.1
```
