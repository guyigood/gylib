package oracledb

import (
	"database/sql"
	"fmt"
	"github.com/guyigood/gylib/common"
	"github.com/guyigood/gylib/common/datatype"
	_ "github.com/mattn/go-oci8"
	"strconv"
	"strings"
	"sync"
	"time"
	//"reflect"
)

type Oraclecon struct {
	Masterdb *sql.DB
	//Slavedb     []*sql.DB
	Db_host     string
	Db_user     string
	Db_port     string
	Db_name     string
	Db_password string
	Where_param []interface{}
	SqlTx       *sql.Tx
	Slock       sync.RWMutex
	Page_start  int
	Page_end    int
	Tablename   string
	Sql_where   string
	Sql_order   string
	Sql_fields  string
	Is_open     bool
	Sql_param   []interface{}
	//Sql_limit   string
	Db_perfix   string
	PRK_editfd  string
	Query_data  []map[string]interface{}
	LastSqltext string
}

var G_Dbcon *Oraclecon

var G_dbtables map[string]interface{}
var G_fd_list map[string]interface{}

func init() {
	G_dbtables = make(map[string]interface{})
	G_fd_list = make(map[string]interface{})

}

func NewOracle_Server_DB(action string) *Oraclecon {
	this := new(Oraclecon)
	if action == "" {
		action = "database"
	}

	data := common.Getini("conf/app.ini", action, map[string]string{"db_user": "root", "db_password": "",
		"db_host": "127.0.0.1", "db_port": "1521", "db_name": "", "db_maxpool": "20", "db_minpool": "5", "db_perfix": "", "db_type": "", "slavedb": ""})
	//"bfcrm8/DHHZDHHZ@10.100.2.202:1521/crmtest"
	if data["db_name"] == "" {
		return nil
	}
	if data["db_type"] != "" && data["db_type"] != "oracle" {
		return nil
	}

	con := fmt.Sprintf("%s/%s@%s:%s/%s", data["db_user"],
		data["db_password"], data["db_host"],
		data["db_port"], data["db_name"])

	//fmt.Println(con)
	var err error
	this.Masterdb, err = sql.Open("oci8", con)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	maxpool, _ := strconv.Atoi(data["db_maxpool"])
	minpool, _ := strconv.Atoi(data["db_minpool"])
	//fmt.Println(this.Masterdb,con,maxpool,minpool)
	this.Masterdb.SetMaxOpenConns(maxpool)
	this.Masterdb.SetMaxIdleConns(minpool)
	this.Masterdb.SetConnMaxLifetime(time.Minute * 5)
	err = this.Masterdb.Ping()
	if err != nil {
		fmt.Println("PING", err)
		return nil
	}
	/*this.Slavedb = make([]*sql.DB, 0)
	if(data["slavedb"]!="") {
		iplist := strings.Split(data["slavedb"], ",")
		for _, v := range iplist {
			con1 := fmt.Sprintf("%s/%s@%s:%s/%s", data["db_user"],
				data["db_password"], v,
				data["db_port"], data["db_name"])
			sqldb1, _ := sql.Open("oci8", con1)
			maxpool, _ := strconv.Atoi(data["db_maxpool"])
			minpool, _ := strconv.Atoi(data["db_minpool"])
			sqldb1.SetMaxOpenConns(maxpool)
			sqldb1.SetMaxIdleConns(minpool)
			sqldb1.SetConnMaxLifetime(time.Minute * 5)
			sqldb1.Ping()
			this.Slavedb = append(this.Slavedb, sqldb1)
		}
	}*/
	this.SqlTx = nil
	this.Db_perfix = data["db_perfix"]
	this.Db_name = data["db_name"]
	this.Db_host = data["db_host"]
	this.Db_port = data["db_port"]
	this.Db_password = data["db_password"]
	this.Db_user = data["db_user"]
	return this
}

