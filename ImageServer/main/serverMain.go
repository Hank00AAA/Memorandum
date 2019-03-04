package main

import (
	"flag"
	"fmt"
	"github.com/Hank00AAA/Memorandum/ImageServer"
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
	flag.StringVar(&confFile, "config", "./ImageServer.json", "制定参数")
	flag.Parse()
}

func main(){

	var(
		err error
	)
	//初始化环境参数
	initArgs()

	//初始化线程
	initEnv()

	//初始化json配置
	if err = ImageServer.InitConfig(confFile);err!=nil{
		goto ERR
	}


	//初始化api http 服务器
	if err = ImageServer.InitApiServer();err!=nil{
		goto ERR
	}

	//服务注册
	if err = ImageServer.InitRegister();err!=nil{
		goto ERR
	}

	for{
		time.Sleep(1*time.Second)
	}

	return

	ERR:
		fmt.Println(err)
}
