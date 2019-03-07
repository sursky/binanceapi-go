// MIT License
//
// Copyright (c) 2019 Cranky Kernel
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package binanceapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

const API_ROOT = "https://api.binance.com"

type RestClient struct {
	apiKey    string
	apiSecret string
}

func NewRestClient() *RestClient {
	return &RestClient{}
}

func (c *RestClient) WithAuth(key string, secret string) *RestClient {
	c.apiKey = key
	c.apiSecret = secret
	return c
}

// Perform an unauthenticated GET request.
func (c *RestClient) Get(endpoint string, params map[string]interface{}) (*http.Response, error) {

	url := fmt.Sprintf("%s%s", API_ROOT, endpoint)
	queryString := ""

	if params == nil {
		params = map[string]interface{}{}
	}

	if params != nil {
		queryString = c.BuildQueryString(params)
		if queryString != "" {
			url = fmt.Sprintf("%s?%s", url, queryString)
		}
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(request)
}

func (c *RestClient) GetWithAuth(endpoint string, params map[string]interface{}) (*http.Response, error) {

	url := fmt.Sprintf("%s%s", API_ROOT, endpoint)
	queryString := ""

	if params == nil {
		params = map[string]interface{}{}
	}

	params["recvWindow"] = 5000
	params["timestamp"] = getTimeMillis()

	if params != nil {
		queryString = c.BuildQueryString(params)
		if queryString != "" {
			url = fmt.Sprintf("%s?%s", url, queryString)
		}
	}

	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write([]byte(queryString))
	signature := hex.EncodeToString(mac.Sum(nil))
	url = fmt.Sprintf("%s&signature=%s",
		url, signature)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("X-MBX-APIKEY", c.apiKey)

	return http.DefaultClient.Do(request)
}

func (c *RestClient) Post(endpoint string, params map[string]interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", API_ROOT, endpoint)
	queryString := ""

	if params == nil {
		params = map[string]interface{}{}
	}

	params["recvWindow"] = 5000
	params["timestamp"] = getTimeMillis()

	if params != nil {
		queryString = c.BuildQueryString(params)
		if queryString != "" {
			url = fmt.Sprintf("%s?%s", url, queryString)
		}
	}

	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write([]byte(queryString))
	signature := hex.EncodeToString(mac.Sum(nil))
	url = fmt.Sprintf("%s&signature=%s",
		url, signature)

	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("X-MBX-APIKEY", c.apiKey)

	return http.DefaultClient.Do(request)
}

// Send a POST request with only the API key and no other authentication.
func (c *RestClient) PostWithApiKey(endpoint string, params map[string]interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", API_ROOT, endpoint)
	queryString := ""

	if params == nil {
		params = map[string]interface{}{}
	}

	if params != nil {
		queryString = c.BuildQueryString(params)
		if queryString != "" {
			url = fmt.Sprintf("%s?%s", url, queryString)
		}
	}

	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	if c.apiKey != "" {
		request.Header.Add("X-MBX-APIKEY", c.apiKey)
	}

	return http.DefaultClient.Do(request)
}

func (c *RestClient) PutWithApiKey(path string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", API_ROOT, path)
	request, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-MBX-APIKEY", c.apiKey)

	return http.DefaultClient.Do(request)
}

func (c *RestClient) Delete(endpoint string, params map[string]interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", API_ROOT, endpoint)
	queryString := ""

	if params == nil {
		params = map[string]interface{}{}
	}

	params["recvWindow"] = 5000
	params["timestamp"] = getTimeMillis()

	if params != nil {
		queryString = c.BuildQueryString(params)
		if queryString != "" {
			url = fmt.Sprintf("%s?%s", url, queryString)
		}
	}

	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write([]byte(queryString))
	signature := hex.EncodeToString(mac.Sum(nil))
	url = fmt.Sprintf("%s&signature=%s",
		url, signature)

	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("X-MBX-APIKEY", c.apiKey)

	return http.DefaultClient.Do(request)
}

func (c *RestClient) BuildQueryString(params map[string]interface{}) string {
	queryString := ""

	keys := func() []string {
		keys := []string{}
		for key, _ := range params {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		return keys
	}()

	for _, key := range keys {
		if queryString != "" {
			queryString = fmt.Sprintf("%s&", queryString)
		}
		queryString = fmt.Sprintf("%s%s=%v", queryString, key, params[key])
	}

	return queryString
}

func (c *RestClient) GetAndDecode(endpoint string, params map[string]interface{}, response interface{}) error {
	httpResponse, err := c.Get(endpoint, params)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	return c.decodeBody(httpResponse, response)
}

func (c *RestClient) AuthGetAndDecode(endpoint string, params map[string]interface{}, response interface{}) error {
	httpResponse, err := c.Get(endpoint, params)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()
	return c.decodeBody(httpResponse, response)
}

func (c *RestClient) decodeBody(r *http.Response, v interface{}) error {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(v); err != nil {
		return err
	}
	return nil
}

type RestApiError struct {
	StatusCode int
	Body       []byte
}

func NewRestApiErrorFromResponse(r *http.Response) *RestApiError {
	body, _ := ioutil.ReadAll(r.Body)
	return &RestApiError{
		StatusCode: r.StatusCode,
		Body:       body,
	}
}

func (e *RestApiError) Error() string {
	return string(e.Body)
}
