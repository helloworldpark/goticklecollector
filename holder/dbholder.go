package holder

import "github.com/helloworldpark/goticklecollector/collector"

type coinBuffer []collector.Coin

type DBHolder struct {
	buffers     []*coinBuffer
	bufferIndex int
	capacity    int
}

func NewBuffer(capacity int) DBHolder {
	buffers := make([]*coinBuffer, 2)
	b1 := make(coinBuffer, 0, capacity)
	b2 := make(coinBuffer, 0, capacity)
	buffers[0] = &b1
	buffers[1] = &b2

	dbHolder := DBHolder{buffers: buffers, bufferIndex: 0, capacity: capacity}
	return dbHolder
}