func Get_New_Main_DB() *Oraclecon {
	if G_Dbcon == nil {
		G_Dbcon = NewOracle_Server_DB("")
		if G_Dbcon != nil {
			G_Dbcon.Init_Redis_Struct()
		}
	}
	this := new(Oraclecon)
	this.Masterdb = G_Dbcon.Masterdb
	//this.Slavedb = G_Dbcon.Slavedb
	this.Db_perfix = G_Dbcon.Db_perfix
	this.Db_name = G_Dbcon.Db_name
	this.Db_host = G_Dbcon.Db_host
	this.Db_port = G_Dbcon.Db_port
	this.Db_password = G_Dbcon.Db_password
	//this.Masterdb=Dbcon.Masterdb
	//this.Slavedb=Dbcon.Slavedb
	//this.Db_perfix = Dbcon.Db_perfix
	//this.Db_name = Dbcon.Db_name
	//this.Db_host = Dbcon.Db_host
	//this.Db_port = Dbcon.Db_port
	//this.Db_password = Dbcon.Db_password
	this.SqlTx = nil
	this.Dbinit()
	return this
}

func (this *Oraclecon) Merge_And_where(where_str, new_str string) string {
	result := where_str
	if where_str != "" {
		result += " and " + new_str
	} else {
		result = new_str
	}
	return result
}

func (this *Oraclecon) Merge_OR_where(where_str, new_str string) string {
	result := where_str
	if where_str != "" {
		result += " or " + new_str
	} else {
		result = new_str
	}
	return result
}

func (this *Oraclecon) BeginStart() bool {
	tx, err := this.Masterdb.Begin()
	if err != nil {
		return false
	}
	this.SqlTx = tx
	return true
}

func (this *Oraclecon) SetPK(str_val string) *Oraclecon {
	this.PRK_editfd = str_val
	return this
}

/*
*
初始化结构
*/
func (this *Oraclecon) Dbinit() {
	this.Tablename = ""
	//this.Sql_limit = ""
	this.Sql_order = ""
	this.Sql_fields = ""
	this.Sql_where = ""
	this.Page_end = 0
	this.Page_start = 0
	this.Is_open = true
	this.PRK_editfd = ""
	this.Slock.Lock()
	this.Query_data = make([]map[string]interface{}, 0)
	this.Sql_param = make([]interface{}, 0)
	this.Where_param = make([]interface{}, 0)
	this.Slock.Unlock()
}

func (this *Oraclecon) SetOpen(flag bool) *Oraclecon {
	this.Is_open = flag
	return this
}

/*
设置数据表
*/
func (this *Oraclecon) Tbname(name string) *Oraclecon {
	//data := common.Getini("conf/app.ini","database",map[string]string{"db_perfix":""})
	this.Dbinit()
	this.Tablename = this.Db_perfix + name
	return this
}

func (this *Oraclecon) Fileds(name string) *Oraclecon {
	this.Sql_fields = name
	return this
}

func (this *Oraclecon) SetWhere(where string, param ...interface{}) *Oraclecon {
	if this.Sql_where == "" {
		this.Sql_where = where
	} else {
		this.Sql_where += " and (" + where + ")"
	}
	for _, val := range param {
		this.Where_param = append(this.Where_param, val)
	}
	return this
}

func (this *Oraclecon) Where(where interface{}) *Oraclecon {
	//kk:= reflect.TypeOf(where)
	//fmt.Println(kk)
	switch where.(type) {
	case string:
		if datatype.Type2str(where) == "" {
			return this
		}
		if this.Sql_where == "" {
			this.Sql_where = where.(string)
		} else {
			this.Sql_where += " and (" + where.(string) + ")"
		}
	default:
		this.Slock.Lock()
		tmp_arr := where.(map[string]interface{})
		if len(tmp_arr) > 0 {
			this.Query_data = append(this.Query_data, tmp_arr)
		}
		this.Slock.Unlock()
		//fmt.Println("query_data", this.Query_data)
	}

	return this
}

func (this *Oraclecon) Order(orderstr string) *Oraclecon {
	this.Sql_order = orderstr
	return this
}

func (this *Oraclecon) PageLimit(startct, endct int) *Oraclecon {
	this.Page_start = startct
	this.Page_end = endct
	return this
}

