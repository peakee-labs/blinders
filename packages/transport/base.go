package transport

type BaseTransport struct {
	ConsumerMap ConsumerMap
}

func (bt BaseTransport) ConsumerID(key Key) string {
	return bt.ConsumerMap[key]
}
