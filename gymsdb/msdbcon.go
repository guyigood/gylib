package gymsdb

import (
	"database/sql"
	"fmt"
	"strconv"
	"gylib/common"
	"gylib/common/datatype"
	"strings"
	"math/rand"
	_ "github.com/denisenkom/go-mssqldb"
	"gylib/public"
	"sync"
	"time"
)

type Mscon struct {
	Db_host     string
	Db_port     string
	Db_name     string
	Db_password string
	Db_perfix   string
	Masterdb    *sql.DB
	Slavedb     []*sql.DB
	SqlTx       *sql.Tx
	Slock       sync.Mutex

	Tablename        string
	Sql_where        string
	Sql_order        string
	Sql_fields       string
	Sql_limit        string
	PageSize, PageNo int
	Query_data       []map[string]interface{}
	Join_arr         map[string]string
	LastSqltext      string
	App_name         string
}


var G_dbtables map[string]interface{}
var G_fd_list map[string]interface{}
var G_tb_dict map[string]interface{}
var G_fd_dict map[string]interface{}



func NewMsSql_Server_Con(action string) (*Mscon) {
	this := new(Mscon)
	if (action == "") {
		action = "database"
	}
	G_dbtables = make(map[string]interface{})
	G_fd_list = make(map[string]interface{})
	G_tb_dict = make(map[string]interface{})
	G_fd_dict = make(map[string]interface{})

	this.App_name = public.App_name
	data := common.Getini("conf/app.ini", action, map[string]string{"db_user": "root", "db_password": "",
		"db_host": "127.0.0.1", "db_port": "3306", "db_name": "", "db_maxpool": "10", "db_minpool": "5", "db_perfix": "", "slavedb": ""})
	//var password = flag.String("password", data["db_password"], "the database password")
	//var port *int = flag.Int("port", datatype.Str2Int(data["db_port"]), "the database port")
	//var server = flag.String("server", data["db_host"], "the database server")
	//var user = flag.String("user", data["db_user"], "the database user")
	//var database = flag.String("database", data["db_name"], "the database name")

	//connString := fmt.Sprintf("server=%v;database=%v;user id=%s;password=%v;port=%v;encrypt=disable", data["db_host"], data["db_name"], data["db_user"], data["db_password"], data["db_port"])
	connString := fmt.Sprintf("server=%v;database=%v;user id=%s;password=%v;port=%v;encrypt=disable", data["db_host"], data["db_name"], data["db_user"], data["db_password"], data["db_port"])
	//connString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%d;encrypt=disable", *server, *database, *user, *password, *port)
	var err error
	this.Masterdb, err = sql.Open("mssql", connString)
	if (err != nil) {
		fmt.Println(err)
	}
	maxpool, _ := strconv.Atoi(data["db_maxpool"])
	minpool, _ := strconv.Atoi(data["db_minpool"])
	fmt.Println(connString)
	this.Masterdb.SetMaxIdleConns(minpool)
	this.Masterdb.SetMaxOpenConns(maxpool)
	this.Masterdb.SetConnMaxLifetime(time.Minute * 5)
	this.Db_perfix = data["db_perfix"]
	this.Db_name = data["db_name"]
	this.Db_password = data["db_password"]
	this.Db_host = data["db_host"]
	this.Db_port = data["db_port"]
	if err != nil {
		fmt.Println("Open connection failed:", err.Error())
		return nil
	}
	err = this.Masterdb.Ping()
	if err != nil {
		fmt.Println("PING:", err)
	}
	this.Slavedb = make([]*sql.DB, 0)
	if (data["slavedb"] != "") {
		iplist := strings.Split(data["slavedb"], ",")
		for _, v := range iplist {
			con1 := fmt.Sprintf("server=%v;database=%v;user id=%v;password=%v;port=%v;encrypt=disable", v, data["db_name"], data["db_name"], data["db_password"], data["db_port"])
			sqldb1, _ := sql.Open("mysql", con1)
			//maxpool, _ := strconv.Atoi(data["db_maxpool"])
			//minpool, _ := strconv.Atoi(data["db_minpool"])
			sqldb1.SetMaxOpenConns(maxpool)
			sqldb1.SetMaxIdleConns(minpool)
			sqldb1.SetConnMaxLifetime(time.Minute * 5)
			sqldb1.Ping()
			this.Slavedb = append(this.Slavedb, sqldb1)
		}
	}
	this.SqlTx = nil
	return this
}


