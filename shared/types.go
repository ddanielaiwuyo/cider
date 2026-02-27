package shared

type MessageType int

const (
	PaintMessage MessageType = iota
	ChatMessage
)

type Message struct {
	From    int         `json:"from"`
	Content string      `json:"content"`
	Dest    int         `json:"dest"`
}
