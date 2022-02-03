package store_test

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
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

func TestGetNewOperationUUID(t *testing.T) {
	uuidVal := store.GetNewOperationUUID()
	if uuidVal == "" {
		t.Error("Error TestGetNewOperationUUID(), UUID is empty")
	}
}

func TestGetDatasetByKey(t *testing.T) {
	conn := test.InitRedisMockConnection()
	key := "thesimplekey"
	wantedDataset := model.Dataset{
		Key:        key,
		Operation:  "anoperation",
		Satellites: []model.SatelliteInfoRequest{},
	}
	wantedMarshaled, _ := json.Marshal(wantedDataset)
	cmd := conn.Command("GET", key).Expect(wantedMarshaled)
	gotDataset := store.GetDatasetByKey(key)

	if conn.Stats(cmd) != 1 {
		t.Errorf("Error TestGetDatasetByKey(), redis command not used.")
	}

	if !reflect.DeepEqual(gotDataset, wantedDataset) {
		t.Errorf("Error TestGetDatasetByKey(), Dataset mismatch.\n---got:\n%v\n---wanted:\n%v", gotDataset, wantedDataset)
	}

}

func TestSaveNewDataset(t *testing.T) {
	conn := test.InitRedisMockConnection()
	operation := store.GetNewOperationUUID()
	dataValue := model.SatelliteInfoRequest{
		Name:     "kenobi",
		Distance: 100,
		Message:  []string{"is", "", "a", "msg"},
	}
	wantKey := fmt.Sprintf("%s:%s", operation, strings.Join(dataValue.Message, " "))
	wantValueStruct := model.Dataset{
		Satellites: []model.SatelliteInfoRequest{dataValue},
		Key:        wantKey,
		Operation:  operation,
	}
	wantValue, _ := json.Marshal(wantValueStruct)
	cmd := conn.Command("SET", wantKey, wantValue, "NX").Expect("OK")

	store.SaveNewDataset(operation, dataValue)

	if operation == "" {
		t.Fatalf("Error TestSaveNewDataset(), dataset not saved. values:%v", dataValue)
	}

	if conn.Stats(cmd) != 1 {
		t.Fatalf("Error TestSaveNewDataset(), redis command not used.")
	}
}

func TestUpdateDataset(t *testing.T) {
	conn := test.InitRedisMockConnection()
	operation := store.GetNewOperationUUID()
	previousValue := model.SatelliteInfoRequest{
		Name:     "kenobi",
		Distance: 100,
		Message:  []string{"is", "", "a", "msg"},
	}
	previousKey := fmt.Sprintf("%s:%s", operation, strings.Join(previousValue.Message, " "))
	previousValueStruct := model.Dataset{
		Satellites: []model.SatelliteInfoRequest{previousValue},
		Key:        previousKey,
		Operation:  operation,
	}
	previousValueMsl, _ := json.Marshal(previousValueStruct)
	cmdGET := conn.Command("GET", previousKey).Expect(previousValueMsl)
	consMsg := "is  a msg"
	wantedKey := fmt.Sprintf("%s:%s", operation, consMsg)
	dataValue := model.SatelliteInfoRequest{Name: "sato", Distance: 100, Message: []string{"is", "", "a", "msg"}}
	dataValueStruct := model.Dataset{Key: wantedKey, Operation: operation, Satellites: []model.SatelliteInfoRequest{previousValue, dataValue}}
	dataValueMsl, _ := json.Marshal(dataValueStruct)

	cmdSET := conn.Command("SET", wantedKey, dataValueMsl, "NX").Expect("OK")
	cmdDEL := conn.Command("DEL", previousKey).Expect(int64(1))

	saved := store.UpdateDataset(operation, consMsg, previousKey, dataValue)

	if !saved {
		t.Errorf("Error TestUpdateDataset(), dataset not updated. values:%v", dataValue)
	}
	if conn.Stats(cmdGET) != 2 {
		t.Errorf("Error TestUpdateDataset(), redis command GET not used.")
	}
	if conn.Stats(cmdSET) != 1 {
		t.Errorf("Error TestUpdateDataset(), redis command SET not used.")
	}
	if conn.Stats(cmdDEL) != 1 {
		t.Errorf("Error TestUpdateDataset(), redis command DEL not used.")
	}
}

