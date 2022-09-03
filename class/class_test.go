package class_test

import (
	"os"
	"testing"

	. "github.com/thara/godiva/class"
)

func TestParse(t *testing.T) {
	f, err := os.Open("../testdata/HelloWorld.class")
	if err != nil {
		t.Fatal(err)
	}

	_, err = Parse(f)
	if err != nil {
		t.Fatal(err)
	}
}
