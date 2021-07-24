package utils

import (
	"encoding/json"
	"fmt"
	pb "github.com/sakiib/grpc-api/gen/go/book"
	"io/ioutil"
	"math/rand"
	"strconv"
)

// RandomBookWithID returns a random book item with the given ID
func RandomBookWithID(id string) *pb.Book {
	book := &pb.Book{
		Id:    id,
		Name:  fmt.Sprintf("Book - %s", id),
		Genre: pb.Genre(RandomInt(0, 3)),
		Author: &pb.Author{
			Name:       RandomName(),
			Popularity: int32(RandomInt(1, 10)),
		},
		Rating: int32(RandomInt(1, 10)),
		Price:  float64(RandomInt(1, 100)),
	}
	return book
}

// RandomInt returns a random integer in a given range [min-max]
func RandomInt(min int, max int) int {
	return rand.Intn(max-min) + min
}

// RandomName returns a random name from the provided list
func RandomName() string {
	name := []string{"sakib", "messi", "ronaldo", "harrykane"}
	return name[RandomInt(0, 3)]
}

// GenerateTestData generates json array & stores it to testdata/testdata.json file
func GenerateTestData() {
	var testdata []*pb.Book

	for i := 0; i < 100; i++ {
		data := RandomBookWithID(strconv.Itoa(i))
		testdata = append(testdata, data)
	}

	byt, _ := json.MarshalIndent(testdata, "", "    ")
	_ = ioutil.WriteFile("testdata/testdata.json", byt, 0644)
}
