// +build windows
package xufsnotify

import (
	"sync"
	"syscall"
)

type inode struct {
	handle syscall.Handle
	volume uint32
	index  uint64
}

type watch struct {
	ov     syscall.Overlapped
	ino    *inode
	path   string
	mask   uint64
	names  map[string]uint64
	rename string
	buf    [4096]byte
}

type indexMap map[uint64]*watch
type watchMap map[uint32]indexMap

type input struct {
	op    int
	path  string
	flags uint32
	reply chan error
}

type Watcher struct {
	Events   chan Event
	Errors   chan error
	isClosed bool
	mu       sync.Mutex
	port     syscall.Handle
	watches  watchMap
	input    chan *input
	quit     chan chan<- error
}

func NewWatcher() (*Watch, error) {

}
