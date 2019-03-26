package LogicProcessor

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"time"
)

//产生token
func CreateToken()( string){

	crutime := time.Now().Unix()

	h := md5.New() //创造md5对象

	io.WriteString(h, strconv.FormatInt(crutime, 10)) //将unix时间转化成字符串，然后写入md5对象

	token := fmt.Sprintf("%x", h.Sum(nil)) //生成md5串

	return token
}

//检查token是否过期/存在
//如果过期则删除,返回false
//如果没过期则返回true
func CheckToken(email string, token string)(isOK bool, err error){
	isOK, err = G_register.checkEmail_Token(email, token)
	return
}

//将邮箱和token绑定
//将邮箱-token键值对写入etcd，加入lease
func BindEmailWithToken(email string, token string)(isOK bool, err error){
	isOK, err = G_register.writeToken(email, token)
	return
}

// 将邮箱和新token绑定
func BindEmailWithNewToken(email string)(token string, isOK bool, err error){
	token = CreateToken()
	isOK, err = BindEmailWithToken(email, token)
	return
}

