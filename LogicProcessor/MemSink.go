package LogicProcessor

import (
	"github.com/Hank00AAA/Memorandum/Common"
	"github.com/mongodb/mongo-go-driver/bson"
	"gopkg.in/mgo.v2"
)

//mongoDB存储日志
type MemSink struct {
	MemSession *mgo.Session
	MC_User *mgo.Collection
	MC_PMemList *mgo.Collection
	MC_Entry 	*mgo.Collection
	MC_Step		*mgo.Collection
	MC_TMemList *mgo.Collection


}

//单例
var(
	G_memSink *MemSink
)

//登陆：查询是否有此帐号，以及密码是否对应
func (memSink *MemSink)checkWithEmail_Password(email string, password string)(
	UserID string  , isRight bool, err error){

	var(
		results []Common.User
	)

	isRight = false
	if err = G_memSink.MC_User.Find(bson.M{"email":email, "password":password}).All(&results);err!=nil{
		return
	}

	if len()

}

//登陆：根据这个人的Email获取个人清单列表
func (memSink *MemSink)getPMListByUserID(userID string)(plist []Common.PMemList, err error){

}


func InitMemSink()(err error){

	var(
		session *mgo.Session
		user	*mgo.Collection
		pmemlist *mgo.Collection
		entry  *mgo.Collection
		step   *mgo.Collection
		tmemlist *mgo.Collection
	)

	if session, err = mgo.Dial(G_config.MongodbUri);err!=nil{
		return
	}

	//选择db和collection
	user= session.DB("memdb").C("User")
	pmemlist = session.DB("memdb").C("PMemList")
	entry	 = session.DB("memdb").C("Entry")
	step     = session.DB("memdb").C("Step")
	tmemlist = session.DB("memdb").C("TMemList")

	//单例
	G_memSink = &MemSink{
		MC_User:user,
		MC_PMemList:pmemlist,
		MC_Entry:entry,
		MC_Step:step,
		MC_TMemList:tmemlist,
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

	fmt.Println(result)
*/








	return
}