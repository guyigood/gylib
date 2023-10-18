package mysqllib

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/guyigood/gylib/common"
	"github.com/guyigood/gylib/common/datatype"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type MysqlIni struct {
	DbUser     string   `json:"db_user"`
	DbPassword string   `json:"db_password"`
	DbHost     string   `json:"db_host"`
	DbPort     string   `json:"db_port"`
	DbName     string   `json:"db_name"`
	DbMaxpool  int      `json:"db_maxpool"`
	DbMinpool  int      `json:"db_minpool"`
	DbPerfix   string   `json:"db_perfix"`
	Slavedb    []string `json:"slavedb"`
	Maxtime    int      `json:"maxtime"`
}

type Mysqlcon struct {
	Masterdb    *sql.DB
	Slavedb     []*sql.DB
	Db_host     string
	Db_port     string
	Db_name     string
	Db_password string
	SqlTx       *sql.Tx
	Slock       sync.RWMutex
	Sql_param   []interface{}
	Tablename   string
	Sql_where   string
	Sql_order   string
	Sql_fields  string
	Sql_limit   string
	Db_perfix   string
	PRK_editfd  string
	Query_data  []map[string]interface{}
	Join_arr    map[string]string
	LastSqltext string
	ActionName  string
}

var G_Dbcon *Mysqlcon
var G_MysqlIni MysqlIni
var tb_lock, fd_lock sync.RWMutex
var G_dbtables map[string]interface{}
var G_fd_list map[string]interface{}
var G_tb_dict map[string]interface{}
var G_fd_dict map[string]interface{}

func init() {
	G_dbtables = make(map[string]interface{})
	G_fd_list = make(map[string]interface{})
	G_tb_dict = make(map[string]interface{})
	G_fd_dict = make(map[string]interface{})
}

func NewMySql_Server_DB() *Mysqlcon {
	that := new(Mysqlcon)

	con := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", G_MysqlIni.DbUser,
		G_MysqlIni.DbPassword, G_MysqlIni.DbHost,
		G_MysqlIni.DbPort, G_MysqlIni.DbName)
	//fmt.Println(con)
	var err error
	that.Masterdb, err = sql.Open("mysql", con)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	that.Masterdb.SetMaxOpenConns(G_MysqlIni.DbMaxpool)
	that.Masterdb.SetMaxIdleConns(G_MysqlIni.DbMinpool)
	that.Masterdb.SetConnMaxLifetime(time.Minute * time.Duration(G_MysqlIni.Maxtime))
	err = that.Masterdb.Ping()
	if err != nil {
		fmt.Println("PING", err)
		return nil
	}
	that.Slavedb = make([]*sql.DB, 0)
	if len(G_MysqlIni.Slavedb) > 0 {

		for i := 0; i < len(G_MysqlIni.Slavedb); i++ {
			con1 := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", G_MysqlIni.DbUser,
				G_MysqlIni.DbPassword, G_MysqlIni.Slavedb[i],
				G_MysqlIni.DbPort, G_MysqlIni.DbName)
			sqldb1, _ := sql.Open("mysql", con1)
			sqldb1.SetMaxOpenConns(G_MysqlIni.DbMaxpool)
			sqldb1.SetMaxIdleConns(G_MysqlIni.DbMaxpool)
			sqldb1.SetConnMaxLifetime(time.Minute * time.Duration(G_MysqlIni.Maxtime))
			sqldb1.Ping()
			that.Slavedb = append(that.Slavedb, sqldb1)
		}
	}
	that.SqlTx = nil
	that.Db_perfix = G_MysqlIni.DbPerfix
	that.Db_name = G_MysqlIni.DbName
	that.Db_host = G_MysqlIni.DbHost
	that.Db_port = G_MysqlIni.DbPort
	that.Db_password = G_MysqlIni.DbPassword
	//go that.CheckConnectStatus()
	return that
}

