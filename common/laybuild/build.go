package laybuild

import (
	"fmt"
	"github.com/guyigood/gylib/common"
	"github.com/guyigood/gylib/gydblib"
	"io/ioutil"
	"os"
	"strings"
)

type Build_GO struct {
	Id         string
	Lb         string
	Path       string
	GO_Path    string
	Model_path string
	H5_Path    string
}

func NewBuildGo() (qb *Build_GO) {
	qb = new(Build_GO)
	return
}

func (this *Build_GO) Build_all(id, lb string) {
	db := gydblib.Get_New_Main_DB()
	data := db.Tbname("db_tb_dict").Where("id=" + id).Find()
	this.Id = id
	this.Lb = lb
	this.Model_path = "model"
	if this.Lb == "PCTPL" {
		this.Path = "views/pctpl"
		this.GO_Path = "controller/admin"
		this.H5_Path = "views/admin/"

		//this.Build_DB_model(data["name"])
	} else {
		this.Path = "views/waptpl"
		this.GO_Path = "controller/wap"
		this.H5_Path = "views/wap/"

	}
	if len(data) > 0 {
		this.build_gofile(data)
		this.build_pcindex(data)
		this.build_edit_html(data)
	}

}

func (this *Build_GO) Build_html_all(id, lb string) {
	db := gydblib.Get_New_Main_DB()
	data := db.Tbname("db_tb_dict").Where("id=" + id).Find()
	this.Id = id
	this.Lb = lb
	this.Model_path = "model"

	this.Path = "views/pctpl"
	this.GO_Path = "controller/admin"
	this.H5_Path = "views/admin/"
	//this.Build_DB_model(data["name"])
	if len(data) > 0 {
		this.build_pcindex(data)
		this.build_edit_html(data)
	}

}

func (this *Build_GO) build_gofile(data map[string]string) {
	db := gydblib.Get_New_Main_DB()
	tbname := strings.Replace(data["name"], db.Db_perfix, "", -1)
	tbname = common.ToUP_map_name(tbname)
	contpl := read_file(this.Path + "/pc.gotpl")
	contpl = strings.Replace(contpl, "{{tb_name}}", strings.Replace(data["name"], db.Db_perfix, "", -1), -1)
	contpl = strings.Replace(contpl, "{{pack_name}}", strings.ToLower(tbname), -1)
	contpl = strings.Replace(contpl, "{{module_name}}", strings.Title(strings.ToLower(tbname)), -1)
	list := db.Tbname("db_fd_dict").Where("list_tb_name<>'0' and t_id=" + this.Id).Select()
	for _, v := range list {
		tmp_str := read_file(this.Path + "/db.gotpl")
		tb_name := strings.Replace(v["list_tb_name"], db.Db_perfix, "", -1)
		tmp_str = strings.Replace(tmp_str, "{{tb_name}}", tb_name, -1)
		tmp_str = strings.Replace(tmp_str, "{{tbname}}", v["list_tb_name"], -1)
		tmp_str = strings.Replace(tmp_str, "{{where}}", v["list_where"], -1)
		contpl = strings.Replace(contpl, "{{edit_list}}", tmp_str+"{{edit_list}}", -1)
	}
	contpl = strings.Replace(contpl, "{{edit_list}}", "", -1)
	dirpath := this.GO_Path
	if !common.PathExists(dirpath) {
		os.Mkdir(dirpath, os.ModePerm)
	}
	write_file(dirpath+"/"+strings.ToLower(tbname)+".go", contpl)
}

func read_file(filename string) string {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(dat)
}

func write_file(filename string, memo string) {
	err := ioutil.WriteFile(filename, []byte(memo), 0777)
	if err != nil {
		fmt.Println(err)
	}
}

