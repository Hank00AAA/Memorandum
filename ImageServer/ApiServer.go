package ImageServer

import (
	bytes2 "bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/Hank00AAA/Memorandum/Common"
	"hash"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type ApiServer struct{
	httpServer *http.Server
}

var (
	G_apiServer *ApiServer
)

//上传文件
//post请求
//http://localhost:8980/Upload
//POST form-data 加文件
func handleUploadImage(resp http.ResponseWriter, req *http.Request){
	var(
		file multipart.File
		head *multipart.FileHeader
		err error
		filePath string
		fileExtension string
		lastIndex int
		fw *os.File
		bytes []byte
		md5_hash hash.Hash
		fileData *bytes2.Buffer
		md5_str []byte
		fileMd5Name string
		pathRespArr []string
	)

	fileData = bytes2.NewBuffer([]byte{})

	resp.Header().Set("Access-Control-Allow-Origin", "*")
	resp.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp.Header().Set("content-type", "application/json")

	fmt.Println(req.Method)

	if req.Method == "OPTIONS"{
		if bytes, err = Common.BuildUploadFileResp(0,[]string{"nil"});err==nil{
			resp.Write(bytes)
			fmt.Println("return Options")
			return
		}
	}else{
		fmt.Println(req.Method)
	}

	//获取文件
	file, head, err = req.FormFile("file")
	if err != nil{
		fmt.Println(err.Error() + "get file fail")
		goto ERR
	}
	defer file.Close()

	//截取拓展名
	lastIndex = strings.LastIndex(head.Filename, ".")
	fileExtension = head.Filename[lastIndex:]


	//读取文件
	if _ , err= io.Copy(fileData, file);err!=nil{
		goto ERR
	}

	fmt.Println(fileData.Bytes())

	//md5改名，文件不允许相同的md5
	md5_hash = md5.New()
	md5_hash.Write(fileData.Bytes())
	md5_str = md5_hash.Sum(nil)
	fmt.Println(md5_str)

	// 变成16进制
	fileMd5Name = hex.EncodeToString(md5_str[0:len(md5_str)])

	fmt.Println(fileMd5Name)

	//获取保存路径
	filePath = G_config.AbsDirPath + fileMd5Name + fileExtension
	fmt.Println(filePath)

	//创建同名空文件
	fw, err = os.Create(filePath)
	if err!=nil{
		fmt.Println(err.Error() + "create fail")
		goto ERR
	}
	defer fw.Close()

	//保存文件
	_, err = io.Copy(fw, fileData)
	if err != nil{
		fmt.Println(err.Error() + "copy failed")
		goto ERR
	}

	pathRespArr = make([]string, 0)
	pathRespArr = append(pathRespArr, "http://192.168.1.198:8980/download/"+fileMd5Name + fileExtension)

	//返回信息
	if bytes, err = Common.BuildUploadFileResp(0,pathRespArr);err==nil{
		resp.Write(bytes)
	}
	return

	ERR:
		if bytes, err = Common.BuildUploadFileResp(-1, nil);err==nil{
			resp.Write(bytes)
		}
}



//初始化服务
func InitApiServer()(err error){
	//配置路由
	var(
		mux *http.ServeMux
		listener net.Listener
		httpServer *http.Server
		fs http.Handler
		)

	mux = http.NewServeMux()
	//路由
	mux.HandleFunc("/upload", handleUploadImage)

	//文件下载路由
	fs = http.FileServer(http.Dir(G_config.AbsDirPath))
	mux.Handle("/download/",http.StripPrefix("/download",fs))

	//启动tcp监听
	if listener, err = net.Listen("tcp",":"+strconv.Itoa(G_config.ApiPort));err!=nil{
		return
	}

	//创建http服务器
	httpServer = &http.Server{
		ReadTimeout:time.Duration(G_config.ApiReadTimeout)*time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriterTimeout)*time.Millisecond,
		Handler:mux,
	}

	G_apiServer = &ApiServer{
		httpServer,
	}

	//启动服务器
	go httpServer.Serve(listener)

	return


}