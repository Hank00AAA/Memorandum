package LogicProcessor

import (
	"context"
	"github.com/Hank00AAA/Memorandum/Common"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"strconv"
	"time"
)

//注册到etcd: /Mem/LogicProcessor/ + IP + port
type Register struct{
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
	localIP string //本机地址
	localPort string //服务器端口
}

var(
	G_register *Register
)

//注册到 /Mem/LogicalProcessor/+ IP + port
func (register *Register)KeepOnline(){

	var(
		regKey string
		leaseGrantedResp *clientv3.LeaseGrantResponse
		err error
		keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
		keepAliveResp *clientv3.LeaseKeepAliveResponse
		cancelCtx context.Context
		cancelFunc context.CancelFunc
	)

	for{
		//注册路径 IP:PORT
		regKey = Common.LOGIC_PROCESSOR_DIR + register.localIP + ":" + register.localPort

		cancelFunc = nil

		//创建租约
		if leaseGrantedResp, err = register.lease.Grant(context.TODO(), 10);err!=nil{
			goto RETRY
		}

		//自动续租
		if keepAliveChan, err = register.lease.KeepAlive(context.TODO(), leaseGrantedResp.ID);err!=nil{
			goto RETRY
		}

		cancelCtx, cancelFunc = context.WithCancel(context.TODO())

		//注册etcd
		if _, err = register.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseGrantedResp.ID));err!=nil{
			goto RETRY
		}

		//处理续租应答
		for{
			select {
			case keepAliveResp = <- keepAliveChan:
				if keepAliveResp==nil{
					goto RETRY
				}
			}
		}

		RETRY:
			time.Sleep(1*time.Second)
			if cancelFunc!=nil{
				cancelFunc()
			}
	}
}

//查询email和token是否存在
func (register *Register)checkEmail_Token(email string, token string)(isExist bool, err error){
	var(
		getResp *clientv3.GetResponse
		kvPair *mvccpb.KeyValue
	)

	if getResp, err  =register.kv.Get(context.TODO(), Common.TOKEN_SAVE_DIT+email);err!=nil{
		return false, err
	}

	if getResp.Count == 0{
		//没有
		return false, nil
	}

	for _, kvPair = range getResp.Kvs{
		if string(kvPair.Value) == token{
			return true, nil
		}
	}
	return  false, nil
}

//将token写入etcd
func (register *Register)writeToken(email string, token string)(isOK bool, err error){

	var(
		tokenKey string
		leaseGrantedResp *clientv3.LeaseGrantResponse
	)

	//注册路径
	tokenKey = Common.TOKEN_SAVE_DIT + email

	//创建lease
	if leaseGrantedResp, err = register.lease.Grant(context.TODO(), G_config.TokenLease);err!=nil{
		return false, err
	}

	//写入到token目录
	if _, err = register.kv.Put(context.TODO(), tokenKey, token, clientv3.WithLease(leaseGrantedResp.ID));err!=nil{
		return false, err
	}

	return true, nil
}

func InitRegister()(err error){
	var(
		config clientv3.Config
		client *clientv3.Client
		kv clientv3.KV
		lease clientv3.Lease
		localIP string
		localPort string
	)

	//初始化配置
	config = clientv3.Config{
		Endpoints:G_config.EtcdEndpoints,
		DialTimeout:time.Duration(G_config.EtcdDialTimeout)*time.Millisecond,
	}

	//建立链接
	if client, err = clientv3.New(config);err!=nil{
		return
	}

	//得到KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	//本机IP
	if localIP, err = Common.GetLocalIP();err!=nil{
		return
	}

	//服务端口
	localPort = strconv.Itoa(G_config.ApiPort)

	//单例
	G_register = &Register{
		client:client,
		kv:kv,
		lease:lease,
		localIP:localIP,
		localPort:localPort,
	}

	go G_register.KeepOnline()

	/*
	var isOK bool
	//test
	G_register.writeToken("111@qq.com", "222")
	if isOK, err = G_register.checkEmail_Token("1111@qq.com","222");err!=nil{
		fmt.Println(err)
	}

	fmt.Println(isOK)
	*/

	return
}