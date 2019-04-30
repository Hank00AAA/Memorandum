package LogicProcessor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Hank00AAA/Memorandum/Common"
	"io/ioutil"
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

	var(
		err error
		bytes []byte
		respbytes []byte
		password string
		email string
		token string
		isOK bool
		url string
		body string
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	//解析表单
	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	email = req.PostForm.Get("email")
	password = req.PostForm.Get("password")


	if token, isOK, err  = BindEmailWithNewToken(email);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("绑定token失败")
		}
		goto ERR
	}

/*
	if err = G_memSink.SingUp(email, password);err!=nil||isOK==false{
		goto ERR
	}
*/

	//不需要发token给前端，直接发邮件给用户即可
	//用户点击邮件后跳转到前端，前端将token，email，password发给我
	//到时候再调用激活帐号接口写入数据库
	//qqmail激活码:fcnevkixyahtbgbg

	url = "http://"+G_config.REdirectURL+"?token="+token+"&email="+email+"&password="+password
	fmt.Println(token)
	fmt.Println(email)
	fmt.Println(password)
	fmt.Println(url)

	body=url

	if err = SendToMail("154630877@qq.com","fcnevkixyahtbgbg","smtp.qq.com:25", email,
		"帐号确认注册邮件",body,"html");err!=nil{
			goto ERR
	}

	if respbytes, err = Common.BuildSignUpResp(0, nil);err==nil{
		resp.Write(respbytes)
	}

	return
ERR:
	fmt.Println(err)
	//异常应答
	if respbytes, err = Common.BuildSignUpResp(-1, err.Error());err==nil{
		resp.Write(respbytes)
	}
}

//18. 激活帐号
//token email password
func handleAccountActivation(resp http.ResponseWriter, req *http.Request){

	var(
		bytes []byte
		err error
		respbytes []byte
		token string
		email string
		password string
		isOK bool
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	//解析表单
	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	token = req.PostForm.Get("token")
	email = req.PostForm.Get("email")
	password = req.PostForm.Get("password")

	//检查token
	if isOK, err = CheckToken(email, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	//写入数据库
	if err = G_memSink.SingUp(email, password);err!=nil||isOK==false{
		goto ERR
	}

	//返回
	if respbytes, err = Common.BuildSignUpResp(0, nil);err==nil{
		resp.Write(respbytes)
	}



	return
ERR:
	fmt.Println(err)
	//异常应答
	if respbytes, err = Common.BuildSignUpResp(-1, err.Error());err==nil{
		resp.Write(respbytes)
	}
}

//2. 登陆服务接口
//流程：
//1. 根据email、password查询是否有对应的记录
//2. 如果有，就查找其对应的个人条目，返回报文
//{ errno: 0 data:[ {plistid: | pListName:  } ]}
//3. 如果没有，返回
//{ errno: -1 data:[nil] }
//POST FORM-DATA
//State: 测试完成
//http://localhost:9000/signin?email=1222&password=4
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
		bytes []byte
		token string
		isOK bool
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}


	plist_resp = make([]Common.PList, 0)

	//解析表单
	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	fmt.Println("get req")

	email = req.PostForm.Get("email")
	password = req.PostForm.Get("password")

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

	if token, isOK, err  = BindEmailWithNewToken(email);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("绑定token失败")
		}
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
		fmt.Println("SignIn finish",temp)
	}

	if respbytes, err = Common.BuildSignInRespWithToken(0, plist_resp, token);err==nil{
		resp.Write(respbytes)
	}

	return

	ERR:
		fmt.Println(err)
		//异常应答
		if respbytes, err = Common.BuildSignInRespWithToken(-1, err.Error(), "");err==nil{
			resp.Write(respbytes)
		}
}

