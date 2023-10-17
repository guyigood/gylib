package resetful

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ResetFulClient struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     int
	IsKeepActive        int
	HeadData            map[string]string
	HttpClient          *http.Client
}

func NewResetFulClient() *ResetFulClient {
	this := new(ResetFulClient)
	this.IdleConnTimeout = 90
	this.IsKeepActive = 1
	this.MaxIdleConnsPerHost = 1000
	this.MaxIdleConns = 1000
	this.HttpClient = this.CreateHTTPClient()
	this.HeadData = make(map[string]string)
	return this
}

func (this *ResetFulClient) CreateHTTPClient() *http.Client {
	if this.IsKeepActive == 1 {
		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				//DisableKeepAlives: true,
				DialContext: (&net.Dialer{
					Timeout:   time.Duration(this.IdleConnTimeout) * time.Second,
					KeepAlive: time.Duration(this.IdleConnTimeout) * time.Second,
				}).DialContext,
				MaxIdleConns:        this.MaxIdleConns,
				MaxIdleConnsPerHost: this.MaxIdleConnsPerHost,
				IdleConnTimeout:     time.Duration(this.IdleConnTimeout) * time.Second,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
				DisableCompression:  true,
			},

			Timeout: time.Duration(this.IdleConnTimeout) * time.Second,
		}
		return client
	} else {
		client := &http.Client{
			Transport: &http.Transport{
				Proxy:             http.ProxyFromEnvironment,
				DisableKeepAlives: true,
				DialContext: (&net.Dialer{
					Timeout: time.Duration(this.IdleConnTimeout) * time.Second,
					//KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:        this.MaxIdleConns,
				MaxIdleConnsPerHost: this.MaxIdleConnsPerHost,
				IdleConnTimeout:     time.Duration(this.IdleConnTimeout) * time.Second,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
				DisableCompression:  true,
			},

			Timeout: time.Duration(this.IdleConnTimeout) * time.Second,
		}
		return client
	}
}

func (this *ResetFulClient) SetHeader(data map[string]string) *ResetFulClient {
	this.HeadData = make(map[string]string)
	for key, val := range data {
		this.HeadData[key] = val
	}
	return this
}

func (this *ResetFulClient) HttpGet(url_add string) string {
	resp, err := this.HttpClient.Get(url_add)
	if err != nil {
		fmt.Println(err)
		return "" // handle error
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		// handle error
		fmt.Println(err)
		return ""
	}
	return string(body)
}

func (this *ResetFulClient) PostJson(url_add string, data string) string {
	body := bytes.NewBuffer([]byte(data))
	request, err := http.NewRequest("POST", url_add, body)
	if err != nil {
		return ""
	}
	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	for key, val := range this.HeadData {
		request.Header.Set(key, val)
	}

	var resp *http.Response

	resp, err = this.HttpClient.Do(request)
	if err != nil {
		fmt.Println("post json", err)
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("end", string(b))
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("post josn get body", err)
		return ""
	}
	return string(b)
}

func (this *ResetFulClient) PostBase64Json(url_add string, data string) string {
	//body := bytes.NewBuffer([]byte(base64.StdEncoding.EncodeToString([]byte(data))))
	body := bytes.NewBuffer([]byte(data))
	request, err := http.NewRequest("POST", url_add, body)
	if err != nil {
		return ""
	}
	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	for key, val := range this.HeadData {
		request.Header.Set(key, val)
	}

	var resp *http.Response

	resp, err = this.HttpClient.Do(request)
	if err != nil {
		fmt.Println("post json", err)
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("end", string(b))
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("post josn get body", err)
		return ""
	}
	return string(b)
}

func (this *ResetFulClient) PostData(url_add string, data url.Values) string {
	request, err := http.NewRequest("POST", url_add, strings.NewReader(data.Encode()))
	if err != nil {
		return ""
	}
	if this.IsKeepActive == 1 {
		request.Header.Set("Connection", "Keep-Alive")
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	for key, val := range this.HeadData {
		request.Header.Add(key, val)
	}
	var resp *http.Response
	resp, err = this.HttpClient.Do(request)
	if err != nil {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return ""
	}
	return string(b)
}
