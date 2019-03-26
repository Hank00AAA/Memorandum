package main

import (
	"flag"
	"fmt"
	"github.com/Hank00AAA/Memorandum/LogicProcessor"
	"runtime"
	"time"
)

func initEnv(){
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var(
	confFile string
)

func initArgs(){
	flag.StringVar(&confFile, "config", "./LPConf.json", "制定参数")
	flag.Parse()
}

func main(){

	var(
		err error
	)

	//初始化环境参数
	initArgs()

	//初始化json配置
	if err = LogicProcessor.InitConfig(confFile);err!=nil{
		goto ERR
	}

	//初始化http服务器
	if err = LogicProcessor.InitApiServer();err!=nil{
		goto ERR
	}

	//初始化雪花算法参数
	if err = LogicProcessor.InitNode(LogicProcessor.G_config.NodeID);err!=nil{
		goto ERR
	}

	//初始化数据库
	if err = LogicProcessor.InitMemSink();err!=nil{
		goto ERR
	}

	//服务注册
	if err = LogicProcessor.InitRegister();err!=nil{
		goto ERR
	}

	LogicProcessor.CreateToken()

	//TEST INSERT
	//LogicProcessor.Insert_Data()





	for{
		time.Sleep(1*time.Second)
	}

	return


	ERR:
		fmt.Println(err)
}
