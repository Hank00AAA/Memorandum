package LogicProcessor

import (
	"fmt"
	"github.com/Hank00AAA/Memorandum/Common"
	"go.mongodb.org/mongo-driver/bson"
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

	fmt.Println(results)

	if len(results) == 0{
		//匹配失败
		fmt.Println("匹配失败")
		return "", false, nil
	}else if len(results) >=1{
		//匹配成功
		if (len(results)<1){
			fmt.Println("存在多条账户记录，请检查")
		}

		return results[0].UserID, true, nil
	}

	return "", false, nil

}

//登陆：根据这个人的userID获取个人清单列表
func (memSink *MemSink)getPMListByUserID(userID string)(plist *[]Common.PMemList, err error){

	var(
		results []Common.PMemList
	)

	if err = G_memSink.MC_PMemList.Find(bson.M{"userid":userID}).All(&results);err!=nil{
		return nil, err
	}

	return &results, nil
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

/* test2
	if err = pmemlist.Insert(&Common.PMemList{
		ID: "1",
		ListID:"1",
		ListName:"test1",
		UserID:"1",
	});err!=nil{
		fmt.Println(err)
	}

	if err = pmemlist.Insert(&Common.PMemList{
		ID: "2",
		ListID:"2",
		ListName:"test2",
		UserID:"1",
	});err!=nil{
		fmt.Println(err)
	}

*/








	return
}