func (this *Mscon) BeginStart() (bool) {
	tx, err := this.Masterdb.Begin()
	if (err != nil) {
		return false
	}
	this.SqlTx = tx
	return true
}

func (this *Mscon) Dbinit() {
	this.Slock.Lock()
	defer this.Slock.Unlock()
	this.Tablename = ""
	this.Sql_limit = ""
	this.Sql_order = ""
	this.Sql_fields = ""
	this.Sql_where = ""
	this.PageNo = 0
	this.PageSize = 0
	this.Join_arr = make(map[string]string)
	this.Query_data = make([]map[string]interface{}, 0)
}

func (this *Mscon) MapContains_str(src map[string]string, key string) bool {
	if _, ok := src[key]; ok {
		return true
	}
	return false
}

func (this *Mscon) set_sql(flag int) (string) {
	sqltext := ""
	if (flag == 0) {
		if (this.MapContains_str(this.Join_arr, "tbname")) {
			if (this.Join_arr["fields"] != "") {
				sqltext = "select " + this.Join_arr["fields"] + " from " + this.Join_arr["tbname"]
			} else {
				sqltext = "select " + this.Tablename + ".* from " + this.Tablename
			}
		} else {
			sqltext = "select  * from " + this.Tablename
		}
	} else {
		if (this.MapContains_str(this.Join_arr, "tbname")) {
			sqltext = "select count(" + this.Tablename + ".*) as ct " + " from " + this.Join_arr["tbname"]
		} else {
			sqltext = "select count(*) as ct from " + this.Tablename
		}
	}
	return sqltext
}

func (this *Mscon) Find() map[string]string {
	sqltext := this.set_sql(0)
	param_data := make([]interface{}, 0)
	tmpstr := ""
	tmpstr, param_data = this.Build_where()
	sqltext += tmpstr
	if this.Sql_order != "" {
		sqltext += " order by " + this.Sql_order
	}
	this.LastSqltext = strings.Replace(sqltext, "select", "select top 1 ", -1)
	sqltext = this.LastSqltext
	var rows *sql.Rows
	var err error
	if (this.SqlTx != nil) {
		rows, err = this.SqlTx.Query(sqltext, param_data...)
	} else {
		sqldbcon := this.Get_read_dbcon()
		rows, err = sqldbcon.Query(sqltext, param_data...)
	}
	//fmt.Println("rows",rows,err)
	if (err != nil) {
		fmt.Println(this.LastSqltext)
		return nil
	}
	if (rows == nil) {
		return nil
	}

	defer rows.Close()
	//????????????
	//??????scanArgs???values???????????????scanArgs??????????????????values??????????????????
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	record := make(map[string]string)
	for rows.Next() {
		//?????????????????????record??????
		_ = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = datatype.Type2str(col)

			} else {
				record[columns[i]] = ""
			}
		}

	}
	if (len(record) == 0) {
		return nil
	}
	return record
}

func (this *Mscon) Get_read_dbcon() (*sql.DB) {
	read_ct := len(this.Slavedb)
	//fmt.Println("read_ct",read_ct)
	if (read_ct == 0) {
		return this.Masterdb
	} else {
		result := rand.Intn(read_ct)
		//fmt.Println("readcon",result)
		return this.Slavedb[result]

	}
}

