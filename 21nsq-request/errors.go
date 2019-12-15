package xunsq

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

var ErrStoppted = errors.New("stopped")
var ErrNotConnected = errors.New("not connected")

type ErrIdentify struct {
	Reason string
}
func (e ErrIdentify) Error() string {
	return fmt.Sprintf("failed to IDENTIFY - %s", e.Reason)
}