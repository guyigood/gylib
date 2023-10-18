package weixinmp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/guyigood/gylib/common/webclient"
	"io/ioutil"
)

// response from weixinmp
type response struct {
	// error fields
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	// token fields
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	// media fields
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt int64  `json:"created_at"`
	// ticket fields
	Ticket        string `json:"ticket"`
	ExpireSeconds int64  `json:"expire_seconds"`
}

func Http_post_json(url_add string, data string) (string, error) {
	client := webclient.NewHttpClient()
	body := bytes.NewBuffer([]byte(data))
	res, err := client.Client.Post(url_add, "application/json;charset=utf-8", body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return "", err
	}
	return string(result), nil
}

func post(url string, bodyType string, body *bytes.Buffer) (*response, error) {
	client := webclient.NewHttpClient()
	resp, err := client.Client.Post(url, bodyType, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rtn response
	if err := json.Unmarshal(data, &rtn); err != nil {
		return nil, err
	}
	if rtn.ErrCode != 0 {
		return nil, errors.New(fmt.Sprintf("%d %s", rtn.ErrCode, rtn.ErrMsg))
	}
	return &rtn, nil
}

func get(url string) (*response, error) {
	client := webclient.NewHttpClient()
	resp, err := client.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rtn response
	if err := json.Unmarshal(data, &rtn); err != nil {
		return nil, err
	}
	if rtn.ErrCode != 0 {
		return nil, errors.New(fmt.Sprintf("%d %s", rtn.ErrCode, rtn.ErrMsg))
	}
	return &rtn, nil
}