func (this *Mscon) Build_where() (string, []interface{}) {
	is_where := false
	sqltext := ""
	if this.Sql_where != "" {
		sqltext += " where " + this.Sql_where
		is_where = true
	}
	param_data := make([]interface{}, 0)
	if (len(this.Query_data) > 0) {
		if (is_where) {
			sqltext += " and "

		} else {
			sqltext += " where "
		}
		i := 0
		for _, v := range this.Query_data {
			for key, val := range v {
				//if (this.Check_data_fields(key) == false) {
				//	continue
				//}
				if (i > 0) {
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
					sqltext += key + "=(?) "

				}
			}
		}

	}
	return sqltext, param_data
}

func (this *Mscon) Tbname(name string) (*Mscon) {
	//data := common.Getini("conf/app.ini","database",map[string]string{"db_perfix":""})
	this.Dbinit()
	this.Tablename = this.Db_perfix + name
	return this
}

func (this *Mscon) PageSet(pagesize, pageno int) (*Mscon) {
	//data := common.Getini("conf/app.ini","database",map[string]string{"db_perfix":""})
	this.PageSize = pagesize
	this.PageNo = pageno
	return this
}

func (this *Mscon) Where(where interface{}) (*Mscon) {
	//kk:= reflect.TypeOf(where)
	//fmt.Println(kk)
	switch where.(type) {
	case string:
		if (datatype.Type2str(where) == "") {
			return this
		}
		if this.Sql_where == "" {
			this.Sql_where = where.(string)
		} else {
			this.Sql_where += " and (" + where.(string) + ")"
		}
	default:
		this.Slock.Lock()
		defer this.Slock.Unlock()
		tmp_arr := where.(map[string]interface{})
		if (len(tmp_arr) > 0) {
			this.Query_data = append(this.Query_data, tmp_arr)
		}
		//fmt.Println("query_data", this.Query_data)
	}

	return this
}

func (this *Mscon) Order(orderstr string) (*Mscon) {
	this.Sql_order = orderstr
	return this
}

func (this *Mscon) Limit(limitstr string) (*Mscon) {
	this.Sql_limit = limitstr
	return this
}

func (this *Mscon) MapContains(src map[string]interface{}, key string) bool {
	if _, ok := src[key]; ok {
		return true
	}
	return false
}

func (this *Mscon) Insert(postdata map[string]interface{}) (sql.Result, error) {
	var sqltext string
	sqltext = "insert into " + this.Tablename + " ("
	values := " values ("
	i := 0
	param_data := make([]interface{}, 0)
	for k, v := range postdata {
		if (this.Check_data_fields(k) == false) {
			continue
		}
		if (i > 0) {
			sqltext += ","
			values += ","
		}
		i++
		sqltext += "[" + k + "]"
		values += " ? "
		param_data = append(param_data, v)
	}
	sqltext += ") " + values + ")"
	this.LastSqltext = sqltext
	//fmt.Println(i,sqltext)
	//fmt.Println(len(param_data),param_data)
	var result sql.Result
	var err error
	if (this.SqlTx != nil) {
		result, err = this.SqlTx.Exec(sqltext, param_data...)

	} else {
		result, err = this.Masterdb.Exec(sqltext, param_data...)
	}
	return result, err

}

func (this *Mscon) Update(postdata map[string]interface{}) (sql.Result, error) {
	var sqltext string
	sqltext = fmt.Sprintf("update %v set ", this.Tablename)
	i := 0
	param_data := make([]interface{}, 0)
	for k, v := range postdata {
		if (this.Check_data_fields(k) == false) {
			continue
		}
		if (i > 0) {
			sqltext += ","

		}
		i++
		sqltext += "[" + k + "]" + "= ?"
		param_data = append(param_data, v)
	}
	sqlwhere, param := this.Build_where()
	for _, v := range param {
		param_data = append(param_data, v)
	}
	sqltext += sqlwhere
	this.LastSqltext = sqltext
	//fmt.Println(sqltext, param_data)
	var result sql.Result
	var err error
	if (this.SqlTx != nil) {
		result, err = this.SqlTx.Exec(sqltext, param_data...)
	} else {
		result, err = this.Masterdb.Exec(sqltext, param_data...)
	}
	//fmt.Println(err)
	return result, err
}

