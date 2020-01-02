package xunsq

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

var byteNewLine = []byte("\n")
var byteSpace = []byte(" ")

type Command struct {
	Name   []byte
	Params [][]byte
	Body   []byte
}

func Publish(topic string, body []byte) *Command {

	var params = [][]byte{[]byte(topic)}
	return &Command{[]byte("PUB"), params, body}
}

func Identify(js map[string]interface{}) (*Command, error) {
	body, err := json.Marshal(js)
	if err != nil {
		return nil, err
	}

	return &Command{[]byte("IDENTIFY"), nil, body}, nil
}

func (c *Command) WriteTo(w io.Writer) (int64, error) {
	var total int64
	var buf [4]byte

	n, err := w.Write(c.Name)
	total += int64(n)
	if err != nil {
		return total, err
	}

	for _, param := range c.Params {
		n, err := w.Write(byteSpace)
		total := int64(n)
		if err != nil {
			return total, err
		}

		n, err = w.Write(param)
		total += int64(n)
		if err != nil {
			return total, err
		}
	}

	n, err = w.Write(byteNewLine)
	total += int64(n)
	if err != nil {
		return total, err
	}

	if c.Body != nil {
		bufs := buf[:]
		binary.BigEndian.PutUint32(bufs, uint32(len(c.Body)))
		n, err := w.Write(bufs)
		total += int64(n)
		if err != nil {
			return total, err
		}
		n, err = w.Write(c.Body)
		total += int64(n)
		if err != nil {
			return total, err
		}

	}

	return total, nil
}
