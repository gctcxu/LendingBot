package poloniex

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"lending/common"
	"net/http"
	"net/url"
)

var baseUrl string = "https://poloniex.com"

type PoloHttpClient struct {
	apiKey    string
	apiSecret string
}

func NewPoloHttpClient(apikey string, apiSecret string) *PoloHttpClient {
	return &PoloHttpClient{apikey, apiSecret}
}

func (client *PoloHttpClient) Call(method string, apiUrl string, params map[string]string) string {

	var req *http.Request

	var sign string

	httpClient := http.Client{}

	switch method {
	case "GET", "DELETE":
		req, _ = http.NewRequest(method, baseUrl+apiUrl, nil)
	default:
		values := url.Values{}
		for k, v := range params {
			values.Add(k, v)
		}

		params := values.Encode()

		req, _ = http.NewRequest(method, baseUrl+apiUrl, bytes.NewBuffer([]byte(params)))

		sign = common.SignWithHmac512ToHex(client.apiSecret, params)

		req.Header.Set("Sign", sign)
		req.Header.Set("Key", client.apiKey)
		req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	}

	httpRes, _ := httpClient.Do(req)

	res_byte, error := ioutil.ReadAll(httpRes.Body)

	if error != nil {
		fmt.Println(error.Error())
		return ""
	}

	return string(res_byte)
}

func (client *PoloHttpClient) Get(url string) string {
	return client.Call("GET", url, nil)
}

func (client *PoloHttpClient) Delete(url string) string {
	return client.Call("DELETE", url, nil)
}

func (client *PoloHttpClient) Post(url string, params map[string]string) string {
	return client.Call("POST", url, params)
}
