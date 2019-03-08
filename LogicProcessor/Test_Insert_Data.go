package LogicProcessor

import (
	"fmt"
	"github.com/Hank00AAA/Memorandum/Common"
)

func Insert_Data()(err error){

	var(
		user_id string = "111@qq.com"
		pMmemList_id string = "PTL1"
		pMemList_id_2 string = "PTL2"
		tMemList_id string = "TML1"
		entryID string = "test_entry_1"
		entryID_2 string = "test_entry_2"
		outOfWeekEntry string = "test_out_of_week_entry"
		tentryID  string = "TMem_test_1"
		step_ID string = "test_step_1"
		step_ID_2 string = "test_step_2"
		step_ID_3 string = "test_step_3"
		step_ID_4 string = "test_step_4"
		step_ID_outofWeek = "test_outof_week_step"
		step_Team_ID_1 string = "test_TMemList_step_1"

	)

	//Insert User
	if err = G_memSink.MC_User.Insert(&Common.User{
		ID: user_id,
		UserID:user_id,
		Email:user_id,
		PassWord:"111",
	});err!=nil{
		fmt.Println(err)
	}

	//Insert PMemlist
	//list1
	if err = G_memSink.MC_PMemList.Insert(&Common.PMemList{
		ID:pMmemList_id,
		ListID:pMmemList_id,
		ListName:pMmemList_id,
		UserID:user_id,
	});err!=nil{
		fmt.Println(err)
	}

	//list2
	if err = G_memSink.MC_PMemList.Insert(&Common.PMemList{
		ID:pMemList_id_2,
		ListID:pMemList_id_2,
		ListName:pMemList_id_2,
		UserID:user_id,
	});err!=nil{
		fmt.Println(err)
	}

	//Insert Entry
	//PEntry:1
	if err = G_memSink.MC_Entry.Insert(&Common.Entry{
		ID:entryID,
		EntryID:entryID,
		EntryName:entryID,
		ListID:pMmemList_id,
		State:0,
		Version:0,
	});err!=nil{
		fmt.Println(err)
	}

	//pEntry:2
	if err = G_memSink.MC_Entry.Insert(&Common.Entry{
		ID:entryID_2,
		EntryID:entryID_2,
		EntryName:entryID_2,
		ListID:pMemList_id_2,
		State:0,
		Version:0,
	});err!=nil{
		fmt.Println(err)
	}

	//pEntry:outofWeek 3
	if err = G_memSink.MC_Entry.Insert(&Common.Entry{
		ID:outOfWeekEntry,
		EntryID:outOfWeekEntry,
		EntryName:outOfWeekEntry,
		ListID:pMmemList_id,
		State:0,
		Version:0,
	});err!=nil{
		fmt.Println(err)
	}

	//tEntry:1
	if err = G_memSink.MC_Entry.Insert(&Common.Entry{
		ID:tentryID,
		EntryID:tentryID,
		EntryName:"test Team Member List 1",
		ListID:tMemList_id,
		State:0,
		Version:0,
	});err!=nil{
		fmt.Println(err)
	}

	//Insert Step
	if err = G_memSink.MC_Step.Insert(&Common.Step{
		ID:	step_ID,
		StepID: step_ID,
		EntryID:entryID,
		Sequence:0,
		StepName:"test",
		Date:"2019-03-07",
		Importance:1,
		Done:0,
		Content:"AAAA",
	});err!=nil{
		fmt.Println(err)
	}

	if err = G_memSink.MC_Step.Insert(&Common.Step{
		ID:	step_ID_2,
		StepID: step_ID_2,
		EntryID:entryID,
		Sequence:0,
		StepName:"test2",
		Date:"2019-03-07",
		Importance:1,
		Done:0,
		Content:"BBB",
	});err!=nil{
		fmt.Println(err)
	}

	if err = G_memSink.MC_Step.Insert(&Common.Step{
		ID:	step_ID_3,
		StepID: step_ID_3,
		EntryID:entryID,
		Sequence:0,
		StepName:"test",
		Date:"2019-03-08",
		Importance:1,
		Done:0,
		Content:"CCC",
	});err!=nil{
		fmt.Println(err)
	}

	if err = G_memSink.MC_Step.Insert(&Common.Step{
		ID:	step_ID_4,
		StepID: step_ID_4,
		EntryID:entryID_2,
		Sequence:0,
		StepName:step_ID_4,
		Date:"2019-03-07",
		Importance:1,
		Done:0,
		Content:"CCCBBBBB",
	});err!=nil{
		fmt.Println(err)
	}

	//outofWeekStep
	if err = G_memSink.MC_Step.Insert(&Common.Step{
		ID:	step_ID_outofWeek,
		StepID: step_ID_outofWeek,
		EntryID:outOfWeekEntry,
		Sequence:0,
		StepName:step_ID_outofWeek,
		Date:"2019-03-07",
		Importance:1,
		Done:0,
		Content:"out out out",
	});err!=nil{
		fmt.Println(err)
	}

	//TEAM ENTRY
	if err = G_memSink.MC_Step.Insert(&Common.Step{
		ID:step_Team_ID_1,
		StepID: step_Team_ID_1,
		EntryID:tentryID,
		Sequence:1,
		StepName:"test_tmemlist_step_1",
		Date:"2019-03-07",
		Importance:1,
		Done:0,
		Content:"CCCAAA",
	});err!=nil{
		fmt.Println(err)
	}

	//Insert TMemList
	if err = G_memSink.MC_TMemList.Insert(&Common.TMemList{
		ID:tMemList_id,
		ListID:tMemList_id,
		ListName:"Test Tmemlist_1",
		UserID:user_id,
	});err!=nil{
		fmt.Println(err)
	}

	return

}
