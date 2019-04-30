package LogicProcessor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct{
	ApiPort int `json:"apiPort"`
	ApiReadTimeout int `json:"apiReadTimeout"`
	ApiWriterTimeout int `json:"apiWriteTimeout"`
	EtcdEndpoints []string `json:"etcdEndpoints"`
	EtcdDialTimeout int `json:"etcdDialTimeout"`
	MongodbUri string `json:"mongodbUri"`
	MongodbConnectionTimeout int `json:"mongodbConnectionTimeout"`
	NodeID 					 int64 `json:"nodeId"`
	TokenLease  int64        `json:"tokenLease"`
	REdirectURL string `json:"redirecturl"`
}

var(
	G_config *Config
)

func InitConfig(filenames string)(err error){

	var(
		content []byte
		Conf Config
	)

	//1. 读取配置文件
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
		G_config.ApiReadTimeout,
		G_config.ApiWriterTimeout,
		G_config.NodeID,
		G_config.TokenLease)

	return

}