func (this *Oraclecon) Check_data_fields(fieldname string) bool {
	flag := false
	fd_list, ok := G_dbtables[this.Db_name+this.Tablename]
	if ok {
		for _, v := range fd_list.([]map[string]string) {
			record := v
			if this.Check_PK(fieldname) {
				continue
			}
			/*if (record["key"] == "PRI" && record["extra"] == "auto_increment") {
				continue
			}*/
			if record["column_name"] == fieldname {
				flag = true
				break
			}

		}
		return flag
	} else {
		this.Update_redis(this.Tablename)
		rows, _ := this.Masterdb.Query(`select  t.*,c.COMMENTS 
from user_tab_columns  t,user_col_comments  c
 where t.table_name = c.table_name and t.column_name = c.column_name and t.table_name ='` + this.Tablename + "'")
		if rows == nil {
			return false
		}
		defer rows.Close()
		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		for rows.Next() {
			//将行数据保存到record字典
			record := make(map[string]string)
			_ = rows.Scan(scanArgs...)
			for i, col := range values {
				if col != nil {
					record[strings.ToLower(columns[i])] = this.Type2str(col)
				}
			}
			if this.Check_PK(fieldname) {
				continue
			}
			/*if (record["key"] == "PRI" && record["extra"] == "auto_increment") {
				continue
			}*/
			if record["column_name"] == fieldname {
				flag = true
				break
			}

		}
		return flag
	}
}

func (this *Oraclecon) Check_PK(fdname string) bool {
	if this.PRK_editfd == "" {
		return false
	}
	list := strings.Split(this.PRK_editfd, ",")
	flag := false
	for _, val := range list {
		if val == "" {
			continue
		}
		if fdname == val {
			flag = true
			break
		}
	}
	return flag
}

func (this *Oraclecon) Type2str(val interface{}) string {
	//fmt.Println(fmt.Sprintf("%T,%v",val,val))
	var result string = ""
	switch val.(type) {
	case []string:
		strArray := val.([]string)
		result = strings.Join(strArray, "")
	case []uint8:
		result = string(val.([]uint8))
	default:
		result = fmt.Sprintf("%v", val)
	}
	return result
}

func (this *Oraclecon) Insert(postdata map[string]interface{}) (sql.Result, error) {
	//this.Wlock.Lock()
	//defer this.Wlock.Unlock()
	var sqltext string
	sqltext = "insert into \"" + this.Tablename + "\" ("
	values := " values ("
	i := 0
	this.Sql_param = make([]interface{}, 0)
	for k, v := range postdata {

		if this.Check_data_fields(k) == false {
			continue
		}

		if i > 0 {
			sqltext += ","
			values += ","
		}
		i++
		sqltext += "\"" + k + "\""

		fun_val, fun_flag := this.Get_fun_fields(k, datatype.Type2str(v))
		values += fun_val
		if fun_flag {
			continue
		}
		if datatype.Type2str(v) != "" {
			this.Sql_param = append(this.Sql_param, v)
		} else {
			this.Sql_param = append(this.Sql_param, nil)
		}
	}
	sqltext += ") " + values + ")"
	this.LastSqltext = sqltext
	if !this.Is_open {
		return nil, nil
	}
	//fmt.Println(i,sqltext)
	//fmt.Println(postdata,len(param_data),param_data)
	var result sql.Result
	var err error
	if this.SqlTx != nil {
		result, err = this.SqlTx.Exec(sqltext, this.Sql_param...)
	} else {
		result, err = this.Masterdb.Exec(sqltext, this.Sql_param...)
	}
	//fmt.Println(err)
	return result, err

}

