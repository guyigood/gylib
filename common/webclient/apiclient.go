package webclient

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"gylib/common/datatype"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func API_Web_Form_POST_Header(url_add string, data url.Values, header string) []byte {
	//return Web_Form_POST(url_add,data)
	fmt.Println(data)
	//fmt.Println("start", strings.NewReader(data.Encode()))
	request, err := http.NewRequest("POST", url_add, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if header != "" {
		request.Header.Add("userinfo", header)
	}
	var resp *http.Response
	resp, err = G_Http_client.Client.Do(request)
	fmt.Println(request)
	fmt.Println(resp)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("end", string(b))
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}

func API_Web_Form_GET_Header(url_add string, header string) []byte {
	//return Web_Form_POST(url_add,data)
	//fmt.Println(data)
	//fmt.Println("start", strings.NewReader(data.Encode()))
	request, err := http.NewRequest("GET", url_add, nil)
	if err != nil {
		return nil
	}
	request.Header.Set("Connection", "Keep-Alive")
	if header != "" {
		request.Header.Add("userinfo", header)
	}
	var resp *http.Response
	resp, err = G_Http_client.Client.Do(request)
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

func API_Web_Form_POST(url_add string, data url.Values) []byte {
	//s_data:=url.Values{}
	//for k,v:=range data{
	//	s_data.Set(k,datatype.Type2str(v))
	//}
	res, err := G_Http_client.Client.PostForm(url_add, data)
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

func API_HttpGet(url_add string) []byte {
	resp, err := G_Http_client.Client.Get(url_add)
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

func API_Http_post_json(url_add string, data string) []byte {
	body := bytes.NewBuffer([]byte(data))
	res, err := G_Http_client.Client.Post(url_add, "application/json;charset=utf-8", body)
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

func API_Https_post(url_add string, data map[string]interface{}) []byte {
	s_data := url.Values{}
	for k, v := range data {
		s_data.Set(k, datatype.Type2str(v))
	}
	var resp *http.Response
	var err error
	var result []byte

	//启用cookie
	//client.Jar, _ = cookiejar.New(nil)
	resp, err = G_Http_client.Client.PostForm(url_add, s_data)
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

func API_Https_post_json(url_add string, r_data string) []byte {
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
	resp, err = G_Http_client.Client.Do(req)
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

func API_DoBytesPost(url string, data []byte) ([]byte, error) {
	body := bytes.NewReader(data)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		//log.Println("http.NewRequest,[err=%s][url=%s]", err, url)
		return nil, err
	}
	//request.Header.Set("Connection", "Keep-Alive")
	//request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var resp *http.Response
	resp, err = G_Http_client.Client.Do(request)
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

func API_DoJBytesPost_Header(url string, data []byte, header string) ([]byte, error) {

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
	resp, err = G_Http_client.Client.Do(request)
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
func API_SendXmlRequest(method, path string, xmlString []byte, tlsConfig *tls.Config, timeout time.Duration) (body []byte, err error) {
	inbody := bytes.NewReader(xmlString)
	//req, err := http.NewRequest(method, path, bytes.NewBufferString(xmlString))
	req, err := http.NewRequest(method, path, inbody)
	if err != nil {
		return
	}

	if timeout > 0 {
		G_Http_client.Client.Timeout = timeout * time.Second
	}

	if tlsConfig != nil {
		G_Http_client.Client.Transport = &http.Transport{TLSClientConfig: tlsConfig}
	}

	resp, err := G_Http_client.Client.Do(req)
	if err != nil {
		err = errors.New("request fail")
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}
