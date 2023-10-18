package curl

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/guyigood/gylib/common/datatype"
	"github.com/guyigood/gylib/common/webclient"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type CurlGo struct {
	Url_add    string
	Header     map[string]string
	Body       map[string]interface{}
	_tlsConfig *tls.Config
	HttpClient *webclient.Http_Client
	CertPath   string
	KeyPath    string
	CAPath     string
}

func NewCurlGo() *CurlGo {
	this := new(CurlGo)
	this.Header = make(map[string]string)
	this.Body = make(map[string]interface{})
	this.HttpClient = webclient.NewHttpClient()
	return this
}

func (this *CurlGo) SetHead(data map[string]string) *CurlGo {
	this.Header = make(map[string]string)
	this.Header = data
	return this
}

func (this *CurlGo) SetBody(data map[string]interface{}) *CurlGo {
	this.Body = make(map[string]interface{})
	this.Body = data
	return this
}

func (this *CurlGo) Post() []byte {
	s_data := url.Values{}
	for k, v := range this.Body {
		s_data.Set(k, datatype.Type2str(v))
	}
	request, err := http.NewRequest("POST", this.Url_add, strings.NewReader(s_data.Encode()))
	if err != nil {
		return nil
	}
	request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range this.Header {
		request.Header.Add(k, v)
	}
	var resp *http.Response

	resp, err = this.HttpClient.Client.Do(request)
	if err != nil {
		return nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("end", string(b))
	defer resp.Body.Close()

	if err != nil {
		return nil
	}
	return b
}

func (this *CurlGo) Get() []byte {
	s_data := url.Values{}
	for k, v := range this.Body {
		s_data.Set(k, datatype.Type2str(v))
	}

	request, err := http.NewRequest("GET", this.Url_add, nil)
	if err != nil {
		return nil
	}
	request.Header.Set("Connection", "Keep-Alive")
	for k, v := range this.Header {
		request.Header.Add(k, v)
	}
	var resp *http.Response
	resp, err = this.HttpClient.Client.Do(request)
	if err != nil {
		return nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("end", string(b))
	defer resp.Body.Close()

	if err != nil {
		return nil
	}
	return b
}

func (this *CurlGo) PostJson() []byte {
	s_data, err := json.Marshal(&this.Body)
	if err != nil {
		return nil
	}

	request, err := http.NewRequest("POST", this.Url_add, bytes.NewBuffer(s_data))
	if err != nil {
		return nil
	}
	request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	for k, v := range this.Header {
		request.Header.Add(k, v)
	}
	var resp *http.Response
	resp, err = this.HttpClient.Client.Do(request)
	if err != nil {
		return nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("end", string(b))
	defer resp.Body.Close()

	if err != nil {
		return nil
	}
	return b
}

func (this *CurlGo) PostXML(xmlstr string) []byte {
	xml_data := []byte(xmlstr)

	request, err := http.NewRequest("POST", this.Url_add, bytes.NewBuffer(xml_data))
	if err != nil {
		return nil
	}
	for k, v := range this.Header {
		request.Header.Add(k, v)
	}
	var resp *http.Response
	resp, err = this.HttpClient.Client.Do(request)
	if err != nil {
		return nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("end", string(b))
	defer resp.Body.Close()

	if err != nil {
		return nil
	}
	return b
}

func (this *CurlGo) getTLSConfig() (*tls.Config, error) {
	if this._tlsConfig != nil {
		return this._tlsConfig, nil
	}
	// load cert
	cert, err := tls.LoadX509KeyPair(this.CertPath, this.KeyPath)
	if err != nil {
		fmt.Println("load wechat keys fail", err)
		return nil, err
	}
	// load root ca
	caData, err := ioutil.ReadFile(this.CAPath)
	if err != nil {
		fmt.Println("read wechat ca fail", err)
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	this._tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	return this._tlsConfig, nil
}

func (this *CurlGo) PostWXSSL(xmlstr string) []byte {
	xml_data := []byte(xmlstr)
	tlsConfig, err := this.getTLSConfig()
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	this.HttpClient.Client = &http.Client{Transport: tr}
	request, err := http.NewRequest("POST", this.Url_add, bytes.NewBuffer(xml_data))
	if err != nil {
		return nil
	}
	request.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	for k, v := range this.Header {
		request.Header.Add(k, v)
	}
	var resp *http.Response
	resp, err = this.HttpClient.Client.Do(request)
	if err != nil {
		return nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("end", string(b))
	defer resp.Body.Close()
	this.HttpClient.Init_HTTPClient()
	if err != nil {
		return nil
	}
	return b
}