func NewMysqlDB(jsonstr string) *Mysqlcon {
	err1 := json.Unmarshal([]byte(jsonstr), &MysqlIni{})
	if err1 != nil {
		return nil
	}
	if G_Dbcon == nil {
		G_Dbcon = NewMySql_Server_DB()
		G_Dbcon.Init_Redis_Struct()
	}
	that := new(Mysqlcon)
	that.Masterdb = G_Dbcon.Masterdb
	that.Slavedb = G_Dbcon.Slavedb
	that.Db_perfix = G_Dbcon.Db_perfix
	that.Db_name = G_Dbcon.Db_name
	that.Db_host = G_Dbcon.Db_host
	that.Db_port = G_Dbcon.Db_port
	that.Db_password = G_Dbcon.Db_password
	//that.Masterdb=Dbcon.Masterdb
	//that.Slavedb=Dbcon.Slavedb
	//that.Db_perfix = Dbcon.Db_perfix
	//that.Db_name = Dbcon.Db_name
	//that.Db_host = Dbcon.Db_host
	//that.Db_port = Dbcon.Db_port
	//that.Db_password = Dbcon.Db_password
	that.SqlTx = nil
	that.Dbinit()
	return that
}

func (that *Mysqlcon) CheckConnectStatus() {
	for {
		err := that.Masterdb.Ping()
		if err != nil {
			//fmt.Println(err)
			NewMySql_Server_DB()
			break
		}
		time.Sleep(time.Millisecond)
	}

}

func (that *Mysqlcon) Merge_And_where(where_str, new_str string) string {
	result := where_str
	if where_str != "" {
		result += " and " + new_str
	} else {
		result = new_str
	}
	return result
}

func (that *Mysqlcon) Merge_OR_where(where_str, new_str string) string {
	result := where_str
	if where_str != "" {
		result += " or " + new_str
	} else {
		result = new_str
	}
	return result
}

func (that *Mysqlcon) BeginStart() bool {
	tx, err := that.Masterdb.Begin()
	if err != nil {
		return false
	}
	that.SqlTx = tx
	return true
}

/*
*
初始化结构
*/
func (that *Mysqlcon) Dbinit() {
	that.Tablename = ""
	that.Sql_limit = ""
	that.Sql_order = ""
	that.Sql_fields = ""
	that.Sql_where = ""
	that.Slock.Lock()
	that.Join_arr = make(map[string]string)
	that.Query_data = make([]map[string]interface{}, 0)
	that.Sql_param = make([]interface{}, 0)
	that.Slock.Unlock()
}

/*
设置数据表
*/
func (that *Mysqlcon) Tbname(name string) *Mysqlcon {
	//data := common.Getini("conf/app.ini","database",map[string]string{"db_perfix":""})
	that.Dbinit()
	that.Tablename = that.Db_perfix + name
	return that
}

func (that *Mysqlcon) Fileds(name string) *Mysqlcon {
	that.Sql_fields = name
	return that
}

func (that *Mysqlcon) SetWhere(where string, param ...interface{}) *Mysqlcon {
	for _, val := range param {
		that.Sql_param = append(that.Sql_param, val)
	}
	if that.Sql_where == "" {
		that.Sql_where = where
	} else {
		that.Sql_where += " and (" + where + ")"
	}
	return that
}

func (that *Mysqlcon) Where(where interface{}) *Mysqlcon {
	//kk:= reflect.TypeOf(where)
	//fmt.Println(kk)
	if where == nil {
		return that
	}
	switch where.(type) {
	case string:
		if datatype.Type2str(where) == "" {
			return that
		}
		if that.Sql_where == "" {
			that.Sql_where = where.(string)
		} else {
			that.Sql_where += " and (" + where.(string) + ")"
		}
	default:
		that.Slock.Lock()
		tmp_arr := where.(map[string]interface{})
		if len(tmp_arr) > 0 {
			that.Query_data = append(that.Query_data, tmp_arr)
		}
		that.Slock.Unlock()
		//fmt.Println("query_data", that.Query_data)
	}

	return that
}

func (that *Mysqlcon) Order(orderstr string) *Mysqlcon {
	that.Sql_order = orderstr
	return that
}

func (that *Mysqlcon) Limit(limitstr string) *Mysqlcon {
	that.Sql_limit = limitstr
	return that
}

