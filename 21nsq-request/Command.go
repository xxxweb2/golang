package xunsq

import "encoding/json"

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
