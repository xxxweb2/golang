package xufsnotify

type Op uint32

type Event struct {
	Name string
	Op   Op
}

const (
	Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)