//3. 根据标签查询条目
//根据tag进行查询 0:今天 1:最近七天 2：已完成的 3：已删除的 4：已过期的
//根据查询结果返回
//POST方法 FORM_DATA
//url:http://localhost:9000/searchByTag?email=111@qq.com&tag=1
func handleSearchByTag(resp http.ResponseWriter, req *http.Request){

	var(
		err error
		bytes []byte
		tag string
		email string
		searchArr *[]Common.SearchRespData
		token string
		isOK bool
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	//解析表单
	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	//email   tag
	email = req.PostForm.Get("email")
	tag = req.PostForm.Get("tag")
	token = req.PostForm.Get("token")
	//检查token
	if isOK, err = CheckToken(email, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	searchArr = nil
	if tag == "0"{
		if searchArr, err = G_memSink.getTodayEntry(email);err!=nil{
			goto ERR
		}
	}else if tag == "1"{
		if searchArr,err = G_memSink.getWeekEntry(email);err!=nil{
			goto ERR
		}
	}else if tag == "2"{
		if searchArr, err = G_memSink.getHasDoneEntry(email);err!=nil{
			goto ERR
		}
	}else if tag == "3"{
		if searchArr, err = G_memSink.getHasDeleteData(email);err!=nil{
			goto ERR
		}
	}else if tag =="4"{
		if searchArr, err = G_memSink.getHasTimeOutData(email);err!=nil{
			goto ERR
		}
	}

	if bytes, err = Common.BuildSearchByTagResp(0, searchArr);err==nil{
		resp.Write(bytes)
	}

	return

	ERR:
		if bytes, err = Common.BuildSearchByTagResp(-1, err.Error());err==nil{
			fmt.Println(err)
			resp.Write(bytes)
		}
}

//4. 根据日期查询条目
//state:finish
//POST url: http://localhost:9000/searchByDate?email=111@qq.com&date=2019-03-08
func handleSearchByDate(resp http.ResponseWriter, req *http.Request){

	var(
		err error
		bytes []byte
		email string
		date string
		searchArr *[]Common.SearchRespData
		token string
		isOK bool
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	//解析表单
	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	email = req.PostForm.Get("email")
	date  = req.PostForm.Get("date")
	token = req.PostForm.Get("token")
	//检查token
	if isOK, err = CheckToken(email, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	if searchArr, err = G_memSink.getDateEntry(email, date);err!=nil{
		goto ERR
	}

	if bytes, err = Common.BuildSearchByTagResp(0, searchArr);err==nil{
		resp.Write(bytes)
	}

	return

	ERR:
		if bytes, err = Common.BuildSearchByTagResp(-1, err.Error());err==nil{
			fmt.Println(err)
			resp.Write(bytes)
		}

}

//5. 查询个人备忘录清单
//state:finish
//POST url:http://localhost:9000/getPMemList?email=111@qq.com
func handleGetPMemList(resp http.ResponseWriter, req *http.Request){

	var(
		err error
		bytes []byte
		email string
		result *[]Common.PListInfo
		token string
		isOK bool
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	//解析表单
	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	email = req.PostForm.Get("email")
	token = req.PostForm.Get("token")
	//检查token
	if isOK, err = CheckToken(email, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	if result ,err = G_memSink.getPMemList(email);err!=nil{
		goto ERR
	}

	if bytes, err = Common.BuildGetPMemListResp(0, result);err==nil{
		resp.Write(bytes)
	}

	return

	ERR:
		if bytes, err = Common.BuildGetPMemListResp(-1, err.Error());err==nil{
			fmt.Println(err)
			resp.Write(bytes)
		}

}

//6. 查询团队备忘录清单
//state: finish
//GET URL: http://localhost:9000/getTMemList?email=111@qq.com
func handleGetTMemList(resp http.ResponseWriter, req *http.Request){
	var(
		err error
		bytes []byte
		email string
		result *[]Common.TListInfo
		token string
		isOK bool
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	//解析表单
	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	email = req.PostForm.Get("email")
	token = req.PostForm.Get("token")
	//检查token
	if isOK, err = CheckToken(email, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	if result ,err = G_memSink.getTMemList(email);err!=nil{
		goto ERR
	}

	if bytes, err = Common.BuildGetTMemListResp(0, result);err==nil{
		resp.Write(bytes)
	}

	return

ERR:
	if bytes, err = Common.BuildGetTMemListResp(-1, err.Error());err==nil{
		fmt.Println(err)
		resp.Write(bytes)
	}

}

//7. 查询条目
//传入listid，返回Entry和step
//state: finish
//GET url:http://localhost:9000/getEntry?listid=PTL1
func handleGetEntry(resp http.ResponseWriter, req *http.Request){
	var(
		err error
		listID string
		resps *[]Common.EntryAndStep
		bytes []byte
		token string
		isOK bool
		email string
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	//解析表单
	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	listID = req.PostForm.Get("listid")
	token = req.PostForm.Get("token")
	email = req.PostForm.Get("email")
	//检查token
	if isOK, err = CheckToken(email, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	if resps, err = G_memSink.getEntryAndStep(listID);err!=nil{
		goto ERR
	}

	if bytes, err = Common.BuildGetEntryResp(0, resps);err==nil{
		resp.Write(bytes)
	}

	return

	ERR:
		if bytes, err = Common.BuildGetEntryResp(-1, err.Error());err==nil{
			resp.Write(bytes)
		}

}

//8. 添加个人清单
//State: finish
// URL:http://localhost:9000/addPMemList
//POST form-data
//email
//listname
func handleAddPMemList(resp http.ResponseWriter, req *http.Request){

	var(
		err error
		bytes []byte
		email string
		listname string
		respData *Common.AddPMLResp
		token string
		isOK bool
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	email = req.PostForm.Get("email")
	listname = req.PostForm.Get("listname")
	token = req.PostForm.Get("token")
	//检查token
	if isOK, err = CheckToken(email, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	if respData,  err = G_memSink.addPMemList(email, listname);err!=nil{
		goto ERR
	}

	if bytes, err = Common.BuildAddPMemListResp(0, respData);err==nil{
		resp.Write(bytes)
	}

	return

	ERR:
		if bytes, err = Common.BuildAddPMemListResp(-1, err.Error());err==nil{
			resp.Write(bytes)
		}

}

//9. 添加团队清单
//state:finish
//POST URL:http://localhost:9000/addTMemList
//email
//listname

func handleAddTMemList(resp http.ResponseWriter, req *http.Request){

	var(
		err error
		bytes []byte
		email string
		listname string
		respData *Common.AddTMLResp
		token string
		isOK bool
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	email = req.PostForm.Get("email")
	listname = req.PostForm.Get("listname")
	token = req.PostForm.Get("token")
	//检查token
	if isOK, err = CheckToken(email, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	if respData, err = G_memSink.addTMemList(email, listname);err!=nil{
		goto ERR
	}

	if bytes, err = Common.BuildAddTMemListResp(0, respData);err==nil{
		resp.Write(bytes)
	}

	return

	ERR:
		if bytes, err = Common.BuildAddTMemListResp(-1, err.Error());err==nil{
			resp.Write(bytes)
		}
}

//10. 根据条目id获取步骤
//state:finish
//http://localhost:9000/getSteps?entryid=test_entry_1
//POST
//entryid
func handleGetStep(resp http.ResponseWriter, req *http.Request){
	var(
		err error
		bytes []byte
		entryID string
		result []Common.Step
		token string
		isOK bool
		email string
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	entryID = req.PostForm.Get("entryid")
	token = req.PostForm.Get("token")
	email = req.PostForm.Get("email")
	//检查token
	if isOK, err = CheckToken(email, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	if result, err = G_memSink.getSteps(entryID);err!=nil{
		goto ERR
	}

	if bytes, err = Common.BuildResp(0, result);err==nil{
		resp.Write(bytes)
	}

	return

	ERR:
		if bytes, err = Common.BuildAddTMemListResp(-1, err.Error());err==nil{
			resp.Write(bytes)
		}
}

//json test:
//

//11. 条目保存
//
func handleSaveEntry(resp http.ResponseWriter, req *http.Request){

	var(
		reqContent []byte
		err error
		bytes []byte
		dataUnMar Common.ReqData
		isOK bool
		result Common.ReqResp
		entryID string
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}


	if reqContent, err = ioutil.ReadAll(req.Body);err!=nil{
		goto ERR
	}

	if err = json.Unmarshal(reqContent, &dataUnMar);err!=nil{
		goto ERR
	}

	//检查token
	if isOK, err = CheckToken(dataUnMar.Email, dataUnMar.Token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	fmt.Println(dataUnMar)

	if isOK, entryID, err = G_memSink.saveEntry(&dataUnMar);err!=nil{
		goto ERR
	}

	if isOK == false{
		goto ERR
	}

	result.Version = dataUnMar.Version
	result.EntryID = entryID
	result.Steps = dataUnMar.StepArr

	if bytes, err = Common.BuildAddTMemListResp(0, result);err==nil{
		resp.Write(bytes)
	}

	return
ERR:
	if bytes, err = Common.BuildAddTMemListResp(-1, err.Error());err==nil{
		resp.Write(bytes)
	}
}

//12. 删除条目
func handleDeleteEntry(resp http.ResponseWriter, req *http.Request){

	var(
		err error
		bytes []byte
		entryID string
		isOK bool
		token string
		email string
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	entryID = req.PostForm.Get("entryid")
	token = req.PostForm.Get("token")
	//检查token
	if isOK, err = CheckToken(email, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	if isOK , err = G_memSink.deleteEntry(entryID);err!=nil{
		goto ERR
	}

	if isOK{
		if bytes, err = Common.BuildResp(0, nil);err==nil{
			resp.Write(bytes)
		}
	}

	return

ERR:
	if bytes, err = Common.BuildAddTMemListResp(-1, err.Error());err==nil{
		resp.Write(bytes)
	}
}

//13. 查询团队成员
func handleGetMember(resp http.ResponseWriter, req *http.Request){

	var(
		err error
		bytes []byte
		tMemListID string
		email Common.EmailResult
		token string
		isOK bool
		emails string
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	tMemListID = req.PostForm.Get("tmemlistid")
	token = req.PostForm.Get("token")
	emails = req.PostForm.Get("email")
	//检查token
	if isOK, err = CheckToken(emails, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	if email.Email, err = G_memSink.getTMemberByListID(tMemListID);err!=nil{
		goto ERR
	}

	if bytes, err = Common.BuildAddTMemListResp(0, email);err==nil{
		resp.Write(bytes)
	}

	return

ERR:
	if bytes, err = Common.BuildAddTMemListResp(-1, err.Error());err==nil{
		resp.Write(bytes)
	}
}

//14. 添加团队成员
//state:finish
//POST  http://localhost:9000/addMember?tmemlistid=TML2&email=222@qq.com
//tmemlistid
//email
func handleAddMember(resp http.ResponseWriter, req *http.Request){
	var(
		tmemlistid string
		email string
		email_ string
		err error
		bytes []byte
		isOK bool
		token string
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	tmemlistid = req.PostForm.Get("tmemlistid")
	email      = req.PostForm.Get("email") //自己的email
	email_	   = req.PostForm.Get("email_") //添加人的email
	token = req.PostForm.Get("token")
	//检查token
	if isOK, err = CheckToken(email, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	if isOK, err = G_memSink.addTMember(tmemlistid, email_);!isOK{
		goto ERR
	}

	if bytes, err = Common.BuildAddTMemListResp(0, nil);err==nil{
		resp.Write(bytes)
	}

	return

ERR:
	if bytes, err = Common.BuildAddTMemListResp(-1, err.Error());err==nil{
		resp.Write(bytes)
	}

}

//15. 删除团队成员
//state:finish
//POST http://localhost:9000/deleteMember?tmemlistid=TML1&email=111@qq.com
//tmemlist
//email
func handleDeleteMember(resp http.ResponseWriter, req *http.Request){
	var(
		tmemlistid string
		email string
		email_ string
		err error
		bytes []byte
		isOK bool
		token string
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	tmemlistid = req.PostForm.Get("tmemlistid")
	email      = req.PostForm.Get("email")
	email_	   = req.PostForm.Get("email_")
	token = req.PostForm.Get("token")
	//检查token
	if isOK, err = CheckToken(email, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	if isOK, err = G_memSink.deleteTMember(tmemlistid, email_);!isOK{
		goto ERR
	}

	if bytes, err = Common.BuildAddTMemListResp(0, nil);err==nil{
		resp.Write(bytes)
	}

	return
ERR:
	if bytes, err = Common.BuildAddTMemListResp(-1, err.Error());err==nil{
		resp.Write(bytes)
	}
}

//16.删除个人清单
func handleDeletePMemList(resp http.ResponseWriter, req *http.Request){

	var(
		err error
		bytes []byte
		pMemListId string
		isok bool
		email string
		token string
		isOK bool
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	pMemListId = req.PostForm.Get("pmemlistid")
	email 	   = req.PostForm.Get("email")
	token = req.PostForm.Get("token")
	//检查token
	if isOK, err = CheckToken(email, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	if isok, err = G_memSink.deletePMemList(pMemListId);err!=nil||isok==false{
		goto ERR
	}

	if bytes, err = Common.BuildAddTMemListResp(0, nil);err==nil{
		resp.Write(bytes)
	}


	return
	ERR:
		if bytes, err = Common.BuildAddTMemListResp(-1, err.Error());err==nil{
			resp.Write(bytes)
		}


}

//17. 删除团队清单
func handleDeleteTMemList(resp http.ResponseWriter, req *http.Request){
	var(
		err error
		bytes []byte
		tMemListId string
		isok bool
		email string
		token string
		isOK bool
	)

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")
	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	if err = req.ParseMultipartForm(32<<20);err!=nil{
		goto ERR
	}

	tMemListId = req.PostForm.Get("tmemlistid")
	token = req.PostForm.Get("token")
	email = req.PostForm.Get("email")
	//检查token
	if isOK, err = CheckToken(email, token);err!=nil||isOK==false{
		if isOK == false{
			err = errors.New("Token isOK = false")
		}
		goto ERR
	}

	if isok, err = G_memSink.deleteTMemList(tMemListId);err!=nil||isok==false{
		goto ERR
	}

	if bytes, err = Common.BuildAddTMemListResp(0, nil);err==nil{
		resp.Write(bytes)
	}

	return
ERR:
	if bytes, err = Common.BuildAddTMemListResp(-1, err.Error());err==nil{
		resp.Write(bytes)
	}
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
	//16. 删除个人清单
	mux.HandleFunc("/deletePMemList", handleDeletePMemList)
	//17. 删除团队清单
	mux.HandleFunc("/deleteTMemList", handleDeleteTMemList)
	//18. 激活帐号
	mux.HandleFunc("/accountActivation", handleAccountActivation)

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