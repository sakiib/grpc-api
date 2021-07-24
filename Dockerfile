FROM golang:latest as build
WORKDIR /app
COPY . ./
RUN go build -o bin/grpc server/main.go
EXPOSE 8000
ENTRYPOINT ["./bin/grpc"]