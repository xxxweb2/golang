package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

//创建租约注册服务
type ServiceReg struct {
	client        *clientv3.Client
	lease         clientv3.Lease
	leaseResp     *clientv3.LeaseGrantResponse
	cancleFunc    func()
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
}

func NewServiceReg(addr []string, timeNum int64) (*ServiceReg, error) {
	conf := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	}

	var (
		client *clientv3.Client
	)

	if clientTem, err := clientv3.New(conf); err == nil {
		client = clientTem
	} else {
		return nil, err
	}

	ser := &ServiceReg{
		client: client,
	}
	if err := ser.setLease(timeNum); err != nil {
		return nil, err
	}

	go ser.ListenLeaseRespChan()

	return ser, nil
}

func (this *ServiceReg) setLease(timeNum int64) error {
	lease := clientv3.Lease(this.client)

	leaseResp, err := lease.Grant(context.TODO(), timeNum)
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithCancel(context.TODO())
	leaseRespChan, err := lease.KeepAlive(ctx, leaseResp.ID)

	if err != nil {
		return err
	}

	this.lease = lease
	this.leaseResp = leaseResp
	this.cancleFunc = cancelFunc
	this.keepAliveChan = leaseRespChan

	return nil
}

func (this *ServiceReg) ListenLeaseRespChan() {
	for {
		select {
		case leaseKeepResp := <-this.keepAliveChan:
			if leaseKeepResp == nil {
				fmt.Println("已经关闭续租功能")
				return
			} else {
				fmt.Printf("续租成功\n")
			}
		}
	}
}

//通过租约 注册服务
func (this *ServiceReg) PutService(key, val string) error {
	kv := clientv3.KV(this.client)
	_, err := kv.Put(context.TODO(), key, val, clientv3.WithLease(this.leaseResp.ID))

	return err
}

//撤销租约
func (this *ServiceReg) RevokeLease() error {
	this.cancleFunc()
	time.Sleep(2 * time.Second)
	_, err := this.lease.Revoke(context.TODO(), this.leaseResp.ID)
	return err
}

func main() {
	ser, _ := NewServiceReg([]string{"198.13.34.80:2379"}, 5)
	ser.PutService("/node/111", "heiheihei")
	select {}
}
