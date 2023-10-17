package main

/*type Nav struct {
	Id          int    `xorm:"not null pk autoincr INT(11)"`
	ParentId    int    `xorm:"default 0 INT(11)"`
	NavName     string `xorm:"VARCHAR(200)"`
	NavCode     string `xorm:"VARCHAR(200)"`
	NavModule   string `xorm:"VARCHAR(200)"`
	NavImage    string `xorm:"VARCHAR(300)"`
	IsDisplay   int    `xorm:"default 0 INT(11)"`
	OrderNumber int    `xorm:"INT(11)"`
	IsTel       int    `xorm:"default 0 INT(11)"`
}

type Nav_struct struct {
	Id       int    `json:"id"`
	ParentId int    `json:"parent_id"`
	NavName  string `json:"nav_name"`
}*/

func main() {
	/*str1 := "guyi_bhkdjfkllsls;slsldk"
	pre := "guyi_"
	strlen := len(pre)
	if str1[:strlen] == pre {
		fmt.Println(pre)
	}
	fmt.Println(common.Get_UUID())*/

}

/*func oracle_test() {
	fmt.Println("Oracle test start....")
	db, err := sql.Open("oci8", "bfcrm8/DHHZDHHZ@10.100.2.202:1521/crmtest")
	if err != nil {
		fmt.Println(err)
	}
	rows, err := db.Query("select * from hykdef")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	if (rows == nil) {
		fmt.Println("stop test....")
		return
	}

	defer rows.Close()
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
	fmt.Println(result)

}

func redis_test(key string) {
	redis := rediscomm.NewRedisComm()
	fmt.Println(redis.SetKey(key).HasKey())
}

func test() {
	nav := new(Nav)
	db := gydblib.Get_New_Main_DB()
	rows := db.Tbname("nav").Find()
	common.DataToStruct(rows, nav)
	fmt.Println(nav)
	fmt.Println(common.Struct2Map(*nav))
	//时间转换测试
	timestamp := time.Now().Unix()
	fmt.Println(common.Int2Time_str(timestamp))
	fmt.Println(common.String2Time("2006-01-02 15:04:05"))
	fmt.Println(path.Ext("guyi.txt"))
}

func sum(values []int, resultChan chan int) {
	sum := 0
	fmt.Println(values)
	for _, value := range values {
		sum += value
	}
	// 将计算结果发送到channel中
	resultChan <- sum

}

func Subs() { //订阅者
	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("connect redis error :", err)
		return
	}
	defer conn.Close()
	psc := redis.PubSubConn{conn}
	psc.Subscribe("channel1") //订阅channel1频道
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			fmt.Printf("%s: message: %s\n", v.Channel, v.Data)
		case redis.Subscription:
			fmt.Printf("Subscr:%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			fmt.Println(v)
			return
		}
	}
}
*/
