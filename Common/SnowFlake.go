package Common

import (
	"sync"
	"time"
)

//生成全局唯一id
//解决高并发时id生成不重复的问题

const(
	nodebits uint8 = 10 //节点ID的位数 2^10 = 1024
	stepBits uint8 = 12 //序列号位数 2^12 = 4096 每一毫秒可以生成4096个独立ID
	nodeMax  int64 = -1^(-1<<nodebits)
	stepMax  int64 = -1^(-1<<stepBits)
	timeShift uint8 = nodebits + stepBits //时间戳向左偏移量
	nodeShift uint8 = stepBits //节点ID向左偏移量
)

var (
	Epoch int64 = 1300000000000  //随便给个
)



type Node struct{
	mu sync.Mutex //添加互斥锁保证并发安全
	timestamp int64 //时间戳部分
	node      int64 //节点id
	step      int64 //序列号id
}

func NewNode(node int64)(*Node, error){
	//如果超出节点作答范围，产生一个error
	if node < 0||node > nodeMax{
		return nil, ERR_NODE_NUMBER_OVER_MAX
	}

	//生成并返回节点实例指针
	return &Node{
		timestamp:	0,
		node:		node,
		step:		0,
	}, nil
}

func (n *Node) Generate() int64{
	n.mu.Lock()
	defer n.mu.Unlock()

	//获取当前时间
	now := time.Now().UnixNano()/1e6

	if n.timestamp ==now{
		//step 步进1
		n.step ++

		//当前step用完
		if n.step > stepMax{
			//等待本毫秒结束
			for now <= n.timestamp{
				now = time.Now().UnixNano()/1e6
			}
		}
	}else{
		//当前毫秒用完
		n.step = 0
	}

	n.timestamp = now

	//移位运算
	return    int64((now - Epoch) << timeShift | (n.node << nodeShift) | (n.step))
}


