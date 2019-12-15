package xunsq

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type BackoffStrategy interface {
	Calculate(attempt int) time.Duration
}

type configHandler interface {
	HandlesOption(c *Config, option string) bool
	Set(c *Config, option string, value interface{}) error
	Validate(c *Config) error
}

type defaultsHandler interface {
	SetDefaults(c *Config) error
}
type ExponentialStrategy struct {
	cfg *Config
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

	TlsV1     bool        `opt:"tls_v1"`

	Deflate      bool `opt:"deflate"`
	DeflateLevel int  `opt:"deflate_level" min:"1" max:"9" default:"6"`
	Snappy       bool `opt:"snappy"`

	HeartbeatInterval time.Duration `opt:"heartbeat_interval" default:"30s"`

	SampleRate int32 `opt:"sample_rate" min:"0" max:"99"`
	OutputBufferTimeout time.Duration `opt:"output_buffer_timeout" default:"250ms"`
	OutputBufferSize int64 `opt:"output_buffer_size" default:"16384"`

	MsgTimeout time.Duration `opt:"msg_timeout" min:"0"`
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

func (h *structTagsConfig) Validate(c *Config) error {
	return nil
}

func (h *structTagsConfig) Set(c *Config, option string, value interface{}) error {
	val := reflect.ValueOf(c).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		opt := field.Tag.Get("opt")

		if option != opt {
			continue
		}

		min := field.Tag.Get("min")
		max := field.Tag.Get("max")

		fieldVal := val.FieldByName(field.Name)
		dest := unsafeValueOf(fieldVal)
		coercedVal, err := coerce(value, field.Type)
		if err != nil {
			return fmt.Errorf("failed to coerce option %s (%v) - %s",
				option, value, err)
		}

		if min != "" {
			coercedMinVal, _ := coerce(min, field.Type)
			if valueCompare(coercedVal, coercedMinVal) == -1 {
				return fmt.Errorf("invalid %s ! %v < %v",
					option, coercedVal.Interface(), coercedMinVal.Interface())
			}
		}

		if max != "" {
			coercedMaxVal, _ := coerce(max, field.Type)
			if valueCompare(coercedVal, coercedMaxVal) == 1 {
				return fmt.Errorf("invalid %s ! %v > %v",
					option, coercedVal.Interface(), coercedMaxVal.Interface())
			}
		}

		dest.Set(coercedVal)
		return nil
	}

	return fmt.Errorf("unknown option %s", option)
}

func (h *structTagsConfig) HandlesOption(c *Config, option string) bool {
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

func (t *tlsConfig) Validate(c *Config) error {
	return nil
}

func (t *tlsConfig) Set(c *Config, option string, value interface{}) error {
	return nil
}
func (t *tlsConfig) HandlesOption(c *Config, option string) bool {
	switch option {
	case "tls_root_ca_file", "tls_insecure_skip", "tls_cert", "tls_key", "tls_min_version":
		return true
	}

	return false
}

func unsafeValueOf(val reflect.Value) reflect.Value {
	uptr := unsafe.Pointer(val.UnsafeAddr())

	return reflect.NewAt(val.Type(), uptr).Elem()
}

func coerce(v interface{}, typ reflect.Type) (reflect.Value, error) {
	var err error
	if typ.Kind() == reflect.Ptr {
		return reflect.ValueOf(v), nil
	}

	switch typ.String() {
	case "string":
		v, err = coerceString(v)
	case "int", "int16", "int32", "int64":
		v, err = coerceInt64(v)
	case "uint", "uint16", "uint32", "uint64":
		v, err = coerceUint64(v)
	case "float32", "float64":
		v, err = coerceFloat64(v)
	case "bool":
		v, err = coerceBool(v)
	case "time.Duration":
		v, err = coerceDuration(v)
	case "net.Addr":
		v, err = coerceAddr(v)
	default:
		v = nil
		err = fmt.Errorf("invalid type %s", typ.String())
	}
	return valueTypeCoerce(v, typ), err
}

func valueTypeCoerce(v interface{}, typ reflect.Type) reflect.Value {
	val := reflect.ValueOf(v)
	if reflect.TypeOf(v) == typ {
		return val
	}
	tval := reflect.New(typ).Elem()
	switch typ.String() {
	case "int", "int16", "int32", "int64":
		tval.SetInt(val.Int())
	case "uint", "uint16", "uint32", "uint64":
		tval.SetUint(val.Uint())
	case "float32", "float64":
		tval.SetFloat(val.Float())
	default:
		tval.Set(val)
	}
	return tval
}
func coerceAddr(v interface{}) (net.Addr, error) {
	switch v := v.(type) {
	case string:
		return net.ResolveTCPAddr("tcp", v)
	case net.Addr:
		return v, nil
	}
	return nil, errors.New("invalid value type")
}
func coerceString(v interface{}) (string, error) {
	switch v := v.(type) {
	case string:
		return v, nil
	case int, int16, int32, int64, uint, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v), nil
	case float32, float64:
		return fmt.Sprintf("%f", v), nil
	}
	return fmt.Sprintf("%s", v), nil
}

