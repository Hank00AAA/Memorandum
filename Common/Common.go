package Common

import (
	"encoding/json"
	"net"
)

//HTTP接口应答
type Response struct{
	Errno int `json:"errno"` //OK 0
	Msg   string `json:"msg"`
	Data  interface{}`json:"data"`
}

//文件上传接口应答
type UploadFileResp struct{
	Errno int `json:"errno"`
	MsgArr []string `json:"data"`
}

//数据批次
type DataBatch struct{
	Data []interface{} //多条日志
}

//1. 注册接口应答


//2. 登陆接口应答
type PList struct{
	PListID string `json:"plistid"`
	PListName string `json:"plistname"`
}

type SignInResp struct{
	Errno int  `json:"errno"`
	Data  interface{} `json:"data"`
}

func BuildSignInResp(errno int, listArr []PList)(resp []byte, err error){

	//1. 定义Resp
	var(
		response SignInResp
	)

	response.Errno = errno
	response.Data = listArr

	//2. 序列化
	resp, err = json.Marshal(response)

	return

}

//3. 根据标签查询条目应答

//文件上传接口应答
func BuildUploadFileResp(errno int, msg []string)(resp []byte, err error){

	var(
		response UploadFileResp
	)

	response.Errno = errno
	response.MsgArr = msg

	//序列化
	resp, err = json.Marshal(response)
	return

}

func BuildResponse(errno int, msg string, data interface{})(resp []byte, err error){

	var(
		response Response
	)

	response.Errno = errno
	response.Msg = msg
	response.Data = data

	//序列化
	resp, err = json.Marshal(response)

	return
}

//获取本机网卡ip
func GetLocalIP()(ipv4 string, err error){

	var(
		addrs []net.Addr
		addr net.Addr
		ipNet *net.IPNet
		isIPNet bool
	)

	//获取所有网卡
	if addrs, err = net.InterfaceAddrs();err!=nil{
		return
	}

	//取第一个非localhost的网卡IP
	for _, addr = range addrs{
		//ipv4, ipv6
		//判断网络地址是否为ip
		//过滤掉环回地址
		if ipNet, isIPNet = addr.(*net.IPNet);isIPNet&&!ipNet.IP.IsLoopback(){
			//跳过ipv6
			if ipNet.IP.To4()!=nil {
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}

	err = ERR_NO_LOCAL_ANY_IP_FOUND

	return
}
