package nirvana

// https://medium.com/@marcus.olsson/writing-a-go-client-for-your-restful-api-c193a2f4998c

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

// NirvanaClient ...
type NirvanaClient struct {
	BaseURL    *url.URL
	authToken  string
	apiType    string
	appID      string
	appVersion string
	userAgent  string
	httpClient *http.Client
}

// NewNirvanaClient ...
func NewNirvanaClient(hc *http.Client) (*NirvanaClient, error) {
	if hc == nil {
		cookieJar, _ := cookiejar.New(nil)
		hc = &http.Client{
			Jar: cookieJar,
		}
	}

	url, err := url.Parse("https://api.nirvanahq.com/")
	if err != nil {
		return nil, err
	}

	ret := NirvanaClient{
		BaseURL:    url,
		appID:      "nirvana-sdk-golang",
		apiType:    "rest",
		appVersion: "1",
		httpClient: hc,
	}

	return &ret, err
}

// Close ...
func (c *NirvanaClient) Close() {
	c.httpClient.CloseIdleConnections()
}

// Authenticate ...
func (c *NirvanaClient) Authenticate(login string, passwordMd5 string) error {

	data := make(map[string]string)
	data["method"] = "auth.new"
	data["u"] = login
	data["p"] = passwordMd5

	req, err := c.newRequestURL(http.MethodPost, "/", data, nil, nil)
	if err != nil {
		return err
	}

	var response NirvanaResponse
	_, err = c.do(req, &response)
	if err != nil {
		return err
	}

	if err := response.ResultError(); err != nil {
		return err
	}

	token, found := response.AuthToken()
	if found {
		c.authToken = token
		return nil
	}
	return errors.New("Nirvan Client error: no authentification token")
}

// RetrieveSince ...
func (c *NirvanaClient) RetrieveSince(since int64) (*NirvanaResponse, error) {
	strSince := strconv.FormatInt(since, 10)
	params := map[string]string{
		"method": "everything",
		"since":  strSince,
	}
	req, err := c.newRequestJSON(http.MethodGet, "/", nil, params, nil)
	if err != nil {
		return nil, err
	}

	var response NirvanaResponse
	_, err = c.do(req, &response)
	if err != nil {
		return nil, err
	}

	if err := response.ResultError(); err != nil {
		return nil, err
	}

	return &response, err
}

func (c *NirvanaClient) newRequestURL(method string, path string, body map[string]string, extraParams map[string]string, extraHeaders map[string]string) (*http.Request, error) {

	data := url.Values{}
	if body != nil {
		for k, v := range body {
			data.Set(k, v)
		}
	}
	buf := data.Encode()

	if extraHeaders == nil {
		extraHeaders = make(map[string]string)
	}
	extraHeaders["Content-Type"] = "application/x-www-form-urlencoded"
	extraHeaders["Content-Length"] = strconv.Itoa(len(buf))

	return c.newRequest(method, path, strings.NewReader(buf), extraParams, extraHeaders)

}

func (c *NirvanaClient) newRequestJSON(method string, path string, body map[string]string, extraParams map[string]string, extraHeaders map[string]string) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	if extraHeaders == nil {
		extraHeaders = make(map[string]string)
		extraHeaders["Content-Type"] = "application/json"
	}
	return c.newRequest(method, path, buf, extraParams, extraHeaders)
}

func (c *NirvanaClient) newRequest(method string, path string, body io.Reader, extraParams map[string]string, extraHeaders map[string]string) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}
	c.addQueryParams(req, extraParams)
	c.addHeaders(req, extraHeaders)
	return req, nil
}

func (c *NirvanaClient) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("HTTP reponse error %d", resp.StatusCode)
	}

	// err = json.NewDecoder(resp.Body).Decode(v)
	// return resp, err

	// DEBUG - dump la réponse
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("nirvana_response.dump.json", bodyBytes, 0600)
	err = json.Unmarshal(bodyBytes, v)
	return resp, err
}

func (c *NirvanaClient) addQueryParams(req *http.Request, extraParams map[string]string) {
	params := map[string]string{
		"api":        c.apiType,
		"requestid":  uuid.New().String(),
		"clienttime": strconv.FormatInt(time.Now().Unix(), 10),
		"appid":      c.appID,
		"appversion": c.appVersion,
	}
	if len(c.authToken) > 0 {
		params["authtoken"] = c.authToken
	}
	if extraParams != nil {
		for k, v := range extraParams {
			params[k] = v
		}
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
}

func (c *NirvanaClient) addHeaders(req *http.Request, extraHeaders map[string]string) {
	headers := map[string]string{
		"Accept": "application/json",
	}
	if len(c.userAgent) > 0 {
		headers["User-Agent"] = c.userAgent
	}
	if extraHeaders != nil {
		for k, v := range extraHeaders {
			headers[k] = v
		}
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}
}