func TestUpdateDatasetByOperation(t *testing.T) {
	conn := test.InitRedisMockConnection()
	operation := store.GetNewOperationUUID()
	previousValue := model.SatelliteInfoRequest{
		Name:     "kenobi",
		Distance: 100,
		Message:  []string{"is", "", "a", "msg"},
	}
	previousKey := fmt.Sprintf("%s:%s", operation, strings.Join(previousValue.Message, " "))
	previousValueStruct := model.Dataset{
		Satellites: []model.SatelliteInfoRequest{previousValue},
		Key:        previousKey,
		Operation:  operation,
	}
	previousValueMsl, _ := json.Marshal(previousValueStruct)
	cmdGET := conn.Command("GET", previousKey).Expect(previousValueMsl)
	//consMsg := "is  a msg"
	wantedKey := operation
	dataValue := model.SatelliteInfoRequest{Name: "sato", Distance: 100, Message: []string{"is", "", "a", "msg"}}
	dataValueStruct := model.Dataset{Key: wantedKey, Operation: operation, Satellites: []model.SatelliteInfoRequest{previousValue, dataValue}}
	dataValueMsl, _ := json.Marshal(dataValueStruct)

	cmdSET := conn.Command("SET", wantedKey, dataValueMsl, "NX").Expect("OK")
	cmdDEL := conn.Command("DEL", previousKey).Expect(int64(1))

	saved := store.UpdateDataset(operation, "", previousKey, dataValue)

	if !saved {
		t.Errorf("Error TestUpdateDatasetByOperation(), dataset not updated. values:%v", dataValue)
	}
	if conn.Stats(cmdGET) != 2 {
		t.Errorf("Error TestUpdateDatasetByOperation(), redis command GET not used.")
	}
	if conn.Stats(cmdSET) != 1 {
		t.Errorf("Error TestUpdateDatasetByOperation(), redis command SET not used.")
	}
	if conn.Stats(cmdDEL) != 1 {
		t.Errorf("Error TestUpdateDatasetByOperation(), redis command DEL not used.")
	}
}
func TestGetDatasetByMessage(t *testing.T) {
	conn := test.InitRedisMockConnection()
	dataValue := model.SatelliteInfoRequest{Name: "kenobi", Distance: 100, Message: []string{"es", "", "msg"}}
	operation := "123"
	key := operation + ":es  msg"
	want := model.Dataset{
		Key:        key,
		Operation:  operation,
		Satellites: []model.SatelliteInfoRequest{dataValue},
	}

	rslScan1 := make([]interface{}, 2)
	rslScan1[0] = "210"
	rslScan1[1] = []interface{}{}
	cmdSCAN1 := conn.Command("SCAN", "0", "MATCH", "*:es * msg").Expect(rslScan1)

	rslScan2 := make([]interface{}, 2)
	rslScan2[0] = "0"
	rslScan2[1] = []interface{}{key}
	cmdSCAN2 := conn.Command("SCAN", "210", "MATCH", "*:es * msg").Expect(rslScan2)

	wantMsl, _ := json.Marshal(want)
	cmdGET := conn.Command("GET", key).Expect(wantMsl)

	message := []string{"es", "", "msg"}
	got := store.GetDatasetByMessage(message)

	if conn.Stats(cmdSCAN1) != 1 {
		t.Errorf("Error TestGetDatasetByMessage(), redis command SCAN1 not used.")
	}

	if conn.Stats(cmdSCAN2) != 1 {
		t.Errorf("Error TestGetDatasetByMessage(), redis command SCAN2 not used.")
	}

	if conn.Stats(cmdGET) != 1 {
		t.Errorf("Error TestGetDatasetByMessage(), redis command GET not used.")
	}

	if got.Key == "" {
		t.Errorf("Error TestGetDatasetByMessage(), returns an empty dataset.")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error TestGetDatasetByMessage(), result mismatch.\n---got:\n%v\n---want:\n%v", got, want)
	}

}

