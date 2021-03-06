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
type SignUpResp struct{
	Errno int `json:"errno"`
	Data  interface{} `json:"data"`
}

func BuildSignUpResp(errno int, data interface{})(resp []byte, err error){

	var(
		response SignUpResp
	)

	response.Errno = errno
	response.Data  = data

	resp, err = json.Marshal(response)
	return

}


//2. 登陆接口应答
type PList struct{
	PListID string `json:"plistid"`
	PListName string `json:"plistname"`
}

type SignInResp struct{
	Errno int  `json:"errno"`
	Data  interface{} `json:"data"`
	Token string `json:"token"`
}

func BuildSignInResp(errno int, listArr interface{})(resp []byte, err error){

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

func BuildSignInRespWithToken(errno int, listArr interface{}, token string)(resp []byte, err error){

	//1. 定义Resp
	var(
		response SignInResp
	)

	response.Errno = errno
	response.Data = listArr
	response.Token = token

	//2. 序列化
	resp, err = json.Marshal(response)

	return

}

//3. 根据标签查询条目应答
type SearchByTagResp struct{
	Errno int `json:"errno"`
	Data interface{} `json:"data"`
}

type StepResp struct{
	StepID string `json:"stepid"`
	Date   string `json:"date"`
	Importance int `json:"importance"`
}

type EntryResp struct{
	EntryID string `json:"entryid"`
	Entryname string `json:"entryname"`
	Entryversion int `json:"entryversion"`
}

type SearchRespData struct{
	Entryresp EntryResp  `json:"entryresp"`
	Stepresp  []StepResp  `json:"stepresp"`
}

func BuildSearchByTagResp(errno int, data interface{})(resp []byte, err error){
	var(
		response SearchByTagResp
	)

	response.Errno = errno
	response.Data = data

	//序列化
	resp ,err = json.Marshal(response)
	return
}

//5. 查询个人备忘录清单
type PListInfo struct{
	PListID string  `json:"plistid"`
	PListName string `json:"plistname"`
}

func BuildGetPMemListResp(errno int, data interface{})(resp []byte, err error){
	var(
		response SearchByTagResp
	)

	response.Errno = errno
	response.Data = data

	//序列化
	resp ,err = json.Marshal(response)
	return
}

//6. 查询团队备忘录清单
type TListInfo struct{
	TListID string `json:"tlistid"`
	TListName string `json:"tlistname"`
}

type Resp struct{
	Errno int    `json:"errno"`
	Data  interface{} `json:"data"`
}

func BuildGetTMemListResp(errno int, data interface{})(resp []byte, err error){
	var (
		response Resp
	)

	response.Errno = errno
	response.Data = data

	resp, err = json.Marshal(response)
	return
}

//7. 查询条目
type SingleStep struct{
	StepID string  `json:"stepid"`
	Date   string  `json:"date"`
	Importance int `json:"importance"`
}

type EntryAndStep struct{
	EntryID string      `json:"entryid"`
	EntryName string	`json:"entryname"`
	EntryVersion int    `json:"version"`
	StepArr []SingleStep `json:"steparr"`
}

func BuildGetEntryResp(errno int, data interface{})(resp []byte, err error){
	var (
		response Resp
	)

	response.Errno = errno
	response.Data = data

	resp, err = json.Marshal(response)
	return
}

//8. 添加个人清单
type AddPMLResp struct{
	ListID string   `json:"listid"`
	ListName string `json:"listname"`
}

func BuildAddPMemListResp(errno int, data interface{})(resp []byte, err error){
	var(
		response Resp
	)
	response.Errno = errno
	response.Data = data

	resp, err = json.Marshal(response)
	return
}

//9. 添加团队清单
type AddTMLResp struct{
	ListID string `json:"listid"`
	ListName string `json:"listname"`
}

func BuildAddTMemListResp(errno int, data interface{})(resp []byte, err error){

	var(
		response Resp
	)

	response.Errno = errno
	response.Data  = data

	//序列化
	resp, err = json.Marshal(response)
	return
}

//10. 根据entryID查询step
func BuildResp(errno int, data interface{})(resp []byte, err error){

	var(
		response Resp
	)

	response.Errno = errno
	response.Data  = data

	//序列化
	resp, err = json.Marshal(response)
	return
}

//11. 条目保存
type ReqData struct{
	EntryID string `json:"entryid"`
	EntryName string `json:"entryname"`
	ListID string `json:"listid"`
	//State int `json:"state"`
	Version int    `json:"version"`
	StepArr []Step  `json:"steps"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type ReqResp struct{
	Version int `json:"version"`
	EntryID string `json:"entryid"`
	Steps []Step `json:"steps"`
}


//13.根据tmemlistID查询成员
type EmailResult struct{
	Email []string `json:"email" bson:"email"`
}


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
