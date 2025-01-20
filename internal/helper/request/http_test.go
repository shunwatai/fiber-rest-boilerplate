//go:build integration
// to run this test: go test -tags=integration -v ./internal/helper/request/... -run TestGet -count=1
package request

import (
	zlog "golang-api-starter/internal/helper/logger/zap_log"
	"golang-api-starter/internal/helper/utils"
	"testing"
)

type httpReqTest struct {
	name  string
	input req
	want  int
}

type req struct {
	method      string
	url         string
	body        *string
	header      map[string]string
	bearerToken *string
}

func TestGet(t *testing.T) {
	cfg.LoadEnvVariables()
	zlog.NewZlog()
	tests := []httpReqTest{
		// {
		// 	name:  "request to google",
		// 	input: req{method: "GET", url: "http://google.com"}, want: 200,
		// },
		// {
		// 	name:  "request to jsonplacehoder",
		// 	input: req{method: "GET", url: "https://jsonplaceholder.typicode.com/todos"}, want: 200,
		// },
		// {
		// 	name:  "request to jsonplacehoder",
		// 	input: req{method: "GET", url: "https://jsonplaceholder.typicode.com/todos/1"}, want: 200,
		// },
		{
			name:  "GET /ping",
			input: req{method: "GET", url: "http://localhost:7000/ping"}, want: 200,
		},
		{
			name:  "POST /auth/login",
			input: req{method: "POST", url: "http://localhost:7000/api/auth/login", body: utils.ToPtr(`{ "name": "admin@example.com", "password": "admin" }`), header: map[string]string{"Content-Type": "application/json"}}, want: 200,
		},
	}

	retries := 0

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := HttpReq(testCase.input.method, testCase.input.url, testCase.input.body, testCase.input.header, &retries)
			if err != nil {
				t.Fatal(err.Error())
			}

			// logger.Debugf("statusCode: %+v", got.StatusCode)
			// if strings.Contains(got.ContenType, "json") {
			// 	jsonMap, _ := JsonToMap(got.BodyBytes)
			// 	logger.Debugf("json val: %+v", jsonMap)
			// } else if strings.Contains(got.ContenType, "text") {
			// 	logger.Debugf("text val: %+v", string(got.BodyBytes))
			// } else {
			// 	logger.Debugf("contentType: %+v, val: %+v", got.ContenType, string(got.BodyBytes))
			// }

			if got.StatusCode != testCase.want {
				t.Errorf("got %q want %q", got, testCase.want)
			}
		})
	}
}
