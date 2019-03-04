package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	uploadPath = "/Users/hhs/go/src/github.com/Hank00AAA/Memorandum/test/main/Files/"
)

func main() {
	http.HandleFunc("/upload", uploadHandle)
	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/Files/", http.StripPrefix("/Files", fs))
	log.Fatal(http.ListenAndServe(":8037", nil))
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {
	file, head, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	filePath := uploadPath + head.Filename
	fW, err := os.Create(filePath)
	if err != nil {
		fmt.Println("文件创建失败", err)
		return
	}
	defer fW.Close()
	_, err = io.Copy(fW, file)
	if err != nil {
		fmt.Println("文件保存失败", err)
		return
	}
	io.WriteString(w, "save to "+filePath)
}

//test
/*
if err = user.Insert(&Common.User{
	ID:"3",
	UserID: "1",
	Email: "1222",
	PassWord:"4",
});err!=nil{
	fmt.Println(err)
}
*/

/*

if _, err = user.RemoveAll(bson.M{"password":"1"});err!=nil{
	fmt.Println(err)
}else{
	fmt.Println("233")
}
*/

/*
var result []Common.User
if err = user.Find(bson.M{"userid":"1"}).All(&result);err!=nil{
	return
}

for _, tmp := range result{
	fmt.Println(tmp.PassWord)
}
*/