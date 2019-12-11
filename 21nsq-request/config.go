package xunsq

import (
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
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

	ClientID  string `opt:"client_id"`
	Hostname  string `opt:"hostname"`
	UserAgent string `opt:"user_agent"`
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

//设置值
func (c *Config) Set(option string, value interface{}) error {
	c.assertInitialiaed()
	option = strings.Replace(option, "-", "_", -1)
	for _, h := range c.configHandlers {
		//参数验证
		if h.HandlesOption(c, option) {
			return h.Set(c, option, value)
		}
	}

	return fmt.Errorf("invalid option %s", option)
}

func (c *Config) assertInitialiaed() {
	if !c.initialized {
		panic("Config{} must be created with NewConfig()")
	}
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

func (h *structTagsConfig) HandlersOption(c *Config, option string) bool {
	val := reflect.ValueOf(c).Elem()
	typ := val.Type()

	//如果是标签里的值就验证通过
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		opt := field.Tag.Get("opt")
		if opt == option {
			return true
		}
	}

	return false
}

func (h *structTagsConfig) SetDefaults(c *Config) error {
	val := reflect.ValueOf(c).Elem()
	typ := val.Type()

	//structTagsConfig 把值和默认值 复制给 Config
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		opt := field.Tag.Get("opt")
		defaultVal := field.Tag.Get("default")
		if defaultVal == "" || opt == "" {
			continue
		}

		if err := c.Set(opt, defaultVal); err != nil {
			return err
		}
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("error: unable to get hostname %s", err.Error())
	}

	c.ClientID = strings.Split(hostname, ".")[0]
	c.Hostname = hostname
	c.UserAgent = fmt.Sprintf("go-nsq%s", VERSION)

	return nil
}

type tlsConfig struct {
	certFile string
	keyFile  string
}

func (t *tlsConfig) HandlesOption(c *Config, option string) bool {
	switch option {
	case "tls_root_ca_file", "tls_insecure_skip", "tls_cert", "tls_key", "tls_min_version":
		return true
	}

	return false
}
