package LogicProcessor

import (
	"fmt"
	"github.com/Hank00AAA/Memorandum/Common"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct{
	httpServer *http.Server
}

var(
	G_apiServer *ApiServer
)

// 1. 注册服务接口
func handleSignUp(resp http.ResponseWriter, req *http.Request){

}

//2. 登陆服务接口
//流程：
//1. 根据email、password查询是否有对应的记录
//2. 如果有，就查找其对应的个人条目，返回报文
//{ errno: 0 data:[ {plistid: | pListName:  } ]}
//3. 如果没有，返回
//{ errno: -1 data:[nil] }
//不知道为什么post方法得不出结果，改成get了
//State: 测试完成
func handleSignIn(resp http.ResponseWriter, req *http.Request){

	var(
		err error
		email string
		password string
		userID string
		isRight bool
		respbytes []byte
		plist *[]Common.PMemList
		plist_temp Common.PMemList
		plist_resp []Common.PList
		temp Common.PList

	)

	plist_resp = make([]Common.PList, 0)

	//解析表单
	if err = req.ParseForm();err!=nil{
		goto ERR
	}

	fmt.Println("get req")

	email = req.Form.Get("email")
	password = req.Form.Get("password")

	fmt.Println(email)
	fmt.Println(password)

	//查询数据库：帐号密码是否对应
	if userID, isRight, err = G_memSink.checkWithEmail_Password(email, password);err!=nil{
		goto ERR
	}

	//帐号密码不对应
	if isRight == false{
		err = Common.ERR_NO_FOUND_ACCOUNT
		goto ERR
	}

	//帐号存在，查询对应的个人清单
	if userID == ""{
		err = Common.ERR_ACCOUNT_IS_NIL
		goto ERR
	}

	//查询个人清单
	if plist, err = G_memSink.getPMListByUserID(userID);err!=nil{
		goto ERR
	}

	//生成应答报文
	fmt.Println(*plist)
	for _, plist_temp = range *plist{
		temp.PListID = plist_temp.ListID
		temp.PListName = plist_temp.ListName
		plist_resp = append(plist_resp, temp)
		fmt.Println("finish",temp)
	}

	if respbytes, err = Common.BuildSignInResp(0, plist_resp);err==nil{
		resp.Write(respbytes)
	}

	return

	ERR:
		fmt.Println(err)
		//异常应答
		if respbytes, err = Common.BuildSignInResp(-1, nil);err==nil{
			resp.Write(respbytes)
		}

}

//3. 根据标签查询条目
func handleSearchByTag(resp http.ResponseWriter, req *http.Request){

}

//4. 根据日期查询条目
func handleSearchByDate(resp http.ResponseWriter, req *http.Request){

}

//5. 查询个人备忘录清单
func handleGetPMemList(resp http.ResponseWriter, req *http.Request){

}

//6. 查询团队备忘录清单
func handleGetTMemList(resp http.ResponseWriter, req *http.Request){

}

//7. 查询条目
func handleGetEntry(resp http.ResponseWriter, req *http.Request){

}

//8. 添加个人清单
func handleAddPMemList(resp http.ResponseWriter, req *http.Request){

}

//9. 添加团队清单
func handleAddTMemList(resp http.ResponseWriter, req *http.Request){

}

//10. 根据条目id获取步骤
func handleGetStep(resp http.ResponseWriter, req *http.Request){

}

//11. 条目保存
func handleSaveEntry(resp http.ResponseWriter, req *http.Request){

}

//12. 删除条目
func handleDeleteEntry(resp http.ResponseWriter, req *http.Request){

}

//13. 查询团队成员
func handleGetMember(resp http.ResponseWriter, req *http.Request){

}

//14. 添加团队成员
func handleAddMember(resp http.ResponseWriter, req *http.Request){

}

//15. 删除团队成员
func handleDeleteMember(resp http.ResponseWriter, req *http.Request){

}

//初始化服务
func InitApiServer()(err error){

	var(
		mux *http.ServeMux
		listener net.Listener
		httpServer *http.Server
	)

	//配置路由
	mux = http.NewServeMux()
	//1. 登陆
	mux.HandleFunc("/signin", handleSignIn)
	//2. 注册
	mux.HandleFunc("/signup", handleSignUp)
	//3. 根据标签查询条目
	 mux.HandleFunc("/searchByTag", handleSearchByTag)
	//4. 根据日期查询条目
	mux.HandleFunc("/searchByDate", handleSearchByDate)
	//5. 查询个人备忘录清单
	mux.HandleFunc("/getPMemList", handleGetPMemList)
	//6. 查询团队备忘录清单
	mux.HandleFunc("/getTMemList", handleGetTMemList)
	//7. 查询条目
	mux.HandleFunc("/getEntry", handleGetEntry)
	//8. 添加个人清单
	mux.HandleFunc("/addPMemList", handleAddPMemList)
	//9. 添加团队清单
	mux.HandleFunc("/addTMemList", handleAddTMemList)
	//10. 根据条目id获取步骤
	mux.HandleFunc("/getSteps", handleGetStep)
	//11. 条目保存,本质是覆盖，先删除掉原来的，然后再创建
	mux.HandleFunc("/saveEntry", handleSaveEntry)
	//12. 删除条目
	mux.HandleFunc("/deleteEntry", handleDeleteEntry)
	//13. 查询团队成员
	mux.HandleFunc("/getMember", handleGetMember)
	//14. 添加团队成员
	mux.HandleFunc("/addMember", handleAddMember)
	//15. 删除团队成员
	mux.HandleFunc("/deleteMember", handleDeleteMember)

	//启动tcp监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort));err!=nil{
		return
	}

	//创建http服务器
	httpServer = &http.Server{
		ReadTimeout:time.Duration(G_config.ApiReadTimeout)*time.Millisecond,
		WriteTimeout:time.Duration(G_config.ApiWriterTimeout)*time.Millisecond,
		Handler:mux,
	}

	G_apiServer = &ApiServer{
		httpServer:httpServer,
	}

	//启动服务器
	go httpServer.Serve(listener)

	return

}