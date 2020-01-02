package xunsq

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type IdentifyResponse struct {
	MaxRdyCount  int  `json:"max_rdy_count"`
	TLSv1        bool `json:"tls_v1"`
	Deflate      bool `json:"deflate"`
	Snappy       bool `json:"snappy"`
	AuthRequired bool `json:"auth_required"`
}

type msgResponse struct {
	msg     *Message
	cmd     *Command
	success bool
	backoff bool
}

type Conn struct {
	messageInFlight  int64
	maxRdyCount      int64
	rdyCount         int64
	lastTdyTimestamp int64
	lastMsgTimestamp int64

	mtx sync.Mutex

	config *Config

	conn    *net.TCPConn
	tlsConn *tls.Conn
	addr    string

	delegate ConnDelegate

	logGuard sync.RWMutex

	r io.Reader
	w io.Writer

	cmdChan         chan *Command
	msgResponseChan chan *msgResponse
	exitChan        chan int
	drainReady      chan int

	closeFlag int32
	stopper   sync.Once
	wg        sync.WaitGroup

	readLoopRunning int32

}

func (c *Conn) String() string {
	return ""
}

func (c *Conn) Connect() (*IdentifyResponse, error) {
	dialer := &net.Dialer{
		LocalAddr: c.config.LocalAddr,
		Timeout:   c.config.DialTimeout,
	}

	conn, err := dialer.Dial("tcp", c.addr)
	if err != nil {
		return nil, err
	}

	c.conn = conn.(*net.TCPConn)
	c.r = conn
	c.w = conn

	_, err = c.Write(MagicV2)
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("[%s] failed to write magic - %s", c.addr, err)
	}

	resp, err := c.identify()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (c *Conn) Close() error {
	return nil
}

func (c *Conn) WriteCommand(cmd *Command) error {
	c.mtx.Lock()

	_,err := cmd.WriteTo(c)
	if err != nil {
		goto exit
	}
	err = c.Flush()

exit:



	return nil
}

func (c *Conn) Write(p []byte) (int, error) {
	c.conn.SetWriteDeadline(time.Now().Add(c.config.WriteTimeout))
	return c.w.Write(p)
}

func NewConn(addr string, config *Config, delegate ConnDelegate) *Conn {
	if !config.initialized {
		panic("Config must be created with NewConfig()")
	}

	return &Conn{
		addr:     addr,
		config:   config,
		delegate: delegate,

		maxRdyCount:      2500,
		lastMsgTimestamp: time.Now().Unix(),

		cmdChan:         make(chan *Command),
		msgResponseChan: make(chan *msgResponse),
		exitChan:        make(chan int),
		drainReady:      make(chan int),
	}
}

func (c *Conn)identify()(*IdentifyResponse,error)  {
	ci := make(map[string]interface{})
	ci["client_id"] = c.config.ClientID
	ci["hostname"] = c.config.Hostname
	ci["user_agent"] = c.config.UserAgent
	ci["short_id"] = c.config.ClientID // deprecated
	ci["long_id"] = c.config.Hostname  // deprecated
	ci["tls_v1"] = c.config.TlsV1
	ci["deflate"] = c.config.Deflate
	ci["deflate_level"] = c.config.DeflateLevel
	ci["snappy"] = c.config.Snappy
	ci["feature_negotiation"] = true

	if c.config.HeartbeatInterval == -1 {
		ci["heartbeat_interval"] = -1
	}else {
		ci["heartbeat_interval"] = int64(c.config.HeartbeatInterval / time.Millisecond)
	}

	ci["sample_rate"] = c.config.SampleRate
	ci["output_buffer_size"] = c.config.OutputBufferSize
	if c.config.OutputBufferTimeout == -1 {
		ci["output_buffer_timeout"] = -1
	} else {
		ci["output_buffer_timeout"] = int64(c.config.OutputBufferTimeout / time.Millisecond)
	}
	ci["msg_timeout"] = int64(c.config.MsgTimeout / time.Millisecond)

	cmd,err := Identify(ci)
	if err != nil {
		return nil,ErrIdentify{err.Error()}
	}

	err = c.WriteCommand(cmd)
	if err != nil {
		return nil,ErrIdentify{err.Error()}
	}

	frameType,data,err := ReadUnpackedResponse(c)


}