func (this *Mscon) Delete() (sql.Result, error) {
	sqlwhere, param := this.Build_where()
	sqltext := fmt.Sprintf(" delete from %v %v", this.Tablename, sqlwhere)
	this.LastSqltext = sqltext
	var result sql.Result
	var err error
	if (this.SqlTx != nil) {
		result, err = this.SqlTx.Exec(sqltext, param...)
	} else {
		result, err = this.Masterdb.Exec(sqltext, param...)
	}
	return result, err
}

func (this *Mscon) SetDec(fdname string, quantity int) (sql.Result, error) {
	sqltext := fmt.Sprintf("update %v set %v=%v-%v", this.Tablename, fdname, fdname, quantity)
	sqlwhere, param := this.Build_where()
	sqltext += sqlwhere
	this.LastSqltext = sqltext
	var result sql.Result
	var err error
	if (this.SqlTx == nil) {
		result, err = this.Masterdb.Exec(sqltext, param...)
	} else {
		result, err = this.SqlTx.Exec(sqltext, param...)
	}
	return result, err
}

func (this *Mscon) SetInc(fdname string, quantity int) (sql.Result, error) {
	sqlwhere, param := this.Build_where()
	sqltext := fmt.Sprintf("update %v set %v=%v+%v  %v", this.Tablename, fdname, fdname, quantity, sqlwhere)
	this.LastSqltext = sqltext
	var result sql.Result
	var err error
	if (this.SqlTx == nil) {
		result, err = this.Masterdb.Exec(sqltext, param...)
	} else {
		result, err = this.SqlTx.Exec(sqltext, param...)
	}
	return result, err
}

func (this *Mscon) Query(sqltext string, param []interface{}) []map[string]string {
	this.LastSqltext = sqltext
	var rows *sql.Rows
	var err error
	if (this.SqlTx != nil) {
		rows, err = this.SqlTx.Query(sqltext, param...)
	} else {
		sqldbcon := this.Get_read_dbcon()
		rows, err = sqldbcon.Query(sqltext, param...)
	}
	if (err != nil) {
		fmt.Println(err, this.LastSqltext)
		return nil
	}
	defer rows.Close()
	//????????????
	//??????scanArgs???values???????????????scanArgs??????????????????values??????????????????
	columns, _ := rows.Columns()
	result := make([]map[string]string, 0)
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		//?????????????????????record??????
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
	if (len(result) == 0) {
		return nil
	}
	return result
}

func (this *Mscon) Excute(sqltext string, param []interface{}) (sql.Result, error) {
	this.LastSqltext = sqltext
	var result sql.Result
	var err error
	if (this.SqlTx == nil) {
		result, err = this.Masterdb.Exec(sqltext, param...)
	} else {
		result, err = this.Masterdb.Exec(sqltext, param...)
	}
	return result, err
}

func (this *Mscon) Check_data_fields(fieldname string) (bool) {

	fd_list, ok := G_dbtables[this.Db_name+this.Tablename]
	flag := false
	if (ok) {
		for _, v := range fd_list.([]map[string]string) {
			record := v
			if (record["??????"] == "1" && record["??????"] == "1") {
				continue
			}
			if (record["?????????"] == fieldname) {
				flag = true
				break
			}

		}
		return flag
	} else {
		this.Update_redis(this.Tablename)
		sqltext := this.Get_Fileds_sql(this.Tablename)
		rows, _ := this.Masterdb.Query(sqltext)
		defer rows.Close()
		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		for rows.Next() {
			//?????????????????????record??????
			record := make(map[string]string)
			_ = rows.Scan(scanArgs...)
			for i, col := range values {
				if col != nil {
					record[strings.ToLower(columns[i])] = datatype.Type2str(col)
				}
			}
			if (record["??????"] == "1" && record["??????"] == "1") {
				continue
			}

			if (record["?????????"] == fieldname) {
				flag = true
				break
			}

		}
		return flag
	}
}

