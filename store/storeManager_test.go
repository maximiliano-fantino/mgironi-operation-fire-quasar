package store_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/mgironi/operation-fire-quasar/model"
	"github.com/mgironi/operation-fire-quasar/store"

	test "github.com/mgironi/operation-fire-quasar/_test"
)

func TestInitializeSatelitesInfo(t *testing.T) {
	test.CleanSatelitesInfoEnvs()

	store.InitializeSatelitesInfo()

	wanted := map[int]model.SateliteInfo{
		0: {Name: "kenobi", Location: model.Point{X: -500, Y: -200}},
		1: {Name: "skywalker", Location: model.Point{X: 100, Y: -100}},
		2: {Name: "sato", Location: model.Point{X: 500, Y: 100}},
	}

	if !reflect.DeepEqual(store.Satelites, wanted) {
		t.Errorf("Satelites default info don't match. Is %v, wanted %v", store.Satelites, wanted)
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