func (this *Oraclecon) Update(postdata map[string]interface{}) (sql.Result, error) {
	//this.Wlock.Lock()
	//defer this.Wlock.Unlock()
	var sqltext string
	sqltext = fmt.Sprintf("update \"%v\" set ", this.Tablename)
	i := 0
	this.Sql_param = make([]interface{}, 0)
	for k, v := range postdata {
		if this.Check_data_fields(k) == false {
			continue
		}
		if i > 0 {
			sqltext += ","

		}
		i++
		fun_val, fun_flag := this.Get_fun_fields(k, datatype.Type2str(v))
		sqltext += "\"" + k + "\"= " + fun_val
		if fun_flag {
			continue
		}
		if datatype.Type2str(v) != "" {
			this.Sql_param = append(this.Sql_param, v)
		} else {
			this.Sql_param = append(this.Sql_param, nil)
		}
	}
	sqlwhere, param := this.Build_where()
	for _, v := range param {
		this.Sql_param = append(this.Sql_param, v)
	}
	sqltext += sqlwhere
	this.LastSqltext = sqltext
	if !this.Is_open {
		return nil, nil
	}
	//fmt.Println(sqltext, param_data)
	var result sql.Result
	var err error
	if this.SqlTx != nil {
		result, err = this.SqlTx.Exec(sqltext, this.Sql_param...)
	} else {
		result, err = this.Masterdb.Exec(sqltext, this.Sql_param...)
	}
	//fmt.Println(err,sqltext)
	return result, err
}

func (this *Oraclecon) Delete() (sql.Result, error) {
	//this.Wlock.Lock()
	//defer this.Wlock.Unlock()
	sqlwhere := ""
	this.Sql_param = make([]interface{}, 0)
	sqlwhere, this.Sql_param = this.Build_where()
	sqltext := fmt.Sprintf(" delete from \"%v\" %v", this.Tablename, sqlwhere)
	this.LastSqltext = sqltext
	if !this.Is_open {
		return nil, nil
	}
	var result sql.Result
	var err error
	if this.SqlTx != nil {
		result, err = this.SqlTx.Exec(sqltext, this.Sql_param...)
	} else {
		result, err = this.Masterdb.Exec(sqltext, this.Sql_param...)
	}
	return result, err
}

func (this *Oraclecon) SetDec(fdname string, quantity int) (sql.Result, error) {
	sqltext := fmt.Sprintf("update \"%v\" set \"%v\"=\"%v\"-%v", this.Tablename, fdname, fdname, quantity)
	sqlwhere, param := this.Build_where()
	sqltext += sqlwhere
	this.LastSqltext = sqltext
	var result sql.Result
	var err error
	if this.SqlTx != nil {
		result, err = this.SqlTx.Exec(sqltext, param...)
	} else {
		result, err = this.Masterdb.Exec(sqltext, param...)
	}
	return result, err
}

func (this *Oraclecon) SetInc(fdname string, quantity int) (sql.Result, error) {
	sqlwhere, param := this.Build_where()
	sqltext := fmt.Sprintf("update  \"%v\" set \"%v\"=\"%v\"+%v  %v", this.Tablename, fdname, fdname, quantity, sqlwhere)
	this.LastSqltext = sqltext
	var result sql.Result
	var err error
	if this.SqlTx != nil {
		result, err = this.SqlTx.Exec(sqltext, param...)
	} else {
		result, err = this.Masterdb.Exec(sqltext, param...)
	}
	return result, err
}

func (this *Oraclecon) Query(sqltext string, param []interface{}) []map[string]string {
	this.LastSqltext = sqltext
	var rows *sql.Rows
	var err error
	if this.SqlTx != nil {
		rows, err = this.SqlTx.Query(sqltext, param...)
	} else {

		rows, err = this.Masterdb.Query(sqltext, param...)
	}
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	//字典类型
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, _ := rows.Columns()
	result := make([]map[string]string, 0)
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		//将行数据保存到record字典
		record := make(map[string]string)
		_ = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = datatype.Type2str(col)
			} else {
				record[columns[i]] = ""
			}
		}
		result = append(result, record)
	}
	//fmt.Print(result)
	if len(result) == 0 {
		return nil
	}
	return result
}

