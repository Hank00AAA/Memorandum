package LogicProcessor

import (
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
func handleSignIn(resp http.ResponseWriter, req *http.Request){




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