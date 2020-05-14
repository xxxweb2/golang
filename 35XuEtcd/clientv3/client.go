package clientv3

import "errors"

var (
	ErrNoAvailableEndpoints = errors.New("etcdclient: no available endpoints")
	ErrOldCluster           = errors.New("etcdclient: old cluster version")
)

type Client struct {
}

func New(cfg Config) (*Client, error) {
	if len(cfg.Endpoints) == 0 {
		//如果没有填url 直接返回错误
		return nil, ErrNoAvailableEndpoints
	}

	return newClient(&cfg)
}

func newClient(cfg *Config) (*Client, error) {
	//如果是空
	if cfg == nil {
		//创建一个
		cfg = &Config{}
	}

    var creds grpccredentials.TransportCredentials

}
