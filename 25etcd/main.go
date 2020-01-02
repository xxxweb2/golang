package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	config := clientv3.Config{
		Endpoints:   []string{"198.13.34.80:2379"},
		DialTimeout: 10 * time.Second,
	}

	client, err := clientv3.New(config)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	kv := clientv3.NewKV(client)

	//put
	putResp, err := kv.Put(context.TODO(), "/test/key1", "hello world2")
	_, err = kv.Put(context.TODO(), "testspam", "spam")
	fmt.Println("putResp:", putResp, err)

	//get
	getResp, err := kv.Get(context.TODO(), "/test/key1")
	fmt.Println(getResp.Kvs)
	fmt.Println(string(getResp.Kvs[0].Value))

	//获取test目录下所有子元素
	rangeResp, err := kv.Get(context.TODO(), "/test/", clientv3.WithPrefix())
	fmt.Println("rangeResp", rangeResp)

	/////////////lease///////////////////
	lease := clientv3.Lease(client)
	grantResp, err := lease.Grant(context.TODO(), 10)
	kv.Put(context.TODO(), "/test/vanish", "vanish in 10s", clientv3.WithLease(grantResp.ID))
	keepResp, err := lease.KeepAliveOnce(context.TODO(), grantResp.ID)
	fmt.Println("keepResp:", keepResp)

	////////////txn////////////
	k1 := "test/k1"
	v1 := "v1"
	k2 := "test/k1"
	v2 := "v1"
	k3 := "test/k1"
	v3 := "v1"
	k4 := "test/k1"
	v4 := "v1"
	k5 := "test/k1"
	v5 := "v1"

	//txn := kv.Txn(context.TODO())
	kv.Txn(context.TODO()).If(
		clientv3.Compare(clientv3.Value(k1), ">", v1),
		clientv3.Compare(clientv3.Version(k1), "=", 2),
	).Then(
		clientv3.OpPut(k2, v2), clientv3.OpPut(k3, v3),
	).Else(
		clientv3.OpPut(k4, v4), clientv3.OpPut(k5, v5),
	).Commit()

	getResp4, err := kv.Get(context.TODO(), k4)
	getResp2, err := kv.Get(context.TODO(), k2)
	fmt.Println("k4:", getResp4.Kvs, getResp2.Kvs, err)

	watchConfig(client, "config_key", &appConfig)
	appConfig.config1 = "xuxinxin"
	appConfig.config2 = "xuxinxin2"
	value, _ := json.Marshal(appConfig)
	kv.Put(context.TODO(), "config_key", string(value))

	select {

	}
}

type AppConfig struct {
	config1 string
	config2 string
}

var appConfig AppConfig

func watchConfig(clt *clientv3.Client, key string, ss interface{}) {
	watchCh := clt.Watch(context.TODO(), key)
	go func() {
		for res := range watchCh {
			value := res.Events[0].Kv.Value
			if err := json.Unmarshal(value, ss); err != nil {
				fmt.Println("now", time.Now(), "watchConfig err", err)
				continue
			}
			fmt.Println("now", time.Now(), "watchConfig", ss)
		}
	}()
}
