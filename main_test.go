package main

import (
	"math"
	"testing"

	"os"

	"github.com/mgironi/operation-fire-quasar/location"
)

func TestMain(t *testing.T) {
	oldsArgs := os.Args
	os.Args = []string{"cmd", "-distances=500,424.26,707.10"}
	main()
	os.Args = oldsArgs
}

func TestParseArgs(t *testing.T) {
	oldsArgs := os.Args
	os.Args = []string{"cmd", "-distances=500,424.26,707.10"}
	wantedLength := 3
	wantedValues := []float32{500, 424.26, 707.10}

	distances, messages, err := ParseArgs()

	if err != nil {
		t.Fatalf("Error parsing args %e", err)
	}
	if len(distances) != wantedLength {
		t.Fatalf("Error parsing distances args length is %d, wanted %d", len(distances), wantedLength)
	}
	for i, wantedVal := range wantedValues {
		if !areEquals(wantedVal, distances[i]) {
			t.Fatalf("Error parsing distances args value in list position %d is %f, wanted %f", i, distances[i], wantedVal)
		}
	}
	if len(messages) != 0 {
		t.Fatal("Messages spected empty")
	}
	os.Args = oldsArgs
}

func areEquals(a float32, b float32) bool {
	diff := math.Abs(float64(a) - float64(b))
	return diff < location.FLOAT_COMPARISION_TOLERANCE
}