func (this *Oraclecon) Query_One(sqltext string, param []interface{}) map[string]string {
	sqltext = "select * from (" + sqltext + ") A WHERE ROWNUM =1"
	this.LastSqltext = sqltext
	var rows *sql.Rows
	var err error
	if this.SqlTx != nil {
		rows, err = this.SqlTx.Query(sqltext, param...)
	} else {
		rows, err = this.Masterdb.Query(sqltext, param...)
	}
	if err != nil {
		return nil
	}
	defer rows.Close()
	//字典类型
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, _ := rows.Columns()
	result := make([]map[string]string, 0)
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		//将行数据保存到record字典
		record := make(map[string]string)
		_ = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = datatype.Type2str(col)
			} else {
				record[columns[i]] = ""
			}
		}
		result = append(result, record)
	}
	//fmt.Print(result)
	if len(result) == 0 {
		return nil
	}
	return result[0]
}

func (this *Oraclecon) Excute(sqltext string, param []interface{}) (sql.Result, error) {
	this.LastSqltext = sqltext
	var result sql.Result
	var err error
	if this.SqlTx != nil {
		result, err = this.SqlTx.Exec(sqltext, param...)
	} else {
		result, err = this.Masterdb.Exec(sqltext, param...)
	}
	return result, err
}

func (this *Oraclecon) set_sql(flag int) string {
	this.Slock.RLock()
	defer this.Slock.RUnlock()
	sqltext := ""
	if flag == 0 {

		if this.Sql_fields != "" {
			sqltext = "select " + this.Sql_fields + " from \"" + this.Tablename + "\""
		} else {
			sqltext = "select * from \"" + this.Tablename + "\""
		}
	} else {

		sqltext = "select count(1) as ct from \"" + this.Tablename + "\""

	}
	return sqltext
}

func (this *Oraclecon) Build_where() (string, []interface{}) {
	is_where := false
	sqltext := ""
	if this.Sql_where != "" {
		sqltext += " where " + this.Sql_where
		is_where = true
	}
	param_data := make([]interface{}, 0)
	if len(this.Query_data) > 0 {
		if is_where {
			sqltext += " and "

		} else {
			sqltext += " where "
		}
		i := 0
		this.Slock.RLock()
		for _, v := range this.Query_data {
			for key, val := range v {
				//if (this.Check_data_fields(key) == false) {
				//	continue
				//}
				if i > 0 {
					sqltext += " and "
				}
				i++
				switch val.(type) {
				//data["name"]=" %v like ?"
				//data["name"]=" %v>=(?)"
				//data["name"]="locate(?,`"+this.Tablename+"`.`%v`)>0"
				case map[string]interface{}:
					param_data = append(param_data, val.(map[string]interface{})["value"])
					sqltext += datatype.Type2str(val.(map[string]interface{})["name"])
				default:
					param_data = append(param_data, val)
					sqltext += "\"" + key + "\"=(:" + key + ") "

				}
			}
		}
		this.Slock.RUnlock()

	}
	if len(this.Where_param) > 0 {
		for _, val := range this.Sql_param {
			param_data = append(param_data, val)
		}
	}

	return sqltext, param_data
}

func (this *Oraclecon) Find() map[string]string {
	sqltext := this.set_sql(0)
	param_data := make([]interface{}, 0)
	tmpstr := ""
	tmpstr, param_data = this.Build_where()
	sqltext += tmpstr
	if this.Sql_order != "" {
		sqltext += " order by " + this.Sql_order
	}
	sqltext = "select * from (" + sqltext + ") A WHERE ROWNUM =1"
	this.LastSqltext = sqltext
	if !this.Is_open {
		//this.LastSqltext=sqltext
		this.Sql_param = make([]interface{}, 0)
		for _, p_val := range param_data {
			this.Sql_param = append(this.Sql_param, p_val)
		}
		return nil
	}
	var rows *sql.Rows
	var err error
	if this.SqlTx != nil {
		rows, err = this.SqlTx.Query(sqltext, param_data...)
	} else {
		rows, err = this.Masterdb.Query(sqltext, param_data...)
	}
	//fmt.Println("rows",rows,err)
	if err != nil {
		//fmt.Println("find",err)
		return nil
	}
	if rows == nil {
		return nil
	}

	defer rows.Close()
	//字典类型
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	record := make(map[string]string)
	for rows.Next() {
		//将行数据保存到record字典
		_ = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = datatype.Type2str(col)

			} else {
				record[columns[i]] = ""
			}
		}

	}
	if len(record) == 0 {
		return nil
	}
	return record
}

