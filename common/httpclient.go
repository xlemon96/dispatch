package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/navieboy/dispatch/model/api"
)

const (
	defaultConnectTimeout        int = 500
	defaultMaxIdleConns          int = 100
	defaultTimerInterval         int = 60
	defaultResponseHeaderTimeout int = 10000
	defaultRequestTotalTimeout   int = 20000
)

type HttpClient struct {
	client *http.Client
}

//默认HttpClient
func NewDefaultClient() *HttpClient {
	client := &HttpClient{}
	tr := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Duration(defaultConnectTimeout) * time.Millisecond,
		}).DialContext,
		MaxIdleConnsPerHost:   defaultMaxIdleConns,
		MaxIdleConns:          defaultMaxIdleConns,
		IdleConnTimeout:       time.Duration(defaultTimerInterval*2) * time.Second,
		ResponseHeaderTimeout: time.Duration(defaultResponseHeaderTimeout) * time.Millisecond,
	}
	client.client = &http.Client{Transport: tr, Timeout: time.Duration(defaultRequestTotalTimeout) * time.Millisecond}
	return client
}

func (c *HttpClient) CallHttpResponse(url string, action string, param interface{}, result interface{}) error {
	url = url + "?Action=" + action
	response := &api.CommonResponse{}
	response.Data = result
	if err := c.DoRequest("POST", url, nil, param, response); err != nil {
		return err
	}
	if response.Code != 0 {
		return errors.New(response.Message + response.Detail)
	}
	return nil
}

func (c *HttpClient) DoRequest(method, url string, header map[string]string, params, response interface{}) error {
	//序列化参数param
	input, err := json.Marshal(params)
	if err != nil {
		return err
	}
	var req *http.Request
	//初始化request
	if string(input) == "null" {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, bytes.NewReader(input))
	}
	if err != nil {
		return err
	}
	//设置头部
	for key := range header {
		req.Header.Set(key, header[key])
	}
	req.Header.Set("User-Agent", "JvirtClient")
	req.Header.Set("Accept", "application/json")
	//获取响应结果
	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode == 404 {
		return errors.New(fmt.Sprintf("ItemNotFound"))
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New(fmt.Sprintf("HTTP Error Code: %d", resp.StatusCode))
	}
	if len(body) == 0 {
		return nil
	}
	//反序列化
	if response != nil {
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.UseNumber()
		err = decoder.Decode(response)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *HttpClient) Post(url string, header map[string]string, params, response interface{}) error {
	return c.DoRequest("POST", url, header, params, response)
}

func (c *HttpClient) Get(url string, header map[string]string, params, response interface{}) error {
	return c.DoRequest("GET", url, header, params, response)
}