package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type ServerConfig struct {
	SqlUser string
	SqlPw   string
	SqlAddr string
	Schema  string
}

func printConfigContent() {
	fmt.Printf(
		"请在运行目录下创建 config.yaml,并填写内容:\n\n" +
			"#数据库用户名\nsql_user:\n" +
			"#数据库密码\nsql_pw:\n" +
			"#数据库地址 如 1.1.1.1:3306, hhh.com:3306 \nsql_addr:\n" +
			"schema:\n")
}
func loadConfig() (error, *ServerConfig) {
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
	err = yaml.Unmarshal(dataBytes, &config)
	if err != nil {
		fmt.Println("解析 config.yaml 文件失败：", err)
		return err, nil
	}
	return nil, &config
}
