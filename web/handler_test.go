package web_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	test "github.com/mgironi/operation-fire-quasar/_test"
	"github.com/mgironi/operation-fire-quasar/model"
	"github.com/mgironi/operation-fire-quasar/store"
	"github.com/mgironi/operation-fire-quasar/web"
)

func TestPingHandler(t *testing.T) {
	rPath := "/ping"
	router := gin.Default()

	router.GET(rPath, web.PingHandler)
	req, _ := http.NewRequest("GET", rPath, strings.NewReader(""))
	w := httptest.NewRecorder()
	wantedStatusCode := http.StatusOK
	wantedBodyStr := "echo"
	router.ServeHTTP(w, req)

	if w.Code != wantedStatusCode {
		t.Errorf("HTTP resposne status code mismatch got %d, want %d.", w.Code, wantedStatusCode)
	}

	if w.Body.String() != wantedBodyStr {
		t.Errorf("HTTP resposne body mismatch got '%s', want '%s'.", w.Body.String(), wantedBodyStr)
	}
}

func TestTopSecretHandler(t *testing.T) {
	type args struct {
		routerPath string
		rqFilename string
	}
	tests := []struct {
		name           string
		args           args
		wantJSONFile   string
		wantStatusCode int
		wantError      bool
	}{
		{name: "test1", args: args{routerPath: "/topsecret/", rqFilename: "../_test/topSecret_test1_request.json"}, wantJSONFile: "../_test/topSecret_test1_response.json", wantStatusCode: http.StatusOK, wantError: false},
		{name: "test2", args: args{routerPath: "/topsecret/", rqFilename: "../_test/topSecret_test2_request.json"}, wantJSONFile: "../_test/topSecret_test2_response.json", wantStatusCode: http.StatusNotFound, wantError: true},
		{name: "test3", args: args{routerPath: "/topsecret/", rqFilename: "../_test/topSecret_test3_request.json"}, wantJSONFile: "../_test/topSecret_test3_response.json", wantStatusCode: http.StatusNotFound, wantError: true},
		{name: "test4", args: args{routerPath: "/topsecret/", rqFilename: "../_test/topSecret_test4_request.json"}, wantJSONFile: "../_test/topSecret_test4_response.json", wantStatusCode: http.StatusNotFound, wantError: true},
		{name: "test5", args: args{routerPath: "/topsecret/", rqFilename: "../_test/topSecret_test5_request.json"}, wantJSONFile: "../_test/topSecret_test5_response.json", wantStatusCode: http.StatusBadRequest, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()

			// reads JSON body parameter from file
			jsonData, rErr := ioutil.ReadFile(tt.args.rqFilename)
			if rErr != nil {
				t.Fatalf("Error reading file %s.\n\ttrace: %s", tt.args.rqFilename, rErr.Error())
			}

			// loads POST URL and test handler
			router.POST(tt.args.routerPath, web.TopSecretHandler)
			// gets request
			request, _ := http.NewRequest("POST", tt.args.routerPath, bytes.NewReader(jsonData))

			// loads want response
			wantRsp := httptest.NewRecorder()

			// reads JSON response file
			wantJSON, rErr := ioutil.ReadFile(tt.wantJSONFile)
			if rErr != nil {
				t.Fatalf("Error reading file %s.\n\ttrace: %s", tt.wantJSONFile, rErr.Error())
			}

			// run http call
			router.ServeHTTP(wantRsp, request)

			if wantRsp.Code != tt.wantStatusCode {
				t.Errorf("HTTP response status code mismatch got %d, want %d.", wantRsp.Code, tt.wantStatusCode)
			}

			if !tt.wantError {
				// unmarshal JSON response
				var want model.TopSecretResponse
				wantedUMErr := json.Unmarshal(wantJSON, &want)
				if wantedUMErr != nil {
					t.Fatalf("Wanted value with format error.\n----wanted:\n%s\n%s", wantJSON, wantedUMErr.Error())
				}

				var got model.TopSecretResponse
				gotUMErr := json.Unmarshal(wantRsp.Body.Bytes(), &got)
				if gotUMErr != nil {
					t.Errorf("HTTP response body format error.\n----got:\n%s\n%s", wantRsp.Body.String(), gotUMErr.Error())
				}

				if !reflect.DeepEqual(got, want) {
					t.Errorf("HTTP response body mismatch...\n----got:\n%v\n----want:\n%v\n\n", got, want)
				}

			} else {
				// unmarshal JSON response
				var want model.ErrorResponse
				wantedUMErr := json.Unmarshal(wantJSON, &want)
				if wantedUMErr != nil {
					t.Fatalf("Wanted value with format error.\n----wanted:\n%s\n%s", wantJSON, wantedUMErr.Error())
				}

				var got model.ErrorResponse
				gotUMErr := json.Unmarshal(wantRsp.Body.Bytes(), &got)
				if gotUMErr != nil {
					t.Errorf("HTTP response body format error.\n----got:\n%s\n%s", wantRsp.Body.String(), gotUMErr.Error())
				}

				if !reflect.DeepEqual(got, want) {
					t.Errorf("HTTP response body mismatch...\n----got:\n%v\n----want:\n%v\n\n", got, want)
				}
			}

		})
	}
}

