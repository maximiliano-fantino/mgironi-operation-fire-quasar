package main

import (
	"os"
	"reflect"
	"testing"

	test "github.com/mgironi/operation-fire-quasar/_test"
)

func TestParseArgs(t *testing.T) {
	oldsArgs := os.Args
	os.Args = []string{"cmd", "-distances=500,424.26,707.10", "-messages=this..the.complete.message,.is.the..message,.is...message"}
	wantedDistances := []float32{500, 424.26, 707.10}
	wantedMessages := [][]string{
		{"this", "", "the", "complete", "message"},
		{"", "is", "the", "", "message"},
		{"", "is", "", "", "message"},
	}

	distances, messages, err := ParseArgs()

	if err != nil {
		t.Errorf("Error parsing args %e", err)
	}
	if len(distances) != len(wantedDistances) {
		t.Errorf("Error parsing distances args length is %d, wanted %d", len(distances), len(wantedDistances))
	}
	for i, wantedVal := range wantedDistances {
		if !test.AreFloatsEquals(float64(wantedVal), float64(distances[i])) {
			t.Errorf("Error parsing distances args value in list position %d is %f, wanted %f", i, distances[i], wantedVal)
		}
	}
	if !reflect.DeepEqual(messages, wantedMessages) {
		t.Errorf("Slices are diferent.\ncalc: %v\nwant: %v", messages, wantedMessages)
	}
	os.Args = oldsArgs
}

func TestAskForHelp(t *testing.T) {
	oldsArgs := os.Args
	os.Args = []string{"cmd", "-h"}
	wanted := true
	got := AskForHelp()
	if got != wanted {
		t.Errorf("Test AskForHelp() with presence result error, got %t wanted %t", got, wanted)
	}
	os.Args = oldsArgs
	os.Args = []string{"cmd", "help"}
	wanted = true
	got = AskForHelp()
	if got != wanted {
		t.Errorf("Test AskForHelp() with presence result error, got %t wanted %t", got, wanted)
	}
	os.Args = oldsArgs

}

func TestIsProfileServerArgPresent(t *testing.T) {
	oldsArgs := os.Args
	os.Args = []string{"cmd", "-profile=server"}
	wanted := true
	got := IsProfileServerArgPresent()
	if got != wanted {
		t.Errorf("Test IsProfileServerArgPresent() with presence result error, got %t wanted %t", got, wanted)
	}

	// restores previous args
	os.Args = oldsArgs

	// run again with no presence
	wanted = false
	got = IsProfileServerArgPresent()
	if got != wanted {
		t.Errorf("Test IsProfileServerArgPresent() wiout presence result error, got %t wanted %t", got, wanted)
	}
}