func coerceInt64(v interface{}) (int64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseInt(v, 10, 64)
	case int, int16, int32, int64:
		return reflect.ValueOf(v).Int(), nil
	case uint, uint16, uint32, uint64:
		return int64(reflect.ValueOf(v).Uint()), nil
	}
	return 0, errors.New("invalid value type")
}

func coerceUint64(v interface{}) (uint64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseUint(v, 10, 64)
	case int, int16, int32, int64:
		return uint64(reflect.ValueOf(v).Int()), nil
	case uint, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint(), nil
	}
	return 0, errors.New("invalid value type")
}

func coerceBool(v interface{}) (bool, error) {
	switch v := v.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	case int, int16, int32, int64:
		return reflect.ValueOf(v).Int() != 0, nil
	case uint, uint16, uint32, uint64:
		return reflect.ValueOf(v).Uint() != 0, nil
	}
	return false, errors.New("invalid value type")
}

func coerceFloat64(v interface{}) (float64, error) {
	switch v := v.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case int, int16, int32, int64:
		return float64(reflect.ValueOf(v).Int()), nil
	case uint, uint16, uint32, uint64:
		return float64(reflect.ValueOf(v).Uint()), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	}
	return 0, errors.New("invalid value type")
}

func coerceDuration(v interface{}) (time.Duration, error) {
	switch v := v.(type) {
	case string:
		return time.ParseDuration(v)
	case int, int16, int32, int64:
		// treat like ms
		return time.Duration(reflect.ValueOf(v).Int()) * time.Millisecond, nil
	case uint, uint16, uint32, uint64:
		// treat like ms
		return time.Duration(reflect.ValueOf(v).Uint()) * time.Millisecond, nil
	case time.Duration:
		return v, nil
	}
	return 0, errors.New("invalid value type")
}

func valueCompare(v1 reflect.Value, v2 reflect.Value) int {
	switch v1.Type().String() {
	case "int", "int16", "int32", "int64":
		if v1.Int() > v2.Int() {
			return 1
		} else if v1.Int() < v2.Int() {
			return -1
		}

		return 0

	case "uint", "uint16", "uint32", "uint64":
		if v1.Uint() > v2.Uint() {
			return 1
		} else if v1.Uint() < v2.Uint() {
			return -1
		}
		return 0
	case "float32", "float64":
		if v1.Float() > v2.Float() {
			return 1
		} else if v1.Float() < v2.Float() {
			return -1
		}
		return 0

	case "time.Duration":
		if v1.Interface().(time.Duration) > v2.Interface().(time.Duration) {
			return 1
		} else if v1.Interface().(time.Duration) < v2.Interface().(time.Duration) {
			return -1
		}
		return 0
	}

	panic("impossible")
}