func (this *Build_GO) build_pcindex(data map[string]string) {
	db := gydblib.Get_New_Main_DB()
	tbname := strings.Replace(data["name"], db.Db_perfix, "", -1)
	tbname = common.ToUP_map_name(tbname)
	contpl := read_file(this.Path + "/index.html")
	contpl = strings.Replace(contpl, "{{tpl_title}}", data["cn_name"], -1)
	url_add := "/admin/" + strings.ToLower(tbname) + "_list/"
	contpl = strings.Replace(contpl, "{{tpl_url_list}}", url_add, -1)
	url_add = "/admin/" + strings.ToLower(tbname) + "_edit/"
	contpl = strings.Replace(contpl, "{{tpl_url_edit}}", url_add, -1)
	url_add = "/admin/" + strings.ToLower(tbname) + "_del/"
	contpl = strings.Replace(contpl, "{{tpl_url_del}}", url_add, -1)
	/*生成其他按钮和对应的js语句*/
	link_btn := strings.Split(data["url_name"], "\n")
	link_url := strings.Split(data["url_add"], "\n")
	other_btn := ""
	for k, v := range link_btn {
		if v == "" {
			continue
		}
		other_btn += fmt.Sprintf("<button class=\"layui-btn\" data-type=\"other_btn_%v\">%v</button>\n", k, v)

	}
	contpl = strings.Replace(contpl, "{{tpl_other_btn}}", other_btn, -1)
	other_btn = ""
	for k, v := range link_url {
		if v == "" {
			continue
		}
		other_btn += fmt.Sprintf("other_btn_%v:function(){\nperContent = layer.open({type: 2,title:'管理', maxmin: true,content: '%v'});\nlayer.full(perContent);\n},", k, v)
	}
	contpl = strings.Replace(contpl, "{{tpl_other_url_link}}", other_btn, -1)
	contpl = this.build_pcseach_form(contpl)
	contpl = this.build_table(contpl)
	dirpath := this.H5_Path + strings.ToLower(tbname)
	if !common.PathExists(dirpath) {
		os.Mkdir(dirpath, os.ModePerm)
	}
	write_file(dirpath+"/index.html", contpl)
}

func (this *Build_GO) build_pcseach_form(content string) string {
	db := gydblib.Get_New_Main_DB()
	data := db.Tbname("db_fd_dict").Where("f_is_search=1 and t_id=" + this.Id).Select()
	search_form := ""
	search_value := ""
	for _, v := range data {
		if v["f_type"] == "下拉框" {
			cn_name := v["cn_name"]
			en_name := v["name"]
			search_form += fmt.Sprintf("<div class=\"layui-inline layui-col-space2\"></div><div class=\"layui-inline\">\n <select name='S_%v' id='S_%v' class=\"layui-select\"><option value=''>请选择%v</option>", en_name, en_name, cn_name)

			if v["list_tb_name"] == "0" {

				if search_value == "" {
					search_value = fmt.Sprintf("S_%v:$('#S_%v').val()", en_name, en_name)
				} else {
					search_value += fmt.Sprintf(",S_%v:$('#S_%v').val()", en_name, en_name)
				}
				list_val := strings.Split(v["list_val"], "|")
				list_display := strings.Split(v["list_display"], "|")
				for k, val := range list_val {
					search_form += "<option value='" + val + "'>" + list_display[k] + "</option>\n"
				}

			} else {
				/*获取数据库中对应的条件，生成下拉框*/
				db.Dbinit()
				list_val := strings.Split(v["list_val"], "|")
				list_display := strings.Split(v["list_display"], "|")
				search_form += "{{range $key,$vo:=." + strings.Replace(v["list_tb_name"], db.Db_perfix, "", -1) + "}}\n"
				for k, val := range list_val {
					search_form += "<option value='{{$vo." + val + "}}'>{{$vo." + list_display[k] + "}}</option>\n"
				}
				search_form += "{{end}}\n"

			}
			search_form += "</select> </div>"
		} else if v["f_type"] == "单选框" {
			search_form += fmt.Sprintf("<div class=\"layui-inline layui-col-space2\"></div><div class=\"layui-inline\"> \n<select name='S_%v' id='S_%v' class=\"layui-select\"><option value=''>请选择%v</option>", v["name"], v["name"], v["cn_name"])
			if search_value == "" {
				search_value = fmt.Sprintf("S_%v:$('#S_%v').val()", v["name"], v["name"])
			} else {
				search_value += fmt.Sprintf(",S_%v:$('#S_%v').val()", v["name"], v["name"])
			}
			list_val := strings.Split(v["list_val"], "|")
			list_display := strings.Split(v["list_display"], "|")
			for k, v := range list_val {
				search_form += "<option value='" + v + "'>" + list_display[k] + "</option>\n"
			}
			search_form += "</select> </div>\n"
		} else {
			search_form += fmt.Sprintf("\n<div class=\"layui-inline layui-col-space2\"></div><div class=\"layui-inline\"><input class=\"layui-input\" name=\"I_%v\" id=\"I_%v\" autocomplete=\"off\"  placeholder=\"%v\">\n</div>", v["name"], v["name"], v["cn_name"])
			if search_value == "" {
				search_value = fmt.Sprintf("I_%v:$('#I_%v').val()", v["name"], v["name"])
			} else {
				search_value += fmt.Sprintf(",I_%v:$('#I_%v').val()", v["name"], v["name"])
			}
		}

	}
	content = strings.Replace(content, "{{tpl_search}}", search_form, -1)
	content = strings.Replace(content, "{{tpl_search_value}}", search_value, -1)
	return content
}

