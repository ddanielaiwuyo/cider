package protocol

type UserId uint64
type MessageType int

const (
	ServerPaintMessage MessageType = iota
	NormalMessage
	ServerErrorResponse
)

type Request struct {
	Recipient UserId `json:"recipient"`
	Msg       string `json:"msg"`
}

type Response struct {
	From UserId      `json:"from"`
	Code MessageType `json:"code"`
	Msg  string      `json:"msg"`
}
