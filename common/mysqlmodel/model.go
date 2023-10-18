package mysqlmodel

import (
	"fmt"
	"github.com/guyigood/gylib/gydblib"
	"strings"
)

type MysqlModel struct {
	TableName string
}

func NewMysql_model() *MysqlModel {
	this := new(MysqlModel)
	return this
}

func (this *MysqlModel) Build() string {
	db := gydblib.Get_New_Main_DB()
	db.Query("use information_schema", nil)
	sqltext := fmt.Sprintf(`SELECT
        COLUMN_NAME,DATA_TYPE, COLUMN_COMMENT,
        COLUMN_DEFAULT,COLUMN_KEY,EXTRA
        FROM COLUMNS
        WHERE TABLE_NAME = '%v%v'  and TABLE_SCHEMA = '%v'`, db.Db_perfix, this.TableName, db.Db_name)
	list := db.Query(sqltext, nil)
	//fmt.Println(sqltext,list)
	fdlist := ""
	//fmt.Println(list,sqltext)
	model := fmt.Sprintf("type %v struct{\n", strings.Title(this.TableName))
	for _, v := range list {

		fdlist = strings.Title(v["COLUMN_NAME"]) + " " + this.TypeConvert(v["DATA_TYPE"]) + " `json:\"" + v["COLUMN_NAME"] + "\"`"
		model += fdlist + "\n"
		//fmt.Println(v)
		//fmt.Println(strings.Title(v["COLUMN_NAME"])+" "+this.TypeConvert(v["DATA_TYPE"]))
	}
	model += "}"
	db.Query("use "+db.Db_name, nil)
	return model

}

func (this *MysqlModel) TypeConvert(str string) string {

	switch str {
	case "smallint", "tinyint":
		return "int8"

	case "varchar", "text", "longtext", "char", "mediumtext":
		return "string"

	case "date":
		return "string"

	case "int":
		return "int"

	case "timestamp", "datetime":
		return "time.Time"

	case "bigint", "mediumint":
		return "int64"

	case "float", "double", "decimal":
		return "float64"

	default:
		return str
	}
}