type tssArgs struct {
	routerPath string
	url        string
	rqFilename string
	handler    gin.HandlerFunc
	method     string
}

type tssTest struct {
	name           string
	args           tssArgs
	wantJSONFile   string
	got            interface{}
	want           interface{}
	wantStatusCode int
	wantError      bool
}

func TestTopSecretSplitHandler(t *testing.T) {
	conn := test.InitRedisMockConnection()

	basepathPOST := "/topsecret_split/"
	basepathGET := "/topsecret_split/results/"
	baseTestName := "test1"

	tPOST := tssTest{
		name: "",
		args: tssArgs{
			routerPath: basepathPOST,
			url:        basepathPOST,
			rqFilename: "",
			handler:    web.TopSecretSplitPOSTHandler,
			method:     http.MethodPost,
		},
		wantJSONFile:   "",
		got:            model.TopSecretSplitPOSTResponse{},
		want:           model.TopSecretSplitPOSTResponse{},
		wantStatusCode: http.StatusOK,
		wantError:      false,
	}

	tGET := tssTest{
		name: "",
		args: tssArgs{
			routerPath: "",
			url:        "",
			rqFilename: "",
			handler:    web.TopSecretSplitGETHandler,
			method:     http.MethodGet,
		},
		wantJSONFile:   "",
		got:            model.ErrorResponse{},
		want:           model.ErrorResponse{},
		wantStatusCode: http.StatusNotFound,
		wantError:      true,
	}

	dataValue := model.SatelliteInfoRequest{Name: "kenobi", Distance: 500, Message: []string{"este", "", "", "mensaje", ""}}
	operation := "123"
	store.GetNewOperationUUID = func() string {
		return operation
	}
	wantKey := operation + ":este   mensaje "
	want := model.Dataset{
		Key:        wantKey,
		Operation:  operation,
		Satellites: []model.SatelliteInfoRequest{dataValue},
	}

	rslScan := make([]interface{}, 2)
	rslScan[0] = "0"
	rslScan[1] = []interface{}{}
	cmdSCAN1 := conn.Command("SCAN", "0", "MATCH", "*:este * * mensaje *").Expect(rslScan)
	cmdSCAN2 := conn.Command("SCAN", "0", "MATCH", "*").Expect(rslScan)

	wantMsl, _ := json.Marshal(want)
	cmdSET := conn.Command("SET", wantKey, wantMsl, "NX").Expect("OK")

	tPOST.name = baseTestName + "-POST1"
	tPOST.args.rqFilename = "../_test/topSecretSplit_test1-POST1_request.json"
	gotP1 := model.TopSecretSplitPOSTResponse{}
	wantP1 := model.TopSecretSplitPOSTResponse{}
	runAsIt(tPOST, &gotP1, &wantP1, t)
	if gotP1.Operation == "" {
		t.Fatalf("Error in test %s. Operation token is empty.", tPOST.name)
	}

	if conn.Stats(cmdSCAN1) != 1 && conn.Stats(cmdSCAN2) != 1 {
		t.Errorf("Error TestTopSecretSplitHandler(), redis command SCAN not used.")
	}

	if conn.Stats(cmdSET) != 1 {
		t.Errorf("Error TestTopSecretSplitHandler(), redis command SET not used.")
	}

	cmdGET1GET := conn.Command("GET", operation).Expect("")

	tGET.name = baseTestName + "-GET1"
	tGET.args.routerPath = basepathGET + ":operation"
	tGET.args.url = basepathGET + gotP1.Operation
	tGET.wantJSONFile = "../_test/topSecretSplit_test1-GET1_response.json"
	gotG1 := model.ErrorResponse{}
	wantG1 := model.ErrorResponse{}
	runAsIt(tGET, &gotG1, &wantG1, t)

	if conn.Stats(cmdGET1GET) != 1 {
		t.Errorf("Error TestTopSecretSplitHandler(), redis command GET not used.")
	}

	rslScanPOST2SCAN1 := make([]interface{}, 2)
	rslScanPOST2SCAN1[0] = "36"
	rslScanPOST2SCAN1[1] = []interface{}{}
	cmdPOST2SCAN1 := conn.Command("SCAN", "0", "MATCH", "*:* es * * secreto").Expect(rslScanPOST2SCAN1)

	rslScanPOST2SCAN2 := make([]interface{}, 2)
	rslScanPOST2SCAN2[0] = "0"
	rslScanPOST2SCAN2[1] = []interface{}{wantKey}
	cmdPOST2SCAN2 := conn.Command("SCAN", "36", "MATCH", "*:* es * * secreto").Expect(rslScanPOST2SCAN2)

	cmdPOST2GET := conn.Command("GET", wantKey).Expect(wantMsl)
	dtValuePOST2 := model.SatelliteInfoRequest{Name: "skywalker", Distance: 424.26, Message: []string{"", "es", "", "", "secreto"}}
	wantkeyPOST2 := operation + ":este es  mensaje secreto"
	want.Key = wantkeyPOST2
	want.Satellites = append(want.Satellites, dtValuePOST2)
	wantMsl, _ = json.Marshal(want)
	cmdPOST2SET := conn.Command("SET", wantkeyPOST2, wantMsl, "NX").Expect("OK")

	cmdPOST2DEL := conn.Command("DEL", wantKey).Expect(1)

	tPOST.name = baseTestName + "-POST2"
	tPOST.args.rqFilename = "../_test/topSecretSplit_test1-POST2_request.json"
	gotP2 := model.TopSecretSplitPOSTResponse{}
	wantP2 := model.TopSecretSplitPOSTResponse{}
	runAsIt(tPOST, &gotP2, &wantP2, t)
	if gotP2.Operation == "" {
		t.Fatalf("Error in test %s. Operation token is empty.", tPOST.name)
	}

	if conn.Stats(cmdPOST2SCAN1) != 1 && conn.Stats(cmdPOST2SCAN2) != 1 {
		t.Errorf("Error TestTopSecretSplitHandler(), redis command SCAN not used.")
	}

	if conn.Stats(cmdPOST2GET) != 3 {
		t.Errorf("Error TestTopSecretSplitHandler(), redis command GET not used.")
	}

	if conn.Stats(cmdPOST2SET) != 1 {
		t.Errorf("Error TestTopSecretSplitHandler(), redis command SET not used.")
	}

	if conn.Stats(cmdPOST2DEL) != 1 {
		t.Errorf("Error TestTopSecretSplitHandler(), redis command DEL not used.")
	}

	//cmdGET2GET := conn.Command("GET", operation).Expect("")

	tGET.name = baseTestName + "-GET2"
	tGET.args.routerPath = basepathGET + ":operation"
	tGET.args.url = basepathGET + gotP2.Operation
	tGET.wantJSONFile = "../_test/topSecretSplit_test1-GET2_response.json"
	gotG2 := model.ErrorResponse{}
	wantG2 := model.ErrorResponse{}
	runAsIt(tGET, &gotG2, &wantG2, t)

	if conn.Stats(cmdGET1GET) != 2 {
		t.Errorf("Error TestTopSecretSplitHandler(), redis command GET not used.")
	}

	rslScanPOST3SCAN1 := make([]interface{}, 2)
	rslScanPOST3SCAN1[0] = "0"
	rslScanPOST3SCAN1[1] = []interface{}{}
	cmdPOST3SCAN1 := conn.Command("SCAN", "0", "MATCH", "*:este * un * *").Expect(rslScanPOST3SCAN1)

	rslScanPOST3SCAN2 := make([]interface{}, 2)
	rslScanPOST3SCAN2[0] = "0"
	rslScanPOST3SCAN2[1] = []interface{}{wantkeyPOST2}
	cmdPOST3SCAN2 := conn.Command("SCAN", "0", "MATCH", "*:este * un * *").Expect(rslScanPOST3SCAN2)

	cmdPOST3GET := conn.Command("GET", wantkeyPOST2).Expect(wantMsl)
	dtValuePOST3 := model.SatelliteInfoRequest{Name: "sato", Distance: 707.10, Message: []string{"este", "", "un", "", ""}}
	wantkeyPOST3 := operation
	want.Key = wantkeyPOST3
	want.Satellites = append(want.Satellites, dtValuePOST3)
	wantMsl, _ = json.Marshal(want)
	cmdPOST3SET := conn.Command("SET", wantkeyPOST3, wantMsl, "NX").Expect("OK")

	cmdPOST3DEL := conn.Command("DEL", wantkeyPOST2).Expect(1)

	tPOST.name = baseTestName + "-POST3"
	tPOST.args.rqFilename = "../_test/topSecretSplit_test1-POST3_request.json"
	gotP3 := model.TopSecretSplitPOSTResponse{}
	wantP3 := model.TopSecretSplitPOSTResponse{}
	runAsIt(tPOST, &gotP3, &wantP3, t)
	if gotP3.Operation == "" {
		t.Fatalf("Error in test %s. Operation token is empty.", tPOST.name)
	}

	if conn.Stats(cmdPOST3SCAN1) != 1 && conn.Stats(cmdPOST3SCAN2) != 1 {
		t.Errorf("Error TestTopSecretSplitHandler(), redis command SCAN not used.")
	}

	if conn.Stats(cmdPOST3GET) != 3 {
		t.Errorf("Error TestTopSecretSplitHandler(), redis command GET not used.")
	}

	if conn.Stats(cmdPOST3SET) != 1 {
		t.Errorf("Error TestTopSecretSplitHandler(), redis command SET not used.")
	}

	if conn.Stats(cmdPOST3DEL) != 1 {
		t.Errorf("Error TestTopSecretSplitHandler(), redis command DEL not used.")
	}

	cmdGET1GET = conn.Command("GET", wantkeyPOST3).Expect(wantMsl)

	tGET.name = baseTestName + "-GET3"
	tGET.args.routerPath = basepathGET + ":operation"
	tGET.args.url = basepathGET + gotP3.Operation
	tGET.wantJSONFile = "../_test/topSecretSplit_test1-GET3_response.json"
	tGET.got = model.TopSecretResponse{}
	tGET.want = model.TopSecretResponse{}
	tGET.wantStatusCode = http.StatusOK
	tGET.wantError = false
	gotG3 := model.TopSecretResponse{}
	wantG3 := model.TopSecretResponse{}
	runAsIt(tGET, &gotG3, &wantG3, t)

	if conn.Stats(cmdGET1GET) != 3 {
		t.Errorf("Error TestTopSecretSplitHandler(), redis command GET not used.")
	}
}

