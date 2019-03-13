package Common

//标签字段全部小写
// 查询的时候是用bson小写查询的
//需要人为添加主键

//用户表User
type User struct{
	ID   string    		`json:"_id"      bson:"_id"`
	//mongodb唯一主键 ID = userID = Email
	UserID string  		`json:"userid"   bson:"userid"`
	Email  string 		`json:"email"    bson:"email"`
	PassWord string 	`json:"password" bson:"password"`
}

//个人清单列表
type PMemList struct{
	ID   string    		`json:"_id"      bson:"_id"`
	//mongodb唯一主键 ID = ListID
	ListID string `json:"listid" bson:"listid"`
	ListName string `json:"listname" bson:"listname"`
	UserID   string `json:"userid" bson:"userid"`
}

//条目
type Entry struct{
	ID   string    		`json:"_id"      bson:"_id"`
	//mongodb唯一主键 ID = EntryID
	EntryID string `json:"entryid" bson:"entryid"`
	EntryName string  `json:"entryname" bson:"entryname"`
	ListID   string `json:"listid" bson:"listid"`
	State    int  `json:"state" bson:"state"` // Normal:0 Delete:1
	Version  int  `json:"version" bson:"version"`
}

//步骤
type Step struct{
	ID   string    		`json:"_id"      bson:"_id"`
	//mongodb唯一主键 ID = StepID
	StepID string `json:"stepid" bson:"stepid"`
	EntryID string `json:"entryid" bson:"entryid"`
	Sequence int   `json:"sequence" bson:"sequence"`
	StepName string `json:"stepname" bson:"stepname"`
	Date    string `json:"date" bson:"date"`
	Importance int `json:"importance" bson:"importance"`
	Done       int `json:"done" bson:"done"`
	Content    string `json:"content" bson:"content"`
}

//成员列表
type TMemList struct{
	ID   string    		`json:"_id"      bson:"_id"`
	//mongodb唯一主键  ID = ListID + UserID
	ListID string `json:"listid" bson:"listid"`
	ListName string `json:"listname" bson:"listname"`
	UserID   string `json:"userid" bson:"userid"`
}