package LogicProcessor

import (
	"context"
	"github.com/Hank00AAA/Memorandum/Common"
	"go.etcd.io/etcd/clientv3"
)

type EntryLock struct{
	//etcd客户端
	kv clientv3.KV
	lease clientv3.Lease
	entryId string //锁的entryID
	cancelFunc context.CancelFunc
	leaseId clientv3.LeaseID
	isLocked bool

}

//初始化锁
func InitEntryLock(entryID string, kv clientv3.KV, lease clientv3.Lease)(entryLock *EntryLock){
	entryLock = &EntryLock{
		kv:kv,
		lease:lease,
		entryId:entryID,
	}
	return
}

//上锁函数
func(entryLock *EntryLock)TryLock()(err error){

	var(
		leaseGrantResp *clientv3.LeaseGrantResponse
		cancelCtx context.Context
		cancelFunc context.CancelFunc
		leaseId clientv3.LeaseID
		keepRespChan <- chan *clientv3.LeaseKeepAliveResponse
		txn clientv3.Txn
		lockKey string
		txnResp *clientv3.TxnResponse
	)

	//1. 创建租约
	if leaseGrantResp, err = entryLock.lease.Grant(context.TODO(), 5);err!=nil{
		return
	}

	//2. 自动续租
	//创建context用于取消自动续租
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())

	//租约id
	leaseId = leaseGrantResp.ID
	if keepRespChan, err = entryLock.lease.KeepAlive(cancelCtx, leaseId);err!=nil{
		goto FAIL
	}

	//3. 处理续租应答协程
	go func() {
		var(
			keepResp *clientv3.LeaseKeepAliveResponse
		)

		for{
			select{
			case keepResp = <- keepRespChan:
				if keepResp==nil{
					goto END
				}
			}
		}

		END:
	}()

	//4. 创建事务txn
	txn = entryLock.kv.Txn(context.TODO())

	//锁路径
	lockKey = Common.ENTRY_LOCK + entryLock.entryId

	//5. 事务抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey),"=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	//提交事务
	if txnResp, err = txn.Commit();err!=nil{
		goto FAIL
	}

	//6. 成功则返回，失败回滚，释放租约
	if !txnResp.Succeeded{//锁被占用，if条件不成立
		err = Common.ERR_LOCK_HAS_EXISTED
		goto FAIL
	}

	//抢锁成功
	entryLock.leaseId = leaseId
	entryLock.cancelFunc = cancelFunc
	entryLock.isLocked = true
	return

	FAIL:
		cancelFunc()//取消自动续租
		entryLock.lease.Revoke(context.TODO(), leaseId) //释放租约
		return
}

func (entryLock *EntryLock)Unlock(){
	if entryLock.isLocked{
		entryLock.cancelFunc()
		entryLock.lease.Revoke(context.TODO(), entryLock.leaseId)
	}
}