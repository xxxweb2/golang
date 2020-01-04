package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type EtcdMutex struct {
	Ttl     int64 //租约时间
	Conf    clientv3.Config
	Key     string
	cancel  context.CancelFunc
	lease   clientv3.Lease
	leaseId clientv3.LeaseID
	txn     clientv3.Txn
}

//初始化锁
func (em *EtcdMutex) init() error {
	var err error
	var ctx context.Context
	client, err := clientv3.New(em.Conf)
	if err != nil {
		return err
	}

	em.txn = clientv3.NewKV(client).Txn(context.TODO())
	em.lease = clientv3.NewLease(client)
	leaseResp, err := em.lease.Grant(context.TODO(), em.Ttl)
	if err != nil {
		return err
	}

	ctx, em.cancel = context.WithCancel(context.TODO())
	em.leaseId = leaseResp.ID
	_, err = em.lease.KeepAlive(ctx, em.leaseId)

	return err
}

//获取锁
func (em *EtcdMutex) Lock() error {
	err := em.init()
	if err != nil {
		return err
	}

	em.txn.If(clientv3.Compare(clientv3.CreateRevision(em.Key), "=", 0)).Then(
		clientv3.OpPut(em.Key, "", clientv3.WithLease(em.leaseId))).Else()
	txnResp, err := em.txn.Commit()
	if err != nil {
		return err
	}
	if !txnResp.Succeeded {
		return fmt.Errorf("抢锁失败")
	}

	return nil
}

//释放锁
func (em *EtcdMutex) UnLock() {
	em.cancel()
}

func main() {
	var conf = clientv3.Config{
		Endpoints:   []string{"198.13.34.80:2379"},
		DialTimeout: 5 * time.Second,
	}
	eMutex1 := &EtcdMutex{
		Ttl:  10,
		Conf: conf,
		Key:  "lock",
	}

	eMutex2 := &EtcdMutex{
		Ttl:  10,
		Conf: conf,
		Key:  "lock",
	}

	go func() {
		err := eMutex1.Lock()
		if err != nil {
			fmt.Println("groutine1抢锁失败")
			fmt.Println(err)
			return
		}
		fmt.Println("groutine1抢锁成功")
		defer eMutex1.UnLock()
		time.Sleep(10 * time.Second)
	}()

	go func() {
		err := eMutex2.Lock()
		if err != nil {
			fmt.Println("groutine2抢锁失败")
			fmt.Println(err)
			return
		}
		fmt.Println("groutine2抢锁成功")
		defer eMutex2.UnLock()
	}()
	time.Sleep(30 * time.Second)

}
