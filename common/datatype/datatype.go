package datatype

import (
	"time"
	"reflect"
	"strings"
	"strconv"
	"fmt"
	"github.com/satori/go.uuid"
	"math/rand"
	"encoding/json"
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding/simplifiedchinese"
)

//结构转map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		fd_name := Tolow_map_name(t.Field(i).Name)
		data[fd_name] = v.Field(i).Interface()
	}
	return data
}

func Struct2DBMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	//fmt.Println(t.NumField())
	//fmt.Println(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		is_skip := t.Field(i).Tag.Get("skip")
		if (is_skip == "1") {
			continue
		}
		fd_name := t.Field(i).Tag.Get("json")
		if (fd_name == "") {
			fd_name = t.Field(i).Name
		}
		data[fd_name] = v.Field(i).Interface()
	}
	return data
}

//结构转map
func Type2Map(obj interface{}) map[string]interface{} {
	result := obj.(map[string]interface{})
	/*var data=make(map[string]interface{})
	for k,v:=range result{
		fd_name := Tolow_map_name(k)
		data[fd_name] = v
	}*/
	return result
}

//驼峰写法转下划线写法
func Tolow_map_name(name string) (string) {
	result := ""
	for k, v := range name {
		if (v >= 'A' && v <= 'Z') {
			if (k > 0) {
				result += "_" + strings.ToLower(string(v))
			} else {
				result += strings.ToLower(string(v))
			}
		} else {
			result += strings.ToLower(string(v))
		}

	}
	return result
}

//map转结构体
func DataToDBStruct(data map[string]string, out interface{}) {
	ss := reflect.ValueOf(out).Elem()
	for i := 0; i < ss.NumField(); i++ {
		name := ss.Type().Field(i).Tag.Get("json")
		if (name == "") {
			name = ss.Type().Field(i).Name
		}
		val, ok := data[name]
		if (ok == false) {
			continue
		}
		switch ss.Field(i).Kind() {
		case reflect.String:
			ss.FieldByName(name).SetString(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			ss.FieldByName(name).SetInt(int64(i))
		case reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			ss.FieldByName(name).SetUint(uint64(i))
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				continue
			}
			ss.FieldByName(name).SetFloat(f)
		case reflect.Struct:
			fmt.Println(ss.Field(i), ss.Field(i).NumField())
			//f,err:=time.Parse("2006-01-02 15:04:05", val)
			//ss.FieldByName(name).Set(f)
		default:
			fmt.Println("unknown type:%+v", ss.Field(i).Kind())
		}
	}
	return
}

//map转结构体
func DataToStruct(data map[string]string, out interface{}) {
	ss := reflect.ValueOf(out).Elem()
	for i := 0; i < ss.NumField(); i++ {
		name := ss.Type().Field(i).Name
		val, ok := data[Tolow_map_name(name)]
		if (ok == false) {
			continue
		}
		switch ss.Field(i).Kind() {
		case reflect.String:
			ss.FieldByName(name).SetString(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			ss.FieldByName(name).SetInt(int64(i))
		case reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i, err := strconv.Atoi(val)
			if err != nil {
				continue
			}
			ss.FieldByName(name).SetUint(uint64(i))
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				continue
			}
			ss.FieldByName(name).SetFloat(f)
		case reflect.Struct:
			fmt.Println(ss.Field(i), ss.Field(i).NumField())
			//f,err:=time.Parse("2006-01-02 15:04:05", val)
			//ss.FieldByName(name).Set(f)
		default:
			fmt.Println("unknown type:%+v", ss.Field(i).Kind())
		}
	}
	return
}

func Has_map_index(name string, data map[string]interface{}) bool {
	_, ok := data[name]
	return ok

}

func MapArray2interface(data []map[string]string) ([]map[string]interface{}) {
	list := make([]map[string]interface{}, 0)
	for _, v := range data {
		tmp := make(map[string]interface{})
		for key, val := range v {
			tmp[key] = val
		}
		list = append(list, tmp)
	}
	return list
}

func Type2int(val interface{}) (int) {
	val_str := Type2str(val)
	result, err := strconv.Atoi(val_str)
	if (err != nil) {
		return -1
	} else {
		return result
	}
}

func Type2str(val interface{}) (string) {
	//fmt.Println(fmt.Sprintf("%T,%v",val,val))
	if (val == nil) {
		return ""
	}
	var result string = ""
	switch val.(type) {
	case []string:
		strArray := val.([]string)
		result = strings.Join(strArray, "")
	case []uint8:
		result = string(val.([]uint8))
		//case []byte:
		//	result = string(val.([]byte))
	default:
		result = fmt.Sprintf("%v", val)
	}
	return result
}

func Byte2str(postdata []map[string][]byte) []map[string]interface{} {
	data := make([]map[string]interface{}, 0)
	for _, val := range postdata {
		temp := make(map[string]interface{})
		for k, v := range val {
			temp[strings.ToLower(k)] = string(v[:])
		}
		data = append(data, temp)
	}
	return data
}

