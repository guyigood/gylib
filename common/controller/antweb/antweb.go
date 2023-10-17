package antweb

import (
	"encoding/json"
	"net/http"
)

type Controller struct {
	W            http.ResponseWriter
	R            *http.Request
	Data         map[string]interface{}
	JsonData     map[string]interface{}
	Postdata     map[string]interface{}
	Body         []byte
	Site_name    string
	Masterdb     string
	Access_token string
	Tokenname    string
	Jsonmsg      Json_msg
	Err_status   int
}

type Json_msg struct {
	Status int         `json:"status"`
	Msg    string      `json:"message"`
	Data   interface{} `json:"result"`
}

func (this *Controller) Error_return(msg string) {
	if this.Err_status == 0 {
		this.Err_status = 101
	}
	this.Jsonmsg.Status = this.Err_status
	this.Jsonmsg.Msg = msg
	this.Jsonmsg.Data = nil
	b, _ := json.Marshal(&this.Jsonmsg)
	this.W.Header().Set("content-type", "application/json")
	this.W.Write(b)
}

func (this *Controller) Success_return(msg string, data interface{}) {
	this.Jsonmsg.Status = 200
	this.Jsonmsg.Msg = msg
	this.Jsonmsg.Data = data
	b, _ := json.Marshal(&this.Jsonmsg)
	//fmt.Println(string(b),jsonstr)
	this.W.Header().Set("content-type", "application/json")
	this.W.Write(b)
}
