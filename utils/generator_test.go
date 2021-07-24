package utils

import (
	"testing"
)

func TestRandomBookWithID(t *testing.T) {
	book := RandomBookWithID("1")
	if !(book.Id == "1" && (book.Name == "Book - 1")) {
		t.Error()
	}
}

func TestRandomInt(t *testing.T) {
	for i := 1; i < 100; i++ {
		val := RandomInt(0, 1000000000)
		if !(val >= 1 && val <= 1000000000) {
			t.Error()
		}
	}
}
