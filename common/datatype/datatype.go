package datatype

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 结构转map
func StructJson2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("json")
		if tag == "" {
			tag = Tolow_map_name(t.Field(i).Name)
		}
		fd_name := tag
		data[fd_name] = v.Field(i).Interface()
	}
	return data
}

func Struct2ListMap(obj interface{}, perfix string) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		fd_name := perfix + Tolow_map_name(t.Field(i).Name)
		data[fd_name] = v.Field(i).Interface()
	}
	return data
}

// 结构转map
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
	//fmt.Println(t.Kind(),v.Kind())
	//fmt.Println(t.NumField())
	//fmt.Println(obj)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		is_skip := t.Field(i).Tag.Get("skip")
		if is_skip == "1" {
			continue
		}
		fd_name := t.Field(i).Tag.Get("json")
		if fd_name == "" {
			fd_name = t.Field(i).Name
		}
		data[fd_name] = v.Field(i).Interface()
	}
	return data
}

/*
func Map2Struct(obj interface{}, data map[string]interface{}, tag string) (interface{}) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if(tag==""){
		tag="json"
	}
	for key,val:=range data{
		ind:=getstructindex(obj,key,tag)
		if(ind!=-1){
			v.Field(ind)
		}
		fd_name := t.Field(i).Tag.Get("json")
		if (fd_name == "") {
			fd_name = t.Field(i).Name
		}
		data[fd_name] = v.Field(i).Interface()
	}
	return obj
}

func getstructindex(obj interface{}, name, tag string) (int) {
	t := reflect.TypeOf(obj)
	result := -1
	for i := 0; i < t.NumField(); i++ {
		is_skip := t.Field(i).Tag.Get("skip")
		if (is_skip == "1") {
			continue
		}
		fd_name := t.Field(i).Tag.Get("json")
		if (fd_name != "") {
			result = i
			break
		}

	}
	return result
}*/

// 结构转map
func Type2Map(obj interface{}) map[string]interface{} {
	result := obj.(map[string]interface{})
	/*var data=make(map[string]interface{})
	for k,v:=range result{
		fd_name := Tolow_map_name(k)
		data[fd_name] = v
	}*/
	return result
}