func readJSONFile(filename string, t *testing.T) (jsonData []byte) {
	if filename == "" {
		return jsonData
	}
	jsonData, rErr := ioutil.ReadFile(filename)
	if rErr != nil {
		t.Fatalf("Error reading file %s.\n\ttrace: %s", filename, rErr.Error())
	}
	return jsonData
}

func unmarshalJSONWithFatal(operation string, jsonData []byte, v interface{}, t *testing.T) {
	jsonUMErr := json.Unmarshal(jsonData, &v)
	if jsonUMErr != nil {
		t.Fatalf("Error in %s, JSON to unamrshal with format error.\n----json:\n%s\n%s", operation, jsonData, jsonUMErr.Error())
	}
}

func unmarshalJSONWithError(operation string, jsonData []byte, v interface{}, t *testing.T) {
	jsonUMErr := json.Unmarshal(jsonData, &v)
	if jsonUMErr != nil {
		t.Errorf("Error in %s, JSON to unamrshal with format error.\n----json:\n%s\n%s", operation, jsonData, jsonUMErr.Error())
	}
}

func compareValuesWithError(operation string, got int, want int, t *testing.T) {
	if got != want {
		t.Errorf("%s mismatch got %d, want %d.", operation, got, want)
	}
}

func compareResponsesByStructure(operation string, got interface{}, want interface{}, t *testing.T) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s mismatch...\n----got:\n%v\n----want:\n%v\n\n", operation, got, want)
	}
}

