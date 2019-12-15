package xunsq

import (
	"golang.org/x/tools/go/ssa/interp/testdata/src/fmt"
	"sync"
	"sync/atomic"
)

type producerConn interface {
	String() string
	//SetLogger(logger, LogLevel, string)
	Connect() (*IdentifyResponse, error)
	Close() error
	WriteCommand(*Command) error
}

type ProducerTransaction struct {
	cmd      *Command
	doneChan chan *ProducerTransaction
	Error    error
	Args     []interface{}
}

type Producer struct {
	id     int64
	addr   string
	conn   producerConn
	config Config

	logGuard sync.RWMutex

	responseChan chan []byte
	errorChan    chan []byte
	closeChan    chan int

	transactionChan chan *ProducerTransaction
	transactions    [] *ProducerTransaction
	state           int32

	concurrentProducers int32
	stopFlag            int32
	exitChan            chan int
	wg                  sync.WaitGroup
	guard               sync.Mutex
}

func NewProducer(addr string, config *Config) (*Producer, error) {
	config.assertInitialiaed()
	//err := config.Validate()

	p := &Producer{
		id:     atomic.AddInt64(&instCount, 1),
		addr:   addr,
		config: *config,

		transactionChan: make(chan *ProducerTransaction),
		exitChan:        make(chan int),
		responseChan:    make(chan []byte),
		errorChan:       make(chan []byte),
	}

	return p, nil
}

func (w *Producer) Publish(topic string, body []byte) error {
	return w.sendCommand(Publish(topic, body))
}

func (w *Producer) sendCommand(cmd *Command) error {
	doneChan := make(chan *ProducerTransaction)
	err := w.sendCommandAsync(cmd, doneChan, nil)
	if err != nil {
		close(doneChan)
		return err
	}

	t := <-doneChan
	return t.Error
}

func (w *Producer) sendCommandAsync(cmd *Command, doneChan chan *ProducerTransaction, args []interface{}) error {
	atomic.AddInt32(&w.concurrentProducers, 1)
	defer atomic.AddInt32(&w.concurrentProducers, -1)

	if atomic.LoadInt32(&w.state) != StateConnected {
		err := w.connect()
		if err != nil {
			return err
		}
	}

	t := &ProducerTransaction{
		cmd:      cmd,
		doneChan: doneChan,
		Args:     args,
	}

	select {
	case w.transactionChan <- t:
	case <-w.exitChan:
		return ErrStoppted
	}

	return nil
}

func (w *Producer) connect() error {
	w.guard.Lock()
	defer w.guard.Unlock()

	if atomic.LoadInt32(&w.stopFlag) == 1 {
		return ErrStoppted
	}

	switch state := atomic.LoadInt32(&w.state); state {
	case StateInit:
	case StateConnected:
		return nil
	default:
		return ErrNotConnected
	}

	fmt.Println("(%s) connecting to nsqd", w.addr)

	w.conn = NewConn(w.addr, &w.config, &producerConnDelegate{w})
	_, err := w.conn.Connect()
	if err != nil {
		w.conn.Close()
		return err
	}

	atomic.StoreInt32(&w.state, StateConnected)
	w.closeChan = make(chan int)
	w.wg.Add(1)
	go w.router()

	return nil
}
