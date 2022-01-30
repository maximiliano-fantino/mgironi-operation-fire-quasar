package store_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	test "github.com/mgironi/operation-fire-quasar/_test"
	"github.com/mgironi/operation-fire-quasar/model"
	"github.com/mgironi/operation-fire-quasar/store"
)

func TestInitializeSatelitesInfo(t *testing.T) {
	test.CleanSatelitesInfoEnvs()

	store.InitializeSatelitesInfo()
	wantCount := 3
	want := []model.SateliteInfo{
		{Name: "kenobi", Location: model.Point{X: -500, Y: -200}},
		{Name: "skywalker", Location: model.Point{X: 100, Y: -100}},
		{Name: "sato", Location: model.Point{X: 500, Y: 100}},
	}
	gotCount := store.GetSatellitesInfoCount()
	if gotCount != wantCount {
		t.Errorf("Satelites default info size don't match. got: %d, wanted: %d.", gotCount, wantCount)
	}

	got := store.GetSatellitesInfo()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Satelites default info don't match.\n---Is:\n%v\n---wanted:\n%v\n", got, want)
	}
}

func TestGetKnownReferenceCoordinates(t *testing.T) {
	want := []model.Point{{X: -500, Y: -200}, {X: 100, Y: -100}, {X: 500, Y: 100}}
	got := store.GetKnownReferenceCoordinates()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Test GetKnownReferenceCoordinates() mismatch.\n---got:\n%v\n---want:\n%v\n", got, want)
	}
}

// Tests ParseSatelitesInfo using envs
func TestParseSatelitesInfoFromEnvs(t *testing.T) {
	envs := []string{store.SATELITE_KENOBI_ENV, store.SATELITE_SKYWALKER_ENV, store.SATELITE_SATO_ENV}
	names := []string{"kenobi", "skywalker", "sato"}
	coords := []model.Point{{X: 100, Y: -100}, {X: -200, Y: 200}, {X: 300, Y: 300}}
	envValueStrFormatPattern := "%s_%f,%f"
	for i, envKey := range envs {
		os.Setenv(envKey, fmt.Sprintf(envValueStrFormatPattern, names[i], coords[i].X, coords[i].Y))
	}
	satelitesInfo, hasErrors := store.ParseSatelitesInfoFromEnvs(envs)
	if hasErrors {
		t.Error("Error parsing env variables")
	}
	evaluateCoords := func(idx int, a, b model.Point) {
		if !test.AreFloatsEquals(a.X, b.X) {
			t.Errorf("Coordinate X mismatch on idx %d. Is %f, wanted %f", idx, a.X, b.X)
		}

		if !test.AreFloatsEquals(a.Y, b.Y) {
			t.Errorf("Coordinate Y mismatch on idx %d. Is %f, wanted %f", idx, a.X, b.X)
		}
	}
	for i, sat := range satelitesInfo {
		if sat.Name != names[i] {
			t.Errorf("Name mismatch. Is %s, wanted %s", sat.Name, names[i])
		}
		evaluateCoords(i, sat.Location, coords[i])
	}

	//clean envs
	test.CleanSatelitesInfoEnvs()
}

// Tests ConvertSateliteInfo
func TestConvertSateliteInfo(t *testing.T) {
	satInfo, err := store.ConvertSateliteInfo("kenobi_340.21,637.84")
	wantedName := "kenobi"
	wantedCoordX := float64(340.21)
	wantedCoordY := float64(637.84)

	if err != nil {
		t.Error(err)
	}
	if satInfo.Name != wantedName {
		t.Errorf("Name mismatch. Is %s, wanted %s", satInfo.Name, wantedName)
	}
	if !test.AreFloatsEquals(satInfo.Location.X, wantedCoordX) {
		t.Errorf("Location X mismatch. Is %f, wanted %f, ", satInfo.Location.X, wantedCoordX)
	}
	if !test.AreFloatsEquals(satInfo.Location.Y, wantedCoordY) {
		t.Errorf("Location Y mismatch. Is %f, wanted %f, ", satInfo.Location.Y, wantedCoordY)
	}
}

func TestGetSatelliteInfoIndex(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		wantIndex int
	}{
		{name: "test1", args: args{name: ""}, wantIndex: -1},
		{name: "test2", args: args{name: "any"}, wantIndex: -1},
		{name: "test3", args: args{name: "kenobi"}, wantIndex: 0},
		{name: "test4", args: args{name: "skywalker"}, wantIndex: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIndex := store.GetSatelliteInfoIndex(tt.args.name); gotIndex != tt.wantIndex {
				t.Errorf("GetSatelliteInfoIndex() = %v, want %v", gotIndex, tt.wantIndex)
			}
		})
	}
}
