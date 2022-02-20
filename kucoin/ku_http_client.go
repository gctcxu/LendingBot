package kucoin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lending/common"
	"net/http"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
)

var base = "https://api.kucoin.com"

type KuHttpClient struct {
	apiKey        string
	apiSecret     string
	apiPassphrase string
}

func NewKuHttpClient(apiKey string, apiSecret string, apiPassphrase string) *KuHttpClient {
	return &KuHttpClient{
		apiKey,
		apiSecret,
		apiPassphrase,
	}
}

func (client *KuHttpClient) Call(method string, url string, params map[string]interface{}) string {

	var body_signature string
	//var passphrase_signature string
	var req *http.Request

	body_byte, _ := json.Marshal(params)

	now := time.Now()
	timestamp := now.Unix() * 1000

	switch method {
	case "GET", "DELETE":
		body_signature = strconv.FormatInt(timestamp, 10) + method + url
	default:
		body_signature = strconv.FormatInt(timestamp, 10) + method + url + string(body_byte)
	}

	body_signature = common.SignWithHmac256(client.apiSecret, body_signature)
	//passphrase_signature = SignWithHmac256(client.apiSecret, client.apiPassphrase)
	httpClient := http.Client{}

	switch method {
	case "GET", "DELETE":
		req, _ = http.NewRequest(method, base+url, nil)
	default:
		req, _ = http.NewRequest(method, base+url, bytes.NewBuffer(body_byte))
	}

	req.Header.Add("KC-API-KEY", client.apiKey)
	req.Header.Add("KC-API-SIGN", body_signature)
	req.Header.Add("KC-API-TIMESTAMP", strconv.FormatInt(timestamp, 10))
	req.Header.Add("KC-API-PASSPHRASE", client.apiPassphrase)
	req.Header.Add("Content-Type", "application/json")

	httpRes, _ := httpClient.Do(req)

	res_byte, error := ioutil.ReadAll(httpRes.Body)

	if error != nil {
		fmt.Println(error.Error())
		return ""
	}

	res_json, _ := simplejson.NewJson(res_byte)
	data_byte, _ := res_json.Get("data").MarshalJSON()

	return string(data_byte)
}

func (client *KuHttpClient) Get(url string) string {
	return client.Call("GET", url, nil)
}

func (client *KuHttpClient) Delete(url string) string {
	return client.Call("DELETE", url, nil)
}

func (client *KuHttpClient) Post(url string, params map[string]interface{}) string {
	return client.Call("POST", url, params)
}
