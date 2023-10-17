package webclient

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"gylib/common/datatype"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/*var (
	HttpClient *http.Client
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     int
)*/

type Http_Client struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     int
	IsKeepActive        int
	Client              *http.Client
}

func NewHttpClient() *Http_Client {
	this := new(Http_Client)
	this.IdleConnTimeout = 90
	this.IsKeepActive = 1
	this.MaxIdleConnsPerHost = 1000
	this.MaxIdleConns = 1000
	this.Init_HTTPClient()
	return this
}

// init HTTPClient
func (this *Http_Client) SetMaxconns(id int) *Http_Client {
	this.MaxIdleConns = id
	return this
}

func (this *Http_Client) SetMaxperHost(id int) *Http_Client {
	this.MaxIdleConnsPerHost = id
	return this
}

func (this *Http_Client) SetTimeOut(id int) *Http_Client {
	this.IdleConnTimeout = id
	return this
}

func (this *Http_Client) Init_HTTPClient() {
	this.Client = this.CreateHTTPClient()

}

// createHTTPClient for connection re-use
func (this *Http_Client) CreateHTTPClient() *http.Client {
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

func (this *Http_Client) Web_Form_POST_Header(url_add string, data url.Values, header string) string {
	//return Web_Form_POST(url_add,data)
	//fmt.Println(data)
	//fmt.Println("start", strings.NewReader(data.Encode()))
	request, err := http.NewRequest("POST", url_add, strings.NewReader(data.Encode()))
	if err != nil {
		return ""
	}
	if this.IsKeepActive == 1 {
		request.Header.Set("Connection", "Keep-Alive")
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36")
	if header != "" {
		request.Header.Add("userinfo", header)
	}
	var resp *http.Response
	fmt.Println(request)
	resp, err = this.Client.Do(request)
	if err != nil {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("end", string(b))
	defer resp.Body.Close()

	if err != nil {
		return ""
	}
	return string(b)
}

func (this *Http_Client) Web_Form_GET_Header(url_add string, header string) string {
	//return Web_Form_POST(url_add,data)
	//fmt.Println(data)
	//fmt.Println("start", strings.NewReader(data.Encode()))
	request, err := http.NewRequest("GET", url_add, nil)
	if err != nil {
		return ""
	}
	if this.IsKeepActive == 1 {
		request.Header.Set("Connection", "Keep-Alive")
	}
	if header != "" {
		request.Header.Add("userinfo", header)
	}
	var resp *http.Response

	resp, err = this.Client.Do(request)
	if err != nil {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("end", string(b))
	defer resp.Body.Close()

	if err != nil {
		return ""
	}
	return string(b)
}

func (this *Http_Client) Web_Form_POST(url_add string, data url.Values) string {
	//s_data:=url.Values{}
	//for k,v:=range data{
	//	s_data.Set(k,datatype.Type2str(v))
	//}

	res, err := this.Client.PostForm(url_add, data)
	//设置http中header参数，可以再此添加cookie等值
	//res.Header.Add("User-Agent", "***")
	//res.Header.Add("http.socket.timeou", 5000)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	return string(body)
}

func (this *Http_Client) HttpGet(url_add string) string {
	/*client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        MaxIdleConns,
			MaxIdleConnsPerHost: MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Duration(IdleConnTimeout) * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			DisableCompression:  true,
		},

		Timeout: 20 * time.Second,
	}*/

	resp, err := this.Client.Get(url_add)
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

func (this *Http_Client) Http_post_json(url_add string, data string) string {
	body := bytes.NewBuffer([]byte(data))

	res, err := this.Client.Post(url_add, "application/json;charset=utf-8", body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(result)
}

func (this *Http_Client) Http_post_json_Head(url_add string, data string, header map[string]string) string {
	body := bytes.NewBuffer([]byte(data))
	request, err := http.NewRequest("POST", url_add, body)
	if err != nil {
		return ""
	}
	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	for key, val := range header {
		request.Header.Set(key, val)
	}

	var resp *http.Response

	resp, err = this.Client.Do(request)
	if err != nil {
		return ""
	}
	b, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("end", string(b))
	defer resp.Body.Close()

	if err != nil {
		return ""
	}
	return string(b)

}

func (this *Http_Client) Https_post_json(url_add string, r_data string) string {
	//s_data := url.Values{}
	//for k, v := range r_data {
	//	s_data.Set(k, datatype.Type2str(v))
	//}
	var resp *http.Response
	var err error
	var result []byte

	//b, err := json.Marshal(&r_data)
	body := bytes.NewBuffer([]byte(r_data))
	req, _ := http.NewRequest("POST", url_add, body)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Accept", "application/json")

	resp, err = this.Client.Do(req)
	//fmt.Println(resp,err)
	//启用cookie
	//client.Jar, _ = cookiejar.New(nil)
	//resp, err =client.Post(url_add, "application/json;charset=utf-8;", body)
	//resp, err = client.Post(url_add,"application/json", body)
	result, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err == nil {
		//fmt.Printf("%s\n", result)
		return string(result)
	}

	return ""

}

func (this *Http_Client) DoBytesPost(url string, data []byte) (string, error) {
	body := bytes.NewReader(data)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		//log.Println("http.NewRequest,[err=%s][url=%s]", err, url)
		return "", err
	}
	//request.Header.Set("Connection", "Keep-Alive")
	//request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var resp *http.Response

	resp, err = this.Client.Do(request)
	if err != nil {
		//log.Println("http.Do failed,[err=%s][url=%s]", err, url)
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {

		//log.Println("http.Do failed,[err=%s][url=%s]", err, url)
		return "", err
	}
	return string(b), err
}

func (this *Http_Client) DoJBytesPost_Header(url string, data []byte, header string) (string, error) {

	body := bytes.NewReader(data)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		//log.Println("http.NewRequest,[err=%s][url=%s]", err, url)
		return "", err
	}
	if header != "" {
		request.Header.Add("userinfo", header)
	}
	//request.Header.Set("Connection", "Keep-Alive")
	//request.Header.Set("Content-Type", "application/json")
	var resp *http.Response

	resp, err = this.Client.Do(request)
	if err != nil {
		//log.Println("http.Do failed,[err=%s][url=%s]", err, url)
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {

		//log.Println("http.Do failed,[err=%s][url=%s]", err, url)
		return "", err
	}
	return string(b), err
}

// 发送Xml请求
// 调用示例sendXmlRequest("POST", WXPAY_UNIFIEDORDER_URL, xmlString, config.TlsConfig, config.Timeout)
func (this Http_Client) SendXmlRequest(method, path string, xmlString []byte, tlsConfig *tls.Config, timeout time.Duration) (body []byte, err error) {
	inbody := bytes.NewReader(xmlString)
	//req, err := http.NewRequest(method, path, bytes.NewBufferString(xmlString))
	req, err := http.NewRequest(method, path, inbody)
	if err != nil {
		return
	}

	if timeout > 0 {
		this.Client.Timeout = timeout * time.Second
	}

	if tlsConfig != nil {
		this.Client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	resp, err := this.Client.Do(req)
	if err != nil {
		err = errors.New("request fail")
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func (this *Http_Client) API_Web_Form_POST_Header(url_add string, data url.Values, header string) []byte {
	//return Web_Form_POST(url_add,data)
	//fmt.Println(data)
	//fmt.Println("start", strings.NewReader(data.Encode()))
	rd_data := strings.NewReader(data.Encode())
	request, err := http.NewRequest("POST", url_add, rd_data)
	if err != nil {
		return nil
	}
	if this.IsKeepActive == 1 {
		request.Header.Set("Connection", "Keep-Alive")
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if header != "" {
		request.Header.Add("userinfo", header)
	}
	var resp *http.Response
	resp, err = this.Client.Do(request)
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

func (this *Http_Client) ApiPostSetHeader(url_add string, data []byte, header map[string]string) []byte {
	//return Web_Form_POST(url_add,data)
	//fmt.Println(data)
	//fmt.Println("start", strings.NewReader(data.Encode()))
	rd_data := bytes.NewReader(data)
	request, err := http.NewRequest("POST", url_add, rd_data)
	if err != nil {
		return nil
	}
	for key, val := range header {
		request.Header.Set(key, val)
	}
	var resp *http.Response
	resp, err = this.Client.Do(request)
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

func (this *Http_Client) ApiGetSetHeader(url_add string, header map[string]string) []byte {
	//return Web_Form_POST(url_add,data)
	//fmt.Println(data)
	//fmt.Println("start", strings.NewReader(data.Encode()))
	request, err := http.NewRequest("GET", url_add, nil)
	if err != nil {
		return nil
	}
	if this.IsKeepActive == 1 {
		request.Header.Set("Connection", "Keep-Alive")
	}
	request.Header.Set("Content-Type", "application/json")
	for key, val := range header {
		request.Header.Set(key, val)
	}
	var resp *http.Response
	resp, err = this.Client.Do(request)
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

func (this *Http_Client) API_Web_Form_GET_Header(url_add string, header string) []byte {
	//return Web_Form_POST(url_add,data)
	//fmt.Println(data)
	//fmt.Println("start", strings.NewReader(data.Encode()))
	request, err := http.NewRequest("GET", url_add, nil)
	if err != nil {
		return nil
	}
	if this.IsKeepActive == 1 {
		request.Header.Set("Connection", "Keep-Alive")
	}
	if header != "" {
		request.Header.Add("userinfo", header)
	}
	var resp *http.Response
	resp, err = this.Client.Do(request)
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

func (this *Http_Client) API_Web_Form_POST(url_add string, data url.Values) []byte {
	//s_data:=url.Values{}
	//for k,v:=range data{
	//	s_data.Set(k,datatype.Type2str(v))
	//}

	res, err := this.Client.PostForm(url_add, data)
	//设置http中header参数，可以再此添加cookie等值
	//res.Header.Add("User-Agent", "***")
	//res.Header.Add("http.socket.timeou", 5000)

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	return body
}

func (this *Http_Client) API_HttpGet(url_add string) []byte {
	resp, err := this.Client.Get(url_add)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil
		// handle error
	}
	return body
}

func (this *Http_Client) API_Http_post_json(url_add string, data string) []byte {
	body := bytes.NewBuffer([]byte(data))
	res, err := this.Client.Post(url_add, "application/json;charset=utf-8", body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {

		return nil
	}
	return result
}

func (this *Http_Client) API_Https_post(url_add string, data map[string]interface{}) []byte {
	s_data := url.Values{}
	for k, v := range data {
		s_data.Set(k, datatype.Type2str(v))
	}
	var resp *http.Response
	var err error
	var result []byte
	//启用cookie
	//client.Jar, _ = cookiejar.New(nil)
	resp, err = this.Client.PostForm(url_add, s_data)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	if result, err = ioutil.ReadAll(resp.Body); err == nil {
		//fmt.Printf("%s\n", data)
		return result
	}
	return nil

}

func (this *Http_Client) Https_post(url_add string, data map[string]interface{}) string {
	s_data := url.Values{}
	for k, v := range data {
		s_data.Set(k, datatype.Type2str(v))
	}
	var resp *http.Response
	var err error
	var result []byte

	//启用cookie
	//client.Jar, _ = cookiejar.New(nil)
	resp, err = this.Client.PostForm(url_add, s_data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer resp.Body.Close()
	if result, err = ioutil.ReadAll(resp.Body); err == nil {
		//fmt.Printf("%s\n", data)
		return string(result)
	}

	return ""

}

func (this *Http_Client) API_Https_post_json(url_add string, r_data string) []byte {
	//s_data := url.Values{}
	//for k, v := range r_data {
	//	s_data.Set(k, datatype.Type2str(v))
	//}
	var resp *http.Response
	var err error
	var result []byte

	//b, err := json.Marshal(&r_data)
	body := bytes.NewBuffer([]byte(r_data))
	req, _ := http.NewRequest("POST", url_add, body)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Accept", "application/json")
	resp, err = this.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//fmt.Println(resp,err)
	//启用cookie
	//client.Jar, _ = cookiejar.New(nil)
	//resp, err =client.Post(url_add, "application/json;charset=utf-8;", body)
	//resp, err = client.Post(url_add,"application/json", body)
	result, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err == nil {
		//fmt.Printf("%s\n", result)
		return result
	}

	return nil

}

func (this *Http_Client) API_DoBytesPost(url string, data []byte) ([]byte, error) {
	body := bytes.NewReader(data)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		//log.Println("http.NewRequest,[err=%s][url=%s]", err, url)
		return nil, err
	}
	//request.Header.Set("Connection", "Keep-Alive")
	//request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var resp *http.Response
	resp, err = this.Client.Do(request)
	if err != nil {
		//log.Println("http.Do failed,[err=%s][url=%s]", err, url)
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {

		//log.Println("http.Do failed,[err=%s][url=%s]", err, url)
		return nil, err
	}
	return b, err
}

func (this *Http_Client) API_DoJBytesPost_Header(url string, data []byte, header string) ([]byte, error) {

	body := bytes.NewReader(data)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		//log.Println("http.NewRequest,[err=%s][url=%s]", err, url)
		return nil, err
	}
	if header != "" {
		request.Header.Add("userinfo", header)
	}
	//request.Header.Set("Connection", "Keep-Alive")
	//request.Header.Set("Content-Type", "application/json")
	var resp *http.Response
	resp, err = this.Client.Do(request)
	if err != nil {
		//log.Println("http.Do failed,[err=%s][url=%s]", err, url)
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {

		//log.Println("http.Do failed,[err=%s][url=%s]", err, url)
		return nil, err
	}
	return b, err
}

// 发送Xml请求
// 调用示例sendXmlRequest("POST", WXPAY_UNIFIEDORDER_URL, xmlString, config.TlsConfig, config.Timeout)
func (this *Http_Client) API_SendXmlRequest(method, path string, xmlString []byte, tlsConfig *tls.Config, timeout time.Duration) (body []byte, err error) {
	inbody := bytes.NewReader(xmlString)

	//req, err := http.NewRequest(method, path, bytes.NewBufferString(xmlString))
	req, err := http.NewRequest(method, path, inbody)
	if err != nil {
		return
	}
	if timeout > 0 {
		this.Client.Timeout = timeout * time.Second
	}

	if tlsConfig != nil {
		this.Client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	resp, err := this.Client.Do(req)
	if err != nil {
		err = errors.New("request fail")
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	this.Init_HTTPClient()
	return
}
