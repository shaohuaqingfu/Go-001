package model

type MessageKind int32
type MessageType int32

const (
	ReplyKind MessageKind = 1

	ReplyType MessageType = 1
)

type Message struct {
	ID      int64       `json:"id"`
	Title   string      `json:"title"`
	Content string      `json:"content"`
	Kind    MessageKind `json:"kind"`
	Type    MessageType `json:"type"`
}
