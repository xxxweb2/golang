package xunsq

import (
	"net"
	"time"
)

type configHandler interface {
	HandlesOption(c *Config, option string) bool
	Set(c *Config, option string, value interface{}) error
	Validate(c *Config) error
}

type defaultsHandler interface {
	SetDefaults(c *Config) error
}

type Config struct {
	initialized    bool
	configHandlers []configHandler
	DialTimeout    time.Duration `opt:"dial_timeout" default:"1s"`

	ReadTimeout  time.Duration `opt:"read_timeout" min:"100ms" max:"5m" default:"60s"`
	WriteTimeout time.Duration `opt:"write_timeout" min:"100ms" max:"5m" default:"1s"`

	LocalAddr net.Addr `opt:"local_addr"`

	LookupPollInterval time.Duration `opt:"lookupd_poll_interval" min:"10ms" max:"5m" default:"60s"`
	LookupPollJitter   float64       `opt:"lookupd_poll_jitter" min:"0" max:"1" default:"0.3"`
}

func NewConfig() *Config {
	c := &Config{
		configHandlers: []configHandler{&structTagsConfig{}, &tlsConfig{}},
		initialized:    true,
	}

	if err := c.setDefaults(); err != nil {
		panic(err)
	}

	return c
}

func (c *Config) setDefaults() error {
	for _, h := range c.configHandlers {
		//谁有自己的设置默认值函数 就调用谁
		hh, ok := h.(defaultsHandler)
		if ok {
			if err := hh.SetDefaults(c); err != nil {
				return err
			}
		}
	}

	return nil
}

type structTagsConfig struct{}

func (h *structTagsConfig) SetDefaults(c *Config) error {

}

type tlsConfig struct {
	certFile string
	keyFile  string
}
