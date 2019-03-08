package LogicProcessor

import (
	"errors"
	"fmt"
	"github.com/Hank00AAA/Memorandum/Common"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/mgo.v2"
	"time"
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

//根据标签查询
//查询今天
//先获取今天的日期,
//首先寻找email对应的所有条目
//然后看找到的条目里是否有步骤是今天
//有则加入查询结果

func (memSink *MemSink)getTodayEntry(email string)(searchArr *[]Common.SearchRespData, err error){
	var(
		today string
		now_time time.Time
		pmList []Common.PMemList
		tmList []Common.TMemList
		pEntryList []Common.Entry
		tEntryList []Common.Entry
		step_list   []Common.Step
		pList_tmp   Common.PMemList
		tlist_tmp   Common.TMemList
		entry_tmp   Common.Entry
		step_tmp    Common.Step
		isToday		bool
		//resp
		singleSearchResp Common.SearchRespData
		singleStepResp   Common.StepResp
		searchRespArr    []Common.SearchRespData
	)

	//获取今天日期 如2015-02-01
	now_time = time.Now()
	today = now_time.String()[:10]
	fmt.Println("today:",today)

	//根据email从PMemList和TmemList里查询清单id，因为email=id=userid，所以可以这样操作
	if err = G_memSink.MC_PMemList.Find(bson.M{"userid":email}).All(&pmList);err!=nil{
		return nil , err
	}

	if err = G_memSink.MC_TMemList.Find(bson.M{"userid":email}).All(&tmList);err!=nil{
		return nil, err
	}

	//然后根据获得的清单id，获取对应的entry
	//pMemList
	for _, pList_tmp = range pmList{
		//根据pmList获取对应的entry
		if err = G_memSink.MC_Entry.Find(bson.M{"listid":pList_tmp.ListID}).All(&pEntryList);err!=nil{
			return nil, err
		}

		//根据清单获得的Entry，寻找所有step，看是否有step的日期是今天，如果是就留下这条entry
		for _, entry_tmp = range pEntryList{
			//根据entryid查询所有step
			if err = G_memSink.MC_Step.Find(bson.M{"entryid":entry_tmp.EntryID}).All(&step_list);err!=nil {
				return nil, err
			}
			singleSearchResp.Stepresp = make([]Common.StepResp, 0)
			isToday = false
			for _, step_tmp = range step_list{
				if step_tmp.Date == today{
					isToday = true
					break
				}
			}

			//如果存在今天的step，那么将entry和对应的step打包返回
			if isToday {
				singleSearchResp.Entryresp.EntryID = entry_tmp.EntryID
				singleSearchResp.Entryresp.Entryversion = entry_tmp.Version
				singleSearchResp.Entryresp.Entryname = entry_tmp.EntryName

				for _, step_tmp = range step_list{
					singleStepResp.Date = step_tmp.Date
					singleStepResp.StepID = step_tmp.StepID
					singleStepResp.Importance = step_tmp.Importance
					singleSearchResp.Stepresp = append(singleSearchResp.Stepresp, singleStepResp)
				}

				searchRespArr = append(searchRespArr, singleSearchResp)
			}
		}
	}



	//tMemList
	for _, tlist_tmp = range tmList{
		//根据pmList获取对应的entry
		if err = G_memSink.MC_Entry.Find(bson.M{"listid":tlist_tmp.ListID}).All(&tEntryList);err!=nil{
			return nil, err
		}

		//根据清单获得的Entry，寻找所有step，看是否有step的日期是今天，如果是就留下这条entry
		for _, entry_tmp = range tEntryList{
			//如果已经删除则不显示
			if entry_tmp.State == 1{
				continue
			}
			singleSearchResp.Stepresp = make([]Common.StepResp, 0)
			//根据entryid查询所有step
			if err = G_memSink.MC_Step.Find(bson.M{"entryid":entry_tmp.EntryID}).All(&step_list);err!=nil {
				return nil, err
			}

			isToday = false
			for _, step_tmp = range step_list{
				if step_tmp.Date == today{
					isToday = true
					break
				}
			}

			//如果存在今天的step，那么将entry和对应的step打包返回
			if isToday {
				singleSearchResp.Entryresp.EntryID = entry_tmp.EntryID
				singleSearchResp.Entryresp.Entryversion = entry_tmp.Version
				singleSearchResp.Entryresp.Entryname = entry_tmp.EntryName

				for _, step_tmp = range step_list{
					singleStepResp.Date = step_tmp.Date
					singleStepResp.StepID = step_tmp.StepID
					singleStepResp.Importance = step_tmp.Importance
					singleSearchResp.Stepresp = append(singleSearchResp.Stepresp, singleStepResp)
				}
				searchRespArr = append(searchRespArr, singleSearchResp)
			}
		}
	}

	//返回数据
	return &searchRespArr, nil

}

func checkIfInAWeek(date string)(isInWeek bool, err error){

	var(
		timeTemplate string = "2006-01-02"
		convTime time.Time
		nowTime time.Time
		timeBeforeWeek time.Time
		timeOffset time.Duration
	)

	//将date字符串转化成时间
	if convTime, err = time.ParseInLocation(timeTemplate, date, time.Local);err!=nil{
		return false, err
	}

	//从现在起往前一周
	nowTime = time.Now()
	if timeOffset, err = time.ParseDuration("-168h");err!=nil{
		return false, err
	}
	timeBeforeWeek = nowTime.Add(timeOffset)

	if convTime.After(timeBeforeWeek) {
		fmt.Println(date, convTime.After(timeBeforeWeek))
		return true, nil
	}else{
		fmt.Println(fmt.Println(date, convTime.After(timeBeforeWeek)))
		return false, nil
	}

	return false, errors.New("Unknown time transform")
}

//根据标签查询
//最近一周
func (memSink *MemSink)getWeekEntry(email string)(searchArr *[]Common.SearchRespData, err error){
	var(
		today string
		now_time time.Time
		pmList []Common.PMemList
		tmList []Common.TMemList
		pEntryList []Common.Entry
		tEntryList []Common.Entry
		step_list   []Common.Step
		pList_tmp   Common.PMemList
		tlist_tmp   Common.TMemList
		entry_tmp   Common.Entry
		step_tmp    Common.Step
		isWeek		bool
		//resp
		singleSearchResp Common.SearchRespData
		singleStepResp   Common.StepResp
		searchRespArr    []Common.SearchRespData
	)

	//获取今天日期 如2015-02-01
	now_time = time.Now()
	today = now_time.String()[:10]
	fmt.Println("today:",today)

	//根据email从PMemList和TmemList里查询清单id，因为email=id=userid，所以可以这样操作
	if err = G_memSink.MC_PMemList.Find(bson.M{"userid":email}).All(&pmList);err!=nil{
		return nil , err
	}

	if err = G_memSink.MC_TMemList.Find(bson.M{"userid":email}).All(&tmList);err!=nil{
		return nil, err
	}

	//然后根据获得的清单id，获取对应的entry
	//pMemList
	for _, pList_tmp = range pmList{
		//根据pmList获取对应的entry
		if err = G_memSink.MC_Entry.Find(bson.M{"listid":pList_tmp.ListID}).All(&pEntryList);err!=nil{
			return nil, err
		}

		//根据清单获得的Entry，寻找所有step，看是否有step的日期是今天，如果是就留下这条entry
		for _, entry_tmp = range pEntryList{
			//根据entryid查询所有step
			if err = G_memSink.MC_Step.Find(bson.M{"entryid":entry_tmp.EntryID}).All(&step_list);err!=nil {
				return nil, err
			}
			singleSearchResp.Stepresp = make([]Common.StepResp, 0)


			isWeek = false
			//判断是否在一周内
			for _, step_tmp = range step_list{
				if isWeek, err = checkIfInAWeek(step_tmp.Date);err==nil&&isWeek{
					break
				}
			}

			//如果存在一周的step，那么将entry和对应的step打包返回
			if isWeek {
				singleSearchResp.Entryresp.EntryID = entry_tmp.EntryID
				singleSearchResp.Entryresp.Entryversion = entry_tmp.Version
				singleSearchResp.Entryresp.Entryname = entry_tmp.EntryName

				for _, step_tmp = range step_list{
					singleStepResp.Date = step_tmp.Date
					singleStepResp.StepID = step_tmp.StepID
					singleStepResp.Importance = step_tmp.Importance
					singleSearchResp.Stepresp = append(singleSearchResp.Stepresp, singleStepResp)
				}

				searchRespArr = append(searchRespArr, singleSearchResp)
			}
		}
	}



	//tMemList
	for _, tlist_tmp = range tmList{
		//根据pmList获取对应的entry
		if err = G_memSink.MC_Entry.Find(bson.M{"listid":tlist_tmp.ListID}).All(&tEntryList);err!=nil{
			return nil, err
		}

		//根据清单获得的Entry，寻找所有step，看是否有step的日期是今天，如果是就留下这条entry
		for _, entry_tmp = range tEntryList{
			//如果已经删除则不显示
			if entry_tmp.State == 1{
				continue
			}
			singleSearchResp.Stepresp = make([]Common.StepResp, 0)
			//根据entryid查询所有step
			if err = G_memSink.MC_Step.Find(bson.M{"entryid":entry_tmp.EntryID}).All(&step_list);err!=nil {
				return nil, err
			}

			isWeek = false
			//判断是否在一周内
			for _, step_tmp = range step_list{
				if isWeek, err = checkIfInAWeek(step_tmp.Date);err==nil&&isWeek{
					break
				}
			}

			//如果存在一周的step，那么将entry和对应的step打包返回
			if isWeek {
				singleSearchResp.Entryresp.EntryID = entry_tmp.EntryID
				singleSearchResp.Entryresp.Entryversion = entry_tmp.Version
				singleSearchResp.Entryresp.Entryname = entry_tmp.EntryName

				for _, step_tmp = range step_list{
					singleStepResp.Date = step_tmp.Date
					singleStepResp.StepID = step_tmp.StepID
					singleStepResp.Importance = step_tmp.Importance
					singleSearchResp.Stepresp = append(singleSearchResp.Stepresp, singleStepResp)
				}
				searchRespArr = append(searchRespArr, singleSearchResp)
			}
		}
	}

	//返回数据
	return &searchRespArr, nil

}


//根据日期查询
func (memSink *MemSink)getDateEntry(email string, date string)(searchArr *[]Common.SearchRespData, err error){
	var(
		pmList []Common.PMemList
		tmList []Common.TMemList
		pEntryList []Common.Entry
		tEntryList []Common.Entry
		step_list   []Common.Step
		pList_tmp   Common.PMemList
		tlist_tmp   Common.TMemList
		entry_tmp   Common.Entry
		step_tmp    Common.Step
		isToday		bool
		//resp
		singleSearchResp Common.SearchRespData
		singleStepResp   Common.StepResp
		searchRespArr    []Common.SearchRespData
	)

	//根据email从PMemList和TmemList里查询清单id，因为email=id=userid，所以可以这样操作
	if err = G_memSink.MC_PMemList.Find(bson.M{"userid":email}).All(&pmList);err!=nil{
		return nil , err
	}

	if err = G_memSink.MC_TMemList.Find(bson.M{"userid":email}).All(&tmList);err!=nil{
		return nil, err
	}

	//然后根据获得的清单id，获取对应的entry
	//pMemList
	for _, pList_tmp = range pmList{
		//根据pmList获取对应的entry
		if err = G_memSink.MC_Entry.Find(bson.M{"listid":pList_tmp.ListID}).All(&pEntryList);err!=nil{
			return nil, err
		}

		//根据清单获得的Entry，寻找所有step，看是否有step的日期是今天，如果是就留下这条entry
		for _, entry_tmp = range pEntryList{
			//根据entryid查询所有step
			if err = G_memSink.MC_Step.Find(bson.M{"entryid":entry_tmp.EntryID}).All(&step_list);err!=nil {
				return nil, err
			}
			singleSearchResp.Stepresp = make([]Common.StepResp, 0)
			isToday = false
			for _, step_tmp = range step_list{
				if step_tmp.Date == date{
					isToday = true
					break
				}
			}

			//如果存在今天的step，那么将entry和对应的step打包返回
			if isToday {
				singleSearchResp.Entryresp.EntryID = entry_tmp.EntryID
				singleSearchResp.Entryresp.Entryversion = entry_tmp.Version
				singleSearchResp.Entryresp.Entryname = entry_tmp.EntryName

				for _, step_tmp = range step_list{
					singleStepResp.Date = step_tmp.Date
					singleStepResp.StepID = step_tmp.StepID
					singleStepResp.Importance = step_tmp.Importance
					singleSearchResp.Stepresp = append(singleSearchResp.Stepresp, singleStepResp)
				}

				searchRespArr = append(searchRespArr, singleSearchResp)
			}
		}
	}



	//tMemList
	for _, tlist_tmp = range tmList{
		//根据pmList获取对应的entry
		if err = G_memSink.MC_Entry.Find(bson.M{"listid":tlist_tmp.ListID}).All(&tEntryList);err!=nil{
			return nil, err
		}

		//根据清单获得的Entry，寻找所有step，看是否有step的日期是今天，如果是就留下这条entry
		for _, entry_tmp = range tEntryList{
			//如果已经删除则不显示
			if entry_tmp.State == 1{
				continue
			}
			singleSearchResp.Stepresp = make([]Common.StepResp, 0)
			//根据entryid查询所有step
			if err = G_memSink.MC_Step.Find(bson.M{"entryid":entry_tmp.EntryID}).All(&step_list);err!=nil {
				return nil, err
			}

			isToday = false
			for _, step_tmp = range step_list{
				if step_tmp.Date == date{
					isToday = true
					break
				}
			}

			//如果存在今天的step，那么将entry和对应的step打包返回
			if isToday {
				singleSearchResp.Entryresp.EntryID = entry_tmp.EntryID
				singleSearchResp.Entryresp.Entryversion = entry_tmp.Version
				singleSearchResp.Entryresp.Entryname = entry_tmp.EntryName

				for _, step_tmp = range step_list{
					singleStepResp.Date = step_tmp.Date
					singleStepResp.StepID = step_tmp.StepID
					singleStepResp.Importance = step_tmp.Importance
					singleSearchResp.Stepresp = append(singleSearchResp.Stepresp, singleStepResp)
				}
				searchRespArr = append(searchRespArr, singleSearchResp)
			}
		}
	}

	//返回数据
	return &searchRespArr, nil

}

//5. 查询个人清单列表
func (memSink *MemSink)getPMemList(email string)(Resps *[]Common.PListInfo, err error){

	var(
		pListInfo Common.PListInfo
		pMemList_tmp Common.PMemList
		pMemListResult []Common.PMemList
		resp []Common.PListInfo
	)

	resp = make([]Common.PListInfo, 0)

	if err = memSink.MC_PMemList.Find(bson.M{"userid":email}).All(&pMemListResult);err!=nil{
		return nil, err
	}

	for _, pMemList_tmp = range pMemListResult{
		pListInfo.PListName = pMemList_tmp.ListName
		pListInfo.PListID   = pMemList_tmp.ListID
		resp = append(resp, pListInfo)
	}

	return &resp, nil
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

/*
	var(
		test_time time.Time
	)
	test_time = time.Now()
	fmt.Println(test_time)
	fmt.Println(test_time.String()[:10])

*/






	return
}