func (this *Build_GO) build_table(content string) string {
	db := gydblib.Get_New_Main_DB()
	data := db.Tbname("db_fd_dict").Where("f_is_display=1 and t_id=" + this.Id).Order("f_list_order").Select()
	table_str := ""
	for _, v := range data {
		sort_str := ""
		if v["f_sort"] == "1" {
			sort_str = ",sort:true"
		}
		tmplet_str := ""
		if v["f_type"] == "单选框" || (v["f_type"] == "下拉框" && v["list_tbname"] == "0") {
			list_val := strings.Split(v["list_val"], "|")
			list_display := strings.Split(v["list_display"], "|")
			where_str := ""
			for k, val := range list_val {
				where_str += fmt.Sprintf("\nif (d.%v == %v) {\n", v["name"], val)
				where_str += fmt.Sprintf("return '<span class=\"layui-badge layui-bg-green\">%v</span>';\n}", list_display[k])
			}
			tmplet_str = ",templet: function (d) {\n" + where_str + "\n}"

		}
		table_str += fmt.Sprintf(", \n{field: '%v', title: '%v' %v %v}", v["name"], v["cn_name"], sort_str, tmplet_str)
	}
	content = strings.Replace(content, "{{tpl_table}}", table_str, -1)
	return content
}

func (this *Build_GO) build_edit_html(data map[string]string) {
	db := gydblib.Get_New_Main_DB()
	tbname := strings.Replace(data["name"], db.Db_perfix, "", -1)
	tbname = common.ToUP_map_name(tbname)
	contpl := read_file(this.Path + "/edit.html")
	contpl = strings.Replace(contpl, "{{tpl_title}}", data["cn_name"], -1)
	url_add := "/admin/" + strings.ToLower(tbname) + "_save/"
	contpl = strings.Replace(contpl, "{{tpl_url_save}}", url_add, -1)
	contpl = this.build_form(contpl)
	dirpath := this.H5_Path + strings.ToLower(tbname)
	if !common.PathExists(dirpath) {
		os.Mkdir(dirpath, os.ModePerm)
	}
	write_file(dirpath+"/edit.html", contpl)
}

