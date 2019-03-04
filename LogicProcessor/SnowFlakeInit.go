package LogicProcessor

import "github.com/Hank00AAA/Memorandum/Common"

var(
	G_Node *Common.Node
)

func InitNode(node int64)(err error){
	var(
		node_instance *Common.Node
	)

	if node_instance, err = Common.NewNode(node);err!=nil{
		return
	}

	//单例
	G_Node = node_instance

	return
}