func (this *Oraclecon) Count() int64 {
	sqltext := this.set_sql(1)
	sqlwhere, param := this.Build_where()
	sqltext += sqlwhere
	this.LastSqltext = sqltext
	rows := this.Masterdb.QueryRow(sqltext, param...)
	var record int64
	rows.Scan(&record)

	return record
}

func (this *Oraclecon) Sum(fd string) float64 {
	var result float64
	sqltext := this.set_sql(1)
	sqltext = strings.Replace(sqltext, "count(1)", "sum(\""+fd+"\")", -1)
	sqlwhere, param := this.Build_where()
	sqltext += sqlwhere
	this.LastSqltext = sqltext
	rows := this.Masterdb.QueryRow(sqltext, param...)
	rows.Scan(&result)
	return result
}

func (this *Oraclecon) Select() []map[string]string {
	sqltext := this.set_sql(0)
	sqlwhere, param := this.Build_where()
	sqltext += sqlwhere
	if this.Sql_order != "" {
		sqltext += " order by " + this.Sql_order
	}
	if this.Page_end > 0 {
		sqltext = fmt.Sprintf(`SELECT * FROM (SELECT A.*, ROWNUM ROWN FROM (`+sqltext+`) A  
WHERE ROWNUM <=%v)  WHERE ROWN >= %v `, this.Page_end, this.Page_start)
	}
	this.LastSqltext = sqltext
	if !this.Is_open {
		this.Sql_param = make([]interface{}, 0)
		for _, p_val := range param {
			this.Sql_param = append(this.Sql_param, p_val)
		}
		return nil
	}
	//fmt.Println(sqltext)
	rows, err := this.Masterdb.Query(sqltext, param...)
	if err != nil {
		return nil
	}
	if rows == nil {
		return nil
	}

	defer rows.Close()

	//字典类型
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, _ := rows.Columns()
	result := make([]map[string]string, 0)
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	j := 0
	for rows.Next() {
		//将行数据保存到record字典
		record := make(map[string]string)
		_ = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = this.Type2str(col)
				//record[columns[i]] = col.([]byte)
			} else {
				record[columns[i]] = ""
			}
		}
		result = append(result, record)
		//result[j] = record
		j++

	}
	if len(result) == 0 {
		return nil
	}
	return result
}

func (this *Oraclecon) GetLastSql() string {
	return this.LastSqltext
}

func (this *Oraclecon) Get_new_add() map[string]string {
	fd_list, ok := G_dbtables[this.Db_name+this.Tablename]
	if ok {
		//fmt.Println(fd_list)
		//fmt.Println(reflect.TypeOf(fd_list))
		result := make(map[string]string)
		for _, v := range fd_list.([]map[string]string) {
			fd_name := v["column_name"]
			result[fd_name] = ""
		}
		return result
	} else {
		this.Update_redis(this.Tablename)
		rows, _ := this.Masterdb.Query(`select  t.*,c.COMMENTS 
		from user_tab_columns  t,user_col_comments  c
		where t.table_name = c.table_name and t.column_name = c.column_name and t.table_name ='` + this.Tablename + "'")
		defer rows.Close()
		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		result := make(map[string]string)
		for i := range values {
			scanArgs[i] = &values[i]
		}
		for rows.Next() {
			//将行数据保存到record字典
			record := make(map[string]string)
			_ = rows.Scan(scanArgs...)
			for i, col := range values {
				if col != nil {
					record[strings.ToLower(columns[i])] = string(col.([]byte))
					result[record["column_name"]] = ""
				}
			}
		}

		return result
	}
}

