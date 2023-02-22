package utils

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type ServerConfig struct {
	SqlUser      string
	SqlPw        string
	SqlAddr      string
	Schema       string
	RedisAddr    string
	RedisPw      string
	ServerDomain string
}

func printConfigContent() {
	fmt.Printf(
		"#数据库用户名\n" +
			"SqlUser: \n" +
			"#数据库密码\n" +
			"SqlPw: \n" +
			"#数据库地址 如 1.1.1.1:3306, hhh.com:3306\n" +
			"SqlAddr: \n" +
			"Schema: \n" +
			"#redis地址\n" +
			"RedisAddr: \n" +
			"#redis密码\n" +
			"RedisPw: \n" +
			"ServerDomain: \n")
}
func LoadConfig() (error, *ServerConfig) {
	dataBytes, err := os.ReadFile("config.yaml")
	if err != nil {
		//打印配置文件模版
		fmt.Printf("\n在当前运行位置未找到文件" +
			"\n")
		printConfigContent()
		return err, nil
	}
	//return nil, nil
	//fmt.Println("yaml 文件的内容: \n", string(dataBytes))
	config := ServerConfig{}
	configmap := map[string]string{}
	err = yaml.Unmarshal(dataBytes, &configmap)
	if err != nil {
		fmt.Println("解析 config.yaml 文件失败：", err)
		return err, nil
	}

	config.SqlPw = configmap["SqlPw"]
	config.SqlUser = configmap["SqlUser"]
	config.SqlAddr = configmap["SqlAddr"]
	config.Schema = configmap["Schema"]
	config.RedisAddr = configmap["RedisAddr"]
	config.RedisPw = configmap["RedisPw"]
	config.ServerDomain = configmap["ServerDomain"]

	//fmt.Printf("config:%+v\n", configmap)
	return nil, &config
}