// 驼峰写法转下划线写法
func Tolow_map_name(name string) string {
	result := ""
	for k, v := range name {
		if v >= 'A' && v <= 'Z' {
			if k > 0 {
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

// map转结构体
func DataToDBStruct(data map[string]string, out interface{}) {
	ss := reflect.ValueOf(out).Elem()
	for i := 0; i < ss.NumField(); i++ {
		name := ss.Type().Field(i).Tag.Get("json")
		if name == "" {
			name = ss.Type().Field(i).Name
		}
		val, ok := data[name]
		if ok == false {
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

// map转结构体
func DataToStruct(data map[string]string, out interface{}) {
	ss := reflect.ValueOf(out).Elem()
	for i := 0; i < ss.NumField(); i++ {
		name := ss.Type().Field(i).Name
		val, ok := data[Tolow_map_name(name)]
		if ok == false {
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

func MapArray2interface(data []map[string]string) []map[string]interface{} {
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

func Type2int(val interface{}) int {
	val_str := Type2str(val)
	result, err := strconv.Atoi(val_str)
	if err != nil {
		return -1
	} else {
		return result
	}
}

func Type2str(val interface{}) string {
	//fmt.Println(fmt.Sprintf("%T,%v",val,val))
	if val == nil {
		return ""
	}
	var result string = ""
	switch val.(type) {
	case []string:
		strArray := val.([]string)
		result = strings.Join(strArray, "")
	case []uint8:
		tmp := val.([]uint8)
		if len(tmp) == 1 {
			for _, v := range tmp {
				result = fmt.Sprintf("%v", v)
			}
			if Str2Int(result) > 31 {
				result = string(val.([]uint8))
			}
		} else {
			result = string(val.([]uint8))
		}

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

func DateStringToGoTime(tm string) time.Time {
	t, err := time.ParseInLocation("2006-01-02", tm, time.Local)
	if nil == err && !t.IsZero() {
		return t
	}
	return time.Time{}
}

func FormatDate(key string) string {
	return Int2Time_str(int64(Str2Int(key)))
}
func Replace_map(data map[string]string, memo string) string {
	msg := memo
	for k, v := range data {
		//fmt.Println(k,v)
		msg = strings.Replace(msg, "{"+k+"}", v, -1)
		//fmt.Println(msg)
	}
	return msg
}

func Str2Int(key string) int {
	value, err := strconv.Atoi(key)
	if err != nil {
		return -1
	}
	return value
}

func Str2Int64(key string) int64 {
	result, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return -1
	}
	return result
}

func Str2Float(key string) float64 {
	value, _ := strconv.ParseFloat(key, 64)
	return value
}

func Map2str(postdata map[string]interface{}) map[string]string {
	data := make(map[string]string)
	for key, val := range postdata {
		data[key] = Type2str(val)
	}
	return data
}

func Get_UUID() string {
	uuid := uuid.NewV4()
	return uuid.String()

}

// 下划线转驼峰写法转写法
func ToUP_map_name(name string) string {
	result := ""
	flag := false
	for _, v := range name {
		if v == '_' {
			flag = true
		} else {
			if flag {
				result += strings.ToUpper(string(v))
				flag = false
			} else {
				result += strings.ToLower(string(v))
			}
		}

	}
	return strings.Title(result)
}

func String2Time(date string) int64 {
	toBeCharge := date
	timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	sr := theTime.Unix()                                            //转化为时间戳 类型是int64
	return sr
}

func String2date(date string) int64 {
	toBeCharge := date
	timeLayout := "2006-01-02"                                      //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	sr := theTime.Unix()                                            //转化为时间戳 类型是int64
	return sr
}

func String2TimeStamp(date string) time.Time {
	toBeCharge := date
	timeLayout := "2006-01-02"                                      //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	return theTime
}

func GetCalMonth(date string, mon_sub int) string {
	toBeCharge := date
	timeLayout := "2006-01-02"                                      //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	tm := theTime.AddDate(0, mon_sub, 0)                            //转化为时间戳 类型是int64
	return tm.String()[:10]
}

func Int2Time_str(date int64) string {
	//格式化为字符串,tm为Time类型
	tm := time.Unix(date, 0)
	return tm.Format("2006-01-02 15:04:05")

}

func Int2Date_str(date int64) string {
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
	if err1 != nil {
		return -1
	}
	return data
}

func Type2List(data interface{}) []map[string]interface{} {
	b, ok := data.([]map[string]interface{})
	if ok {
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

func String2Json(strin string) map[string]interface{} {
	if strin == "" {
		return nil
	}
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(strin), &data)
	if err != nil {
		fmt.Println("string2json", strin)
		fmt.Println("string2json", err)
		return nil
	} else {
		return data
	}
}

func InterfaceToStringList(data []map[string]interface{}) []map[string]string {
	list := make([]map[string]string, 0)
	for _, val := range data {
		temp := make(map[string]string)
		for key, d_val := range val {
			temp[key] = Type2str(d_val)
		}
		list = append(list, temp)
	}
	return list
}

func String2JsonList(strin string) []map[string]interface{} {
	if strin == "" {
		return nil
	}
	data := make([]map[string]interface{}, 0)
	err := json.Unmarshal([]byte(strin), &data)
	if err != nil {
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

func JSON2maplist(arr interface{}) []map[string]interface{} {
	list, ok := arr.([]interface{})
	if !ok {
		return nil
	}
	data := make([]map[string]interface{}, 0)
	flag := false
	for _, v := range list {
		temp, ok := v.(map[string]interface{})
		if !ok {
			flag = true
			break
		}
		data = append(data, temp)
	}
	if flag {
		return nil
	} else {
		return data
	}
}

func JSON2map(arr interface{}) map[string]interface{} {
	list, ok := arr.(map[string]interface{})
	if !ok {
		return nil
	}
	return list

}

func Map2Json(json_data map[string]interface{}) []byte {
	list, err := json.Marshal(&json_data)
	if err != nil {
		fmt.Println(json_data)
		return nil
	} else {
		return list
	}
}

func MapList2Json(json_data []map[string]interface{}) []byte {
	list, err := json.Marshal(&json_data)
	if err != nil {
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

func Array_to_map(data interface{}) []map[string]interface{} {
	list, ok := data.([]interface{})
	if !ok {
		return nil
	}
	result := make([]map[string]interface{}, 0)
	for _, val := range list {
		tmp, ok := val.(map[string]interface{})
		if ok {
			result = append(result, tmp)
		}

	}
	return result

}

func SetAutoBh(strlen int, str_qz string, str_val string) string {
	j := strlen - len(str_qz) - len(str_val)
	result := str_qz
	for i := 0; i < j; i++ {
		result += "0"
	}
	result += str_val
	return result
}

// 10进制字符串转16进制字符串
func DecStringTOHexString(str1 string) string {
	num := Str2Int(str1)
	result := fmt.Sprintf("%x", num)
	if len(result)%2 != 0 {
		result = "0" + result
		//return nil
	}
	if num > 255 {

	}
	return result
}

// 字符串转byte数组
func HexStringToByte(str string) []byte {
	slen := len(str)
	data := str
	if slen%2 != 0 {
		data = "0" + str
		//return nil
	}
	slen = len(data)
	//fmt.Println(data)
	bHex := make([]byte, len(data)/2)
	ii := 0
	for i := 0; i < len(data); i = i + 2 {
		if slen != 1 {
			ss := string(data[i]) + string(data[i+1])
			//fmt.Println("str",ss)
			bt, _ := strconv.ParseUint(ss, 16, 32)
			bHex[ii] = byte(bt)
			//fmt.Println("bhex",bHex[ii])
			ii = ii + 1
			slen = slen - 2
		}
	}
	return bHex
}

func BinStringToByte(str string) int32 {
	bt, _ := strconv.ParseUint(str, 2, 32)
	x := int32(bt)
	return x
}

// 16进制字节转字符串
func HexByteToString(list []byte) string {
	result := ""
	for i := 0; i < len(list); i++ {
		tmp := fmt.Sprintf("%x", list[i])
		if len(tmp) == 1 {
			tmp = "0" + tmp
		}
		result += tmp
	}
	return strings.ToUpper(result)
}

// 16进制字节转字符串
func StringToHexString(data string) string {
	list := Int32ToBytes(Str2Int(data), true)
	result := ""
	for i := 0; i < len(list); i++ {
		tmp := fmt.Sprintf("%x", list[i])
		if len(tmp) == 1 {
			tmp = "0" + tmp
		}
		result += tmp

	}
	return strings.ToUpper(result)
}

// 整形转换成字节
func Int32ToBytes(n int, big_flag bool) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	if big_flag {
		binary.Write(bytesBuffer, binary.BigEndian, x)
	} else {
		binary.Write(bytesBuffer, binary.LittleEndian, x)
	}
	return bytesBuffer.Bytes()
}

// 整形转换成字节
func Int64ToBytes(n int64, big_flag bool) []byte {
	x := n
	bytesBuffer := bytes.NewBuffer([]byte{})
	if big_flag {
		binary.Write(bytesBuffer, binary.BigEndian, x)
	} else {
		binary.Write(bytesBuffer, binary.LittleEndian, x)
	}
	return bytesBuffer.Bytes()
}

// 整形转换成字节
func Int8ToBytes(n int, big_flag bool) []byte {
	x := int8(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	if big_flag {
		binary.Write(bytesBuffer, binary.BigEndian, x)
	} else {
		binary.Write(bytesBuffer, binary.LittleEndian, x)
	}
	return bytesBuffer.Bytes()
}

func UInt8ToBytes(n uint8) []byte {
	x := n
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.LittleEndian, x)
	return bytesBuffer.Bytes()
}

func BytesToUInt8(buf []byte) uint8 {
	b_buf := bytes.NewBuffer(buf)
	var x uint8
	binary.Read(b_buf, binary.BigEndian, &x)
	return x
}

// 字节转换成整形
func BytesToInt(buf []byte, big_flag bool) int {
	b_buf := bytes.NewBuffer(buf)
	var x int32
	if big_flag {
		binary.Read(b_buf, binary.BigEndian, &x)
	} else {
		binary.Read(b_buf, binary.LittleEndian, &x)
	}

	return int(x)
}

// 字节转换成整形
func BytesToInt16(buf []byte, big_flag bool) int16 {
	b_buf := bytes.NewBuffer(buf)
	var x int16
	if big_flag {
		binary.Read(b_buf, binary.BigEndian, &x)
	} else {
		binary.Read(b_buf, binary.LittleEndian, &x)
	}
	return x
}

func BytesToInt64(buf []byte) uint64 {
	if len(buf) < 8 {
		return 0
	}
	return binary.LittleEndian.Uint64(buf)
}

// 整形转换成字节
func Int16ToBytes(n int, big_flag bool) []byte {
	x := int16(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	if big_flag {
		binary.Write(bytesBuffer, binary.BigEndian, x)
	} else {
		binary.Write(bytesBuffer, binary.LittleEndian, x)
	}
	return bytesBuffer.Bytes()
}

// 16进制字节转10进制字符串
func HexByteToDecString(list []byte) string {
	result := ""
	for i := 0; i < len(list); i++ {
		result += fmt.Sprintf("%d", list[i])
	}
	return strings.ToUpper(result)
}

// 16进制字节转8位二进制字符串
func HexByteToBinString(list []byte) string {
	result := ""
	for i := 0; i < len(list); i++ {
		tmp := fmt.Sprintf("%b", list[i])
		binlen := len(tmp)
		if binlen < 8 {
			for i := binlen; i < 8; i++ {
				tmp = "0" + tmp
			}
		}
		result += tmp
	}
	return strings.ToUpper(result)
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// 检测是否有注入代码，过滤Sql语句的注入代码
func FilteredSQLInject(to_match_str string) bool {
	//过滤 ‘
	//ORACLE 注解 --  /**/
	//关键字过滤 update ,delete
	// 正则的字符串, 不能用 " " 因为" "里面的内容会转义
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	re, err := regexp.Compile(str)
	if err != nil {
		//panic(err.Error())
		return false
	}
	return re.MatchString(to_match_str)
}

func GetRedisValueMap(src interface{}) map[string]interface{} {
	data, ok := src.(map[string]interface{})
	if ok {
		return data
	} else {
		return nil
	}
}

func FormatFloat(val float64, ws int) string {
	v_format := "%." + Type2str(ws) + "f"
	return fmt.Sprintf(v_format, val)
}

func SafeConvertMap(m interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	switch v := m.(type) {
	case map[string]interface{}:
		// 如果值是map类型，则递归调用该函数进行转换
		for key, val := range v {
			result[key] = val
		}
	case map[string]string:
		for key, val := range v {
			result[key] = val
		}
	default:
		// 如果值不是map类型，则直接赋值给结果map
		result = nil
	}

	return result
}