func TestGetDatasetByMessageForcingFuzzy(t *testing.T) {
	conn := test.InitRedisMockConnection()
	dataValue := model.SatelliteInfoRequest{Name: "kenobi", Distance: 100, Message: []string{"es", "", "msg"}}
	operation := "123"
	key := operation + ":es  msg"
	want := model.Dataset{
		Key:        key,
		Operation:  operation,
		Satellites: []model.SatelliteInfoRequest{dataValue},
	}

	rslScan1 := make([]interface{}, 2)
	rslScan1[0] = "0"
	rslScan1[1] = []interface{}{}
	cmdSCAN1 := conn.Command("SCAN", "0", "MATCH", "*:es un msg").Expect(rslScan1)

	rslScan2 := make([]interface{}, 2)
	rslScan2[0] = "0"
	rslScan2[1] = []interface{}{key}
	cmdSCAN2 := conn.Command("SCAN", "0", "MATCH", "*").Expect(rslScan2)

	wantMsl, _ := json.Marshal(want)
	cmdGET := conn.Command("GET", key).Expect(wantMsl)

	message := []string{"es", "un", "msg"}
	got := store.GetDatasetByMessage(message)

	if conn.Stats(cmdSCAN1) != 1 {
		t.Errorf("Error TestGetDatasetByMessageForcingFuzzy(), redis command SCAN1 not used.")
	}

	if conn.Stats(cmdSCAN2) != 1 {
		t.Errorf("Error TestGetDatasetByMessageForcingFuzzy(), redis command SCAN2 not used.")
	}

	if conn.Stats(cmdGET) != 1 {
		t.Errorf("Error TestGetDatasetByMessageForcingFuzzy(), redis command GET not used.")
	}

	if got.Key == "" {
		t.Errorf("Error TestGetDatasetByMessageForcingFuzzy(), returns an empty dataset.")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error TestGetDatasetByMessageForcingFuzzy(), result mismatch.\n---got:\n%v\n---want:\n%v", got, want)
	}

}

func TestScanWithFuzzy(t *testing.T) {
	conn := test.InitRedisMockConnection()

	message := []string{"", "is", "message"}
	want := "123:this  message"
	similar := "751:this  mess"

	rslScan1 := make([]interface{}, 2)
	rslScan1[0] = "37"
	rslScan1[1] = []interface{}{want, "156:other msg"}
	cmdSCAN1 := conn.Command("SCAN", "0", "MATCH", "*").Expect(rslScan1)

	rslScan2 := make([]interface{}, 2)
	rslScan2[0] = "0"
	rslScan2[1] = []interface{}{similar, "1586:other msg"}
	cmdSCAN2 := conn.Command("SCAN", "37", "MATCH", "*").Expect(rslScan2)

	got := store.ScanWithFuzzy(message)

	if conn.Stats(cmdSCAN1) != 1 {
		t.Errorf("Error TestScanWithFuzzy(), redis command SCAN1 not used.")
	}

	if conn.Stats(cmdSCAN2) != 1 {
		t.Errorf("Error TestScanWithFuzzy(), redis command SCAN2 not used.")
	}

	if got != want {
		t.Errorf("Error TestScanWithFuzzy(), result mismatch. got: '%s', want:'%s'", got, want)
	}
}

func TestGetDatasetByOperation(t *testing.T) {
	conn := test.InitRedisMockConnection()

	dataValue := model.SatelliteInfoRequest{Name: "kenobi", Distance: 100, Message: []string{"es", "", "msg"}}
	operation := "123-456"
	key := operation + ":es  msg"
	want := model.Dataset{
		Key:        key,
		Operation:  operation,
		Satellites: []model.SatelliteInfoRequest{dataValue},
	}

	rslScan := make([]interface{}, 2)
	rslScan[0] = "0"
	rslScan[1] = []interface{}{key}
	cmdSCAN := conn.Command("SCAN", "0", "MATCH", "123-456:*").Expect(rslScan)

	wantMsl, _ := json.Marshal(want)
	cmdGET := conn.Command("GET", key).Expect(wantMsl)

	got := store.GetDatasetByOperation(operation)

	if conn.Stats(cmdSCAN) != 1 {
		t.Errorf("Error TestGetDatasetByOperation(), redis command SCAN not used.")
	}

	if conn.Stats(cmdGET) != 1 {
		t.Errorf("Error TestGetDatasetByOperation(), redis command GET not used.")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error TestGetDatasetByOperation(), result mismatch.\n---got:\n%v\n---want:\n%v", got, want)
	}
}
