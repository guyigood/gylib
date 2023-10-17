package modelbuild

import (
	"fmt"
	"gylib/common"
	"gylib/common/mysqlmodel"
	"gylib/gydblib"
	"io/ioutil"
	"os"
	"strings"
)

type Build_GO struct {
	Id         string
	Model_path string
}

func NewModelBuild() (qb *Build_GO) {
	qb = new(Build_GO)
	return
}

func (this *Build_GO) Build_all(id string) {
	db := gydblib.Get_New_Main_DB()
	data := db.Tbname("db_tb_dict").Where("id=" + id).Find()
	this.Model_path = "model"
	this.Build_DB_model(data["name"])
}

func (this *Build_GO) Build_DB_model(tbname string) {
	db := gydblib.Get_New_Main_DB()
	model := mysqlmodel.NewMysql_model()
	model.TableName = strings.Replace(tbname, db.Db_perfix, "", -1)
	struct_str := `package model

              `
	struct_str += model.Build()
	dirpath := this.Model_path
	if !common.PathExists(dirpath) {
		os.Mkdir(dirpath, os.ModePerm)
	}
	write_file(dirpath+"/"+model.TableName+".go", struct_str)
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