func (this *Build_GO) build_form(content string) string {
	db := gydblib.Get_New_Main_DB()
	data := db.Tbname("db_fd_dict").Where("f_ed_display=1 and t_id=" + this.Id).Order("f_ed_order").Select()
	form_str := ""
	java_str := ""
	for _, v := range data {
		if v["f_type"] == "单行文本" || v["f_type"] == "密码" || v["f_type"] == "日期时间" {
			form_str += this.set_inputbox(v) + "\n"
		}
		if v["f_type"] == "单选框" {
			form_str += this.set_radio(v) + "\n"
		}
		if v["f_type"] == "多选框" {
			form_str += this.set_checkbox(v) + "\n"
		}
		if v["f_type"] == "下拉框" {
			form_str += this.set_select(v) + "\n"
		}
		if v["f_type"] == "多行文本" || v["f_type"] == "富文本" {
			form_str += this.set_textarea(v) + "\n"
		}
		if v["f_type"] == "单图片" {
			form_str += this.set_img(v) + "\n"
		}
		form_str = this.set_common_text(v, form_str)
		java_str += this.set_java(v) + "\n"

	}
	content = strings.Replace(content, "{{tpl_edit_form}}", form_str, -1)
	content = strings.Replace(content, "{{js}}", java_str, -1)
	return content
}

func (this *Build_GO) set_common_text(data map[string]string, contpl string) string {
	result := ""
	if data["f_isnull"] == "1" {
		if data["f_check"] != "" {
			result += fmt.Sprintf("required lay-verify=\"required|%v\"", data["f_check"])
		} else {
			result += "required lay-verify=\"required\""
		}
	}
	contpl = strings.Replace(contpl, "{{other_input}}", result, -1)
	return contpl
}

func (this *Build_GO) set_inputbox(data map[string]string) string {
	contpl := read_file(this.Path + "/inputbox.html")
	contpl = strings.Replace(contpl, "{{name}}", data["name"], -1)
	contpl = strings.Replace(contpl, "{{cn_name}}", data["cn_name"], -1)
	return contpl
}

func (this *Build_GO) set_radio(data map[string]string) string {
	contpl := read_file(this.Path + "/radio.html")
	contpl = strings.Replace(contpl, "{{cn_name}}", data["cn_name"], -1)
	contpl = strings.Replace(contpl, "{{name}}", data["name"], -1)
	ck_list := strings.Split(data["list_display"], "|")
	ck_val := strings.Split(data["list_val"], "|")
	h_str := ""
	for k, dis := range ck_list {
		h_str += fmt.Sprintf("\n<input type=\"radio\" name=\"%s\" value=\"%s\" title=\"%s\" {{if eq .data.%s \"%s\"}}checked=\"\"{{end}}>",
			data["name"], ck_val[k], dis, data["name"], ck_val[k])
	}
	contpl = strings.Replace(contpl, "{{radio}}", h_str, -1)
	return contpl
}

func (this *Build_GO) set_checkbox(data map[string]string) string {
	db := gydblib.Get_New_Main_DB()
	contpl := read_file(this.Path + "/checkbox.html")
	contpl = strings.Replace(contpl, "{{cn_name}}", data["cn_name"], -1)
	contpl = strings.Replace(contpl, "{{name}}", data["name"], -1)
	c_str := "<input type=\"checkbox\" name=\"{{name}}\"  {{check}}>"
	c_str = strings.Replace(c_str, "{{name}}", data["name"], -1)
	ch_box := ""
	if data["list_tb_name"] == "0" {
		ck_list := strings.Split(data["list_display"], "|")
		ck_val := strings.Split(data["list_val"], "|")
		for k, dis := range ck_list {
			h_str := "title=\"{{title}}\" {{ischeckbox .data.{{name}} \"{{value}}\" }}"
			h_str = strings.Replace(h_str, "{{name}}", data["name"], -1)
			h_str = strings.Replace(h_str, "{{title}}", dis, -1)
			h_str = strings.Replace(h_str, "{{value}}", ck_val[k], -1)
			ch_box += strings.Replace(c_str, "{{check}}", h_str, -1) + "\n"
		}
	} else {
		ch_box = strings.Replace("{{range  $key,$vo:=.{{tbname}} }}", "{{tbname}}", strings.Replace(data["list_tb_name"], db.Db_perfix, "", -1), -1)
		h_str := "title=\"{{ $vo.{{title}} }}\" {{ischeckbox $vo.{{vname}} {{.data.{{name}}}} }}"
		h_str = strings.Replace(h_str, "{{title}}", data["list_display"], -1)
		h_str = strings.Replace(h_str, "{{name}}", data["name"], -1)
		h_str = strings.Replace(h_str, "{{vname}}", data["list_val"], -1)
		ch_box += strings.Replace(c_str, "{{check}}", h_str, -1) + "\n{{end}}\n"
	}
	contpl = strings.Replace(contpl, "{{checkbox}}", ch_box, -1)
	return contpl
}