func (this *Mscon) Update_redis(tbname string) {
	sqltext := this.Get_Fileds_sql(tbname)
	list := this.Query(sqltext, nil)
	if (list != nil) {
		data_list := make([]map[string]string, 0)
		for _, val := range list {
			col := make(map[string]string)
			for key, _ := range val {
				col[common.Tolow_map_name(key)] = val[key]
			}
			data_list = append(data_list, col)
		}
		G_dbtables[this.Db_name+tbname] = data_list
	}
	//this.Dbinit()
}

func (this *Mscon) Count() int64 {
	sqltext := this.set_sql(1)
	sqlwhere, param := this.Build_where()
	sqltext += sqlwhere
	this.LastSqltext = sqltext
	var rows *sql.Row
	if (this.SqlTx != nil) {
		rows = this.SqlTx.QueryRow(sqltext, param...)
	} else {
		sqldbcon := this.Get_read_dbcon()
		rows = sqldbcon.QueryRow(sqltext, param...)
	}
	var record int64
	rows.Scan(&record)

	return record
}

func (this *Mscon) Sum(fd string) (float64) {
	var result float64
	sqltext := this.set_sql(1)
	sqltext = strings.Replace(sqltext, "count(*)", "sum("+fd+")", -1)
	sqlwhere, param := this.Build_where()
	sqltext += sqlwhere
	this.LastSqltext = sqltext
	var rows *sql.Row
	if (this.SqlTx != nil) {
		rows = this.SqlTx.QueryRow(sqltext, param...)
	} else {
		sqldbcon := this.Get_read_dbcon()
		rows = sqldbcon.QueryRow(sqltext, param...)
	}
	rows.Scan(&result)
	return result
}

func (this *Mscon) Get_Tbname_Sql() (string, []interface{}) {
	sqltext := this.set_sql(0)
	sqlwhere, param := this.Build_where()
	sqltext += sqlwhere
	if this.Sql_order != "" {
		sqltext += " order by " + this.Sql_order
	}
	if this.Sql_limit != "" {
		sqltext = strings.Replace(sqltext, "select", "select top  "+this.Sql_limit, -1)
	}
	return sqltext, param

}

func (this *Mscon) Build_page_sql() (string, []interface{}) {
	pagect := fmt.Sprintf("%d", this.PageSize*this.PageNo)
	last_ct := fmt.Sprintf("%d", this.PageSize*(this.PageNo-1))
	sqlwhere, param := this.Build_where()
	sqltext := "select * from (select *,rows_id from (select  *,row_number() over(order by" + this.Sql_order + ") as rows_id from " + this.Tablename + sqlwhere + " ) tmpdb1 "
	sqltext += " where rows_id<=" + pagect + ") tmpdb2 where rows_id>" + last_ct
	return sqltext, param

}