func (this *Oraclecon) Update_redis(tbname string) {
	list := this.Query(`select  t.*,c.COMMENTS 
from user_tab_columns  t,user_col_comments  c
 where t.table_name = c.table_name and t.column_name = c.column_name and t.table_name ='`+tbname+"'", nil)
	if list != nil {
		data_list := make([]map[string]string, 0)
		for _, val := range list {
			col := make(map[string]string)
			for key, _ := range val {
				col[strings.ToLower(key)] = val[key]
			}
			data_list = append(data_list, col)
			//fmt.Println(val)
		}
		G_dbtables[this.Db_name+tbname] = data_list
	}

	//fmt.Println(G_dbtables)

}

func (this *Oraclecon) Get_fun_fields(fd_name, val_name string) (result string, flag bool) {
	fd_list, ok := G_dbtables[this.Db_name+this.Tablename]
	flag = false
	if ok {
		for _, v := range fd_list.([]map[string]string) {
			record := v
			if fd_name == record["column_name"] {
				if record["data_type"] == "DATE" {
					result = "to_date('" + val_name + "','yyyy-mm-dd hh24:mi:ss')"
					flag = true
				} else {
					result = ":" + fd_name
				}
				break
			}
		}
	}

	return result, flag

}

func (this *Oraclecon) Get_fields_sql(fd_name, val_name string) (result string) {
	fd_list, ok := G_dbtables[this.Db_name+this.Tablename]
	if ok {
		for _, v := range fd_list.([]map[string]string) {
			record := v
			if fd_name == record["column_name"] {
				if record["data_type"] == "DATE" {
					result = "\"" + record["column_name"] + "\"=to_date('" + val_name + "','yyyy-mm-dd hh24:mi:ss')"
				} else {
					result = "\"" + record["column_name"] + "\"=" + this.checkstr(record["data_type"], val_name)
				}
				break
			}
		}
	}

	return result

}

func (this *Oraclecon) checkstr(fdtype string, fdvalue string) string {
	if fdvalue == "" {
		return "null"
	}
	flag := false
	var fd_list = [...]string{"CHAR"}
	for _, val := range fd_list {
		if strings.Contains(fdtype, val) {
			flag = true
			break
		}
	}

	if flag {
		result := "'" + strings.Replace(fdvalue, "'", "\\'", -1) + "'"
		return result
	} else {
		//result :=strings.Replace(fdvalue, "\\", "\\\\", -1)
		//result = "'" + strings.Replace(result, "'", "\\'", -1) + "'"
		return fdvalue
	}

}

func (this *Oraclecon) Get_where_data(postdata map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, val := range postdata {
		val_str := strings.TrimSpace(this.Type2str(val))
		if val_str != "" {
			if strings.Contains(key, "S_") {
				key1 := strings.Replace(key, "S_", "", -1)
				result[key1] = val_str
			}

			if strings.Contains(key, "I_") {
				key1 := strings.Replace(key, "I_", "", -1)
				result[key1] = map[string]interface{}{"name": "\"" + this.Tablename + "\".\"" + key1 + "\" like '%?%'", "value": val_str}
			}
		}
	}
	return (result)
}

func (this *Oraclecon) Rollback() {
	if this.SqlTx == nil {
		return
	}
	this.SqlTx.Rollback()
	this.SqlTx = nil
}

func (this *Oraclecon) Commit() {
	if this.SqlTx == nil {
		return
	}
	this.SqlTx.Commit()
	this.SqlTx = nil
}

func (this *Oraclecon) Init_Redis_Struct() {
	data := this.Query("select * from user_tables", nil)
	for _, v := range data {
		tbname := v["Tables_in_"+this.Db_name]
		list := this.Query(`select  t.*,c.COMMENTS 
from user_tab_columns  t,user_col_comments  c
 where t.table_name = c.table_name and t.column_name = c.column_name and t.table_name ='`+tbname+"'", nil)
		if list != nil {
			data_list := make([]map[string]string, 0)
			for _, val := range list {
				col := make(map[string]string)
				for key, _ := range val {
					col[strings.ToLower(key)] = val[key]
				}

				data_list = append(data_list, col)
			}
			G_dbtables[this.Db_name+tbname] = data_list
			tbname = strings.Replace(tbname, this.Db_perfix, "", -1)
		}
	}

}
