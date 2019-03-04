package ImageServer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct{
	ApiPort int `json:"apiPort"`
	ApiReadTimeout int `json:"apiReadTimeout"`
	ApiWriterTimeout int `json:"apiWriteTimeout"`
	AbsDirPath string `json:"absDirPath"`
	EtcdEndpoints []string `json:"etcdEndpoints"`
	EtcdDialTimeout int `json:"etcdDialTimeout"`
}

var(
	G_config *Config
)

func InitConfig(filenames string)(err error){

	var(
		content []byte
		Conf Config
	)

	//1. 把配置文件读取进来
	if content, err = ioutil.ReadFile(filenames);err!=nil{
		return
	}

	//2. 反序列化
	if err = json.Unmarshal(content, &Conf);err!=nil{
		return
	}

	//3. 赋值单例
	G_config = &Conf

	fmt.Println(
		G_config.ApiPort,
		G_config.ApiWriterTimeout,
		G_config.AbsDirPath,
		G_config.EtcdEndpoints,
		G_config.EtcdDialTimeout)

	return


}
