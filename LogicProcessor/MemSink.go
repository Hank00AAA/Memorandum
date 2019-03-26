package LogicProcessor

import (
	"errors"
	"fmt"
	"github.com/Hank00AAA/Memorandum/Common"
	"go.mongodb.org/mongo-driver/bson"
	"gopkg.in/mgo.v2"
	"strconv"
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

//6. 查询团队清单列表
func (memSink *MemSink)getTMemList(email string)(Resps *[]Common.TListInfo, err error){

	var(
		tListInfo Common.TListInfo
		tMemList_tmp Common.TMemList
		tMemListResult []Common.TMemList
		resp []Common.TListInfo
	)

	resp = make([]Common.TListInfo, 0)

	if err = memSink.MC_TMemList.Find(bson.M{"userid":email}).All(&tMemListResult);err!=nil{
		return nil, err
	}

	for _, tMemList_tmp = range tMemListResult{
		tListInfo.TListName = tMemList_tmp.ListName
		tListInfo.TListID   = tMemList_tmp.ListID
		resp = append(resp, tListInfo)
	}

	return &resp, nil
}

//7. 查询条目
func (memSink *MemSink)getEntryAndStep(listid string)(resp *[]Common.EntryAndStep, err error){

	var(
		entryResult []Common.Entry
		entryTmp    Common.Entry
		stepResult  []Common.Step
		stepTmp     Common.Step
		singleStep  Common.SingleStep
		resp_       []Common.EntryAndStep
		singleResp  Common.EntryAndStep
	)

	//查到listid 对应的entry
	if err = memSink.MC_Entry.Find(bson.M{"listid":listid}).All(&entryResult);err!=nil{
		return nil,err
	}

	//根据entry查询对应的step，构造返回参数
	for _, entryTmp = range entryResult{

		singleResp.EntryID = entryTmp.EntryID
		singleResp.EntryName = entryTmp.EntryName
		singleResp.EntryVersion = entryTmp.Version
		singleResp.StepArr = make([]Common.SingleStep, 0)

		if err = memSink.MC_Step.Find(bson.M{"entryid":entryTmp.EntryID}).All(&stepResult);err!=nil{
			return nil, err
		}

		for _, stepTmp = range stepResult{
			singleStep.Date = stepTmp.Date
			singleStep.Importance = stepTmp.Importance
			singleStep.StepID = stepTmp.StepID

			singleResp.StepArr = append(singleResp.StepArr, singleStep)
		}

		resp_ = append(resp_, singleResp)
	}

	return &resp_, nil
}

//8. 添加个人清单
//查询是否存在唯一用户，如果不存在则返回错误（用户不存在）
//用雪花算法生成listI，注意
// 1. ID = LISTID
// 2. UserName = Email
//插入完成后返回listid和listname
func (memSink *MemSink)addPMemList(email string, listName string)(pmListResp *Common.AddPMLResp, err error){

	var(
		userResult []Common.User
		listID     string
		id         int64
		insertPMemList Common.PMemList
	)

	//查询是否存在唯一用户
	if err = memSink.MC_User.Find(bson.M{"email":email}).All(&userResult);err!=nil{
		return nil, err
	}

	if len(userResult) == 0{
		return nil, Common.ERR_ACCOUNT_DONT_EXIST
	}else if len(userResult) > 1{
		return nil, Common.ERR_MULTI_EMAIL_EXIST
	}

	//雪花算法生成唯一I
	id = G_Node.Generate()
	listID = strconv.FormatInt(id, 10)
	fmt.Println("listid: ",listID)

	//写入
	insertPMemList.ListID = listID
	insertPMemList.ID 	  = listID
	insertPMemList.ListName = listName
	insertPMemList.UserID  = email

	if err = memSink.MC_PMemList.Insert(&insertPMemList);err!=nil{
		return nil, err
	}

	return	&Common.AddPMLResp{ListName:listName, ListID:listID}, nil
}

//9. 添加团队清单
func (memSink *MemSink)addTMemList(email string, listname string)(resp *Common.AddTMLResp, err error){

	var(
		userResult []Common.User
		listID     string
		id         int64
		insertTMemList Common.TMemList
		)

	//查询是否存在用户
	if err = memSink.MC_User.Find(bson.M{"email":email}).All(&userResult);err!=nil{
		return nil, err
	}

	if len(userResult) == 0{
		return nil, Common.ERR_ACCOUNT_DONT_EXIST
	}else if len(userResult) > 1{
		return nil, Common.ERR_MULTI_EMAIL_EXIST
	}

	//雪花算法生成唯一I
	id = G_Node.Generate()
	listID = strconv.FormatInt(id, 10)
	fmt.Println("listid: ",listID)

	//构造团队清单
	insertTMemList.ListID = listID
	insertTMemList.ID	  = listID
	insertTMemList.ListName = listname
	insertTMemList.UserID  = email

	//插入数据库
	if err = memSink.MC_TMemList.Insert(&insertTMemList);err!=nil{
		return nil, err
	}

	return &Common.AddTMLResp{ListName:listname, ListID:listID}, nil
}

//10. 根据条目id获取步骤
func (memSink *MemSink)getSteps(entryID string)(result []Common.Step, err error){

	var(

	)

	if err = memSink.MC_Step.Find(bson.M{"entryid":entryID}).All(&result);err!=nil{
		return nil, err
	}

	return result, err
}

//11. 保存条目
//首先看是否是 -1 ， 如果是-1，说明是新条目，则按照添加流程走
//不是-1
//则用乐观锁，删除所有之前的step，然后插入现在的step，version+1
//如果isOK = false，err！=nil，说明已经有人先更新了

func (memSink *MemSink)saveEntry(data *Common.ReqData)(isOK bool , entryID string, err error){

	var(
		step_temp Common.Step
		entryLock *EntryLock
		iEntry Common.Entry
		queryEntry []Common.Entry
		oldSteps []Common.Step
	)

	err = nil

	for i:=0;i < len(data.StepArr);i++{
		data.StepArr[i].StepID = strconv.FormatInt(G_Node.Generate(), 10)
		data.StepArr[i].ID = data.StepArr[i].StepID
	}

	if(data.EntryID == "-1"){
		//新建entryI，然后把step的entrid全赋值上
		//抢分布式锁，插入新step
		data.EntryID = strconv.FormatInt(G_Node.Generate(),10)
		fmt.Println(data.EntryID)

		for i:=0; i < len(data.StepArr); i++{
			data.StepArr[i].EntryID = data.EntryID
		}

		//抢锁
		entryLock = InitEntryLock(data.EntryID, G_register.kv, G_register.lease)
		if err = entryLock.TryLock();err!=nil{
			goto ERR
		}
		defer entryLock.Unlock()

		//如果锁被抢了，返回false
		if entryLock.isLocked == false{
			err = Common.ERR_LOCK_HAS_EXISTED
			goto ERR
		}

		//写入数据
		//写入entry
		iEntry.ID = data.EntryID
		iEntry.EntryID = data.EntryID
		iEntry.ListID = data.ListID   //要不要检查下存不存在？
		iEntry.EntryName = data.EntryName
		iEntry.Version = data.Version
		iEntry.State = data.State
		if err = G_memSink.MC_Entry.Insert(&iEntry);err!=nil{
			goto ERR
		}
		//写入step
		for _, step_temp = range data.StepArr{
			if err = G_memSink.MC_Step.Insert(step_temp);err!=nil{
				goto ERR
			}
		}
		return true, data.EntryID,nil
	} else{
		//防止前端打错
		for i:=0;i < len(data.StepArr);i++{
			data.StepArr[i].EntryID = data.EntryID
		}

		//先检查entryid是否存在，不存在则可能被删除，让前端确认
		//抢分布式锁，删除旧step，然后插入新step
		if err = G_memSink.MC_Entry.Find(bson.M{"entryid":data.EntryID}).All(&queryEntry);err!=nil{
			goto ERR
		}

		if len(queryEntry) == 0{
			err = Common.ERR_ENTRY_DONT_EXIST
			goto ERR
		}

		if len(queryEntry) > 1 {
			fmt.Println("entry count > 1", queryEntry)
		}

		//抢锁
		entryLock = InitEntryLock(data.EntryID, G_register.kv, G_register.lease)
		if err = entryLock.TryLock();err!=nil{
			goto ERR
		}
		defer entryLock.Unlock()

		//如果锁被抢了，返回false
		if entryLock.isLocked == false{
			err = Common.ERR_LOCK_HAS_EXISTED
			goto ERR
		}

		//先检查version字段
		//必须大于当前version才能更新
		if queryEntry[0].Version >= data.Version{
			err = Common.ERR_VERSION_IS_SMALLER
			fmt.Println(err.Error())
			goto ERR
		}

		//查询旧step
		if err = G_memSink.MC_Step.Find(bson.M{"entryid":data.EntryID}).All(&oldSteps);err!=nil{
			goto ERR
		}

		//更新新step
		for _, step_temp = range data.StepArr{
			if err = G_memSink.MC_Step.Insert(step_temp);err!=nil{
				goto ERR
			}
		}

		//更新entry的version
		if err = G_memSink.MC_Entry.Update(bson.M{"entryid":data.EntryID}, bson.M{
			"entryid":queryEntry[0].EntryID,
			"entryname":queryEntry[0].EntryName,
			"listid":queryEntry[0].ListID,
			"state":queryEntry[0].State,
			"version":data.Version});
		err!=nil{
			goto ERR
		}

		//删除旧step
		for _, step_temp = range oldSteps{
			if _, err = G_memSink.MC_Step.RemoveAll(bson.M{"stepid":step_temp.StepID});err!=nil{
				goto ERR
			}
		}
	}
	return true, data.EntryID, err

	ERR:
		return false, data.EntryID, err

}


//12. 删除条目
func (memSink *MemSink)deleteEntry(entryID string)(isOk bool, err error){

	var(

	)

	if _, err = memSink.MC_Entry.RemoveAll(bson.M{"entryid":entryID});err!=nil{
		return false, err
	}

	if _, err = memSink.MC_Step.RemoveAll(bson.M{"entryid":entryID});err!=nil{
		return false, err
	}

	return true, nil

}

//13. 查询团队成员
func (memSink *MemSink)getTMemberByListID(tMemListID string)(emails []string, err error){
	var(
		tMemList []Common.TMemList
		email_tmp Common.TMemList
	)

	emails = make([]string, 0)

	if err = memSink.MC_TMemList.Find(bson.M{"listid":tMemListID}).All(&tMemList);err!=nil{
		return nil, err
	}

	for _, email_tmp = range tMemList{
		emails = append(emails, email_tmp.UserID)
	}

	return emails, nil
}

//14. 添加团队成员
func (memSink *MemSink)addTMember(tMemListID string, email string)(isok bool , err error){
	var(
		user []Common.User
		tMemList []Common.TMemList
	)

	if err = memSink.MC_User.Find(bson.M{"userid":email}).All(&user);err!=nil{
		return false, err
	}

	if len(user)==0{
		return false, errors.New("帐号不存在")
	}

	if len(user) > 1{
		fmt.Println("存在多个邮箱相同的帐号")
	}

	if err = memSink.MC_TMemList.Find(bson.M{"listid":tMemListID}).All(&tMemList);err!=nil{
		return false, err
	}

	if len(tMemList) == 0{
		return false, errors.New("ID不存在")
	}

	if err = memSink.MC_TMemList.Insert(&Common.TMemList{
		ID: tMemListID+email,
		ListID:tMemListID,
		ListName:tMemList[0].ListName,
		UserID:email,
	});err!=nil{
		return false, err
	}

	return true, nil
}

//15. 删除团队成员
func (memSink *MemSink)deleteTMember(tMemListID string, email string)(isok bool , err error){
	var(
		user []Common.User
	)

	if err = memSink.MC_User.Find(bson.M{"userid":email}).All(&user);err!=nil{
		return false, err
	}

	if len(user)==0{
		return false, errors.New("帐号不存在")
	}

	if len(user) > 1{
		fmt.Println("存在多个邮箱相同的帐号")
	}

	if _, err = memSink.MC_TMemList.RemoveAll(bson.M{"listid":tMemListID, "userid":email});err!=nil{
		return false, err
	}

	return true, nil
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

	G_memSink.addPMemList("111@qq.com","11")

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