func (this *Mscon) Select() []map[string]string {
	sqltext := this.set_sql(0)
	sqlwhere, param := this.Build_where()
	sqltext += sqlwhere
	if this.Sql_order != "" {
		sqltext += " order by " + this.Sql_order
	}
	if this.Sql_limit != "" {
		sqltext = strings.Replace(sqltext, "select", "select top  "+this.Sql_limit, -1)
	}
	if (this.PageSize > 0) {
		param=make([]interface{},0)
		sqltext, param = this.Build_page_sql()
	}
	this.LastSqltext = sqltext
	//fmt.Println(sqltext)
	var rows *sql.Rows
	var err error
	if (this.SqlTx != nil) {
		rows, err = this.SqlTx.Query(sqltext, param...)
	} else {
		sqldbcon := this.Get_read_dbcon()
		rows, err = sqldbcon.Query(sqltext, param...)
	}
	if (err != nil) {
		fmt.Println(this.LastSqltext)
		return nil
	}
	if (rows == nil) {
		return nil
	}

	defer rows.Close()
	//????????????
	//??????scanArgs???values???????????????scanArgs??????????????????values??????????????????
	columns, _ := rows.Columns()
	result := make([]map[string]string, 0)
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	j := 0
	for rows.Next() {
		//?????????????????????record??????
		record := make(map[string]string)
		_ = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = datatype.Type2str(col)
				//record[columns[i]] = col.([]byte)
			} else {
				record[columns[i]] = ""
			}
		}
		result = append(result, record)
		//result[j] = record
		j++

	}
	if (len(result) == 0) {
		return nil
	}
	return result
}

func (this *Mscon) GetLastSql() string {
	return this.LastSqltext
}

func (this *Mscon) Get_Fileds_sql(tbname string) (string) {
	sqltext := fmt.Sprintf(`SELECT
a.colorder,a.name ?????????,
(case when COLUMNPROPERTY( a.id,a.name,'IsIdentity')=1 then 1 else 0 end) ??????,
(case when (SELECT count(*) FROM sysobjects
WHERE (name in (SELECT name FROM sysindexes
WHERE (id = a.id) AND (indid in
(SELECT indid FROM sysindexkeys
WHERE (id = a.id) AND (colid in
(SELECT colid FROM syscolumns WHERE (id = a.id) AND (name = a.name)))))))
AND (xtype = 'PK'))>0 then 1 else 0 end) ??????,b.name ??????,a.length ???????????????,
COLUMNPROPERTY(a.id,a.name,'PRECISION') as ??????,
isnull(COLUMNPROPERTY(a.id,a.name,'Scale'),0) as ????????????,(case when a.isnullable=1 then 1 else 0 end) ?????????,
isnull(e.text,'') ?????????,isnull(g.[value], ' ') AS [??????]
FROM  syscolumns a
left join systypes b on a.xtype=b.xusertype
inner join sysobjects d on a.id=d.id and d.xtype='U' and d.name<>'dtproperties'
left join syscomments e on a.cdefault=e.id
left join sys.extended_properties g on a.id=g.major_id AND a.colid=g.minor_id
left join sys.extended_properties f on d.id=f.class and f.minor_id=0
--where b.name is not null
WHERE d.name='%v' --????????????????????????,???????????????
order by a.id,a.colorder`, tbname)
	return sqltext
}
func (this *Mscon) Get_new_add() map[string]string {
	fd_list, ok := G_dbtables[this.Db_name+this.Tablename]
	if (ok) {
		//fmt.Println(fd_list)
		result := make(map[string]string)
		for _, v := range fd_list.([]map[string]string) {
			fd_name := v["field"]
			result[fd_name] = ""
		}
		return result
	} else {
		this.Update_redis(this.Tablename)
		sqltext := this.Get_Fileds_sql(this.Tablename)
		rows, _ := this.Masterdb.Query(sqltext)
		defer rows.Close()
		columns, _ := rows.Columns()
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		result := make(map[string]string)
		for i := range values {
			scanArgs[i] = &values[i]
		}
		for rows.Next() {
			//?????????????????????record??????
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

func (this *Mscon) Rollback() {
	//fmt.Println(this.SqlTx, "rollback")
	if (this.SqlTx == nil) {
		return
	}
	this.SqlTx.Rollback()
	this.SqlTx = nil
}

func (this *Mscon) Commit() {
	//fmt.Println("commit", this.SqlTx)
	if (this.SqlTx == nil) {
		return
	}
	this.SqlTx.Commit()
	this.SqlTx = nil
}