func (this *Build_GO) set_select(data map[string]string) string {
	db := gydblib.Get_New_Main_DB()
	contpl := read_file(this.Path + "/select.html")
	contpl = strings.Replace(contpl, "{{cn_name}}", data["cn_name"], -1)
	contpl = strings.Replace(contpl, "{{name}}", data["name"], -1)
	ch_box := ""
	if data["list_tb_name"] == "0" {
		ck_list := strings.Split(data["list_display"], "|")
		ck_val := strings.Split(data["list_val"], "|")
		for k, dis := range ck_list {
			h_str := "<option value='{{value}}' {{if eq $.data.{{name}} \"{{value}}\" }}selected{{end}}>{{title}}</option>"
			h_str = strings.Replace(h_str, "{{name}}", data["name"], -1)
			h_str = strings.Replace(h_str, "{{title}}", dis, -1)
			h_str = strings.Replace(h_str, "{{value}}", ck_val[k], -1)
			ch_box += h_str + "\n"
		}
	} else {
		ch_box = strings.Replace("{{range  $key,$vo:=.{{tbname}} }}\n", "{{tbname}}", strings.Replace(data["list_tb_name"], db.Db_perfix, "", -1), -1)
		h_str := "<option value='{{$vo.{{vname}}}}' {{if eq $.data.{{name}} $vo.{{vname}} }}selected{{end}}>{{$vo.{{title}}}}</option>\n"
		h_str = strings.Replace(h_str, "{{title}}", data["list_display"], -1)
		h_str = strings.Replace(h_str, "{{name}}", data["name"], -1)
		h_str = strings.Replace(h_str, "{{vname}}", data["list_val"], -1)
		ch_box += h_str + "{{end}}\n"
	}
	contpl = strings.Replace(contpl, "{{option}}", ch_box, -1)
	return contpl
}

func (this *Build_GO) set_textarea(data map[string]string) string {

	contpl := read_file(this.Path + "/textarea.html")
	contpl = strings.Replace(contpl, "{{cn_name}}", data["cn_name"], -1)
	contpl = strings.Replace(contpl, "{{name}}", data["name"], -1)
	if data["f_type"] == "富文本" {
		contpl = strings.Replace(contpl, "{{style}}", "", -1)
	} else {
		contpl = strings.Replace(contpl, "{{style}}", " class=\"layui-textarea\"", -1)
	}
	return contpl
}

func (this *Build_GO) set_img(data map[string]string) string {

	contpl := read_file(this.Path + "/sing_img.html")
	contpl = strings.Replace(contpl, "{{cn_name}}", data["cn_name"], -1)
	contpl = strings.Replace(contpl, "{{name}}", data["name"], -1)
	return contpl
}

func (this *Build_GO) set_java(data map[string]string) string {
	contpl := ""
	if data["f_type"] == "日期时间" {
		contpl = read_file(this.Path + "/date.html")
		contpl = strings.Replace(contpl, "{{cn_name}}", data["cn_name"], -1)
		contpl = strings.Replace(contpl, "{{name}}", data["name"], -1)
	}
	if data["f_type"] == "富文本" {
		contpl = read_file(this.Path + "/editmemo.html")
		contpl = strings.Replace(contpl, "{{cn_name}}", data["cn_name"], -1)
		contpl = strings.Replace(contpl, "{{name}}", data["name"], -1)
	}
	if data["f_type"] == "单图片" {
		contpl = read_file(this.Path + "/uploadfile.html")
		contpl = strings.Replace(contpl, "{{cn_name}}", data["cn_name"], -1)
		contpl = strings.Replace(contpl, "{{name}}", data["name"], -1)
	}

	return contpl
}
