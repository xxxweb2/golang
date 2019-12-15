package xunsq

const MsgIDLength = 16
type MessageID [MsgIDLength]byte


type Message struct {
	ID        MessageID
	Body      []byte
	Timestamp int64
	Attempts  uint16

	NSQAddress string

	Delegate MessageDelegate

	autoResponseDisabled int32
	responded            int32
}
