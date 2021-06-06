package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/garyburd/redigo/redis"
	"github.com/go-ini/ini"
	_ "github.com/go-sql-driver/mysql"
)

/**
 * 数据库配置结构
 */
type config struct {
	host     string
	port     string
	password string
	db       string
	timeout  string
	insid    string
}

/**
 * 获取redis配置信息
 */
func GetRedisconfig(inifile string, redisname string) config {
	//获取配置
	cfg, err := ini.Load(inifile)
	if err != nil {
		fmt.Printf("fail to read file: "+inifile, err)
	}
	return config{
		host:     cfg.Section("production").Key(redisname + ".host").String(),
		port:     cfg.Section("production").Key(redisname + ".port").String(),
		password: cfg.Section("production").Key(redisname + ".password").String(),
		db:       cfg.Section("production").Key(redisname + ".db").String(),
		timeout:  cfg.Section("production").Key(redisname + ".timeout").String(),
		insid:    cfg.Section("production").Key(redisname + "system.instanceid").String(),
	}

}

/**
 * mysql数据库配置
 */
type dbconfig struct {
	host     string
	port     string
	username string
	password string
	db       string
	timeout  string
	charset  string
}

/**
 * 获取mysql配置
 */
func Getdbconfig(inifile string, dbsection string) dbconfig {
	//获取配置
	cfg, err := ini.Load(inifile)
	if err != nil {
		fmt.Printf("fail to read file: "+inifile, err)
	}
	return dbconfig{
		host:     cfg.Section("production").Key(dbsection + ".host").String(),
		port:     cfg.Section("production").Key(dbsection + ".port").String(),
		password: cfg.Section("production").Key(dbsection + ".password").String(),
		username: cfg.Section("production").Key(dbsection + ".username").String(),
		db:       cfg.Section("production").Key(dbsection + ".dbname").String(),
		charset:  cfg.Section("production").Key(dbsection + ".charset").String(),
		timeout:  cfg.Section("production").Key(dbsection + ".timeout").String(),
	}

}

/*获取当前文件执行的路径*/
func GetCurPath() string {
	file, _ := exec.LookPath(os.Args[0])

	//得到全路径，比如在windows下E:\\golang\\test\\a.exe
	path, _ := filepath.Abs(file)

	rst := filepath.Dir(path)

	return rst
}

/*
 * redis 链接
 */
func GetRedis() redis.Conn {
	config := GetRedisconfig("config.ini", "shardredis")
	c, err := redis.Dial("tcp", config.host+":"+config.port)
	if err != nil {
		fmt.Println("Connect to redis error", err)
	}
	tpass := config.password
	//解密密码
	//auth认证
	if _, err := c.Do("AUTH", tpass); err != nil {
		fmt.Println("auth error to redis ", err)
	}

	if _, err := c.Do("select", config.db); err != nil {
		fmt.Println("auth error to redis ", err)
	}
	return c
}

/***
 * 获取mysql 链接
 */
func Getdb(secname string) (db *sql.DB, errstr error) {

	//db, err := sql.Open("mysql", "root:123456qwe@tcp(127.0.0.1:3306)/common?parseTime=true")

	config := Getdbconfig("config.ini", "commondb")
	//fmt.Println(config)
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", config.username, config.password, config.host, config.port, config.db, config.charset))
	if err != nil {
		fmt.Println("\nconnection mysql error")
		return db, err
	} //defer db.Close()
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	err = db.Ping()
	if err != nil {
		fmt.Println("\nopen mysql error ", err)
		return db, err
	}
	return db, nil

}

func HttpGetcheck(url string) (status string, err error) {
	client := &http.Client{}
	resp, err := client.Get(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return string(body), nil
}

//通用json 返回定义
type commonRes struct {
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	Data     string `json:"data"`
	HttpCode int    `json:"httpcode"`
}

//公共json 返回
func ComRes(code int, msg, data string, httpcode int) string {
	res1 := commonRes{
		Code:     code,
		Msg:      msg,
		Data:     data,
		HttpCode: httpcode,
	}

	r1j, err := json.Marshal(res1)

	if err != nil {
		return "get json error"
	}
	return string(r1j)
}
