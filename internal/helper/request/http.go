package request

import (
	"encoding/json"
	"fmt"
	"golang-api-starter/internal/config"
	logger "golang-api-starter/internal/helper/logger/zap_log"
	"io"
	"net/url"
	"slices"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

var cfg = config.Cfg

type HttpResp struct {
	StatusCode int
	ContenType string
	BodyBytes  []byte
}

var validMethods = []string{
	"GET",
	"POST",
	"PATCH",
	"PUT",
	"DELETE",
}

func HttpReq(reqMethod, reqUrl string, body *string, header map[string]string, retries *int) (*HttpResp, error) {
	if !slices.Contains(validMethods, reqMethod) {
		return nil, logger.Errorf("unrecognise req methods: %+v", reqMethod)
	}

	url, err := url.ParseRequestURI(reqUrl)
	if err != nil {
		return nil, logger.Errorf("failed to ParseRequestURI, err: %+v", err.Error())
	}

	retryClient := retryablehttp.NewClient()
	if retries != nil {
		retryClient.RetryMax = *retries
	} else {
		retryClient.RetryMax = 1
	}

	payload := &strings.Reader{}
	if body != nil {
		payload = strings.NewReader(*body)
	}

	req, err := retryablehttp.NewRequest(reqMethod, url.String(), payload)
	if err != nil {
		return nil, logger.Errorf("failed to NewRequest, err: %+v", err.Error())
	}

	if len(header) > 0 {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}

	resp, err := retryClient.Do(req)
	if err != nil {
		return nil, logger.Errorf("failed to %+v, err: %+v", fmt.Sprintf("%s %s", reqMethod, reqUrl), err.Error())
	}
	defer resp.Body.Close()

	logger.Debugf("resp code: %+v", resp.Status)
	// for k, v := range resp.Header {
	// 	logger.Debugf("k: %+v, v: %+v", k, v)
	// }

	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	logger.Debugf("content-type: %+v", contentType)

	bodyBytes, err := io.ReadAll(resp.Body)
	// logger.Debugf("resp body: %+v", string(bodyBytes))
	if err != nil {
		return nil, logger.Errorf("failed to ReadAll, err: %+v", err.Error())
	}

	return &HttpResp{StatusCode: resp.StatusCode, ContenType: contentType, BodyBytes: bodyBytes}, nil
}

func JsonToMap(bodyBytes []byte) (map[string]interface{}, error) {
	obj := map[string]interface{}{}
	arr := []map[string]interface{}{}
	result := map[string]interface{}{}

	// try handle normal json: {...}
	if err := json.Unmarshal(bodyBytes, &obj); err != nil {
		logger.Errorf("failed Unmarshal to object: %s", err.Error())
	} else {
		// result["data"] = obj
		return obj, nil
	}

	// try handle array of json: [{...}, {...}]
	if err := json.Unmarshal(bodyBytes, &arr); err != nil {
		logger.Errorf("failed Unmarshal to array object: %s", err.Error())
	} else {
		result["data"] = arr
		return result, nil
	}

	result["data"] = nil
	return result, logger.Errorf("failed to handle the response's JSON...")
}
