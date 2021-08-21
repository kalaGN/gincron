package main

import (
	"encoding/json"
	"fmt"
	"os"

	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/kalaGN/gincron/src/common"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//定义结构
	type Job struct {
		Id              int
		src_config      string
		operate_content string
		task_name       string
		cycle_time      string
	}

	type Cycletime struct {
		Cron string `cron`
		Time string `time`
	}

	type Srcconfig struct {
		Url string `url`
	}

	type Opcontent struct {
		Type  string      `type`
		Path  string      `path`
		Param interface{} `json:"params"`
	}

	var (
		job  Job
		jobs []Job
	)

	//日志
	file := "./" + time.Now().Format("20060102") + ".txt"

	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if nil != err {
		panic(err)
	}
	loger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	// 获取配置

	for {
		db, err := common.Getdb("commondb")
		if err != nil {

			loger.Fatal("open commondb fail:", err)
			log.Fatalln("open commondb fail:", err)
		}
		rows, err := db.Query("select id, src_config, operate_content, task_name, cycle_time from com_schedule_task_define where is_run=0;")
		if err != nil {
			log.Fatalln("Query select fail:", err)
			loger.Fatal("Query select fail:", err)

		}

		for rows.Next() {
			err = rows.Scan(&job.Id, &job.src_config, &job.operate_content, &job.task_name, &job.cycle_time)
			jobs = append(jobs, job)
			if err != nil {
				fmt.Print(err.Error())
			}
		}
		defer rows.Close()
		//gin.DisableConsoleColor()

		for _, v := range jobs {
			loger.Printf("jobs length: %d", len(jobs))
			temp := v.cycle_time
			obj := Cycletime{}
			srcconfig := Srcconfig{}
			opcontent := Opcontent{}

			// 程序执行时间
			err := json.Unmarshal([]byte(temp), &obj)
			if err != nil {
				loger.Printf("json error: %s", err.Error())
			}
			// url
			err2 := json.Unmarshal([]byte(v.src_config), &srcconfig)
			if err2 != nil {
				loger.Printf("json error2: %s", err2.Error())
			}
			err3 := json.Unmarshal([]byte(v.operate_content), &opcontent)
			if err3 != nil {
				loger.Printf("json error3: %s", err3.Error())
			}
			loger.Printf("operate_content:%s", v.operate_content)

			// 获取请求参数
			param := jsoniter.Get([]byte(v.operate_content), "param")
			m := JSONToMap(param.ToString())
			var pmStr = "/?"
			for key, value := range m {
				if reflect.TypeOf(value).String() == "float64" {
					value = Strval(value)
				}
				pmStr = pmStr + fmt.Sprintf("%s=%s&", key, value)

			}
			loger.Printf("url:%s", srcconfig.Url+opcontent.Path+pmStr)

			// 获取当前时间
			now := time.Now().Format("2006-01-02 15:04:05")

			// 比较时间
			t1, err1 := time.Parse("2006-01-02 15:04:05", obj.Time)
			t2, err2 := time.Parse("2006-01-02 15:04:05", now)
			if err1 == nil && err2 == nil && t1.Before(t2) {
				loger.Printf("run task id:%d", v.Id)
				//执行请求任务
				res := requestGet(srcconfig.Url + opcontent.Path + pmStr)
				loger.Printf("run result :%s", res)
				//修改执行结果
				if err != nil {
					loger.Printf("request error:%s", err.Error())
				}
				updateSql := fmt.Sprintf("update com_schedule_task_define set is_run=1 where id='%d';", v.Id)
				loger.Printf("sql:%s", updateSql)

				result, err := db.Exec(updateSql)
				if err != nil {
					loger.Printf("exec failed:%s sql:%s", err, updateSql)
				}
				idAff, err := result.RowsAffected()
				if err != nil {
					loger.Printf("RowsAffected failed:%s", err)
				}
				loger.Printf("RowsAffected:%d", idAff)

			} else {
				//fmt.Println("time not ")
			}

		}
		// 清空jobs
		jobs = jobs[0:0]
		// 休眠5s
		time.Sleep(time.Duration(5) * time.Second)
		//fmt.Printf("\n\r")
		db.Close()
	}
}

//**
//
func requestGet(url string) string {
	timeout := time.Duration(2 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	res, err := client.Get(url)

	if err != nil {
		//panic(err)
		return err.Error()
	} else {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		return string(body[:])
	}

}

// json转map函数，通用
func JSONToMap(str string) map[string]interface{} {

	var tempMap map[string]interface{}

	err := json.Unmarshal([]byte(str), &tempMap)

	if err != nil {
		panic(err)
	}

	return tempMap
}

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
