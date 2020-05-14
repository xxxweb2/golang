package clientv3

import "time"

type Config struct {
	//url列表
	Endpoints []string `json:"endpoints"`
	//超时时间
	DialTimeout time.Duration `json:"dial-timeout"`
}
