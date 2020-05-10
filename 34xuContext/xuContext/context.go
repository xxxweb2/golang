package xuContext

import (
	"errors"
	"sync"
	"time"
)

var DeadlineExceeded error = deadlineExceededError{}

type deadlineExceededError struct{}

func (deadlineExceededError) Error() string   { return "context deadline exceeded" }
func (deadlineExceededError) Timeout() bool   { return true }
func (deadlineExceededError) Temporary() bool { return true }

var Canceled = errors.New("context canceled")

type Context interface {
	Done() <-chan struct{}
	Err() error
	Deadline() (deadline time.Time, ok bool)
}

type emptyCtx int

func (*emptyCtx) Done() <-chan struct{} {
	return nil
}
func (*emptyCtx) Err() error {
	return nil
}

type CancelFunc func()

var (
	background = new(emptyCtx)
)

func Background() Context {
	return background
}

func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	c := newCancelCtx(parent)
	propagateCancel(parent, &c)

	return &c, func() {
		c.cancel(true, Canceled)
	}
}

func newCancelCtx(parent Context) cancelCtx {
	return cancelCtx{Context: parent}
}

type canceler interface {
	cancel(removeFromparent bool, err error)
	Done() <-chan struct{}
}

type cancelCtx struct {
	Context

	mu       sync.Mutex
	done     chan struct{}
	children map[canceler]struct{}
	err      error
}

var closedchan = make(chan struct{})

func (c *cancelCtx) cancel(removeFromParent bool, err error) {
	//Context关闭管道时 要透传error
	if err == nil {
		panic("context: internal error: missing cancel error")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	//如果错误不为空 说明已经取消了
	if c.err != nil {
		return
	}
	c.err = err

	if c.done == nil {
		//重新赋值一个空的管道
		c.done = closedchan
	} else {
		//关闭通道, 接受 <-ctx.Done() 处就知道要关闭这个方法了

		close(c.done)

	}

	for child := range c.children {
		//透传这个err, 关闭子上下文通道
		child.cancel(false, err)
	}
	c.children = nil

	//将自己从父children数组中删除
	if removeFromParent {
		removeChild(c.Context, c)
	}

}

func (c *cancelCtx) Done() <-chan struct{} {
	c.mu.Lock()
	if c.done == nil {
		c.done = make(chan struct{})
	}
	d := c.done
	c.mu.Unlock()

	return d
}

func removeChild(parent Context, child canceler) {
	p, ok := parentCancelCtx(parent)
	if !ok {
		return
	}
	p.mu.Lock()
	if p.children != nil {
		delete(p.children, child)
	}

	p.mu.Unlock()
}

func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	for {
		switch c := parent.(type) {
		case *cancelCtx:
			return c, true
		default:
			return nil, false
		}
	}
}

func (*emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func propagateCancel(parent Context, child canceler) {
	//父类emptyCtx 返回的都是nil, 这样父类永远不会被取消了
	if parent.Done() == nil {
		return
	}

	//查找父类是否有cancelCtx
	if p, ok := parentCancelCtx(parent); ok {
		//如果有
		p.mu.Lock()
		//err 为nil, 说明上下文的管道未关闭, 不是nil时，说明已经关闭 err未关闭的错误信息,例如时间到了的错误信息
		if p.err != nil {
			//父类管道关闭了，子类也关闭管道
			child.cancel(false, p.err)
		} else {
			//父类管道是关闭的
			if p.children == nil {
				p.children = make(map[canceler]struct{})
			}
			//将child放入到父类的children数组中
			p.children[child] = struct{}{}
		}
		p.mu.Unlock()
	} else {
		//走到这里说明父类没有cancelCtx
		go func() {
			select {
			case <-parent.Done():
				child.cancel(false, parent.Err())
			case <-child.Done():
			}
		}()
	}
}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	//调用WithDeadline， 当前时间+timeout，时间到了 done管道收到信息，函数自行关闭
	return WithDeadline(parent, time.Now().Add(timeout))
}

func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	if cur, ok := parent.Deadline(); ok && cur.Before(d) {
		//时间过了 cancelCtx
		return WithCancel(parent)
	}
	//实例timeCtx
	c := &timeCtx{
		//cancel的上下文
		cancelCtx: newCancelCtx(parent),
		deadline:  d,
	}

	//cancelCtx也调用了这个函数
	propagateCancel(parent, c)
	//时间差
	dur := time.Until(d)
	if dur <= 0 {
		c.cancel(true, DeadlineExceeded)
		return c, func() {
			c.cancel(false, Canceled)
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	//管道没有关闭
	if c.err == nil {
		//一定时间后触发关闭函数
		c.timer = time.AfterFunc(dur, func() {
			c.cancel(true, DeadlineExceeded)
		})
	}

	return c, func() {
		c.cancel(true, Canceled)
	}
}

type timeCtx struct {
	cancelCtx
	timer    *time.Timer
	deadline time.Time
}

func (c *timeCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, ok
}

func (c *timeCtx) cancel(removeFromParent bool, err error) {
	//调用cancelCtx的cancel进行关闭
	c.cancelCtx.cancel(false, err)
	if removeFromParent {
		removeChild(c.cancelCtx.Context, c)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	//关闭
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}

}