func runAsIt(tt tssTest, got interface{}, want interface{}, t *testing.T) {
	t.Logf("Test case :\t%s", tt.name)
	router := gin.Default()

	// reads JSON body parameter from file
	jsonData := readJSONFile(tt.args.rqFilename, t)

	// gets request
	var request *http.Request
	var rqErr error
	if tt.args.method == http.MethodPost {
		// loads POST URL and test handler
		router.POST(tt.args.routerPath, tt.args.handler)

		request, rqErr = http.NewRequest(tt.args.method, tt.args.routerPath, bytes.NewReader(jsonData))
	} else {
		// loads GET URL and test handler
		router.GET(tt.args.routerPath, tt.args.handler)

		request, rqErr = http.NewRequest(tt.args.method, tt.args.url, strings.NewReader(""))
	}
	if rqErr != nil {
		t.Fatalf("Error loading new request. Trace: \n%s", rqErr)
	}

	// loads want response
	gotRsp := httptest.NewRecorder()

	// run http call
	router.ServeHTTP(gotRsp, request)

	// compares http response status code
	compareValuesWithError("HTTP response status code", gotRsp.Code, tt.wantStatusCode, t)

	// unmarshal JSON got response
	unmarshalJSONWithError("Got response", gotRsp.Body.Bytes(), &got, t)

	if tt.wantJSONFile != "" {
		// reads JSON response file
		wantJSON := readJSONFile(tt.wantJSONFile, t)

		// unmarshal JSON want response
		unmarshalJSONWithFatal("Wanted response", wantJSON, &want, t)

		// compare responses
		compareResponsesByStructure("HTTP response body", got, want, t)
	}

}