func (that *Mysqlcon) Get_read_dbcon() *sql.DB {
	read_ct := len(that.Slavedb)
	if read_ct == 0 {
		return that.Masterdb
	} else {
		result := rand.Intn(read_ct)
		return that.Slavedb[result]
	}
}

func (that *Mysqlcon) Check_data_fields(fieldname string) bool {
	if that.Check_PK(fieldname) {
		return false
	}
	flag := false
	tb_lock.Lock()
	defer tb_lock.Unlock()
	fd_list, ok := G_dbtables[that.Db_name+that.Tablename]
	if ok {
		for _, v := range fd_list.([]map[string]string) {
			record := v
			if record["key"] == "PRI" && record["extra"] == "auto_increment" {
				continue
			}

			if record["field"] == fieldname {
				flag = true
				break
			}

		}
		return flag
	} else {
		that.Update_redis(that.Tablename)
		rows, _ := that.Masterdb.Query("SHOW full COLUMNS FROM " + that.Tablename)
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
					record[strings.ToLower(columns[i])] = that.Type2str(col)
				}
			}
			if record["key"] == "PRI" && record["extra"] == "auto_increment" {
				continue
			}

			if record["field"] == fieldname {
				flag = true
				break
			}

		}
		return flag
	}
}

func (that *Mysqlcon) Type2str(val interface{}) string {
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

func (that *Mysqlcon) Insert(postdata map[string]interface{}) (sql.Result, error) {
	//that.Wlock.Lock()
	//defer that.Wlock.Unlock()
	var sqltext string
	sqltext = "insert into " + that.Tablename + " ("
	values := " values ("
	i := 0
	param_data := make([]interface{}, 0)
	for k, v := range postdata {
		if that.Check_data_fields(k) == false {
			continue
		}
		if i > 0 {
			sqltext += ","
			values += ","
		}
		i++
		sqltext += "`" + k + "`"
		values += " ? "
		if datatype.Type2str(v) != "" {
			param_data = append(param_data, v)
		} else {
			param_data = append(param_data, nil)
		}
	}
	sqltext += ") " + values + ")"
	that.LastSqltext = sqltext
	//fmt.Println(i,sqltext)
	//fmt.Println(len(param_data),param_data)
	var result sql.Result
	var err error
	if that.SqlTx != nil {
		result, err = that.SqlTx.Exec(sqltext, param_data...)
	} else {
		result, err = that.Masterdb.Exec(sqltext, param_data...)
	}
	//fmt.Println(err)
	return result, err

}

func (that *Mysqlcon) Update(postdata map[string]interface{}) (sql.Result, error) {
	//that.Wlock.Lock()
	//defer that.Wlock.Unlock()
	var sqltext string
	sqltext = fmt.Sprintf("update %v set ", that.Tablename)
	i := 0
	param_data := make([]interface{}, 0)
	for k, v := range postdata {
		if that.Check_data_fields(k) == false {
			continue
		}
		if i > 0 {
			sqltext += ","

		}
		i++
		sqltext += "`" + k + "`" + "= ?"
		if datatype.Type2str(v) != "" {
			param_data = append(param_data, v)
		} else {
			param_data = append(param_data, nil)
		}
	}
	sqlwhere, param := that.Build_where()
	for _, v := range param {
		param_data = append(param_data, v)
	}
	sqltext += sqlwhere
	that.LastSqltext = sqltext
	//fmt.Println(sqltext, param_data)
	var result sql.Result
	var err error
	if that.SqlTx != nil {
		result, err = that.SqlTx.Exec(sqltext, param_data...)
	} else {
		result, err = that.Masterdb.Exec(sqltext, param_data...)
	} //fmt.Println(err)
	return result, err
}

func (that *Mysqlcon) Delete() (sql.Result, error) {
	//that.Wlock.Lock()
	//defer that.Wlock.Unlock()
	sqlwhere, param := that.Build_where()
	sqltext := fmt.Sprintf(" delete from %v %v", that.Tablename, sqlwhere)
	that.LastSqltext = sqltext
	var result sql.Result
	var err error
	if that.SqlTx != nil {
		result, err = that.SqlTx.Exec(sqltext, param...)
	} else {
		result, err = that.Masterdb.Exec(sqltext, param...)
	}
	return result, err
}

func (that *Mysqlcon) SetDec(fdname string, quantity int) (sql.Result, error) {
	sqltext := fmt.Sprintf("update %v set %v=%v-%v", that.Tablename, fdname, fdname, quantity)
	sqlwhere, param := that.Build_where()
	sqltext += sqlwhere
	that.LastSqltext = sqltext
	var result sql.Result
	var err error
	if that.SqlTx != nil {
		result, err = that.SqlTx.Exec(sqltext, param...)
	} else {
		result, err = that.Masterdb.Exec(sqltext, param...)
	}
	return result, err
}

func (that *Mysqlcon) SetInc(fdname string, quantity int) (sql.Result, error) {
	sqlwhere, param := that.Build_where()
	sqltext := fmt.Sprintf("update %v set %v=%v+%v  %v", that.Tablename, fdname, fdname, quantity, sqlwhere)
	that.LastSqltext = sqltext
	var result sql.Result
	var err error
	if that.SqlTx != nil {
		result, err = that.SqlTx.Exec(sqltext, param...)
	} else {
		result, err = that.Masterdb.Exec(sqltext, param...)
	}
	return result, err
}

func (that *Mysqlcon) Query(sqltext string, param []interface{}) []map[string]string {
	that.LastSqltext = sqltext
	//fmt.Println(sqltext)
	var rows *sql.Rows
	var err error
	if that.SqlTx != nil {
		rows, err = that.SqlTx.Query(sqltext, param...)
	} else {
		sqldbcon := that.Get_read_dbcon()
		rows, err = sqldbcon.Query(sqltext, param...)
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
	return result
}

func (that *Mysqlcon) Query_One(sqltext string, param []interface{}) map[string]string {
	that.LastSqltext = sqltext
	var rows *sql.Rows
	var err error
	if that.SqlTx != nil {
		rows, err = that.SqlTx.Query(sqltext, param...)
	} else {
		sqldbcon := that.Get_read_dbcon()
		rows, err = sqldbcon.Query(sqltext, param...)
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

func (that *Mysqlcon) Excute(sqltext string, param []interface{}) (sql.Result, error) {
	that.LastSqltext = sqltext
	var result sql.Result
	var err error
	if that.SqlTx != nil {
		result, err = that.SqlTx.Exec(sqltext, param...)
	} else {
		result, err = that.Masterdb.Exec(sqltext, param...)
	}
	return result, err
}

// func (that *Mysqlcon) delete(fields ...string) (*Mysqlcon) {
// that.Tokens = append(that.Tokens, "SELECT", strings.Join(fields,","))
//
//		return qb
//	}
func (that *Mysqlcon) Join(tbname string, jointype string, where string, fileds string) *Mysqlcon {
	that.Slock.Lock()
	defer that.Slock.Unlock()
	if that.Join_arr["tbname"] == "" {
		that.Join_arr["tbname"] = that.Tablename + " " + jointype + " " + that.Db_perfix + tbname + " on " + where
		if fileds != "" {
			that.Join_arr["fields"] = that.Tablename + ".*," + fileds
		} else {
			that.Join_arr["fields"] = that.Tablename + ".*"
		}
	} else {
		that.Join_arr["tbname"] += " " + jointype + " " + that.Db_perfix + tbname + " on " + where
		if fileds != "" {
			that.Join_arr["fields"] += "," + fileds
		}
	}
	// that.Slock.Unlock()
	return that
}

func (that *Mysqlcon) set_sql(flag int) string {
	that.Slock.RLock()
	defer that.Slock.RUnlock()
	sqltext := ""
	if flag == 0 {
		if datatype.Has_map_index("tbname", datatype.MapString2interface(that.Join_arr)) {
			if that.Join_arr["fields"] != "" {
				sqltext = "select " + that.Join_arr["fields"] + " from " + that.Join_arr["tbname"]
			} else {
				if that.Sql_fields != "" {
					sqltext = "select " + that.Sql_fields + " from " + that.Tablename
				} else {
					sqltext = "select " + that.Tablename + ".* from " + that.Tablename
				}
			}
		} else {
			sqltext = "select  * from " + that.Tablename
		}
	} else {
		if datatype.Has_map_index("tbname", datatype.MapString2interface(that.Join_arr)) {
			sqltext = "select count(" + that.Tablename + ".*) as ct " + " from " + that.Join_arr["tbname"]
		} else {
			sqltext = "select count(*) as ct from " + that.Tablename
		}
	}
	return sqltext
}

func (that *Mysqlcon) Build_where() (string, []interface{}) {
	is_where := false
	sqltext := ""
	if that.Sql_where != "" {
		sqltext += " where " + that.Sql_where
		is_where = true
	}
	param_data := make([]interface{}, 0)
	if len(that.Query_data) > 0 {
		if is_where {
			sqltext += " and "

		} else {
			sqltext += " where "
		}
		i := 0
		that.Slock.RLock()
		for _, v := range that.Query_data {
			for key, val := range v {
				//if (that.Check_data_fields(key) == false) {
				//	continue
				//}
				if i > 0 {
					sqltext += " and "
				}
				i++
				switch val.(type) {
				//data["name"]=" %v like ?"
				//data["name"]=" %v>=(?)"
				//data["name"]="locate(?,`"+that.Tablename+"`.`%v`)>0"
				case map[string]interface{}:
					param_data = append(param_data, val.(map[string]interface{})["value"])
					sqltext += datatype.Type2str(val.(map[string]interface{})["name"])
				default:
					param_data = append(param_data, val)
					sqltext += key + "=(?) "

				}
			}
		}
		that.Slock.RUnlock()

	}
	if len(that.Sql_param) > 0 {
		for _, val := range that.Sql_param {
			param_data = append(param_data, val)
		}
	}

	return sqltext, param_data
}

func (that *Mysqlcon) Find() map[string]string {
	sqltext := that.set_sql(0)
	param_data := make([]interface{}, 0)
	tmpstr := ""
	tmpstr, param_data = that.Build_where()
	sqltext += tmpstr
	if that.Sql_order != "" {
		sqltext += " order by " + that.Sql_order
	}
	that.LastSqltext = sqltext + " limit 1"
	var rows *sql.Rows
	var err error
	if that.SqlTx != nil {
		rows, err = that.SqlTx.Query(sqltext+" limit 1", param_data...)
	} else {
		sqldbcon := that.Get_read_dbcon()
		rows, err = sqldbcon.Query(sqltext+" limit 1", param_data...)
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

func (that *Mysqlcon) Count() int64 {
	sqltext := that.set_sql(1)
	sqlwhere, param := that.Build_where()
	sqltext += sqlwhere
	that.LastSqltext = sqltext
	sqldbcon := that.Get_read_dbcon()
	rows := sqldbcon.QueryRow(sqltext, param...)
	var record int64
	rows.Scan(&record)

	return record
}

func (that *Mysqlcon) Sum(fd string) float64 {
	var result float64
	sqltext := that.set_sql(1)
	sqltext = strings.Replace(sqltext, "count(*)", "IFNULL(sum("+fd+"),0)", -1)
	sqlwhere, param := that.Build_where()
	sqltext += sqlwhere
	that.LastSqltext = sqltext
	sqldbcon := that.Get_read_dbcon()
	rows := sqldbcon.QueryRow(sqltext, param...)
	rows.Scan(&result)
	return result
}

func (that *Mysqlcon) Max(fd string) float64 {
	var result float64
	sqltext := that.set_sql(1)
	sqltext = strings.Replace(sqltext, "count(*)", "IFNULL(max("+fd+"),0)", -1)
	sqlwhere, param := that.Build_where()
	sqltext += sqlwhere
	that.LastSqltext = sqltext
	sqldbcon := that.Get_read_dbcon()
	rows := sqldbcon.QueryRow(sqltext, param...)
	rows.Scan(&result)
	return result
}

func (that *Mysqlcon) IMax(fd string) int64 {
	var result int64
	sqltext := that.set_sql(1)
	sqltext = strings.Replace(sqltext, "count(*)", "IFNULL(max("+fd+"),0)", -1)
	sqlwhere, param := that.Build_where()
	sqltext += sqlwhere
	that.LastSqltext = sqltext
	sqldbcon := that.Get_read_dbcon()
	rows := sqldbcon.QueryRow(sqltext, param...)
	rows.Scan(&result)
	return result
}
func (that *Mysqlcon) IMin(fd string) int64 {
	var result int64
	sqltext := that.set_sql(1)
	sqltext = strings.Replace(sqltext, "count(*)", "IFNULL(min("+fd+"),0)", -1)
	sqlwhere, param := that.Build_where()
	sqltext += sqlwhere
	that.LastSqltext = sqltext
	sqldbcon := that.Get_read_dbcon()
	rows := sqldbcon.QueryRow(sqltext, param...)
	rows.Scan(&result)
	return result
}

func (that *Mysqlcon) Min(fd string) float64 {
	var result float64
	sqltext := that.set_sql(1)
	sqltext = strings.Replace(sqltext, "count(*)", "IFNULL(min("+fd+"),0)", -1)
	sqlwhere, param := that.Build_where()
	sqltext += sqlwhere
	that.LastSqltext = sqltext
	sqldbcon := that.Get_read_dbcon()
	rows := sqldbcon.QueryRow(sqltext, param...)
	rows.Scan(&result)
	return result
}

func (that *Mysqlcon) Select() []map[string]string {
	sqltext := that.set_sql(0)
	sqlwhere, param := that.Build_where()
	sqltext += sqlwhere
	if that.Sql_order != "" {
		sqltext += " order by " + that.Sql_order
	}
	if that.Sql_limit != "" {
		sqltext += " limit " + that.Sql_limit
	}
	that.LastSqltext = sqltext
	//fmt.Println(sqltext)
	sqldbcon := that.Get_read_dbcon()
	rows, err := sqldbcon.Query(sqltext, param...)
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
				record[columns[i]] = that.Type2str(col)
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

func (that *Mysqlcon) GetLastSql() string {
	return that.LastSqltext
}

func (that *Mysqlcon) SetPK(pkfd string) *Mysqlcon {
	that.PRK_editfd = pkfd
	return that
}

func (that *Mysqlcon) Check_PK(fdname string) bool {
	if that.PRK_editfd == "" {
		return false
	}
	list := strings.Split(that.PRK_editfd, ",")
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

func (that *Mysqlcon) Get_new_add() map[string]string {
	tb_lock.Lock()
	defer tb_lock.Unlock()
	fd_list, ok := G_dbtables[that.Db_name+that.Tablename]
	if ok {
		//fmt.Println(fd_list)
		//fmt.Println(reflect.TypeOf(fd_list))
		result := make(map[string]string)
		for _, v := range fd_list.([]map[string]string) {
			fd_name := v["field"]
			result[fd_name] = ""
		}
		return result
	} else {
		that.Update_redis(that.Tablename)
		rows, _ := that.Masterdb.Query("SHOW full COLUMNS FROM " + that.Tablename)
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
					result[record["field"]] = ""
				}
			}
		}

		return result
	}
}

func (that *Mysqlcon) Update_redis(tbname string) {
	list := that.Query("SHOW full COLUMNS FROM "+tbname, nil)
	if list != nil {
		data_list := make([]map[string]string, 0)
		for _, val := range list {
			col := make(map[string]string)
			for key, _ := range val {
				col[common.Tolow_map_name(key)] = val[key]
			}
			data_list = append(data_list, col)
		}
		G_dbtables[that.Db_name+tbname] = data_list
	}

}

func (that *Mysqlcon) Get_fields_sql(fd_name, val_name string) (result string) {
	tb_lock.Lock()
	defer tb_lock.Unlock()
	fd_list, ok := G_dbtables[that.Db_name+that.Tablename]
	if ok {
		for _, v := range fd_list.([]map[string]string) {
			record := v
			if fd_name == record["field"] {
				result = "`" + record["field"] + "`=" + that.checkstr(record["type"], val_name)
				break
			}
		}
	}

	return result

}

func (that *Mysqlcon) checkstr(fdtype string, fdvalue string) string {
	if fdvalue == "" {
		return "null"
	}
	flag := false
	var fd_list = [...]string{"date", "time", "datetime"}
	for _, val := range fd_list {
		if strings.ToLower(fdtype) == val {
			flag = true
			break
		}
	}
	if !flag {
		var fd_list = [...]string{"char", "text", "linestring"}
		for _, val := range fd_list {
			if strings.Contains(fdtype, val) {
				flag = true
				break
			}
		}
	}

	if flag {
		result := "'" + strings.Replace(fdvalue, "'", "\\'", -1) + "'"
		return result
	} else {
		return fdvalue
	}

	/*if (strings.Contains(fdtype, "tinyint") ||
		strings.Contains(fdtype, "double") ||
		strings.Contains(fdtype, "float") ||
		strings.Contains(fdtype, "int") ||
		strings.Contains(fdtype, "decimal")) {
		return fdvalue
	} else {
		//result :=strings.Replace(fdvalue, "\\", "\\\\", -1)
		//result = "'" + strings.Replace(result, "'", "\\'", -1) + "'"
		result := "'" + strings.Replace(fdvalue, "'", "\\'", -1) + "'"
		return result
	}*/

}

//func (that *Mysqlcon) Get_select_data(d_data map[string]string, masterdb string) (map[string]string) {
//	client := rediscomm.NewRedisComm()
//	client.Key = "fd_list"
//	client.Field = masterdb
//	data := client.Hget_map()
//	if (data != nil) {
//		for _, v := range data.([]interface{}) {
//			listname := strings.Replace(v.(map[string]interface{})["list_tb_name"].(string), that.Db_perfix, "", -1)
//			tbname := strings.Replace(v.(map[string]interface{})["list_tb_name"].(string), that.Db_perfix, "", -1)
//			listname = strings.Replace(listname, "_", "", -1)
//			where := v.(map[string]interface{})["list_where"].(string)
//			list_val := v.(map[string]interface{})["list_val"].(string)
//			list_display := datatype.Type2str(v.(map[string]interface{})["list_display"])
//			if (where != "") {
//				where += " and " + that.Tbname(tbname).Get_fields_sql(list_val, d_data[v.(map[string]interface{})["name"].(string)])
//			} else {
//				where = that.Tbname(tbname).Get_fields_sql(list_val, d_data[v.(map[string]interface{})["name"].(string)])
//			}
//			list_data := that.Tbname(tbname).Where(where).Find()
//			//fmt.Println(v,that.GetLastSql())
//			//fmt.Println(list_data)
//			if (list_data != nil) {
//				d_data[v.(map[string]interface{})["name"].(string)+"_name"] = list_data[list_display]
//			} else {
//				d_data[v.(map[string]interface{})["name"].(string)+"_name"] = ""
//			}
//		}
//	}
//	//fmt.Println(d_data)
//	return d_data
//}

func (that *Mysqlcon) GetWherePostFrom(postdata map[string]interface{}, masterdb string) map[string]interface{} {
	fd_lock.Lock()
	defer fd_lock.Unlock()
	data, ok := G_fd_list[that.Db_name+masterdb].([]map[string]string)
	if !ok {
		that.Get_mysql_dict(masterdb)
		data, ok = G_fd_list[that.Db_name+masterdb].([]map[string]string)
		if !ok {
			return nil
		}
	}
	result := make(map[string]interface{})
	for key, val := range postdata {
		val_str := strings.TrimSpace(that.Type2str(val))
		if val_str != "" {
			for i := 0; i < len(data); i++ {
				if data[i]["name"] == key {
					if data[i]["f_like"] == "1" {
						result[key] = map[string]interface{}{"name": "locate(?,`" + that.Tablename + "`.`" + key + "`)>0", "value": val_str}
					} else {
						result[key] = val_str
					}
				}
			}
		}
	}
	return (result)
}

func (that *Mysqlcon) Get_where_data(postdata map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, val := range postdata {
		val_str := strings.TrimSpace(that.Type2str(val))
		if val_str != "" {
			if strings.Contains(key, "S_") {
				key1 := strings.Replace(key, "S_", "", -1)
				result[key1] = val_str
			}

			if strings.Contains(key, "I_") {
				key1 := strings.Replace(key, "I_", "", -1)
				result[key1] = map[string]interface{}{"name": "locate(?,`" + that.Tablename + "`.`" + key1 + "`)>0", "value": val_str}
			}
		}
	}
	return (result)
}

func (that *Mysqlcon) Rollback() {
	if that.SqlTx == nil {
		return
	}
	that.SqlTx.Rollback()
	that.SqlTx = nil
}

func (that *Mysqlcon) Commit() {
	if that.SqlTx == nil {
		return
	}
	that.SqlTx.Commit()
	that.SqlTx = nil
}

func (that *Mysqlcon) Init_Redis_Struct() {
	tb_lock.Lock()
	defer tb_lock.Unlock()
	data := that.Query("show TABLES", nil)
	for _, v := range data {
		tbname := v["Tables_in_"+that.Db_name]
		list := that.Query("SHOW full COLUMNS FROM "+tbname, nil)
		if list != nil {
			data_list := make([]map[string]string, 0)
			for _, val := range list {
				col := make(map[string]string)
				for key, _ := range val {
					col[common.Tolow_map_name(key)] = val[key]
				}

				data_list = append(data_list, col)
			}
			G_dbtables[that.Db_name+tbname] = data_list
			tbname = strings.Replace(tbname, that.Db_perfix, "", -1)
			that.Get_mysql_dict(tbname)
		}
	}

}

func (that *Mysqlcon) Get_select_data(d_data map[string]string, masterdb string) map[string]string {
	fd_lock.Lock()
	defer fd_lock.Unlock()
	data, ok := G_fd_list[that.Db_name+masterdb].([]map[string]string)
	if !ok {
		that.Get_mysql_dict(masterdb)
		data, ok = G_fd_list[that.Db_name+masterdb].([]map[string]string)
		if !ok {
			return d_data
		}
	}
	where := make(map[string]interface{})
	for _, v := range data {
		list_val := v["list_val"]
		list_display := datatype.Type2str(v["list_display"])
		if v["list_tb_name"] == "0" { //没有数据源的时候
			val_arr := strings.Split(list_val, "|") //分割为数组
			dis_arr := strings.Split(list_display, "|")
			if len(val_arr) == len(dis_arr) { //值和显示标签数组长度相等
				for i := 0; i < len(val_arr); i++ {
					if d_data[v["name"]] == val_arr[i] {
						d_data[v["name"]+"_name"] = dis_arr[i]
					}
				}
			}
			continue
		}
		tbname := strings.Replace(v["list_tb_name"], that.Db_perfix, "", -1)
		where1 := v["list_where"]
		where[list_val] = d_data[v["name"]]

		list_data := make(map[string]string)
		if where1 != "" {
			list_data = that.Tbname(tbname).Where(where).Where(where1).Find()
		} else {
			list_data = that.Tbname(tbname).Where(where).Find()
		}
		//fmt.Println(that.GetLastSql())
		//fmt.Println(v,that.GetLastSql())
		//fmt.Println(list_data)
		if list_data != nil {
			d_data[v["name"]+"_name"] = list_data[list_display]
		} else {
			d_data[v["name"]+"_name"] = ""
		}
	}

	//fmt.Println(d_data)
	return d_data

}

func (db *Mysqlcon) Get_mysql_dict(tbname string) {
	data := db.Tbname("db_tb_dict").Where(fmt.Sprintf("name='%v'", db.Db_perfix+tbname)).Find()
	if data == nil {
		return
	}
	fd_lock.Lock()
	defer fd_lock.Unlock()
	fd_data := db.Tbname("db_fd_dict").Where(fmt.Sprintf("t_id=%v", data["id"])).Select()
	list_data := db.Tbname("db_fd_dict").Where(fmt.Sprintf("t_id=%v and list_tb_name<>'0'", data["id"])).Select()

	G_tb_dict[db.Db_name+tbname] = data
	if fd_data != nil {
		G_fd_dict[db.Db_name+tbname] = fd_data
	}
	if list_data != nil {
		G_fd_list[db.Db_name+tbname] = list_data
	}

}
