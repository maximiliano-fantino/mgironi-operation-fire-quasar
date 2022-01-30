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
	"github.com/mgironi/operation-fire-quasar/model"
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
		{name: "test5", args: args{routerPath: "/topsecret/", rqFilename: "../_test/topSecret_test5_request.json"}, wantJSONFile: "../_test/topSecret_test5_response.json", wantStatusCode: http.StatusNotFound, wantError: true},
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
