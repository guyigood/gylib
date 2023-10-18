package mysqlsplit

import (
	"github.com/guyigood/gylib/gydblib"
	"strings"
	"time"
)

type Table_Split struct {
	conn   *gydblib.Mysqlcon
	Dbname string
}

/*
	数据表结构

CREATE TABLE `sl_table_split` (

	`id` int(11) NOT NULL AUTO_INCREMENT,
	`tbname` varchar(255) DEFAULT NULL,
	`cur_name` varchar(255) DEFAULT NULL,
	`px` int(11) DEFAULT '0',
	`is_cur` int(11) DEFAULT '0',
	PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/
const mast_tbname = "table_split"

func NewGylib_Table_Split(dbname string) *Table_Split {
	this := new(Table_Split)
	this.Dbname = dbname
	this.conn = gydblib.NewMySql_Server_DB(dbname)
	return this
}

func (this *Table_Split) GetCurName(name string) string {
	data := this.conn.Tbname(mast_tbname).Where(map[string]interface{}{"name": name, "is_cur": 1}).Find()
	if data != nil {
		return data["cur_name"]
	} else {
		return ""
	}
}

func (this *Table_Split) Split_Table(name string, auto_inc_pk string) bool {
	src_name := this.conn.Db_perfix + name
	dt := time.Now().String()[:19]
	dt = strings.Replace(dt, "-", "", -1)
	dt = strings.Replace(dt, ":", "", -1)
	des_name := src_name + dt
	sqltext := "create table " + des_name + " Like " + src_name
	ct := this.conn.Tbname(mast_tbname).Where(map[string]interface{}{"name": name}).Count()
	ct++
	var err error
	this.conn.BeginStart()
	_, err = this.conn.Excute(sqltext, nil)
	if err != nil {
		this.conn.Rollback()
		return false
	}
	_, err = this.conn.Tbname(mast_tbname).Insert(map[string]interface{}{"tbname": name, "cur_name": des_name, "px": ct, "is_cur": 1})
	if err != nil {
		this.conn.Rollback()
		return false
	}
	sqltext = "update " + this.conn.Db_perfix + mast_tbname + " set is_cur=0 where cur_name<>'" + des_name + "'"
	_, err = this.conn.Excute(sqltext, nil)
	if err != nil {
		this.conn.Rollback()
		return false
	}
	if auto_inc_pk != "" {
		//自动增长主键的操作，获取当前最大的增长值，设置新表的增长值
		sqltext = "select max(" + auto_inc_pk + ")+1 as id from " + src_name
		max_data := this.conn.Query(sqltext, nil)
		if max_data == nil {
			this.conn.Rollback()
			return false
		}
		new_id := max_data[0]["id"]
		sqltext = "alter table " + des_name + " auto_increment = " + new_id
		_, err = this.conn.Excute(sqltext, nil)
		if err != nil {
			this.conn.Rollback()
			return false
		}
	}

	this.conn.Commit()
	return true
}