func FormatDate(key string) (string) {
	return Int2Time_str(int64(Str2Int(key)))
}
func Replace_map(data map[string]string, memo string) string {
	msg := memo
	for k, v := range data {
		msg = strings.Replace(msg, "{"+k+"}", v, -1)
	}
	return msg
}

func Str2Int(key string) int {
	value, err := strconv.Atoi(key)
	if (err != nil) {
		return -1
	}
	return value
}

func Str2Float(key string) (float64) {
	value, _ := strconv.ParseFloat(key, 64)
	return value
}

func Map2str(postdata map[string]interface{}) (map[string]string) {
	data := make(map[string]string)
	for key, val := range postdata {
		data[key] = Type2str(val)
	}
	return data
}

func Get_UUID() (string) {
	uuid, _ := uuid.NewV4()
	return uuid.String()

}

//下划线转驼峰写法转写法
func ToUP_map_name(name string) (string) {
	result := ""
	flag := false
	for _, v := range name {
		if (v == '_') {
			flag = true
		} else {
			if (flag) {
				result += strings.ToUpper(string(v))
				flag = false
			} else {
				result += strings.ToLower(string(v))
			}
		}

	}
	return strings.Title(result)
}

func String2Time(date string) (int64) {
	toBeCharge := date
	timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	sr := theTime.Unix()                                            //转化为时间戳 类型是int64
	return sr
}

func String2date(date string) (int64) {
	toBeCharge := date
	timeLayout := "2006-01-02"                                      //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	sr := theTime.Unix()                                            //转化为时间戳 类型是int64
	return sr
}

func Int2Time_str(date int64) (string) {
	//格式化为字符串,tm为Time类型
	tm := time.Unix(date, 0)
	return tm.Format("2006-01-02 15:04:05")

}

func Int2Date_str(date int64) (string) {
	//格式化为字符串,tm为Time类型
	tm := time.Unix(date, 0)
	return tm.Format("2006-01-02")

}

func MapString2interface(data map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for key, val := range data {
		result[key] = val
	}
	return result
}

func Int64toint(val int64) int {
	result := strconv.FormatInt(val, 10)
	data, err1 := strconv.Atoi(result)
	if (err1 != nil) {
		return -1
	}
	return data
}

func Type2List(data interface{}) []map[string]interface{} {
	b, ok := data.([]map[string]interface{})
	if (ok) {
		return b
	} else {
		return nil
	}
}

func Type2map(data interface{}) map[string]interface{} {
	b, ok := data.(map[string]interface{})
	if ok {
		return b
	} else {
		return nil
	}

}
func GenerateRangeNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max - min)
	randNum = randNum + min
	//fmt.Printf("rand is %v\n", randNum)
	return randNum
}

func String2Json(strin string) (map[string]interface{}) {
	if (strin == "") {
		return nil
	}
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(strin), &data)
	if (err != nil) {
		return nil
	} else {
		return data
	}
}

func String2JsonList(strin string) ([]map[string]interface{}) {
	if (strin == "") {
		return nil
	}
	data := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(strin), &data)
	if (err != nil) {
		return nil
	} else {
		return data
	}
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s)) //使用zhifeiya名字做散列值，设定后不要变
	return hex.EncodeToString(h.Sum(nil))
}

func JSON2maplist(arr interface{}) ([]map[string]interface{}) {
	list, ok := arr.([]interface{})
	if (!ok) {
		return nil
	}
	data := make([]map[string]interface{}, 0)
	flag := false
	for _, v := range list {
		temp, ok := v.(map[string]interface{})
		if (!ok) {
			flag = true
			break;
		}
		data = append(data, temp)
	}
	if (flag) {
		return nil
	} else {
		return data
	}
}

func JSON2map(arr interface{}) (map[string]interface{}) {
	list, ok := arr.(map[string]interface{})
	if (!ok) {
		return nil
	}
	return list

}

func Map2Json(json_data map[string]interface{}) ([]byte) {
	list, err := json.Marshal(&json_data)
	if (err != nil) {
		fmt.Println(json_data)
		return nil
	} else {
		return list
	}
}

func Utf8ToGBK(utf8str string) string {
	result, _, _ := transform.String(simplifiedchinese.GBK.NewEncoder(), utf8str)
	return result
}

func Array_to_map(data interface{}) ([]map[string]interface{}) {
	list, ok := data.([]interface{})
	if (!ok) {
		return nil
	}
	result := make([]map[string]interface{}, 0)
	for _, val := range list {
		tmp, ok := val.(map[string]interface{})
		if (ok) {
			result = append(result, tmp)
		}

	}
	return result

}

func SetAutoBh(strlen int, str_qz string, str_val string) (string) {
	j := strlen - len(str_qz) - len(str_val)
	result := str_qz
	for i := 0; i < j; i++ {
		result += "0"
	}
	result += str_val
	return result
}
