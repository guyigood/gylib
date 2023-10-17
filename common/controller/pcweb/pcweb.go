package pcweb

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
	W            http.ResponseWriter
	R            *http.Request
	Tplname      string
	Data         map[string]interface{}
	Url_list     map[string]interface{}
	JsonData     map[string]interface{}
	Postdata     map[string]interface{}
	Body         []byte
	Site_name    string
	Api_url      string
	Get_url      string
	Masterdb     string
	Access_token string
	T_data       map[string]string
	F_data       []map[string]string
	Tokenname    string
	Jsonmsg      Json_msg
	Err_status   int
}

type Json_msg struct {
	Code   int         `json:"code"`
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func set_funcmap() template.FuncMap {
	tempfunc := make(template.FuncMap)
	tempfunc["str2html"] = Strtohtml
	tempfunc["ischeckbox"] = ISCheckbox
	tempfunc["date2local"] = Date2Local
	tempfunc["date2int"] = Date2Int
	tempfunc["round"] = Round
	return (tempfunc)
}

func (this *Controller) Rander() {
	tempfunc := set_funcmap()
	filename := strings.Split(this.Tplname, "/")
	tmpname := filename[len(filename)-1]
	t := template.New(tmpname)
	//t := template.New(tmpname).Delims("<<",">>")
	t = t.Funcs(tempfunc)
	t, _ = t.ParseFiles(this.Tplname)
	t.Execute(this.W, this.Data)
}

func (this *Controller) MuitplRander(arg ...string) {
	tempfunc := set_funcmap()
	filename := strings.Split(this.Tplname, "/")
	tmpname := filename[len(filename)-1]
	t := template.New(tmpname)
	//t := template.New(tmpname).Delims("<<",">>")
	t = t.Funcs(tempfunc)
	t, _ = t.ParseFiles(this.Tplname)
	t, _ = t.ParseFiles(arg...)

	t.Execute(this.W, this.Data)
}

func Strtohtml(htmlstr interface{}) interface{} {
	return template.HTML(htmlstr.(string))
}

func ISCheckbox(code interface{}, qxz interface{}) interface{} {
	if qxz == nil {
		return template.HTML("")
	}
	qxzmemo := fmt.Sprintf("%v", qxz)
	code_str := fmt.Sprintf("%v", code)
	result := ""
	qxzarr := strings.Split(qxzmemo, ",")
	for _, val := range qxzarr {
		if code_str == val {
			result = "checked"
			break
		}
	}
	return template.HTML(result)
}

func Date2Local(htmlstr interface{}) interface{} {
	result := fmt.Sprintf("%v", htmlstr)
	result = strings.Replace(result, "+0800 CST", "", -1)
	result = strings.Replace(result, "+0000 UTC", "", -1)
	return template.HTML(result)
}

func Date2Int(htmlstr interface{}) interface{} {
	result := time.Unix(htmlstr.(int64), 0).Format("2006-01-02 15:04:05")
	return template.HTML(result)
}

func Round(htmlstr interface{}, s int) interface{} {
	result := fmt.Sprintf("%."+strconv.Itoa(s)+"f", htmlstr)
	return template.HTML(result)
}

//func (this *Controller) Get_redis_tbname(tbname string) (map[string]string, []map[string]string) {
//
//	client := rediscomm.NewRedisComm()
//	r_tb := client.SetKey("tb_dict").SetFiled(tbname).Hget_map()
//	if (r_tb == nil) {
//		return this.Get_mysql_dict(tbname)
//	}
//	f_db := client.SetKey("fd_dict").SetFiled(tbname).Hget_map()
//	tb_data := r_tb.(map[string]interface{})
//	fd_data := make([]map[string]string, 0)
//	for _, v := range f_db.([]interface{}) {
//		tmp := datatype.Map2str(v.(map[string]interface{}))
//		fd_data = append(fd_data, tmp)
//	}
//
//	return datatype.Map2str(tb_data), fd_data
//
//}

//func (this *Controller) Get_mysql_dict(tbname string) (map[string]string, []map[string]string) {
//	db := lib.NewQuerybuilder()
//	data := db.Tbname("db_tb_dict").Where(fmt.Sprintf("name='%v'", lib.Db_perfix+tbname)).Find();
//	if (data == nil) {
//		return nil, nil
//	}
//	db.Dbinit()
//	fd_data := db.Tbname("db_fd_dict").Where(fmt.Sprintf("t_id=%v", data["id"])).Select()
//	list_data := db.Tbname("db_fd_dict").Where(fmt.Sprintf("t_id=%v and list_tb_name<>'0'", data["id"])).Select()
//	client := rediscomm.NewRedisComm()
//	client.SetKey("tb_dict").SetFiled(tbname).SetData(data).Hset_map()
//	if (fd_data != nil) {
//		client.SetKey("fd_dict").SetFiled(tbname).SetData(fd_data).Hset_map()
//	}
//	if (list_data != nil) {
//		client.SetKey("fd_list").SetFiled(tbname).SetData(list_data).Hset_map()
//
//	}
//	return data, fd_data
//
//}

func (this *Controller) Error_return(msg string) {
	this.Jsonmsg.Code = 0
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
	this.Jsonmsg.Code = 0
	this.Jsonmsg.Status = 100
	this.Jsonmsg.Msg = msg
	this.Jsonmsg.Data = data
	b, _ := json.Marshal(&this.Jsonmsg)
	//fmt.Println(string(b),jsonstr)
	this.W.Header().Set("content-type", "application/json")
	this.W.Write(b)
}

//func (this *Controller) Get_city_name(id string) string {
//	db := gydblib.Dbcon
//	data := db.Tbname("china").Where("id=" + id).Find()
//	//fmt.Println(db.GetLastSql(),data)
//	return data["name"]
//}
